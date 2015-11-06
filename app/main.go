package main
import (
	"dicomsend/service"
	"log"
)
func main(){
 ups:=filereciversrv.UpSrv{}
	log.Println("info:try start service at http://localhost:9982")
	if err:=ups.Start(9982);err!=nil{
		log.Println("error: can't start service with reason"+err.Error())
		return
	}
}
