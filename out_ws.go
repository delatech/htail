package main

import (
	"encoding/json"
	"sync"

	"golang.org/x/net/websocket"
)

type websocketHub struct {
	conns []*websocket.Conn
	rw    *sync.RWMutex
}

func newWebsocketHub() *websocketHub {
	return &websocketHub{
		conns: make([]*websocket.Conn, 0),
		rw:    &sync.RWMutex{},
	}
}

func (h *websocketHub) addConn(conn *websocket.Conn) {
	h.rw.Lock()
	defer h.rw.Unlock()

	h.conns = append(h.conns, conn)
}

func (h *websocketHub) WriteLine(l Line) error {
	j, err := json.Marshal(l)
	if err != nil {
		return err
	}
	_, err = h.Write(j)
	return err
}

// Write implements the io.Writer interface
func (h *websocketHub) Write(d []byte) (n int, err error) {
	h.rw.RLock()
	defer h.rw.RUnlock()

	for i := range h.conns {
		c, e := h.conns[i].Write(d)
		n += c
		if e != nil {
			err = e
		}
	}
	return
}
