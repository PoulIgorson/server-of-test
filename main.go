package main

import (
	"encoding/json"
	"fmt"

	//firebase "firebase.google.com/go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/csrf"
	"github.com/gofiber/template/django"

	mauth "dom50b_fiberWoodMonitor/auth"
	db "dom50b_fiberWoodMonitor/datebase"
	router "dom50b_fiberWoodMonitor/router"
)

// GetResponse returns response of url
func GetResponse(method, curl string, data fiber.Map) (fiber.Map, []error) {
	agent := fiber.AcquireAgent()
	request := agent.Request()
	request.Header.SetMethod(method)
	request.SetRequestURI(curl)
	agent.JSON(data)

	if err := agent.Parse(); err != nil {
		return nil, []error{err}
	}

	_, body, errs := agent.Bytes()
	if len(errs) != 0 {
		return nil, errs
	}

	var response fiber.Map
	json.Unmarshal([]byte(fmt.Sprintf("%s", body)), &response)
	return response, nil
}

func GetOfIndex(val interface{}, ind interface{}) interface{} {
	lst := val.([]interface{})
	index := ind.(int)
	return lst[index]
}

func GetOfMapIndex(val interface{}, typ1 interface{}, typ2 interface{}, ind interface{}) interface{} {
	switch typ1.(string) {
	case "string":
		switch typ2.(string) {
		case "string":
			lst := val.(map[string]string)
			index := ind.(string)
			return lst[index]
		case "int":
			lst := val.(map[string]int)
			index := ind.(string)
			return lst[index]
		case "uint":
			lst := val.(map[string]uint)
			index := ind.(string)
			return lst[index]
		case "uint8":
			lst := val.(map[string]uint8)
			index := ind.(string)
			return lst[index]
		}
	case "int":
		switch typ2.(string) {
		case "string":
			lst := val.(map[int]string)
			index := ind.(int)
			return lst[index]
		case "int":
			lst := val.(map[int]int)
			index := ind.(int)
			return lst[index]
		case "uint":
			lst := val.(map[int]uint)
			index := ind.(int)
			return lst[index]
		case "uint8":
			lst := val.(map[int]uint8)
			index := ind.(int)
			return lst[index]
		}
	case "uint":
		switch typ2.(string) {
		case "string":
			lst := val.(map[uint]string)
			index := ind.(uint)
			return lst[index]
		case "int":
			lst := val.(map[uint]int)
			index := ind.(uint)
			return lst[index]
		case "uint":
			lst := val.(map[uint]uint)
			index := ind.(uint)
			return lst[index]
		case "uint8":
			lst := val.(map[uint]uint8)
			index := ind.(uint)
			return lst[index]
		}
	case "uint8":
		switch typ2.(string) {
		case "string":
			lst := val.(map[uint8]string)
			index := ind.(uint8)
			return lst[index]
		case "int":
			lst := val.(map[uint8]int)
			index := ind.(uint8)
			return lst[index]
		case "uint":
			lst := val.(map[uint8]uint)
			index := ind.(uint8)
			return lst[index]
		case "uint8":
			lst := val.(map[uint8]uint8)
			index := ind.(uint8)
			return lst[index]
		}
	}
	return ""
}

func main() {
	engine := django.New("./templates", ".django")
	engine.AddFunc("index", GetOfIndex)
	engine.AddFunc("indexMap", GetOfMapIndex)

	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Use(csrf.New(csrf.ConfigDefault))

	app.Static("/static", "./static")

	db_, err := db.Open("wood_monitor_db.db")
	app.Use(mauth.New(db_))
	if err != nil {
		panic(err)
	}
	defer db_.Close()

	router.Router(app, db_)

	app.Listen("")
}
