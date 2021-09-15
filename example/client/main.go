package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gopherjs/gopherjs/js"
	nanojs "github.com/revzim/gopherjs-nano"
	vue "github.com/revzim/gopherjs-vue"
	"honnef.co/go/js/dom"
)

type (
	Model struct {
		*js.Object
		MountEl      string                   `js:"mountEl"`
		SendBtnLabel string                   `js:"sendBtnLabel"`
		Username     string                   `js:"username"`
		InputMessage string                   `js:"inputMsg"`
		Messages     []map[string]interface{} `js:"messages"`
	}
)

const (
	VueAppMountElement = "#app"
	NanoPort           = 8081
	NanoPath           = "/ws"
	SendBtnLabel       = "SEND"
)

var ()

func (m *Model) SendMessage() {
	log.Println("sending msg:", m.InputMessage)
	nanoJS := js.Global.Get("nano")
	nanoJS.Call("notify", "room.message", map[string]interface{}{
		"name":    m.Username,
		"content": m.InputMessage,
	})
	m.InputMessage = ""
}

func (m *Model) AddMessage(data map[string]interface{}) {
	m.Messages = append(m.Messages, data)
	chatBoxElem := dom.GetWindow().Document().QuerySelector("#chat-box").Underlying()
	dom.GetWindow().SetTimeout(func() {
		chatBoxElem.Set("scrollTop", chatBoxElem.Get("scrollHeight"))
	}, 100)
}

func InitNanoClient(m *Model, port int, gsPath string) *js.Object {
	loc := js.Global.Get("window").Get("location")
	nanoJS := nanojs.New()

	nanoJSCallback := func() {
		log.Println("nano init")
		onNewUser := func(data map[string]interface{}) {
			msgData := map[string]interface{}{
				"name":    "system",
				"content": data["content"].(string),
			}
			m.AddMessage(msgData)
		}
		onMembers := func(data map[string]interface{}) {
			var content string
			if data["members"] != nil {
				content = fmt.Sprintf("active members: %v", data["members"].([]interface{}))
			} else {
				content = "welcome"
			}
			msgData := map[string]interface{}{
				"name":    "system",
				"content": content,
			}
			m.AddMessage(msgData)
		}

		onJoin := func(data map[string]interface{}) {
			dataCode := data["code"].(float64)
			if dataCode == 0 {
				m.Username = data["username"].(string)
				msgData := map[string]interface{}{
					"name":    "system",
					"content": data["result"].(string),
				}
				m.AddMessage(msgData)
				nanoJS.On("onMessage", m.AddMessage)
			}
		}

		nanoJS.On("onNewUser", onNewUser)
		nanoJS.On("onMembers", onMembers)
		nanoJS.Request("room.join", nil, onJoin)
	}

	nanoOpts := &nanojs.Opts{
		Host:         loc.Get("hostname").String(),
		Port:         port,
		Path:         gsPath,
		CallbackFunc: nanoJSCallback,
	}
	nanoJS.Init(nanoOpts)

	return nanoJS.Object
}

func InitVuetify() *js.Object {
	Vuetify := js.Global.Get("Vuetify")

	vue.Use(Vuetify)

	vuetifyCFG := js.Global.Get("Object").New()
	vuetifyCFG.Set("theme", map[string]interface{}{
		"dark": true,
	})

	vtfy := Vuetify.New(vuetifyCFG)

	return vtfy
}

func InitVueOpts(m *Model) *vue.Option {
	o := vue.NewOption()

	o.SetDataWithMethods(m)

	o.AddComputed("test", func(vm *vue.ViewModel) interface{} {
		return strings.ToUpper(vm.Data.Get("test").String())
	})

	o = o.Mixin(js.M{
		"vuetify": InitVuetify(),
		"nano":    InitNanoClient(m, NanoPort, NanoPath),
	})

	return o
}

func main() {
	m := &Model{
		Object: js.Global.Get("Object").New(),
	}

	m.MountEl = VueAppMountElement
	m.SendBtnLabel = SendBtnLabel
	m.Messages = []map[string]interface{}{}
	m.Username = ""
	m.InputMessage = ""

	o := InitVueOpts(m)

	v := o.NewViewModel()

	v.Mount(VueAppMountElement)

}
