package main

import (
        "database/sql"
        "encoding/json"
        "time"
)

type Store struct {
        db *sql.DB
}

func NewStore(db *sql.DB) *Store {
        s := &Store{db: db}
        s.seed()
        return s
}

func (s *Store) seed() {
        var count int
        s.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
        if count > 0 {
                return
        }

        users := []struct {
                User
                plain string
        }{
                {User: User{ID: "u1", Name: "Администратор", Email: "admin@konst-adm.ru", Role: "admin", Organization: "Администрация Константиновского района"}, plain: "admin123"},
                {User: User{ID: "u2", Name: "Иван Петров", Email: "user@konst-adm.ru", Role: "user", Organization: "Администрация Константиновского района"}, plain: "user123"},
        }
        for _, u := range users {
                s.db.Exec(`INSERT INTO users (id, name, email, password, role, organization) VALUES ($1,$2,$3,$4,$5,$6) ON CONFLICT DO NOTHING`,
                        u.ID, u.Name, u.Email, hashPassword(u.plain), u.Role, u.Organization)
        }

        components := []Component{
                {ID: "c1", Name: "Сервер HP ProLiant DL380 Gen10", Category: CatServer, Brand: "HP", Model: "DL380 Gen10",
                        Price: 485000, Description: "2U сервер для высоких нагрузок с поддержкой до 3TB RAM",
                        Specs:   map[string]string{"CPU": "2x Intel Xeon Gold 6226R", "RAM": "64GB DDR4", "HDD": "2x 600GB SAS", "Блоки питания": "2x 800W"},
                        InStock: true, Image: "/static/img/server1.svg"},
                {ID: "c2", Name: "Сервер Dell PowerEdge R740", Category: CatServer, Brand: "Dell", Model: "R740",
                        Price: 520000, Description: "2U двухпроцессорный сервер для критических приложений",
                        Specs:   map[string]string{"CPU": "2x Intel Xeon Silver 4210R", "RAM": "128GB DDR4", "HDD": "4x 1.2TB SAS", "Блоки питания": "2x 750W"},
                        InStock: true, Image: "/static/img/server2.svg"},
                {ID: "c3", Name: "Сервер Lenovo ThinkSystem SR650", Category: CatServer, Brand: "Lenovo", Model: "SR650",
                        Price: 445000, Description: "Универсальный 2U сервер для виртуализации",
                        Specs:   map[string]string{"CPU": "2x Intel Xeon Gold 5218", "RAM": "96GB DDR4", "HDD": "8x 600GB SAS", "Блоки питания": "2x 900W"},
                        InStock: false, Image: "/static/img/server1.svg"},
                {ID: "c4", Name: "СХД NetApp AFF A250", Category: CatStorage, Brand: "NetApp", Model: "AFF A250",
                        Price: 1250000, Description: "All-Flash массив для высокопроизводительных рабочих нагрузок",
                        Specs:   map[string]string{"Ёмкость": "100TB", "Тип": "All-Flash NVMe", "Интерфейс": "25GbE, 32Gb FC", "IOPS": "600K"},
                        InStock: true, Image: "/static/img/storage.svg"},
                {ID: "c5", Name: "СХД Dell EMC PowerStore 500T", Category: CatStorage, Brand: "Dell EMC", Model: "PowerStore 500T",
                        Price: 980000, Description: "Гибридный массив с интеллектуальным управлением данными",
                        Specs:   map[string]string{"Ёмкость": "50TB", "Тип": "Гибридный", "Интерфейс": "10GbE, 16Gb FC", "IOPS": "300K"},
                        InStock: true, Image: "/static/img/storage.svg"},
                {ID: "c6", Name: "Коммутатор Cisco Nexus 9300", Category: CatNetwork, Brand: "Cisco", Model: "Nexus 9300",
                        Price: 380000, Description: "Высокопроизводительный коммутатор для ЦОД",
                        Specs:   map[string]string{"Порты": "48x 10GbE + 6x 100GbE", "Скорость": "10/40/100GbE", "Буфер": "40MB", "Задержка": "< 1мкс"},
                        InStock: true, Image: "/static/img/network.svg"},
                {ID: "c7", Name: "Коммутатор Huawei CloudEngine CE6870", Category: CatNetwork, Brand: "Huawei", Model: "CE6870",
                        Price: 295000, Description: "Коммутатор уровня доступа/агрегации для ЦОД",
                        Specs:   map[string]string{"Порты": "48x 25GbE + 8x 100GbE", "Пропускная способность": "6.4Tbps", "Задержка": "< 2мкс"},
                        InStock: true, Image: "/static/img/network.svg"},
                {ID: "c8", Name: "Маршрутизатор Cisco ASR 1002-HX", Category: CatNetwork, Brand: "Cisco", Model: "ASR 1002-HX",
                        Price: 520000, Description: "Граничный маршрутизатор для ЦОД с функциями безопасности",
                        Specs:   map[string]string{"Пропускная способность": "до 60Gbps", "Слоты": "8x SFP+", "ОЗУ": "16GB", "ОС": "IOS XE"},
                        InStock: true, Image: "/static/img/network.svg"},
                {ID: "c9", Name: "ИБП APC Smart-UPS VT 60kVA", Category: CatPower, Brand: "APC", Model: "Smart-UPS VT 60kVA",
                        Price: 890000, Description: "Источник бесперебойного питания для ЦОД",
                        Specs:   map[string]string{"Мощность": "60kVA / 54kW", "Время автономии": "до 8 мин при полной нагрузке", "КПД": "до 99%", "Форм-фактор": "Напольный"},
                        InStock: true, Image: "/static/img/ups.svg"},
                {ID: "c10", Name: "Прецизионный кондиционер Emerson Liebert PEX", Category: CatCooling, Brand: "Emerson", Model: "Liebert PEX",
                        Price: 650000, Description: "Прецизионная система охлаждения для серверных помещений",
                        Specs:   map[string]string{"Мощность охлаждения": "20kW", "Тип": "Прецизионный", "Хладагент": "R410A", "Управление": "iCOM"},
                        InStock: true, Image: "/static/img/cooling.svg"},
                {ID: "c11", Name: "Система холодных коридоров APC", Category: CatCooling, Brand: "APC", Model: "InRow RC",
                        Price: 420000, Description: "Система изоляции холодных/горячих коридоров",
                        Specs:   map[string]string{"Мощность охлаждения": "15kW", "Тип": "Рядное охлаждение", "Подача": "Фронтальная"},
                        InStock: false, Image: "/static/img/cooling.svg"},
                {ID: "c12", Name: "СКУД Bolid С2000-М", Category: CatSecurity, Brand: "Болид", Model: "С2000-М",
                        Price: 125000, Description: "Система контроля и управления доступом для ЦОД",
                        Specs:   map[string]string{"Считыватели": "до 512", "Контроллеры": "до 127", "Интерфейс": "RS-485", "ОС": "Встроенная"},
                        InStock: true, Image: "/static/img/security.svg"},
                {ID: "c13", Name: "Видеонаблюдение Hikvision DS-2CD2T47G2", Category: CatSecurity, Brand: "Hikvision", Model: "DS-2CD2T47G2",
                        Price: 85000, Description: "Комплект IP-видеонаблюдения для серверного помещения",
                        Specs:   map[string]string{"Разрешение": "4MP", "Количество камер": "8", "Хранилище": "4TB HDD", "Аналитика": "Deepin View"},
                        InStock: true, Image: "/static/img/security.svg"},
        }

        for _, c := range components {
                specs, _ := json.Marshal(c.Specs)
                s.db.Exec(`INSERT INTO components (id, name, category, brand, model, price, description, specs, in_stock, image) 
                        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) ON CONFLICT DO NOTHING`,
                        c.ID, c.Name, string(c.Category), c.Brand, c.Model, c.Price, c.Description, string(specs), c.InStock, c.Image)
        }
}

