package logging

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

func SettingLogger(level log.Level, folder string) {
	log.SetLevel(level)
	log.SetFormatter(&easy.Formatter{
		TimestampFormat: "15:04:05",
		LogFormat:       "[%lvl%]: %time% - %msg%\n",
	})

	path := filepath.Join("logs", folder)

	err := os.MkdirAll(path, os.ModePerm)

	if err != nil {
		log.Fatalf("Couldn't create directory:\n%v\n", err)
	}

	name := filepath.Join(path, fmt.Sprintf("logs - %v.txt", time.Now().Format("2006-01-02-15-04-05")))

	file, err := os.Create(name)

	if err != nil {
		log.Fatalf("Error creating log files: %v", err)
	}

	mw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)
}
