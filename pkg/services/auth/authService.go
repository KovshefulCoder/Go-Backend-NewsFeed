package auth

import (
	"NewsFeed/pkg/models"
	"NewsFeed/pkg/services"
	"github.com/google/uuid"
	"log"
	"strconv"
	"strings"
)

type Storage interface {
	AddNewUser(user models.User) error
	RefreshToken(refreshToken string) (uint32, error)
	AddRefreshToken(refreshToken string, id uint32) error
	UpdateRefreshToken(refreshToken string, id uint32) error
	FinishAuth(userID uint32, nickname string) error
	GetUserByID(userID uint32) (models.User, error)
}

type Service struct {
	storage    Storage
	th         *services.TelethonService
	jwtManager *Manager
}

func NewAuthService(storage Storage, th *services.TelethonService, jwtManager *Manager) *Service {
	return &Service{storage: storage, th: th, jwtManager: jwtManager}
}

// var (
//
//	NameAlreadyExist = errors.New("Name already exist")
//	InvalidTitle     = errors.New("Invalid title")
//	NotUserArea      = errors.New("Not user area")
//
// )

func newUserID() uint32 {
	return uuid.New().ID()
}

func (s *Service) SignUpNumber(number models.Number) (models.Session, error) {
	userID := newUserID()
	signUpNumberBody := models.AuthNumberRequest{Id: userID, Command: "send_code", Number: number}
	err := s.th.SignNumber(signUpNumberBody)
	if err != nil {
		log.Printf("SignUpNumber: can`t signUp with telethon: %v", err)
		return models.Session{}, err
	}
	user := models.NewUser(number, userID)
	err = s.storage.AddNewUser(user) //поменять местами
	if err != nil {
		log.Printf("SignUpNumber: can`t add new user to db: %v", err)
		return models.Session{}, err
	}
	token, err := s.jwtManager.CreateToken(user.ID)
	if err != nil {
		log.Printf("SignUpNumber: can`t create token: %v", err)
		return models.Session{}, err
	}
	refreshToken := s.jwtManager.CreateRefreshToken()
	err = s.storage.AddRefreshToken(refreshToken, user.ID)
	if err != nil {
		log.Printf("SignUpNumber: can`t add refresh token: %v", err)
		return models.Session{}, err
	}
	return models.Session{Token: token, Refresh: refreshToken}, nil
}

func (s *Service) RefreshToken(refreshToken string) (models.Session, error) {
	userID, err := s.storage.RefreshToken(refreshToken)
	if err != nil {
		log.Printf("RefreshToken: can`t check refresh token from db: %v", err)
		return models.Session{}, err
	}
	log.Println("RefreshToken: ID: ", strconv.Itoa(int(userID)))
	token, err := s.jwtManager.CreateToken(userID)
	if err != nil {
		log.Printf("RefreshToken: can`t create token: %v", err)
		return models.Session{}, err
	}
	refreshToken = s.jwtManager.CreateRefreshToken()
	err = s.storage.UpdateRefreshToken(refreshToken, userID)
	if err != nil {
		log.Printf("RefreshToken: can`t add refresh token: %v", err)
		return models.Session{}, err
	}
	return models.Session{Token: token, Refresh: refreshToken}, nil
}

func (s *Service) RetrieveIDFromToken(token string) (uint32, error) {
	userID, err := s.jwtManager.GetIDFromToken(token)
	if err != nil {
		log.Printf("RetrieveIDFromToken: can`t retrieve token: %v", err)
		return 0, err
	}
	return userID, nil
}

func (s *Service) AuthCode(userID uint32, code string) (models.AuthFinishResponse, error) {
	user, err := s.storage.GetUserByID(userID)
	if err != nil {
		log.Printf("AuthCode: can`t get user from db: %v", err)
		return models.AuthFinishResponse{}, err
	}
	signUpCodeBody := models.AuthCodeRequest{UserID: userID, Code: code, Number: user.Number, Command: "sign_in"}
	signUpCodeResult, err := s.th.AuthCode(signUpCodeBody)
	if err != nil {
		return models.AuthFinishResponse{}, err
	}
	err = s.storage.FinishAuth(userID, signUpCodeResult.Details.Nickname)
	if err != nil {
		log.Printf("AuthCode: can`t finish signUp of user: %v", err)
		return models.AuthFinishResponse{}, err
	}
	parts := strings.Split(signUpCodeResult.Details.Nickname, ",")
	return models.AuthFinishResponse{Nickname: parts[0]}, err
}

func (s *Service) Auth2FA(userID uint32, password string) (models.AuthTelethonFinishResponse, error) {
	signUp2FABody := models.Auth2FARequest{UserID: userID, Password: password}
	signUp2FAResponse, err := s.th.Auth2FA(signUp2FABody)
	if err != nil {
		log.Printf("Auth2FA: can`t signUp: %v", err.Error())
		return models.AuthTelethonFinishResponse{}, err
	}
	err = s.storage.FinishAuth(userID, signUp2FAResponse.Details.Nickname)
	if err != nil {
		log.Printf("Auth2FA: can`t finish signUp of user: %v", err)
		return models.AuthTelethonFinishResponse{}, err
	}
	return signUp2FAResponse, err
}

func (s *Service) SignInNumber(number models.Number, userID uint32) error {
	signUpNumberBody := models.AuthNumberRequest{Id: userID, Command: "sign_in", Number: number}
	err := s.th.SignNumber(signUpNumberBody)
	if err != nil {
		log.Printf("SignInNumber: can`t signIn with telethon: %v", err)
		return err
	}
	return nil
}

func (s *Service) CheckClient(userID uint32) error {
	checkClientBody := models.AuthCheckClientRequest{UserID: userID, Command: "has_session"}
	err := s.th.CheckClient(checkClientBody)
	if err != nil {
		log.Printf("CheckClient: can`t check client in pygram server: %v", err)
		return err
	}
	return nil
}
