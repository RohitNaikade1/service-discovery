package models

type Resource struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

type VirtualMachine struct {
	Name   string `json:"name,omitempty"`
	Type   string `json:"type,omitempty"`
	OsType string `json:"ostype,omitempty"`
}

type Properties struct {
	StorageProfile    StorageProfiles `json:"storageProfile,omitempty"`
	ProvisioningState string          `json:"provisioningState,omitempty"`
	VMId              string          `json:"vmId,omitempty"`
}

type StorageProfiles struct {
	DataDisk       []DataDisks     `json:"dataDisks,omitempty"`
	ImageReference ImageReferences `json:"imageReference,omitempty"`
	OSDisk         OSDisks         `json:"osDisk,omitempty"`
}

type DataDisks struct {
}

type ImageReferences struct {
}

type OSDisks struct {
	Caching      string      `json:"caching,omitempty"`
	DiskSizeGB   float64     `json:"diskSizeGB,omitempty"`
	ManageDisk   ManageDisks `json:"managedDisk,omitempty"`
	Name         string      `json:"name,omitempty"`
	OSType       string      `json:"osType,omitempty"`
	CreateOption string      `json:"createOption,omitempty"`
}

type ManageDisks struct {
	ID                 string `json:"id,omitempty"`
	StorageAccountType string `json:"storageAccountType,omitempty"`
}
