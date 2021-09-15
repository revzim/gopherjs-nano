package nanojs

// revzim <https://github.com/revzim>
// nanojs | gopherjs-nano - gopherjs client wrapper for https://github.com/nano-ecosystem/nano-websocket-client

import "github.com/gopherjs/gopherjs/js"

type (
	NanoJS struct {
		*js.Object
	}

	Opts struct {
		Host         string // `js:"host"`
		Port         int    // `js:"port"`
		Path         string // `js:"path"`
		CallbackFunc func() // `js:"cb"`
	}
)

func New() *NanoJS {
	return &NanoJS{
		Object: js.Global.Get("nano"),
	}
}

// Init --
// cb IS WHERE TO INIT ALL CLIENT NANO FUNCTIONALITY
func (njs *NanoJS) Init(opts *Opts) {
	njs.Call("init", map[string]interface{}{
		"host": opts.Host,
		"port": opts.Port,
		"path": opts.Path,
	}, opts.CallbackFunc)
}

// On --
func (njs *NanoJS) On(msgKey string, cb func(map[string]interface{})) {
	njs.Call("on", msgKey, cb)
}

// Request --
func (njs *NanoJS) Request(reqKey string, data interface{}, cb func(data map[string]interface{})) {
	njs.Call("request", reqKey, data, cb)
}
