package main

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"os"
	"time"
)

func configEnv() {
	os.Setenv("PROJECT_NAME", "product")
	os.Setenv("LOGGING_SERVICE", "product")
}

func startTracing() {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Tags: []opentracing.Tag{
			{
				Key:   "server",
				Value: "product",
			},
			{
				Key:   "environment",
				Value: "prod",
			},
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
		},
	}

	tracer, closer, _ := cfg.New(
		"product",
		config.Logger(jaeger.StdLogger),
	)

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	someFunction()
}

func someFunction() {
	parent := opentracing.GlobalTracer().StartSpan("hello")
	defer parent.Finish()
	child := opentracing.GlobalTracer().StartSpan(
		"world", opentracing.ChildOf(parent.Context()))
	defer child.Finish()
}

func main() {
	configEnv()
	startTracing()

	r := gin.Default()

	r.GET("/product", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "response product",
		})
	})

	r.Run(":6001") // listen and serve on 0.0.0.0:8080
}