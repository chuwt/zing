package recorder

import (
	"testing"
	"github.com/chuwt/zing/config"
	"github.com/chuwt/zing/object"
	"github.com/chuwt/zing/recorder/huobi"
)

func TestRecorderHuobi(t *testing.T) {

	recorder := NewRecorder(config.Config.Redis)
	recorder.AddPublisher(object.GatewayHuobi, huobi.NewPublisher)

	if err := recorder.Init(); err != nil {
		t.Error(err)
		return
	}
	recorder.Run()
}
