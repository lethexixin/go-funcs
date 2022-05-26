package mid_gin

import (
	"errors"
	"strings"
	"time"
)

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	TokenExpired     = errors.New("token is expired")
	TokenNotValidYet = errors.New("token not active yet")
	TokenMalformed   = errors.New("that's not even a token")
	TokenInvalid     = errors.New("this token could not handle")
)

// 教程地址: https://www.cnblogs.com/marshhu/p/12639633.html
// CustomClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 我们这里需要额外记录下面的几个字段，所以要自定义结构体
// 如果想要保存更多信息，都可以添加到这个结构体中
type CustomClaims struct {
	ID    int64  `json:"id"`
	Phone string `json:"phone"`
	jwt.StandardClaims
}

// 定义JWT的过期时间，这里以3天为例
const TokenExpireDuration = time.Hour * 24 * 3

// 创建一个jwt token
func CreateToken(claims CustomClaims, key []byte) (string, error) {
	claims.StandardClaims.ExpiresAt = time.Now().Add(TokenExpireDuration).Unix() // 过期时间,3天
	claims.StandardClaims.Issuer = "go-gin-model"                                // 签发人
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的key签名并获得完整的编码后的字符串token
	return token.SignedString(key)
}

// 解析 token
func ParseToken(token string, key []byte) (*CustomClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*CustomClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
		return nil, TokenInvalid

	} else {
		return nil, TokenInvalid

	}
}

func JWT(key []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// 这里的具体实现方式要依据你的实际业务情况决定

		// token := c.Query("token")
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			//todo
			//web.Fail(c, web.HandlerAuthIsEmptyCode)
			c.Abort()
			return
		}

		// 按空格分割
		parts := strings.SplitN(token, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			//todo
			//web.Fail(c, web.HandlerAuthIsErrCode)
			c.Abort()
			return
		}

		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		claims, err := ParseToken(parts[1], key)
		if err != nil {
			//todo
			//web.Fail(c, web.TokenIsInvalidCode)
			c.Abort()
			return
		}

		// 将当前请求的用户信息保存到请求的上下文c上
		c.Set("users", claims)
		c.Next() // 后续的处理函数可以用过c.Get("claims") 之类的方法来获取当前请求的用户信息
	}
}
