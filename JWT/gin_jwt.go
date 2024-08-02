package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

const TokenExpireDuration = time.Hour * 2

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	r.POST("/login", loginHandler)
	r.GET("/book", JM, bookHandler)
	r.Run(":8080")
}

func bookHandler(c *gin.Context) {
	username := c.MustGet("username").(string)
	if username != "" {
		fmt.Println(username)
		c.JSON(http.StatusOK, gin.H{"books": getBook()})
	}
}

func loginHandler(c *gin.Context) {
	// 用户发送用户名和密码过来
	var user User
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 2001,
			"msg":  "无效的参数",
		})
		return
	}
	// 校验用户名和密码是否正确
	if user.Username == "liyang" && user.Password == "123456" {
		// 生成Token
		tokenString, _ := GenToken(user.Username)
		c.JSON(http.StatusOK, gin.H{"token": tokenString})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 2002,
		"msg":  "鉴权失败",
	})
	return
}

// JWTAuthMiddleware 基于JWT的认证中间件
func JM(c *gin.Context) {
	//return func(c *gin.Context) {
	// 这里假设Token放在Header的Authorization中，并使用Bearer开头
	//authHeader := c.Request.Header.Get("Authorization")
	authHeader := c.Request.Header.Get("Token")
	fmt.Println(authHeader)
	if authHeader == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 2003,
			"msg":  "请求头中auth为空",
		})
		c.Abort()
		return
	} else {
		mc, err := ParseToken(authHeader)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"msg": "无效的Token"})
			c.Abort()
			return
		}
		c.Set("username", mc.Username)
		c.Next()
	}
	/*
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, gin.H{
				"code": 2004,
				"msg":  "请求头中auth格式有误",
			})
			c.Abort()
			return
		}
		mc, err := ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 2005,
				"msg":  "无效的Token",
			})
			c.Abort()
			return
		}
		c.Set("username", mc.Username)
		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	*/
	//}
}

var MySecret = []byte("samli008")

func GenToken(username string) (string, error) {
	// 创建一个我们自己的声明
	c := MyClaims{
		username, // 自定义字段
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(), // 过期时间
			Issuer:    "my-project",                               // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(MySecret)
}

func ParseToken(tokenString string) (*MyClaims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func getBook() []Book {
	return []Book{
		Book{
			Name:   "Book1",
			Author: "Author1",
		},
		Book{
			Name:   "Book2",
			Author: "Author2",
		},
		Book{
			Name:   "Book3",
			Author: "Author3",
		},
	}
}

type Book struct {
	Name   string
	Author string
}
