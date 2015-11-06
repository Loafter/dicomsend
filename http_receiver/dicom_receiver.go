package httpreciver
import (
	"godownloader/monitor"
	"log"
	"net/http"
)

func dummyOnFileDownload(path string)error {
	log.Println("info: i do some thing with file", path)
	return nil
}

type DicomReciver struct {
	drs map[string]*monitor.MonitoredWorker
}
func CreateDicomReciver() (*DicomReciver, error) {
	dr:=new(DicomReciver)
	dr.drs = make(map[string]*monitor.MonitoredWorker)
	return dr,nil
}
func (dr *DicomReciver) AddReq(rw_ http.ResponseWriter,req_ *http.Request) (error) {
	if fr, err := CreateReciver(rw_, req_, dummyOnFileDownload); err != nil {
		return err
	}else {
		mw := monitor.MonitoredWorker{Itw:fr}
		dr.drs[mw.GetId()]=&mw
		mw.Start()
	}
	return nil
}