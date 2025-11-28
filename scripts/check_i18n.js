#!/usr/bin/env node

/**
 * i18n Configuration Analyzer and Fixer (Enhanced)
 * 
 * This script performs multiple tasks:
 * 1. Detects and removes redundant i18n keys that are not used in the codebase
 * 2. Detects hardcoded Chinese/English strings and suggests i18n keys
 * 3. Validates dictionary consistency across languages
 * 4. Suggests i18n keys for detected hardcoded strings
 * 5. Generates auto-fix suggestions for common patterns
 * 
 * Usage:
 *   node scripts/check_i18n.js [--fix] [--detect-hardcoded] [--validate] [--suggest] [--verbose]
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
const validateDicts = args.includes('--validate');
const suggestKeys = args.includes('--suggest');
const verbose = args.includes('--verbose');

// Colors for terminal output
const colors = {
    reset: '\x1b[0m',
    bright: '\x1b[1m',
    red: '\x1b[31m',
    green: '\x1b[32m',
    yellow: '\x1b[33m',
    blue: '\x1b[34m',
    magenta: '\x1b[35m',
    cyan: '\x1b[36m',
};

function colorize(text, color) {
    return `${colors[color]}${text}${colors.reset}`;
}

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
    const dynamicKeyPrefixes = new Set(); // Track dynamic key base paths

    // Pattern to match t('key') or t("key") or t(`key`)
    const directPattern = /t\(['"`]([\w.]+)['"`]\)/g;
    // Pattern to match dynamic keys like t(`prefix.${variable}`)
    const dynamicPattern = /t\(['"`]([\w.]+)\.\$\{/g;
    // Pattern to match variable assignments for dynamic keys: const x = 'keyName'
    const variablePattern = /const\s+(\w+)\s*=.*?['"]([\w]+)['"];?/g;

    // Get all TypeScript/JavaScript/TSX/JSX files
    const findCmd = `find "${srcDir}" -type f \\( -name "*.ts" -o -name "*.tsx" -o -name "*.js" -o -name "*.jsx" \\) ${CONFIG.excludePatterns.map(p => `! -path "*/${p}/*"`).join(' ')}`;

    try {
        const files = execSync(findCmd, { encoding: 'utf8' })
            .split('\n')
            .filter(Boolean);

        for (const file of files) {
            const content = fs.readFileSync(file, 'utf8');
            let match;

            // Find direct key usage
            while ((match = directPattern.exec(content)) !== null) {
                usedKeys.add(match[1]);
            }

            // Find dynamic key patterns and analyze them
            directPattern.lastIndex = 0;
            while ((match = dynamicPattern.exec(content)) !== null) {
                const basePath = match[1];
                usedKeys.add(basePath);
                dynamicKeyPrefixes.add(basePath);
                
                // Add common dynamic suffixes (known patterns)
                ['positive', 'neutral', 'negative'].forEach(s => usedKeys.add(`${basePath}.${s}`));
                ['high', 'medium', 'low'].forEach(s => usedKeys.add(`${basePath}.${s}`));
                ['professional', 'casual', 'concise', 'detailed'].forEach(s => usedKeys.add(`${basePath}.${s}`));
            }
            
            // Analyze dynamic key construction patterns
            // Look for patterns like: const timeOfDay = ... 'greetingMorning' : ... 'greetingAfternoon'
            dynamicPattern.lastIndex = 0;
            for (const dynamicBase of dynamicKeyPrefixes) {
                // Extract the last segment of the dynamic base (e.g., 'dashboard' from 'dashboard')
                const baseSegments = dynamicBase.split('.');
                const lastSegment = baseSegments[baseSegments.length - 1];
                
                // Search for variable assignments that might be used with this base
                const contextRegex = new RegExp('[\'"]([\\w]+(?:Morning|Afternoon|Evening|Day|Night|Start|End|Begin|Finish|Success|Error|Warning|Info))[\'"]', 'gi');
                let contextMatch;
                while ((contextMatch = contextRegex.exec(content)) !== null) {
                    const potentialKey = contextMatch[1];
                    // Check if this appears in a conditional expression near the dynamic key usage
                    const escapedBase = dynamicBase.replace(/\./g, '\\.');
                    const keyPattern = new RegExp('t\\([\'"`]' + escapedBase + '\\.\\$\\{', 'g');
                    if (keyPattern.test(content)) {
                        usedKeys.add(dynamicBase + '.' + potentialKey);
                    }
                }
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
    // 更精准的英文 UI 模式 - 只检测真正的 UI 文本
    const englishUIPattern = /["'`]((Loading|Error|Success|Failed|Submit|Cancel|Delete|Edit|Save|Search|Filter|Settings|Dashboard|Profile|Logout|Login|Sign up|Sign in|Welcome|Hello|Goodbye|Confirm|Back|Next|Previous|Close|Open|New|Create|Update|Remove|Add|View|Show|Hide|Clear|Reset|Apply|Download|Upload|Export|Import|Print|Share|Copy|Paste|Cut|Refresh|Reload)[^"'`]{0,30})["'`]/gi;

    // 需要排除的模式
    const excludePatterns = [
        // API 方法和 HTTP 动词
        /^(GET|POST|PUT|DELETE|PATCH|HEAD|OPTIONS)$/i,
        // 导入路径
        /^[a-z@][a-z0-9/-]*$/,
        // 状态常量和类型
        /^(success|error|warning|info|new|active|pending|completed|failed|idle|loading|search|chat)$/i,
        // 技术常量
        /^(true|false|null|undefined)$/i,
        // CSS 类名或标识符
        /^[a-z][a-z0-9-_]*$/,
        // 文件扩展名或 MIME 类型
        /^\.(json|js|ts|tsx|jsx|css|html|xml|txt)$/,
        /^(application|text|image|video)\//,
        // 数据库字段或 API 键名
        /^[a-z_][a-z0-9_]*$/,
        // 单个技术词汇
        /^(id|uuid|url|uri|email|password|token|key|data|type|name|value|status|code)$/i,
    ];

    const findCmd = `find "${srcDir}" -type f \\( -name "*.ts" -o -name "*.tsx" -o -name "*.js" -o -name "*.jsx" \\) ${CONFIG.excludePatterns.map(p => `! -path "*/${p}/*"`).join(' ')}`;

    try {
        const files = execSync(findCmd, { encoding: 'utf8' })
            .split('\n')
            .filter(Boolean);

        for (const file of files) {
            const content = fs.readFileSync(file, 'utf8');
            const lines = content.split('\n');

            // Skip dictionary files themselves and type definition files
            if (file.includes('/i18n/dictionaries/') || file.endsWith('.d.ts')) continue;

            lines.forEach((line, lineNum) => {
                // Skip comment lines
                const trimmedLine = line.trim();
                if (trimmedLine.startsWith('//') || trimmedLine.startsWith('/*') || trimmedLine.startsWith('*')) {
                    return;
                }

                // Check for Chinese strings
                let match;
                while ((match = chinesePattern.exec(line)) !== null) {
                    // Skip if it's already in a t() call
                    const beforeMatch = line.substring(0, match.index);
                    if (!beforeMatch.match(/t\(['"`]*$/)) {
                        hardcodedFindings.push({
                            file: path.relative(srcDir, file),
                            line: lineNum + 1,
                            text: match[1],
                            type: 'chinese',
                            context: line.trim().substring(0, 100)
                        });
                    }
                }

                // Check for English UI strings
                chinesePattern.lastIndex = 0;
                while ((match = englishUIPattern.exec(line)) !== null) {
                    const text = match[1];
                    const beforeMatch = line.substring(0, match.index);
                    
                    // Skip if already in t() call
                    if (beforeMatch.match(/t\(['"`]*$/)) continue;
                    
                    // Skip if matches any exclude pattern
                    if (excludePatterns.some(pattern => pattern.test(text))) continue;
                    
                    // Skip if it's in an import statement
                    if (line.includes('import ') && line.includes('from')) continue;
                    
                    // Skip if it's a property key (key: 'value' or {key: 'value'})
                    if (beforeMatch.match(/[{,]\s*\w+\s*:\s*['"`]*$/)) continue;
                    
                    // Skip if it's in console.log or console.error
                    if (line.match(/console\.(log|error|warn|info|debug)/)) continue;
                    
                    // Skip if word count is too high (likely a sentence)
                    if (text.split(' ').length > 6) continue;
                    
                    hardcodedFindings.push({
                        file: path.relative(srcDir, file),
                        line: lineNum + 1,
                        text: text,
                        type: 'english',
                        context: line.trim().substring(0, 100)
                    });
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
 * Validate dictionary consistency across languages
 */
