package dist

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

//type StandardClaims struct {
//	Audience  string `json:"aud,omitempty"`
//	ExpiresAt int64  `json:"exp,omitempty"`
//	Id        string `json:"jti,omitempty"`
//	IssuedAt  int64  `json:"iat,omitempty"`
//	Issuer    string `json:"iss,omitempty"`
//	NotBefore int64  `json:"nbf,omitempty"`
//	Subject   string `json:"sub,omitempty"`
//}

type JwtTokener struct {
	key string
}

func NewJwtTokener(key string) *JwtTokener {
	return &JwtTokener{
		key: key,
	}
}

func (self *JwtTokener) Token(m map[string]interface{}) (string, error) {
	claims := make(jwt.MapClaims)
	for k, v := range m {
		claims[k] = v
	}

	claim := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := claim.SignedString([]byte(self.key))
	if err != nil {
		return "", err
	}
	return token, nil
}

func (self *JwtTokener) Parse(token string) (jwt.MapClaims, error) {
	tk, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(self.key), nil
	})
	if nil != err {
		return nil, err
	}
	claims, ok := tk.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("bad token: %v", token)
	}
	return claims, nil
}
