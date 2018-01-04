package tests

import (
	"github.com/gin-gonic/gin"
	"github.com/last911/tools/server"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestGinServer(t *testing.T) {
	go func() {
		server, err := server.NewGinServer("dev", "{}")
		if err != nil {
			t.Fatal(err)
		}
		server.Engine.GET("/", func(c *gin.Context) {
			c.JSON(200, map[string]interface{}{
				"name": "scnjl",
				"age":  38,
			})
		})

		if err := server.Run("127.0.0.1:8008"); err != nil {
			t.Fatal(err)
		}
	}()

	res, err := http.Get("http://127.0.0.1:8008/")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(b))
}
