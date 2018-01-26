package server

import (
	"encoding/json"
	"github.com/bitly/go-simplejson"
	log "github.com/last911/tools/log"
	"io/ioutil"
	"net/http"
)

// HandleFunc route handle function
type HandleFunc func(*Context) *JSONResponse

// jsonObj is JSONRequest and JSONResponse parent
type jsonObj struct {
	ID      int    `json:"id"`
	JSONRPC string `json:"jsonrpc"`
}

func (j *jsonObj) Marshal() ([]byte, error) {
	j.JSONRPC = "2.0"
	return json.Marshal(j)
}

func (j *jsonObj) Unmarshal(b []byte) error {
	return json.Unmarshal(b, j)
}

// JSONRequest request struct
type JSONRequest struct {
	jsonObj
	Method string      `json:"method"`
	Params interface{} `json:"params,omitempty"`
}

// JSONResponse response struct
type JSONResponse struct {
	jsonObj
	Result interface{} `json:"result,omitempty"`
	Error  *JSONError  `json:"error,omitempty"`
}

// JSONRPCError error struct
type JSONRPCError struct {
	Code    int
	Err     error
	Message string
}

// Error return JSONRPCError message string
func (j *JSONRPCError) Error() string {
	return "JSONRPCError Code: " + string(j.Code) + " Error: [" + j.Err.Error() + "] message: [" + j.Message + "]"
}

// JSONError error of response struct
type JSONError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Marshal JSONError to []byte
func (j *JSONError) Marshal() ([]byte, error) {
	return json.Marshal(j)
}

// Unmarshal []byte to JSONError
func (j *JSONError) Unmarshal(b []byte) error {
	return json.Unmarshal(b, j)
}

// Context web request and response context
type Context struct {
	Response http.ResponseWriter
	Request  *http.Request

	RequestData  *JSONRequest
	ResponseData *JSONResponse
}

// NewContext return context
func NewContext(w http.ResponseWriter, r *http.Request) (*Context, error) {
	body, err := r.GetBody()
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	reqData := &JSONRequest{}
	if err := reqData.Unmarshal(b); err != nil {
		return nil, err
	}

	jsonResponse := &JSONResponse{}
	jsonResponse.ID = reqData.ID
	jsonResponse.JSONRPC = "2.0"
	return &Context{
		Response:     w,
		Request:      r,
		RequestData:  reqData,
		ResponseData: jsonResponse,
	}, nil
}

// JSONRPCServer jsonrpc server struct
type JSONRPCServer struct {
	Server
	handles     map[string]HandleFunc // handles
	middlewares []HandleFunc          // middleware
	error404    *JSONResponse
	error500    *JSONResponse
}

// NewJSONRPCServer jsonrpc protocol web server
func NewJSONRPCServer(env string, config ...*simplejson.Json) (*JSONRPCServer, error) {
	server := &JSONRPCServer{handles: make(map[string]HandleFunc)}
	server.Env = env
	var err error

	server.AppPath, server.Config, err = initialize(env, config...)
	if err != nil {
		return nil, err
	}

	// init error page
	server.error404 = &JSONResponse{
		Error: &JSONError{
			Code:    60404,
			Message: "Not found method",
		},
	}

	server.error500 = &JSONResponse{
		Error: &JSONError{
			Code:    60500,
			Message: "Service unavailable",
		},
	}

	return server, nil
}

// Error404 set error 404 page
func (s *JSONRPCServer) Error404(page404 *JSONResponse) {
	s.error404 = page404
}

// Error500 set error 500 page
func (s *JSONRPCServer) Error500(page500 *JSONResponse) {
	s.error500 = page500
}

// Use add middleware
func (s *JSONRPCServer) Use(handles ...HandleFunc) {
	s.middlewares = append(s.middlewares, handles...)
}

// AddHandle add handle path function
func (s *JSONRPCServer) AddHandle(method string, handle HandleFunc) {
	s.handles[method] = handle
}

// ServeHTTP
func (s *JSONRPCServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			respData := s.error500
			switch t := err.(type) {
			case JSONRPCError:
				respData.Error.Code = t.Code
				respData.Error.Message = t.Message
				log.Debug("request error:", err)
			default:
				log.Error("request error:", err)
			}
			b, _ := respData.Marshal()
			w.Write(b)
		}
	}()
	r.ParseForm()
	c, err := NewContext(w, r)

	if err != nil {
		panic(JSONRPCError{Code: 501, Err: err, Message: "Context error"})
	}

	var respData *JSONResponse
	// handle middleware
	for _, handle := range s.middlewares {
		respData = handle(c)
		if respData != nil {
			break
		}
	}

	if respData == nil {
		handleFunc := s.handles[c.RequestData.Method]
		if handleFunc != nil {
			respData = handleFunc(c)
		} else {
			log.Debug("request method[%s] not found", c.RequestData.Method)
			respData = s.error404
		}
	}

	b, err := respData.Marshal()
	if err != nil {
		panic(JSONRPCError{Code: 502, Err: err, Message: "Response decode error"})
	}
	w.Write(b)
}

// Run start JSONRPCServer run
func (s *JSONRPCServer) Run(addr ...string) error {
	var address string
	if len(addr) == 0 {
		address = s.Config.Get("app").Get("addr").MustString()
	} else {
		address = addr[0]
	}
	err := http.ListenAndServe(address, s)
	if err != nil {
		return err
	}

	return nil
}
