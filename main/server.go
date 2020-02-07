package main

import "log"

var sessionLog Logger

func main() {

	sessionLog.NewSessionLogger()
	log.Println("Start!")

	atoll, err := NewAtollV3()
	if err != nil {
		log.Fatalln(err)
		return
	}

	sendData := []byte{0xc1, 5, 0x10, 0x97, 0x4c, 0x31, 0x32, 0x33}

	res, _ := atoll.Send(&sendData)

	log.Println(res)

	//status, err := GetECRStatus()
	//if err.p {
	//	log.Fatalln(err)
	//}
	//
	//log.Println(status)
}
