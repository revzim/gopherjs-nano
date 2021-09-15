package gopherjs-nano

// revzim <https://github.com/revzim>
// nanojs | gopherjs-nano - gopherjs client wrapper for https://github.com/nano-ecosystem/nano-websocket-client

import "github.com/gopherjs/gopherjs/js"

type (
	NanoJS struct {
		*js.Object
	}
)

func New() *NanoJS {
	return &NanoJS{
		Object: js.Global.Get("nano"),
	}
}

// Init --
// cb IS WHERE TO INIT ALL CLIENT NANO FUNCTIONALITY
func (njs *NanoJS) Init(host string, port int, path string, cb func()) {
	njs.Call("init", map[string]interface{}{
		"host": host,
		"port": port,
		"path": path,
	}, cb)
}

// On --
func (njs *NanoJS) On(msgKey string, cb func(map[string]interface{})) {
	njs.Call("on", msgKey, cb)
}

func (njs *NanoJS) Request(reqKey string, data interface{}, cb func(data map[string]interface{})) {
	njs.Call("request", reqKey, data, cb)
}