// ─── User methods ─────────────────────────────────────────────────────────────

func (s *Store) getUserByEmail(email string) *User {
        row := s.db.QueryRow(`SELECT id, name, email, password, role, organization, created_at FROM users WHERE email=$1`, email)
        u := &User{}
        err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Role, &u.Organization, &u.CreatedAt)
        if err != nil {
                return nil
        }
        return u
}

func (s *Store) getUserByID(id string) *User {
        row := s.db.QueryRow(`SELECT id, name, email, password, role, organization, created_at FROM users WHERE id=$1`, id)
        u := &User{}
        err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Role, &u.Organization, &u.CreatedAt)
        if err != nil {
                return nil
        }
        return u
}

func (s *Store) getUsers() []*User {
        rows, err := s.db.Query(`SELECT id, name, email, password, role, organization, created_at FROM users ORDER BY created_at`)
        if err != nil {
                return nil
        }
        defer rows.Close()
        var list []*User
        for rows.Next() {
                u := &User{}
                rows.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Role, &u.Organization, &u.CreatedAt)
                list = append(list, u)
        }
        return list
}

func (s *Store) createUser(u *User) error {
        _, err := s.db.Exec(`INSERT INTO users (id, name, email, password, role, organization) VALUES ($1,$2,$3,$4,$5,$6)`,
                u.ID, u.Name, u.Email, u.Password, u.Role, u.Organization)
        return err
}

