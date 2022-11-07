// Package functions implements handlers for registration pages.
package functions

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"

	//gofiberfirebaseauth "github.com/sacsand/gofiber-firebaseauth"

	user "dom50b_fiberWoodMonitor/buckets/user"
	db "dom50b_fiberWoodMonitor/datebase"
	. "dom50b_fiberWoodMonitor/define"
)

func dropPort(host string) string {
	index := strings.IndexRune(host, ':')
	if index != -1 {
		host = host[:index]
	}
	return host
}

// LoginPage returns handler for login page.
func LoginPage(db_ *db.DB, urls interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		context := fiber.Map{
			"pagename": "Login",
			"menu":     urls,
		}
		errors := map[string]string{}
		if c.Method() == "GET" {
			userStr := c.Cookies("userCookie")
			if user.CheckUser(db_, userStr) != nil {
				return c.Redirect("/")
			}
		} else if c.Method() == "POST" {
			var data map[string]string
			json.Unmarshal(c.Request().Body(), &data)
			context["login"] = data["login"]
			context["password"] = data["password"]

			users, err := db_.Bucket("users")
			if err != nil {
				context["error"] = err.Error()
			} else {
				value, err := users.GetOfField("login", data["login"])
				if err != nil {
					errors["login"] = "Логин не существует"
				}

				if len(errors) == 0 && context["error"] == nil {
					var cuser user.User
					json.Unmarshal([]byte(value), &cuser)

					if Hash([]byte(data["password"])) != cuser.Password {
						errors["password"] = "Неверный пароль"
					} else {
						cookie := fiber.Cookie{
							Name:   "userCookie",
							Value:  fmt.Sprintf(`{"login": "%v", "password": "%v"}`, data["login"], cuser.Password),
							Path:   "/",
							Domain: dropPort(c.Hostname()),
							//Expires:     time.Now().Add(time.Hour),
							Secure:      false,
							HTTPOnly:    false,
							SessionOnly: false,
						}
						c.Cookie(&cookie)
						//return c.Redirect("/")
						url := "/"
						if cuser.Role == user.Admin {
							url = "/admin"
						}
						if cuser.Role == user.Worker {
							url = "/monitor"
						}
						return c.Status(302).JSON(fiber.Map{
							"redirectURL": url,
						})
					}
				}
			}
		}
		context["errors"] = errors
		return c.Render("registration/login", context)
	}
}

func APILogout(db_ *db.DB, urls interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.ClearCookie("userCookie")
		return c.Redirect("/")
	}
}

func RegistrationPage(db_ *db.DB, urls interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		context := fiber.Map{
			"pagename": "Регистрация",
			"menu":     urls,
		}
		errors := map[string]string{}
		if c.Method() == "GET" {
			userStr := c.Cookies("userCookie")
			if user.CheckUser(db_, userStr) != nil {
				return c.Redirect("/")
			}
		} else if c.Method() == "POST" {
			var data map[string]string
			json.Unmarshal(c.Request().Body(), &data)
			context["login"] = data["login"]
			context["email"] = data["email"]
			context["password1"] = data["password1"]
			context["password2"] = data["password2"]

			users, err := db_.Bucket("users")
			if err != nil {
				context["error"] = err.Error()
			} else {
				if len(data["login"]) < 4 {
					fmt.Println(data["login"])
					errors["login"] = "Слишком короткий логин"
				}

				if len(data["password1"]) < 8 {
					errors["password1"] = "Слишком короткий пароль"
				}

				if len(data["password2"]) < 8 {
					errors["password2"] = "Слишком короткий пароль"
				}

				_, err := users.GetOfField("login", data["login"])
				if err == nil && errors["login"] == "" {
					errors["login"] = "Логин существует"
				}

				if !VerifyEmail(data["email"]) {
					errors["email"] = "Некоректный адрес элекстронной почты"
				}

				if data["password1"] != data["password2"] && errors["password1"] == errors["password2"] && errors["password1"] == "" {
					errors["password2"] = "Пароли не совпадают"
				}

				if len(errors) == 0 && context["error"] == nil {
					cuser := user.User{
						Login:    data["login"],
						Email:    data["email"],
						Password: Hash([]byte(data["password1"])),
					}
					cuser.Save(users)
					return c.Status(302).JSON(fiber.Map{"redirectURL": "/login"})
				}
			}
		}
		context["errors"] = errors
		return c.Render("registration/registration", context)
	}
}
