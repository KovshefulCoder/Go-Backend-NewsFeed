package services

import (
	"NewsFeed/pkg/models"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type TelethonService struct {
	client *http.Client
}

func NewTelethonService(client *http.Client) *TelethonService {

	return &TelethonService{client: client}
}

// Yes, this it wicked and wrong, wasn`t my accountability
const sharedUrl = ""
const sharedUrlMethod = "GET"

func (ts *TelethonService) SignNumber(authSignUpNumberRequest models.AuthNumberRequest) error {
	//Convert the request body to a JSON byte array
	jsonBytes, err := authSignUpNumberRequest.Marshall()
	if err != nil {
		log.Println("telethon SignNumber: can`t marshall")
		return err
	}
	// Create a new request object
	request, err := http.NewRequest(sharedUrlMethod, sharedUrl, bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Println("telethon SignNumber: can`t create NewRequest")
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := ts.client.Do(request)
	if err != nil {
		log.Println("telethon SignNumber: error during request")
		return err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusOK {
		return nil
	} else {
		log.Println("telethon SignUpNumber: request return code "+response.Status, ", body:", response.Body)
		return errors.New("request returned code " + response.Status)
	}
}

func (ts *TelethonService) AuthCode(requestBody models.AuthCodeRequest) (models.AuthTelethonFinishResponse, error) {

	jsonBytes, err := requestBody.Marshall()
	if err != nil {
		log.Println("telethon AuthCode: can`t marshall")
		return models.AuthTelethonFinishResponse{}, err
	}

	request, err := http.NewRequest(sharedUrlMethod, sharedUrl, bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Println("telethon AuthCode: can`t create NewRequest")
		return models.AuthTelethonFinishResponse{}, err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := ts.client.Do(request)
	if err != nil {
		log.Println("telethon AuthCode: error during request")
		return models.AuthTelethonFinishResponse{}, err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusForbidden { //403 - request for 2FA password to continue telegram login
		log.Println("telethon AuthCode: request returned code " + response.Status)
		return models.AuthTelethonFinishResponse{}, &ErrorUserHas2FA{message: "Need 2FA password", code: 403}
	} else if response.StatusCode != http.StatusOK {
		log.Println("telethon AuthCode: request returned code " + response.Status)
		return models.AuthTelethonFinishResponse{}, errors.New("request returned code " + response.Status)
	}
	//successful login to telegram, no 2FA password needed
	var codeResponse models.AuthTelethonFinishResponse
	err = json.NewDecoder(response.Body).Decode(&codeResponse)
	if err != nil {
		log.Println("telethon AuthCode: error during reading response body")
		return models.AuthTelethonFinishResponse{}, err
	}
	return codeResponse, nil
}

func (ts *TelethonService) Auth2FA(requestBody models.Auth2FARequest) (models.AuthTelethonFinishResponse, error) {

	jsonBytes, err := requestBody.Marshall()
	if err != nil {
		log.Println("telethon Auth2FA: can`t marshall")
		return models.AuthTelethonFinishResponse{}, err
	}
	request, err := http.NewRequest("POST", "http://localhost:8081/auth/2fa", bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Println("telethon Auth2FA: can`t create NewRequest")
		return models.AuthTelethonFinishResponse{}, err
	}
	request.Header.Set("Content-Type", "application/json")
	//response, err := ts.client.Do(request)
	response := http.Response{StatusCode: http.StatusOK}
	if err != nil {
		log.Println("telethon Auth2FA: error during request")
		return models.AuthTelethonFinishResponse{}, err
	}
	//defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Println("telethon Auth2FA: request returned code " + response.Status)
		return models.AuthTelethonFinishResponse{}, errors.New("request returned code " + response.Status)
	}
	var signUpResponse models.AuthTelethonFinishResponse
	err = json.NewDecoder(response.Body).Decode(&signUpResponse)
	if err != nil {
		log.Println("telethon Auth2FA: error during reading response body")
		return models.AuthTelethonFinishResponse{}, err
	}
	return signUpResponse, nil
}

func (ts *TelethonService) CheckClient(requestBody models.AuthCheckClientRequest) error {
	jsonBytes, err := requestBody.Marshall()
	if err != nil {
		log.Println("telethon CheckClient: can`t marshall")
		return err
	}
	request, err := http.NewRequest(sharedUrlMethod, sharedUrl, bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Println("telethon CheckClient: can`t create NewRequest")
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := ts.client.Do(request)
	if err != nil {
		log.Println("telethon CheckClient: error during request")
		return err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusForbidden {
		log.Println("telethon CheckClient: client is gone, request returned code " + response.Status)
		return &ErrorUserClientIsGone{message: "Client is gone"}
	}
	if response.StatusCode != http.StatusOK {
		log.Println("telethon CheckClient: request returned code " + response.Status)
		return errors.New("request returned code " + response.Status)
	}
	return nil
}

func (ts *TelethonService) GetPublicUserChannels(requestBody models.GetChannelsRequest) ([]models.Channel, error) {
	jsonBytes, err := requestBody.Marshall()
	if err != nil {
		log.Println("telethon GetPublicUserChannels: can`t marshall")
		return []models.Channel{}, err
	}
	request, err := http.NewRequest(sharedUrlMethod, sharedUrl, bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Println("telethon GetPublicUserChannels: can`t create NewRequest")
		return []models.Channel{}, err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := ts.client.Do(request)
	if err != nil {
		log.Println("telethon GetPublicUserChannels: error during request")
		return []models.Channel{}, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		log.Println("telethon GetPublicUserChannels: request returned code " + response.Status)
		return []models.Channel{}, errors.New("request returned code " + response.Status)
	}
	var result models.GetPublicUserChannelsResult
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		log.Println("telethon GetPublicUserChannels: error during reading response body")
		return []models.Channel{}, err
	}
	return result.Details, nil
}

func (ts *TelethonService) FormFeed(requestBody models.FormFeedRequest) (models.FormFeedResponse, error) {
	jsonBytes, err := json.Marshal(requestBody)
	if err != nil {
		log.Println("telethon FormFeed: can`t marshall")
		return models.FormFeedResponse{}, err
	}
	request, err := http.NewRequest(sharedUrlMethod, sharedUrl, bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Println("telethon FormFeed: can`t create NewRequest")
		return models.FormFeedResponse{}, err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := ts.client.Do(request)
	if err != nil {
		log.Println("telethon FormFeed: error during request")
		return models.FormFeedResponse{}, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		log.Println("telethon FormFeed: request returned code " + response.Status)
		return models.FormFeedResponse{}, errors.New("request returned code " + response.Status)
	}
	var feed models.FormFeedResponse
	err = json.NewDecoder(response.Body).Decode(&feed)
	if err != nil {
		log.Println("telethon FormFeed: error during reading response body")
		return models.FormFeedResponse{}, err
	}
	return feed, nil
}

func (ts *TelethonService) GetRecommendation(requestBody models.GetRecommendationRequest) ([]string, error) {
	jsonBytes, err := json.Marshal(requestBody)
	if err != nil {
		log.Println("telethon GetRecommendation: can`t marshall")
		return nil, err
	}
	request, err := http.NewRequest(sharedUrlMethod, sharedUrl, bytes.NewBuffer(jsonBytes))
	if err != nil {
		log.Println("telethon GetRecommendation: can`t create NewRequest")
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := ts.client.Do(request)
	if err != nil {
		log.Println("telethon GetRecommendation: error during request")
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		log.Println("telethon GetRecommendation: request returned code " + response.Status)
		return nil, errors.New("request returned code " + response.Status)
	}
	var result models.GetRecommendationRequestResult
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		log.Println("telethon GetRecommendation: error during reading response body")
		return nil, err
	}
	return result.Details, nil
}

type ErrorUserHas2FA struct {
	message string
	code    int
}

type ErrorUserClientIsGone struct {
	message string
}

func (e ErrorUserHas2FA) Error() string {
	return fmt.Sprintf("ErrorUserHas2FA: %s (code: %d)", e.message, e.code)
}

func (e ErrorUserClientIsGone) Error() string {
	return fmt.Sprintf("ErrorUserClientIsGone: %s", e.message)
}
