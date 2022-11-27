package server

import (
	pb "FormSubmit/grpc"
	"errors"
	"fmt"
	"net/mail"
	"reflect"
	"runtime"
)

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func EmailValidator(in *pb.FormData) bool {
	_, err := mail.ParseAddress(in.GetEmail())
	return err == nil
}
func AgeValidator(in *pb.FormData) bool {
	return in.GetAge() > 0
}
func HeightValidator(in *pb.FormData) bool {
	return in.GetHeight() > 0
}
func checkValidators(in *pb.FormData, validators []func(data *pb.FormData) bool) error {
	for _, validator := range validators {
		if !validator(in) {
			failedValidator := GetFunctionName(validator)
			errorText := fmt.Sprintf("validator %v failed", failedValidator)
			return errors.New(errorText)
		}
	}
	return nil
}
