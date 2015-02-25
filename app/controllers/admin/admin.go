package admin

import "github.com/revel/revel"

type Admin struct {
	*revel.Controller
}

func (c Admin) Index() revel.Result {
	return c.RenderTemplate("App/admin/index.html")
}
