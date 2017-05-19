package olderhc

import (
	"github.com/giskook/smarthome-interface/olderhc/pbgo"
	"log"
	"sync"
)

type HttpRouter struct {
	HttpRespones     map[string]chan *Report.ControlReport
	HttpRequestAdd   chan *HttpRequestPair
	HttpRequestDel   chan string
	HttpResponseChan chan *HttpResponsePair
}

var G_HttpRouter *HttpRouter = nil
var G_HttpRouterMutex sync.Mutex

func GetHttpRouter() *HttpRouter {
	defer G_HttpRouterMutex.Unlock()
	G_HttpRouterMutex.Lock()

	if G_HttpRouter == nil {
		G_HttpRouter = &HttpRouter{
			HttpRespones:     make(map[string]chan *Report.ControlReport),
			HttpRequestAdd:   make(chan *HttpRequestPair),
			HttpRequestDel:   make(chan string),
			HttpResponseChan: make(chan *HttpResponsePair),
		}
	}

	return G_HttpRouter
}

type HttpRequestPair struct {
	Key     string
	Command chan *Report.ControlReport
}

type HttpResponsePair struct {
	Key     string
	Command *Report.ControlReport
}

func (h *HttpRouter) AddRequest(req *HttpRequestPair) {
	h.HttpRequestAdd <- req
}

func (h *HttpRouter) DelRequest(key string) {
	h.HttpRequestDel <- key
}

func (h *HttpRouter) Run() {
	for {
		select {
		case add := <-h.HttpRequestAdd:
			h.HttpRespones[add.Key] = add.Command
		case key := <-h.HttpRequestDel:
			close(h.HttpRespones[key])
			delete(h.HttpRespones, key)
		case res := <-h.HttpResponseChan:
			chan_resp, ok := h.HttpRespones[res.Key]
			if ok {
				res.Command.SerialNumber = h.GetSerialID(uint16(res.Command.SerialNumber))
				chan_resp <- res.Command
			} else {
				log.Println("nokey")
			}

		}
	}
}

func (h *HttpRouter) DoResponse(resp *HttpResponsePair) {
	h.HttpResponseChan <- resp
}

func (h *HttpRouter) SendRequest(key string) chan *Report.ControlReport {
	chan_response := make(chan *Report.ControlReport)
	h.AddRequest(&HttpRequestPair{
		Key:     key,
		Command: chan_response,
	})

	return chan_response
}
