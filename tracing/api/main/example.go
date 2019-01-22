package main

import (
	"fmt"
	"github.com/go-resty/resty"
	"os"
	"queen.com/x/tracing"

	"github.com/gin-gonic/gin"
	"queen.com/x/tracing/api"
)

func main() {
	r := gin.Default()
	hostName, err := os.Hostname()
	if err != nil {
		hostName = "unknown"
	}
	hostName = hostName + "29090"
	//192.168.1.52:16686
	_, closer, err := tracing.NewTracer(hostName, "192.168.1.52:6831")
	if err == nil {
		fmt.Println("Setting global tracer")
		defer closer.Close()
	} else {
		fmt.Println("Can't enable tracing: ", err.Error())
	}

	p := apitracing.ApiTracer([]byte("api-request-"))
	r.Use(p)

	r.GET("/", func(c *gin.Context) {
		r := resty.R()
		apitracing.InjectTraceID(tracing.GetSpanID(c), r.Header)
		req, err := r.Get("http://localhost:29091")
		if err != nil {
			// handle error
		}
		c.JSON(200, "FROM:"+string(req.Body()))
	})

	r.Run(":29090")
}
