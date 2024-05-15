package shared

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
)

func GetIntArrFromOriginalURL(c *fiber.Ctx, key string) (*[]int, error) {
	originalURL := c.OriginalURL()
	check := c.Query(key)
	if check == "" {
		return nil, errors.New(fmt.Sprintf("no param %s", key))
	}
	firstSplit := strings.Split(originalURL, "?")
	secondSplit := strings.Split(firstSplit[1], "&")
	var intArr []int
	for _, str := range secondSplit {
		idStringArr := strings.Split(str, "=")
		idInt, err := strconv.Atoi(idStringArr[1])
		if err == nil && idStringArr[0] == key {
			intArr = append(intArr, idInt)
		}
	}

	return &intArr, nil
}
