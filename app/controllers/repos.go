package controllers

import "github.com/revel/revel"

type Repos struct {
	App
}

func (r Repos) Index() revel.Result {
	return r.Render()
}
