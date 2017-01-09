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
		defconfiguration.RedisPubSubIp = "127.0.0.1"
		defconfiguration.RedisPubSubPort = "6389"
		defconfiguration.RedisPubSubPassword = "DuoS123"
		defconfiguration.RedisIp = "127.0.0.1"
		defconfiguration.RedisPort = "6389"
		defconfiguration.RedisDb = 8
		defconfiguration.ArdsRedisDb = 6
		defconfiguration.RedisPassword = "DuoS123"
		defconfiguration.Port = "2226"
		defconfiguration.StatsDIp = "45.55.142.207"
		defconfiguration.StatsDPort = 8125
		defconfiguration.PgUser = "duo"
		defconfiguration.PgPassword = "DuoS123"
		defconfiguration.PgDbname = "dvpdb"
		defconfiguration.PgHost = "104.131.105.222"
		defconfiguration.PgPort = 5432
		defconfiguration.SecurityIp = "45.55.142.207"
		defconfiguration.SecurityPort = "6389"
		defconfiguration.SecurityPassword = "DuoS123"
		defconfiguration.MongoIp = "45.55.142.207"
		defconfiguration.MongoPort = "27017"
		defconfiguration.MongoDbname = "dvpdb"
		defconfiguration.MongoPassword = "DuoS123"
		defconfiguration.MongoUser = "duo"
		defconfiguration.CacheMachenism = "redis"
	}

	return defconfiguration
}

func LoadDefaultConfig() {
	confPath := filepath.Join(dirPath, "conf.json")
	fmt.Println("LoadDefaultConfig config path: ", confPath)

	content, operr := ioutil.ReadFile(confPath)
	if operr != nil {
		fmt.Println(operr)
	}

	defconfiguration := Configuration{}
	deferr := json.Unmarshal(content, &defconfiguration)

	if deferr != nil {
		fmt.Println("error:", deferr)
		redisPubSubIp = "127.0.0.1"
		redisPubSubPort = "6389"
		redisPubSubPassword = "DuoS123"
		redisIp = "127.0.0.1"
		redisPort = "6389"
		redisDb = 8
		ardsRedisDb = 6
		redisPassword = "DuoS123"
		port = "2226"
		statsDIp = "45.55.142.207"
		statsDPort = 8125
		pgUser = "duo"
		pgPassword = "DuoS123"
		pgDbname = "dvpdb"
		pgHost = "104.131.105.222"
		pgPort = 5432
		securityIp = "45.55.142.207"
		securityPort = "6389"
		securityPassword = "DuoS123"
		mongoIp = "45.55.142.207"
		mongoPort = "27017"
		mongoDbname = "dvpdb"
		mongoPassword = "DuoS123"
		mongoUser = "duo"
		cacheMachenism = "redis"
	} else {
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
		securityIp = defconfiguration.SecurityIp
		securityPort = defconfiguration.SecurityPort
		securityPassword = defconfiguration.SecurityPassword
		mongoIp = defconfiguration.MongoIp
		mongoPort = defconfiguration.MongoPort
		mongoDbname = defconfiguration.MongoDbname
		mongoPassword = defconfiguration.MongoPassword
		mongoUser = defconfiguration.MongoUser
		cacheMachenism = defconfiguration.CacheMachenism
	}
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

		redisIp = fmt.Sprintf("%s:%s", redisIp, redisPort)
		redisPubSubIp = fmt.Sprintf("%s:%s", redisPubSubIp, redisPubSubPort)
		securityIp = fmt.Sprintf("%s:%s", securityIp, securityPort)
	}

	fmt.Println("redisPubSubIp:", redisPubSubIp)
	fmt.Println("redisIp:", redisIp)
	fmt.Println("redisDb:", redisDb)
	fmt.Println("redisDb:", securityIp)
}
