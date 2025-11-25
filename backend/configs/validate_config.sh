#!/bin/bash

# 配置验证脚本
# 使用方法: ./validate_config.sh [配置文件路径]

CONFIG_FILE=${1:-"config.yaml"}
EXAMPLE_FILE="config.example.yaml"

echo "🔍 配置验证工具"
echo "=================="

# 检查配置文件是否存在
if [ ! -f "$CONFIG_FILE" ]; then
    echo "❌ 错误: 配置文件 '$CONFIG_FILE' 不存在"
    echo "💡 提示: 运行 'cp $EXAMPLE_FILE $CONFIG_FILE' 创建配置文件"
    exit 1
fi

echo "✅ 找到配置文件: $CONFIG_FILE"

# 检查YAML语法
echo "📋 检查YAML语法..."
if python3 -c "import yaml; yaml.safe_load(open('$CONFIG_FILE'))" 2>/dev/null; then
    echo "✅ YAML语法正确"
else
    echo "❌ YAML语法错误"
    exit 1
fi

# 检查敏感数据
echo "🔒 检查敏感数据..."
SENSITIVE_PATTERNS=("sk-" "AIza" "dummy-jwt-secret" "dummy-openai-key" "dummy-moonshot-key")

for pattern in "${SENSITIVE_PATTERNS[@]}"; do
    if grep -q "$pattern" "$CONFIG_FILE"; then
        echo "⚠️  警告: 发现可能的敏感数据 (包含 '$pattern')"
        echo "💡 提示: 请确保这是你的私有配置文件，不是示例文件"
    fi
done

# 检查必需字段
echo "🔧 检查必需配置字段..."

REQUIRED_FIELDS=(
    "server.port"
    "server.jwt.secret"
    "security.encryption_key"
    "database.dsn"
    "ai.active_services.chat"
    "ai.active_services.embedding"
)

MISSING_FIELDS=()
for field in "${REQUIRED_FIELDS[@]}"; do
    if ! grep -q "$field:" "$CONFIG_FILE"; then
        MISSING_FIELDS+=("$field")
    fi
done

if [ ${#MISSING_FIELDS[@]} -eq 0 ]; then
    echo "✅ 所有必需字段都已配置"
else
    echo "❌ 缺少必需字段:"
    for field in "${MISSING_FIELDS[@]}"; do
        echo "   - $field"
    done
fi

# 检查AI提供商配置
echo "🤖 检查AI提供商配置..."

CHAT_PROVIDER=$(grep -A 2 "active_services:" "$CONFIG_FILE" | grep "chat:" | awk '{print $2}' | tr -d ' ')
EMBEDDING_PROVIDER=$(grep -A 3 "active_services:" "$CONFIG_FILE" | grep "embedding:" | awk '{print $2}' | tr -d ' ')

if [ -n "$CHAT_PROVIDER" ] && [ -n "$EMBEDDING_PROVIDER" ]; then
    echo "✅ AI提供商配置:"
    echo "   - Chat: $CHAT_PROVIDER"
    echo "   - Embedding: $EMBEDDING_PROVIDER"

    # 检查提供商是否存在
    if grep -q "  $CHAT_PROVIDER:" "$CONFIG_FILE"; then
        echo "✅ Chat提供商 '$CHAT_PROVIDER' 已配置"
    else
        echo "❌ Chat提供商 '$CHAT_PROVIDER' 配置缺失"
    fi

    if grep -q "  $EMBEDDING_PROVIDER:" "$CONFIG_FILE"; then
        echo "✅ Embedding提供商 '$EMBEDDING_PROVIDER' 已配置"
    else
        echo "❌ Embedding提供商 '$EMBEDDING_PROVIDER' 配置缺失"
    fi
else
    echo "❌ AI提供商配置不完整"
fi

# 检查嵌入维度配置
echo "📏 检查嵌入维度配置..."
if grep -q "embedding_dimensions:" "$CONFIG_FILE"; then
    echo "✅ 找到嵌入维度配置"

    # 提取所有嵌入维度
    DIMENSIONS=$(grep "embedding_dimensions:" "$CONFIG_FILE" | awk '{print $2}')
    echo "   配置的维度: $DIMENSIONS"

    # 检查是否有冲突
    UNIQUE_DIMENSIONS=$(echo "$DIMENSIONS" | sort -u)
    if [ $(echo "$DIMENSIONS" | wc -l) -ne $(echo "$UNIQUE_DIMENSIONS" | wc -l) ]; then
        echo "⚠️  警告: 发现不同的嵌入维度配置，确保数据库schema匹配"
    fi
else
    echo "⚠️  警告: 未找到嵌入维度配置"
fi

# 检查数据库连接字符串
echo "🗄️ 检查数据库配置..."
DB_DSN=$(grep "database.dsn:" "$CONFIG_FILE" | awk '{print $2}' | tr -d '"')

if [ -n "$DB_DSN" ]; then
    if echo "$DB_DSN" | grep -q "password"; then
        echo "⚠️  警告: 数据库连接字符串包含密码，请确保安全"
    fi

    if echo "$DB_DSN" | grep -q "sslmode=disable"; then
        echo "⚠️  警告: SSL已禁用，仅适用于开发环境"
    fi
else
    echo "❌ 数据库连接字符串未配置"
fi

echo ""
echo "🎯 验证完成!"
echo "=================="

# 给出配置建议
echo "💡 配置建议:"
echo "1. 生产环境请使用真实的JWT密钥和加密密钥"
echo "2. 确保所有API密钥都已配置"
echo "3. 检查嵌入维度与数据库schema是否匹配"
echo "4. 生产环境请启用SSL"
echo ""
echo "📚 更多信息请查看: configs/README.md"