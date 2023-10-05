package handler

import (
	"github.com/helloferdie/golib/libecho"
	"github.com/helloferdie/pusher/pkg/deploy/request"
	"github.com/helloferdie/pusher/pkg/deploy/service"
	"github.com/labstack/echo/v4"
)

// Deploy -
func Deploy(c echo.Context) (err error) {
	r := new(request.Deploy)
	if err = c.Bind(r); err != nil {
		return
	}

	resp := service.Deploy(r)
	return libecho.ParseResponse(c, resp)
}
