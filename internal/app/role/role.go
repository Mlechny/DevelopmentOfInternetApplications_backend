package role

type Role int

const (
	NotAuthorized Role = iota // 0
	Student                   // 1
	Moderator                 // 2
)
