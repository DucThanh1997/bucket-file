package model

import (
	"bucket_file/constant"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"gopkg.in/mgo.v2"
)

type jsonOwner struct {
	ID          string
	DisplayName string
}

type jsonContent struct {
	Key          string
	LastModified string
	Etag         string
	Size         int
	StorageClass string
	Owner        jsonOwner
}

type jsonListBucket struct {
	Name        string
	Prefix      string
	Marker      string
	MaxKeys     string
	IsTruncated string
	Contents    []jsonContent
}

type Owner struct {
	XMLName     xml.Name `xml:"Owner"`
	ID          string   `xml:"ID"`
	DisplayName string   `xml:"DisplayName"`
}

type Contents struct {
	XMLName      xml.Name `xml:"Contents"`
	Key          string   `xml:"Key"`
	LastModified string   `xml:"LastModified"`
	ETag         string   `xml:"ETag"`
	Size         int      `xml:"Size"`
	StorageClass string   `xml:"StorageClass"`
	Owner        Owner    `xml:"Owner"`
}

type ListBucketResult struct {
	XMLName     xml.Name   `xml:"ListBucketResult"`
	Name        string     `xml:"Name"`
	Prefix      string     `xml:"Prefix"`
	Marker      string     `xml:"Marker"`
	MaxKeys     string     `xml:"MaxKeys"`
	IsTruncated string     `xml:"IsTruncated"`
	Contents    []Contents `xml:"Contents"`
}

type ErrorResponse struct {
	XMLName   xml.Name `xml:"Error"`
	Code      string   `xml:"Code"`
	RequestId string   `xml:"RequestId"`
	HostId    string   `xml:"HostId"`
}

type JsonErrorResponse struct {
	Code      string
	RequestId string
	HostId    string
}

type bucketModel struct {
	*mgo.Collection
}

func (db *DB) Bucket() *bucketModel {
	return &bucketModel{db.C(constant.MONGO_COLLECTION_BUCKET)}
}

func (db *bucketModel) BucketRequest(httpVerb string, bucketName string) (string, int) {
	var access_key = constant.ACCESS_KEY
	var secret_key = constant.SECRET_KEY

	now := time.Now().In(time.UTC)
	date := now.Format(constant.DAY_DD_MM_YYYY_HH_MM_SS)
	var HTTP_Verb = httpVerb
	var Content_MD5 = ""
	var Content_Length = ""
	var CanonicalizedAmzHeaders = ""
	var CanonicalizedResource = fmt.Sprintf("/%s", bucketName)
	// var CanonicalizedResource = "/hunghung123"

	var StringToSign = fmt.Sprintf("%s\n%s\n%s\n%s\n%s%s", HTTP_Verb, Content_MD5, Content_Length, date, CanonicalizedAmzHeaders, CanonicalizedResource)
	h := hmac.New(sha1.New, []byte(secret_key))
	h.Write([]byte(StringToSign))
	signature := base64.URLEncoding.EncodeToString(h.Sum(nil))
	authorization_header := fmt.Sprintf("AWS %s:%s", access_key, signature)

	url := fmt.Sprintf(constant.HOST+"%s", CanonicalizedResource)

	req, err := http.NewRequest(HTTP_Verb, url, nil)
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
	return string(body), res.StatusCode
}

func (db *bucketModel) ErrorXml2Json(xmlData string) JsonErrorResponse {
	var rawError ErrorResponse
	err := xml.Unmarshal([]byte(xmlData), &rawError)
	fmt.Println(err)
	var errorData JsonErrorResponse
	errorData.Code = rawError.Code
	errorData.RequestId = rawError.RequestId
	errorData.HostId = rawError.HostId

	return errorData
}

func (db *bucketModel) Xml2Json(xmlData string) (jsonListBucket, error) {
	var bucket ListBucketResult
	var bucketData jsonListBucket
	err := xml.Unmarshal([]byte(xmlData), &bucket)
	if err != nil {
		fmt.Println(err)
		return bucketData, err
	}

	bucketData.Name = bucket.Name
	bucketData.Prefix = bucket.Prefix
	bucketData.Marker = bucket.Marker
	bucketData.MaxKeys = bucket.MaxKeys
	bucketData.IsTruncated = bucket.IsTruncated

	var contentData jsonContent
	var contentDatas []jsonContent

	for _, value := range bucket.Contents {
		contentData.Key = value.Key
		contentData.LastModified = value.LastModified
		contentData.Etag = value.ETag
		contentData.Size = value.Size
		contentData.StorageClass = value.StorageClass
		contentData.Owner.ID = value.Owner.ID
		contentData.Owner.DisplayName = value.Owner.DisplayName

		contentDatas = append(contentDatas, contentData)
	}

	bucketData.Contents = contentDatas

	return bucketData, nil
}
