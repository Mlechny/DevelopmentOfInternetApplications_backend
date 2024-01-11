package ds

import (
	"time"
)

type User struct {
	UserId   string `gorm:"primary_key"`
	Name      string `gorm:"size:50;not null"`
	Login    string `gorm:"size:30;not null"`
	Password string `gorm:"size:40;not null"`
	Moderator bool  `gorm:"not null"`
}

type Language struct {
	LanguageId  string  `gorm:"primary_key"`
	Name        string  `gorm:"size:20;not null"`
	Subject     string  `gorm:"size:70;not null"`
	ImageURL    *string `gorm:"size:100"`
	Task        string  `gorm:"size:30;not null"`
	Description string  `gorm:"size:1000;not null"`
	IsDeleted   bool    `gorm:"not null;default:false"`
}

type Form struct {
	FormId           string   `gorm:"primary_key"`
	CreationDate   time.Time  `gorm:"not null;type:timestamp"`
	FormationDate  *time.Time `gorm:"type:timestamp"`
	CompletionDate *time.Time `gorm:"type:timestamp"`
	ModeratorId    *string    `json:"-"`
	StudentId      string     `gorm:"not null"`
	Status         string     `gorm:"size:30;not null"`
	Comments       *string    `gorm:"size:300"`

	Moderator *User
	Student   User
}

type Code struct {
	LanguageId string  `gorm:"primary_key"`
	FormId     string  `gorm:"primary_key"`
	Github     *string `gorm:"size:50"`

	Language *Language `gorm:"foreignKey:LanguageId"`
	Form     *Form     `gorm:"foreignKey:FormId"`
}
