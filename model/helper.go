package model

import "log"

func checkFatalError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
