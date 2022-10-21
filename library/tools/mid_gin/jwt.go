package mid_gin

import (
	"errors"
	"net/http"
	"strings"
	"time"
)

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var (
	TokenExpired     = errors.New("token is expired")
	TokenNotValidYet = errors.New("token not active yet")
	TokenMalformed   = errors.New("that's not even a token")
	TokenInvalid     = errors.New("this token could not handle")
)

// CustomClaims
// 教程地址: https://www.cnblogs.com/marshhu/p/12639633.html
// CustomClaims 自定义声明结构体并内嵌jwt.StandardClaims
// jwt.StandardClaims.ExpiresAt 设置过期时间, jwt.StandardClaims.Issuer 设置签发人
// jwt包自带的jwt.StandardClaims只包含了官方字段
// 我们这里需要额外记录下面的几个字段, 所以要自定义结构体
// 如果想要保存更多信息, 都可以添加到这个结构体中
type CustomClaims struct {
	ID    int64  `json:"id"`
	Phone string `json:"phone"`
	jwt.StandardClaims
}

// TokenExpireDuration 定义JWT的过期时间, 这里以1天为例
const TokenExpireDuration = time.Hour * 24

// CreateToken 创建一个jwt token
func CreateToken(claims CustomClaims, key []byte) (string, error) {
	if claims.StandardClaims.ExpiresAt == 0 {
		claims.StandardClaims.ExpiresAt = time.Now().Add(TokenExpireDuration).Unix() // 过期时间,1天
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的key签名并获得完整的编码后的字符串token
	return token.SignedString(key)
}

// ParseToken 解析 token
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
		// 这里假设Token放在Header的Authorization中, 并使用Bearer开头
		// 这里的具体实现方式要依据你的实际业务情况决定

		// token := c.Query("token")
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.String(http.StatusOK, "Authorization Token is empty")
			c.Abort()
			return
		}

		// 按空格分割
		parts := strings.SplitN(token, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.String(http.StatusOK, "Authorization Bearer is empty")
			c.Abort()
			return
		}

		// parts[1]是获取到的tokenString, 我们使用之前定义好的解析JWT的函数来解析它
		claims, err := ParseToken(parts[1], key)
		if err != nil {
			c.String(http.StatusOK, "Token is invalid")
			c.Abort()
			return
		}

		// 将当前请求的用户信息保存到请求的上下文c上
		c.Set("users", claims)
		c.Next() // 后续的处理函数可以用过c.Get("users") 之类的方法来获取当前请求的用户信息
	}
}
