package model

import (
	"bucket_file/constant"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserForm struct {
	ID       string `json:"id" bson:"id"`
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
	Fullname string `json:"fullname" bson:"fullname"`
	IsAdmin  bool   `json:"is_admin" bson:"is_admin"`
}

type UserResponse struct {
	Username string `json:"username" bson:"username"`
	Fullname string `json:"fullname" bson:"fullname"`
}

type UpdateUserForm struct {
	Fullname string `json:"fullname" bson:"fullname"`
}

type ChangePasswordForm struct {
	OldPassword   string `json:"old_password"`
	NewPassword   string `json:"new_password"`
	RenewPassword string `json:"renew_password"`
}

type userModel struct {
	*mgo.Collection
}

func (db *DB) User() *userModel {
	return &userModel{db.C(constant.MONGO_COLLECTION_USER)}
}

func (db *userModel) CreateAdmin() {
	_, found, _ := db.FindByUsername("admin")
	if !found {
		var newUser UserForm
		hash, _ := bcrypt.GenerateFromPassword([]byte("admin"), 10)
		newUser.Password = string(hash)
		newUser.Username = "admin"
		newUser.Fullname = "admin"
		newUser.IsAdmin = true
		newUser.ID = genUUID()
		db.Insert(&newUser)
	}
}

func (db *userModel) FindByUsername(username string) (UserForm, bool, error) {
	return db.Get(bson.M{"username": username})
}

func (db *userModel) FindById(id string) (UserForm, bool, error) {
	return db.Get(bson.M{"id": id})
}

func (db *userModel) Get(condition bson.M) (UserForm, bool, error) {
	document := UserForm{}
	found := true
	err := db.Find(condition).One(&document)
	if err != nil {
		found = false
		if err == mgo.ErrNotFound {
			err = nil
		}
	}
	return document, found, err
}

func (db *userModel) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (db *userModel) ChangePassword(userId string, newPass string) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(newPass), 10)
	change := bson.M{"$set": bson.M{"password": hash}}
	query := bson.M{"id": userId}
	return db.Update(query, change)
}

func (db *userModel) Create(newUser UserForm) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), 10)
	newUser.Password = string(hash)
	newUser.ID = genUUID()
	newUser.IsAdmin = false
	err := db.Insert(&newUser)
	return err
}

func (db *userModel) Authen(username, password string) (UserForm, bool, error) {
	u, found, err := db.FindByUsername(username)
	if err != nil || !found {
		return u, found, err
	}

	return u, db.CheckPassword(password, u.Password), nil
}

func (db *userModel) GetListUsers(limit int, page int, search string) ([]UserResponse, int, error) {
	var results []UserResponse
	query := db.Find(bson.M{
		"username": bson.RegEx{
			Pattern: search,
			Options: "i",
		},
	})
	err := query.Limit(limit).Skip((page - 1) * limit).All(&results)

	totals, _ := query.Count()
	return results, totals, err
}

func (db *userModel) UpdateUserById(id string, updateData UpdateUserForm) error {
	change := bson.M{"$set": bson.M{"fullname": updateData.Fullname}}
	query := bson.M{"id": id}
	return db.Update(query, change)
}
