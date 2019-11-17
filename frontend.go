package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Article struct {
	Title         string
	Content       string
	Author        string
	TimePublished string
}

func home(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func about(c *gin.Context) {
	c.HTML(http.StatusOK, "about.html", nil)
}

func contactUs(c *gin.Context) {
	c.HTML(http.StatusOK, "contact.html", nil)
}

func login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func articles(c *gin.Context) {
	c.HTML(http.StatusOK, "articles.html", nil)
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
	router.GET("/articles", articles)
	router.GET("/logout", logout)
	err := router.Run(":9090")
	if err != nil {
		log.Fatal(err)
	}
}
