package ds

import (
	"time"
)

type User struct {
	UserId    uint   `gorm:"primaryKey"`
	Login     string `gorm:"size:20;not null"`
	Password  string `gorm:"size:20;not null"`
	FIO       string `gorm:"size:50;not null"`
	Phone     string `gorm:"size:13;not null"`
	Email     string `gorm:"size:30;not null"`
	Moderator bool   `gorm:"not null"`
}

type Language struct {
	LanguageId  uint   `gorm:"primaryKey"`
	Name        string `gorm:"size:20;not null"`
	Subject     string `gorm:"size:70;not null"`
	ImageURL    string `gorm:"size:100;not null"`
	Task        string `gorm:"size:30;not null"`
	Description string `gorm:"size:1000;not null"`
	IsDeleted   bool   `gorm:"not null"`
}

type Form struct {
	FormId         uint       `gorm:"primaryKey"`
	CreationDate   time.Time  `gorm:"not null;type:date"`
	FormationDate  *time.Time `gorm:"type:date"`
	CompletionDate *time.Time `gorm:"type:date"`
	ModeratorId    uint       `gorm:"not null"`
	StudentId      uint       `gorm:"not null"`
	Status         string     `gorm:"size:30;not null"`
	Comments       string     `gorm:"size:300"`

	Moderator User `gorm:"foreignKey:ModeratorId"`
	Student   User `gorm:"foreignKey:StudentId"`
}

type Code struct {
	LanguageId uint `gorm:"primaryKey;not null;autoIncrement:false"`
	FormId     uint `gorm:"primaryKey;not null;autoIncrement:false"`

	Language *Language `gorm:"foreignKey:LanguageId"`
	Form     *Form     `gorm:"foreignKey:FormId"`
}
