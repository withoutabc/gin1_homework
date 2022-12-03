package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

type User struct {
	username string
	password string
}

func checkErr(err error) {
	if err != nil {
		log.Println(err)
		return
	}
}
func visit() *sql.DB {
	//访问数据库
	db, err := sql.Open("mysql", "root:224488@tcp(127.0.0.1:3306)/relo")
	if err != nil {
		panic("open failed")
	}
	return db
}
func connect(a *sql.DB) {
	//查询是否建立连接
	err := a.Ping()
	if err != nil {
		panic("connect failed")
		return
	}
}
func send(c *gin.Context) (string, string) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	return username, password
}
func QueryRow(a *sql.DB, x string) string {
	//单行查看
	row := a.QueryRow("select * from relo1 where username=?", x)
	var U User
	err := row.Scan(&U.username, &U.password)
	checkErr(err)
	return U.password
}
func SelectExist(a string) bool {
	if a == "" {
		return false
	}
	return true
}
func SelectEqual(a, b string) bool {
	if a == b {
		return true
	}
	return false
}
func insert(a *sql.DB, b, c string) {
	_, err := a.Exec("insert into relo1 (username,password) value (?,?)", b, c)
	checkErr(err)
}

func register(c *gin.Context) {
	//访问、测试连接
	db := visit()
	connect(db)
	//获取用户名和密码、密保问题和答案
	username, password := send(c)
	//从数据库获取正确密码
	correctPassword := QueryRow(db, username)
	//通过是否有密码 判断用户名是否已存在
	b1 := SelectExist(correctPassword)
	if b1 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "repeated name",
		})
		return
	}
	//向数据库添加注册的用户信息，添加密保
	insert(db, username, password)
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "success register",
	})
}
func login(c *gin.Context) {
	//访问、测试连接
	db := visit()
	connect(db)
	//获取用户名、密码
	username, password := send(c)
	correctPassword := QueryRow(db, username)
	//判断用户名是否存在
	b1 := SelectExist(correctPassword)
	if !b1 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "user don't exist",
		})
		return
	}
	//判断密码是否正确
	b2 := SelectEqual(password, correctPassword)
	//错误提示错误
	if !b2 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "wrong password",
		})
		return
	}
	//正确设置cookie
	c.SetCookie("gin_cookie", "test", 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "login successful",
	})
}
func main() {
	r := gin.Default()
	r.POST("/register", register)
	r.POST("/login", login)
	r.POST("/protect", protect)
	r.Run()
}
