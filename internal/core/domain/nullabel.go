package domain

type Nullable[T any] struct {
	Value *T
	Set   bool
}

/*
JSON: {
	"phone_number": null
}
даже если передали null то мы сможем понять установлено ли там хоть что то
NullableString:
 - Value: *nil
 - Set: true
*/
