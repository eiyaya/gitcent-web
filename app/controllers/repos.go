package controllers

import "github.com/revel/revel"

type Repos struct {
	*revel.Controller
}

func (r Repos) Index() revel.Result {
	return r.Render()
}
