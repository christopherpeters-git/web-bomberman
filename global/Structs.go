package global

type User struct {
	UserID       uint64
	Username     string
	PasswordHash string
	SessionID    string
	Image        []byte
}
