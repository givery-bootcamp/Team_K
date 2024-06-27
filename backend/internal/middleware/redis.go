package middleware

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"myapp/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetupRedis() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		redis := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", config.RedisHost, config.RedisPort),
			Password: config.RedisPassword,
			DB:       0,                    // use default DB
			TLSConfig: &tls.Config{
				RootCAs:    x509.NewCertPool(),},
		})
		ctx.Set("redis", redis)
		ctx.Next()
	}
}
