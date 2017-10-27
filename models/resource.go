package models

type Resource struct {
	ID                              string         `json:"_id"`
	ParentID                        string         `json:"parentId"`
	Modified                        int64          `json:"modified"`
	Created                         int64          `json:"created"`
	Name                            string         `json:"name"`
	Description                     string         `json:"description,omitempty"`
	Certificates                    []interface{}  `json:"certificates,omitempty"`
	Type                            string         `json:"_type"`
	Data                            DataUrl        `json:"data,omitempty"`
	Color                           interface{}    `json:"color,omitempty"`
	Cookies                         []interface{}  `json:"cookies,omitempty"`
	Environment                     struct{}       `json:"environment,omitempty"`
	URL                             string         `json:"url,omitempty"`
	Method                          string         `json:"method,omitempty"`
	Body                            EntityBody     `json:"body,omitempty"`
	Parameters                      []interface{}  `json:"parameters,omitempty"`
	Headers                         []EntityHeader `json:"headers,omitempty"`
	Authentication                  struct{}       `json:"authentication,omitempty"`
	SettingStoreCookies             bool           `json:"settingStoreCookies,omitempty"`
	SettingSendCookies              bool           `json:"settingSendCookies,omitempty"`
	SettingDisableRenderRequestBody bool           `json:"settingDisableRenderRequestBody,omitempty"`
	SettingEncodeURL                bool           `json:"settingEncodeUrl,omitempty"`
	InsomniaParams                  []string       `json:"insomnia_params,omitempty"`
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

type DataUrl struct {
	BaseURL string `json:"base_url"`
	RefURL  string `json:"ref_url"`
}
