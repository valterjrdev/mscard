package entity

const (
	OperationTableName = "operation"
)

type Operation struct {
	ID          uint   `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	Description string `json:"description" gorm:"type:varchar(80);column:description"`
	Debit       bool   `json:"debit" gorm:"type:boolean;column:debit;default:false"`
}

func (a *Operation) TableName() string {
	return OperationTableName
}
