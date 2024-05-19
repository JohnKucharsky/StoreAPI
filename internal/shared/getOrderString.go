package shared

import "fmt"

func GetOrderString(orderBy string, sortOrder string) string {
	if orderBy != "" {
		return fmt.Sprintf("order by %s %s", orderBy, sortOrder)
	}
	return ""
}
