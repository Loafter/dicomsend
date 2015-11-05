package FileReciverSrv
import (
	"net/http"
	"strconv"
	"io/ioutil"
	"log"
	"encoding/base64"
	/*"os"
	"io"*/
	"crypto/rand"
	"fmt"
)
type UpSrv struct {

}
const  htmlData=""
const FlushDiskSize = 1024 * 1024

func genUid() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
func (srv *UpSrv) Start(listenPort int) error {
	http.HandleFunc("/", srv.Redirect)
	http.HandleFunc("/index.html", srv.index)
	http.HandleFunc("/upload_dicom", srv.uploadDicom)
	if err := http.ListenAndServe(":"+strconv.Itoa(listenPort), nil); err != nil {
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
	log.Println("try upload")
	rd,err :=r.MultipartReader()
	if err!=nil{
		log.Println("error:"+ err.Error())
	}
	//buffer := make([]byte, FlushDiskSize, FlushDiskSize)
	for p,err:=rd.NextPart();err==nil;p,err=rd.NextPart(){

		/*c,er:=p.Read(buffer){
			if er!=nil{
				return
			}
		}*/
		log.Println(p.)

	}


	/* err != nil {
		log.Println(w,err)
		return
	}
	formdata := r.MultipartForm  // ok, no problem so far, read the Form data
	//get the *fileheaders
	files := formdata.File["multiplefiles"]  // grab the filenames
	for i, _ := range files {  // loop through the files one by one
		log.Println(i)
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			return
		}

		out, err := os.Create("/tmp/" + files[i].Filename)

		defer out.Close()
		if err != nil {

			return
		}

		_, err = io.Copy(out, file) // file not files[i] !

		if err != nil {

			return
		}

		log.Println(w,"Files uploaded successfully : ")
		log.Println(w, files[i].Filename + "\n")

	}*/
}

