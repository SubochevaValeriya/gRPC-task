package main

import (
	"errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	gRPC_task "rusprofile/proto"
	"testing"
)

type TestCaseGetInfo struct {
	INN           int64
	Name          string
	ExpectedError error
}

func TestGetInfo(t *testing.T) {

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	s := server{}

	testCases := []TestCaseGetInfo{
		{INN: 7802836667, Name: "Valid INN", ExpectedError: nil},
		{INN: 81039203331131, Name: "NotValid", ExpectedError: errors.New("incorrect format for INN")},
	}

	for _, cse := range testCases {
		cse := cse
		t.Run(cse.Name, func(t *testing.T) {
			req := gRPC_task.Request{INN: cse.INN}
			_, err := s.GetInfo(context.Background(), &req)
			if cse.ExpectedError == nil {
				if err != nil {
					t.Errorf("should not be error")
				}
			} else if err == nil {
				t.Errorf("should be an error")
			}
		})
	}
}
