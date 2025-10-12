package user

// type LoginUseCase struct {
// 	Repo user.Repository
// }

// func NewLoginUseCase(repo user.Repository) *RegisterUseCase {
// 	return &RegisterUseCase{Repo: repo}
// }

// func (uc *LoginUseCase) generateSecureToken(length int) (string, error) {
// 	bytes := make([]byte, length)
// 	if _, err := rand.Read(bytes); err != nil {
// 		return "", err
// 	}
// 	return base64.URLEncoding.EncodeToString(bytes), nil
// }

// func (uc *LoginUseCase) generateVerificationToken() (string, error) {
// 	return uc.generateSecureToken(32)
// }

// func (uc *LoginUseCase) generateToken(cfg *config.Config, userID uuid.UUID) (string, error) {
// 	claims := jwt.MapClaims{
// 		"user_id": userID.String(),
// 		"exp":     time.Now().Add(15 * time.Minute).Unix(),
// 		"iat":     time.Now().Unix(),
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	return token.SignedString(cfg.Auth.JWTSecret)
// }
