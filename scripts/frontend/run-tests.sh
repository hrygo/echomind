#!/bin/bash

# 端到端测试执行脚本
# 此脚本运行所有测试并生成测试报告

set -e  # 遇到错误时退出

# 获取脚本所在目录和项目根目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
FRONTEND_DIR="$PROJECT_ROOT/frontend"

# 切换到前端目录
cd "$FRONTEND_DIR"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查前置条件
check_prerequisites() {
    log_info "检查前置条件..."
    
    # 检查 Node.js
    if ! command -v node &> /dev/null; then
        log_error "未找到 Node.js，请先安装 Node.js"
        exit 1
    fi
    log_success "Node.js 已安装: $(node --version)"
    
    # 检查 pnpm
    if ! command -v pnpm &> /dev/null; then
        log_error "未找到 pnpm，请先安装 pnpm"
        exit 1
    fi
    log_success "pnpm 已安装: $(pnpm --version)"
    
    # 检查依赖是否已安装
    if [ ! -d "node_modules" ]; then
        log_warning "依赖未安装，正在安装..."
        pnpm install
    fi
    
    log_success "前置条件检查完成"
}

# 运行类型检查
run_type_check() {
    log_info "运行 TypeScript 类型检查..."
    
    if pnpm type-check; then
        log_success "类型检查通过"
        return 0
    else
        log_error "类型检查失败"
        return 1
    fi
}

# 运行单元测试
run_unit_tests() {
    log_info "运行单元测试..."
    
    if pnpm test; then
        log_success "单元测试通过"
        return 0
    else
        log_error "单元测试失败"
        return 1
    fi
}

# 运行 E2E 测试
run_e2e_tests() {
    log_info "运行端到端测试..."
    
    # 检查后端服务
    log_info "检查后端服务..."
    if ! curl -s http://localhost:8080/health > /dev/null 2>&1; then
        log_warning "后端服务未运行，某些 E2E 测试可能失败"
        log_warning "请确保后端服务运行在 http://localhost:8080"
        
        read -p "是否继续运行 E2E 测试? (y/n) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            log_info "跳过 E2E 测试"
            return 2
        fi
    else
        log_success "后端服务正常运行"
    fi
    
    # 运行 Playwright 测试
    if pnpm playwright test; then
        log_success "E2E 测试通过"
        return 0
    else
        log_error "E2E 测试失败"
        return 1
    fi
}

# 生成测试报告
generate_report() {
    log_info "生成测试报告..."
    
    # 生成 Playwright 报告
    if [ -d "playwright-report" ]; then
        log_info "Playwright 报告可用，运行 'pnpm playwright show-report' 查看"
    fi
    
    # 生成 Jest 覆盖率报告
    if [ -d "coverage" ]; then
        log_info "Jest 覆盖率报告生成在 coverage/ 目录"
    fi
    
    log_success "测试报告生成完成"
}

# 清理函数
cleanup() {
    log_info "清理临时文件..."
    # 在这里添加清理逻辑
}

# 主函数
main() {
    echo "======================================"
    echo "   AI Native 架构 - 端到端测试"
    echo "======================================"
    echo ""
    
    # 注册清理函数
    trap cleanup EXIT
    
    # 运行检查和测试
    check_prerequisites
    
    TYPE_CHECK_RESULT=0
    UNIT_TEST_RESULT=0
    E2E_TEST_RESULT=0
    
    # 类型检查
    if ! run_type_check; then
        TYPE_CHECK_RESULT=1
    fi
    
    # 单元测试
    if ! run_unit_tests; then
        UNIT_TEST_RESULT=1
    fi
    
    # E2E 测试
    run_e2e_tests
    E2E_TEST_RESULT=$?
    
    # 生成报告
    generate_report
    
    # 汇总结果
    echo ""
    echo "======================================"
    echo "          测试结果汇总"
    echo "======================================"
    
    if [ $TYPE_CHECK_RESULT -eq 0 ]; then
        log_success "✓ TypeScript 类型检查通过"
    else
        log_error "✗ TypeScript 类型检查失败"
    fi
    
    if [ $UNIT_TEST_RESULT -eq 0 ]; then
        log_success "✓ 单元测试通过"
    else
        log_error "✗ 单元测试失败"
    fi
    
    if [ $E2E_TEST_RESULT -eq 0 ]; then
        log_success "✓ E2E 测试通过"
    elif [ $E2E_TEST_RESULT -eq 2 ]; then
        log_warning "⊘ E2E 测试被跳过"
    else
        log_error "✗ E2E 测试失败"
    fi
    
    echo ""
    
    # 返回总体结果
    if [ $TYPE_CHECK_RESULT -eq 0 ] && [ $UNIT_TEST_RESULT -eq 0 ] && [ $E2E_TEST_RESULT -lt 2 ]; then
        log_success "所有测试通过！🎉"
        exit 0
    else
        log_error "部分测试失败，请检查日志"
        exit 1
    fi
}

# 运行主函数
main
