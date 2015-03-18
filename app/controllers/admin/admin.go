package admin

import (
	"fmt"
	"gitcent-web/app/controllers"

	"github.com/revel/revel"
)

type Admin struct {
	controllers.App
}

func (c Admin) Index() revel.Result {
	user, _ := c.RequireUser()

	fmt.Println(user.Name)
	return c.RenderTemplate("App/admin/index.html")
}
