package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	. "./models"
)

type ParentID string

var groups map[string]string

type Swagger struct{}

func (s Swagger) Generate(insomniaFile string, configFile string) {
	raw, err := ioutil.ReadFile(insomniaFile)
	if err != nil {
		log.Fatal(err)
	}
	var insomnia Insomnia
	if err := json.Unmarshal(raw, &insomnia); err != nil {
		log.Fatal(err)
	}

	groups = make(map[string]string)
	list := make(map[string]map[string][]Resource)
	for _, resource := range insomnia.Resources {
		if resource.Type == "request_group" {
			groups[resource.ID] = resource.Name
			list[resource.ID] = make(map[string][]Resource, 0)
		}
		if resource.Type == "request" {
			fetchVariables(&resource)
			list[resource.ParentID][resource.URL] = append(list[resource.ParentID][resource.URL], resource)
		}
	}

	GenerateSwagger(list, configFile)
}

func GenerateSwagger(data map[string]map[string][]Resource, configFileName string) {
	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
		"RemovePathPrefix": func(path string) string {
			re := regexp.MustCompile("{{(.*?)}}")
			for _, param := range re.FindAllStringSubmatch(path, -1) {
				path = strings.Replace(path, param[0], "", -1)
			}
			return path
		},
		"GetGroupName": func(id string) string {
			return groups[id]
		},
	}
	tmpl, err := template.New("swagger.yaml").Funcs(funcMap).ParseFiles("tmpl/swagger.yaml")
	if err != nil {
		fmt.Println(err)
	}
	f, err := os.Create("swagger.yaml")
	if err != nil {
		fmt.Println(err)
	}
	config := ParseConfig(configFileName)
	res := struct {
		Config   ConfigSwagger
		Entities map[string]map[string][]Resource
	}{
		config,
		data,
	}
	err = tmpl.Execute(f, res)
	if err != nil {
		fmt.Println(err)
	}
}

func ParseConfig(configFileName string) ConfigSwagger {
	raw, err := ioutil.ReadFile(configFileName)
	if err != nil {
		log.Fatal(err)
	}
	var config ConfigSwagger
	if err := json.Unmarshal(raw, &config); err != nil {
		log.Fatal(err)
	}
	return config
}

func fetchVariables(resource *Resource) {
	re := regexp.MustCompile("/{(.*?)}#*")
	for _, param := range re.FindAllStringSubmatch(resource.URL, -1) {
		resource.InsomniaParams = append(resource.InsomniaParams, param[1])
	}
}
