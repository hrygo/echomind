# 📊 WeChat SDK 选型评估报告

> **日期**: 2025-11-26
> **状态**: ✅ 已决策
> **决策**: 采用 `silenceper/wechat`

## 1. 候选方案

我们调研了 Go 语言生态中主流的微信 SDK，主要关注以下指标：
*   **成熟度**: GitHub Stars, Issues 响应速度, 版本更新频率。
*   **功能覆盖**: 公众号 (Official Account), 微信支付 (WeChat Pay), 小程序 (Mini Program)。
*   **文档质量**: 是否有清晰的文档和示例。

### 1.1 [silenceper/wechat](https://github.com/silenceper/wechat)
*   **Stars**: 5.2k+
*   **License**: MIT
*   **特点**:
    *   Go 生态中最流行的微信 SDK。
    *   模块化设计，支持公众号、小程序、微信支付、企业微信。
    *   内置缓存接口 (Redis/Memcache)，方便集成。
    *   活跃度高，最近发布于 1 个月前 (v2.1.10)。
*   **适用性**: ⭐⭐⭐⭐⭐ (完全符合)

### 1.2 [PowerWechat](https://github.com/ArtisanCloud/PowerWechat)
*   **Stars**: 1.5k+
*   **License**: MIT
*   **特点**:
    *   全功能覆盖，API 设计参考了 PHP 的 EasyWeChat。
    *   文档较详细。
*   **适用性**: ⭐⭐⭐⭐ (备选)

### 1.3 [eatmoreapple/openwechat](https://github.com/eatmoreapple/openwechat)
*   **Stars**: 5.5k+
*   **特点**:
    *   主要用于**个人号** (模拟网页版微信客户端)。
    *   **不适用**于本项目。我们需要的是基于**公众号/服务号**的官方接口开发，而非个人号 Hook。

## 2. 详细评估: silenceper/wechat

### 2.1 核心优势
1.  **开箱即用**: 提供了标准的 `OfficialAccount` 结构体，直接封装了消息接收 (`Serve`) 和被动回复逻辑。
2.  **Redis 集成**: 我们的架构中已经包含 Redis，该 SDK 的 `cache` 接口可以直接对接，用于存储 `access_token`，无需重复造轮子。
3.  **扩展性**: 支持自定义 `Context`，方便我们在处理消息时注入 Trace ID 或 User Info。

### 2.2 代码示例 (验证)

```go
package main

import (
    "github.com/silenceper/wechat/v2"
    "github.com/silenceper/wechat/v2/cache"
    "github.com/silenceper/wechat/v2/officialaccount"
    "github.com/silenceper/wechat/v2/officialaccount/config"
)

func main() {
    // 1. 初始化 Redis 缓存
    redisCache := cache.NewRedis(&cache.RedisOpts{Host: "localhost:6379"})

    // 2. 初始化 WeChat 实例
    wc := wechat.NewWechat()
    wc.SetCache(redisCache)

    // 3. 配置公众号参数
    cfg := &config.Config{
        AppID:          "your_app_id",
        AppSecret:      "your_app_secret",
        Token:          "your_token",
        EncodingAESKey: "your_encoding_aes_key",
    }
    officialAccount := wc.GetOfficialAccount(cfg)

    // 4. 处理消息 (在 Controller 中调用)
    // server := officialAccount.GetServer(req, writer)
    // server.SetMessageHandler(func(msg *message.MixMessage) *message.Reply {
    //     if msg.Content == "ping" {
    //         return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText("pong")}
    //     }
    //     return nil
    // })
    // server.Serve()
}
```

## 3. 结论

确认在 **WeChat OS** 模块中使用 `github.com/silenceper/wechat/v2` 作为核心 SDK。它成熟、稳定且功能完备，能够满足我们对公众号消息处理、OAuth 授权和模板消息推送的所有需求。
