package api

import (
	"encoding/json"
	"fmt"
	"io"
	"mokku/constants"
	"net"
	"time"
)

type Data struct {
	Name       string
	Address    string
	Info       map[string]interface{}
	Version    map[string]interface{}
	Containers map[string]map[string]interface{}
}

type Client struct {
	Instance
	Name    string
	Address string
	Data    *Data
	docker  *DockerApi
}

func NewClient(settings *Settings) *Client {
	return &Client{
		Instance: Instance{settings: settings},
		Data: &Data{
			Info:       make(map[string]interface{}),
			Version:    make(map[string]interface{}),
			Containers: make(map[string]map[string]interface{}),
		},
	}
}

func (self *Client) collect() *Data {
	// collect client info
	self.Data.Name = self.Name
	self.Data.Address = self.Address

	// collect docker info
	info, err := self.docker.GetInfo()
	HandleError(err)
	self.Data.Info = (*info)

	// collect docker version
	version, err := self.docker.GetVersion()
	HandleError(err)
	self.Data.Version = (*version)

	// collect docker containers
	var cpuParser = NewCpuParser()
	var memoryParser = NewMemoryParser()

	c, err := self.docker.GetContainers()
	if err != nil {
		fmt.Println("Error getting containers", err)
		return nil
	}

	var container_ids = make([]string, 0)
	for _, cont := range *c {
		container_ids = append(container_ids, cont.(map[string]interface{})["Id"].(string))
	}

	// collect containers data
	for _, container_id := range container_ids {

		// initialise entry map
		container, exist := self.Data.Containers[container_id]
		if !exist {
			container = make(map[string]interface{})
			self.Data.Containers[container_id] = container
		}

		// collect cpu data
		cpu, err := cpuParser.Parse(&container_id)
		HandleError(err)
		container["cpu"] = cpu

		// collect memory data
		mem, err := memoryParser.Parse(&container_id)
		HandleError(err)
		container["memory"] = mem

	}
	return self.Data
}

func (self *Client) connect() (*net.TCPConn, error) {
	server, err := net.ResolveTCPAddr(
		constants.CONN_TYPE,
		fmt.Sprintf("%s:%d", self.settings.Host, self.settings.Port),
	)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP(constants.CONN_TYPE, nil, server)
	if err != nil {
		return nil, err
	}

	fmt.Println("connected to server")
	return conn, nil
}

func (self *Client) Start() {
	fmt.Println("starting client")

	var reconnects int = 0
	var max_reconnects int = 20
	var buf = make([]byte, 1024)

	self.docker = NewDockerApi(*self.settings.Docker_host)

	// connect
	conn, err := self.connect()
	if conn != nil {
		defer conn.Close()

		// set settings
		self.Address = conn.RemoteAddr().String()
		if self.Instance.settings.Name == "" {
			self.Instance.settings.Name = fmt.Sprintf("client:%s", self.Address)
		}
		self.Name = self.Instance.settings.Name

	}
	HandleError(err)

	for {

		// connected
		if conn != nil {

			// collect data
			stats, err := json.Marshal(self.collect())
			HandleError(err)

			// send data
			_, err = conn.Write([]byte(stats))
			HandleError(err)

			// read response
			length, err := conn.Read(buf)
			if err != nil {
				if err == io.EOF {
					fmt.Println("server disconnected")
				} else {
					HandleError(err)
				}
				conn.SetReadDeadline(time.Now().Add(time.Second * constants.CONN_TIMEOUT))
				conn = nil
				continue
			}

			// parse received data from client
			var data map[string]interface{}
			var byt = []byte(string(buf[:length]))
			err = json.Unmarshal(byt, &data)
			HandleError(err)

		} else {
			// try to reconnect
			conn, err = self.connect()
			if err != nil {
				reconnects += 1
				if reconnects >= max_reconnects {
					fmt.Println("exiting")
					break
				}
				if reconnects == 1 {
					fmt.Println("reconnecting")
				}
				fmt.Println(reconnects)
			} else {
				reconnects = 0
			}
		}

		// sleep
		time.Sleep(time.Millisecond * time.Duration(self.settings.Sleep*1000))
	}
}
