package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func setupHttp() *gin.Engine {
	// setup
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// middleware
	r.Use(
		otelgin.Middleware(serviceName),
	)

	setupEndpoints(r)

	return r
}

func setupEndpoints(r *gin.Engine) {
	// endpoints
	r.POST("/ping", func(c *gin.Context) {
		type body struct {
			Hosts []string `json:"hosts"`
		}

		var b body

		err := c.ShouldBind(&b)
		if err != nil {
			logger.ErrorfC(c.Request.Context(), "error", zap.Error(err))
		}

		target := b.Hosts[0]
		b.Hosts = b.Hosts[1:]

		var (
			req *http.Request
		)
		if len(b.Hosts) == 0 {
			req, _ = http.NewRequestWithContext(c.Request.Context(), http.MethodGet, fmt.Sprintf("%s/ping", target), http.NoBody)
		} else {
			rawb, _ := json.Marshal(b)
			logger.InfofC(c.Request.Context(), "posting ping: '%s' with %s", target, string(rawb))
			req, _ = http.NewRequestWithContext(c.Request.Context(), http.MethodPost, fmt.Sprintf("%s/ping", target), bytes.NewBuffer(rawb))
			req.Header.Set("Content-Type", "application/json")
		}

		defaultClient := &http.Client{
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		}

		resp, err := defaultClient.Do(req)
		if resp != nil {
			defer resp.Body.Close()
		}
		if err != nil {
			logger.ErrorfC(c.Request.Context(), "error", zap.Error(err))
			c.Status(500)
			return
		}

		bodyResBytes, _ := io.ReadAll(resp.Body)

		var bodyRes map[string]interface{}
		json.Unmarshal(bodyResBytes, &bodyRes)

		c.JSON(200, bodyRes)
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":  "pong",
			"trace_id": trace.SpanFromContext(c.Request.Context()).SpanContext().TraceID().String(),
		})
	})
}
