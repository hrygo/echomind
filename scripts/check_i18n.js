#!/usr/bin/env node

/**
 * i18n Configuration Analyzer and Fixer
 * 
 * This script performs two main tasks:
 * 1. Detects and removes redundant i18n keys that are not used in the codebase
 * 2. Detects hardcoded Chinese/English strings and suggests i18n keys
 * 
 * Usage:
 *   node scripts/check_i18n.js [--fix] [--detect-hardcoded]
 */

const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

// Configuration
const CONFIG = {
    dictionariesDir: path.join(__dirname, '../frontend/src/lib/i18n/dictionaries'),
    frontendSrcDir: path.join(__dirname, '../frontend/src'),
    backupDir: path.join(__dirname, '../frontend/src/lib/i18n/backups'),
    languages: ['en', 'zh'],
    excludePatterns: [
        'node_modules',
        '.next',
        'dist',
        'build',
        '.git',
        '*.test.*',
        '*.spec.*'
    ]
};

// Parse command line arguments
const args = process.argv.slice(2);
const shouldFix = args.includes('--fix');
const detectHardcoded = args.includes('--detect-hardcoded');

/**
 * Load JSON file
 */
function loadJSON(filePath) {
    try {
        return JSON.parse(fs.readFileSync(filePath, 'utf8'));
    } catch (error) {
        console.error(`Error loading ${filePath}:`, error.message);
        process.exit(1);
    }
}

/**
 * Save JSON file with formatting
 */
function saveJSON(filePath, data) {
    fs.writeFileSync(filePath, JSON.stringify(data, null, 4) + '\n', 'utf8');
}

/**
 * Flatten nested JSON keys into dot-notation paths
 * Example: { common: { appName: "EchoMind" } } => { "common.appName": "EchoMind" }
 */
function flattenKeys(obj, prefix = '') {
    const result = {};

    for (const [key, value] of Object.entries(obj)) {
        const fullKey = prefix ? `${prefix}.${key}` : key;

        if (value && typeof value === 'object' && !Array.isArray(value)) {
            Object.assign(result, flattenKeys(value, fullKey));
        } else {
            result[fullKey] = value;
        }
    }

    return result;
}

/**
 * Scan source files for i18n key usage
 */
