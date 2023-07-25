package api

import (
	"NewsFeed/pkg/models"
	"NewsFeed/pkg/services"
	"NewsFeed/pkg/services/auth"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type AuthAPI struct {
	AuthService *auth.Service
}

func NewAuthAPI(AuthService *auth.Service) *AuthAPI {
	return &AuthAPI{AuthService: AuthService}
}

const userIDContextKey = "uid"

func (api *AuthAPI) SignUpNumber(ctx *gin.Context, r *http.Request) {
	if r.Method != http.MethodPost {
		ctx.AbortWithStatusJSON(http.StatusMethodNotAllowed, "Invalid request method")
		return
	}

	var signUpNumberBody AuthNumberBody
	err := json.NewDecoder(r.Body).Decode(&signUpNumberBody)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, "Invalid request body")
		log.Println("Error:", err)
		return
	}
	tokens, err := api.AuthService.SignUpNumber(signUpNumberBody.Number)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, tokens)
}

func (api *AuthAPI) RefreshToken(ctx *gin.Context, r *http.Request) {
	if r.Method != http.MethodPost {
		ctx.AbortWithStatusJSON(http.StatusMethodNotAllowed, "Invalid request method")
		return
	}
	var refreshTokenBody RefreshTokenBody
	err := json.NewDecoder(r.Body).Decode(&refreshTokenBody)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, "Invalid request body")
		return
	}
	tokens, err := api.AuthService.RefreshToken(refreshTokenBody.RefreshToken)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, tokens)
}

func (api *AuthAPI) AuthMW() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.Request.Header.Get("Authorization")
		if strings.HasPrefix(token, "Bearer ") {
			token = strings.TrimPrefix(token, "Bearer ")
		} else {
			log.Printf("AuthMW: can`t get token from header")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		userID, err := api.AuthService.RetrieveIDFromToken(token)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Set(userIDContextKey, userID)
	}
}

func popUserIDFromContext(ctx *gin.Context) (uint32, error) {
	userIDValue, ok := ctx.Get(userIDContextKey)
	if !ok {
		return 0, errors.New("no user id in ctx")
	}
	userID, err := strconv.ParseUint(fmt.Sprint(userIDValue), 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(userID), nil
}

func (api *AuthAPI) SignInNumber(ctx *gin.Context, r *http.Request) {
	if r.Method != http.MethodPost {
		ctx.AbortWithStatusJSON(http.StatusMethodNotAllowed, "Invalid request method")
		return
	}
	var signUpNumberBody AuthNumberBody
	err := json.NewDecoder(r.Body).Decode(&signUpNumberBody)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, "Invalid request body")
		return
	}
	userID, err := popUserIDFromContext(ctx)
	if err != nil {
		log.Println("AuthAPI SignInNumber: can`t retrieve user id from gin.Context")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	err = api.AuthService.SignInNumber(signUpNumberBody.Number, userID)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.Status(http.StatusOK)
}

func (api *AuthAPI) AuthCode(ctx *gin.Context, r *http.Request) {
	if r.Method != http.MethodPost {
		ctx.AbortWithStatusJSON(http.StatusMethodNotAllowed, "Invalid request method")
		return
	}
	var signUpCodeBody AuthCodeBody
	err := json.NewDecoder(r.Body).Decode(&signUpCodeBody)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, "Invalid request body")
		return
	}
	userID, err := popUserIDFromContext(ctx)
	if err != nil {
		log.Println("AuthAPI AuthSignUpCode: can`t retrieve user id from gin.Context")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	userInfo, err := api.AuthService.AuthCode(userID, signUpCodeBody.Code)
	if err != nil {
		var userHas2FAErr *services.ErrorUserHas2FA
		if errors.As(err, &userHas2FAErr) {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, userInfo)
}

func (api *AuthAPI) Auth2FA(ctx *gin.Context, r *http.Request) {
	if r.Method != http.MethodPost {
		ctx.AbortWithStatusJSON(http.StatusMethodNotAllowed, "Invalid request method")
		return
	}
	var signUp2FABody Auth2FABody
	err := json.NewDecoder(r.Body).Decode(&signUp2FABody)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, "Invalid request body")
		return
	}
	userID, err := popUserIDFromContext(ctx)
	if err != nil {
		log.Println("AuthAPI AuthSignUpCode: can`t retrieve user id from gin.Context")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	userInfo, err := api.AuthService.Auth2FA(userID, signUp2FABody.Password)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.JSON(http.StatusOK, userInfo)
}

func (api *AuthAPI) CheckClient(ctx *gin.Context, r *http.Request) {
	if r.Method != http.MethodGet {
		ctx.AbortWithStatusJSON(http.StatusMethodNotAllowed, "Invalid request method")
		return
	}
	userID, err := popUserIDFromContext(ctx)
	if err != nil {
		log.Println("AuthAPI CheckClient: can`t retrieve user id from gin.Context")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	err = api.AuthService.CheckClient(userID)
	if err != nil {
		var errorUserClientIsGone *services.ErrorUserClientIsGone
		if errors.As(err, &errorUserClientIsGone) {
			ctx.AbortWithStatus(http.StatusForbidden)
			return
		}
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.AbortWithStatus(http.StatusOK)
}

type RefreshTokenBody struct {
	RefreshToken string `json:"refresh_token"`
}

type AuthNumberBody struct {
	Number models.Number `json:"number"`
}

type AuthCodeBody struct {
	Code string `json:"code"`
}

type Auth2FABody struct {
	Password string `json:"password"`
}
