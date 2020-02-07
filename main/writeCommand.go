package main

import (
	"fmt"
	"github.com/tarm/serial"
	"strconv"
	"time"
)

const ENC = 0x05
const STX = 0x02
const ACK = 0x06
const NAK = 0x15

const EOT = 0x04
const ETX = 0x03
const DLE = 0x10

const T1 = 500
const T2 = 2000
const T3 = 500
const T4 = 500
const T5 = 10000
const T6 = 500
const T7 = 500
const T8 = 1000

type errorCom struct {
	p bool
	c byte
}

func Send(data []byte, responseLength int) (r []byte, err errorCom) {
	sessionLog.PrintL(LOG_METOD_START_INT, "--SEND-- Старт Send")
	// вычисляем контрольную сумму
	length := byte(len(data))
	lrc := length
	var i byte
	for i = 0; i < length; i++ {
		lrc ^= data[i]
	}

	sendData := append(append([]byte{STX, length}, data...), lrc)

	//println(sendData)
	r, err = writeCommand(sendData, responseLength)
	sessionLog.PrintL(LOG_METOD_RESULT_INT, "--SEND-- Стоп Send, длина ответа: "+fmt.Sprint(responseLength)+", Ответ: "+fmt.Sprint(r))
	return r, err
}

func writeCommandV3(data []byte, responseLength int) (r []byte, err errorCom) {
	sessionLog.PrintL(LOG_METOD_START_INT, "--WRITE-COMMAND-- Старт  writeCommand")
	c := &serial.Config{Name: "COM3", Baud: 115200, ReadTimeout: time.Second * 5}
	s, e := serial.OpenPort(c)
	if e != nil {
		sessionLog.PrintL(LOG_ERROR, "--WRITE-COMMAND-- Cannot open COM port")
		sessionLog.PrintL(LOG_METOD_RESULT, "--WRITE-COMMAND-- Стоп  writeCommand")
		return nil, errorCom{false, 0}
	}
	defer s.Close()
	sessionLog.PrintL(LOG_METOD_INSIDE, "--WRITE-COMMAND-- Порт открыт")
	attempts := 10

	//  Отправляем команду
	for attempts > 0 {
		sessionLog.PrintL(LOG_METOD_INSIDE, "--WRITE-COMMAND-- Попытка "+fmt.Sprint(11-attempts))
		attempts--
		n, err := s.Write([]byte{ENC})
		if err != nil {
			sessionLog.PrintL(LOG_ERROR, "--WRITE-COMMAND-- Ошибка записи в COM порт")
			continue
		}
		sessionLog.PrintL(LOG_PORT_DATA, "--COM-PORT-SEND-- Отправлен ENC "+fmt.Sprint([]byte{ENC}))

		buf := make([]byte, 1024)
		n, err = s.Read(buf)
		//commandSent := false
		if n > 0 {

		}
		return buf, errorCom{}
	}
	return []byte{}, errorCom{}

}

