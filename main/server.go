package main

import "log"

var sessionLog Logger

func main() {

	sessionLog.NewSessionLogger()
	log.Println("Start!")

	status, err := GetECRStatus()
	if err.p {
		log.Fatalln(err)
	}

	log.Println(status)
}
