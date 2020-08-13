package models

const (
	HostsTable = "host"
)

type Hosts struct {
	Id       int64  `gorm:"auto" json:"id,omitempty"`
	Mid      int64  `json:"mid"`
	Hostname string `gorm:"size:255" json:"hostname"`
	Model
}

func (*Hosts) TableName() string {
	return HostsTable
}
