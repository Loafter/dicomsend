package httpreciver
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
const BufferSize = 1024 * 1024

type OnDownloadFunc func(path string)(error)

type FilesReceiver struct {
	buffer []byte
	mprd *multipart.Reader
	rw http.ResponseWriter
	req *http.Request
	onDown OnDownloadFunc

}

func CreateReciver(rw_ http.ResponseWriter,req_ *http.Request,odf OnDownloadFunc)(*FilesReceiver,error){
	if req_==nil{
		return nil,errors.New("error: empty http request")
	}
	var fs FilesReceiver
	fs.rw=rw_
	fs.req=req_
	fs.onDown=odf
	return &fs,nil
}
func (fs *FilesReceiver) DoWork() (bool, error) {

	p, err := fs.mprd.NextPart()
	if err!=nil{
		fs.req.Body.Close()
		return false,err
	}
	if p.FormName() == "files" {
		//if f, er := os.Create(os.TempDir() + sep() + genUid()); er != nil {
		if f,er:=os.Create("C:\\Users\\andre\\Desktop\\Target"+sep()+genUid()); er!=nil{
			log.Println("error: can't create temp file")
			fs.req.Body.Close()
			return false,er

		}else {
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
		fs.req.Body.Close()
		return err
	}
	fs.buffer = make([]byte, BufferSize)

	return nil
}
func (fs *FilesReceiver) AfterStop() error {
	return nil
}