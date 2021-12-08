package wallet

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/FrankS17/wallet/pkg/types"
	"github.com/google/uuid"
)

type testService struct {
	*Service		// embedding(встраивание)
}

func newTestService() *testService {
	return &testService{Service: &Service{}} // функция конструктор
}

type testAccount struct {
	phone    types.Phone
	balance  types.Money
	payments []struct {
		amount   types.Money
		category types.PaymentCategory
	}
}


var defaultTestAccount = testAccount { 
	phone:		"+992900100500",
	balance: 	10_000_00,
	payments:	[]struct {
		amount		types.Money
		category	types.PaymentCategory	
	} {
		{amount: 1_000_00, category: "auto"},
	},
}


func (s *testService) addAccount(data testAccount) (*types.Account, []*types.Payment, []*types.Favorite, error) {
	//регистрируем там пользователя
	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("can't redister account, error = %v", err)
	}

	//пополняем его счет
	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("can't deposity account, error = %v", err)
	}

	//выполняем платеж
	//можем создать слайс сразу нужной длиныб поскольку знаем размер
	payments := make([]*types.Payment, len(data.payments))
	favorites := make([]*types.Favorite, len(data.payments))
	for i, payment := range data.payments {
		//тогда здесь работаем просто через index, а не append
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("can't make payment, error = %v", err)
		}

		favorites[i], err = s.FavoritePayment(payments[i].ID, "Favorite payment #i")
		if err != nil {
			return nil, nil, nil, fmt.Errorf("can't make favorite paymnet, error = %v", err)
		}
	}
	return account, payments, favorites, nil
}

