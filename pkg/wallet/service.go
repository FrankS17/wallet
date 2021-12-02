package wallet

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/FrankS17/wallet/pkg/types"
	"github.com/google/uuid"
)


func New(text string) error {
	return &errorString{text}
}

type errorString struct {  			// сам тип ошибки не экспортируется
	s string
}

func (e *errorString) Error() string { 		
	return e.s
}

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments 	  []*types.Payment
	favorites	  []*types.Favorite
}

var ErrPhoneRegistered = errors.New("phone already registered")
var ErrAmountMustBePositive = errors.New("amount must be positive")
var ErrAccountNotFound = errors.New("account not found")
var ErrNotEnoughBalance = errors.New("not enough balance")
var ErrPaymentNotFound = errors.New("payment not found")
var ErrFavoriteNotFound = errors.New("favorite not found")


func (s *Service) ExportToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		log.Print(err)
		return err
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()

	data := make([]byte,0)
	lastString := ""
	for _, account := range s.accounts {
		text := []byte(strconv.FormatInt(account.ID,10) + ";" + string(account.Phone) + ";" + strconv.FormatInt(int64(account.Balance),10) + "|")
		data = append(data,text...)
		}
	str := string(data)
	lastString = strings.TrimSuffix(str,"|")
		_, err = file.Write([]byte(lastString))
	if err != nil {
		log.Print(err)
		return err
	}
	//log.Printf("%v",file)
	return nil
}

