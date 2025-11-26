# 🏗️ WeChat OS 技术方案

> **版本**: v1.0 | **状态**: 规划中

## 1. 总体架构

WeChat OS 模块将作为 EchoMind 后端的一个独立服务或模块存在，通过 API Gateway 暴露给微信服务器。

```mermaid
graph TD
    WeChat_Server[微信服务器] -- Webhook (XML/JSON) --> API_Gateway[Nginx / API Gateway]
    API_Gateway --> WeChat_Handler[WeChat Gateway (Go)]
    
    subgraph EchoMind Backend
        WeChat_Handler -- 1. 消息接收 --> Msg_Router[消息路由]
        Msg_Router -- 2. 文本/语音 --> AI_Pipeline[AI 处理管道]
        Msg_Router -- 2. 事件(点击) --> Action_Svc[Action Service]
        
        AI_Pipeline -- 3. 意图识别 --> Intent_Analyzer[LLM Intent Analyzer]
        Intent_Analyzer -- 4. 执行指令 --> Biz_Services[业务服务 (Mail/Calendar/Task)]
        
        Biz_Services -- 5. 结果数据 --> Response_Generator[回复生成器]
        Response_Generator -- 6. 格式化消息 --> WeChat_Sender[WeChat API Client]
        WeChat_Sender --> WeChat_Server
    end
    
    subgraph Data Stores
        WeChat_Handler -- Read/Write --> Redis_FSM[Redis (FSM State)]
        WeChat_Handler -- Read --> DB_User[PostgreSQL (User Bindings)]
    end
```

## 2. 核心模块设计

### 2.1 微信接入网关 (WeChat Gateway)
*   **职责**: 处理微信服务器的验证 (Token Verify)、消息接收 (XML Parse)、被动回复和主动推送。
*   **技术选型**: 采用 `github.com/silenceper/wechat/v2` (v2.1.10+)。
    *   利用其内置的 `cache.Redis` 模块管理 AccessToken，确保多实例部署时的一致性。
    *   使用 `officialaccount` 模块处理消息路由。
*   **Endpoint**: `POST /api/v1/wechat/callback`

### 2.2 多轮对话状态机 (FSM)
由于微信交互是无状态的 HTTP 请求，我们需要在 Redis 中维护用户的会话状态。
*   **存储**: Redis Key `echomind:fsm:{openid}`
*   **状态定义**:
    *   `IDLE`: 空闲状态，等待指令。
    *   `WAIT_CMD_CONFIRM`: 等待用户确认执行（如“是否发送？”）。
    *   `WAIT_SLOT_FILLING`: 等待用户补充信息（如“回复给谁？”）。
*   **过期策略**: 会话状态 TTL 设置为 5-10 分钟，超时自动重置为 `IDLE`。

### 2.3 语音处理管道 (Voice Pipeline)
1.  **接收**: 接收微信 `MsgType=voice` 消息，获取 `MediaId`。
2.  **下载**: 调用微信 API 下载 AMR/MP3 音频文件。
3.  **转录**:
    *   方案 A (优先): 使用微信自带的 `Recognition` 字段（如果开通了语音识别接口）。
    *   方案 B: 调用 OpenAI Whisper API 进行高精度转录。
4.  **处理**: 转录后的文本送入 `Intent Analyzer`。

### 2.4 意图识别与执行 (Intent Analyzer)
使用 LLM (DeepSeek/GPT-4o-mini) 进行 Function Calling / Tool Use。
*   **Prompt**: 定义一组工具 (Tools)，如 `search_emails`, `reply_email`, `check_calendar`, `create_task`.
*   **流程**:
    1.  User Query -> LLM -> Tool Call (JSON)
    2.  Backend Execute Tool -> Result
    3.  Result + History -> LLM -> Final Response (Text)

### 2.5 消息推送 (Push Notification)
*   **机制**: 使用微信“客服消息”接口（48小时内活跃）或“模板消息”（服务号）。
*   **触发源**:
    *   `Morning Briefing`: Cron Job 触发。
    *   `One-Touch Decision`: 邮件分析服务 (`EmailIngestor`) 发现高优先级且需决策邮件时触发。

## 3. 数据模型变更

### 3.1 User 表扩展
需要存储微信 OpenID 和 UnionID。

```go
type User struct {
    // ... existing fields
    WeChatOpenID  string `gorm:"uniqueIndex"`
    WeChatUnionID string `gorm:"index"`
    WeChatState   string // JSON blob for user preferences (e.g., briefing time)
}
```

### 3.2 消息日志
记录微信交互日志用于调试和优化。

```go
type WeChatLog struct {
    ID        uint
    UserID    uint
    OpenID    string
    MsgType   string // text, voice, event
    Content   string // 脱敏后的内容
    Direction string // inbound, outbound
    CreatedAt time.Time
}
```

## 4. 安全设计
*   **签名验证**: 严格验证微信回调的 Signature。
*   **账号绑定**:
    *   用户在 Web 端生成带参数的二维码。
    *   用户微信扫码，回调触发绑定逻辑。
    *   **禁止**直接在微信端输入账号密码。
*   **敏感数据**: 推送消息中不包含完整邮件正文，仅包含摘要和链接。

## 5. 部署依赖
*   **Redis**: 必须高可用，用于存储 FSM 和 AccessToken。
*   **公网域名**: 微信回调需要公网可访问的 HTTPS 域名 (开发环境使用 Ngrok/Cloudflare Tunnel)。
