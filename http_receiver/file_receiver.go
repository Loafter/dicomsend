package http_receiver

//import "godownloader/monitor"

type FilesReceiver struct {

}
func (fs *FilesReceiver) Init(fs interface{}){

}
func (fs *FilesReceiver) DoWork() (bool, error) {

	return false, nil
}
func (fs *FilesReceiver) GetProgress() interface{} {
	return nil
}
func (fs *FilesReceiver) BeforeRun() error {
	return nil
}
func (fs *FilesReceiver) AfterStop() error {
	return nil
}