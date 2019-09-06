package controller

import (
	"bucket_file/constant"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Code      string `json:"Code"`
	RequestId string `json:"RequestId"`
	HostId    string `json:"HostId"`
}

func CreateBucket(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	user, found, _ := container.MongoClient.User().FindById(userID)
	if found == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Can't find this User",
		})
		return
	}
	res, code := container.MongoClient.Bucket().BucketRequest("PUT", user.Username)
	log.Println(res)
	if code == 200 {
		c.JSON(http.StatusOK, gin.H{
			"msg": "Create bucket",
		})
	} else if code == 500 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": constant.SERVER_ERROR_500,
		})
	} else {
		errorResponse := container.MongoClient.Bucket().ErrorXml2Json(res)
		c.JSON(code, gin.H{
			"msg":    "error",
			"result": errorResponse,
		})
	}
}

func GetBucket(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	user, found, _ := container.MongoClient.User().FindById(userID)
	if found == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "Can't find this User",
		})
		return
	}
	res, code := container.MongoClient.Bucket().BucketRequest("GET", user.Username)
	if code == 200 {
		bucketResponse, _ := container.MongoClient.Bucket().Xml2Json(res)
		c.JSON(http.StatusOK, gin.H{
			"msg":     "Get bucket",
			"results": bucketResponse,
		})
		return
	} else {
		errorResponse := container.MongoClient.Bucket().ErrorXml2Json(res)
		c.JSON(code, gin.H{
			"msg":    "error",
			"result": errorResponse,
		})
	}
}