function validateDictionaries(dictionaries) {
    const issues = [];
    const flatDicts = {};
    
    for (const lang of CONFIG.languages) {
        flatDicts[lang] = flattenKeys(dictionaries[lang]);
    }
    
    const allKeys = new Set();
    Object.values(flatDicts).forEach(dict => {
        Object.keys(dict).forEach(key => allKeys.add(key));
    });
    
    // Check for missing translations
    for (const key of allKeys) {
        const missingIn = [];
        for (const lang of CONFIG.languages) {
            if (!(key in flatDicts[lang])) {
                missingIn.push(lang);
            }
        }
        if (missingIn.length > 0) {
            issues.push({
                type: 'missing',
                key: key,
                languages: missingIn
            });
        }
    }
    
    // Check for empty values
    for (const lang of CONFIG.languages) {
        for (const [key, value] of Object.entries(flatDicts[lang])) {
            if (!value || (typeof value === 'string' && value.trim() === '')) {
                issues.push({
                    type: 'empty',
                    key: key,
                    language: lang
                });
            }
        }
    }
    
    return issues;
}

/**
 * Suggest i18n keys for hardcoded strings
 */
function suggestI18nKey(text, context, existingKeys) {
    // Generate key from text
    let baseKey = text
        .toLowerCase()
        .replace(/[^a-z0-9\u4e00-\u9fa5\s]/g, '')
        .trim()
        .split(/\s+/)
        .slice(0, 3)
        .join('');
    
    // Determine module from context
    let module = 'common';
    if (context.file.includes('/components/')) {
        const match = context.file.match(/\/components\/([^\/]+)/);
        if (match) module = match[1];
    } else if (context.file.includes('/app/')) {
        const match = context.file.match(/\/app\/([^\/]+)/);
        if (match) module = match[1];
    }
    
    // Build suggested key
    let suggestedKey = `${module}.${baseKey}`;
    let counter = 1;
    
    // Ensure uniqueness
    while (existingKeys.has(suggestedKey)) {
        suggestedKey = `${module}.${baseKey}${counter}`;
        counter++;
    }
    
    return suggestedKey;
}

