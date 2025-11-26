# Enhanced CI/CD Monitor for EchoMind

## 🚀 Overview

The `check_ci_enhanced.sh` script provides a comprehensive CI/CD monitoring and analysis experience for the EchoMind project, with enhanced interactivity, smart error detection, and actionable insights.

## 📋 Features

### 🎯 Core Capabilities

- **Real-time Status Monitoring**: Check latest CI/CD run status with progress indicators
- **Smart Error Analysis**: Automatically detect and categorize common failure patterns
- **Interactive Mode**: Full-featured menu system for easy navigation
- **Historical Analysis**: Track success rates and identify recurring issues
- **Fix Suggestions**: Get actionable recommendations based on error patterns

### 🛠️ Advanced Functions

- **Live Watching**: Monitor running workflows in real-time
- **Deep Analysis**: Comprehensive CI/CD health assessment
- **Failed Job Rerun**: Quick rerun capabilities for failed workflows
- **Error Pattern Detection**: TypeScript, test, build, and dependency issues
- **Success Rate Analytics**: Track CI/CD performance over time

## 📖 Usage

### Basic Commands

```bash
# Show current CI/CD status
./scripts/check_ci_enhanced.sh

# Show last 5 runs (default)
./scripts/check_ci_enhanced.sh --history

# Show last 10 runs
./scripts/check_ci_enhanced.sh --history 10

# Watch live running workflow
./scripts/check_ci_enhanced.sh --watch

# Deep analysis of recent runs
./scripts/check_ci_enhanced.sh --analyze

# Interactive menu mode
./scripts/check_ci_enhanced.sh --interactive

# Show help
./scripts/check_ci_enhanced.sh --help
```

### Interactive Menu Options

When using `--interactive` mode, you get access to:

1. **Current Status** - Quick overview of latest run
2. **Watch Live** - Real-time log viewing
3. **History** - Browse recent CI/CD runs
4. **Deep Analysis** - Comprehensive health check
5. **Rerun Failed** - Quick retry for failed jobs
6. **Exit** - Leave interactive mode

## 🔧 Smart Error Detection

The script automatically detects and categorizes:

### TypeScript Errors
```bash
# Detects patterns like:
error TS2307: Cannot find module 'next-themes/dist/types'
```

### Test Failures
```bash
# Detects patterns like:
FAIL src/components/Widget.test.tsx
Test Suites: 1 failed, 5 passed
```

### Build/Dependency Issues
```bash
# Detects patterns like:
npm ERR! code ERESOLVE
ELIFECYCLE Command failed with exit code 1
```

### Permission/Docker Issues
```bash
# Detects patterns like:
permission denied: docker
Error response from daemon
```

## 💡 Fix Suggestions

Based on detected errors, the script provides targeted suggestions:

### For TypeScript Issues
- Check import paths in theme components
- Verify type declarations
- Update type dependencies

### For Build Issues
- Run `npm install` or `pnpm install`
- Check package.json dependencies
- Clear node_modules and reinstall

### For Test Failures
- Run tests locally: `make test`
- Check test configuration
- Update test expectations

### For Permission Issues
- Check file permissions
- Verify script execution rights
- Review GitHub Actions permissions

## 📊 Analytics Features

### Success Rate Tracking
```bash
📊 Statistics:
  Total Runs: 10
  Successful: 8
  Failed: 2
  Success Rate: 80%
```

### Failure Pattern Analysis
```bash
🔍 Failure Patterns:
  • Run #19673553783:
    → TypeScript errors
  • Run #19673553782:
    → Test failures
```

### Performance Recommendations
- Low success rate alerts
- Pre-commit hook suggestions
- Parallel execution tips
- Build stability recommendations

## 🎨 Visual Features

### Progress Indicators
```bash
[⏳] ████████████████████████████████████░░░░ 80%
```

### Status Emojis
- ✅ Success
- ❌ Failure
- ⚠️ Warning
- ℹ️ Info
- ⏳ In Progress
- 🚫 Cancelled

