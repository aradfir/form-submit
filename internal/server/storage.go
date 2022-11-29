package server

import (
	pb "FormSubmit/internal/grpc"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

func formToStr(data *pb.FormData) string {
	return fmt.Sprintf("%v ### %v ### %v ### %v ### %v",
		data.GetFirstName(), data.GetLastName(), data.GetEmail(), data.GetAge(), data.GetHeight())
}
func safeClose(f *os.File) {
	err := f.Close()
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("File close error")
	}
}

func WriteToFile(data *pb.FormData) error {
	f, err := os.OpenFile("users.form", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("DB file open failed")
		return err
	}
	defer safeClose(f)

	_, err = f.WriteString(formToStr(data) + "\n")
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("DB file write failed")
		return err
	}
	return nil
}
