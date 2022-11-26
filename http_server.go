package memento

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MementoServerConfig struct {
	*MementoConfig
	Host string
	Port int
}

func RunServer(c *MementoServerConfig) error {
	var cache, err = NewMemento[string](c.MementoConfig)
	defer cache.Close()

	if err != nil {
		return err
	}

	r := gin.Default()
	apiV1 := r.Group("/api/v1")

	apiV1.GET("/cache/:key", func(c *gin.Context) {
		key := c.Param("key")
		status := http.StatusOK

		v, ok := cache.Get(key)
		if !ok {
			status = http.StatusNotFound
		}

		c.Data(status, "text/plain", v)
	})

	apiV1.PUT("/cache/:key/:value", func(c *gin.Context) {
		key := c.Param("key")
		value := c.Param("value")

		go cache.Set(key, []byte(value))

		c.Data(http.StatusNoContent, "text/plain", nil)
	})

	apiV1.DELETE("/cache/:key", func(c *gin.Context) {
		key := c.Param("key")

		go cache.Delete(key)

		c.Data(http.StatusNoContent, "text/plain", nil)
	})

	return r.Run(fmt.Sprintf("%s:%d", c.Host, c.Port))
}
