// Package sawed implements model of bucket.
package sawed

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	bolt "go.etcd.io/bbolt"

	db "dom50b_fiberWoodMonitor/datebase"
	. "dom50b_fiberWoodMonitor/define"
)

// Sawed presents model of bucket.
type Sawed struct {
	ID    int            `json:"id"`
	Date  string         `json:"date"`
	Sizes map[uint]uint8 `json:"sizes"`
}

// Save implements saving model in bucket.
func (this *Sawed) Save(bucket *db.Bucket) error {
	// if object does not exists
	if _, err := bucket.Get(this.ID); err != nil || this.ID == 0 {
		bucket.DB().View(func(tx *bolt.Tx) error {
			bct := tx.Bucket([]byte(bucket.Name()))
			k, _ := bct.Cursor().Last()
			id, _ := Atoi(string(k))
			this.ID = id + 1
			return nil
		})
	}

	buf, err := json.Marshal(this)
	if err != nil {
		return err
	}
	return bucket.Set(this.ID, string(buf))
}

func GetAllSizes(bucket *db.Bucket) []uint {
	sizes := []uint{}
	for inc := 1; true; inc++ {
		v, err := bucket.Get(inc)
		if err != nil {
			break
		}

		var data fiber.Map
		err = json.Unmarshal([]byte(v), &data)
		if err != nil {
			break
		}

		for sizeStr := range data["sizes"].(map[string]interface{}) {
			size, _ := Atoi(sizeStr)
			if !Contains(sizes, uint(size)) {
				sizes = append(sizes, uint(size))
			}
		}
	}
	return sizes
}
