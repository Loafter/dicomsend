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
	"encoding/json"
)

const htmlData = ""
const FlushDiskSize = 1024 * 1024


func genUid() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
type Upst struct {
	Uid      string
	Progress int64
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
	srv.drs = make(map[string]*monitor.MonitoredWorker)
	http.HandleFunc("/", srv.Redirect)
	http.HandleFunc("/index.html", srv.index)
	http.HandleFunc("/upload_dicom", srv.uploadDicom)
	http.HandleFunc("/progress_upload.json", srv.progressJson)
	http.HandleFunc("/delete_upload", srv.deleteUpload)
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



func (srv *UpSrv)uploadDicom(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if fr, err := httpreciver.CreateReciver(w, r,httpreciver.OnDicomDownload); err != nil {
		http.Error(w, "error: can't create reciver", http.StatusInternalServerError)

	}else {
		mw := monitor.MonitoredWorker{Itw:fr}
		srv.drs[mw.GetId()] = &mw
		mw.Start()
		mw.Wait()
		log.Println(mw.GetState())
	}
	log.Println("info: finish upload")
}


func (srv *UpSrv) progressJson(rwr http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	rwr.Header().Set("Access-Control-Allow-Origin", "*")
	jbs := make([]Upst, 0, len(srv.drs))
	for ind, i := range srv.drs {
		if i.GetState() == monitor.Running {


			prs, _ := i.GetProgress().(int64)
			st := Upst{Uid:ind, Progress:prs}
			jbs = append(jbs, st)
		}else {
			delete(srv.drs, ind)
		}

	}
	js, err := json.Marshal(jbs)
	if err != nil {
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		return
	}
	rwr.Write(js)

}
func (srv *UpSrv) deleteUpload(rwr http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	rwr.Header().Set("Access-Control-Allow-Origin", "*")
	bodyData, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	var uid string
	if err := json.Unmarshal(bodyData, &uid); err != nil {
		http.Error(rwr, err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}
	if val, ok := srv.drs[uid]; ok {
		val.Stop()
		delete(srv.drs, uid)
		log.Println("info: remove upload jobs with id",uid)
	}else {
		log.Println("warning: can't find upload with id ", uid)
}
	rwr.Write([]byte{0})

}
