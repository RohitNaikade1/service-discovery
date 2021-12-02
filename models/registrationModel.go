package models

type Registration struct {
	ID         string  `bson:"_id" json:"_id,omitempty"`
	Name       string  `json:"name"`
	Accounts   Account `json:"Accounts,omitempty"`
	URL        string  `json:"url,omitempty"`
	Categories []Category_info
	Created_At string `bson:"created_at,omitempty" json:"created_at,omitempty"`
	Updated_At string `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

type Account struct {
	CredsId string `json:"credsid"`
}

type Category_info struct {
	Category      string `json:"category,omitempty"`
	Resource_info Resource_Infos
}
type Resource_Infos struct {
	Provider  string   `json:"provider,omitempty"`
	Resources []string `json:"resources,omitempty"`
}