func writeCommand(data []byte, responseLength int) (r []byte, err errorCom) {
	sessionLog.PrintL(LOG_METOD_START_INT, "--WRITE-COMMAND-- Старт  writeCommand")
	c := &serial.Config{Name: "COM3", Baud: 115200, ReadTimeout: time.Second * 5}
	s, e := serial.OpenPort(c)

	//	s := serial
	//	var e error = nil

	if e != nil {
		sessionLog.PrintL(LOG_ERROR, "--WRITE-COMMAND-- Cannot open COM port")
		sessionLog.PrintL(LOG_METOD_RESULT, "--WRITE-COMMAND-- Стоп  writeCommand")
		return nil, errorCom{false, 0}
	}
	defer s.Close()
	sessionLog.PrintL(LOG_METOD_INSIDE, "--WRITE-COMMAND-- Порт открыт")
	attempts := 10

	//  Отправляем команду
	for attempts > 0 {
		sessionLog.PrintL(LOG_METOD_INSIDE, "--WRITE-COMMAND-- Попытка "+fmt.Sprint(11-attempts))
		attempts--
		n, err := s.Write([]byte{ENC})
		if err != nil {
			sessionLog.PrintL(LOG_ERROR, "--WRITE-COMMAND-- Ошибка записи в COM порт")
			continue
		}
		sessionLog.PrintL(LOG_PORT_DATA, "--COM-PORT-SEND-- Отправлен ENC "+fmt.Sprint([]byte{ENC}))

		buf := make([]byte, 1024)
		n, err = s.Read(buf)
		commandSent := false
		if n > 0 {
			if buf[0] == NAK {
				sessionLog.PrintL(LOG_PORT_DATA, "--COM-PORT-RECEIVE-- Ожидаемый ответ: "+fmt.Sprint(buf[0:n]))
				// Отправка команды
				sendAttempts := 10
				for sendAttempts > 0 {
					sessionLog.PrintL(LOG_METOD_INSIDE, "--WRITE-COMMAND-- Отправка данных, попытка "+fmt.Sprint(11-sendAttempts))
					sendAttempts--
					n, err = s.Write(data)
					if err != nil {
						sessionLog.PrintL(LOG_ERROR, "--WRITE-COMMAND-- Ошибка записи в COM-порт "+fmt.Sprint(err))
						continue
					}
					sessionLog.PrintL(LOG_PORT_DATA, "--COM-PORT-SEND-- Oтправлена команда: "+fmt.Sprint(data))
					commandSent = true
					break
				}

			} else if buf[0] == ACK {
				// Принтер подготавливает ответное сообщение
				sessionLog.PrintL(LOG_PORT_DATA, "--COM-PORT-RECEIVE-- Получен ACK: "+fmt.Sprint(buf[0:n]))
				sessionLog.PrintL(LOG_METOD_INSIDE, "--WRITE-COMMAND-- Ждем 1 секунду и повторяем ")
				time.Sleep(1)
				continue
			} else {
				if buf[0] == STX {
					sessionLog.PrintL(LOG_PORT_DATA, "--COM-PORT-RECEIVE-- Ответ от предыдущего запроса: "+fmt.Sprint(buf[0:n]))
				} else {
					sessionLog.PrintL(LOG_PORT_DATA, "--COM-PORT-RECEIVE-- Неизвесный ответ: "+fmt.Sprint(buf[0:n]))
				}
				_, err := s.Write([]byte{ACK})
				if err != nil {
					sessionLog.PrintL(LOG_METOD_INSIDE, "--WRITE-COMMAND-- Ошибка записи в COM порт")
					continue
				}
				sessionLog.PrintL(LOG_PORT_DATA, "--COM-PORT-SEND-- Отправлен ACK "+fmt.Sprint([]byte{ACK}))
				continue
			}
		} else {
			sessionLog.PrintL(LOG_METOD_INSIDE, "--WRITE-COMMAND-- Таймаут при чтении порта.")
			_, err := s.Write([]byte{ACK})
			if err != nil {
				sessionLog.PrintL(LOG_ERROR, "--WRITE-COMMAND-- Ошибка записи в COM порт")
				continue
			}
			sessionLog.PrintL(LOG_PORT_DATA, "--COM-PORT-SEND-- Отправлен ACK "+fmt.Sprint([]byte{ACK}))
			continue
		}

		// Ждем подтверждения команды\
		immediateAnswer := false
		if commandSent {
			// ожидаем ACK
			n, err = s.Read(buf)
			if err != nil {
				// Ошибка ком - порта, повторяем
				sessionLog.PrintL(LOG_METOD_INSIDE, "--WRITE-COMMAND-- Ошибка порта, повторяем..."+fmt.Sprint(err))
				continue
			} else if n == 0 {
				sessionLog.PrintL(LOG_METOD_INSIDE, "--WRITE-COMMAND-- Таймаут при чтении порта")
				continue
			}
			sessionLog.PrintL(LOG_PORT_DATA, "--COM-PORT-RECEIVE-- Получен ответ: "+fmt.Sprint(buf[0:n]))

			if buf[0] == NAK {
				// Если принят NAK повторяем команду
				sessionLog.PrintL(LOG_METOD_INSIDE, "--WRITE-COMMAND-- Регистратор вернул NAK ")
				continue
			} else if buf[0] == ACK {
				sessionLog.PrintL(LOG_METOD_INSIDE, "--WRITE-COMMAND-- Команда подтверждена ")
				if n > 1 {
					sessionLog.PrintL(LOG_METOD_INSIDE, "--WRITE-COMMAND-- Регистратор вернул ACK и ответ на команду. "+fmt.Sprint(buf[0:n]))
					for i := 1; i < n; i++ {
						buf[i-1] = buf[i]
					}
					n--
					immediateAnswer = true
				}
			}
		} else {
			// Если команда не отправлена, делаем еще попытку
			continue
		}

		// Принимаем ответ
		commandBuffer := make([]byte, 1024)

		commandBufferId := 0

		var length int = 0
		commandCorrect := false

		for true {
			sessionLog.PrintL(LOG_METOD_INSIDE, "--WRITE-COMMAND-- Читаем ответ")
			if immediateAnswer {
				immediateAnswer = false
			} else {
				n, err = s.Read(buf)
			}

			if n == 0 {
				sessionLog.PrintL(LOG_METOD_INSIDE, "--WRITE-COMMAND-- Пустой ответ - таймаут")
				break
			}

			sessionLog.PrintL(LOG_PORT_DATA, "--COM-PORT-RECEIVE-- Получен ответ: "+fmt.Sprint(buf[0:n]))

			for i := 0; i < n; i++ {
				commandBuffer[commandBufferId] = buf[i]
				commandBufferId++
			}

			sessionLog.PrintL(LOG_METOD_INSIDE, "--WRITE-COMMAND-- Состояние CommandBufferId: "+fmt.Sprint(commandBuffer[0:commandBufferId]))

			if length == 0 {
				if commandBuffer[0] == STX {
					if commandBuffer[1] > 0 {
						// Получена длина пакета
						length = int(commandBuffer[1])
						sessionLog.PrintL(LOG_METOD_INSIDE, "--WRITE-COMMAND-- Длина пакета: "+fmt.Sprint(length))
					}
				} else if commandBuffer[0] == NAK {
					sessionLog.PrintL(LOG_METOD_INSIDE, "[Регистратор вернул NAK "+fmt.Sprint(commandBuffer[0:commandBufferId]))
					//  Отправляем ENC и сбрасываем буферы
					length = 0
					commandBufferId = 0
					_, err := s.Write([]byte{ENC})
					if err != nil {
						sessionLog.PrintL(LOG_ERROR, "--WRITE-COMMAND-- Ошибка записи в COM порт")
					}
					sessionLog.PrintL(LOG_PORT_DATA, "--COM-PORT-SEND-- Отправлен ENC "+fmt.Sprint([]byte{ENC}))
					continue

				} else {
					sessionLog.PrintL(LOG_ERROR, "[Регистратор вернул неизвестный ответ "+fmt.Sprint(commandBuffer[0:commandBufferId]))
					_, err := s.Write([]byte{ACK})
					if err != nil {
						sessionLog.PrintL(LOG_ERROR, "--WRITE-COMMAND-- Ошибка записи в COM порт")
						continue
					}
					sessionLog.PrintL(LOG_PORT_DATA, "--COM-PORT-SEND-- Отправлен ACK "+fmt.Sprint([]byte{ACK}))
					break
				}
			}
			if length > 0 && length <= commandBufferId-3 {
				commandCorrect = true
				sessionLog.PrintL(LOG_METOD_INSIDE, "--WRITE-COMMAND-- Корректный ответ принят: "+fmt.Sprint(commandBuffer[0:commandBufferId]))
				_, err := s.Write([]byte{ACK})
				if err != nil {
					sessionLog.PrintL(LOG_ERROR, "--WRITE-COMMAND-- Ошибка записи в COM порт")
					continue
				}
				sessionLog.PrintL(LOG_PORT_DATA, "--COM-PORT-SEND-- Отправлен ACK "+fmt.Sprint([]byte{ACK}))
				break
			}
		}

		if commandCorrect {
			//  Проверяем контрольную сумму

			var lrc byte = 0
			for i := 0; i < commandBufferId; i++ {
				lrc ^= commandBuffer[i+1]
			}
			if lrc != 0 {
				//				_, _ = s.Write([]byte{NAK})
				sessionLog.PrintL(LOG_ERROR, "--WRITE-COMMAND-- Неправильнаяч контрольная сумма"+fmt.Sprint(commandBuffer[0:commandBufferId]))
				break
			} else {
				sessionLog.PrintL(LOG_METOD_INSIDE, "--WRITE-COMMAND-- Контрольная сумма корректна")
				// проверяем код ошибки
				var errorByte byte = 0

				commandData := commandBuffer[3 : commandBufferId-1]

				if commandBuffer[3] == 0 || (commandBuffer[2] == 255 && commandBuffer[4] == 0) {
					sessionLog.PrintL(LOG_METOD_INSIDE, "--WRITE-COMMAND-- Команда выполнена успешно"+fmt.Sprint(commandData))
				} else {
					if commandBuffer[2] == 255 {
						errorByte = commandBuffer[4]
					} else {
						errorByte = commandBuffer[3]
					}
					sessionLog.PrintL(LOG_ERROR, "--WRITE-COMMAND-- Ошибка при выпоннении команды "+
						fmt.Sprint(errCode(errorByte))+"  "+fmt.Sprint(commandBuffer[0:commandBufferId]))
				}
				sessionLog.PrintL(LOG_METOD_RESULT_INT, "--WRITE-COMMAND-- Стоп writeCommand")
				return commandData, errorCom{true, errorByte}
			}
		} else {
			sessionLog.PrintL(LOG_ERROR, "--WRITE-COMMAND-- Корректный ответ НЕ принят: "+fmt.Sprint(commandBuffer[0:commandBufferId]))
		}
	}

	sessionLog.PrintL(LOG_METOD_RESULT_INT, "--WRITE-COMMAND-- Стоп writeCommand")
	return nil, errorCom{true, 0}
}

func errCode(code byte) (str string) {
	switch code {
	case 0x00:
		return "Нет ошибки"
	case 0x15:
		return "Смена уже открыта"
	case 0x16:
		return "Смена не открыта"
	case 0x45:
		return "Cумма всех типов оплаты меньше итога чека"
	case 0x50:
		return "Принтер занят"
	case 0x73:
		return "Команда не поддерживается в данном режиме"
	case 0x7e:
		return "Неверное значение в поле длины"
	case 0x80:
		return "Ошибка связи с ФП"
	default:
		return "Неизвестная ошибка " + strconv.Itoa(int(code))
	}
}
