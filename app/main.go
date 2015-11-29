package main
import (
	"dicomsend/parralels"
	"log"
	"time"
)
type Dummy struct {

}
func (Dummy)DoAction(pb* parralels.ParralelsBallancer,data interface{}){
	time.Sleep(time.Millisecond*1000)
	log.Println("info: data=",data)

}
func main() {
	pb:=parralels.ParralelsBallancer{}
	pb.MaxParralels=1
	pb.Pb=Dummy{}
 	for i:=0;i<20;i++{
		pb.StartNew(i)
	}
	pb.WaitAll()
	log.Println("done!!!!!!!")
}
/*import (
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
*/