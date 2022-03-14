package entity

type OperationType struct {
	ID          uint   `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	Description string `json:"description" gorm:"type:varchar(80);column:description"`
	Negative    bool   `json:"negative" gorm:"type:boolean;column:negative;default:false"`
}

func (a *OperationType) TableName() string {
	return "operation_type"
}
