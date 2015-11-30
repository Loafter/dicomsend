package testing

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