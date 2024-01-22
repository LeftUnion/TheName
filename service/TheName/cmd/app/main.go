package main

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/LeftUnion/theName/service/TheName/internal/app"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(frame *runtime.Frame) (string, string) {
			file := frame.File[len(path.Dir(os.Args[0]))+1:]
			line := frame.Line
			return "", fmt.Sprintf("%s:%d", file, line)
		},
	})
}

func main() {

	// go func() { // DEBUG Profiler
	// 	fmt.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	app.Run()
}
