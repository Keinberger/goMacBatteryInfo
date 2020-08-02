package main

import "log"

// panicError panics if there is an error
func panicError(err error) {
	if err != nil {
		panic(err)
	}
}

// logError logs the error with an error message
func logError(msg string, err error) {
	if err != nil {
		log.Println(msg, err)
	}
}
