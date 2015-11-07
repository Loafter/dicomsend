package filereciversrv
import (
	"net/http"
	"strconv"
	"io/ioutil"
	"log"
	"encoding/base64"
	"crypto/rand"
	"fmt"
	"os"
	"dicomsend/http_receiver"
	"godownloader/monitor"
)

const htmlData = ""
const FlushDiskSize = 1024 * 1024

func genUid() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func sep() string {
	st := strconv.QuoteRune(os.PathSeparator)
	st = st[1 : len(st) - 1]
	return st
}
type UpSrv struct {
drs map[string]*monitor.MonitoredWorker
}

func (srv *UpSrv) Start(listenPort int) error {
	srv.drs=make(map[string]*monitor.MonitoredWorker)
	http.HandleFunc("/", srv.Redirect)
	http.HandleFunc("/index.html", srv.index)
	http.HandleFunc("/upload_dicom", srv.uploadDicom)
	if err := http.ListenAndServe(":" + strconv.Itoa(listenPort), nil); err != nil {
		return err
	}
	return nil
}
func (srv *UpSrv) Redirect(responseWriter http.ResponseWriter, request *http.Request) {
	http.Redirect(responseWriter, request, "/index.html", 301)
}


func (srv *UpSrv) index(rwr http.ResponseWriter, req *http.Request) {
	rwr.Header().Set("Content-Type: text/html", "*")
	content, err := ioutil.ReadFile("index.html")
	if err != nil {
		log.Println("warning: start page not found, return included page")
		val, _ := base64.StdEncoding.DecodeString(htmlData)
		rwr.Write(val)
		return
	}
	rwr.Write(content)
}

/*func (srv *UpSrv)uploadDicom(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	rd, err := r.MultipartReader()
	if err != nil {
		log.Println("error: can't get MultipartReader")
		return
	}
	buffer := make([]byte, FlushDiskSize)
	for p, err := rd.NextPart(); err == nil; p, err = rd.NextPart() {
		if p.FormName() == "files" {
			//if f, er := os.Create(os.TempDir() + sep() + genUid()); er != nil {
			if f,er:=os.Create("C:\\Users\\andre\\Desktop\\Target"+sep()+genUid()); er!=nil{
				log.Println("error: can't create temp file")
				return
			}else {
				log.Println(p)
				for {
					if count, e := p.Read(buffer); e == io.EOF {
						log.Println("info: Last buffer read!")
						f.Close()
						break
					}else {
						log.Println(count)
						f.Write(buffer[0:count])
					}

				}
			}
		}
	}
}*/

func dummyOnFileDownload(path string)error {
	log.Println("info: i do some thing with file", path)
	return nil
}

func (srv *UpSrv)uploadDicom(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if fr, err := httpreciver.CreateReciver(w, r, dummyOnFileDownload); err != nil {
		http.Error(w, "error: can't create reciver", http.StatusInternalServerError)
	}else {
		mw := monitor.MonitoredWorker{Itw:fr}
		srv.drs[mw.GetId()]=&mw
		mw.Start()
		mw.Wait()
		log.Println(mw.GetState())
	}
	log.Println("info: finish upload")
}


