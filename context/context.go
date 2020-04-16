package context

import (
	"fmt"
	conf "go-webapp-starter/conf"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Context struct
type Context struct {
	DB   *DB
	RD   *RD
	HTTP *http.Client
}

// Init method
func (c *Context) Init() {
	c.DB = &DB{}
	c.DB.Connect()
	c.RD = &RD{}
	c.RD.Connect()
	c.HTTP = CreateHTTPClient(10, 20)
}

// CreateHTTPClient func
func CreateHTTPClient(timeout int, maxIdleConn int) *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: maxIdleConn,
		},
		Timeout: time.Duration(timeout) * time.Second,
	}
	return client
}

// DB Mysql client
type DB struct {
	Client *sqlx.DB
}

// Connect Mysql connect method
func (db *DB) Connect() {
	uri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&autocommit=true&parseTime=true&loc=utc",
		conf.MysqlUser, conf.MysqlPassword, conf.MysqlHost, conf.MysqlPort, conf.MysqlDatabase)
	client, err := sqlx.Connect("mysql", uri)
	if err != nil {
		panic(err)
	}
	db.Client = client
}

// RD Redis client
type RD struct {
	Client *redis.Client
}

// Connect Redis connect method
func (rd *RD) Connect() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	rd.Client = client
}
