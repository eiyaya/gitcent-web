package admin

import (
	"gitcent-web/app/controllers"

	"github.com/revel/revel"
)

type Admin struct {
	controllers.App
}

func (c Admin) Index() revel.Result {
	user, err := c.RequireUser()
	if err != nil {
		return c.RenderText("need login")
	}
	c.RenderArgs["user"] = user
	return c.RenderTemplate("App/admin/index.html")
}
