package wechat

import (
	"context"

	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/spf13/viper"
)

// Gateway encapsulates WeChat Official Account SDK
type Gateway struct {
	wc              *wechat.Wechat
	OfficialAccount *officialaccount.OfficialAccount
}

// NewGateway initializes the WeChat Gateway with configuration from Viper
func NewGateway(redisAddr, redisPassword string, redisDB int) (*Gateway, error) {
	// 1. Initialize Redis Cache for AccessToken storage
	ctx := context.Background()
	redisCache := cache.NewRedis(ctx, &cache.RedisOpts{
		Host:     redisAddr,
		Password: redisPassword,
		Database: redisDB,
	})

	// 2. Initialize WeChat SDK
	wc := wechat.NewWechat()
	wc.SetCache(redisCache)

	// 3. Load Official Account Config
	cfg := &offConfig.Config{
		AppID:          viper.GetString("wechat.app_id"),
		AppSecret:      viper.GetString("wechat.app_secret"),
		Token:          viper.GetString("wechat.token"),
		EncodingAESKey: viper.GetString("wechat.encoding_aes_key"),
	}

	officialAccount := wc.GetOfficialAccount(cfg)

	return &Gateway{
		wc:              wc,
		OfficialAccount: officialAccount,
	}, nil
}
