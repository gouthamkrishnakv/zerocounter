package server

import (
	// stdlib imports
	"context"
	"fmt"
	"net"
	"net/http"

	// third-party imports
	"github.com/gouthamkrishnakv/zerocounter/logging"
	"github.com/rs/zerolog"
	"nhooyr.io/websocket"
)

// -- constants --
const DefaultPort = 4389

// -- variables

// -- struts --

// Server holds an HTTP server connection, required to connect to an HTTP
// server. It also helps us to manage life-cycles easier than the other
// ways.
type Server struct {
	s        http.Server
	wsConn   *websocket.Conn
	logger   zerolog.Logger
	readChan chan []byte
}

func NewServer() *Server {
	// setup logger
	newLogger := logging.L().
		// add timestamp, and stack printing, set "module" as database
		With().Str("module", "database").Logger()
	// create a serve mux
	serveMux := http.ServeMux{}
	// create the server
	newServer := &Server{
		s: http.Server{
			// TODO: Change to accept parameter
			Addr:    net.JoinHostPort("localhost", fmt.Sprintf("%d", DefaultPort)),
			Handler: &serveMux,
			// ReadTimeout:    10 * time.Second,
			// WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
		logger: newLogger,
	}

	newServer.logger.Info().Msg("adding routes")

	// set the method that the serveMux should handler
	serveMux.HandleFunc("/ws", newServer.Handle)

	// set the "/test" method also
	serveMux.HandleFunc("/test", testHandler)

	// return the created server
	return newServer
}

// TODO: remove it later on
func testHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, Test!"))
}

// Handle is the main request handler for WebSocket. All requests come here
// before it gets upgrated to WebSocket connection
func (s *Server) Handle(w http.ResponseWriter, r *http.Request) {
	s.logger.Info().Str("host", r.Host).Msg("handler connected")
	// TODO: all handler functions come here
	wsConn, wsOpenErr := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if wsOpenErr != nil {
		s.logger.Error().Err(wsOpenErr).Msg("websocket open error")
		return
		// w.WriteHeader(http.StatusInternalServerError)
		// w.Write([]byte("Websocket connection not accepted"))
	}

	// create context
	ctx := context.Background()

	// write "Hello, world"
	wsConn.Write(ctx, websocket.MessageText, []byte("Hello, World!"))
	s.wsConn = wsConn

	go s.readLoop(ctx)

	for msg := range s.readChan {
		if msg[0] == 0 {
			break
		}
		wsConn.Write(ctx, websocket.MessageBinary, msg)
	}

	s.wsConn.Close(websocket.StatusNormalClosure, "Thank you, client!")
}

func (s *Server) Serve() error {
	return s.s.ListenAndServe()
}

func (s *Server) readLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			// context is cancelled somewhere else
			return
		default:
			// Read the message
			msgType, data, readErr := s.wsConn.Read(ctx)
			if readErr != nil {
				// TODO: Maybe not the most graceful way to return an error right
				// here considering it's an error in a loop
				if websocket.CloseStatus(readErr) == websocket.StatusNormalClosure {
					s.logger.Info().Msg("connection closed")
					return
				}
				s.logger.Error().Err(readErr).Msg("abnormal read error")
			}
			s.readChan <- data
			// switch by message type
			// TODO: There should be more to do!!
			switch msgType {
			case websocket.MessageText:
				fmt.Printf("Message: %s", string(data))
			case websocket.MessageBinary:
				fmt.Printf("Bytes: %s", string(data))
			}
		}
	}
}
