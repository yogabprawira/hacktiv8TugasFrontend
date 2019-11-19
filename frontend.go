package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"log"
	"net/http"
	"strconv"
	"time"
)

// username : yoga
// password : yoga123

const BackendUrl = "http://localhost:8080"
const PwdSalt = "rIc[@(}sgO>LNyAzaJ?k.RUhYOKZtQ#rlB+$r-e%rr*L-CF+33JTrg@}50E`X/50"
const SessionName = "auth-session"

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
	Session  Session    `json:"session"`
	Username string     `json:"username"`
	Name     string     `json:"name"`
	Resp     RespStatus `json:"resp"`
}

type Login struct {
	Username string `json:"username" form:"inputUsername"`
	Password string `json:"password" form:"inputPassword"`
}

type ReqLogin struct {
	Session Session `json:"session"`
	Login   Login   `json:"login"`
}

type PostId struct {
	Id string `uri:"id" binding:"required"`
}

type SubmitPost struct {
	Session Session `json:"session"`
	Post    Post    `json:"post"`
}

func hashPwd(username string, pwd string, salt string) string {
	s, _ := base64.StdEncoding.DecodeString(salt)
	pwdAdd := append([]byte(username+pwd), s...)
	h := sha256.New()
	h.Write(pwdAdd)
	result := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(result)
}

func createSession(username string, pwd string, salt string) string {
	s, _ := base64.StdEncoding.DecodeString(salt)
	pwdAdd := append([]byte(username+pwd), s...)
	pwdAdd = append(pwdAdd, []byte(time.Now().Format(time.RFC3339Nano))...)
	h := sha256.New()
	h.Write(pwdAdd)
	result := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(result)
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

	session := sessions.Default(c)
	user := ""
	val := session.Get("username")
	if val != nil {
		user = val.(string)
	}

	data["posts"] = posts
	data["user"] = user
	c.HTML(http.StatusOK, "index.html", data)
}

func about(c *gin.Context) {
	data := map[string]interface{}{}
	session := sessions.Default(c)
	user := ""
	val := session.Get("username")
	if val != nil {
		user = val.(string)
	}
	data["user"] = user
	c.HTML(http.StatusOK, "about.html", data)
}

func contactUs(c *gin.Context) {
	data := map[string]interface{}{}
	session := sessions.Default(c)
	user := ""
	val := session.Get("username")
	if val != nil {
		user = val.(string)
	}
	data["user"] = user
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
	session := sessions.Default(c)
	user := ""
	val := session.Get("username")
	if val != nil {
		user = val.(string)
	}
	if len(user) > 0 {
		c.Redirect(http.StatusMovedPermanently, "/")
	}
	c.HTML(http.StatusOK, "login.html", nil)
}

func loginPost(c *gin.Context) {
	var login Login
	var reqLogin ReqLogin
	err := c.ShouldBind(&login)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}
	pwdHash := hashPwd(login.Username, login.Password, PwdSalt)
	sessionId := createSession(login.Username, login.Password, PwdSalt)
	login.Password = pwdHash
	reqLogin.Login = login
	reqLogin.Session.Token = sessionId
	loginJson, err := json.Marshal(&reqLogin)
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
	session := sessions.Default(c)
	session.Set("username", respLogin.Username)
	session.Set("name", respLogin.Name)
	session.Set("sessionId", respLogin.Session.Token)
	_ = session.Save()
	c.Redirect(http.StatusMovedPermanently, "/")
}

func articles(c *gin.Context) {
	var posts Posts
	data := map[string]interface{}{}
	posts.Posts = make([]Post, 0)
	var sessionSt Session

	session := sessions.Default(c)
	user := ""
	token := ""
	val := session.Get("username")
	if val != nil {
		user = val.(string)
	}
	val = session.Get("sessionId")
	if val != nil {
		token = val.(string)
	}

	sessionSt.Token = token
	sessionJson, err := json.Marshal(&sessionSt)
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
	data["user"] = user
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

	user := ""
	token := ""
	session := sessions.Default(c)
	val := session.Get("username")
	if val != nil {
		user = val.(string)
	}
	val = session.Get("sessionId")
	if val != nil {
		token = val.(string)
	}

	var sessionSt Session
	sessionSt.Token = token
	sessionJson, err := json.Marshal(&sessionSt)
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
	data["post"] = post
	data["user"] = user
	data["action"] = "/articles/id/" + postId.Id
	c.HTML(http.StatusOK, "articlesId.html", data)
}

