package tests

import (
	jsonserver "github.com/last911/tools/server"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestJSONRequest(t *testing.T) {
	j := &jsonserver.JSONRequest{
		Method: "get",
		Params: map[string]interface{}{"Name": "scnjl"},
	}
	b, err := j.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(b))
}

func TestJSONResponse(t *testing.T) {
	j := &jsonserver.JSONResponse{
		Result: map[string]interface{}{"Name": "scnjl"},
	}
	b, err := j.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(b))

	j = &jsonserver.JSONResponse{
		Error: &jsonserver.JSONError{
			Code:    404,
			Message: "Not found",
		},
	}
	b, err = j.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(b))
}

func TestJSONRPCServer(t *testing.T) {
	go func() {
		server, err := jsonserver.NewJSONRPCServer("", "{}")
		if err != nil {
			t.Fatal(err)
		}

		server.AddHandle("/", func(c *jsonserver.Context) *jsonserver.JSONResponse {
			return &jsonserver.JSONResponse{
				Result: map[string]string{
					"Name": "last911",
				},
			}
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

	j := &jsonserver.JSONResponse{}
	if err := j.Unmarshal(b); err != nil {
		t.Fatal(err)
	}

	t.Log(j)
}
