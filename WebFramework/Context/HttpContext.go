package Context

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/yoyofx/yoyogo/DependencyInjection"
	"github.com/yoyofx/yoyogo/WebFramework/ActionResult"
	"net/http"
	"strings"
	"sync"
)

const (
	defaultTagName   = "param"
	jsonTagName      = "json"
	defaultMaxMemory = 32 << 20 // 32 MB

)

type M = map[string]string

type HttpContext struct {
	Input            Input
	Output           Output
	RequiredServices DependencyInjection.IServiceProvider
	store            map[string]interface{}
	storeMutex       *sync.RWMutex
	Result           interface{}
}

func NewContext(w http.ResponseWriter, r *http.Request, sp DependencyInjection.IServiceProvider) *HttpContext {
	ctx := &HttpContext{}
	ctx.init(w, r, sp)
	return ctx
}

func (ctx *HttpContext) init(w http.ResponseWriter, r *http.Request, sp DependencyInjection.IServiceProvider) {
	ctx.storeMutex = new(sync.RWMutex)
	ctx.Input = NewInput(r, 20<<32)
	ctx.Output = Output{Response: &responseWriter{w, 0, 0, nil}}
	ctx.RequiredServices = sp
	ctx.storeMutex.Lock()
	ctx.store = nil
	ctx.storeMutex.Unlock()
}

//Set data in context.
func (ctx *HttpContext) SetItem(key string, val interface{}) {
	ctx.storeMutex.Lock()
	if ctx.store == nil {
		ctx.store = make(map[string]interface{})
	}
	ctx.store[key] = val
	ctx.storeMutex.Unlock()
}

// Get data in context.
func (ctx *HttpContext) GetItem(key string) interface{} {
	ctx.storeMutex.RLock()
	v := ctx.store[key]
	ctx.storeMutex.RUnlock()
	return v
}

func (ctx *HttpContext) Bind(i interface{}) (err error) {
	req := ctx.Input.Request
	contentType := req.Header.Get(HeaderContentType)
	if req.Body == nil {
		err = errors.New("request body can't be empty")
		return err
	}
	err = errors.New("request unsupported MediaType -> " + contentType)
	tagName := defaultTagName
	switch {
	case strings.HasPrefix(contentType, MIMEApplicationXML):
		err = xml.Unmarshal(ctx.Input.FormBody, i)
	case strings.HasPrefix(contentType, MIMEApplicationJSON):
		err = json.Unmarshal(ctx.Input.FormBody, i)
	default:
	}
	err = ConvertMapToStruct(tagName, i, ctx.Input.GetAllParam())
	return err
}

// Redirect redirects the request
func (ctx *HttpContext) Redirect(code int, url string) {
	http.Redirect(ctx.Output.GetWriter(), ctx.Input.GetReader(), url, code)
}

// ActionResult writes the response headers and calls render.ActionResult to render data.
func (ctx *HttpContext) Render(code int, r ActionResult.IActionResult) {

	if !bodyAllowedForStatus(code) {
		r.WriteContentType(ctx.Output.GetWriter())
		ctx.Output.SetStatusCodeNow()
		return
	}

	if err := r.Render(ctx.Output.GetWriter()); err != nil {
		panic(err)
	}

	ctx.Output.SetStatusCode(code)
}
