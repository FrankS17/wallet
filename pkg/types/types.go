package types


// Представляет собой денежную сумму в минимальных единицах (копейки и т.д.)
type Money int64


// Category представляет собой категорию, в которой был совершен платеж
type PaymentCategory string

type PaymentStatus string

const (
	PaymentStatusOk PaymentStatus = "Ok"
	PaymentStatusFail PaymentStatus = "FAIL"
	PaymentStatusInProgress PaymentStatus = "INPROGRESS"
)

type Payment struct {
	ID string 
	AccountID int64
	Amount Money
	Category PaymentCategory
	Status PaymentStatus
}

type Phone string

type Account struct {
	ID int64
	Phone Phone
	Balance Money
}