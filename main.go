package main

import (
	"chose-course/common/logger"
	"chose-course/common/natsclient"
	"chose-course/common/utils"
	"chose-course/config"
	"chose-course/service"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"time"
)

func initConfig() (*config.Config, error) {
	viper.SetConfigFile("./config/config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	var cfg config.Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func main() {
	log := logger.InitLogger()
	cfg, err := initConfig()
	utils.Must(err)
	db := initDB(cfg)
	redisClient := initRedis(cfg)
	natsClient := natsclient.NewNatsClient("Course", cfg.NATS.URL, log)
	server := service.InitServer(log, db, redisClient, natsClient, cfg.Server.HttpListen)
	server.Start()
	utils.WaitClose(log, func() {
		server.Stop()
		_ = log.Sync()
	})
}

func initDB(cfg *config.Config) *gorm.DB {
	// 构建 DSN 连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info), // 可以选择调试级别
	})
	if err != nil {
		utils.Must(err)
	}

	// 设置连接池参数
	sqlDB, err := db.DB()
	if err != nil {
		utils.Must(err)
	}
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	return db
}

func initRedis(cfg *config.Config) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{Addr: cfg.Redis.Addr,
		Password: cfg.Redis.Password, DB: cfg.Redis.DB})
	// 使用超时上下文，验证redis
	timeoutCtx, cancelFunc := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancelFunc()
	_, err := redisClient.Ping(timeoutCtx).Result()
	if err != nil {
		utils.Must(err)
	}
	return redisClient
}
