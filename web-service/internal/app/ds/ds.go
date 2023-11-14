package ds

import "time"

const DRAFT string = "черновик"
const FORMED string = "сформирован"
const COMPLETED string = "завершён"
const REJECTED string = "отклонён"
const DELETED string = "удалён"

type User struct {
	UUID      string `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"-"`
	Login     string `gorm:"size:20;not null" json:"-"`
	Password  string `gorm:"size:20;not null" json:"-"`
	FIO       string `gorm:"size:50;not null" json:"fio"`
	Phone     string `gorm:"size:13;not null" json:"phone"`
	Email     string `gorm:"size:30;not null" json:"email"`
	Moderator bool   `gorm:"not null" json:"-"`
}

type Language struct {
	UUID        string  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"uuid" binding:"-"`
	Name        string  `gorm:"size:20;not null" form:"name" json:"name" binding:"required"`
	Subject     string  `gorm:"size:70;not null" form:"subject" json:"subject" binding:"required"`
	ImageURL    *string `gorm:"size:100;not null" json:"image_url" binding:"-"`
	Task        string  `gorm:"size:30;not null" form:"task" json:"task" binding:"required"`
	Description string  `gorm:"size:1000;not null" form:"description" json:"description" binding:"required"`
	IsDeleted   bool    `gorm:"not null;default:false" json:"-" binding:"-"`
}

type Form struct {
	UUID           string     `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreationDate   time.Time  `gorm:"not null;type:timestamp"`
	FormationDate  *time.Time `gorm:"type:timestamp"`
	CompletionDate *time.Time `gorm:"type:timestamp"`
	ModeratorId    *string    `json:"-"`
	StudentId      string     `gorm:"not null"`
	Status         string     `gorm:"size:30;not null"`
	Comments       string     `gorm:"size:300"`

	Moderator *User
	Student   User
}

type Code struct {
	LanguageId string `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"language_id"`
	FormId     string `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"form_id"`

	Language *Language `gorm:"foreignKey:LanguageId" json:"language"`
	Form     *Form     `gorm:"foreignKey:FormId" json:"form"`
}
