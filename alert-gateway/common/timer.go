package common

import (
	"time"
	"log"
)

func UpdateMaintainlist() {
	log.Println("timer: UpdateMaintainlist start ...")
}

func init()  {
	go func ()  {
		for {
			current := time.Now()
			time.Sleep(time.Duration(90-current.Second()) * time.Second)
			UpdateMaintainlist()
		}
	}()
}
