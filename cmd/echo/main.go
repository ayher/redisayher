package main

import (
	"github.com/hdt3213/godis/lib/logger"
	"github.com/hdt3213/godis/tcp"
	"redisayher/internal/echo"
)

func main(){
	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       "godis",
		Ext:        "log",
		TimeFormat: "2006-01-02",
	})

	err:=echo.ListenAndServeWithSignal(&tcp.Config{
		Address: ":8080",
	},echo.MakeEchoHandler())
	if err != nil {
		logger.Error(err)
	}
}
