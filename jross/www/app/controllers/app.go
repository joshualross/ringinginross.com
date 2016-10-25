package controllers

import "github.com/revel/revel"

type App struct {
	*revel.Controller
}

func (c App) Index(test bool) revel.Result {
	return c.Render(test)
}
