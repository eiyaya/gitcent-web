package controllers

import (
	"gitcent-web/app/model"

	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	//fmt.Println(c.Session["username"] + "logined")
	c.FlashParams()
	c.Flash.Success("I See You!")
	return c.Render()
}

func (c App) RequireUser() (*model.User, error) {
	user := &model.User{Name: "eiyaya", Email: "eiyaya@gmail.com"}

	return user, nil
}