func (s *Service) ImportFromFile(path string) error {
	
	file, err := os.Open(path)
	if err != nil {
		log.Print(err)
		return err
	}
	defer func() {		
		if cerr := file.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()

	content := make([]byte, 0)
	buf := make([]byte, 4)

	for {
		read, err := file.Read(buf)
		if err == io.EOF { // файл закончился
			content = append(content, buf[:read]...)
			break
		}
		if err != nil {
			log.Print(err)
			return err
		}
		content = append(content, buf[:read]...)
	}

	data := string(content)
	log.Println("data:", data)

	acc := strings.Split(data, "|")
	log.Println("acc:", acc)

//	var account *types.Account
	for _, operation := range acc {
		strAcc := strings.Split(operation, ";")
		log.Println("strAcc:", strAcc)

		id, _ := strconv.ParseInt(strAcc[0], 10, 64)
		phone := types.Phone(strAcc[1])
		balance, _ := strconv.ParseInt(strAcc[2], 10, 64)

		account := types.Account{
			ID:      id,
			Phone:   phone,
			Balance: types.Money(balance),
		}
		 s.accounts = append(s.accounts, &account)	
		 log.Print(account)	
	}
	return nil
}


func (s *Service) Export(dir string) error {
	if s.accounts != nil {
		fileNameAccounts := "/accounts.dump"
		joinName := filepath.Join(dir,fileNameAccounts)
		src, err := os.Create(joinName)
		if err != nil {
			return err
		}
		defer func() {
			if cerr := src.Close(); cerr != nil {
				if err == nil {              	// тут происходит замыканик, значит, можем добраться до имени err проверяем,
					err = cerr		// если nil и при закрытии произошла ошибка, то записываем ошиьку в err
				}
			}
		}()
		data := make([]byte,0)
		lastStr := ""
		for _, acc := range s.accounts {
			text := []byte(strconv.FormatInt(acc.ID,10) + ";" + string(acc.Phone) + ";" +strconv.FormatInt(int64(acc.Balance),10) + string('\n'))
			data = append(data,text...)
			//dataStr := string(data)
		}
		lastStr = string(data)
		aa := strings.TrimSuffix(lastStr,"\n")
		err = os.WriteFile(joinName, []byte(aa),0666) // 0666 - файл, доступен всем на запись и на чтение
		if err != nil {
			log.Print(err)
			return err
		}

		if s.payments != nil {
			fileNamePayments := "/payments.dump"
			joinNamePayments := filepath.Join(dir,fileNamePayments)
			_, err := os.Create(joinNamePayments)
			if err != nil {
				return err
			}
			defer func() {
				if cerr := src.Close(); cerr != nil {
					if err == nil {              	// тут происходит замыканик, значит, можем добраться до имени err проверяем,
						err = cerr		// если nil и при закрытии произошла ошибка, то записываем ошиьку в err
					}
				}
			}()
			data := make([]byte,0)
			lastStr := ""
			for _, payment := range s.payments {
				text := []byte(payment.ID + ";" + strconv.FormatInt(payment.AccountID,10) + ";" + strconv.FormatInt(int64(payment.Amount),10) +
					";" + string(payment.Category) + ";" + string(payment.Status) + string('\n'))
				data = append(data,text...)
				//dataStr := string(data)
			}
			lastStr = string(data)
			aa := strings.TrimSuffix(lastStr,"\n")
			err = os.WriteFile(joinNamePayments, []byte(aa),0666) // 0666 - файл, доступен всем на запись и на чтение
			if err != nil {
				log.Print(err)
				return err
			}

			if s.favorites != nil {
				fileNameFav := "/favorites.dump"
				joinNameFav := filepath.Join(dir,fileNameFav)

				_, err := os.Create(joinNameFav)
				if err != nil {
					return err
				}
				defer func() {
					if cerr := src.Close(); cerr != nil {
						if err == nil {              	// тут происходит замыканик, значит, можем добраться до имени err проверяем,
							err = cerr		// если nil и при закрытии произошла ошибка, то записываем ошиьку в err
						}
					}
				}()
				data := make([]byte,0)
				lastStr := ""
				for _, favorite := range s.favorites {
					text := []byte(favorite.ID + ";" + strconv.FormatInt(favorite.AccountID,10) + ";" + favorite.Name + ";" + strconv.FormatInt(int64(favorite.Amount),10) +
						";" + string(favorite.Category) + string('\n'))
					data = append(data,text...)
					//dataStr := string(data)
				}
				lastStr = string(data)
				aa := strings.TrimSuffix(lastStr,"\n")
				err = os.WriteFile(joinNameFav, []byte(aa),0666) // 0666 - файл, доступен всем на запись и на чтение
				if err != nil {
					log.Print(err)
					return err
				}
			}
		}
	}
return nil
}

func (s *Service) Import(dir string) error {
if s.accounts != nil {
	fileNameAccounts := "/accounts.dump"
	joinName := filepath.Join(dir, fileNameAccounts)
	src, err := os.Open(joinName)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := src.Close(); cerr != nil {
			if err == nil { // тут происходит замыканик, значит, можем добраться до имени err проверяем,
				err = cerr // если nil и при закрытии произошла ошибка, то записываем ошиьку в err
			}
		}
	}()

	reader := bufio.NewReader(src)
	lastLine := make([]string, 0)

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			strTrim := strings.TrimSpace(line)
			lastLine = append(lastLine, strTrim)
			break
		}

		if err != nil {
			return err
		}
		strTrim := strings.TrimSpace(line)
		lastLine = append(lastLine, strTrim)
	}
	//log.Println(lastLine)
	for _, newAcc := range lastLine {
		str := strings.Split(newAcc, ";")

		id, _ := strconv.ParseInt(str[0], 10, 64)
		phone := str[1]
		balance, _ := strconv.ParseInt(str[2], 10, 64)

		if len(str) == 0 {
			break
		}

		fAccount, _ := s.FindAccountByID(id)
		if fAccount != nil {
			fAccount.Phone = types.Phone(phone)
			fAccount.Balance = types.Money(balance)
		} else {
			s.nextAccountID++
			account := &types.Account{
				ID:      id,
				Phone:   types.Phone(phone),
				Balance: types.Money(balance),
			}
			s.accounts = append(s.accounts, account)
			//log.Println("Added Account: ", account)
			//	log.Println("Find Account: ",fAccount)
		}
	}

	if s.payments != nil && len(s.payments) > 0{
		fileNamePayments := "/payments.dump"
		paymentsFile := filepath.Join(dir, fileNamePayments)
		src, err := os.Open(paymentsFile)
		if err != nil {
			return err
		}
		defer func() {
			if cerr := src.Close(); cerr != nil {
				if err == nil { // тут происходит замыканик, значит, можем добраться до имени err проверяем,
					err = cerr // если nil и при закрытии произошла ошибка, то записываем ошиьку в err
				}
			}
		}()

		readerPayment := bufio.NewReader(src)
		paymentsT := make([]string, 0)

		for {
			line, err := readerPayment.ReadString('\n')
			if err == io.EOF {
				strTrim := strings.TrimSpace(line)
				paymentsT = append(paymentsT, strTrim)
				break
			}

			if err != nil {
				return err
			}
			strTrim := strings.TrimSpace(line)
			paymentsT = append(paymentsT, strTrim)
		}
		//log.Println(lastLine)
		for _, newAcc := range paymentsT {
			str := strings.Split(newAcc, ";")
			//{ID: uuid.New().String(),AccountID: 1,Amount: 100_00,Category: "auto",Status: PaymentStatusInProgress},
			id := str[0]
			accountID, _ := strconv.ParseInt(str[1], 10, 64)
			amount, _ := strconv.ParseInt(str[2], 10, 64)
			category := str[3]
			status := str[4]

			log.Println(id, accountID, amount, category, status)

			if len(str) == 0 {
				break
			}

			fPayment, _ := s.FindPaymentByID(id)
			if fPayment != nil {
				fPayment.AccountID = accountID
				fPayment.Amount = types.Money(amount)
				fPayment.Category = types.PaymentCategory(category)
				fPayment.Status = types.PaymentStatus(status)
			} else {
				payment := &types.Payment{
					ID:        id,
					AccountID: accountID,
					Amount:   types. Money(amount),
					Category:  types.PaymentCategory(category),
					Status:    types.PaymentStatus(status),
				}
				s.payments = append(s.payments, payment)
			}
		}
	}


	if s.favorites != nil {
		fileNameFavorites := "/favorites.dump"
		favoritesFile := filepath.Join(dir, fileNameFavorites)
		src, err := os.Open(favoritesFile)
		if err != nil {
			return err
		}
		defer func() {
			if cerr := src.Close(); cerr != nil {
				if err == nil { // тут происходит замыканик, значит, можем добраться до имени err проверяем,
					err = cerr // если nil и при закрытии произошла ошибка, то записываем ошиьку в err
				}
			}
		}()

		readerFav := bufio.NewReader(src)
		favT := make([]string, 0)

		for {
			line, err := readerFav.ReadString('\n')
			if err == io.EOF {
				strTrim := strings.TrimSpace(line)
				favT = append(favT, strTrim)
				break
			}

			if err != nil {
				return err
			}
			strTrim := strings.TrimSpace(line)
			favT = append(favT, strTrim)
		}
		//log.Println(lastLine)
		for _, newAcc := range favT {
			str := strings.Split(newAcc, ";")
			//{ID: uuid.New().String(),AccountID: 1,Name: "megafon",Amount: 100_00,Category: "auto"}
			id := str[0]
			accountID, _ := strconv.ParseInt(str[1], 10, 64)
			name := str[2]
			amount, _ := strconv.ParseInt(str[3], 10, 64)
			category := str[4]

			log.Println(id, accountID, name, amount, category)

			if len(str) == 0 {
				break
			}

			fFav, _ := s.FindFavoriteByID(id)
			if fFav != nil {
				fFav.AccountID = accountID
				fFav.Name = name
				fFav.Amount = types.Money(amount)
				fFav.Category = types.PaymentCategory(category)
			} else {
				favorites := &types.Favorite{
					ID:        id,
					AccountID: accountID,
					Name: name,
					Amount:    types.Money(amount),
					Category:  types.PaymentCategory(category),
				}
				s.favorites = append(s.favorites, favorites)
			}
		}
	}
}
return nil
}



