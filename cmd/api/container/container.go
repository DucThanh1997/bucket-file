package container

import (
	"bucket_file/constant"
	"bucket_file/model"
	"log"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/natefinch/lumberjack"
	"gopkg.in/mgo.v2"
)

type Config struct {
	Listen string

	//token
	TokenSecretKey   string
	TokenExpiredHour uint64

	//log
	ErrorLog string

	//mongo
	MongoServer   []string
	MongoUser     string
	MongoPassword string
	MongoAuthDB   string
	MongoDatabase string

	// Host
	Host string
}

type Container struct {
	Config      *Config
	MongoClient *model.DB
	Log         *log.Logger
}

func NewContainer() *Container {
	var container = new(Container)
	return container
}

func (container *Container) Setup(pathConfig string) error {
	if _, err := toml.DecodeFile(pathConfig, &container.Config); err != nil {
		return err
	}
	container.loadLog()

	err := container.loadMongo()
	if err != nil {
		return err
	}

	return nil
}

func (container *Container) loadLog() {
	container.Log = log.New(nil, "", log.LstdFlags|log.Lshortfile)
	if len(container.Config.ErrorLog) < 1 {
		container.Config.ErrorLog = constant.DEFAULT_PATH_ERROR_LOG_API
	}

	container.Log.SetOutput(&lumberjack.Logger{
		Filename:   container.Config.ErrorLog,
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	})

}

func (container *Container) loadMongo() error {
	info := &mgo.DialInfo{
		Addrs:    container.Config.MongoServer,
		Timeout:  30 * time.Second,
		Username: container.Config.MongoUser,
		Password: container.Config.MongoPassword,
		Database: container.Config.MongoAuthDB,
	}
	var err error

	sesstion, err := mgo.DialWithInfo(info)
	if err != nil {
		return err
	}

	mongoClient := sesstion.DB(container.Config.MongoDatabase)
	container.MongoClient = &model.DB{mongoClient}
	container.MongoClient.Session.SetSocketTimeout(30 * time.Second)
	return nil
}
