package main

import (
	"encoding/json"
	"net/http"
)

func (s *Server) handleEstimate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		errResp(w, "method not allowed", 405)
		return
	}
	var req struct {
		Servers    int `json:"servers"`
		StorageTB  int `json:"storage_tb"`
		Bandwidth  int `json:"bandwidth_gbps"`
		Redundancy int `json:"redundancy"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	serverCost := float64(req.Servers) * 490000.0
	storageCost := float64(req.StorageTB) * 12000.0
	networkCost := float64(req.Bandwidth) * 8000.0
	redundancyMult := 1.0 + float64(req.Redundancy-1)*0.3
	infra := (serverCost + storageCost + networkCost) * redundancyMult
	cooling := infra * 0.15
	power := infra * 0.12
	security := 250000.0
	installation := infra * 0.08

	type LineItem struct {
		Name string  `json:"name"`
		Cost float64 `json:"cost"`
	}

	resp := struct {
		Items []LineItem `json:"items"`
		Total float64    `json:"total"`
	}{
		Items: []LineItem{
			{"Серверное оборудование", serverCost * redundancyMult},
			{"Системы хранения данных", storageCost * redundancyMult},
			{"Сетевое оборудование", networkCost * redundancyMult},
			{"Системы охлаждения", cooling},
			{"Системы электропитания", power},
			{"Системы безопасности", security},
			{"Монтаж и пусконаладка", installation},
		},
	}
	for _, item := range resp.Items {
		resp.Total += item.Cost
	}
	jsonResp(w, resp, 200)
}
