package testing
/*import (
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
	pb.MaxParralels=30
	pb.Pb=Dummy{}
 	for i:=0;i<200;i++{
		pb.StartNew(i)
	}
	pb.WaitAll()
	log.Println("done!!!!!!!")
}*/
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
}