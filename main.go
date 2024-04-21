package main

import (
	"Ronin/component/appctx"
	"Ronin/middleware"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

func main() {

	dbCon, errDb := connectDBWithRetry(1)
	if errDb != nil {
		log.Println("Err connect redis")
	}

	rds, errRds := connectRedis()
	if errRds != nil {
		log.Println("Err connect redis")
	}

	appCtx := appctx.NewAppContext(dbCon, "nil", rds)

	route := gin.Default()

	route.Static("static", "static")

	route.Use(middleware.Recover(appCtx))

	v1 := route.Group("/booking-service/api/v1")
	setupAuthRoute(appCtx, v1)
	setupBookingRoute(appCtx, v1)

	route.Run()
}

func connectDBWithRetry(times int) (*gorm.DB, error) {
	var e error

	for i := 1; i <= times; i++ {
		dsn := "root:123456@tcp(localhost:3306)/db_food?charset=utf8mb4&parseTime=True&loc=Local"
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

		if err == nil {
			fmt.Println("============>>>>> MYSQL CONNECTED <<<<<===============")
			return db, nil
		}

		e = err
		time.Sleep(time.Second * 2)
	}

	return nil, e
}

func connectRedis() (*redis.Client, error) {
	strConn := fmt.Sprintf("%s:%s", "localhost", "6379")
	rdb := redis.NewClient(&redis.Options{
		Addr:     strConn,
		Password: os.Getenv("REDIS_PASSWORD"), // no password set
		DB:       0,                           // use default DB
	})

	if rdb == nil {
		log.Fatal("error connect redis")
		return nil, nil
	}

	fmt.Println("============>>>>> REDIS CONNECTED <<<<<===============")
	return rdb, nil
}
