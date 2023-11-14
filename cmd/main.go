package main

import (
	"log"

	"github.com/zapscloud/golib-utils/utils"
)

func main() {

	testString := "tickets-dev-db-biz_cl9guv21ev366n9iciag"
	log.Println("New String: ", utils.Left(testString, 10), utils.Right(testString, 10))
}

// func Right(val string, n int) string {
// 	retVal := val

// 	l := len(val)
// 	if n < l {
// 		retVal = val[(l - n):l]
// 	}
// 	return retVal
// }

// func Left(val string, n int) string {
// 	retVal := val

// 	l := len(val)
// 	if n < l {
// 		retVal = val[0:n]
// 	}
// 	return retVal
// }
