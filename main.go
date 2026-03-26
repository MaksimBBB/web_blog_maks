package main

import (
	"fmt"
	"net/http"
	"os"

	api "web-blog/handlers"

	"github.com/subosito/gotenv"
)

func init() {
	_ = gotenv.Load()
	if _, err := os.Stat("articles"); os.IsNotExist(err) {
		os.Mkdir("articles", 0755)
	}
}

func main() {
	http.HandleFunc("/", api.HomeHandler)
	http.HandleFunc("/articles", api.GetArticlesHandler)
	http.HandleFunc("/article/", api.ArticlePageHandler)
	http.HandleFunc("/dashboard", api.DashboardArticleWithAuthI())
	http.HandleFunc("/articles/new", api.CreateArticleWithAuthI())
	http.HandleFunc("/articles/update/", api.UpdateArticleWithAuth())
	http.HandleFunc("/articles/delete/", api.DeleteArticleWithAuth())

	http.HandleFunc("/login", api.LoginHandler)
	http.HandleFunc("/logout", api.LogoutHandler)

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