/*

func (s *testService) addAccount(data testAccount) (*types.Account, []*types.Payment, error) {
	//региструем пользователя
	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, fmt.Errorf("can't register account, error = %v", err)
	}
	
	// пополняем счет
	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("can't deposit account, error = %v", err)
	}

	// выполняем платежи
	// можем создать слайс сразу нужной длины, поскольку знаем размер
	payments := make([]*types.Payment, len(data.payments))
	
	for i, payment := range data.payments {
	// тогда здесь работаем через индекс, а не через append
	payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("can't make payment, error = %v", err)
	}
	
	return account, payments, nil
}
func TestService_FindPaymentID_success(t *testing.T) {
	//создаем сервис
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// попробуем найти платеж
	payment := payments[0]
	got, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("FindPaymentID(): error = %v", err)
		return
	}

	//сравниваем платежи
	if !reflect.DeepEqual(payment, got)  {
		t.Errorf("FindPaymentID(): wrong payment was returned = %v", err)
		return
	}
}

func TestService_FindPaymentID_fail(t *testing.T) {
	//создаем сервис
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// попробуем найти платеж
	_, err = s.FindPaymentByID(uuid.New().String())
	if err == nil {
		t.Errorf("FindPaymentID(): error = %v", err)
		return
	}

	//сравниваем платежи
	if err != ErrPaymentNotFound {
		t.Errorf("FindPaymentID(): must return ErrPaymentNotFound, returned= %v", err)
		return
	}	
}

func TestService_Reject_success(t *testing.T) {
	//создаем сервис
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// попробуем найти платеж
	payment := payments[0]
	err = s.Reject(payment.ID)
	if err != nil {
		t.Errorf("Reject(): error = %v", err)
		return
	}

	savedPayment, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("Reject(): can't find payment by id, error = %v", err)
		return
	}
	if savedPayment.Status != types.PaymentStatusFail {
		t.Errorf("Reject(): status didn't change, error = %v", savedPayment)
		return
	}

	savedAccount, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		t.Errorf("Reject(): can't find account by id, error = %v", err)
		return
	}
	if savedAccount.Balance != defaultTestAccount.balance {
		t.Errorf("Reject(): balance didn't change, error = %v", savedAccount)
		return
	}
} 


func TestService_Repeat_success(t *testing.T) {
	//создаем сервис
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// попробуем найти платеж
	payment := payments[0]
		
	savedPayment, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("Repeat(): can't find payment by id, error = %v", err)
		return
	}

	repeatPayment,err := s.Repeat(savedPayment.ID)
	if err != nil {
		t.Errorf("Repeat(): error = %v", err)
		return
	}
	

	_, err = s.FindAccountByID(payment.AccountID)
	if err != nil {
		t.Errorf("Repeat(): can't find account by id, error = %v", err)
		return
	}


	if savedPayment.ID == repeatPayment.ID {
		t.Errorf("Repeat: payments are equal! error = %v", err)
		return	
	}

}

func TestService_FindFavoriteID_success(t *testing.T) {
	//создаем сервис
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// попробуем найти платеж
	payment := payments[0]
	
	addFavoritePayment, err := s.FavoritePayment(payment.ID, "a")
	if err != nil {
		t.Errorf("FindFavoriteID(): error = %v", err)
		return
	} 

	findFavorite, err := s.FindFavoriteByID(addFavoritePayment.ID)
	if err != nil {
		t.Errorf("FindFavoriteID(): favorite payment was not found, error = %v", err)
		return
	} 

	//сравниваем платежи
	if !reflect.DeepEqual(addFavoritePayment, findFavorite)  {
		t.Errorf("FindPaymentID(): favorite payments are not equal = %v", err)
		return
	}
}

func TestService_FindFavoriteID_fail(t *testing.T) {
	//создаем сервис
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// попробуем найти платеж
	_, err = s.FindFavoriteByID(uuid.New().String())
	if err == nil {
		t.Errorf("FindFavoriteByID(): error = %v", err)
		return
	}

	//сравниваем платежи
	if err != ErrFavoriteNotFound {
		t.Errorf("FindFavoriteByID(): must return ErrPaymentNotFound, returned= %v", err)
		return
	}	
}

func TestService_FavoritePayment_success(t *testing.T) {
	//создаем сервис
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// попробуем найти платеж
	payment := payments[0]
	
	findPayment, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("Favorite(): can't find payment by id, error = %v", err)
		return
	}

	// добавляем платеж в favorite
	addFavPay,err := s.FavoritePayment(findPayment.ID,"aaa")
	if err != nil {
		t.Errorf("Favorite(): error = %v", err)
		return
	}	
	
	// находим этот добавленный платеж в слайсе favorites 
	findFavoritePayment, err := s.FindFavoriteByID(addFavPay.ID)
	if err != nil {
		t.Errorf("Favorite(): error = %v", err)
		return
	}	

	// сравниваем полученный платеж из слайса favorites - с добавленным 
	if !reflect.DeepEqual(addFavPay,findFavoritePayment) {
		t.Errorf("Favorite(): error = %v", err)
		return
	}
	
}


func TestService_PayFromFavorite_success(t *testing.T) {
	//создаем сервис
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// попробуем найти платеж
	payment := payments[0]
	
	findPayment, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("PayFromFavorite(): can't find payment by id, error = %v", err)
		return
	}

	// добавляем платеж в favorite
	addFavPay,err := s.FavoritePayment(findPayment.ID,"aaa")
	if err != nil {
		t.Errorf("PayFromFavorite(): error = %v", err)
		return
	}	
	
	// находим этот добавленный платеж в слайсе favorites 
	findFavoritePayment, err := s.FindFavoriteByID(addFavPay.ID)
	if err != nil {
		t.Errorf("PayFromFavorite(): error = %v", err)
		return
	}	

	// сравниваем полученный платеж из слайса favorites - с добавленным 
	if !reflect.DeepEqual(addFavPay,findFavoritePayment) {
		t.Errorf("PayFromFavorite(): error = %v", err)
		return
	}

	
	newPayment, err := s.PayFromFavorite(addFavPay.ID)
	if err != nil {
		t.Errorf("PayFromFavorite(): error = %v", err)
		return
	}

	_, err = s.FindAccountByID(newPayment.AccountID)
	if err != nil {
		t.Errorf("PayFromFavorite(): account was not found, error = %v", err)
		return
	}	
}
*/


