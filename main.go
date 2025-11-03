package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
)

// render sends the specified HTML file through Go's html/template engine.
// Files are expected to live under the view/ directory.
func render(w http.ResponseWriter, filename string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Safety: only allow .html files and resolve relative to view/
	clean := filepath.Clean(filename)
	if filepath.Ext(clean) != ".html" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if strings.Contains(clean, "..") || filepath.IsAbs(clean) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	// If the template path is under pages/, render with base layout
	fullPath := filepath.Join("view", clean)
	if strings.HasPrefix(clean, "pages/") {
		base := filepath.Join("view", "layout", "base.html")
		if _, err := os.Stat(fullPath); err == nil {
			tmpl, err := template.ParseFiles(base, fullPath)
			if err != nil {
				log.Printf("template parse error for %s: %v", fullPath, err)
				http.Error(w, "template error", http.StatusInternalServerError)
				return
			}
			if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
				log.Printf("template execute error for %s: %v", fullPath, err)
				http.Error(w, "render error", http.StatusInternalServerError)
				return
			}
			return
		}
	}

	// Fallback: render standalone file under view/
	tmpl, err := template.ParseFiles(fullPath)
	if err != nil {
		log.Printf("template parse error for %s: %v", fullPath, err)
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("template execute error for %s: %v", fullPath, err)
		http.Error(w, "render error", http.StatusInternalServerError)
		return
	}
}

func main() {
	r := mux.NewRouter()

	// Basic Auth middleware (applies to all routes)
	r.Use(authMiddleware)

	// Static files under /public/
	fileServer := http.FileServer(http.Dir("public"))
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", fileServer))

	// Routes mapping to existing HTML files
	r.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		render(w, "pages/index.html", nil)
	})

	r.HandleFunc("/about-us", func(w http.ResponseWriter, _ *http.Request) {
		render(w, "pages/about-us.html", nil)
	})

	r.HandleFunc("/services", func(w http.ResponseWriter, _ *http.Request) {
		render(w, "pages/our-services.html", nil)
	})

	r.HandleFunc("/training", func(w http.ResponseWriter, _ *http.Request) {
		render(w, "pages/training.html", nil)
	})

	r.HandleFunc("/blog", func(w http.ResponseWriter, _ *http.Request) {
		render(w, "pages/bloglisting.html", nil)
	})

	// Example dynamic detail route using same template (you can personalize later)
	r.HandleFunc("/blog/{slug}", func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		data := map[string]any{
			"Slug": vars["slug"],
		}
		render(w, "pages/blogDetails.html", data)
	})

	r.HandleFunc("/contact", func(w http.ResponseWriter, _ *http.Request) {
		render(w, "pages/contact_us.html", nil)
	})

	// Linux commands reference page (uses layout)
	r.HandleFunc("/linux-commands", func(w http.ResponseWriter, _ *http.Request) {
		render(w, "pages/linux-commands.html", nil)
	})

	// Linux directory structure page
	r.HandleFunc("/linux-directory-structure", func(w http.ResponseWriter, _ *http.Request) {
		render(w, "pages/linux-directory-structure.html", nil)
	})

	// Linux permissions and user management page
	r.HandleFunc("/linux-permissions", func(w http.ResponseWriter, _ *http.Request) {
		render(w, "pages/linux-permissions.html", nil)
	})

	// Golang project structure page
	r.HandleFunc("/golang-project-structure", func(w http.ResponseWriter, _ *http.Request) {
		render(w, "pages/golang-project-structure.html", nil)
	})

	// Golang create project tutorial page
	r.HandleFunc("/golang-create-project", func(w http.ResponseWriter, _ *http.Request) {
		render(w, "pages/golang-create-project.html", nil)
	})

	// Optional: if you want to expose server.html on /server
	r.HandleFunc("/server", func(w http.ResponseWriter, _ *http.Request) {
		render(w, "pages/server.html", nil)
	})

	// Under development page (standalone, no layout)
	r.HandleFunc("/under-development", func(w http.ResponseWriter, _ *http.Request) {
		render(w, "under-development.html", nil)
	})

	// Roadmaps
	r.HandleFunc("/golang", func(w http.ResponseWriter, _ *http.Request) {
		render(w, "pages/roadmaps/golang-roadmap.html", nil)
	})
	r.HandleFunc("/devops", func(w http.ResponseWriter, _ *http.Request) {
		render(w, "pages/roadmaps/devops-roadmap.html", nil)
	})
	r.HandleFunc("/project-manager", func(w http.ResponseWriter, _ *http.Request) {
		render(w, "pages/project-manager-roadmap.html", nil)
	})
	r.HandleFunc("/ai-ml", func(w http.ResponseWriter, _ *http.Request) {
		render(w, "pages/ai-ml-roadmap.html", nil)
	})

	srv := &http.Server{
		Addr:    ":9090",
		Handler: r,
	}

	log.Printf("listening on http://localhost%s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

// authMiddleware enforces HTTP Basic authentication on all requests.
// Configure credentials via env: BASIC_USER, BASIC_PASS. Defaults to admin/admin.
func authMiddleware(next http.Handler) http.Handler {
	expectedUser := os.Getenv("BASIC_USER")
	if expectedUser == "" {
		expectedUser = "admin"
	}
	expectedPass := os.Getenv("BASIC_PASS")
	if expectedPass == "" {
		expectedPass = "0987654321"
	}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		user, pass, ok := req.BasicAuth()
		if !ok || user != expectedUser || pass != expectedPass {
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, req)
	})
}
