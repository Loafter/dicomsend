package HttpReciver
import (
	"net/http"
	"errors"
	"log"
	"os"
	"io"
	"mime/multipart"
	"strconv"
	"crypto/rand"
	"fmt"
)
func sep() string {
	st := strconv.QuoteRune(os.PathSeparator)
	st = st[1 : len(st) - 1]
	return st
}


func genUid() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

//import "godownloader/monitor"
const FlushDiskSize = 1024 * 1024

func dummyOnFileDownload(path string){
 log.Println("info: i do some thing with file",path)
}

type OnDownloadFunc func(path string)(error)
type FilesReceiver struct {
	buffer []byte
	mprd *multipart.Reader
	rw http.ResponseWriter
	req *http.Request
	onDown OnDownloadFunc

}

func (fs *FilesReceiver) Init(rw_ http.ResponseWriter,req_ *http.Request)error{

	if req_!=nil{
		return errors.New("error: empty http request")
	}
	fs.rw=rw_
	fs.req=req_
	fs.onDown=dummyOnFileDownload
	return nil
}
func (fs *FilesReceiver) DoWork() (bool, error) {
	p, err := fs.mprd.NextPart()
	if err!=nil{
		return false,nil
		fs.req.Body.Close()
	}
	if p.FormName() == "files" {
		//if f, er := os.Create(os.TempDir() + sep() + genUid()); er != nil {
		if f,er:=os.Create("C:\\Users\\212402712\\Desktop\\Target"+sep()+genUid()); er!=nil{
			log.Println("error: can't create temp file")
			return false,er
			fs.req.Body.Close()
		}else {
			log.Println(p)
			for {
				if count, e := p.Read(fs.buffer); e == io.EOF {
					log.Println("info: Last buffer read!")

					f.Close()
					break
				}else {
					log.Println(count)
					f.Write(fs.buffer[0:count])
				}

			}
		}
	}
return true,nil
}
func (fs *FilesReceiver) GetProgress() interface{} {
	return nil
}
func (fs *FilesReceiver) BeforeRun() error {
	var err error
	fs.mprd, err = fs.req.MultipartReader()
	if err != nil {
		return err
		fs.req.Body.Close()
	}
	fs.buffer = make([]byte, FlushDiskSize)
	return nil
}
func (fs *FilesReceiver) AfterStop() error {
	return nil
}