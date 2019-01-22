package main

import (
	"fmt"
	"os"
	"queen.com/m/merchant/character_address/pb"
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
	hostName = hostName + "29091"
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
		client := pb.New(c)
		response, err := client.Gets(&pb.CharacterAddressGetsRequest{
			AddressIds: []int64{2009},
			PageIndex:  1,
			PageSize:   10,
		})
		if err != nil {
			c.JSON(200, err.Error())
		} else {
			c.JSON(200, response)
		}
	})

	r.Run(":29091")
}
