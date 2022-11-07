// Package support implements support functions.
package support

import (
	"encoding/json"
	"fmt"
	"strings"

	fiber "github.com/gofiber/fiber/v2"

	sawed "dom50b_fiberWoodMonitor/buckets/sawed"
	db "dom50b_fiberWoodMonitor/datebase"
	. "dom50b_fiberWoodMonitor/define"
)

// GetLastData return last record in bucket.
func GetLastData(b *db.Bucket) (fiber.Map, error) {
	today := GetToday()
	value, err := b.GetOfField("date", today)
	if err != nil {
		if value == "ErrorJSON" {
			return nil, err
		}
		sizesStr, err := b.Get(-1)
		if err != nil {
			if err.Error() == "Key `-1` is not exists" {
				return nil, fmt.Errorf("`Sizes` does not exists")
			}
			return nil, err
		}
		if sizesStr == "" {
			return nil, fmt.Errorf("`Sizes` does not exists")
		}
		sizesList := strings.Split(sizesStr, " ")
		sizes := map[uint]uint8{}
		for _, size := range sizesList {
			intSize, _ := Atoi(size)
			sizes[uint(intSize)] = 0
		}
		s := sawed.Sawed{
			ID:    0,
			Date:  today,
			Sizes: sizes,
		}
		if err := s.Save(b); err != nil {
			return nil, err
		}
		value, err = b.GetOfField("date", today)
		if err != nil {
			return nil, err
		}
	}
	var saw sawed.Sawed
	err = json.Unmarshal([]byte(value), &saw)
	if err != nil {
		return nil, err
	}
	sizes_dict := saw.Sizes
	sizes := fiber.Map{}
	for k, v := range sizes_dict {
		sizes[Itoa(int(k))] = Itoa(int(v))
	}
	return sizes, nil
}