func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error){
	for _, account := range s.accounts {
	if account.Phone == phone {
		return nil, ErrPhoneRegistered
	}
}
	s.nextAccountID++
	account := &types.Account{
		ID: 		s.nextAccountID,
		Phone:		phone,
		Balance: 	0,
	}
	s.accounts = append(s.accounts,account)

	return account, nil
}


func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <=0 {
		return ErrAmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}
	
	if account == nil {
		return ErrAccountNotFound
	}
	
	// зачисление средств пока не рассматриваем как платеж
	account.Balance += amount
	return nil
}


func (s *Service) FindAccountByID(accountID int64) (*types.Account,error) {
	
	for _, account := range s.accounts {
		if account.ID == accountID {
			return account, nil
		}
	}

	return nil, ErrAccountNotFound
}


func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory)(*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		account = acc
		break
	}
	if account == nil {
		return nil, ErrAccountNotFound
	}
	if account.Balance < amount {
		return nil, ErrNotEnoughBalance
	}

	account.Balance -= amount
	paymentID := uuid.New().String()
	payment := &types.Payment{
		ID: paymentID,
		AccountID: accountID,
		Amount: amount,
		Category: category,
		Status: types.PaymentStatusInProgress,
	}
	s.payments = append(s.payments, payment)
	return payment, nil
}


