package FileReciverSrv
import (
	"net/http"
	"strconv"
	"io/ioutil"
	"log"
	"encoding/base64"
	"os"
	"io"
)
type UpSrv struct {

}
const  htmlData=""
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
	log.Println("try upload")
	err := r.ParseMultipartForm(200000) // grab the multipart form
	if err != nil {
		log.Println(w,err)
		return
	}


	formdata := r.MultipartForm  // ok, no problem so far, read the Form data

	//get the *fileheaders
	files := formdata.File["multiplefiles"]  // grab the filenames

	for i, _ := range files {  // loop through the files one by one
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

	}
}