func (s *Store) updateUser(u *User) error {
        _, err := s.db.Exec(`UPDATE users SET name=$2, email=$3, password=$4, role=$5, organization=$6 WHERE id=$1`,
                u.ID, u.Name, u.Email, u.Password, u.Role, u.Organization)
        return err
}

// ─── Session methods ───────────────────────────────────────────────────────────

func (s *Store) createSession(token, userID string, expiresAt time.Time) {
        s.db.Exec(`INSERT INTO sessions (token, user_id, expires_at) VALUES ($1,$2,$3)`, token, userID, expiresAt)
}

func (s *Store) getSessionUser(token string) *User {
        var userID string
        var expiresAt time.Time
        err := s.db.QueryRow(`SELECT user_id, expires_at FROM sessions WHERE token=$1`, token).Scan(&userID, &expiresAt)
        if err != nil || expiresAt.Before(time.Now()) {
                return nil
        }
        return s.getUserByID(userID)
}

func (s *Store) deleteSession(token string) {
        s.db.Exec(`DELETE FROM sessions WHERE token=$1`, token)
}

func (s *Store) clearAllSessions() {
        s.db.Exec(`DELETE FROM sessions`)
}

// ─── Component methods ────────────────────────────────────────────────────────

func (s *Store) getComponents(category string) []*Component {
        var rows *sql.Rows
        var err error
        if category != "" {
                rows, err = s.db.Query(`SELECT id, name, category, brand, model, price, description, specs, in_stock, image FROM components WHERE category=$1`, category)
        } else {
                rows, err = s.db.Query(`SELECT id, name, category, brand, model, price, description, specs, in_stock, image FROM components`)
        }
        if err != nil {
                return nil
        }
        defer rows.Close()
        var list []*Component
        for rows.Next() {
                c := &Component{}
                var specsRaw string
                rows.Scan(&c.ID, &c.Name, &c.Category, &c.Brand, &c.Model, &c.Price, &c.Description, &specsRaw, &c.InStock, &c.Image)
                json.Unmarshal([]byte(specsRaw), &c.Specs)
                list = append(list, c)
        }
        return list
}

func (s *Store) getComponent(id string) *Component {
        c := &Component{}
        var specsRaw string
        err := s.db.QueryRow(`SELECT id, name, category, brand, model, price, description, specs, in_stock, image FROM components WHERE id=$1`, id).
                Scan(&c.ID, &c.Name, &c.Category, &c.Brand, &c.Model, &c.Price, &c.Description, &specsRaw, &c.InStock, &c.Image)
        if err != nil {
                return nil
        }
        json.Unmarshal([]byte(specsRaw), &c.Specs)
        return c
}

