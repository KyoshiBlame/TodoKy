package domain



type User struct {
	ID int
	Version int

	FullName string
	PhoneNumber *string//ссылка т.к номер не обязателен
}

func NewUser(
	id int,
	version int,
	fullname string,
	phone_number *string,
) User {
	return User{
		ID: id,
		Version: version,
		FullName: fullname,
		PhoneNumber: phone_number,
	}
}

func NewUserUninitialized(
	fullname string, 
	phone_number *string,
) User {
	return NewUser(
		UninitializedID,
		UninitializedVersion,
		fullname,
		phone_number,
	)
}
