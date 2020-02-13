package dist

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

/**
TODO: 流控，debug，cors，(自动)黑名单
 */

type HttpGin struct {
	addr string
	routes map[string]func(*gin.Context)
	server *http.Server
	enableCors bool
	enableDebug bool
	enableFlowControl bool
	blacklist []string
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

func WithBlacklist(b []string) HttpGinOptions {
	return func(h *HttpGin) {
		h.blacklist = b
	}
}

func NewHttpGinServer(addr string, routes map[string]func(*gin.Context), options ...HttpGinOptions) *HttpGin {
	gin := &HttpGin{
		addr: addr,
		routes: routes,
	}

	for _, opt := range options {
		opt(gin)
	}

	return gin
}

func flow_control() gin.HandlerFunc {
	return func(c *gin.Context){

	}
}

func (self *HttpGin) Start() error {
	router := gin.Default()
	for path, cb := range self.routes {
		router.POST(path, cb)
	}

	if self.enableDebug {
		gin.SetMode(gin.DebugMode)
	}else{
		gin.SetMode(gin.ReleaseMode)
	}

	if self.enableCors {
		router.Use(cors.Default())
	}

	if self.enableFlowControl {
		router.Use(flow_control())
	}

	self.server = &http.Server{
		Addr: self.addr,
		Handler: router,
	}
	go func(){
		if err := self.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
	return nil
}

func (self *HttpGin) Stop() error {
	ctx, canceller := context.WithTimeout(context.Background(), 5 * time.Second)
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