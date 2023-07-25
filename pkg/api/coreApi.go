package api

import (
	"NewsFeed/pkg/config"
)

type Api struct {
	AuthApi *AuthAPI
	UserApi *UserAPI
	FeedAPI *FeedAPI
	cnf     config.ServiceConfiguration
}

func NewAPI(authApi *AuthAPI, UserApi *UserAPI, FeedApi *FeedAPI, cnf config.ServiceConfiguration) *Api {
	return &Api{AuthApi: authApi, UserApi: UserApi, FeedAPI: FeedApi, cnf: cnf}
}
