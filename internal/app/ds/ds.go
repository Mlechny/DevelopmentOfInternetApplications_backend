package ds

import (
	"time"
	"web-service/internal/app/role"
)

const StatusDraft string = "черновик"
const StatusFormed string = "сформирован"
const StatusCompleted string = "завершён"
const StatusRejected string = "отклонён"
const StatusDeleted string = "удалён"

const TestingSuccess string = "код верен, ошибки не обнаружены"
const TestingFailure string = "код неверен, обнаружены ошибки"
const TestingStart string = "проверка начата"

type User struct {
	UUID     string `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"-"`
	Role     role.Role
	Login    string `gorm:"size:30;not null" json:"login"`
	Password string `gorm:"size:40;not null" json:"-"`
	// The SHA-1 hash is 20 bytes. When encoded in hexadecimal, each byte is represented by two characters. Therefore, the resulting hash string will be 40 characters long
}

type Language struct {
	UUID        string  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"uuid" binding:"-"`
	Name        string  `gorm:"size:20;not null" form:"name" json:"name" binding:"required"`
	Subject     string  `gorm:"size:70;not null" form:"subject" json:"subject" binding:"required"`
	ImageURL    *string `gorm:"size:100" json:"image_url" binding:"-"`
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
	Comments       *string    `gorm:"size:300"`
	Autotest       *string    `gorm:"size:50"`

	Moderator *User
	Student   User
}

type Code struct {
	LanguageId string `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"language_id"`
	FormId     string `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"form_id"`
	Github     string `gorm:"size:50"`

	Language *Language `gorm:"foreignKey:LanguageId" json:"language"`
	Form     *Form     `gorm:"foreignKey:FormId" json:"form"`
}