func (s *Store) createComponent(c *Component) error {
        specs, _ := json.Marshal(c.Specs)
        _, err := s.db.Exec(`INSERT INTO components (id, name, category, brand, model, price, description, specs, in_stock, image) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
                c.ID, c.Name, string(c.Category), c.Brand, c.Model, c.Price, c.Description, string(specs), c.InStock, c.Image)
        return err
}

func (s *Store) updateComponent(c *Component) error {
        specs, _ := json.Marshal(c.Specs)
        _, err := s.db.Exec(`UPDATE components SET name=$2, category=$3, brand=$4, model=$5, price=$6, description=$7, specs=$8, in_stock=$9, image=$10 WHERE id=$1`,
                c.ID, c.Name, string(c.Category), c.Brand, c.Model, c.Price, c.Description, string(specs), c.InStock, c.Image)
        return err
}

func (s *Store) deleteComponent(id string) {
        s.db.Exec(`DELETE FROM components WHERE id=$1`, id)
}

// ─── Configuration methods ────────────────────────────────────────────────────

func (s *Store) getConfigurations(userID, role string) []*Configuration {
        var rows *sql.Rows
        var err error
        if role == "admin" {
                rows, err = s.db.Query(`SELECT id, name, user_id, items, total_cost, description, created_at, updated_at FROM configurations ORDER BY created_at DESC`)
        } else {
                rows, err = s.db.Query(`SELECT id, name, user_id, items, total_cost, description, created_at, updated_at FROM configurations WHERE user_id=$1 ORDER BY created_at DESC`, userID)
        }
        if err != nil {
                return nil
        }
        defer rows.Close()
        var list []*Configuration
        for rows.Next() {
                c := &Configuration{}
                var itemsRaw string
                rows.Scan(&c.ID, &c.Name, &c.UserID, &itemsRaw, &c.TotalCost, &c.Description, &c.CreatedAt, &c.UpdatedAt)
                json.Unmarshal([]byte(itemsRaw), &c.Items)
                list = append(list, c)
        }
        return list
}

func (s *Store) getConfiguration(id string) *Configuration {
        c := &Configuration{}
        var itemsRaw string
        err := s.db.QueryRow(`SELECT id, name, user_id, items, total_cost, description, created_at, updated_at FROM configurations WHERE id=$1`, id).
                Scan(&c.ID, &c.Name, &c.UserID, &itemsRaw, &c.TotalCost, &c.Description, &c.CreatedAt, &c.UpdatedAt)
        if err != nil {
                return nil
        }
        json.Unmarshal([]byte(itemsRaw), &c.Items)
        return c
}

func (s *Store) createConfiguration(c *Configuration) error {
        items, _ := json.Marshal(c.Items)
        _, err := s.db.Exec(`INSERT INTO configurations (id, name, user_id, items, total_cost, description) VALUES ($1,$2,$3,$4,$5,$6)`,
                c.ID, c.Name, c.UserID, string(items), c.TotalCost, c.Description)
        return err
}

func (s *Store) updateConfiguration(c *Configuration) error {
        items, _ := json.Marshal(c.Items)
        _, err := s.db.Exec(`UPDATE configurations SET name=$2, items=$3, total_cost=$4, description=$5, updated_at=NOW() WHERE id=$1`,
                c.ID, c.Name, string(items), c.TotalCost, c.Description)
        return err
}

func (s *Store) deleteConfiguration(id string) {
        s.db.Exec(`DELETE FROM configurations WHERE id=$1`, id)
}

func (s *Store) calcCost(items []ConfigItem) float64 {
        total := 0.0
        for _, item := range items {
                var price float64
                s.db.QueryRow(`SELECT price FROM components WHERE id=$1`, item.ComponentID).Scan(&price)
                total += price * float64(item.Quantity)
        }
        return total
}
