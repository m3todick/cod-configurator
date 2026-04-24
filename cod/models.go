package main

import "time"

type Category string

const (
	CatServer   Category = "server"
	CatStorage  Category = "storage"
	CatNetwork  Category = "network"
	CatCooling  Category = "cooling"
	CatPower    Category = "power"
	CatSecurity Category = "security"
)

type Component struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Category    Category          `json:"category"`
	Brand       string            `json:"brand"`
	Model       string            `json:"model"`
	Price       float64           `json:"price"`
	Description string            `json:"description"`
	Specs       map[string]string `json:"specs"`
	InStock     bool              `json:"in_stock"`
	Image       string            `json:"image"`
}

type ConfigItem struct {
	ComponentID string `json:"component_id"`
	Quantity    int    `json:"quantity"`
}

type Configuration struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	UserID      string       `json:"user_id"`
	Items       []ConfigItem `json:"items"`
	TotalCost   float64      `json:"total_cost"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Description string       `json:"description"`
}

type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Password     string    `json:"-"`
	Role         string    `json:"role"`
	Organization string    `json:"organization"`
	CreatedAt    time.Time `json:"created_at"`
}
