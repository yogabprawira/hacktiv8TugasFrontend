package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type Post struct {
	Id            int    `json:"id"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	Author        string `json:"author"`
	TimePublished string `json:"time_published"`
	IsPublished   int    `json:"is_published"`
}

type Posts struct {
	Posts []Post `json:"posts"`
}

func home(c *gin.Context) {
	var posts Posts
	data := map[string]interface{}{}
	posts.Posts = make([]Post, 0)

	var post Post
	post.Id = 2
	post.Title = "Hahahaha"
	post.Author = "Yoga"
	post.Content = "Hihihihihi"
	post.IsPublished = 1
	post.TimePublished = time.Now().Format("1 January 2006")
	posts.Posts = append(posts.Posts, post)
	post.Id = 3
	post.Title = "dfafdafdafa"
	post.Author = "Yoga"
	post.Content = "dafafdafadsfa"
	post.IsPublished = 1
	post.TimePublished = time.Now().Format("1 January 2006")
	posts.Posts = append(posts.Posts, post)

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

func login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func loginPost(c *gin.Context) {

}

func articles(c *gin.Context) {
	data := map[string]interface{}{}

	var posts Posts
	var post Post
	post.Id = 2
	post.Title = "Hahahaha"
	post.Author = "Yoga"
	post.Content = "Hihihihihi"
	post.IsPublished = 1
	post.TimePublished = time.Now().Format("1 January 2006")
	posts.Posts = append(posts.Posts, post)
	post.Id = 3
	post.Title = "dfafdafdafa"
	post.Author = "Yoga"
	post.Content = "dafafdafadsfa"
	post.IsPublished = 1
	post.TimePublished = time.Now().Format("1 January 2006")
	posts.Posts = append(posts.Posts, post)

	data["posts"] = posts
	data["user"] = "yoga"
	c.HTML(http.StatusOK, "articles.html", data)
}

func articlesId(c *gin.Context) {
	data := map[string]interface{}{}
	var post Post
	post.Id = 2
	post.Title = "Hahahaha"
	post.Author = "Yoga"
	post.Content = "Hihihihihi"
	post.IsPublished = 1
	post.TimePublished = time.Now().Format("1 January 2006")
	data["post"] = post
	data["user"] = "yoga"
	c.HTML(http.StatusOK, "articlesId.html", data)
}

func articlesAdd(c *gin.Context) {
	data := map[string]interface{}{}
	data["user"] = "yoga"
	c.HTML(http.StatusOK, "articlesId.html", data)
}

func articlesIdPost(c *gin.Context) {

}

func logout(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "/login")
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/", home)
	router.GET("/about", about)
	router.GET("/contact", contactUs)
	router.GET("/login", login)
	router.POST("/login", loginPost)
	router.GET("/articles", articles)
	router.GET("/articles/id/:id", articlesId)
	router.POST("/articles/id/:id", articlesIdPost)
	router.GET("/articles/add", articlesAdd)
	router.GET("/logout", logout)
	err := router.Run(":9090")
	if err != nil {
		log.Fatal(err)
	}
}
