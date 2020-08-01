package main

import "log"

func panicError(err error) {
	if err != nil {
		panic(err)
	}
}

func logError(msg string, err error) {
	if err != nil {
		log.Fatal(msg, err)
	}
}
