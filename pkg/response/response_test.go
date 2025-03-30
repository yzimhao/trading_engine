package response_test

import (
	"encoding/json"
	"testing"

	"github.com/yzimhao/trading_engine/v2/pkg/response"
)

func TestResponseSuccess(t *testing.T) {
	expect := `{"code":0}`
	res := response.Success()
	raw, _ := json.Marshal(res)

	if string(raw) != expect {
		t.Errorf("expect %s, got %s", expect, string(raw))
	}
}

func TestResponseSuccessWithData(t *testing.T) {
	expect := `{"code":0,"data":"hi"}`
	res := response.Success().WithData("hi")
	raw, _ := json.Marshal(res)

	if string(raw) != expect {
		t.Errorf("expect %s, got %s", expect, string(raw))
	}
}

func TestResponseFail(t *testing.T) {
	expect := `{"code":2,"message":"unknown error"}`
	res := response.Fail()
	raw, _ := json.Marshal(res)

	if string(raw) != expect {
		t.Errorf("expect %s, got %s", expect, string(raw))
	}
}

func TestResponseFailWithCustomError(t *testing.T) {
	expect := `{"code":21000,"message":"custom error"}`
	res := response.Fail().WithError(21000, "custom error")
	raw, _ := json.Marshal(res)

	if string(raw) != expect {
		t.Errorf("expect %s, got %s", expect, string(raw))
	}
}
