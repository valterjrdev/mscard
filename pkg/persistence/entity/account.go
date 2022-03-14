package entity

type Account struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	Document string `json:"document_number" gorm:"type:varchar(11);unique;column:document_number"`
}

func (a *Account) TableName() string {
	return "account"
}
