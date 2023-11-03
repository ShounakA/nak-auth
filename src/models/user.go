package models

type User struct {
	ID      int      `json:"id" sql:"id" gorm:"primaryKey"`
	Name    string   `json:"name" sql:"name"`
	Secret  string   `json:"secret" sql:"secret"`
	Clients []Client `json:"-" sql:"authorizedClients" gorm:"many2many:user_clients"`
	Codes   []Code   `json:"-" sql:"codes" gorm:"foreignKey:UserRefer"`
}

func (User) CreateTable() string {
	return "users"
}
