package parralels
import (
	"container/list"
	"log"
	"sync"
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
func (pb *ParralelsBallancer) ActiveJobs()int {
	return pb.activeJobs

}

func (pb *ParralelsBallancer) SleepedJobs()int {
	if pb.dats==nil{
		return 0
	}
	return pb.dats.Len()

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
		go pb.startParallel(data)
	}else {
		log.Println("info: overflow parralels jobs")
		pb.wgrun.Add(1)
		pb.dats.PushBack(data)
	}

}

func (pb *ParralelsBallancer) WaitAll() {
	log.Println("info: wait all")
	pb.wgrun.Wait()
}