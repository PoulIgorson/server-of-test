// Package router implements setting routing for site.
package router

import (
	"github.com/gofiber/fiber/v2"

	db "dom50b_fiberWoodMonitor/datebase"

	urls "dom50b_fiberWoodMonitor/router/urls"
)

// Router setting handlers on url
func Router(app *fiber.App, db_ *db.DB) {
	for _, url := range urls.UrlPatterns {
		switch url.Method {
		case "Get":
			app.Get(url.Path, url.Handler(db_, urls.UrlPatterns))
		case "All":
			app.All(url.Path, url.Handler(db_, urls.UrlPatterns))
		default:
			app.Add(url.Method, url.Path, url.Handler(db_, urls.UrlPatterns))
		}
	}
}
