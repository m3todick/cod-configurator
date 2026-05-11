package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (s *Server) handleConfigurations(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		user := s.getUserFromCookie(r)
		if user == nil {
			errResp(w, "unauthorized", 401)
			return
		}
		list := s.store.getConfigurations(user.ID, user.Role)
		if list == nil {
			list = []*Configuration{}
		}
		jsonResp(w, list, 200)
	case "POST":
		user := s.getUserFromCookie(r)
		if user == nil {
			errResp(w, "unauthorized", 401)
			return
		}
		var cfg Configuration
		json.NewDecoder(r.Body).Decode(&cfg)
		cfg.ID = genID()
		cfg.UserID = user.ID
		cfg.TotalCost = s.store.calcCost(cfg.Items)
		if err := s.store.createConfiguration(&cfg); err != nil {
			errResp(w, "Ошибка создания конфигурации", 500)
			return
		}
		saved := s.store.getConfiguration(cfg.ID)
		if saved == nil {
			saved = &cfg
		}
		jsonResp(w, saved, 201)
	}
}

func (s *Server) handleExport(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/configurations/export/")
	cfg := s.store.getConfiguration(id)
	if cfg == nil {
		errResp(w, "not found", 404)
		return
	}
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="config_%s.json"`, id))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cfg)
}

func (s *Server) handleImport(w http.ResponseWriter, r *http.Request) {
	user := s.getUserFromCookie(r)
	if user == nil {
		errResp(w, "unauthorized", 401)
		return
	}
	var cfg Configuration
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		errResp(w, "invalid json", 400)
		return
	}
	cfg.ID = genID()
	cfg.UserID = user.ID
	cfg.TotalCost = s.store.calcCost(cfg.Items)
	if err := s.store.createConfiguration(&cfg); err != nil {
		errResp(w, "Ошибка импорта конфигурации", 500)
		return
	}
	saved := s.store.getConfiguration(cfg.ID)
	if saved == nil {
		saved = &cfg
	}
	jsonResp(w, saved, 201)
}

func (s *Server) handleConfiguration(w http.ResponseWriter, r *http.Request) {
	p := strings.TrimPrefix(r.URL.Path, "/api/configurations/")
	id := strings.Split(p, "/")[0]

	switch r.Method {
	case "GET":
		c := s.store.getConfiguration(id)
		if c == nil {
			errResp(w, "not found", 404)
			return
		}
		jsonResp(w, c, 200)
	case "PUT":
		user := s.getUserFromCookie(r)
		if user == nil {
			errResp(w, "unauthorized", 401)
			return
		}
		cfg := s.store.getConfiguration(id)
		if cfg == nil {
			errResp(w, "not found", 404)
			return
		}
		json.NewDecoder(r.Body).Decode(cfg)
		cfg.ID = id
		cfg.TotalCost = s.store.calcCost(cfg.Items)
		if err := s.store.updateConfiguration(cfg); err != nil {
			errResp(w, "Ошибка обновления", 500)
			return
		}
		jsonResp(w, cfg, 200)
	case "DELETE":
		user := s.getUserFromCookie(r)
		if user == nil {
			errResp(w, "unauthorized", 401)
			return
		}
		s.store.deleteConfiguration(id)
		jsonResp(w, map[string]string{"ok": "true"}, 200)
	}
}
