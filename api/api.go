package api

import (
	"crypto/rand"
	"flag"
	"fmt"
	"mokku/constants"
	"net/url"
	"os"
	"time"
)

type Storage struct {
	Server  *StorageServer
	Clients map[string]*StorageClient
}

type StorageServer struct {
	Name    string
	Address string
}

type StorageClient struct {
	Name              string
	Address           string
	Date_last_updated time.Time
	Info              map[string]interface{}
	Version           map[string]interface{}
	Containers        map[string]*StorageContainer
}

type StorageContainer struct {
	Info  map[string]interface{}
	Stats []interface{}
}

func NewStorageClient() *StorageClient {
	return &StorageClient{
		Info:       make(map[string]interface{}),
		Version:    make(map[string]interface{}),
		Containers: make(map[string]*StorageContainer),
	}
}

func NewStorageContainer() *StorageContainer {
	return &StorageContainer{
		Info:  make(map[string]interface{}),
		Stats: make([]interface{}, 0),
	}
}

func NewStorage() *Storage {
	return &Storage{
		Server:  &StorageServer{},
		Clients: make(map[string]*StorageClient),
	}
}

func (self *Storage) Update(data *map[string]interface{}) error {
	var now = time.Now().UTC()

	var client_data = (*data)
	if client_data == nil {
		return nil
	}

	var client_name = (client_data["Name"]).(string)

	client, exist := self.Clients[client_name]
	if !exist {
		client = NewStorageClient()
		self.Clients[client_name] = client
	}
	client.Name = client_name
	client.Date_last_updated = now
	client.Info = (client_data["Info"]).(map[string]interface{})
	client.Version = (client_data["Version"]).(map[string]interface{})

	// containers data
	for container_id, container_data := range (client_data["Containers"]).(map[string]interface{}) {
		container, exist := client.Containers[container_id]
		if !exist {
			container = NewStorageContainer()
			client.Containers[container_id] = container
		}
		container.Info["date_last_updated"] = now
		// shift first element
		if len(container.Stats) >= constants.STORAGE_MAX_LENGTH {
			container.Stats = container.Stats[1:]
		}
		container.Stats = append(container.Stats, container_data)
	}
	return nil
}

type Instance struct {
	settings *Settings
}

type Settings struct {
	Docker_host *url.URL
	Type_of     int
	Name        string
	Sleep       float64
	Host        string
	Port        int
}

func ParseSettings() *Settings {
	var arg_is_server = flag.Bool("server", false, "explain instance type")
	var arg_name = flag.String("name", "", "explain name")
	var arg_sleep = flag.Float64("sleep", constants.DEFAULT_SLEEP, "explain sleep")
	var arg_host = flag.String("host", constants.SERVER_HOST, "explain host")
	var arg_port = flag.Int("port", constants.SERVER_PORT, "explain port")
	flag.Parse()

	var type_of = constants.TYPE_CLIENT
	if *arg_is_server {
		type_of = constants.TYPE_SERVER
	}
	var name = arg_name
	var sleep = arg_sleep
	var host = arg_host
	var port = arg_port

	var docker_host *url.URL
	if type_of == constants.TYPE_CLIENT {
		var env_docker_host = os.Getenv(constants.KEY_DOCKER_HOST)
		var dh, err = url.Parse(env_docker_host)
		if err != nil {
			panic(err)
		}
		if dh.Host == "" {
			panic(fmt.Sprintf("Invalid DOCKER_HOST format: '%s'", docker_host))
		}
		docker_host = dh
	}

	return &Settings{
		Docker_host: docker_host,
		Type_of:     type_of,
		Name:        (*name),
		Sleep:       (*sleep),
		Host:        (*host),
		Port:        (*port),
	}
}

func Uuid() (uuid string) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
