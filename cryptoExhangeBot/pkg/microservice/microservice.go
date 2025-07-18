package microservice

// PaymentResponse соответствует JSON-ответу микросервиса
type PaymentResponse struct {
	PaymentDetail string `json:"payment_detail"`
}

// GetPaymentDetail выполняет GET-запрос к микросервису и возвращает реквизит для оплаты
func GetPaymentDetail(url string) (string, error) {
	return "+79005555555", nil

	//client := http.Client{
	//	Timeout: 10 * time.Second,
	//}
	//resp, err := client.Get(url)
	//if err != nil {
	//	return "", err
	//}
	//defer resp.Body.Close()
	//
	//if resp.StatusCode != http.StatusOK {
	//	return "", errors.New("микросервис вернул неверный статус")
	//}
	//
	//var pr PaymentResponse
	//if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
	//	return "", err
	//}
	//return pr.PaymentDetail, nil
}
