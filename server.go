package main

import (
        "encoding/json"
        "net/http"
)

type Server struct {
        store *Store
        mux   *http.ServeMux
}

func NewServer(store *Store) *Server {
        s := &Server{store: store, mux: http.NewServeMux()}
        s.routes()
        return s
}

func (s *Server) routes() {
        fs := http.FileServer(http.Dir("web/static"))
        s.mux.Handle("/static/", http.StripPrefix("/static/", noCache(fs)))

        s.mux.HandleFunc("/", s.pageHandler("web/templates/index.html"))
        s.mux.HandleFunc("/login", s.pageHandler("web/templates/login.html"))
        s.mux.HandleFunc("/register", s.pageHandler("web/templates/register.html"))
        s.mux.HandleFunc("/cabinet", s.pageHandler("web/templates/cabinet.html"))
        s.mux.HandleFunc("/configurator", s.pageHandler("web/templates/configurator.html"))
        s.mux.HandleFunc("/calculator", s.pageHandler("web/templates/calculator.html"))
        s.mux.HandleFunc("/admin", s.pageHandler("web/templates/admin.html"))
        s.mux.HandleFunc("/terms", s.pageHandler("web/templates/terms.html"))
        s.mux.HandleFunc("/privacy", s.pageHandler("web/templates/privacy.html"))

        s.mux.HandleFunc("/api/auth/login", s.handleLogin)
        s.mux.HandleFunc("/api/auth/logout", s.handleLogout)
        s.mux.HandleFunc("/api/auth/me", s.handleMe)
        s.mux.HandleFunc("/api/auth/register", s.handleRegister)

        s.mux.HandleFunc("/api/components", s.handleComponents)
        s.mux.HandleFunc("/api/components/", s.handleComponent)

        s.mux.HandleFunc("/api/configurations", s.handleConfigurations)
        s.mux.HandleFunc("/api/configurations/export/", s.handleExport)
        s.mux.HandleFunc("/api/configurations/import", s.handleImport)
        s.mux.HandleFunc("/api/configurations/", s.handleConfiguration)

        s.mux.HandleFunc("/api/calculator/estimate", s.handleEstimate)
        s.mux.HandleFunc("/api/profile", s.handleProfile)
        s.mux.HandleFunc("/api/users", s.handleUsers)
        s.mux.HandleFunc("/api/users/", s.handleUserByID)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        if r.Method == "OPTIONS" {
                w.WriteHeader(204)
                return
        }
        s.mux.ServeHTTP(w, r)
}

func (s *Server) pageHandler(path string) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
                w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
                w.Header().Set("Pragma", "no-cache")
                w.Header().Set("Expires", "0")
                http.ServeFile(w, r, path)
        }
}

func noCache(h http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
                w.Header().Set("Pragma", "no-cache")
                w.Header().Set("Expires", "0")
                h.ServeHTTP(w, r)
        })
}

func jsonResp(w http.ResponseWriter, data interface{}, code int) {
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(code)
        json.NewEncoder(w).Encode(data)
}

func errResp(w http.ResponseWriter, msg string, code int) {
        jsonResp(w, map[string]string{"error": msg}, code)
}

func (s *Server) getUserFromCookie(r *http.Request) *User {
        c, err := r.Cookie("session")
        if err != nil {
                return nil
        }
        return s.store.getSessionUser(c.Value)
}
