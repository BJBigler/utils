package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"time"
)

//Log string to file
func Log(message ...interface{}) {

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(message...)
		log.Printf("%v", err)
		return
	}

	logname := time.Now().Format("20060102_log.txt")
	directory := path.Join(pwd, "logs")
	logfile := path.Join(pwd, "logs", logname)

	//Does the directory exist?
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		fmt.Println("utils.Log(): directory ", directory, "does not exist; skipping logging ; message is ", message)
	}

	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)

	defer closeFile(f)

	if err != nil {
		log.Fatal(err)
	}

	text := fmt.Sprintf("%v%v%v%v", time.Now().Format("2006-01-02 15:04:05"), "\r\n", message, "\r\n\r\n")

	if _, err := f.WriteString(text); err != nil {
		log.Fatal(err)
	}
}

func closeFile(f *os.File) {
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

//LogToStdError ...
func LogToStdError(namePostfix string, message interface{}) (err error) {

	dir, err := os.Getwd()

	//Do we have a logs folder?
	logfolder := path.Join(dir, "logs")
	err = CreateFolder(logfolder)

	logname := time.Now().Format("20060102")

	if namePostfix == "" {
		logname += "_log"
	} else {
		logname += "_" + namePostfix
	}

	logname += ".txt"

	logfile := path.Join(dir, "logs", logname)

	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	multi := io.MultiWriter(f, os.Stdout)
	logger := log.New(multi, "", log.Ldate|log.Ltime|log.Lshortfile)

	logger.Println(message)

	return err
}

//TimeTrack should be called like
//defer TimeTrack(time.Now(), "caller")
func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	LogToStdError("timer", fmt.Sprintf("%s took %s", name, elapsed))
}