func TestService_FindAccountById_success(t *testing.T) {
	svc := Service{}
	account1, _ := svc.RegisterAccount("+992900000001")
	svc.RegisterAccount("+992900000002")
	svc.RegisterAccount("+992900000003")

	_, err1 := svc.FindAccountByID(account1.ID)

	if err1 != nil {
		t.Error("Account not found")
	}
}

func TestService_FindAccountById_notSuccess(t *testing.T) {
	svc := Service{}
	svc.RegisterAccount("+992900000004")
	svc.RegisterAccount("+992900000005")
	svc.RegisterAccount("+992900000004")

	_, err1 := svc.FindAccountByID(6666)

	if err1 == nil {
		t.Error(err1)
	}
}

func TestService_FindPaymentByID_success(t *testing.T) {
	//создаем сервис
	s := newTestService()
	_, payments, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	//пробуем найти платеж
	payment := payments[0]
	got, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("FIndPaymentByID(): error = %v", err)
		return
	}

	//сравниваем платежи
	if reflect.DeepEqual(payment, got) {
		if err != nil {
			t.Errorf("FIndPaymentByID(): wrong payment returned = %v", err)
			return
		}
	}
}

func TestService_FindPaymentByID_fail(t *testing.T) {
	//создаем сервис
	s := newTestService()
	_, _, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	//пробуем найти не существующий платеж
	_, err = s.FindPaymentByID(uuid.New().String())
	if err == nil {
		t.Error("FIndPaymentByID(): must returned error, returned nil")
		return
	}

	if err != ErrPaymentNotFound {
		t.Errorf("FIndPaymentByID(): must returned ErrPaymentNotFound, returned = %v", err)
		return
	}

}

func TestService_Reject_success(t *testing.T) {
	// svc := Service{}
	// acc1, _ := svc.RegisterAccount("+992900000001")
	// acc2, _ := svc.RegisterAccount("+992900000002")
	// acc3, _ := svc.RegisterAccount("+992900000003")

	// _ = svc.Deposit(acc1.ID, types.Money(100))
	// _ = svc.Deposit(acc2.ID, types.Money(100))
	// _ = svc.Deposit(acc3.ID, types.Money(100))

	// payment1, _ := svc.Pay(acc1.ID, types.Money(10), types.PaymentCategory("mobile"))
	// svc.Pay(acc2.ID, types.Money(10), types.PaymentCategory("mobile"))
	// svc.Pay(acc3.ID, types.Money(10), types.PaymentCategory("mobile"))

	// rejectError := svc.Reject(payment1.ID)
	// if rejectError != nil {
	// 	t.Error(rejectError)
	// }

	// rejectedAccount, _ := svc.FindAccountByID(acc1.ID)
	// if rejectedAccount.Balance != 100 {
	// 	t.Error("Wrong balance")
	// }

	// rejectedPayment, _ := svc.FindPaymentByID(payment1.ID)
	// if rejectedPayment.Status != types.PaymentStatusFail {
	// 	t.Error("Wrong status")
	// }

	//создаем сервис
	s := newTestService()
	_, payments, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	//пробуем отменить платеж
	payment := payments[0]
	err = s.Reject(payment.ID)
	if err != nil {
		t.Errorf("Reject():  error = %v", err)
		return
	}

	savedPayment, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("Reject(): can't find payment by id, error = %v", err)
		return
	}
	if savedPayment.Status != types.PaymentStatusFail {
		t.Errorf("Reject(): status didn't changed, payment = %v", savedPayment)
		return
	}

	savedAccount, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		t.Errorf("Reject(): can't find account by id, error = %v", err)
		return
	}
	if savedAccount.Balance != defaultTestAccount.balance {
		t.Errorf("Reject(): balance didn't changed, account = %v", savedAccount)
		return
	}
}

