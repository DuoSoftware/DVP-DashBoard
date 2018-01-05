package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

var dirPath string
var redisPubSubIp string
var redisPubSubPort string
var redisPubSubPassword string
var redisIp string
var redisPort string
var redisDb int
var ardsRedisDb int
var redisPassword string
var port string
var statsDIp string
var statsDPort int
var pgUser string
var pgPassword string
var pgDbname string
var pgHost string
var pgPort int
var securityIp string
var securityPort string
var securityPassword string
var mongoIp string
var mongoPort string
var mongoDbname string
var mongoPassword string
var mongoUser string
var cacheMachenism string
var dashboardServiceHost string
var dashboardServicePort string
var accessToken string
var redisClusterName string
var redisMode string
var sentinelHosts string
var sentinelPort string
var decrRetryCount string
var decrRetryDelay string
var rabbitMQIp string
var rabbitMQPort string
var rabbitMQUser string
var rabbitMQPassword string
var useMsgQueue string

func GetDirPath() string {
	envPath := os.Getenv("GO_CONFIG_DIR")
	if envPath == "" {
		envPath = "./"
	}
	fmt.Println(envPath)
	return envPath
}

func GetDefaultConfig() Configuration {
	confPath := filepath.Join(dirPath, "conf.json")
	fmt.Println("GetDefaultConfig config path: ", confPath)
	content, operr := ioutil.ReadFile(confPath)
	if operr != nil {
		fmt.Println(operr)
	}

	defconfiguration := Configuration{}
	deferr := json.Unmarshal(content, &defconfiguration)

	if deferr != nil {
		fmt.Println("error:", deferr)
		panic(deferr)
	}

	return defconfiguration
}

func LoadDefaultConfig() {

	defconfiguration := GetDefaultConfig()

	redisPubSubIp = fmt.Sprintf("%s:%s", defconfiguration.RedisPubSubIp, defconfiguration.RedisPort)
	redisIp = fmt.Sprintf("%s:%s", defconfiguration.RedisIp, defconfiguration.RedisPort)
	redisPubSubPort = defconfiguration.RedisPubSubPort
	redisPubSubPassword = defconfiguration.RedisPubSubPassword
	redisPort = defconfiguration.RedisPort
	redisDb = defconfiguration.RedisDb
	ardsRedisDb = defconfiguration.ArdsRedisDb
	redisPassword = defconfiguration.RedisPassword
	port = defconfiguration.Port
	statsDIp = defconfiguration.StatsDIp
	statsDPort = defconfiguration.StatsDPort
	pgUser = defconfiguration.PgUser
	pgPassword = defconfiguration.PgPassword
	pgDbname = defconfiguration.PgDbname
	pgHost = defconfiguration.PgHost
	pgPort = defconfiguration.PgPort
	securityIp = fmt.Sprintf("%s:%s", defconfiguration.SecurityIp, defconfiguration.SecurityPort)
	securityPort = defconfiguration.SecurityPort
	securityPassword = defconfiguration.SecurityPassword
	mongoIp = defconfiguration.MongoIp
	mongoPort = defconfiguration.MongoPort
	mongoDbname = defconfiguration.MongoDbname
	mongoPassword = defconfiguration.MongoPassword
	mongoUser = defconfiguration.MongoUser
	cacheMachenism = defconfiguration.CacheMachenism
	dashboardServiceHost = defconfiguration.DashboardServiceHost
	dashboardServicePort = defconfiguration.DashboardServicePort
	accessToken = defconfiguration.AccessToken
	redisClusterName = defconfiguration.RedisClusterName
	redisMode = defconfiguration.RedisMode
	sentinelHosts = defconfiguration.SentinelHosts
	sentinelPort = defconfiguration.SentinelPort
	decrRetryCount = defconfiguration.DecrRetryCount
	decrRetryDelay = defconfiguration.DecrRetryDelay
	rabbitMQIp = defconfiguration.RabbitMQIp
	rabbitMQPort = defconfiguration.RabbitMQPort
	rabbitMQUser = defconfiguration.RabbitMQUser
	rabbitMQPassword = defconfiguration.RabbitMQPassword
	useMsgQueue = defconfiguration.UseMsgQueue
}

