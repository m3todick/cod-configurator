package main

import (
        "encoding/json"
        "net/http"
        "os"
        "strings"
)

func (s *Server) handleProfile(w http.ResponseWriter, r *http.Request) {
        user := s.getUserFromCookie(r)
        if user == nil {
                errResp(w, "unauthorized", 401)
                return
        }
        if r.Method != "PUT" {
                errResp(w, "method not allowed", 405)
                return
        }
        var req struct {
                Name         string `json:"name"`
                Email        string `json:"email"`
                Organization string `json:"organization"`
                OldPassword  string `json:"old_password"`
                NewPassword  string `json:"new_password"`
        }
        json.NewDecoder(r.Body).Decode(&req)

        if req.NewPassword != "" {
                if !checkPassword(req.OldPassword, user.Password) {
                        errResp(w, "Неверный текущий пароль", 400)
                        return
                }
                user.Password = hashPassword(req.NewPassword)
        }
        if req.Name != "" {
                user.Name = req.Name
        }
        if req.Organization != "" {
                user.Organization = req.Organization
        }
        if req.Email != "" && req.Email != user.Email {
                existing := s.store.getUserByEmail(req.Email)
                if existing != nil && existing.ID != user.ID {
                        errResp(w, "Email уже занят другим пользователем", 409)
                        return
                }
                user.Email = req.Email
        }
        if err := s.store.updateUser(user); err != nil {
                errResp(w, "Ошибка обновления профиля", 500)
                return
        }
        jsonResp(w, user, 200)
}

func (s *Server) handleUsers(w http.ResponseWriter, r *http.Request) {
        admin := s.getUserFromCookie(r)
        if admin == nil || admin.Role != "admin" {
                errResp(w, "forbidden", 403)
                return
        }
        if r.Method != "GET" {
                errResp(w, "method not allowed", 405)
                return
        }
        list := s.store.getUsers()
        if list == nil {
                list = []*User{}
        }
        jsonResp(w, list, 200)
}

func (s *Server) handleUserByID(w http.ResponseWriter, r *http.Request) {
        admin := s.getUserFromCookie(r)
        if admin == nil || admin.Role != "admin" {
                errResp(w, "forbidden", 403)
                return
        }
        id := strings.TrimPrefix(r.URL.Path, "/api/users/")

        switch r.Method {
        case "GET":
                found := s.store.getUserByID(id)
                if found == nil {
                        errResp(w, "not found", 404)
                        return
                }
                jsonResp(w, found, 200)
        case "PUT":
                var req struct {
                        Name         string `json:"name"`
                        Email        string `json:"email"`
                        Organization string `json:"organization"`
                        Password     string `json:"password"`
                        Role         string `json:"role"`
                        RolePassword string `json:"role_password"`
                }
                json.NewDecoder(r.Body).Decode(&req)
                found := s.store.getUserByID(id)
                if found == nil {
                        errResp(w, "not found", 404)
                        return
                }
                if req.Email != "" && req.Email != found.Email {
                        existing := s.store.getUserByEmail(req.Email)
                        if existing != nil && existing.ID != id {
                                errResp(w, "Email уже занят", 409)
                                return
                        }
                        found.Email = req.Email
                }
                if req.Name != "" {
                        found.Name = req.Name
                }
                if req.Organization != "" {
                        found.Organization = req.Organization
                }
                if req.Password != "" {
                        found.Password = hashPassword(req.Password)
                }
                if (req.Role == "admin" || req.Role == "user") && req.Role != found.Role {
                        secret := strings.TrimSpace(os.Getenv("ROLE_CHANGE_PASSWORD"))
                        if secret == "" {
                                secret = "admin"
                        }
                        if strings.TrimSpace(req.RolePassword) != secret {
                                errResp(w, "Неверный пароль для смены роли", 403)
                                return
                        }
                        found.Role = req.Role
                }
                if err := s.store.updateUser(found); err != nil {
                        errResp(w, "Ошибка обновления пользователя", 500)
                        return
                }
                jsonResp(w, found, 200)
        case "DELETE":
                if id == admin.ID {
                        errResp(w, "Нельзя удалить собственный аккаунт", 400)
                        return
                }
                found := s.store.getUserByID(id)
                if found == nil {
                        errResp(w, "not found", 404)
                        return
                }
                if _, err := s.store.db.Exec(`DELETE FROM users WHERE id=$1`, id); err != nil {
                        errResp(w, "Ошибка удаления пользователя", 500)
                        return
                }
                jsonResp(w, map[string]string{"status": "ok"}, 200)
        default:
                errResp(w, "method not allowed", 405)
        }
}
