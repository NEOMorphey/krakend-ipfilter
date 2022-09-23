package gin

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luraproject/lura/config"
	"github.com/luraproject/lura/logging"
	ipfilter "github.com/NEOMorphey/krakend-ipfilter"
)

// Register register a ip filter middleware at gin
func Register(extraConfig config.ExtraConfig, logger logging.Logger, engine *gin.Engine) {
	logPrefix := "[Service: Gin][IPFilter]"
	filterCfg := ipfilter.ParseConfig(ExtraConfig, logger)
	if filterCfg == nil {
		return
	}
	logger.Debug(logPrefix, "Middleware is now ready")
	ipFilter := ipfilter.NewIPFilter(filterCfg)
	engine.Use(middleware(ipFilter, logger))
}

func middleware(ipFilter ipfilter.IPFilter, logger logging.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()
		if ipFilter.Deny(ip) {
			logger.Error(fmt.Sprintf("krakend-ipfilter deny request from: %s", ip))
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}

		ctx.Next()
	}
}