/**
 * Group hardcoded findings by file and priority
 */
function prioritizeFindings(findings) {
    const prioritized = {
        high: [],    // UI components with Chinese
        medium: [],  // UI components with English
        low: []      // Backend/utility files
    };
    
    for (const finding of findings) {
        if (finding.type === 'chinese') {
            if (finding.file.includes('/components/') || finding.file.includes('/app/')) {
                prioritized.high.push(finding);
            } else {
                prioritized.medium.push(finding);
            }
        } else {
            if (finding.file.includes('/components/') || finding.file.includes('/app/')) {
                prioritized.medium.push(finding);
            } else {
                prioritized.low.push(finding);
            }
        }
    }
    
    return prioritized;
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
    console.log(colorize('\n🔍 i18n Configuration Analyzer (Enhanced)\n', 'cyan'));

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

    console.log(colorize('📚 Loaded dictionaries:', 'blue'));
    for (const lang of CONFIG.languages) {
        console.log(`   ${lang.toUpperCase()}: ${colorize(Object.keys(flatDictionaries[lang]).length, 'green')} keys`);
    }
    console.log();

    // Validate dictionaries if requested
    if (validateDicts) {
        console.log(colorize('🔍 Validating dictionary consistency...\n', 'cyan'));
        const validationIssues = validateDictionaries(dictionaries);
        
        if (validationIssues.length === 0) {
            console.log(colorize('✅ All dictionaries are consistent!\n', 'green'));
        } else {
            console.log(colorize(`⚠️  Found ${validationIssues.length} validation issues:\n`, 'yellow'));
            
            const missingKeys = validationIssues.filter(i => i.type === 'missing');
            const emptyValues = validationIssues.filter(i => i.type === 'empty');
            
            if (missingKeys.length > 0) {
                console.log(colorize('Missing translations:', 'yellow'));
                missingKeys.forEach(issue => {
                    console.log(`   ${colorize(issue.key, 'red')} - missing in: ${issue.languages.join(', ')}`);
                });
                console.log();
            }
            
            if (emptyValues.length > 0) {
                console.log(colorize('Empty values:', 'yellow'));
                emptyValues.forEach(issue => {
                    console.log(`   ${colorize(issue.key, 'red')} - empty in: ${issue.language}`);
                });
                console.log();
            }
        }
    }

    // Find used keys
    console.log(colorize('🔎 Scanning source files for i18n key usage...', 'cyan'));
    const usedKeys = findUsedKeys(CONFIG.frontendSrcDir);
    console.log(`   Found ${colorize(usedKeys.size, 'green')} unique keys in use\n`);

    // Find redundant keys
    const redundantKeys = {};
    for (const lang of CONFIG.languages) {
        const allKeys = Object.keys(flatDictionaries[lang]);
        redundantKeys[lang] = allKeys.filter(key => !usedKeys.has(key));
    }

    // Report redundant keys
    console.log(colorize('📊 Redundant Keys Report:\n', 'cyan'));
    for (const lang of CONFIG.languages) {
        console.log(colorize(`${lang.toUpperCase()}:`, 'blue'));
        if (redundantKeys[lang].length === 0) {
            console.log(colorize('   ✅ No redundant keys found!', 'green'));
        } else {
            console.log(colorize(`   ⚠️  Found ${redundantKeys[lang].length} redundant keys:`, 'yellow'));
            if (verbose) {
                redundantKeys[lang].forEach(key => {
                    console.log(`      ${colorize('-', 'yellow')} ${key}`);
                });
            } else {
                redundantKeys[lang].slice(0, 5).forEach(key => {
                    console.log(`      ${colorize('-', 'yellow')} ${key}`);
                });
                if (redundantKeys[lang].length > 5) {
                    console.log(`      ${colorize(`... and ${redundantKeys[lang].length - 5} more`, 'yellow')}`);
                }
            }
        }
        console.log();
    }

    // Fix redundant keys if requested
    if (shouldFix && Object.values(redundantKeys).some(arr => arr.length > 0)) {
        console.log(colorize('🔧 Fixing redundant keys...\n', 'cyan'));
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
                console.log(colorize(`✅ Removed ${removedCount} redundant keys from ${lang}.json`, 'green'));
            }
        }
        console.log();
    }

    // Detect hardcoded strings if requested
    if (detectHardcoded) {
        console.log(colorize('🔍 Detecting hardcoded strings...\n', 'cyan'));
        const findings = detectHardcodedStrings(CONFIG.frontendSrcDir);

        if (findings.length === 0) {
            console.log(colorize('✅ No hardcoded strings detected!\n', 'green'));
        } else {
            console.log(colorize(`⚠️  Found ${findings.length} potential hardcoded strings\n`, 'yellow'));

            const prioritized = prioritizeFindings(findings);
            const allExistingKeys = new Set(Object.keys(flatDictionaries.en));
            
            // High priority findings
            if (prioritized.high.length > 0) {
                console.log(colorize(`🔴 HIGH PRIORITY (${prioritized.high.length} items) - Chinese in UI components:\n`, 'red'));
                const grouped = prioritized.high.reduce((acc, f) => {
                    if (!acc[f.file]) acc[f.file] = [];
                    acc[f.file].push(f);
                    return acc;
                }, {});

                for (const [file, items] of Object.entries(grouped)) {
                    console.log(colorize(`📄 ${file}:`, 'yellow'));
                    items.forEach(item => {
                        console.log(`   Line ${item.line}: ${colorize(`"${item.text}"`, 'red')}`);
                        if (suggestKeys) {
                            const suggested = suggestI18nKey(item.text, { file: item.file }, allExistingKeys);
                            console.log(`   ${colorize('💡 Suggested key:', 'cyan')} ${suggested}`);
                            allExistingKeys.add(suggested);
                        }
                        if (verbose) {
                            console.log(`   Context: ${item.context}`);
                        }
                        console.log();
                    });
                }
            }
            
            // Medium priority findings
            if (prioritized.medium.length > 0 && verbose) {
                console.log(colorize(`🟡 MEDIUM PRIORITY (${prioritized.medium.length} items) - English in UI components:\n`, 'yellow'));
                const grouped = prioritized.medium.slice(0, 10).reduce((acc, f) => {
                    if (!acc[f.file]) acc[f.file] = [];
                    acc[f.file].push(f);
                    return acc;
                }, {});

                for (const [file, items] of Object.entries(grouped)) {
                    console.log(colorize(`📄 ${file}:`, 'yellow'));
                    items.forEach(item => {
                        console.log(`   Line ${item.line}: ${colorize(`"${item.text}"`, 'yellow')}`);
                        if (suggestKeys) {
                            const suggested = suggestI18nKey(item.text, { file: item.file }, allExistingKeys);
                            console.log(`   ${colorize('💡 Suggested key:', 'cyan')} ${suggested}`);
                            allExistingKeys.add(suggested);
                        }
                        console.log();
                    });
                }
                if (prioritized.medium.length > 10) {
                    console.log(colorize(`   ... and ${prioritized.medium.length - 10} more\n`, 'yellow'));
                }
            }
            
            if (!verbose) {
                console.log(colorize(`\n💡 Use --verbose to see all findings\n`, 'cyan'));
            }
            
            console.log(colorize(`💡 Suggestion: Add these strings to your i18n dictionaries and replace with t() calls.`, 'cyan'));
        }
    }

    // Summary
    console.log('\n' + colorize('='.repeat(70), 'cyan'));
    console.log(colorize('📊 Summary:', 'bright'));
    console.log(`  ${colorize('Total keys in dictionaries:', 'blue')} ${Object.keys(flatDictionaries.en).length}`);
    console.log(`  ${colorize('Keys in use:', 'green')} ${usedKeys.size}`);
    console.log(`  ${colorize('Redundant keys:', 'yellow')} ${redundantKeys.en.length}`);
    if (validateDicts) {
        const validationIssues = validateDictionaries(dictionaries);
        const missingCount = validationIssues.filter(i => i.type === 'missing').length;
        const emptyCount = validationIssues.filter(i => i.type === 'empty').length;
        console.log(`  ${colorize('Missing translations:', 'yellow')} ${missingCount}`);
        console.log(`  ${colorize('Empty values:', 'yellow')} ${emptyCount}`);
    }
    if (detectHardcoded) {
        const findings = detectHardcodedStrings(CONFIG.frontendSrcDir);
        const prioritized = prioritizeFindings(findings);
        console.log(`  ${colorize('Hardcoded strings (High):', 'red')} ${prioritized.high.length}`);
        console.log(`  ${colorize('Hardcoded strings (Medium):', 'yellow')} ${prioritized.medium.length}`);
        console.log(`  ${colorize('Hardcoded strings (Low):', 'blue')} ${prioritized.low.length}`);
    }
    console.log(colorize('='.repeat(70), 'cyan') + '\n');

    // Tips
    const tips = [];
    if (!shouldFix && redundantKeys.en.length > 0) {
        tips.push('Run with --fix to automatically remove redundant keys (backups will be created)');
    }
    if (!detectHardcoded) {
        tips.push('Run with --detect-hardcoded to scan for hardcoded Chinese/English strings');
    }
    if (!validateDicts) {
        tips.push('Run with --validate to check dictionary consistency across languages');
    }
    if (!suggestKeys && detectHardcoded) {
        tips.push('Run with --suggest to get i18n key suggestions for hardcoded strings');
    }
    if (!verbose && detectHardcoded) {
        tips.push('Run with --verbose to see detailed context for all findings');
    }
    
    if (tips.length > 0) {
        console.log(colorize('💡 Tips:', 'cyan'));
        tips.forEach(tip => console.log(`   • ${tip}`));
        console.log();
    }
}

// Run main function
main();
