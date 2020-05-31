package dist

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

/**
TODO: 流控，debug，cors，(自动)黑名单
*/
func http_flow_control_middleware() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func http_blacklist_middleware(bl map[string]bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if _, ok := bl[ip]; ok {
			c.AbortWithStatus(http.StatusProxyAuthRequired)
		}
	}
}

func http_whitelist_middleware(wl map[string]bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if _, ok := wl[ip]; !ok {
			c.AbortWithStatus(http.StatusProxyAuthRequired)
		}
	}
}

type HttpGin struct {
	addr   string
	server *http.Server

	enableCors        bool
	enableDebug       bool
	enableFlowControl bool
	enableSession     bool
	sessionStore      sessions.Store

	enableHTML bool
	htmlPtn    string

	blacklist map[string]bool
	whitelist map[string]bool
}

type HttpGinOptions func(*HttpGin)

func EnableDebug(enable bool) HttpGinOptions {
	return func(h *HttpGin) {
		h.enableDebug = enable
	}
}

func EnableCors(enable bool) HttpGinOptions {
	return func(h *HttpGin) {
		h.enableCors = enable
	}
}

func EnableFlowControl(enable bool) HttpGinOptions {
	return func(h *HttpGin) {
		h.enableFlowControl = enable
	}
}

func EnableSession(enable bool, store sessions.Store) HttpGinOptions {
	return func(h *HttpGin) {
		h.enableSession = enable
		h.sessionStore = store
	}
}

func EnableHTML(enable bool, ptn string) HttpGinOptions {
	return func(h *HttpGin) {
		h.enableHTML = enable
		h.htmlPtn = ptn
	}
}

func EnableBlacklist(b []string) HttpGinOptions {
	return func(h *HttpGin) {
		for _, ip := range b {
			h.blacklist[ip] = true
		}
	}
}

func EnableWhitelist(b []string) HttpGinOptions {
	return func(h *HttpGin) {
		for _, ip := range b {
			h.whitelist[ip] = true
		}
	}
}

func NewHttpGinServer(addr string, options ...HttpGinOptions) *HttpGin {
	hg := &HttpGin{
		addr:      addr,
		blacklist: make(map[string]bool),
		whitelist: make(map[string]bool),
	}

	for _, opt := range options {
		opt(hg)
	}

	return hg
}

func (self *HttpGin) Start(fnSetupRouter func(*gin.Engine)) error {
	router := gin.Default()

	if self.enableDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	if self.enableCors {
		router.Use(cors.Default())
	}
	if self.enableSession {
		router.Use(sessions.Sessions("default", self.sessionStore))
	}

	if self.enableFlowControl {
		router.Use(http_flow_control_middleware())
	}

	if len(self.blacklist) > 0 {
		router.Use(http_blacklist_middleware(self.blacklist))
	}

	if len(self.whitelist) > 0 {
		router.Use(http_whitelist_middleware(self.whitelist))
	}

	if self.enableHTML {
		router.LoadHTMLGlob(self.htmlPtn)
	}

	self.server = &http.Server{
		Addr:    self.addr,
		Handler: router,
	}
	fnSetupRouter(router)

	go func() {
		if err := self.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
	return nil
}

func (self *HttpGin) Stop() error {
	ctx, canceller := context.WithTimeout(context.Background(), 5*time.Second)
	defer canceller()

	if err := self.server.Shutdown(ctx); err != nil {
		fmt.Println("HttpGin Stop error:", err)
	}
	select {
	case <-ctx.Done():
		fmt.Println("HttpGin quit")
	}
	return nil
}