func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	
	for _, payment := range s.payments {
		if payment.ID == paymentID {
			return payment, nil
		}
	}
	return nil, ErrPaymentNotFound
}


func (s *Service) Reject(paymentID string) error {
	
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return ErrPaymentNotFound
	}
	
	account, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		return ErrAccountNotFound
	}

	payment.Status = types.PaymentStatusFail
	account.Balance += payment.Amount

	return nil
}


func (s *Service) Repeat(paymentID string) (*types.Payment, error) {

	/* payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, ErrPaymentNotFound
	}

	account, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		return nil, ErrPaymentNotFound
	}
	
	if account.Balance < payment.Amount {
		return nil, ErrNotEnoughBalance
	}
	account.Balance -= payment.Amount
	aaa := uuid.New().String()
	newP := &types.Payment{
		ID: aaa,
		AccountID: payment.AccountID,
		Amount: payment.Amount,
		Category: payment.Category,
		Status: types.PaymentStatusInProgress,
	}
	s.payments = append(s.payments, newP)
	return newP, nil */

	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	return s.Pay(payment.AccountID, payment.Amount, payment.Category)
}


func (s *Service) FavoritePayment(paymentID, name string) (*types.Favorite, error) {
	
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, ErrPaymentNotFound
	}

	favoritePayment := &types.Favorite{
		ID: uuid.New().String(),
		AccountID: payment.AccountID,
		Name: name,
		Amount: payment.Amount,
		Category: payment.Category,
	}

	s.favorites = append(s.favorites, favoritePayment)	
	return favoritePayment, nil
}


func (s *Service) FindFavoriteByID(favoriteID string) (*types.Favorite, error) {
	
	for _, favorite := range s.favorites {
		if favorite.ID == favoriteID {
			return favorite, nil
		}
	}
	return nil, ErrFavoriteNotFound
}


func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
/*	
	favorite, err := s.FindFavoriteByID(favoriteID)
	if err != nil {
		return nil, ErrFavoriteNotFound 
	}
	
	account, err := s.FindAccountByID(favorite.AccountID)
	if err != nil {
		return nil, ErrAccountNotFound 
	}


	account.Balance -= favorite.Amount
	payment := &types.Payment{
		ID: uuid.New().String(),
		AccountID: favorite.AccountID,
		Amount: favorite.Amount,
		Category: favorite.Category,
		Status: types.PaymentStatusInProgress,
	}	

	s.payments = append(s.payments, payment)	
	return payment, nil
*/
	favorite, err := s.FindFavoriteByID(favoriteID)
	if err != nil {
		return nil, err
	}

	payment, err := s.Pay(favorite.AccountID, favorite.Amount, favorite.Category)
	if err != nil {
		return nil, err
	}

	return payment, nil
}