package controller

import (
	"bucket_file/model"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"bucket_file/constant"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Id       string `json:"id"`
	Fullname string `json:"fullname"`
	jwt.StandardClaims
}

func Signin(c *gin.Context) {
	var loginData Login
	err := c.BindJSON(&loginData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "something went wrong",
		})
		return
	}
	user, _, _ := container.MongoClient.User().FindByUsername(loginData.Username)

	expirationTime := time.Now().Add(time.Duration(container.Config.TokenExpiredHour) * time.Hour)
	claims := &Claims{
		Id:       user.ID,
		Fullname: user.Fullname,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString([]byte(container.Config.TokenSecretKey))
	log.Println(tokenString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": constant.SERVER_ERROR_500,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":   "Signed in",
		"token": tokenString,
	})
}

func RefreshToken(c *gin.Context) {
	token := c.MustGet("token").(string)
	claims := &Claims{}
	parseToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(container.Config.TokenSecretKey), nil
	})

	if !parseToken.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"msg": "Token is valid",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": constant.SERVER_ERROR_500,
		})
		return
	}

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) < 30*time.Second {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Token expired",
		})
		return
	}

	expirationTime := time.Now().Add(time.Duration(container.Config.TokenExpiredHour) * time.Hour)
	claims.ExpiresAt = expirationTime.Unix()
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := newToken.SignedString([]byte(container.Config.TokenSecretKey))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": constant.SERVER_ERROR_500,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":   "Refresh token",
		"token": tokenString,
	})
}

func CreateUser(c *gin.Context) {
	var postData model.UserForm
	err := c.BindJSON(&postData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": constant.SERVER_ERROR_500,
		})
		return
	}
	found := false
	_, found, err = container.MongoClient.User().FindByUsername(postData.Username)
	if found == true {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "User already exist",
		})
		return
	}

	container.MongoClient.User().Create(postData)
	c.JSON(http.StatusCreated, gin.H{
		"msg": "Create user success",
	})
}

func UpdateUser(c *gin.Context) {
	var postData model.UpdateUserForm
	err := c.BindJSON(&postData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "something went wrong",
		})
		return
	}

	userId := c.MustGet("userID").(string)

	var user model.UserForm
	err = container.MongoClient.User().UpdateUserById(userId, postData)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": constant.SERVER_ERROR_500,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":    "update success",
		"result": user,
	})
}

func ChangePassword(c *gin.Context) {
	var postData model.ChangePasswordForm
	err := c.BindJSON(&postData)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": constant.SERVER_ERROR_500,
		})
		return
	}

	if postData.NewPassword != postData.RenewPassword {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Mật khẩu nhập lại không khớp",
		})
		return
	}

	userID := c.MustGet("userID").(string)
	user, _, err := container.MongoClient.User().FindById(userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": constant.SERVER_ERROR_500,
		})
		return
	}

	if isMatch := container.MongoClient.User().CheckPassword(postData.OldPassword, user.Password); isMatch == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Sai mật khẩu hiện tại",
		})
		return
	}

	err = container.MongoClient.User().ChangePassword(user.ID, postData.NewPassword)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": constant.SERVER_ERROR_500,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "Đổi mật khẩu thành công",
	})
}

func GetListUsers(c *gin.Context) {
	var limit, page int
	limit, _ = strconv.Atoi(c.Query("limit"))
	if limit <= 0 {
		limit = 10
	}

	page, _ = strconv.Atoi(c.Query("page"))
	if page <= 0 {
		page = 1
	}

	search := c.Query("search")

	results, totals, err := container.MongoClient.User().GetListUsers(limit, page, search)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": constant.SERVER_ERROR_500,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":     "get list users",
		"results": results,
		"totals":  totals,
	})
}

func login(c *gin.Context) (map[string]interface{}, error) {
	var e = fmt.Errorf("missing username and password")
	var data map[string]string
	err := c.BindJSON(&data)
	if err != nil {
		return nil, e
	}

	email, ok := data["username"]
	if !ok {
		return nil, e
	}

	password, ok := data["password"]
	if !ok {
		return nil, e
	}

	e = fmt.Errorf("username or password is incorrect")
	user, found, err := container.MongoClient.User().Authen(email, password)
	if err != nil {
		container.Log.Println(err)
		return nil, e
	}

	if !found {
		return nil, e
	}

	var payload = map[string]interface{}{
		"userID":   user.ID,
		"username": user.Username,
	}

	return payload, nil
}

func userVerify(data map[string]interface{}) (bool, error) {
	id, ok := data["userID"]
	if !ok {
		return false, nil
	}

	_, found, err := container.MongoClient.User().FindById(fmt.Sprintf("%v", id))
	if err != nil {
		container.Log.Println(err)
		return false, nil
	}

	if !found {
		return false, nil
	}

	return true, nil
}