func LoadConfiguration() {
	dirPath = GetDirPath()
	confPath := filepath.Join(dirPath, "custom-environment-variables.json")
	fmt.Println("InitiateRedis config path: ", confPath)

	content, operr := ioutil.ReadFile(confPath)
	if operr != nil {
		fmt.Println(operr)
	}

	envconfiguration := EnvConfiguration{}
	enverr := json.Unmarshal(content, &envconfiguration)
	if enverr != nil {
		fmt.Println("error:", enverr)
		LoadDefaultConfig()
	} else {
		var converr error
		defConfig := GetDefaultConfig()
		redisPubSubIp = os.Getenv(envconfiguration.RedisPubSubIp)
		redisPubSubPort = os.Getenv(envconfiguration.RedisPubSubPort)
		redisPubSubPassword = os.Getenv(envconfiguration.RedisPubSubPassword)
		redisIp = os.Getenv(envconfiguration.RedisIp)
		redisPort = os.Getenv(envconfiguration.RedisPort)
		redisDb, converr = strconv.Atoi(os.Getenv(envconfiguration.RedisDb))
		ardsRedisDb, converr = strconv.Atoi(os.Getenv(envconfiguration.ArdsRedisDb))
		redisPassword = os.Getenv(envconfiguration.RedisPassword)
		port = os.Getenv(envconfiguration.Port)
		statsDIp = os.Getenv(envconfiguration.StatsDIp)
		statsDPort, converr = strconv.Atoi(os.Getenv(envconfiguration.StatsDPort))
		pgUser = os.Getenv(envconfiguration.PgUser)
		pgPassword = os.Getenv(envconfiguration.PgPassword)
		pgDbname = os.Getenv(envconfiguration.PgDbname)
		pgHost = os.Getenv(envconfiguration.PgHost)
		pgPort, converr = strconv.Atoi(os.Getenv(envconfiguration.PgPort))
		securityIp = os.Getenv(envconfiguration.SecurityIp)
		securityPort = os.Getenv(envconfiguration.SecurityPort)
		securityPassword = os.Getenv(envconfiguration.SecurityPassword)
		mongoIp = os.Getenv(envconfiguration.MongoIp)
		mongoPort = os.Getenv(envconfiguration.MongoPort)
		mongoDbname = os.Getenv(envconfiguration.MongoDbname)
		mongoPassword = os.Getenv(envconfiguration.MongoPassword)
		mongoUser = os.Getenv(envconfiguration.MongoUser)
		cacheMachenism = os.Getenv(envconfiguration.CacheMachenism)
		dashboardServiceHost = os.Getenv(envconfiguration.DashboardServiceHost)
		dashboardServicePort = os.Getenv(envconfiguration.DashboardServicePort)
		accessToken = os.Getenv(envconfiguration.AccessToken)
		redisClusterName = os.Getenv(envconfiguration.RedisClusterName)
		redisMode = os.Getenv(envconfiguration.RedisMode)
		sentinelHosts = os.Getenv(envconfiguration.SentinelHosts)
		sentinelPort = os.Getenv(envconfiguration.SentinelPort)
		decrRetryCount = os.Getenv(envconfiguration.DecrRetryCount)
		decrRetryDelay = os.Getenv(envconfiguration.DecrRetryDelay)
		rabbitMQIp = os.Getenv(envconfiguration.RabbitMQIp)
		rabbitMQPort = os.Getenv(envconfiguration.RabbitMQPort)
		rabbitMQUser = os.Getenv(envconfiguration.RabbitMQUser)
		rabbitMQPassword = os.Getenv(envconfiguration.RabbitMQPassword)
		useMsgQueue = os.Getenv(envconfiguration.UseMsgQueue)

		if redisPubSubIp == "" {
			redisPubSubIp = defConfig.RedisPubSubIp
		}
		if redisPubSubPort == "" {
			redisPubSubPort = defConfig.RedisPubSubPort
		}
		if redisPubSubPassword == "" {
			redisPubSubPassword = defConfig.RedisPubSubPassword
		}
		if redisIp == "" {
			redisIp = defConfig.RedisIp
		}
		if redisPort == "" {
			redisPort = defConfig.RedisPort
		}
		if converr != nil {
			redisDb = defConfig.RedisDb
		}
		if converr != nil {
			ardsRedisDb = defConfig.ArdsRedisDb
		}
		if redisPassword == "" {
			redisPassword = defConfig.RedisPassword
		}
		if port == "" {
			port = defConfig.Port
		}
		if statsDIp == "" {
			statsDIp = defConfig.StatsDIp
		}
		if statsDPort == 0 || converr != nil {
			statsDPort = defConfig.StatsDPort
		}
		if pgUser == "" {
			pgUser = defConfig.PgUser
		}
		if pgPassword == "" {
			pgPassword = defConfig.PgPassword
		}
		if pgDbname == "" {
			pgDbname = defConfig.PgDbname
		}
		if pgHost == "" {
			pgHost = defConfig.PgHost
		}
		if pgPort == 0 || converr != nil {
			pgPort = defConfig.PgPort
		}
		if securityIp == "" {
			securityIp = defConfig.SecurityIp
		}
		if securityPort == "" {
			securityPort = defConfig.SecurityPort
		}
		if securityPassword == "" {
			securityPassword = defConfig.SecurityPassword
		}
		if mongoIp == "" {
			mongoIp = defConfig.MongoIp
		}
		if mongoPort == "" {
			mongoPort = defConfig.MongoPort
		}
		if mongoDbname == "" {
			mongoDbname = defConfig.MongoDbname
		}
		if mongoUser == "" {
			mongoUser = defConfig.MongoUser
		}
		if mongoPassword == "" {
			mongoPassword = defConfig.MongoPassword
		}
		if cacheMachenism == "" {
			cacheMachenism = defConfig.CacheMachenism
		}
		if dashboardServiceHost == "" {
			dashboardServiceHost = defConfig.DashboardServiceHost
		}
		if dashboardServicePort == "" {
			dashboardServicePort = defConfig.DashboardServicePort
		}
		if accessToken == "" {
			accessToken = defConfig.AccessToken
		}
		if redisClusterName == "" {
			redisClusterName = defConfig.RedisClusterName
		}

		if redisMode == "" {
			redisMode = defConfig.RedisMode
		}
		if sentinelHosts == "" {
			sentinelHosts = defConfig.SentinelHosts
		}
		if sentinelPort == "" {
			sentinelPort = defConfig.SentinelPort
		}
		if decrRetryCount == "" {
			decrRetryCount = defConfig.DecrRetryCount
		}
		if decrRetryDelay == "" {
			decrRetryDelay = defConfig.DecrRetryDelay
		}
		if rabbitMQIp == "" {
			rabbitMQIp = defConfig.RabbitMQIp
		}
		if rabbitMQPort == "" {
			rabbitMQPort = defConfig.RabbitMQPort
		}
		if rabbitMQUser == "" {
			rabbitMQUser = defConfig.RabbitMQUser
		}
		if rabbitMQPassword == "" {
			rabbitMQPassword = defConfig.RabbitMQPassword
		}
		if useMsgQueue == "" {
			useMsgQueue = defConfig.UseMsgQueue
		}

		redisIp = fmt.Sprintf("%s:%s", redisIp, redisPort)
		redisPubSubIp = fmt.Sprintf("%s:%s", redisPubSubIp, redisPubSubPort)
		securityIp = fmt.Sprintf("%s:%s", securityIp, securityPort)
	}

	fmt.Println("redisMode:", redisMode)
	fmt.Println("sentinelHosts:", sentinelHosts)
	fmt.Println("sentinelPort:", sentinelPort)
	fmt.Println("redisPubSubIp:", redisPubSubIp)
	fmt.Println("redisIp:", redisIp)
	fmt.Println("redisDb:", redisDb)
	fmt.Println("securityIp:", securityIp)
	fmt.Println("decrRetryCount:", decrRetryCount)
	fmt.Println("decrRetryDelay:", decrRetryDelay)
	fmt.Println("rabbitMQIp:", rabbitMQIp)
	fmt.Println("rabbitMQPort:", rabbitMQPort)
	fmt.Println("rabbitMQUser:", rabbitMQUser)
	fmt.Println("rabbitMQPassword:", rabbitMQPassword)
	fmt.Println("useMsgQueue:", useMsgQueue)
}
