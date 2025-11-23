# i18n 配置管理脚本

## 快速开始

检测 i18n 配置中的冗余项和代码中的硬编码字符串。

```bash
# 1. 检查冗余配置（只读）
node scripts/check_i18n.js

# 2. 自动清理冗余配置（创建备份）
node scripts/check_i18n.js --fix

# 3. 检测硬编码的中英文字符串
node scripts/check_i18n.js --detect-hardcoded

# 4. 完整检查和清理
node scripts/check_i18n.js --fix --detect-hardcoded
```

## 当前状态

- **字典总键数**: 256
- **使用中的键**: 204
- **冗余键**: 81 ⚠️

主要冗余模块：
- Onboarding（27个键）- 未实现的引导流程
- Settings 子模块（12个键）
- Dashboard 细粒度字段（18个键）

## 功能特性

✅ 自动检测未使用的 i18n 配置项  
✅ 自动备份后删除冗余配置  
✅ 检测硬编码的中文字符串  
✅ 检测硬编码的英文 UI 文本  
✅ 双语配置同步检查（en.json + zh.json）

详细文档请参考工件中的 `i18n_guide.md`
