package services

import (
	"gitcent-web/app/model"
)

// UserService user agent
type UserService struct{}

// RequireUser query session
func (u UserService) RequireUser() (*model.User, error) {
	user := &model.User{}
	return user, nil
}
