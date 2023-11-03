package models

type Client struct {
	Name      string `sql:"name" gorm:"primaryKey"`
	Secret    string `sql:"secret"`
	GrantType string `sql:"grant_type"`

	RedirectURI string  `sql:"redirect_uri"`
	Scopes      []Scope `sql:"scopes" gorm:"many2many:client_scopes"`
}

func (Client) CreateTable() string {
	return "clients"
}

func (client Client) Equals(other interface{}) bool {
	if other == nil {
		return false
	}
	for i := 0; i < len(client.Scopes); i++ {
		if !client.Scopes[i].Equals(other.(Client).Scopes[i]) {
			return false
		}
	}
	return client.Name == other.(Client).Name &&
		client.Secret == other.(Client).Secret &&
		client.GrantType == other.(Client).GrantType &&
		client.RedirectURI == other.(Client).RedirectURI

}

func (client Client) From() ClientJson {
	scopes := []string{}
	for i := 0; i < len(client.Scopes); i++ {
		scopes = append(scopes, client.Scopes[i].Name)
	}
	client_json := ClientJson{
		Name:        client.Name,
		GrantType:   client.GrantType,
		Secret:      client.Secret,
		RedirectURI: client.RedirectURI,
		Scopes:      scopes,
	}
	return client_json
}

func (cJson ClientJson) From() Client {
	scopes := []Scope{}
	for i := 0; i < len(cJson.Scopes); i++ {
		scope := Scope{
			Name: cJson.Scopes[i],
		}
		scopes = append(scopes, scope)
	}
	client := Client{
		Name:        cJson.Name,
		GrantType:   cJson.GrantType,
		Secret:      cJson.Secret,
		RedirectURI: cJson.RedirectURI,
		Scopes:      scopes,
	}
	return client
}

func ListOfClientsToListOfClientJson(clients []Client) []ClientJson {
	var client_json = []ClientJson{}
	for i := 0; i < len(clients); i++ {
		cJson := clients[i].From()
		client_json = append(client_json, cJson)
	}
	return client_json
}

type ClientJson struct {
	Name      string `json:"name"`
	GrantType string `json:"grant_type"`
	Secret    string `json:"secret",omitempty`

	RedirectURI string   `json:"redirect_uri"`
	Scopes      []string `json:"scope"`
}
