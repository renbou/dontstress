package utils

import (
	"encoding/json"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"io/ioutil"
	"os"
)

var LangsSet = getLangs()

func getLangs() map[string]bool {
	jsonFile, err := os.Open("/opt/langs.json")
	if err != nil {
		return nil
	}
	defer jsonFile.Close()

	content, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil
	}

	var langs []string
	err = json.Unmarshal(content, &langs)
	if err != nil {
		return nil
	}
	langsSet := map[string]bool{}
	for _, lang := range langs {
		langsSet[lang] = true
	}
	return langsSet
}

func ValidLang(lang string) bool {
	return LangsSet != nil && LangsSet[lang]
}

func GetId() string {
	return uuid.New().String()
}

func Check(c *fiber.Ctx, err error) bool {
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
		return false
	}
	return true
}

func Validate(c *fiber.Ctx, s interface{}) bool {
	validate := validator.New()
	if err := validate.Struct(s); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid data sent in the request (unsupported, required fields missing etc)",
		})
		return false
	}
	return true
}
