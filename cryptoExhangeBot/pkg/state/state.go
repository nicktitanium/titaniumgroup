package state

import "sync"

// Step представляет этап диалога
type Step string

const (
	StepCaptcha      Step = "captcha"
	StepCoinChoice   Step = "coin_choice"  // выбор монеты (будет осуществляться через inline‑кнопки)
	StepAmount       Step = "amount"       // ввод суммы оплаты
	StepWallet       Step = "wallet"       // ввод адреса кошелька
	StepConfirmation Step = "confirmation" // этап подтверждения заявки
)

// UserState хранит данные состояния для пользователя
type UserState struct {
	Step          Step
	CaptchaAnswer int
	Coin          string
	PaymentAmount float64
	Wallet        string
	// Дополнительные поля для редактирования сообщения:
	ChatID        int64
	LastMessageID int
}

var (
	// states хранит состояние для каждого пользователя
	states sync.Map // map[int64]*UserState
)

func GetState(userID int64) *UserState {
	if s, ok := states.Load(userID); ok {
		return s.(*UserState)
	}
	return nil
}

func SetState(userID int64, s *UserState) {
	states.Store(userID, s)
}

func DeleteState(userID int64) {
	states.Delete(userID)
}
