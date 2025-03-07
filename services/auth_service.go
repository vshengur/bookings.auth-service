package services

import (
	"context"
	"fmt"
	"log"

	goauth2 "golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleoauth "google.golang.org/api/oauth2/v2"
	option "google.golang.org/api/option"

	"github.com/vshengur/bookings.auth-service/models"
)

type AuthService struct {
	oauthConfig *goauth2.Config
}

func NewAuthService(clientID, clientSecret, redirectURL string) *AuthService {
	return &AuthService{
		oauthConfig: &goauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{googleoauth.UserinfoEmailScope, googleoauth.UserinfoProfileScope},
			Endpoint:     google.Endpoint,
		},
	}
}

func (a *AuthService) GenerateAuthURL() string {
	return a.oauthConfig.AuthCodeURL("state-token", goauth2.AccessTypeOffline)
}

func (a *AuthService) HandleGoogleCallback(code string) (*models.User, error) {
	// Обмен кода на токен доступа
	token, err := a.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	log.Println("Google auth token :", token)

	// Создание HTTP клиента с токеном
	client := a.oauthConfig.Client(context.Background(), token)
	service, err := googleoauth.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("failed to create oauth2 service: %w", err)
	}

	// Получение информации о пользователе
	userInfo, err := service.Userinfo.V2.Me.Get().Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	return &models.User{
		Email:    userInfo.Email,
		FullName: userInfo.Name,
	}, nil
}
