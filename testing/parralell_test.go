package testing
import (
	"dicomsend/parralels"
	"log"
	"time"
"testing"
)
type Dummy struct {

}
func (Dummy)DoAction(pb* parralels.ParralelsBallancer,data interface{}){
time.Sleep(time.Millisecond*1000)
log.Println("info: data=",data)

}
func TestParralelsBallancer(t *testing.T) {

	pb:=parralels.ParralelsBallancer{}
	pb.MaxParralels=20
	pb.MaxQuied=40
	pb.Pb=Dummy{}
	pb.Done=make(chan bool)
 	for i:=0;i<2000;i++{
		pb.StartNew(i)
	}
	pb.WaitAll()
	log.Println("done!!!!!!!")
}/*
import (
"testing"
	"dicomsend/parralels"
	"log"
)
type Dummy struct {

}
func (Dummy)DoAction(pb* parralels.ParralelsBallancer,data interface{}){
	log.Println("info: data=",data)
}
func TestMultiPartDownloadPool(t *testing.T) {
}*/