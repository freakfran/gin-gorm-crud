package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-jwt/initializers"
	"go-jwt/models"
	"net/http"
	"os"
	"time"
)

// RequireAuth 是一个中间件函数，用于验证用户身份。
// 它通过读取请求中的cookie来获取JWT令牌，并验证该令牌的有效性。
// 如果令牌有效且未过期，并且对应的用户存在于数据库中，则允许请求通过，并将用户信息设置为上下文的一部分。
// 参数:
// - c *gin.Context: Gin框架的上下文对象，用于访问请求信息、设置响应状态和传递用户信息。
func RequireAuth(c *gin.Context) {
	// 尝试从cookie中读取Authorization令牌
	token, err := c.Cookie("Authorization")
	if err != nil {
		// 如果无法读取cookie，则中止请求并返回未授权状态
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	// 解析JWT令牌，验证其签名和有效期
	parseToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			// 如果令牌的签名方法不符合预期，则返回未授权状态
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil // 使用环境变量中的密钥来验证签名
	})

	if claims, ok := parseToken.Claims.(jwt.MapClaims); ok && parseToken.Valid {
		// 检查令牌是否已过期
		exp := claims["exp"].(float64)
		if exp < float64(time.Now().Unix()) {
			// 如果令牌已过期，则中止请求并返回未授权状态
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		// 根据令牌中的用户ID查询数据库，验证用户存在性
		var user models.User
		initializers.DB.Where("id = ?", claims["sub"]).First(&user)
		if user.ID == 0 {
			// 如果用户不存在于数据库中，则中止请求并返回未授权状态
			c.AbortWithStatus(http.StatusUnauthorized)
		}
		// 如果令牌有效且用户存在，则将用户信息设置到上下文中，供后续处理使用
		c.Set("user", user)
		c.Next() // 允许请求继续处理
	} else {
		// 如果令牌无效或解析失败，则中止请求并返回未授权状态
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}
