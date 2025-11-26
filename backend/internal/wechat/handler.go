package wechat

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hrygo/echomind/pkg/logger"
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

// Handler handles WeChat callback requests
type Handler struct {
	gateway *Gateway
	logger  logger.Logger
}

// NewHandler creates a new WeChat handler
func NewHandler(gateway *Gateway, log logger.Logger) *Handler {
	return &Handler{
		gateway: gateway,
		logger:  log,
	}
}

// Callback handles the WeChat server callback
// GET: Server verification (token check)
// POST: Message receiving and passive reply
func (h *Handler) Callback(c *gin.Context) {
	server := h.gateway.OfficialAccount.GetServer(c.Request, c.Writer)

	// Set message handler for incoming messages
	server.SetMessageHandler(h.handleMessage)

	// Serve the request (handles both GET verification and POST message)
	err := server.Serve()
	if err != nil {
		h.logger.Error("WeChat server error",
			logger.Field{Key: "error", Value: err},
			logger.Field{Key: "path", Value: c.Request.URL.Path},
		)
		c.String(http.StatusInternalServerError, "Server Error")
		return
	}

	// If all goes well, the SDK has already written the response
	c.String(http.StatusOK, "success")
}

// handleMessage processes incoming messages and returns responses
func (h *Handler) handleMessage(msg *message.MixMessage) *message.Reply {
	h.logger.Info("Received WeChat message",
		logger.Field{Key: "from_user", Value: string(msg.FromUserName)},
		logger.Field{Key: "msg_type", Value: string(msg.MsgType)},
		logger.Field{Key: "content", Value: msg.Content},
	)

	// TODO: Route to Intent Analyzer or FSM
	// For now, echo the message back
	switch msg.MsgType {
	case message.MsgTypeText:
		return &message.Reply{
			MsgType: message.MsgTypeText,
			MsgData: message.NewText("收到：" + msg.Content),
		}
	case message.MsgTypeVoice:
		// TODO: Phase 7.2 - Voice processing
		return &message.Reply{
			MsgType: message.MsgTypeText,
			MsgData: message.NewText("语音功能即将上线！"),
		}
	case message.MsgTypeEvent:
		// Handle events like SCAN, SUBSCRIBE
		return h.handleEvent(msg)
	default:
		return nil
	}
}

// handleEvent processes WeChat events (e.g., SCAN for account binding)
func (h *Handler) handleEvent(msg *message.MixMessage) *message.Reply {
	h.logger.Info("Received WeChat event",
		logger.Field{Key: "event", Value: string(msg.Event)},
		logger.Field{Key: "event_key", Value: msg.EventKey},
		logger.Field{Key: "from_user", Value: string(msg.FromUserName)},
	)

	switch msg.Event {
	case message.EventSubscribe:
		return &message.Reply{
			MsgType: message.MsgTypeText,
			MsgData: message.NewText("欢迎关注 EchoMind 🧠！请前往网页端扫描二维码绑定账号。"),
		}
	case message.EventScan:
		// TODO: Phase 7.1 Week 2 - Implement account binding
		return &message.Reply{
			MsgType: message.MsgTypeText,
			MsgData: message.NewText("账号绑定功能开发中..."),
		}
	default:
		return nil
	}
}
