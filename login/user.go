package login

var Users = map[string]struct {
	Password string
	Role     string
}{
	"alice": {Password: "password123", Role: "admin"},
	"bob":   {Password: "hunter2", Role: "user"},
}
