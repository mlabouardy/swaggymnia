package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

type Insomnia struct {
	Type         string     `json:"_type"`
	ExportFormat int        `json:"__export_format"`
	ExportDate   time.Time  `json:"__export_date"`
	ExportSource string     `json:"__export_source"`
	Resources    []Resource `json:"resources"`
}

type Resource struct {
	ID           string        `json:"_id"`
	ParentID     string        `json:"parentId"`
	Modified     int64         `json:"modified"`
	Created      int64         `json:"created"`
	Name         string        `json:"name"`
	Description  string        `json:"description,omitempty"`
	Certificates []interface{} `json:"certificates,omitempty"`
	Type         string        `json:"_type"`
	Data         struct {
		BaseURL string `json:"base_url"`
		RefURL  string `json:"ref_url"`
	} `json:"data,omitempty"`
	Color       interface{}   `json:"color,omitempty"`
	Cookies     []interface{} `json:"cookies,omitempty"`
	Environment struct {
	} `json:"environment,omitempty"`
	URL            string         `json:"url,omitempty"`
	Method         string         `json:"method,omitempty"`
	Body           EntityBody     `json:"body,omitempty"`
	Parameters     []interface{}  `json:"parameters,omitempty"`
	Headers        []EntityHeader `json:"headers,omitempty"`
	Authentication struct {
	} `json:"authentication,omitempty"`
	SettingStoreCookies             bool     `json:"settingStoreCookies,omitempty"`
	SettingSendCookies              bool     `json:"settingSendCookies,omitempty"`
	SettingDisableRenderRequestBody bool     `json:"settingDisableRenderRequestBody,omitempty"`
	SettingEncodeURL                bool     `json:"settingEncodeUrl,omitempty"`
	InsomniaParams                  []string `json:"insomnia_params,omitempty"`
}

type EntityBody struct {
	MimeType string        `json:"mimeType"`
	Params   []EntityParam `json:"params"`
}

type EntityParam struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	ID       string `json:"id"`
	Disabled bool   `json:"disabled"`
}

type EntityHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type ConfigSwagger struct {
	Title       string `json:"title"`
	Version     string `json:"version"`
	Host        string `json:"host"`
	BasePath    string `json:"bastPath"`
	Schemes     string `json:"schemes"`
	Description string `json:"description"`
}

type ParentID string

var groups map[string]string

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

func main() {
	insomniaFileName := flag.String("insomnia", "", "Insomnia JSON file. (Required)")
	configFileName := flag.String("config", "", "Configuration file. (Required)")
	flag.Parse()

	if *insomniaFileName == "" || *configFileName == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	raw, err := ioutil.ReadFile(*insomniaFileName)
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

	GenerateSwagger(list, *configFileName)
}
