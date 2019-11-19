package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

const BackendUrl = "http://localhost:8080"

type Post struct {
	Id            int    `json:"id"`
	Title         string `json:"title" form:"postTitle"`
	Content       string `json:"content" form:"postContent"`
	Author        string `json:"author"`
	TimePublished string `json:"time_published"`
	IsPublished   int    `json:"is_published"`
}

type Posts struct {
	Posts []Post `json:"posts"`
}

type Message struct {
	FullName string `form:"fullName" json:"full_name"`
	Email    string `form:"email" json:"email"`
	Message  string `form:"message" json:"message"`
}

type Messages struct {
	Messages []Message `json:"messages"`
}

type RespStatus struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type Session struct {
	Token string `json:"token"`
}

type RespLogin struct {
	Session string     `json:"session"`
	Resp    RespStatus `json:"resp"`
}

type Login struct {
	Email    string `json:"email" form:"inputEmail"`
	Password string `json:"password" form:"inputPassword"`
}

type PostId struct {
	Id string `uri:"id" binding:"required"`
}

type SubmitPost struct {
	Session Session `json:"session"`
	Post    Post    `json:"post"`
}

func home(c *gin.Context) {
	var posts Posts
	data := map[string]interface{}{}
	posts.Posts = make([]Post, 0)

	resp, err := http.Get(BackendUrl + "/home")
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&posts)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}

	data["posts"] = posts
	data["user"] = "yoga"
	c.HTML(http.StatusOK, "index.html", data)
}

func about(c *gin.Context) {
	data := map[string]interface{}{}
	data["user"] = "yoga"
	c.HTML(http.StatusOK, "about.html", data)
}

func contactUs(c *gin.Context) {
	data := map[string]interface{}{}
	data["user"] = "yoga"
	c.HTML(http.StatusOK, "contact.html", data)
}

func messagePost(c *gin.Context) {
	var msg Message
	err := c.ShouldBind(&msg)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}
	msgJson, err := json.Marshal(&msg)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}

	resp, err := http.Post(BackendUrl+"/message", "application/json", bytes.NewBuffer(msgJson))
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}
	defer resp.Body.Close()

	var respStatus RespStatus
	err = json.NewDecoder(resp.Body).Decode(&respStatus)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}
	if respStatus.Status != http.StatusOK {
		_ = c.AbortWithError(http.StatusNotFound, fmt.Errorf("%s", respStatus.Message))
		return
	}

	c.Redirect(http.StatusMovedPermanently, "/contact")
}

func login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func loginPost(c *gin.Context) {
	var login Login
	err := c.ShouldBind(&login)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}
	loginJson, err := json.Marshal(&login)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}
	resp, err := http.Post(BackendUrl+"/login", "application/json", bytes.NewBuffer(loginJson))
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}
	defer resp.Body.Close()

	var respLogin RespLogin
	err = json.NewDecoder(resp.Body).Decode(&respLogin)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}
	if respLogin.Resp.Status != http.StatusOK {
		_ = c.AbortWithError(http.StatusNotFound, fmt.Errorf("%s", respLogin.Resp.Message))
		return
	}
}

func articles(c *gin.Context) {
	var posts Posts
	data := map[string]interface{}{}
	posts.Posts = make([]Post, 0)
	var session Session
	session.Token = ""
	sessionJson, err := json.Marshal(&session)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}
	resp, err := http.Post(BackendUrl+"/articles", "application/json", bytes.NewBuffer(sessionJson))
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&posts)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}

	var msgs Messages
	resp2, err := http.Post(BackendUrl+"/contact", "application/json", bytes.NewBuffer(sessionJson))
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}
	defer resp2.Body.Close()
	err = json.NewDecoder(resp2.Body).Decode(&msgs)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}

	data["msgs"] = msgs
	data["posts"] = posts
	data["user"] = "yoga"
	c.HTML(http.StatusOK, "articles.html", data)
}

func articlesId(c *gin.Context) {
	data := map[string]interface{}{}
	var post Post
	var postId PostId
	err := c.ShouldBindUri(&postId)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}
	var session Session
	session.Token = ""
	sessionJson, err := json.Marshal(&session)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}
	resp, err := http.Post(BackendUrl+"/articles/id/"+postId.Id, "application/json", bytes.NewBuffer(sessionJson))
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&post)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}
	log.Println(post)
	data["post"] = post
	data["user"] = "yoga"
	data["action"] = "/articles/id/" + postId.Id
	c.HTML(http.StatusOK, "articlesId.html", data)
}

func articlesAdd(c *gin.Context) {
	data := map[string]interface{}{}
	data["user"] = "yoga"
	data["action"] = "/articles/add"
	c.HTML(http.StatusOK, "articlesId.html", data)
}

func articlesIdPost(c *gin.Context) {
	var postId PostId
	err := c.ShouldBindUri(&postId)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}
	var post Post
	err = c.ShouldBind(&post)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}
	post.Id, _ = strconv.Atoi(postId.Id)
	post.TimePublished = time.Now().Format("1 January 2006")
	post.IsPublished = 1

	var submitPost SubmitPost
	submitPost.Session.Token = ""
	submitPost.Post = post
	submitPostJson, err := json.Marshal(&submitPost)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}
	resp, err := http.Post(BackendUrl+"/articles/add/"+postId.Id, "application/json", bytes.NewBuffer(submitPostJson))
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}
	defer resp.Body.Close()
	var respStatus RespStatus
	err = json.NewDecoder(resp.Body).Decode(&respStatus)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}
	if respStatus.Status != http.StatusOK {
		_ = c.AbortWithError(http.StatusNotFound, fmt.Errorf("%s", respStatus.Message))
		return
	}
	c.Redirect(http.StatusMovedPermanently, "/articles")
}

func articlesAddPost(c *gin.Context) {
	var post Post
	err := c.ShouldBind(&post)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}
	post.TimePublished = time.Now().Format("1 January 2006")
	post.IsPublished = 1

	var submitPost SubmitPost
	submitPost.Session.Token = ""
	submitPost.Post = post
	submitPostJson, err := json.Marshal(&submitPost)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}
	resp, err := http.Post(BackendUrl+"/articles/add", "application/json", bytes.NewBuffer(submitPostJson))
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}
	defer resp.Body.Close()
	var respStatus RespStatus
	err = json.NewDecoder(resp.Body).Decode(&respStatus)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}
	if respStatus.Status != http.StatusOK {
		_ = c.AbortWithError(http.StatusNotFound, fmt.Errorf("%s", respStatus.Message))
		return
	}
	c.Redirect(http.StatusMovedPermanently, "/articles")
}

func logout(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/login")
}

func main() {
	//gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/", home)
	router.GET("/about", about)
	router.GET("/contact", contactUs)
	router.POST("/message", messagePost)
	router.GET("/login", login)
	router.POST("/login", loginPost)
	router.GET("/articles", articles)
	router.GET("/articles/id/:id", articlesId)
	router.POST("/articles/id/:id", articlesIdPost)
	router.GET("/articles/add", articlesAdd)
	router.POST("/articles/add", articlesAddPost)
	router.GET("/logout", logout)
	err := router.Run(":9090")
	if err != nil {
		log.Fatal(err)
	}
}
