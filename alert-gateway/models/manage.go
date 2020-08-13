package models

import "time"

const (
	MaintainsTable = "maintain"
)

type Maintains struct {
	Id        int64      `gorm:"auto" json:"id,omitempty"`
	Flag      bool       `json:"flag"`
	TimeStart string     `gorm:"size:15" json:"time_start"`
	TimeEnd   string     `gorm:"size:15" json:"time_end"`
	Month     int        `json:"month"`
	DayStart  int8       `json:"day_start"`
	DayEnd    int8       `json:"day_end"`
	Valid     *time.Time `json:"valid"`
	Model
}

func (*Maintains) TableName() string {
	return MaintainsTable
}
