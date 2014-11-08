package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type DockerApi struct {
	docker_host url.URL
	client      *http.Client
}

func NewDockerApi(docker_host url.URL) *DockerApi {
	var endpoint = fmt.Sprintf("http://%s/_ping", docker_host.Host)
	var _, err = http.Get(endpoint)
	if err != nil {
		panic(err)
	}
	return &DockerApi{
		docker_host: docker_host,
		client:      &http.Client{},
	}
}

func (self *DockerApi) GetInfo() (*map[string]interface{}, error) {
	var endpoint = fmt.Sprintf("http://%s/info", self.docker_host.Host)
	return self.parseJsonMap(endpoint)
}

func (self *DockerApi) GetVersion() (*map[string]interface{}, error) {
	var endpoint = fmt.Sprintf("http://%s/version", self.docker_host.Host)
	return self.parseJsonMap(endpoint)
}

func (self *DockerApi) GetContainers() (*[]interface{}, error) {
	var endpoint = fmt.Sprintf("http://%s/containers/json", self.docker_host.Host)
	return self.parseJsonList(endpoint)
}

func (self *DockerApi) parseJsonMap(endpoint string) (*map[string]interface{}, error) {

	var result = make(map[string]interface{})
	bytez, err := self.request(endpoint)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(
		bytez,
		&result,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (self *DockerApi) parseJsonList(endpoint string) (*[]interface{}, error) {

	var result = make([]interface{}, 0)
	bytez, err := self.request(endpoint)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(
		bytez,
		&result,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (self *DockerApi) request(endpoint string) ([]byte, error) {

	var request, err = http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	response, err := self.client.Do(request)
	if err != nil {
		return nil, err
	}

	var bytez = make([]byte, 0)
	if response.Body != nil {
		defer response.Body.Close()

		bytez, err = ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
	}
	return bytez, nil
}
