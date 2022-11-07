// Package functions implements handlers for admin pages.
package functions

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"

	"dom50b_fiberWoodMonitor/buckets/sawed"
	user "dom50b_fiberWoodMonitor/buckets/user"
	db "dom50b_fiberWoodMonitor/datebase"
)

var SpecHandlerForBucket = map[string]func(*db.DB, interface{}) fiber.Handler{
	"sawed": BucketSawedPage,
}

// IndexPage returns handler for admin index page.
func IndexPage(db_ *db.DB, urls interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		context := fiber.Map{
			"pagename": "Админ",
			"menu":     urls,
		}
		saw, err := db_.Bucket("sawed")
		if err != nil {
			return c.JSON(fiber.Map{"Error": err.Error()})
		}
		sizes, err := saw.Get(-1)
		if err != nil {
			if err.Error() != "Key `-1` is not exists" {
				return c.JSON(fiber.Map{"Error": err.Error()})
			}
		}
		context["sizes"] = strings.Split(sizes, " ")
		var cuser user.User
		userStr := c.Cookies("userCookie")
		if userStr != "" {
			json.Unmarshal([]byte(userStr), &cuser)
			context["user"] = cuser
		}
		return c.Render("admin/index", context)
	}
}

// BucketPage returns handler for admin bucket page.
func BucketPage(db_ *db.DB, urls interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if true || SpecHandlerForBucket[c.Params("name")] != nil {
			return SpecHandlerForBucket[c.Params("name")](db_, urls)(c)
		}
		var cuser user.User
		userStr := c.Cookies("userCookie")
		json.Unmarshal([]byte(userStr), &cuser)
		if !db_.ExistsBucket(c.Params("name")) {
			return c.Render("admin/bucket", fiber.Map{
				"pagename": "Bucket: " + c.Params("name"),
				"error":    fmt.Sprintf("Bucket with name `%v` does not exists", c.Params("name")),
				"menu":     urls,
				"user":     cuser,
			})
		}
		bct, err := db_.Bucket(c.Params("name"))
		if err != nil {
			return c.Render("admin/bucket", fiber.Map{
				"pagename": "Bucket: " + c.Params("name"),
				"error":    err.Error(),
				"menu":     urls,
				"user":     cuser,
			})
		}
		var data []string
		date := c.Query("date")
		if date == "" {
			for inc := 1; true; inc++ {
				v, err := bct.Get(inc)
				if err != nil {
					break
				}
				data = append(data, v)
			}
		} else {
			datelst := strings.Split(date, "-")
			v, err := bct.GetOfField("date", date)
			if err == nil {
				data = append(data, v)
			}
			date = fmt.Sprintf("%v-%v-%v", datelst[2], datelst[1], datelst[0])
		}
		return c.Render("admin/bucket", fiber.Map{
			"pagename": "Bucket: " + c.Params("name"),
			"data":     data,
			"date":     date,
			"menu":     urls,
			"user":     cuser,
		})
	}
}

// BucketSawedPage returns handler for admin bucket sawed page.
func BucketSawedPage(db_ *db.DB, urls interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var cuser user.User
		userStr := c.Cookies("userCookie")
		json.Unmarshal([]byte(userStr), &cuser)
		bct, err := db_.Bucket("sawed")
		if err != nil {
			return c.Render("admin/bucket", fiber.Map{
				"pagename": "Bucket: " + c.Params("name"),
				"error":    err.Error(),
				"menu":     urls,
				"user":     cuser,
			})
		}
		var data []sawed.Sawed
		date := c.Query("date")
		if date == "" {
			for inc := 1; true; inc++ {
				v, err := bct.Get(inc)
				if err != nil {
					break
				}
				var csawed sawed.Sawed
				json.Unmarshal([]byte(v), &csawed)
				data = append(data, csawed)
			}
		} else {
			datelst := strings.Split(date, "-")
			v, err := bct.GetOfField("date", date)
			if err == nil {
				var csawed sawed.Sawed
				json.Unmarshal([]byte(v), &csawed)
				data = append(data, csawed)
			}
			date = fmt.Sprintf("%v-%v-%v", datelst[2], datelst[1], datelst[0])
		}
		return c.Render("admin/buckets/sawed", fiber.Map{
			"pagename": "Bucket: sawed",
			"data":     data,
			"sizes":    sawed.GetAllSizes(bct),
			"date":     date,
			"menu":     urls,
			"user":     cuser,
		})
	}
}

// APISizes returns handler for admin api page.
func APISizes(db_ *db.DB, urls interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		saw, err := db_.Bucket("sawed")
		if err != nil {
			return c.JSON(fiber.Map{"Error": err.Error()})
		}
		err = saw.Set(-1, c.Query("sizes"))
		if err != nil {
			return c.JSON(fiber.Map{"Error": err.Error()})
		}
		return c.JSON(fiber.Map{"Stutus": "200"})
	}
}
