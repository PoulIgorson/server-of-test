// Package user implements model of bucket.
package user

import (
	"encoding/json"

	bolt "go.etcd.io/bbolt"

	db "dom50b_fiberWoodMonitor/datebase"
	. "dom50b_fiberWoodMonitor/define"
)

// Role implements access to module site
type Role struct {
	Name   string `json:"name"`
	Access uint   `json:"access"`
}

var Admin = Role{"admin", 2}
var Worker = Role{"worker", 1}

// User presents model of bucket.
type User struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Role     Role   `json:"role"` // Больше значение выше уровень доступа
}

// Save implements saving model in bucket.
func (this *User) Save(bucket *db.Bucket) error {
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

func CheckUser(db_ *db.DB, userStr string) *User {
	cuser := User{}
	if userStr == "" {
		return nil
	}
	json.Unmarshal([]byte(userStr), &cuser)
	users, err := db_.Bucket("users")
	if err != nil {
		return nil
	}
	value, err := users.GetOfField("login", cuser.Login)
	if err != nil {
		return nil
	}
	realUser := User{}
	json.Unmarshal([]byte(value), &realUser)
	if cuser.Password != realUser.Password {
		return nil
	}
	return &realUser
}

/*
type Model interface {
	GSID(...int) int
}

func Save(this *Model, bucket *db.Bucket) error {
	// if object does not exists
	if _, err := bucket.Get((*this).GSID()); err != nil || (*this).GSID() == 0 {
		bucket.DB().View(func(tx *bolt.Tx) error {
			bct := tx.Bucket([]byte(bucket.Name()))
			k, _ := bct.Cursor().Last()
			id, _ := Atoi(string(k))
			(*this).GSID(id + 1)
			return nil
		})
	}

	buf, err := json.Marshal(this)
	if err != nil {
		return err
	}
	return bucket.Set((*this).GSID(), string(buf))
}

func (u *User) GSID(newID ...int) int {
	if len(newID) != 0 {
		u.ID = newID[0]
	}
	return u.ID
}

func f() {
	m := Model()
	m = User{}
	Save(&u, &db.Bucket{})
}
*/