### Color Coding
- 🟢 Green: Success/Positive actions
- 🔴 Red: Failures/Errors
- 🔵 Blue: Information/Links
- 🟡 Yellow: Warnings/Progress
- ⚪ Gray: Metadata/Secondary info

## 🛠️ Requirements

### Prerequisites
- GitHub CLI (`gh`)
- jq (JSON processor)

### Installation
```bash
# Install GitHub CLI
brew install gh

# Install jq
brew install jq

# Login to GitHub
gh auth login

# Make script executable
chmod +x scripts/check_ci_enhanced.sh
```

## 🔄 Comparison with Original Script

| Feature | Original (`check_ci.sh`) | Enhanced (`check_ci_enhanced.sh`) |
|---------|--------------------------|-----------------------------------|
| Basic Status | ✅ | ✅ |
| Error Detection | ❌ | ✅ |
| Fix Suggestions | ❌ | ✅ |
| Interactive Mode | ❌ | ✅ |
| Progress Indicators | ❌ | ✅ |
| Historical Analysis | ❌ | ✅ |
| Success Rate Tracking | ❌ | ✅ |
| Color Output | ✅ | ✅ |
| Live Watching | ✅ | ✅ |
| Deep Analysis | ❌ | ✅ |

## 🚀 Examples

### Quick Status Check
```bash
$ ./scripts/check_ci_enhanced.sh

🔍 EchoMind CI/CD Enhanced Monitor
========================================
⏳ Analyzing latest CI/CD run...
[⏳]██████████████████████████████████████████████████ 100%

✅ CI/CD Run #19673553784
Workflow: CI/CD
Branch: main
Run ID: 19673553784
Created: 2025-11-26T11:15:00Z
Status: SUCCESS
Link: https://github.com/hrygo/echomind/actions/runs/19673553784
```

### Error Analysis Example
```bash
❌ CI/CD Run #19673553783
Workflow: CI/CD
Branch: main
Run ID: 19673553783
Created: 2025-11-26T11:10:00Z
Status: FAILURE

⚠️ Failure Analysis:
Failed Jobs:
  • Quality Assurance ❌
    Failed steps: Frontend Lint & Test

ℹ️ Error Analysis:
Key Error Patterns:
  ❌• TypeScript compilation errors
    src/components/theme/ThemeProviderNext.tsx(4,41): error TS2307: Cannot find module 'next-themes/dist/types'

💡 Suggested Fixes:
  ✅• TypeScript import issues detected
    Fix: Check import paths in theme components
  ✅• Check the full logs: gh run view 19673553783 --log-failed
  ✅• Re-run failed jobs: gh run rerun 19673553783
  ✅• Run locally: make test && make build
```

### History View Example
```bash
📊 CI/CD History (Last 5 runs):
----------------------------------------
#19673553784 ✅ CI/CD           main       2025-11-26T11:15
#19673553783 ❌ CI/CD           main       2025-11-26T11:10
#19673553782 ✅ CI/CD           main       2025-11-26T11:05
#19673553781 ✅ CI/CD           main       2025-11-26T11:00
#19673553780 ✅ CI/CD           main       2025-11-26T10:55
```

## 🔮 Future Enhancements

### Planned Features
- [ ] Real-time notifications
- [ ] Integration with Slack/Discord
- [ ] Performance benchmarking
- [ ] Custom alert thresholds
- [ ] Multi-branch support
- [ ] Automated fix application
- [ ] Integration with IDE plugins

### Advanced Analytics
- [ ] CI/CD performance trends
- [ ] Build time optimization suggestions
- [ ] Resource usage analysis
- [ ] Cost optimization insights

## 🤝 Contributing

To enhance the script:

1. Follow the existing code style
2. Add comprehensive error handling
3. Include helpful comments
4. Test with various CI/CD scenarios
5. Update documentation

## 📝 License

This enhanced script is part of the EchoMind project and follows the same license terms.

---

**Pro Tip**: For the best experience, add the script to your shell aliases:
```bash
echo 'alias ci="./scripts/check_ci_enhanced.sh"' >> ~/.zshrc
source ~/.zshrc
```

Now you can simply run `ci` to check your CI/CD status!