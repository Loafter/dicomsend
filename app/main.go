package main
import "dicomsend/service"
func main(){
 ups:=FileReciverSrv.UpSrv{}
	ups.Start(9982)
}
