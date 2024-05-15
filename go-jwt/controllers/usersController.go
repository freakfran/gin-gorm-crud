
package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go-jwt/initializers"
	"go-jwt/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

// SingUp 处理用户注册请求。
// c: Gin上下文，用于处理HTTP请求和响应。
func SingUp(c *gin.Context) {
	// 定义请求体结构体
	var body struct {
		Email    string
		Password string
		Name     string
	}
	// 解析请求体
	err := c.Bind(&body)
	if err != nil {
		// 请求体解析失败，返回错误信息
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// 使用bcrypt算法加密密码
	password, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		// 密码加密失败，返回错误信息
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// 创建用户模型实例并填充数据
	user := models.User{
		Name:     body.Name,
		Email:    body.Email,
		Password: string(password),
	}
	// 在数据库中创建用户
	result := initializers.DB.Create(&user)
	if result.Error != nil {
		// 创建用户失败，返回错误信息
		c.JSON(http.StatusBadRequest, gin.H{
			"error": result.Error.Error(),
		})
		return
	}
	// 用户创建成功，返回成功信息
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

// Login 处理用户登录请求。
// c: Gin上下文，用于处理HTTP请求和响应。
func Login(c *gin.Context) {
	// 定义请求体结构体
	var body struct {
		Email    string
		Password string
	}
	// 解析请求体
	err := c.Bind(&body)
	if err != nil {
		// 请求体解析失败，返回错误信息
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 在数据库中查找对应邮箱的用户
	var userDb models.User
	initializers.DB.First(&userDb, "email = ?", body.Email)
	if userDb.ID == 0 {
		// 用户不存在，返回错误信息
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "email not registered",
		})
		return
	}
	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(userDb.Password), []byte(body.Password)); err != nil {
		// 密码不匹配，返回错误信息
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "wrong password",
		})
		return
	}

	// 创建并签发JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userDb.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	signedString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	c.SetSameSite(http.SameSiteLaxMode)
	// 设置Cookie携带JWT token
	c.SetCookie("Authorization", signedString, 3600*24*30, "", "", false, true)
	if err != nil {
		// Token创建失败，返回错误信息
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 登录成功，返回空响应体
	c.JSON(http.StatusOK, gin.H{})
}

// Validate 验证用户登录状态。
// c: Gin上下文，用于处理HTTP请求和响应。
func Validate(c *gin.Context) {
	// 从上下文中获取用户信息
	value, _ := c.Get("user")
	// 返回登录状态和用户信息
	c.JSON(http.StatusOK, gin.H{
		"message": "i am login",
		"user":    value,
	})
}