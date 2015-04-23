package controllers

import "github.com/revel/revel"

type LoginD struct {
	App
}

func (c *LoginD) UserLogin(user string, password string) revel.Result {
	c.Session["username"] = user
	return c.RenderText(user + password)
}

func (c *LoginD) UserLogout() revel.Result {
	c.Session["username"] = ""
	return c.Redirect("/")
}
