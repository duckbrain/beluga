package actions

import "github.com/gobuffalo/buffalo"

func AuthMiddleware(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		return next(c)
	}
}
