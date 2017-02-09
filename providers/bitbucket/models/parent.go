package models

type Parent struct {
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
		HTML struct {
			Href string `json:"href"`
		} `json:"html"`
		Avatar struct {
			Href string `json:"href"`
		} `json:"avatar"`
	} `json:"links" kallax:"-"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	UUID     string `json:"uuid"`
}
