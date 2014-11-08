package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mokku/constants"
	"net"
	"time"
)

type Server struct {
	Instance
	storage *Storage
}

func NewServer(settings *Settings, storage *Storage) *Server {
	return &Server{
		Instance: Instance{settings: settings},
		storage:  storage,
	}
}

func (self *Server) Start() {
	fmt.Println("starting mokku server")

	// validate address
	server, err := net.ResolveTCPAddr(constants.CONN_TYPE, fmt.Sprintf(":%d", self.settings.Port))
	HandleErrorAndExit(err)

	// start to listen
	ln, err := net.ListenTCP(constants.CONN_TYPE, server)
	HandleError(err)

	for {
		conn, err := ln.Accept()
		if err != nil {
			HandleError(err)
			continue
		}

		if conn != nil {
			fmt.Println(fmt.Sprintf("client connected (%s)", conn.RemoteAddr()))
			defer conn.Close()

			// set settings
			self.storage.Server.Address = conn.RemoteAddr().String()
			if self.Instance.settings.Name == "" {
				self.Instance.settings.Name = fmt.Sprintf("server:%s", self.storage.Server.Address)
			}
			self.storage.Server.Name = self.settings.Name

			go self.handleConnection(conn)
		}
	}
}

func (self *Server) handleConnection(conn net.Conn) {

	for {
		// read from client
		buffer := &bytes.Buffer{}
		for {
			buf := make([]byte, constants.SOCKET_BUFFER_SIZE)
			size, err := conn.Read(buf)
			if err == io.EOF {
				fmt.Println(fmt.Sprintf("client disconnected (%s)", conn.RemoteAddr()))
				break
			}
			HandleError(err)

			buffer.Write(buf[:size])
			if size < constants.SOCKET_BUFFER_SIZE {
				break
			}
		}

		// parse received data from client
		var data = make(map[string]interface{})
		err := json.Unmarshal(buffer.Bytes(), &data)
		if err != nil {
			fmt.Println("error: parsing data from client")
		} else {
			// push to storage
			self.storage.Update(&data)
		}

		// answer to client
		data = make(map[string]interface{})
		data["status"] = "OK"
		data["sleep"] = constants.DEFAULT_SLEEP

		response, err := json.Marshal(data)
		HandleError(err)

		_, err = conn.Write([]byte(response))
		HandleError(err)
		if err != nil {
			break
		}
	}
	conn.SetReadDeadline(time.Now().Add(time.Second * constants.CONN_TIMEOUT))
}
