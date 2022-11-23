package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

type User struct {
	username  string
	password  string
	question1 string
	answer1   string
	question2 string
	answer2   string
	question3 string
	answer3   string
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
func send(c *gin.Context) (string, string, string, string, string, string, string, string) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	question1 := c.PostForm("question1")
	answer1 := c.PostForm("answer1")
	question2 := c.PostForm("question2")
	answer2 := c.PostForm("answer2")
	question3 := c.PostForm("question3")
	answer3 := c.PostForm("answer3")
	return username, password, question1, answer1, question2, answer2, question3, answer3
}
func QueryRow(a *sql.DB, x string) (string, string, string, string, string, string, string) {
	//单行查看
	row := a.QueryRow("select * from relo1 where username=?", x)
	var U User
	err := row.Scan(&U.username, &U.password, &U.question1, &U.answer1, &U.question2, &U.answer2, &U.question3, &U.answer3)
	checkErr(err)
	return U.password, U.question1, U.question2, U.question3, U.answer1, U.answer2, U.answer3
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
func insert(a *sql.DB, b, c, d, e, f, g, h, i string) {
	_, err := a.Exec("insert into relo1 (username,password,question1,answer1,question2,answer2,question3,answer3) value (?,?,?,?,?,?,?,?)", b, c, d, e, f, g, h, i)
	checkErr(err)
}
func protect(c *gin.Context) {
	//访问、测试连接
	db := visit()
	connect(db)
	//获取用户名
	username, _, _, _, _, _, _, _ := send(c)
	//获取问题和答案
	password, question1, question2, question3, correctAnswer1, correctAnswer2, correctAnswer3 := QueryRow(db, username)
	//判断用户名是否存在
	b := SelectExist(password)
	if !b {
		c.JSON(http.StatusOK, gin.H{
			"status":  200,
			"message": "user don't exist",
		})
		return
	}
	//判断密保问题是否都正确
	b1 := SelectEqual(c.PostForm(question1), correctAnswer1)
	b2 := SelectEqual(c.PostForm(question2), correctAnswer2)
	b3 := SelectEqual(c.PostForm(question3), correctAnswer3)
	//正确返回密码
	if b1 && b2 && b3 {
		c.JSON(http.StatusOK, gin.H{
			"status":   200,
			"password": password,
		})
		return
	}
	//错误提示错误
	c.JSON(http.StatusInternalServerError, gin.H{
		"status":  500,
		"message": "incorrect answer",
	})
	return
}
func register(c *gin.Context) {
	//访问、测试连接
	db := visit()
	connect(db)
	//获取用户名和密码、密保问题和答案
	username, password, question1, answer1, question2, answer2, question3, answer3 := send(c)
	//从数据库获取正确密码
	correctPassword, _, _, _, _, _, _ := QueryRow(db, username)
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
	insert(db, username, password, question1, answer1, question2, answer2, question3, answer3)
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
	username, password, _, _, _, _, _, _ := send(c)
	correctPassword, _, _, _, _, _, _ := QueryRow(db, username)
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
