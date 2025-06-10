package services

import (
	"errors"
	"financeAppAPI/internal/models"
	"financeAppAPI/internal/repositories"
	"financeAppAPI/internal/utils"
	"math/rand"
	"strconv"
	"time"
)

type CardService struct {
	cardRepo    *repositories.CardRepository
	accountRepo *repositories.AccountRepository
}

func NewCardService(cardRepo *repositories.CardRepository) *CardService {
	return &CardService{cardRepo: cardRepo}
}

func (s *CardService) SetAccountRepo(accountRepo *repositories.AccountRepository) {
	s.accountRepo = accountRepo
}

func (s *CardService) ValidateAccountOwnership(userID, accountID int) bool {
	if s.accountRepo == nil {
		return false
	}
	account, err := s.accountRepo.GetAccountByID(accountID)
	if err != nil {
		return false
	}
	return account.UserID == userID
}

func generateLuhnCardNumber(prefix string, length int) string {
	num := prefix
	for len(num) < length-1 {
		num += strconv.Itoa(rand.Intn(10))
	}
	// Calculate Luhn check digit
	sum := 0
	alt := false
	for i := len(num) - 1; i >= 0; i-- {
		n, _ := strconv.Atoi(string(num[i]))
		if alt {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		alt = !alt
	}
	checkDigit := (10 - (sum % 10)) % 10
	return num + strconv.Itoa(checkDigit)
}

func (s *CardService) CreateCard(accountID int) (*models.Card, error) {
	cardNumber := generateLuhnCardNumber("4000", 16)
	expiryDate := time.Now().AddDate(3, 0, 0).Format("01/06")
	cvv := "123"
	encryptedCardNumber, err := utils.EncryptPGP(cardNumber)
	if err != nil {
		return nil, err
	}
	encryptedExpiryDate, err := utils.EncryptPGP(expiryDate)
	if err != nil {
		return nil, err
	}
	cvvHash, err := utils.HashCVV(cvv)
	if err != nil {
		return nil, err
	}
	hmacValue, err := utils.GenerateHMAC(cardNumber + expiryDate)
	if err != nil {
		return nil, err
	}
	card := &models.Card{
		AccountID:           accountID,
		CardNumber:          cardNumber,
		ExpiryDate:          expiryDate,
		CardNumberEncrypted: encryptedCardNumber,
		ExpiryDateEncrypted: encryptedExpiryDate,
		CVVHash:             cvvHash,
		HMAC:                hmacValue,
		CreatedAt:           time.Now(),
	}
	if err := s.cardRepo.CreateCard(card); err != nil {
		return nil, err
	}
	card.CVVHash = ""
	card.CardNumberEncrypted = ""
	card.ExpiryDateEncrypted = ""
	card.HMAC = ""
	return card, nil
}

func (s *CardService) GetCards(accountID int) ([]*models.Card, error) {
	cards, err := s.cardRepo.GetCardsByAccountID(accountID)
	if err != nil {
		return nil, err
	}
	for _, card := range cards {
		cardNumber, err := utils.DecryptPGP(card.CardNumberEncrypted)
		if err == nil {
			card.CardNumber = cardNumber
		}
		expiry, err := utils.DecryptPGP(card.ExpiryDateEncrypted)
		if err == nil {
			card.ExpiryDate = expiry
		}
		if !utils.VerifyHMAC(card.CardNumber+card.ExpiryDate, card.HMAC) {
			return nil, errors.New("HMAC карты не совпадает")
		}
		card.CVVHash = ""
		card.CardNumberEncrypted = ""
		card.ExpiryDateEncrypted = ""
		card.HMAC = ""
	}
	return cards, nil
}
