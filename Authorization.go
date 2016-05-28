package main

import (
	"encoding/json"
	"fmt"
	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/fzzy/radix/redis"
	"github.com/gorilla/context"
	"strconv"
	"strings"
	"time"
)

func loadJwtMiddleware() *jwtmiddleware.JWTMiddleware {
	return (jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			fmt.Println(token.Claims["iss"])
			fmt.Println(token.Claims["jti"])
			secretKey := fmt.Sprintf("token:iss:%s:%s", token.Claims["iss"], token.Claims["jti"])
			secret := SecurityGet(secretKey)
			fmt.Println("secret: ", secret)
			if secret == "" {
				return nil, fmt.Errorf("Invalied 'iss' or 'jti' in JWT")
			}
			return []byte(secret), nil
		},
		// When set, the middleware verifies that tokens are signed with the specific signing algorithm
		// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
		// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		SigningMethod: jwt.SigningMethodHS256,
	}))
}

func validateCompanyTenant(dashboardEvent DashBoardEvent) (company, tenant int) {
	internalAccessToken := dashboardEvent.Context.Request().Header.Get("companyinfo")
	if internalAccessToken != "" {
		ids := strings.Split(internalAccessToken, ":")
		if len(ids) == 2 {
			tenant, _ = strconv.Atoi(ids[0])
			company, _ = strconv.Atoi(ids[1])
			return company, tenant
		} else {
			return 0, 0
		}
	} else {
		user := context.Get(dashboardEvent.Context.Request(), "user")
		if user != nil {
			iTenant := user.(*jwt.Token).Claims["tenant"]
			iCompany := user.(*jwt.Token).Claims["company"]
			if iTenant != nil && iCompany != nil {
				tenant := int(iTenant.(float64))
				company := int(iCompany.(float64))
				return company, tenant
			} else {
				dashboardEvent.RB().Write(ResponseGenerator(false, "Invalid company or tenant", "", ""))
				return
			}
		} else {
			dashboardEvent.RB().Write(ResponseGenerator(false, "User data not found in JWT", "", ""))
			return
		}
	}
}

func validateCompanyTenantGraph(dashBoardGraph DashBoardGraph) (company, tenant int) {
	internalAccessToken := dashBoardGraph.Context.Request().Header.Get("companyinfo")
	if internalAccessToken != "" {
		ids := strings.Split(internalAccessToken, ":")
		if len(ids) == 2 {
			tenant, _ = strconv.Atoi(ids[0])
			company, _ = strconv.Atoi(ids[1])
			return company, tenant
		} else {
			return 0, 0
		}
	} else {
		fmt.Println(dashBoardGraph.Context.Request())
		user := context.Get(dashBoardGraph.Context.Request(), "user")
		fmt.Println(user)
		if user != nil {
			iTenant := user.(*jwt.Token).Claims["tenant"]
			iCompany := user.(*jwt.Token).Claims["company"]
			if iTenant != nil && iCompany != nil {
				tenant := int(iTenant.(float64))
				company := int(iCompany.(float64))
				return company, tenant
			} else {
				dashBoardGraph.RB().Write(ResponseGenerator(false, "Invalid company or tenant", "", ""))
				return
			}
		} else {
			dashBoardGraph.RB().Write(ResponseGenerator(false, "User data not found in JWT", "", ""))
			return
		}
	}
}

func ResponseGenerator(isSuccess bool, customMessage, result, exception string) []byte {
	res := Result{}
	res.IsSuccess = isSuccess
	res.CustomMessage = customMessage
	res.Exception = exception
	res.Result = result
	resb, _ := json.Marshal(res)
	return resb
}

func SecurityGet(key string) string {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in RedisGet", r)
		}
	}()
	client, err := redis.DialTimeout("tcp", securityIp, time.Duration(10)*time.Second)
	errHndlr(err)
	defer client.Close()

	//authServer
	authE := client.Cmd("auth", redisPassword)
	errHndlr(authE.Err)

	strObj, _ := client.Cmd("get", key).Str()
	fmt.Println(strObj)
	return strObj
}
