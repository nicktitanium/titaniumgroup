package model

// Order представляет заявку пользователя.
type Order struct {
	ID            string // первичный ключ в виде строки (шестнадцатеричное представление)
	Username      string
	Coin          string
	PaymentAmount float64 // сумма оплаты в рублях
	WalletAddress string
	PaymentDetail string
}
