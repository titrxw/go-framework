package database

type ModelAbstract struct {
	ID uint `json:"id" form:"id" gorm:"primarykey"`
}

func (ModelAbstract) Connection() string {
	return "default"
}