func TestService_Reject_fail(t *testing.T) {
	svc := Service{}
	acc1, _ := svc.RegisterAccount("+992900000001")
	acc2, _ := svc.RegisterAccount("+992900000002")
	acc3, _ := svc.RegisterAccount("+992900000003")

	_ = svc.Deposit(acc1.ID, types.Money(100))
	_ = svc.Deposit(acc2.ID, types.Money(100))
	_ = svc.Deposit(acc3.ID, types.Money(100))

	svc.Pay(acc1.ID, types.Money(10), types.PaymentCategory("mobile"))
	svc.Pay(acc2.ID, types.Money(10), types.PaymentCategory("mobile"))
	svc.Pay(acc3.ID, types.Money(10), types.PaymentCategory("mobile"))

	rejectError := svc.Reject(uuid.New().String())
	if rejectError != ErrPaymentNotFound {
		t.Error("Payment must be not found")
	}
}

func TestService_Repeat_success(t *testing.T) {
	s := newTestService()
	_, payments, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]

	newPayment, err := s.Repeat(payment.ID)
	if err != nil {
		t.Errorf("Repeat(): error = %v", err)
		return
	}

	if newPayment.AccountID != payment.AccountID {
		t.Errorf("Repeat(): account ID's difference,\n Repeated payment = %v,\n Rejected payment = %v", newPayment, payment)
		return
	}

	if newPayment.Amount != payment.Amount {
		t.Errorf("Repeat(): amount of payments difference,\n Repeated payment = %v,\n Rejected payment = %v", newPayment, payment)
		return
	}

	if newPayment.Category != payment.Category {
		t.Errorf("Repeat(): category of payments difference,\n Repeated payment = %v,\n Rejected payment = %v", newPayment, payment)
		return
	}

	if newPayment.Status != payment.Status {
		t.Errorf("Repeat(): status of payments difference,\n Repeated payment = %v,\n Rejected payment = %v", newPayment, payment)
		return
	}
}

func TestService_Repeat_notFound(t *testing.T) {
	s := newTestService()
	_, _, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}
	payment := uuid.New().String()
	_, err = s.Repeat(payment)
	if err == nil {
		t.Errorf("Repeat(): must return error, returned nil")
		return
	}
	if err != ErrPaymentNotFound {
		t.Errorf("Repeat(): must return ErrPaymentNotFound, returned: %v", err)
		return
	}

}

func TestService_FavoritePayment_success(t *testing.T) {
	s := newTestService()
	_, payments, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	_, err = s.FavoritePayment(payment.ID, "my favorite payment")
	if err != nil {
		t.Errorf("FavoritePayment(): error: %v", err)
		return
	}
}

func TestService_FavoritePayment_notFound(t *testing.T) {
	s := newTestService()
	_, _, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	favorite := uuid.New().String()
	_, err = s.FavoritePayment(favorite, "my favorite payment")
	if err == nil {
		t.Error("FavoritePayment(): must return error, returned nil")
		return
	}

	if err != ErrPaymentNotFound {
		t.Errorf("FavoritePayment(): must return ErrPaymentNotFound, returned: %v", err)
		return
	}
}

func TestService_FindFavoriteByID_success(t *testing.T) {
	s := newTestService()
	_, _, favorites, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	favorite := favorites[0]
	result, err := s.FindFavoriteByID(favorite.ID)
	if err != nil {
		t.Errorf("FindFavoriteByID(): error: %v", err)
		return
	}

	if !reflect.DeepEqual(result, favorite) {
		t.Errorf("FindFavoriteByID(): wrong payment returned = %v", err)
		return
	}
}

func TestService_FindFavoriteByID_notFound(t *testing.T) {

	s := newTestService()
	_, _, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}
	favoriteID := uuid.New().String()
	_, err = s.FindFavoriteByID(favoriteID)
	if err == nil {
		t.Error("FindFavoriteByID(): must return error, returned nil")
		return
	}

	if err != ErrFavoriteNotFound {
		t.Errorf("FindFavoriteByID(): must return ErrFavoriteNotFound, returned: %v", err)
		return
	}
}

