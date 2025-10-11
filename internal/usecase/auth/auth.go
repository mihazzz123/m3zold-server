package auth

type TokenGenerator interface {
	Generate(userID string) (string, error)
}
