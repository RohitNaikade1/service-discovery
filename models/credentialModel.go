package models

type Credentials struct {
	ID             string `bson:"_id" json:"_id,omitempty"`
	User           UserID `json:"user,omitempty"`
	CredsID        string `json:"credsid,omitempty"`
	Provider       string `json:"provider,omitempty"`
	UserName       string `json:"username,omitempty"`
	SubscriptionID string `json:"subscriptionid,omitempty"`
	TenantID       string `json:"tenantid,omitempty"`
	Created_At     string `json:"created_at,omitempty"`
	Updated_At     string `json:"updated_at,omitempty"`
}

type UserID struct {
	ID string `json:"id,omitempty"`
}
