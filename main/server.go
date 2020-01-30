package main

import "log"

func main() {
	log.Println("Start!")

	status, err := GetECRStatus()
	if err.p {
		log.Fatalln(err)
	}

	log.Println(status)
}
