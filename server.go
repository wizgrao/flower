package flower

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Server struct {
	Pixels     []*Pixel
	Clients    map[*Client]bool
	PixelChan  chan *Pixel
	ClientChan chan *Client
	CloseChan  chan *Client
}

func NewServer() *Server {
    upgrader.CheckOrigin = func(r *http.Request) bool {
      return true
    }
	s := &Server{
		Pixels:     make([]*Pixel, 0),
		Clients:    make(map[*Client]bool),
		PixelChan:  make(chan *Pixel, 10000),
		ClientChan: make(chan *Client, 10000),
		CloseChan:  make(chan *Client, 10000),
	}
	go func() {
		for {
			select {
			case c := <-s.ClientChan:
				s.Clients[c] = true
				c.send(s.Pixels)
			case c := <-s.CloseChan:
				delete(s.Clients, c)
			case p := <-s.PixelChan:
				s.Pixels = append(s.Pixels, p)
				for k := range s.Clients {
					k.sendOne(p)
				}
			}
		}
	}()
	return s
}

type Client struct {
	*Server
	*websocket.Conn
	pchan chan []*Pixel
}

func (c *Client) send(p []*Pixel) {
	c.pchan <- p
}

func (c *Client) sendOne(p *Pixel) {
	c.send([]*Pixel{p})
}

type Pixel struct {
	Color string
	I     int
	J     int
}

func (h *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Println("fuck", err)
		return
	}
	c := &Client{
		Server: h,
		Conn: conn,
		pchan: make(chan []*Pixel, 10000),
	}
	go func() {
		for {
			a := new(Pixel)
			if err := conn.ReadJSON(a); err == nil {
				c.Server.PixelChan <- a
			} else {
				log.Println(err)
				c.Server.CloseChan <- c
				conn.Close()
				return
			}
		}
	}()
	go func() {
		for {
			p, ok := <-c.pchan
			if !ok {
				err := conn.WriteMessage(websocket.CloseMessage, nil)
				if err != nil {
					log.Println(err)
					return
				}
				conn.Close()
				return
			}
			err := conn.WriteJSON(p)
			if err != nil {
				log.Println(err)
				c.Server.CloseChan <- c
				return
			}
		}
	}()
	h.ClientChan <- c
}
