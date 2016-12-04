package middleware

import (
	"fmt"
	"html/template"
	"strings"
	"time"

	"strconv"

	"github.com/admpub/caddyui/application/library/errors"
	"github.com/admpub/caddyui/application/library/modal"
	"github.com/webx-top/echo"
)

func FuncMap() echo.MiddlewareFunc {
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			c.SetFunc(`Now`, time.Now)
			c.SetFunc(`Add`, func(x int, y int) int {
				return x + y
			})
			c.SetFunc(`Sub`, func(x int, y int) int {
				return x - y
			})
			c.SetFunc(`HasPrefix`, strings.HasPrefix)
			c.SetFunc(`HasSuffix`, strings.HasSuffix)
			c.SetFunc(`HasString`, hasString)
			c.SetFunc(`Date`, date)
			c.SetFunc(`Modal`, func(data interface{}) template.HTML {
				return modal.Render(c, data)
			})
			c.SetFunc(`IsMessage`, errors.IsMessage)
			c.SetFunc(`IsError`, errors.IsError)
			c.SetFunc(`IsOk`, errors.IsOk)
			c.SetFunc(`Message`, errors.Message)
			c.SetFunc(`Ok`, errors.Ok)
			return h.Handle(c)
		})
	}
}

func hasString(slice []string, str string) bool {
	if slice == nil {
		return false
	}
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func date(timestamp interface{}) time.Time {
	if v, y := timestamp.(int64); y {
		return time.Unix(v, 0)
	}
	if v, y := timestamp.(uint); y {
		return time.Unix(int64(v), 0)
	}
	v, _ := strconv.ParseInt(fmt.Sprint(timestamp), 10, 64)
	return time.Unix(v, 0)
}
