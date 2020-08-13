package models

const (
	ReceiversTable = "plan_receiver"
)

type Receivers struct {
	ID                    int64  `gorm:"auto" json:"id,omitempty"`
	Plan                  Plans  `gorm:"foreignkey:PlanID;save_associations:false:false"`
	PlanID                int64  `json:"plan_id;index"`
	StartTime             string `gorm:"size:31" json:"start_time"`
	EndTime               string `gorm:"size:31" json:"end_time"`
	Start                 int    `json:"start"`
	Period                int    `json:"period"`
	Expression            string `gorm:"size:1023" json:"expression"`
	ReversePolishNotation string `gorm:"size:1023" json:"reverse_polish_notation"`
	User                  string `gorm:"size:1023" json:"user"`
	Group                 string `gorm:"size:1023" json:"group"`
	DutyGroup             string `gorm:"size:255" json:"duty_group"`
	Method                string `gorm:"size:255" json:"method"`
	Model
}

func (*Receivers) TableName() string {
	return "plan_receiver"
}
