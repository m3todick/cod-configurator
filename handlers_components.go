package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (s *Server) handleComponents(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		cat := r.URL.Query().Get("category")
		list := s.store.getComponents(cat)
		if list == nil {
			list = []*Component{}
		}
		jsonResp(w, list, 200)
	case "POST":
		user := s.getUserFromCookie(r)
		if user == nil || user.Role != "admin" {
			errResp(w, "forbidden", 403)
			return
		}
		var comp Component
		json.NewDecoder(r.Body).Decode(&comp)
		comp.ID = "c" + genID()
		if err := s.store.createComponent(&comp); err != nil {
			errResp(w, "Ошибка создания компонента", 500)
			return
		}
		jsonResp(w, comp, 201)
	}
}

func (s *Server) handleComponent(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/components/")
	switch r.Method {
	case "GET":
		c := s.store.getComponent(id)
		if c == nil {
			errResp(w, "not found", 404)
			return
		}
		jsonResp(w, c, 200)
	case "PUT":
		user := s.getUserFromCookie(r)
		if user == nil || user.Role != "admin" {
			errResp(w, "forbidden", 403)
			return
		}
		var comp Component
		json.NewDecoder(r.Body).Decode(&comp)
		comp.ID = id
		if err := s.store.updateComponent(&comp); err != nil {
			errResp(w, "Ошибка обновления", 500)
			return
		}
		jsonResp(w, comp, 200)
	case "DELETE":
		user := s.getUserFromCookie(r)
		if user == nil || user.Role != "admin" {
			errResp(w, "forbidden", 403)
			return
		}
		s.store.deleteComponent(id)
		jsonResp(w, map[string]string{"ok": "true"}, 200)
	}
}