func articlesAdd(c *gin.Context) {
	data := map[string]interface{}{}
	session := sessions.Default(c)
	user := ""
	val := session.Get("username")
	if val != nil {
		user = val.(string)
	}
	data["user"] = user
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

	token := ""
	user := ""
	session := sessions.Default(c)
	val := session.Get("sessionId")
	if val != nil {
		token = val.(string)
	}
	val = session.Get("username")
	if val != nil {
		user = val.(string)
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
	post.Author = user

	var submitPost SubmitPost
	submitPost.Session.Token = token
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

	token := ""
	user := ""
	session := sessions.Default(c)
	val := session.Get("sessionId")
	if val != nil {
		token = val.(string)
	}
	val = session.Get("username")
	if val != nil {
		user = val.(string)
	}
	post.Author = user

	var submitPost SubmitPost
	submitPost.Session.Token = token
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

func articlesIdPublish(c *gin.Context) {
	var postId PostId
	err := c.ShouldBindUri(&postId)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}

	token := ""
	session := sessions.Default(c)
	val := session.Get("sessionId")
	if val != nil {
		token = val.(string)
	}

	var sessionSt Session
	sessionSt.Token = token
	sessionJson, err := json.Marshal(&sessionSt)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}
	resp, err := http.Post(BackendUrl+"/articles/id/"+postId.Id+"/publish", "application/json", bytes.NewBuffer(sessionJson))
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
		_ = c.AbortWithError(http.StatusNotFound, fmt.Errorf(respStatus.Message))
		return
	}
	c.Redirect(http.StatusMovedPermanently, "/articles")
}

func articlesIdUnpublish(c *gin.Context) {
	var postId PostId
	err := c.ShouldBindUri(&postId)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}

	token := ""
	session := sessions.Default(c)
	val := session.Get("sessionId")
	if val != nil {
		token = val.(string)
	}

	var sessionSt Session
	sessionSt.Token = token
	sessionJson, err := json.Marshal(&sessionSt)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}
	resp, err := http.Post(BackendUrl+"/articles/id/"+postId.Id+"/unpublish", "application/json", bytes.NewBuffer(sessionJson))
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
		_ = c.AbortWithError(http.StatusNotFound, fmt.Errorf(respStatus.Message))
		return
	}
	c.Redirect(http.StatusMovedPermanently, "/articles")
}

func articlesIdDelete(c *gin.Context) {
	var postId PostId
	err := c.ShouldBindUri(&postId)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}

	token := ""
	session := sessions.Default(c)
	val := session.Get("sessionId")
	if val != nil {
		token = val.(string)
	}

	var sessionSt Session
	sessionSt.Token = token
	sessionJson, err := json.Marshal(&sessionSt)

	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}
	resp, err := http.Post(BackendUrl+"/articles/id/"+postId.Id+"/delete", "application/json", bytes.NewBuffer(sessionJson))
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
		_ = c.AbortWithError(http.StatusNotFound, fmt.Errorf(respStatus.Message))
		return
	}
	c.Redirect(http.StatusMovedPermanently, "/articles")
}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Set("username", "")
	session.Set("name", "")
	session.Set("sessionId", "")
	_ = session.Save()
	c.Redirect(http.StatusMovedPermanently, "/login")
}

func main() {
	//gin.SetMode(gin.ReleaseMode)

	key := securecookie.GenerateRandomKey(32)
	keyUsed := hex.EncodeToString(key)
	log.Println("Key used:", keyUsed)

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	store := cookie.NewStore(key)
	router.Use(sessions.Sessions(SessionName, store))
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
	router.GET("/articles/id/:id/publish", articlesIdPublish)
	router.GET("/articles/id/:id/unpublish", articlesIdUnpublish)
	router.GET("/articles/id/:id/delete", articlesIdDelete)
	err := router.Run(":9090")
	if err != nil {
		log.Fatal(err)
	}
}
