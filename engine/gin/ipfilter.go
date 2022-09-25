package gin

import (
	"net/http"
	"errors"
	krakend "github.com/NEOMorphey/krakend-ipfilter/krakend"
	"github.com/gin-gonic/gin"
	"github.com/luraproject/lura/config"
	"github.com/luraproject/lura/logging"
	"github.com/luraproject/lura/proxy"
	krakendgin "github.com/luraproject/lura/router/gin"
	ipfilter "github.com/NEOMorphey/krakend-ipfilter"
)

func Register(cfg config.ServiceConfig, l logging.Logger, engine *gin.Engine) {
	detectorCfg, err := krakend.ParseConfig(cfg.ExtraConfig)
	if err == krakend.ErrNoConfig {
		l.Debug("ipfilter middleware: ", err.Error())
		return
	}
	if err != nil {
		l.Warning("ipfilter middleware: ", err.Error())
		return
	}
	d, err := ipfilter.NewIPFilter(detectorCfg)
	if err != nil {
		l.Warning("ipfilter middleware: unable to createt the LRU detector:", err.Error())
		return
	}
	engine.Use(middleware(d))
}


// New checks the configuration and, if required, wraps the handler factory with a bot detector middleware
func New(hf krakendgin.HandlerFactory, l logging.Logger) krakendgin.HandlerFactory {
	return func(cfg *config.EndpointConfig, p proxy.Proxy) gin.HandlerFunc {
		next := hf(cfg, p)

		detectorCfg, err := krakend.ParseConfig(cfg.ExtraConfig)
		if err == krakend.ErrNoConfig {
			l.Debug("ipfilter: ", err.Error())
			return next
		}
		if err != nil {
			l.Warning("ipfilter: ", err.Error())
			return next
		}

		d, err := ipfilter.NewIPFilter(detectorCfg)
		if err != nil {
			l.Warning("ipfilter: unable to create the LRU detector:", err.Error())
			return next
		}
		return handler(d, next)
	}
}

func middleware(f ipfilter.IPFilter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if f.Deny(ip) {
			c.AbortWithError(http.StatusForbidden, errIPRejected)
			return
		}

		c.Next()
	}
}

func handler(f ipfilter.IPFilter, next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if f.Deny(ip) {
			c.AbortWithError(http.StatusForbidden, errIPRejected)
			return
		}
		next(c)
	}
}

var errIPRejected = errors.New("Access Forbidden")

