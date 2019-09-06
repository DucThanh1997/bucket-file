package model

import (
	"bucket_file/constant"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"gopkg.in/mgo.v2"
	_ "bytes"
    // "log"
)

type FileForm struct {
	ID       string `json:"id" bson:"id"`
	FileName string `json`
}

type fileModel struct {
	*mgo.Collection
}

func (db *DB) File() *fileModel {
	return &fileModel{db.C(constant.MONGO_COLLECTION_FILE)}
}

func (db *fileModel) FileRequest(httpVerb string, bucketName string, fileName string, pathfile string) (string, int) {
	var access_key = constant.ACCESS_KEY
	var secret_key = constant.SECRET_KEY

	now := time.Now().In(time.UTC)
	date := now.Format(constant.DAY_DD_MM_YYYY_HH_MM_SS)
	var HTTP_Verb = httpVerb
	var Content_MD5 = ""
	var Content_Length = ""
	var CanonicalizedAmzHeaders = ""
	var CanonicalizedResource = fmt.Sprintf("/%s/%s", bucketName, fileName)

	var StringToSign = fmt.Sprintf("%s\n%s\n%s\n%s\n%s%s", httpVerb, Content_MD5, Content_Length, date, CanonicalizedAmzHeaders, CanonicalizedResource)
	h := hmac.New(sha1.New, []byte(secret_key))
	h.Write([]byte(StringToSign))
	signature := base64.URLEncoding.EncodeToString(h.Sum(nil))
	authorization_header := fmt.Sprintf("AWS %s:%s", access_key, signature)

	url := fmt.Sprintf(constant.HOST+"%s", CanonicalizedResource)
	fmt.Println(url)
	// jsonStr,err := ioutil.ReadFile(pathfile)
	// if(err!=nil){
	// 	panic(err)
	// }
	fmt.Println(HTTP_Verb)
	req, err := http.NewRequest(httpVerb, url, nil)
	if err != nil {
		fmt.Println(err)
		return "", 500
	}

	req.Header.Add("Date", date)
	req.Header.Add("Authorization", authorization_header)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", 500
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	
	fmt.Println("Cương: ", string(body))

	return string(body), res.StatusCode
}

