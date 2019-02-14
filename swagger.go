package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/ghodss/yaml"

	. "github.com/mlabouardy/swaggymnia/models"
)

var groupNames map[string]string

const (
	REQUEST_GROUP = "request_group"
	REQUEST       = "request"
	YAML_FORMAT   = "yaml"
	JSON_FORMAT   = "json"
)

type Swagger struct {
	Config   SwaggerConfig
	Entities map[string]map[string][]Resource
}

type SwaggerConfig struct {
	Title       string `json:"title"`
	Version     string `json:"version"`
	Host        string `json:"host"`
	BasePath    string `json:"basePath"`
	Schemes     string `json:"schemes"`
	Description string `json:"description"`
}

func parse(insomnia Insomnia) map[string]map[string][]Resource {
	groupNames = make(map[string]string)
	entities := make(map[string]map[string][]Resource)
	for _, resource := range insomnia.Resources {
		if resource.Type == REQUEST_GROUP {
			groupNames[resource.ID] = resource.Name
			entities[resource.ID] = make(map[string][]Resource, 0)
		}
		if resource.Type == REQUEST {
			fetchVariables(&resource)
			if entities[resource.ParentID] == nil {
				entities[resource.ParentID] = make(map[string][]Resource, 0)
			}
			entities[resource.ParentID][resource.URL] = append(entities[resource.ParentID][resource.URL], resource)
		}
	}
	return entities
}

func readInsomniaExport(fileName string) Insomnia {
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

func readSwaggerConfig(fileName string) SwaggerConfig {
	raw, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	var config SwaggerConfig
	if err := json.Unmarshal(raw, &config); err != nil {
		log.Fatal(err)
	}
	return config
}

func (s *Swagger) Generate(insomniaFile string, configFile string, outputFormat string) {
	s.Config = readSwaggerConfig(configFile)
	s.Entities = parse(readInsomniaExport(insomniaFile))
	switch outputFormat {
	case YAML_FORMAT:
		s.generateYAML()
		break
	case JSON_FORMAT:
		s.generateJSON()
		break
	default:
		log.Fatal("Only json or yaml formats are supported")
	}
}

func (s Swagger) initTemplate() *template.Template {
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
			return groupNames[id]
		},
	}
	data, err := Asset("tmpl/swagger.yaml")
	if err != nil {
		log.Fatal(err)
	}
	tmpl, err := template.New("swagger.yaml").Funcs(funcMap).Parse(string(data))
	if err != nil {
		log.Fatal(err)
	}
	return tmpl
}

func (s Swagger) generateJSON() {
	tmpl := s.initTemplate()

	var tpl bytes.Buffer
	err := tmpl.Execute(&tpl, s)
	if err != nil {
		log.Fatal(err)
	}
	data, err := yaml.YAMLToJSON(tpl.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("swagger.json", data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func (s Swagger) generateYAML() {
	tmpl := s.initTemplate()

	f, err := os.Create("swagger.yaml")
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(f, s)
	if err != nil {
		log.Fatal(err)
	}
}

func fetchVariables(resource *Resource) {
	re := regexp.MustCompile("/{(.*?)}#*")
	for _, param := range re.FindAllStringSubmatch(resource.URL, -1) {
		resource.InsomniaParams = append(resource.InsomniaParams, param[1])
	}
}
