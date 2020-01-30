package main

import "log"

var opPass = []byte{0x1E, 0, 0, 0}

type ECRStatusShort struct {
	flags 		uint16
	mode		byte
	subMode		byte
	checkOperations uint16
	mainVoltage	byte
	backupVoltage	byte
	frErrCode	byte
	eklzErrCode	byte
	lastprintResult byte
}
type ECRStatus struct {
	version 	uint16
	build	 	uint16
	buildDate	string
	devNumber	byte
	curentDocNum	uint16
	flags 		uint16
	mode		byte
	sunMode		byte
	port 		byte
	versionFM	uint16
	buildFM		uint16
	buildFMDate	string
	date 		string
	time 		string
	flagsFM		uint16
	factoryNum	uint64
	periodNum	uint16
	freeFMRecs	uint16
	resignNum	byte
	resignNumAvlbl	byte
	inn 		uint64
	modeFM		byte
}


func GetECRStatus() (r []byte, err errorCom){
	sessionLog.Print("--GET-ECR_STATUS-- Старт  GetECRStatus")
	result, err := Send(append([]byte{0x11}, opPass...), 48)
	var status ECRStatus

	log.Print("Статус",status)
	sessionLog.Print("--GET-ECR_STATUS-- Стоп GetECRStatus" )

	return result, err
}

