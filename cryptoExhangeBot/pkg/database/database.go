package database

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"cryptoChange/pkg/model"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// InitDB открывает (или создаёт) базу данных и инициализирует таблицу orders.
// В таблице orders поле id теперь имеет тип TEXT и является PRIMARY KEY.
func InitDB(dataSourceName string) error {
	var err error
	DB, err = sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return err
	}

	// Создаём таблицу orders, если её нет.
	query := `
	CREATE TABLE IF NOT EXISTS orders (
		id TEXT PRIMARY KEY,
		username TEXT,
		coin TEXT,
		payment_amount REAL,
		wallet_address TEXT,
		payment_detail TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	if _, err = DB.Exec(query); err != nil {
		return err
	}

	log.Println("База данных инициализирована")
	return nil
}

// getNextOrderID вычисляет следующий id заявки в шестнадцатеричной системе.
// Он ищет максимальный id в таблице, преобразует его из шестнадцатеричной строки в число,
// увеличивает на 1 и возвращает в виде строки в верхнем регистре.
// Если заявок ещё нет, возвращается "1".
func getNextOrderID() (string, error) {
	var lastID string
	err := DB.QueryRow("SELECT id FROM orders ORDER BY created_at DESC LIMIT 1").Scan(&lastID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Нет ни одной заявки — начинаем с 1.
			return "1", nil
		}
		return "", err
	}

	// Преобразуем lastID из шестнадцатеричного представления в число.
	lastNum, err := strconv.ParseInt(lastID, 16, 64)
	if err != nil {
		return "", fmt.Errorf("ошибка преобразования lastID: %v", err)
	}

	nextNum := lastNum + 1
	// Преобразуем число в шестнадцатеричное представление (например, 10 -> "A").
	nextID := fmt.Sprintf("%X", nextNum)
	return nextID, nil
}

// InsertOrder вставляет новую заявку в таблицу orders. При этом функция генерирует следующий id заявки.
func InsertOrder(order *model.Order) error {
	// Генерируем следующий id.
	nextID, err := getNextOrderID()
	if err != nil {
		return err
	}
	order.ID = nextID

	query := `INSERT INTO orders (id, username, coin, payment_amount, wallet_address, payment_detail) VALUES (?, ?, ?, ?, ?, ?)`
	stmt, err := DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(order.ID, order.Username, order.Coin, order.PaymentAmount, order.WalletAddress, order.PaymentDetail)
	return err
}
