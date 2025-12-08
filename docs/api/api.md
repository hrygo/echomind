# EchoMind API 文档

本文档是 EchoMind API 的主要参考。为了获得更好的交互式体验和最准确的定义，请参考我们的 OpenAPI 3.0 规范。

## 📚 相关文档

- **[OpenAPI 规范文件 (openapi.yaml)](./openapi.yaml)** - 完整的 API 接口定义
- **[技术架构文档 (architecture.md)](./architecture.md)** - 系统架构、向量搜索、AI 服务等技术细节
- **[产品需求文档 (prd.md)](./prd.md)** - 产品规划和功能说明
- **[产品设计文档 (product-design.md)](./product-design.md)** - 产品设计和用户体验

## 核心技术特性



### 🤖 AI 服务抽象层
- **统一接口**: 通过标准化的 Provider 接口支持多种 AI 服务
- **配置驱动**: 通过配置文件动态选择 AI 供应商
- **协议支持**: OpenAI 协议、Gemini 协议、Mock 测试协议

## 使用方法

您可以使用任何兼容 OpenAPI 3.0 的工具来查看和交互此规范，例如：

- [Swagger Editor](https://editor.swagger.io/): 一个在线编辑器，用于查看、编辑和测试 OpenAPI 规范。
- [Redocly](https://redocly.github.io/redoc/): 生成一个三栏式的、响应式的文档页面。
- [Postman](https://www.postman.com/): 可以导入 OpenAPI 规范来自动创建请求集合。

## 🔧 开发指南

### 本地开发环境设置
请参考后端配置文档: [backend/configs/README.md](../backend/configs/README.md)

### API 接口示例
```bash
# 搜索邮件 (RAG)
curl -X POST "http://localhost:8080/api/emails/search" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"query": "项目进展", "limit": 10}'

# 获取邮件摘要
curl -X GET "http://localhost:8080/api/emails/{id}/summary" \
  -H "Authorization: Bearer <token>"
```
