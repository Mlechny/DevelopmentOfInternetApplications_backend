package role

type Role int

const (
	NotAuthorized Role = iota
	Student
	Moderator
)
