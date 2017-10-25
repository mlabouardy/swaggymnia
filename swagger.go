package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/ghodss/yaml"

	. "./models"
)

type ParentID string

var groups map[string]string

const (
	REQUEST_GROUP = "request_group"
	REQUEST       = "request"
)

type Swagger struct {
	Config   ConfigSwagger
	Entities map[string]map[string][]Resource
}

func (c Swagger) parse(insomnia Insomnia) map[string]map[string][]Resource {
	groups = make(map[string]string)
	list := make(map[string]map[string][]Resource)
	for _, resource := range insomnia.Resources {
		if resource.Type == REQUEST_GROUP {
			groups[resource.ID] = resource.Name
			list[resource.ID] = make(map[string][]Resource, 0)
		}
		if resource.Type == REQUEST {
			fetchVariables(&resource)
			list[resource.ParentID][resource.URL] = append(list[resource.ParentID][resource.URL], resource)
		}
	}
	return list
}

func (c Swagger) readInsomniaExport(fileName string) Insomnia {
	raw, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	var insomnia Insomnia
	if err := json.Unmarshal(raw, &insomnia); err != nil {
		log.Fatal(err)
	}
	return insomnia
}

func (c Swagger) readSwaggerConfig(fileName string) ConfigSwagger {
	raw, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	var config ConfigSwagger
	if err := json.Unmarshal(raw, &config); err != nil {
		log.Fatal(err)
	}
	return config
}

func (s Swagger) Generate(insomniaFile string, configFile string, output string) {
	insomnia := s.readInsomniaExport(insomniaFile)
	config := s.readSwaggerConfig(configFile)
	entities := s.parse(insomnia)
	switch output {
	case "yaml":
		s.generateYAML(entities, config)
		break
	case "json":
		s.generateJSON(entities, config)
		break
	default:
		log.Fatal("Format isn't supported !")
	}
}

func (s Swagger) generateJSON(entities map[string]map[string][]Resource, config ConfigSwagger) {
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

	res := struct {
		Config   ConfigSwagger
		Entities map[string]map[string][]Resource
	}{
		config,
		entities,
	}

	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, res)
	if err != nil {
		fmt.Println(err)
	}
	data, err := yaml.YAMLToJSON(tpl.Bytes())
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile("swagger.json", data, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func (s Swagger) generateYAML(entities map[string]map[string][]Resource, config ConfigSwagger) {
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

	data, err := Asset("tmpl/swagger.yaml")
	if err != nil {
		fmt.Println(err)
	}
	tmpl, err := template.New("swagger.yaml").Funcs(funcMap).Parse(string(data))
	if err != nil {
		fmt.Println(err)
	}
	f, err := os.Create("swagger.yaml")
	if err != nil {
		fmt.Println(err)
	}

	res := struct {
		Config   ConfigSwagger
		Entities map[string]map[string][]Resource
	}{
		config,
		entities,
	}
	err = tmpl.Execute(f, res)
	if err != nil {
		fmt.Println(err)
	}
}

func fetchVariables(resource *Resource) {
	re := regexp.MustCompile("/{(.*?)}#*")
	for _, param := range re.FindAllStringSubmatch(resource.URL, -1) {
		resource.InsomniaParams = append(resource.InsomniaParams, param[1])
	}
}