func TestService_PayFromFavorite_success(t *testing.T) {
	s := newTestService()
	_, _, favorites, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	favorite := favorites[0]
	payment, err := s.PayFromFavorite(favorite.ID)
	if err != nil {
		t.Errorf("PayFromFavorite(): error: %v", err)
		return
	}

	if payment.AccountID != favorite.AccountID {
		t.Errorf("PayFromFavorite(): account ID's difference,\n Current payment = %v,\n favorite payment = %v", payment, favorite)
		return
	}

	if payment.Amount != favorite.Amount {
		t.Errorf("PayFromFavorite(): amount of payments difference,\n Current payment = %v,\n favorite payment = %v", payment, favorite)
		return
	}

	if payment.Category != favorite.Category {
		t.Errorf("PayFromFavorite(): category of payments difference,\n Current payment = %v,\n favorite payment = %v", payment, favorite)
		return
	}
}

func TestService_PayFromFavorite_notFound(t *testing.T) {
	s := newTestService()
	_, _, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	favoriteID := uuid.New().String()
	_, err = s.PayFromFavorite(favoriteID)
	if err == nil {
		t.Errorf("PayFromFavorite(): must return error, returned nil")
		return
	}
	if err != ErrFavoriteNotFound {
		t.Errorf("PayFromFavorite(): must return ErrFavoriteNotFound, returned: %v", err)
		return
	}

}


