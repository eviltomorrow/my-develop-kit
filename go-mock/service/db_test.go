package service

import (
	"testing"

	"github.com/golang/mock/gomock"
)

func TestGetFromDB(t *testing.T) {
	var ctrl = gomock.NewController(t)
	defer ctrl.Finish()

	var m = NewMockDB(ctrl)
	m.EXPECT().Get(gomock.Eq("Tom")).Return("Hello, Tom", nil)

	data, err := GetFromDB(m, "Tom")
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	t.Logf("Data: %v", data)

}
