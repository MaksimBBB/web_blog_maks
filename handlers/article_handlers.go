package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"
	middleware "web-blog/handlers/middleware"
	"web-blog/model"
)

func parseTemplate(templateName string) *template.Template {
	path := fmt.Sprintf("templates/%s", templateName)

	tmpl, err := template.ParseFiles(path)
	if err != nil {
		panic(fmt.Sprintf("error parsing template %s, %v", templateName, err))
	}

	return tmpl
}

func createArticle(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := parseTemplate("newArticle.html")
		tmpl.Execute(w, nil)
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusNotFound)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	published := r.FormValue("date")

	if title == "" || len(title) > 100 {
		http.Error(w, `{"error":"title required and  <100 chars"}`, http.StatusBadRequest)
		return
	}

	if content == "" {
		http.Error(w, `{"error":"content required"}`, http.StatusBadRequest)
		return
	}

	if _, err := time.Parse("2006-01-02", published); err != nil {
		http.Error(w, `{"error":"invalid date format YYYY-MM-DD"}`, http.StatusBadRequest)
		return
	}

	files, err := os.ReadDir("articles")
	if err != nil {
		http.Error(w, `{"error":"failed to read articles directory"}`, http.StatusInternalServerError)
		return
	}
	maxID := 0

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			var id int
			fmt.Sscanf(file.Name(), "article%d.json", &id)
			if id > maxID {
				maxID = id
			}
		}
	}

	var a model.Article
	a.Title = title
	a.Content = content
	a.Published = published
	a.Author = user.Username
	a.ID = maxID + 1

	filePath := fmt.Sprintf("articles/article%d.json", maxID+1)
	file, err := os.Create(filePath)
	if err != nil {
		http.Error(w, `{"error":"failed to create article file"}`, http.StatusInternalServerError)
		return
	}

	defer file.Close()
	if err := json.NewEncoder(file).Encode(a); err != nil {
		http.Error(w, `{"error":"failed to save article"}`, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func getArticles() []model.Article {
	files, err := os.ReadDir("articles")
	if err != nil {
		return nil
	}
	var articles []model.Article

	for _, f := range files {
		if filepath.Ext(f.Name()) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join("articles", f.Name()))
		if err != nil {
			continue
		}
		var art model.Article
		if err := json.Unmarshal(data, &art); err != nil {
			continue
		}
		articles = append(articles, art)
	}
	return articles
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	articles := getArticles()
	tmpl := parseTemplate("home.html")
	tmpl.Execute(w, articles)
}

func GetArticlesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Query().Get("format") == "json" {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(getArticles()); err != nil {
			http.Error(w, `{"error":"failed to encode response"}`, http.StatusInternalServerError)
			return
		}
		return
	}

	tmpl := parseTemplate("home.html")
	tmpl.Execute(w, getArticles())
}

// get article по id
func getArticleById(id int) (*model.Article, error) {
	filePath := fmt.Sprintf("articles/article%d.json", id)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var a model.Article
	err = json.Unmarshal(data, &a)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	articles := getArticles()
	tmpl := parseTemplate("dashboard.html")
	tmpl.Execute(w, articles)
}

func updateArticle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/articles/update/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodGet {
		article, err := getArticleById(id)
		if err != nil {
			http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
			return
		}

		tmpl := parseTemplate("updateArticle.html")
		tmpl.Execute(w, article)
		return
	}

	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(w, `{"error":"unauthorized"}`, http.StatusNotFound)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	published := r.FormValue("date")

	if title == "" || len(title) > 100 {
		http.Error(w, `{"error":"title required and <100 chars"}`, http.StatusBadRequest)
		return
	}

	if content == "" {
		http.Error(w, `{"error": "content required"}`, http.StatusBadRequest)
		return
	}

	if _, err := time.Parse("2006-01-02", published); err != nil {
		http.Error(w, `{"error": "invalid date format YYYY-MM-DD"}`, http.StatusBadRequest)
		return
	}

	var a model.Article

	a.Title = title
	a.Content = content
	a.Published = published
	a.Author = user.Username
	a.ID = id

	filePath := fmt.Sprintf("articles/article%d.json", id)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}

	file, err := os.Create(filePath)
	if err != nil {
		http.Error(w, `{"error":"failed to update article"}`, http.StatusInternalServerError)
		return
	}
	defer file.Close()
	if err := json.NewEncoder(file).Encode(a); err != nil {
		http.Error(w, `{"error":"failed to update article"}`, http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// Delete
func deleteArticle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/articles/delete/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	filePath := fmt.Sprintf("articles/article%d.json", id)
	if err := os.Remove(filePath); err != nil {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)

}

func ArticlePageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/article/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	article, err := getArticleById(id)
	if err != nil {
		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
		return
	}

	tmpl := parseTemplate("articlepage.html")
	tmpl.Execute(w, article)
}

func CreateArticleWithAuthI() http.HandlerFunc {
	return middleware.CookieAuthMiddleware(createArticle)
}

func DashboardArticleWithAuthI() http.HandlerFunc {
	return middleware.CookieAuthMiddleware(dashboardHandler)
}

func UpdateArticleWithAuth() http.HandlerFunc {
	return middleware.CookieAuthMiddleware(updateArticle)
}

func DeleteArticleWithAuth() http.HandlerFunc {
	return middleware.CookieAuthMiddleware(deleteArticle)
}
