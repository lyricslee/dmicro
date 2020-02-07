package jwt

import (
	"crypto/rsa"
	"dmicro/common/errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type UserInfo struct {
	Appid    int
	Uid      int64
	Plat     int
	DeviceId string
}

type UserClaims struct {
	jwt.StandardClaims
	Info *UserInfo
}

func Encode(privateKey *rsa.PrivateKey, userInfo *UserInfo, expiresIn int) (string, error) {
	claims := UserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(expiresIn) * time.Second).Unix(),
		},
		Info: userInfo,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

func Decode(publicKey *rsa.PublicKey, tokenString string) (*UserInfo, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return publicKey, nil
		},
	)
	if err == nil {
		if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
			// 验证成功，返回信息
			return claims.Info, nil
		}
	}

	ve := err.(*jwt.ValidationError)
	if ve.Errors == jwt.ValidationErrorExpired {
		return nil, errors.ErrTokenExpired
	} else {
		return nil, errors.ErrInvalidToken
	}
	// 验证失败
	return nil, err
}
