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
	"os/exec"
	"dicomsend/parralels"
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

type DicomSendData struct {
	Server string
	Port string
	AET	string
	FileName string
}

type ParallelDicomSend struct {

}
func (ParallelDicomSend)DoAction(pb* parralels.ParralelsBallancer,data interface{}){
	ds:=data.(DicomSendData)
	gdcmscu := os.Getenv("GDCMSCUP")
	log.Println(gdcmscu, "--store", "-H", ds.Server, "-p", ds.Port,
		"--call",ds.AET,
		"--aetitle", "AE_WEBCLI",
		"-i", ds.FileName,
		"-D")

	if out, err := exec.Command(gdcmscu, "--store", "-H", ds.Server, "-p", ds.Port,
		"--call",ds.AET,
		"--aetitle", "AE_WEBCLI",
		"-i", ds.FileName,
		"-D",
		"-D").Output(); err != nil {
		log.Printf("dicom send status : %s %s \n", out, err)
		return
	} else {
		log.Printf("success: %s\n", out)
	}
	if err := os.Remove(ds.FileName); err != nil {
		return
	}
	return
}

const BufferSize = 1024 * 1024


type DicomReceiver struct {
	buffer   []byte
	mprd     *multipart.Reader
	rw       http.ResponseWriter
	req      *http.Request
	ps parralels.ParralelsBallancer
	dataSize int64
	upPos    int64
	server   string
	port     string
	aet      string

}

func CreateReciver(rw_ http.ResponseWriter, req_ *http.Request) (*DicomReceiver, error) {
	if req_ == nil {
		return nil, errors.New("error: empty http request")
	}
	var fs DicomReceiver
	fs.rw = rw_
	fs.req = req_
	fs.ps = parralels.ParralelsBallancer{Pb:ParallelDicomSend{},MaxParralels:10}
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
		if f, er := os.Create(os.TempDir() + sep() + genUid()); er != nil {
			log.Println("error: can't create temp file")
			return false, er
		}else {
			defer func() {
				f.Close()
				log.Println("info: file closed")
				dsd:=DicomSendData{Server:fs.server,Port:fs.port,AET:fs.aet,FileName:f.Name()}
				fs.ps.StartNew(dsd)
			}()
			for {
				if count, e := p.Read(fs.buffer); e == io.EOF {
					log.Printf("info: file %v writed to disk \n", )
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
			fs.server = s
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
			fs.port = s
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
			fs.aet = s
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
	log.Println("info: body closed")
	fs.rw.Write([]byte{0})
	fs.req.Body.Close()
	return nil
}