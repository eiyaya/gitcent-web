package controllers

import (
	"gitcent-web/app/model"

	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) RequireUser() (*model.User, error) {
	user := &model.User{Name: "eiyaya", Email: "eiyaya@gmail.com"}

	return user, nil
}
