# EchoMind CI/CD Monitor

**Simplified, Robust, and Powerful CI/CD Monitoring Tool**

## 🚀 Quick Start

```bash
# Current status (default)
./scripts/ci.sh

# Create a simple alias
echo 'alias ci="./scripts/ci.sh"' >> ~/.zshrc
source ~/.zshrc

# Now just use:
ci
```

## 📋 Usage

### Basic Commands

```bash
# Show current CI/CD status
./scripts/ci.sh

# Watch live running workflow
./scripts/ci.sh watch

# Show last 5 runs (default)
./scripts/ci.sh history

# Show last N runs
./scripts/ci.sh history 10

# Deep analysis of recent runs
./scripts/ci.sh analyze

# Interactive menu
./scripts/ci.sh interactive

# Show help
./scripts/ci.sh help
```

### Environment Variables

```bash
# Set custom defaults
HISTORY_COUNT=10 ./scripts/ci.sh history
BRANCH=develop ./scripts/ci.sh

# Export for session
export HISTORY_COUNT=10
export BRANCH=feature/new
./scripts/ci.sh
```

## 🎯 Features

### Core Capabilities
- **✅ Robust Error Handling** - Graceful failure recovery
- **⚡ Fast Performance** - Optimized data fetching
- **🎨 Clean Output** - Color-coded, emoji-enhanced display
- **🔧 Smart Defaults** - Sensible configuration out of the box
- **📊 Progress Indicators** - Visual feedback during operations

### Advanced Functions
- **🔍 Smart Error Analysis** - Automatic pattern detection
- **💡 Fix Suggestions** - Actionable recommendations
- **📈 Success Rate Tracking** - CI/CD health monitoring
- **🔄 Live Watching** - Real-time workflow monitoring
- **🎛️ Interactive Menu** - User-friendly navigation

## 🔧 Error Detection

The script automatically detects and categorizes:

### TypeScript Errors
```bash
error TS2307: Cannot find module 'next-themes/dist/types'
```

### Test Failures
```bash
FAIL src/components/Widget.test.tsx
Test Suites: 1 failed, 5 passed
```

### Build/Dependency Issues
```bash
npm ERR! code ERESOLVE
ELIFECYCLE Command failed with exit code 1
```

### Permission/Docker Issues
```bash
permission denied: docker
Error response from daemon
```

## 💡 Smart Fix Suggestions

Based on detected errors, the script provides targeted solutions:

### TypeScript Issues
- Check import paths in theme components
- Verify type declarations
- Update type dependencies

### Build Issues
- Run `npm install` or `pnpm install`
- Clear node_modules and reinstall

### Permission Issues
- Check file permissions
- Verify script execution rights

## 📊 Analytics Example

```bash
$ ./scripts/ci.sh analyze

🔍 EchoMind CI/CD Monitor
========================================
⏳ Performing deep analysis... ████████████████████████████████████ 100%

📊 CI/CD Statistics:
----------------------------------------
Total Runs: 10
Successful: 8
Failed: 2
Success Rate: 80%

💡 Recommendations:
  ✅• Use pre-commit hooks to catch issues early
  ✅• Monitor build times and optimize dependencies
  ✅• Consider parallel test execution
```

## 🎨 Visual Features

### Progress Indicators
```bash
[⏳] Checking latest run [███████████████████████████████████] 100%
```

### Status Display
```bash
✅ CI/CD Run #19691211695
Workflow: CI/CD
Branch: main
Run ID: 19691211695
Created: 2025-11-26T03:10:00Z
Status: SUCCESS
Link: https://github.com/hrygo/echomind/actions/runs/19691211695
```

### History View
```bash
📊 CI/CD History:
----------------------------------------
#19691211695 ✅ CI/CD                   main       2025-11-26T03:10
#19673553783 ❌ CI/CD                   main       2025-11-25T14:46
#19672016464 ❌ CI/CD                   main       2025-11-25T13:57
```

