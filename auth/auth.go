// Package auth implements interface for auth.
package auth

import (
	"github.com/gofiber/fiber/v2"

	user "dom50b_fiberWoodMonitor/buckets/user"
	db "dom50b_fiberWoodMonitor/datebase"
	. "dom50b_fiberWoodMonitor/define"
	urls "dom50b_fiberWoodMonitor/router/urls"
)

var IgnoreUrls = []string{
	"/", "/login", "/logout", "/registration",
}

// New return handler for auth.
func New(db_ *db.DB) fiber.Handler {
	return MyNew(db_)
}

func MyNew(db_ *db.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userStr := c.Cookies("userCookie")
		if user.CheckUser(db_, userStr) != nil {
			return c.Next()
		}
		return Unauthorized(db_, c)
	}
}

func Unauthorized(db_ *db.DB, c *fiber.Ctx) error {
	if Contains(IgnoreUrls, c.Path()) {
		return urls.GetUrlOfPath(c.Path()).Handler(db_, urls.UrlPatterns)(c)
	}
	return c.Redirect("/login")
}
