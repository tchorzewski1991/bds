package user

// User is a business representation of the user entity.
type User struct {
	UUID        string
	Email       string
	Permissions []string
}