## 🛠️ Requirements

### Dependencies
```bash
# Install GitHub CLI
brew install gh

# Install jq (JSON processor)
brew install jq

# Authenticate with GitHub
gh auth login
```

## 🔍 Troubleshooting

### Common Issues

1. **"GitHub CLI not authenticated"**
   ```bash
   gh auth login
   ```

2. **"Missing required dependencies"**
   ```bash
   brew install gh jq
   ```

3. **"No CI/CD runs found"**
   - Push a commit to trigger a workflow
   - Check repository permissions

### Debug Mode
```bash
# Enable verbose output
set -x
./scripts/ci.sh
set +x
```

## 🔄 Migration from Old Scripts

| Old Script | New Command | Description |
|------------|-------------|-------------|
| `./scripts/check_ci.sh` | `./scripts/ci.sh` | Current status |
| `./scripts/check_ci_enhanced.sh --history` | `./scripts/ci.sh history` | Show history |
| `./scripts/check_ci_enhanced.sh --watch` | `./scripts/ci.sh watch` | Watch live |
| `./scripts/check_ci_enhanced.sh --analyze` | `./scripts/ci.sh analyze` | Deep analysis |
| `./scripts/check_ci_enhanced.sh --interactive` | `./scripts/ci.sh interactive` | Interactive mode |

## 📝 Best Practices

### Daily Development
```bash
# Quick status check
ci

# After pushing code
ci watch
```

### CI/CD Issues
```bash
# When CI fails
ci analyze

# Quick rerun of failed job
ci interactive  # Choose option 5
```

### Team Usage
```bash
# Share CI status with team
ci history 10

# Monitor branch performance
BRANCH=feature/api ci history
```

## 🎯 Tips & Tricks

### Aliases for Productivity
```bash
# Add to ~/.zshrc or ~/.bashrc
alias ci='./scripts/ci.sh'
alias ciw='./scripts/ci.sh watch'
alias cih='./scripts/ci.sh history'
alias cia='./scripts/ci.sh analyze'
alias cii='./scripts/ci.sh interactive'

# Custom configurations
export HISTORY_COUNT=15
export BRANCH=main
```

### Git Integration
```bash
# Check CI before merging
ci && git merge feature-branch

# Watch CI after push
git push origin main && ci watch
```

### Automation
```bash
# CI status in prompt (add to .zshrc)
precmd() {
    local ci_status=$(./scripts/ci.sh 2>/dev/null | grep -o "✅\|❌" | head -1)
    [[ -n "$ci_status" ]] && echo -n "$ci_status "
}
```

## 🚀 Performance

### Optimizations
- **Fast Data Fetching**: Single API call with smart filtering
- **Memory Efficient**: Temporary files auto-cleaned
- **Network Optimized**: Minimal GitHub API calls
- **Error Resilient**: Graceful handling of API failures

### Benchmarks
- **Status Check**: ~2 seconds
- **History (10 runs)**: ~3 seconds
- **Deep Analysis**: ~5 seconds
- **Memory Usage**: <10MB peak

## 🔮 Future Enhancements

### Planned Features
- [ ] Real-time notifications
- [ ] Slack/Discord integration
- [ ] Custom alert thresholds
- [ ] Multi-branch comparison
- [ ] Performance benchmarking
- [ ] Automated fix suggestions

## 🤝 Contributing

### Development Setup
```bash
# Clone and setup
git clone <repo>
cd echomind
chmod +x scripts/ci.sh

# Test changes
./scripts/ci.sh help
./scripts/ci.sh history 3
```

### Adding Features
1. Follow existing code style
2. Add comprehensive error handling
3. Include helpful comments
4. Test with various scenarios
5. Update documentation

## 📄 License

This script is part of the EchoMind project and follows the same license terms.

---

**🎉 Pro Tip**: Set up the alias `ci=./scripts/ci.sh` for instant CI/CD monitoring with just two keystrokes!