package main

import (
	"os"
	"time"
	"strings"
	"fmt"
)

type Logger struct{
	file *os.File
	level int
}

const LOG_FATAL_ERROR  		= 	0;
const LOG_ERROR  			= 	5;
const LOG_INFO				=	10;
const LOG_METOD_RESULT		=	15;
const LOG_METOD_START		=	20;
const LOG_METOD_RESULT_INT	=	23;
const LOG_METOD_START_INT	=	24;
const LOG_METOD_INSIDE		=	25;
const LOG_PORT_DATA			=	30;

func (l *Logger)NewSessionLogger()(){
	path := "sessionLogs";
	l.level = LOG_INFO;
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0777)
	}

	currentTime := l.getTime()
	l.file,_ = os.Create(path+string(os.PathSeparator)+currentTime+".log");
	msg := currentTime+" --INIT-- Старт сессии\n"
	fmt.Print(msg)
	_,err := l.file.WriteString(msg)
	if(err != nil){
		fmt.Println("Ошибка записи файла", err)
	}
}

func (l *Logger)PrintL(level int, msg string)() {
	line := l.getTime()+" "+msg+"\n"
	if level <= l.level {
		fmt.Print(line)
	}
	_,err := l.file.WriteString(line)
	if(err != nil){
		fmt.Println("Ошибка записи файла", err)
	}
}

func (l *Logger)Print(msg string)(){
	line := l.getTime()+" "+msg+"\n"
	fmt.Print(line)
	_,err := l.file.WriteString(line)
	if(err != nil){
		fmt.Println("Ошибка записи файла", err)
	}
}

func (l *Logger)Close()() {
	msg := l.getTime()+" --INIT-- Окончание сессии\n"
	fmt.Print(msg)
	_,err := l.file.WriteString(msg)
	if(err != nil){
		fmt.Println("Ошибка записи файла", err)
	}
	l.file.Close()
}

func (l *Logger)getTime()(string){
	return 	strings.Replace(time.Now().Format("2006-01-02 15:04:05.000000"),":",".",-1);
}