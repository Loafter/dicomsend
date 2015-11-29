package parralels
import (
	"container/list"
	"log"
	"sync"
	//"time"
)
type PbAction interface {
	DoAction(pb*ParralelsBallancer, data interface{})
}
type  ParralelsBallancer struct {
	wgrun        sync.WaitGroup
	lmut         sync.Mutex
	activeJobs   int
	MaxParralels int
	Pb           PbAction
	dats         *list.List
}

func (pb *ParralelsBallancer) startParallel(data interface{}) {
	pb.Pb.DoAction(pb, data)
	pb.lmut.Lock()
	defer func(){
		pb.activeJobs -= 1
		pb.lmut.Unlock()
		pb.wgrun.Done()
	}()
	if (pb.activeJobs <= pb.MaxParralels) {
		if el := pb.dats.Front(); el == nil {
			log.Println("info: job list is empty")
			return
		}else {
			pb.dats.Remove(el)
			pb.wgrun.Done()
			go pb.StartNew(el.Value)
		}
	}


}
func (pb *ParralelsBallancer) StartNew(data interface{}) {
	log.Println("info: try start")
	pb.lmut.Lock()
	defer pb.lmut.Unlock()
	if pb.dats == nil {
		log.Println("info: list not inited")
		pb.dats = list.New()
	}
	if pb.activeJobs < pb.MaxParralels {
		pb.activeJobs += 1
		pb.wgrun.Add(1)
		log.Println("info: start parralels job",pb.activeJobs)
		go pb.startParallel(data)
	}else {
		log.Println("info: overflow parralels jobs")
		pb.wgrun.Add(1)
		pb.dats.PushBack(data)
	}

}

func (pb *ParralelsBallancer) WaitAll() {
	log.Println("info: wait all")
	/*for count:=pb.dats.Len();count!=0;count=pb.dats.Len(){
		time.Sleep(100*time.Millisecond)
	}*/
	pb.wgrun.Wait()
}