function findUsedKeys(srcDir) {
    const usedKeys = new Set();

    // Pattern to match t('key') or t("key")
    const pattern = /t\(['"]([\w.]+)['"]\)/g;

    // Get all TypeScript/JavaScript/TSX/JSX files
    const findCmd = `find "${srcDir}" -type f \\( -name "*.ts" -o -name "*.tsx" -o -name "*.js" -o -name "*.jsx" \\) ${CONFIG.excludePatterns.map(p => `! -path "*/${p}/*"`).join(' ')}`;

    try {
        const files = execSync(findCmd, { encoding: 'utf8' })
            .split('\n')
            .filter(Boolean);

        for (const file of files) {
            const content = fs.readFileSync(file, 'utf8');
            let match;

            while ((match = pattern.exec(content)) !== null) {
                usedKeys.add(match[1]);
            }
        }
    } catch (error) {
        console.error('Error scanning files:', error.message);
    }

    return usedKeys;
}

/**
 * Detect hardcoded Chinese and English strings in source files
 */
function detectHardcodedStrings(srcDir) {
    const hardcodedFindings = [];

    // Patterns for Chinese characters and common English UI strings
    const chinesePattern = /["'`]([\u4e00-\u9fa5]+[\u4e00-\u9fa5\s\w]*?)["'`]/g;
    const englishUIPattern = /["'`]((?:Loading|Error|Success|Failed|Submit|Cancel|Delete|Edit|Save|Search|Filter|Settings|Dashboard|Profile|Logout|Login|Sign up|Sign in|Welcome|Hello|Goodbye|Yes|No|OK|Confirm|Back|Next|Previous|Close|Open|New|Create|Update|Remove|Add|View|Show|Hide)[^"'`]{0,50})["'`]/gi;

    const findCmd = `find "${srcDir}" -type f \\( -name "*.ts" -o -name "*.tsx" -o -name "*.js" -o -name "*.jsx" \\) ${CONFIG.excludePatterns.map(p => `! -path "*/${p}/*"`).join(' ')}`;

    try {
        const files = execSync(findCmd, { encoding: 'utf8' })
            .split('\n')
            .filter(Boolean);

        for (const file of files) {
            const content = fs.readFileSync(file, 'utf8');
            const lines = content.split('\n');

            // Skip dictionary files themselves
            if (file.includes('/i18n/dictionaries/')) continue;

            lines.forEach((line, lineNum) => {
                // Check for Chinese strings
                let match;
                while ((match = chinesePattern.exec(line)) !== null) {
                    // Skip if it's already in a t() call
                    const beforeMatch = line.substring(0, match.index);
                    if (!beforeMatch.match(/t\(['"]*$/)) {
                        hardcodedFindings.push({
                            file: path.relative(srcDir, file),
                            line: lineNum + 1,
                            text: match[1],
                            type: 'chinese',
                            context: line.trim()
                        });
                    }
                }

                // Check for English UI strings
                chinesePattern.lastIndex = 0;
                while ((match = englishUIPattern.exec(line)) !== null) {
                    const beforeMatch = line.substring(0, match.index);
                    // Skip if already in t() or if it looks like a code identifier
                    if (!beforeMatch.match(/t\(['"]*$/) && !match[1].includes('_') && match[1].split(' ').length <= 5) {
                        hardcodedFindings.push({
                            file: path.relative(srcDir, file),
                            line: lineNum + 1,
                            text: match[1],
                            type: 'english',
                            context: line.trim()
                        });
                    }
                }
            });
        }
    } catch (error) {
        console.error('Error detecting hardcoded strings:', error.message);
    }

    return hardcodedFindings;
}

/**
 * Remove keys from nested object based on dot-notation path
 */
function removeKey(obj, keyPath) {
    const keys = keyPath.split('.');
    const lastKey = keys.pop();

    let current = obj;
    for (const key of keys) {
        if (!current[key]) return false;
        current = current[key];
    }

    if (current[lastKey]) {
        delete current[lastKey];
        return true;
    }
    return false;
}

/**
 * Create backup of dictionary files
 */
function createBackup() {
    if (!fs.existsSync(CONFIG.backupDir)) {
        fs.mkdirSync(CONFIG.backupDir, { recursive: true });
    }

    const timestamp = new Date().toISOString().replace(/[:.]/g, '-').split('T')[0];

    for (const lang of CONFIG.languages) {
        const srcFile = path.join(CONFIG.dictionariesDir, `${lang}.json`);
        const backupFile = path.join(CONFIG.backupDir, `${lang}.${timestamp}.json`);
        fs.copyFileSync(srcFile, backupFile);
        console.log(`✅ Backed up ${lang}.json to ${backupFile}`);
    }
}

/**
 * Main function
 */
function main() {
    console.log('🔍 i18n Configuration Analyzer\n');

    // Load dictionaries
    const dictionaries = {};
    for (const lang of CONFIG.languages) {
        const filePath = path.join(CONFIG.dictionariesDir, `${lang}.json`);
        dictionaries[lang] = loadJSON(filePath);
    }

    // Flatten keys
    const flatDictionaries = {};
    for (const lang of CONFIG.languages) {
        flatDictionaries[lang] = flattenKeys(dictionaries[lang]);
    }

    console.log(`📚 Loaded dictionaries:`);
    for (const lang of CONFIG.languages) {
        console.log(`   ${lang}: ${Object.keys(flatDictionaries[lang]).length} keys`);
    }
    console.log();

    // Find used keys
    console.log('🔎 Scanning source files for i18n key usage...');
    const usedKeys = findUsedKeys(CONFIG.frontendSrcDir);
    console.log(`   Found ${usedKeys.size} unique keys in use\n`);

    // Find redundant keys
    const redundantKeys = {};
    for (const lang of CONFIG.languages) {
        const allKeys = Object.keys(flatDictionaries[lang]);
        redundantKeys[lang] = allKeys.filter(key => !usedKeys.has(key));
    }

    // Report redundant keys
    console.log('📊 Redundant Keys Report:\n');
    for (const lang of CONFIG.languages) {
        console.log(`${lang.toUpperCase()}:`);
        if (redundantKeys[lang].length === 0) {
            console.log('   ✅ No redundant keys found!');
        } else {
            console.log(`   ⚠️  Found ${redundantKeys[lang].length} redundant keys:`);
            redundantKeys[lang].forEach(key => {
                console.log(`      - ${key}`);
            });
        }
        console.log();
    }

    // Fix redundant keys if requested
    if (shouldFix && Object.values(redundantKeys).some(arr => arr.length > 0)) {
        console.log('🔧 Fixing redundant keys...\n');
        createBackup();

        for (const lang of CONFIG.languages) {
            if (redundantKeys[lang].length > 0) {
                const dictCopy = JSON.parse(JSON.stringify(dictionaries[lang]));
                let removedCount = 0;

                for (const key of redundantKeys[lang]) {
                    if (removeKey(dictCopy, key)) {
                        removedCount++;
                    }
                }

                const filePath = path.join(CONFIG.dictionariesDir, `${lang}.json`);
                saveJSON(filePath, dictCopy);
                console.log(`✅ Removed ${removedCount} redundant keys from ${lang}.json`);
            }
        }
        console.log();
    }

    // Detect hardcoded strings if requested
    if (detectHardcoded) {
        console.log('🔍 Detecting hardcoded strings...\n');
        const findings = detectHardcodedStrings(CONFIG.frontendSrcDir);

        if (findings.length === 0) {
            console.log('✅ No hardcoded strings detected!\n');
        } else {
            console.log(`⚠️  Found ${findings.length} potential hardcoded strings:\n`);

            const grouped = findings.reduce((acc, f) => {
                if (!acc[f.file]) acc[f.file] = [];
                acc[f.file].push(f);
                return acc;
            }, {});

            for (const [file, items] of Object.entries(grouped)) {
                console.log(`📄 ${file}:`);
                items.forEach(item => {
                    console.log(`   Line ${item.line} (${item.type}): "${item.text}"`);
                    console.log(`   Context: ${item.context.substring(0, 80)}${item.context.length > 80 ? '...' : ''}`);
                    console.log();
                });
            }

            console.log(`💡 Suggestion: Add these strings to your i18n dictionaries and replace with t() calls.`);
        }
    }

    // Summary
    console.log('\n' + '='.repeat(60));
    console.log('Summary:');
    console.log(`  Total keys in dictionaries: ${Object.keys(flatDictionaries.en).length}`);
    console.log(`  Keys in use: ${usedKeys.size}`);
    console.log(`  Redundant keys: ${redundantKeys.en.length}`);
    if (detectHardcoded) {
        const findings = detectHardcodedStrings(CONFIG.frontendSrcDir);
        console.log(`  Hardcoded strings detected: ${findings.length}`);
    }
    console.log('='.repeat(60) + '\n');

    if (!shouldFix && redundantKeys.en.length > 0) {
        console.log('💡 Run with --fix to automatically remove redundant keys (backups will be created)');
    }

    if (!detectHardcoded) {
        console.log('💡 Run with --detect-hardcoded to scan for hardcoded Chinese/English strings');
    }
}

// Run main function
main();
