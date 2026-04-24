package main

import (
	"encoding/json"
	"net/http"
	"time"
)

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		errResp(w, "method not allowed", 405)
		return
	}
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	user := s.store.getUserByEmail(req.Email)
	if user == nil || !checkPassword(req.Password, user.Password) {
		errResp(w, "Неверный email или пароль", 401)
		return
	}

	token := genID()
	s.store.createSession(token, user.ID, time.Now().Add(24*time.Hour))
	http.SetCookie(w, &http.Cookie{Name: "session", Value: token, Path: "/", MaxAge: 86400})
	jsonResp(w, user, 200)
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session")
	if err == nil {
		s.store.deleteSession(c.Value)
	}
	http.SetCookie(w, &http.Cookie{Name: "session", Value: "", MaxAge: -1, Path: "/"})
	jsonResp(w, map[string]string{"ok": "true"}, 200)
}

func (s *Server) handleMe(w http.ResponseWriter, r *http.Request) {
	user := s.getUserFromCookie(r)
	if user == nil {
		errResp(w, "unauthorized", 401)
		return
	}
	jsonResp(w, user, 200)
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		errResp(w, "method not allowed", 405)
		return
	}
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	if req.Name == "" || req.Email == "" || len(req.Password) < 6 {
		errResp(w, "Заполните все поля. Пароль — минимум 6 символов", 400)
		return
	}

	if s.store.getUserByEmail(req.Email) != nil {
		errResp(w, "Пользователь с таким email уже зарегистрирован", 409)
		return
	}

	newUser := &User{
		ID:           "u" + genID(),
		Name:         req.Name,
		Email:        req.Email,
		Password:     hashPassword(req.Password),
		Role:         "user",
		Organization: "Администрация Константиновского района",
	}
	if err := s.store.createUser(newUser); err != nil {
		errResp(w, "Ошибка при регистрации", 500)
		return
	}

	token := genID()
	s.store.createSession(token, newUser.ID, time.Now().Add(24*time.Hour))
	http.SetCookie(w, &http.Cookie{Name: "session", Value: token, Path: "/", MaxAge: 86400})
	jsonResp(w, newUser, 201)
}
