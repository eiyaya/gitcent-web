package controllers

import (
	"errors"
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
	user := &model.User{Name: "", Email: ""}
	user.Name = c.Session["username"]
	if user.Name == "" {
		return nil, errors.New("not a login user")
	} else {
		return user, nil
	}
}
