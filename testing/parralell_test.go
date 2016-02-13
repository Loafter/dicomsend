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
	pb.MaxParralels=3
	pb.MaxQuied=60
	pb.Pb=Dummy{}
 	for i:=0;i<200;i++{
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