package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

type Food struct {
	Id           int64   `json:"id" gorm:"column:id"`
	RestaurantId int64   `json:"restaurant_id" gorm:"column:restaurant_id"`
	CategoryId   int64   `json:"category_id" gorm:"column:category_id"`
	Name         string  `json:"name" gorm:"column:name"`
	Description  string  `json:"description" gorm:"column:description"`
	Price        float64 `json:"price" gorm:"column:price"`
}

const KEY_REDIS = "ronin_foods"

func main() {

	dbCon, errDb := connectDBWithRetry(1)
	if errDb != nil {
		log.Println("Err connect redis")
	}

	rds, err := connectRedis()
	if err != nil {
		log.Println("Err connect redis")
	}
	var foods []Food

	// Get data from redis first
	errGetFromRds := GetDatRedis(context.Background(), KEY_REDIS, &foods, rds)
	if errGetFromRds != nil {
		log.Println("Error get data from redis")
	}

	if len(foods) > 0 {
		fmt.Println("=======>>>>>>> DATA get from cache")
		// Response data
		return
	}

	errGet := dbCon.Model(Food{}).Find(&foods).Error
	if errGet != nil {
		fmt.Println("Err get list")
	}

	// caching data to redis
	errRds := SetDataToRedis(context.Background(), foods, KEY_REDIS, time.Hour*24, rds)
	if errRds != nil {
		log.Fatal("Err set data to redis")
	}
	fmt.Println("=======>>>>>>>>>> Data get from database")
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