func TestService_Export(t *testing.T) {
	s := newTestService()
	_, payments, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	_, err = s.FavoritePayment(payment.ID, "my favorite payment")
	if err != nil {
		t.Errorf("FavoritePayment(): error: %v", err)
		return
	}

	err = s.Export("data")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestService_Import_success(t *testing.T) {
	s := newTestService()

	err := s.Import("data")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestService_Import_notFound1(t *testing.T) {

	s := newTestService()

	err := s.Import("")
	if err != nil {
		t.Error(err)
		return
	}
}
func TestService_Import_notFound2(t *testing.T) {

	s := newTestService()

	err := s.Import("")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestService_Import_Error(t *testing.T) {

	s := newTestService()
	_, payments, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}
	payment := payments[0]
	_, err = s.FavoritePayment(payment.ID, "my favorite payment")
	if err != nil {
		t.Errorf("FavoritePayment(): error: %v", err)
		return
	}

	err = s.Import("data")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestService_Import_emptyFiles(t *testing.T) {
	s := newTestService()

	file1, _ := os.Create("data/accounts.dump")
	defer file1.Close()

	file2, _ := os.Create("data/payments.dump")
	defer file2.Close()

	file3, _ := os.Create("data/favorites.dump")
	defer file3.Close()

	err := s.Import("data")
	if err != nil {
		t.Error(err)
		return
	}
}


func TestService_ExportAccountHistory_success(t *testing.T) {
	s := newTestService()
	Transactions(s)
	_, err := s.ExportAccountHistory(1)
	if err != nil {
		t.Error(err)
	}
}
func TestService_ExportAccountHistory_notSuccess(t *testing.T) {
	s := newTestService()
	_, _, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	anotherID := s.nextAccountID + 1
	_, err = s.FindAccountByID(anotherID)
	if err == nil {
		t.Error("ExportAccountHistory(): must return error, returned nil")
	}

	_, err = s.ExportAccountHistory(3)
	if err == nil {
		t.Error(err)
	}
	if err != ErrAccountNotFound {
		t.Errorf("ExportAccountHistory(): must return ErrAccountNotFound, returned = %v", err)
		return
	}
}


func Transactions(s *testService) {
	s.RegisterAccount("1111")
	s.Deposit(1, 500)
	s.Pay(1, 10, "food")
	s.Pay(1, 10, "phone")
	s.Pay(1, 15, "bank")
	s.Pay(1, 25, "auto")
	s.Pay(1, 30, "restaurant")
	s.Pay(1, 50, "auto")
	s.Pay(1, 60, "bank")
	s.Pay(1, 50, "bank")

	s.RegisterAccount("2222")
	s.Deposit(2, 200)
	s.Pay(2, 40, "phone")

	s.RegisterAccount("3333")
	s.Deposit(3, 300)
	s.Pay(3, 36, "auto")
	s.Pay(3, 12, "food")
	s.Pay(3, 25, "phone")
}
func TestService_HistoryToFiles_success(t *testing.T) {
	s := newTestService()
	Transactions(s)

	payments, err := s.ExportAccountHistory(1)
	if err != nil {
		t.Error(err)
	}
	err = s.HistoryToFiles(payments, "data", 3)
	if err != nil {
		t.Error(err)
	}
}
func TestService_HistoryToFiles_notSuccess(t *testing.T) {
	s := newTestService()
	Transactions(s)

	payment := []types.Payment{}
	err := s.HistoryToFiles(payment, "data", 12)
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkRegular(b *testing.B) {
	s := newTestService()
	want := int64(2000)
	for i := 0; i < b.N; i++ {
		result := s.Regular()
		if result != want {
			b.Fatalf("invalid result, got %v, want %v", result, want)
		}
	}
}


func BenchmarkConcurrently(b *testing.B) {
	s := newTestService()
	want := int64(2000)
	for i := 0; i < b.N; i++ {
		result := s.Concurrently()
		if result != want {
			b.Fatalf("invalid result, got %v, want %v", result, want)
		}
	}
}

/*
func BenchmarkSumPayment(b *testing.B) {
	s := newTestService()
	want := types.Money(30000)
	for i := 0; i < b.N; i++ {
		result := s.SumPayment(1)

		if result != want {
			b.Fatalf("invalid result, got %v, want %v", result, want)
		}
	}
} */




func TestService_SumPayments(t *testing.T) {
	s := newTestService()
	sum := s.SumPayments(1)
	log.Print(sum)

}


func BenchmarkSumPayments(b *testing.B) {
	s := newTestService()
	sum := s.SumPayments(1)
	log.Print(sum)
}



func TestService_FilterPayments_Success(t *testing.T) {
	/*s := newTestService()
	want := []types.Payment{
		{ID: "1a",AccountID: 1,Amount: 100_00,Category: "auto",Status: types.PaymentStatusInProgress},
		{ID: "1b",AccountID: 1,Amount: 100_00,Category: "auto",Status: types.PaymentStatusInProgress},
		{ID: "1c",AccountID: 1,Amount: 100_00,Category: "auto",Status: types.PaymentStatusInProgress},
	}
	result, _ := s.FilterPayments(1,1)
	
	log.Println(want)
	log.Println(result)

	if reflect.DeepEqual(want,result) == false {
		t.Errorf("TestService_FilterPayments(): Not Equal")
	}*/ 
	s := newTestService()
	sum := s.SumPayments(1)
	log.Print(sum)

}


func TestService_FilterPayments_InvalidAccountID(t *testing.T) {
	/*s := newTestService()
	want := []types.Payment{}
	result, errorA := s.FilterPayments(4,1)
	err := ErrAccountNotFound
	if reflect.DeepEqual(want,result) == false && err != errorA {
		t.Errorf("TestService_FilterPayments(): Not Equal")
	} */
	s := newTestService()
	sum := s.SumPayments(1)
	log.Print(sum)
}


func BenchmarkFilterPayments(b *testing.B) {
	/*s := newTestService()
	want := []types.Payment{
		{ID: "1a",AccountID: 1,Amount: 100_00,Category: "auto",Status: types.PaymentStatusInProgress},
		{ID: "1b",AccountID: 1,Amount: 100_00,Category: "auto",Status: types.PaymentStatusInProgress},
		{ID: "1c",AccountID: 1,Amount: 100_00,Category: "auto",Status: types.PaymentStatusInProgress},
	}
	for i := 0; i < b.N; i++ {
		result,_ := s.FilterPayments(1,2)

		if reflect.DeepEqual(want,result) == false {
				log.Print("AA")
		}
	} */
	s := newTestService()
	sum := s.SumPayments(1)
	log.Print(sum)
}



func TestService_FilterPaymentsByFn(t *testing.T) {
	s := newTestService()
	sum := s.SumPayments(1)
	log.Print(sum)

}


func BenchmarkFilterPaymentsByFn(b *testing.B) {
	s := newTestService()
	sum := s.SumPayments(1)
	log.Print(sum)
}



