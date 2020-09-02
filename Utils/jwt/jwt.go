package jwt

type jwtCustomClaims struct {
	StandardClaims

	// addition
	Uid   uint `json:"uid"`
	Admin bool `json:"admin"`
}

/**
 * 生成 token
 * SecretKey 是一个 const 常量
 */
func CreateToken(SecretKey []byte, userName string, Uid uint, expiresAt int64) (string, int64) {
	claims := &jwtCustomClaims{
		StandardClaims{
			ExpiresAt: expiresAt,
			Issuer:    userName,
		},
		Uid,
		false,
	}

	token := NewWithClaims(SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(SecretKey)
	return tokenString, claims.ExpiresAt
}

/**
 * 生成自定义Claims token
 * SecretKey []byte("Your Secret Key")
 * customClaims
 */
func CreateCustomToken(SecretKey []byte, customClaims Claims) (tokenString string, err error) {
	token := NewWithClaims(SigningMethodHS256, customClaims)
	tokenString, err = token.SignedString(SecretKey)
	return
}

func ParseToken(tokenSrt string, SecretKey []byte) (claims Claims, err error) {
	var token *Token
	token, err = Parse(tokenSrt, func(*Token) (interface{}, error) {
		return SecretKey, nil
	})
	claims = token.Claims
	return
}
