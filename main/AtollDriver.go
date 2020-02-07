package main

import (
	"github.com/tarm/serial"
	"log"
)

// v3 constants
const STX_v3 = 0xFE
const ESC_v3 = 0xFD
const TSTX_v3 = 0xEE
const TESC_v3 = 0xED

type command struct {
	id byte
}

type Atoll struct {
	CurrentId byte
	query     []command
	Port      *serial.Port
}

type atollError struct {
	error string
}

func (e *atollError) Error() string {
	return e.error
}

func NewAtollV3() (*Atoll, *atollError) {
	//c := &serial.Config{Name: "COM3", Baud: 115200, ReadTimeout: time.Second * 5}
	//port, e := serial.OpenPort(c)
	//if e != nil{
	//	return nil, &atollError{"Unable to open port"}
	//}

	//atoll := Atoll{CurrentId:0, query: []command{}, Port:port}
	atoll := Atoll{CurrentId: 0, query: []command{}, Port: nil}
	//
	//go atoll.readPort()

	return &atoll, nil
}

func (a *Atoll) readPort() {
	inf := true
	buffer := make([]byte, 1024)
	for inf {

		n, err := a.Port.Read(buffer)
		if err != nil {
			log.Println("ReadPort", err)
		} else if n > 0 {
			result := buffer[0:n]
			log.Println("ReadPort", result)
		} else {
			log.Println("ReadPort - Zero Result")
		}
	}
}

func (a *Atoll) Send(data *[]byte) ([]byte, *errorCom) {

	command := command{id: a.CurrentId}
	a.CurrentId++
	if a.CurrentId >= 0xDF {
		a.CurrentId = 0
	}

	lenData := len(*data)

	len0 := byte(lenData & 0x7f)
	len1 := byte(lenData >> 7)

	crcData := append([]byte{command.id}, (*data)...)

	commandData := append([]byte{STX_v3, len0, len1}, crcData...)
	commandData = append(commandData, crc8(&crcData))

	stuffedData := []byte{}

	for n, d := range commandData {
		if n > 3 && d == STX_v3 {
			stuffedData = append(stuffedData, ESC_v3, TSTX_v3)
		} else if d == ESC_v3 {
			stuffedData = append(stuffedData, ESC_v3, TESC_v3)
		} else {
			stuffedData = append(stuffedData, d)
		}
	}

	log.Println("Send", commandData)
	log.Println("Send", stuffedData)

	n, err := a.Port.Write(stuffedData)

	if err != nil {
		log.Println("Send", err)
		return nil, &errorCom{true, 1}
	}
	log.Println("Send", "Sent bytes", n)
	return []byte{}, nil
}

const CRC8INIT = 0xFF
const CRC8POLY = 0x31 // = X^8+X^5+X^4+X^0

func crc8(data *[]byte) byte {
	//var uint8_t crc, i;
	var crc byte = CRC8INIT
	size := len(*data)
	dataId := 0
	for size >= 0 {

		crc ^= (*data)[dataId]

		for i := 0; i < 8; i++ {
			if crc&0x80 != 0 {
				crc = (crc << 1) ^ CRC8POLY
			} else {
				crc <<= 1
			}
		}
		size--
	}

	return crc
}
