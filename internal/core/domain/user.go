package domain


type User struct {
	ID int
	Version int

	FullName string
	PhoneNumber *string//ссылка т.к номер не обязателен
}