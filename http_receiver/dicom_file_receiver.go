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

type OnDownloadFunc func(path string) (error)

type DicomReceiver struct {
	buffer   []byte
	mprd     *multipart.Reader
	rw       http.ResponseWriter
	req      *http.Request
	onDown   OnDownloadFunc
	dataSize int64
	upPos    int64

}

func CreateReciver(rw_ http.ResponseWriter, req_ *http.Request, odf OnDownloadFunc) (*DicomReceiver, error) {
	if req_ == nil {
		return nil, errors.New("error: empty http request")
	}
	var fs DicomReceiver
	fs.rw = rw_
	fs.req = req_
	fs.onDown = odf
	return &fs, nil
}

func (fs *DicomReceiver) DoWork() (bool, error) {
	defer
	func() {
		log.Printf("info: uploaded data %v percent \n", (fs.upPos * 100) / fs.dataSize)
	}()
	p, err := fs.mprd.NextPart()
	if err != nil {
		log.Println("info: part is out of")
		fs.upPos = fs.dataSize
		return true, nil
	}

	fn := p.FormName()
	switch fn {
	case "files":{
		//if f, er := os.Create(os.TempDir() + sep() + genUid()); er != nil {
		if f, er := os.Create("/home/andrew/Desktop/da/" + sep() + genUid() + ".jpeg"); er != nil {
			defer f.Close()
			log.Println("error: can't create temp file")
			return false, er
		}else {
			for {
				if count, e := p.Read(fs.buffer); e == io.EOF {
					log.Printf("info: file %v writed to disk \n", p.FileName())
					return false, nil
				}else {
					f.Write(fs.buffer[0:count])
					fs.upPos += int64(count)
				}

			}
		}
	}
	case "server":{
		var c int
		if c, e := p.Read(fs.buffer); e == io.EOF {
			log.Printf("error: can't read server form \n")
			return false, e
		}else {
			s := string(fs.buffer[:c])
			log.Println("info: server is ", s)
		}
		fs.upPos += int64(c)
	}
	case "port":{
		var c int
		if c, e := p.Read(fs.buffer); e == io.EOF {
			log.Printf("error: can't read port form \n")
			return false, e
		}else {
			s := string(fs.buffer[:c])
			log.Println("info: port is ", s)
		}
		fs.upPos += int64(c)
	}
	case "aetitle":{
		var c int
		if c, e := p.Read(fs.buffer); e == io.EOF {
			log.Printf("error: can't read aetitle form \n")
			return false, e
		}else {
			s := string(fs.buffer[:c])
			log.Println("info: aetitle is ", s)
		}
		fs.upPos += int64(c)
	}
	default:{
		log.Println("warning: unsupported form")
	}
	}

	return false, nil
}
func (fs *DicomReceiver) GetProgress() interface{} {
	return (fs.upPos * 100) / fs.dataSize
}
func (fs *DicomReceiver) BeforeRun() error {
	var err error
	fs.mprd, err = fs.req.MultipartReader()
	fs.dataSize = fs.req.ContentLength
	if err != nil {
		return err
	}
	fs.buffer = make([]byte, BufferSize)
	return nil
}
func (fs *DicomReceiver) AfterStop() error {
	return nil
}