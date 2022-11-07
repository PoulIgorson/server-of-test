// Package functions implements handlers for pages.
package functions

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"

	sawed "dom50b_fiberWoodMonitor/buckets/sawed"
	user "dom50b_fiberWoodMonitor/buckets/user"
	db "dom50b_fiberWoodMonitor/datebase"
	. "dom50b_fiberWoodMonitor/define"
	support "dom50b_fiberWoodMonitor/support"
)

// updateData update value in record of bucket of key.
func updateData(b *db.Bucket, key, value string) error {
	datastr, err := b.GetOfField("date", GetToday())
	if err != nil {
		return err
	}
	var saw sawed.Sawed
	json.Unmarshal([]byte(datastr), &saw)
	ckey, _ := Atoi(key)
	cvalue, _ := Atoi(value)
	saw.Sizes[uint(ckey)] = uint8(cvalue)
	return saw.Save(b)
}

// IndexPage returns handler for index page.
func IndexPage(db_ *db.DB, urls interface{}) fiber.Handler {

	return func(c *fiber.Ctx) error {
		context := fiber.Map{
			"pagename": "Главная",
			"menu":     urls,
		}
		var cuser user.User
		userStr := c.Cookies("userCookie")
		if userStr != "" {
			json.Unmarshal([]byte(userStr), &cuser)
			context["user"] = cuser
		}
		return c.Render("index", context)
	}
}

// MonitorPage returns handler for monitor page.
func MonitorPage(db_ *db.DB, urls interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var cuser user.User
		userStr := c.Cookies("userCookie")
		json.Unmarshal([]byte(userStr), &cuser)
		saw, err := db_.Bucket("sawed")
		if err != nil {
			return c.Render("monitor", fiber.Map{
				"pagename": "Монитор",
				"error":    err.Error(),
				"menu":     urls,
				"user":     cuser,
			})
		}

		sizes, err := support.GetLastData(saw)
		if err != nil {
			return c.Render("monitor", fiber.Map{
				"pagename": "Монитор",
				"error":    err.Error(),
				"menu":     urls,
				"user":     cuser,
			})
		}
		return c.Render("monitor", fiber.Map{
			"pagename": "Монитор",
			"sizes":    sizes,
			"menu":     urls,
			"user":     cuser,
		})
	}
}

// APIUpdate returns handler for api update.
func APIUpdate(db_ *db.DB, urls interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		saw, err := db_.Bucket("sawed")
		if err != nil {
			return c.JSON(fiber.Map{"errors": []string{err.Error()}})
		}
		err = updateData(saw, c.Query("key"), c.Query("value"))
		status := fiber.StatusOK
		desc := "OK"
		if err != nil {
			status = fiber.StatusBadRequest
			desc = err.Error()
		}
		return c.JSON(fiber.Map{
			"statusCode":  status,
			"description": desc,
		})
	}
}

// BucketPage returns handler for bucket page.
func BucketPage(db_ *db.DB, urls interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var cuser user.User
		userStr := c.Cookies("userCookie")
		json.Unmarshal([]byte(userStr), &cuser)
		var lines []string
		if !db_.ExistsBucket(c.Params("name")) {
			return c.Render("bucket", fiber.Map{
				"pagename": "Bucket: " + c.Params("name"),
				"error":    fmt.Sprintf("Bucket with name `%v` does not exists", c.Params("name")),
				"menu":     urls,
				"user":     cuser,
			})
		}
		bct, err := db_.Bucket(c.Params("name"))
		if err != nil {
			return c.Render("bucket", fiber.Map{
				"pagename": "Bucket: " + c.Params("name"),
				"error":    err.Error(),
				"menu":     urls,
				"user":     cuser,
			})
		}
		for inc := 1; true; inc++ {
			v, err := bct.Get(inc)
			if err != nil {
				break
			}
			lines = append(lines, v)
		}
		return c.Render("bucket", fiber.Map{
			"pagename": "Bucket: " + c.Params("name"),
			"lines":    lines,
			"menu":     urls,
			"user":     cuser,
		})
	}
}
