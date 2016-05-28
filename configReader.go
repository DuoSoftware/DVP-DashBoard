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
var redisIp string
var redisPort string
var redisDb int
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
		defconfiguration.RedisIp = "127.0.0.1"
		defconfiguration.RedisPort = "6389"
		defconfiguration.RedisDb = 8
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
		redisIp = "127.0.0.1"
		redisPort = "6389"
		redisDb = 8
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
	} else {
		redisIp = fmt.Sprintf("%s:%s", defconfiguration.RedisIp, defconfiguration.RedisPort)
		redisPort = defconfiguration.RedisPort
		redisDb = defconfiguration.RedisDb
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
		redisIp = os.Getenv(envconfiguration.RedisIp)
		redisPort = os.Getenv(envconfiguration.RedisPort)
		redisDb, converr = strconv.Atoi(os.Getenv(envconfiguration.RedisDb))
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

		if redisIp == "" {
			redisIp = defConfig.RedisIp
		}
		if redisPort == "" {
			redisPort = defConfig.RedisPort
		}
		if redisDb == 0 || converr != nil {
			redisDb = defConfig.RedisDb
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

		redisIp = fmt.Sprintf("%s:%s", redisIp, redisPort)
		securityIp = fmt.Sprintf("%s:%s", securityIp, securityPort)
	}

	fmt.Println("redisIp:", redisIp)
	fmt.Println("redisDb:", redisDb)
	fmt.Println("redisDb:", securityIp)
}
