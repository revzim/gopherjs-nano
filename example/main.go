// https://github.com/lonng/nano/tree/master/examples/demo/chat
package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lonng/nano"
	"github.com/lonng/nano/component"
	"github.com/lonng/nano/pipeline"
	"github.com/lonng/nano/scheduler"
	"github.com/lonng/nano/serialize/json"
	"github.com/lonng/nano/session"
)

type (
	Room struct {
		group *nano.Group
	}

	// RoomManager represents a component that contains a bundle of room
	RoomManager struct {
		component.Base
		timer *scheduler.Timer
		rooms map[int]*Room
	}

	// UserMessage represents a message that user sent
	UserMessage struct {
		Name    string `json:"name"`
		Content string `json:"content"`
	}

	// NewUser message will be received when new user join room
	NewUser struct {
		Content string `json:"content"`
	}

	// AllMembers contains all members uid
	AllMembers struct {
		Members []int64 `json:"members"`
	}

	// JoinResponse represents the result of joining room
	JoinResponse struct {
		Code     int    `json:"code"`
		Result   string `json:"result"`
		Username string `json:"username"`
	}

	stats struct {
		component.Base
		timer         *scheduler.Timer
		outboundBytes int
		inboundBytes  int
	}
)

func (stats *stats) outbound(s *session.Session, msg *pipeline.Message) error {
	stats.outboundBytes += len(msg.Data)
	return nil
}

func (stats *stats) inbound(s *session.Session, msg *pipeline.Message) error {
	stats.inboundBytes += len(msg.Data)
	return nil
}

func (stats *stats) AfterInit() {
	stats.timer = scheduler.NewTimer(time.Minute, func() {
		println("OutboundBytes", stats.outboundBytes)
		println("InboundBytes", stats.outboundBytes)
	})
}

const (
	testRoomID = 1
	roomIDKey  = "ROOM_ID"
)

func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: map[int]*Room{},
	}
}

// AfterInit component lifetime callback
func (mgr *RoomManager) AfterInit() {
	session.Lifetime.OnClosed(func(s *session.Session) {
		if !s.HasKey(roomIDKey) {
			return
		}
		room := s.Value(roomIDKey).(*Room)
		room.group.Leave(s)
	})
	mgr.timer = scheduler.NewTimer(time.Minute, func() {
		for roomId, room := range mgr.rooms {
			println(fmt.Sprintf("UserCount: RoomID=%d, Time=%s, Count=%d",
				roomId, time.Now().String(), room.group.Count()))
		}
	})
}

// Join room
func (mgr *RoomManager) Join(s *session.Session, msg []byte) error {
	// NOTE: join test room only in demo
	log.Println(msg)
	room, found := mgr.rooms[testRoomID]
	if !found {
		room = &Room{
			group: nano.NewGroup(fmt.Sprintf("room-%d", testRoomID)),
		}
		mgr.rooms[testRoomID] = room
	}

	fakeUID := s.ID() //just use s.ID as uid !!!
	uid := uuid.New().String()[:6]
	s.Bind(fakeUID) // binding session uids.Set(roomIDKey, room)
	s.Set(roomIDKey, room)
	s.Set(fmt.Sprintf("%d", fakeUID), uid)
	s.Push("onMembers", &AllMembers{Members: room.group.Members()})
	// notify others
	room.group.Broadcast("onNewUser", &NewUser{Content: fmt.Sprintf("New user: %s", uid)})
	// new user join group
	room.group.Add(s) // add session to group
	return s.Response(&JoinResponse{Result: "success", Username: uid})
}

// Message sync last message to all members
func (mgr *RoomManager) Message(s *session.Session, msg *UserMessage) error {
	if !s.HasKey(roomIDKey) {
		return fmt.Errorf("not join room yet")
	}
	room := s.Value(roomIDKey).(*Room)
	return room.group.Broadcast("onMessage", msg)
}

func initGameServer(port int) {
	components := &component.Components{}
	components.Register(
		NewRoomManager(),
		component.WithName("room"), // rewrite component and handler name
		component.WithNameFunc(strings.ToLower),
	)

	// traffic stats
	pip := pipeline.New()
	var stats = &stats{}
	pip.Outbound().PushBack(stats.outbound)
	pip.Inbound().PushBack(stats.inbound)

	log.SetFlags(log.LstdFlags | log.Llongfile)
	// http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))

	nano.Listen(fmt.Sprintf(":%d", port),
		nano.WithIsWebsocket(true),
		nano.WithPipeline(pip),
		nano.WithCheckOriginFunc(func(_ *http.Request) bool { return true }),
		nano.WithWSPath("/ws"),
		nano.WithDebugMode(),
		nano.WithSerializer(json.NewSerializer()), // override default serializer
		nano.WithComponents(components),
	)
}

// WEBSERVER FOR EXAMPLE
func initWebServer(port int) {
	e := echo.New()

	e.Debug = true
	const publicPath = "./public"

	e.Static("", publicPath)

	e.Use(middleware.Logger())

	e.Use(middleware.Recover())

	DefaultCORSConfig := middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"http://127.0.0.1", "http://192.168.1.159", "*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}

	e.Use(middleware.CORSWithConfig(DefaultCORSConfig))
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))

}

func main() {

	go initWebServer(8080)

	initGameServer(8081)
}
