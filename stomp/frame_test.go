package stomp

import (
	"fmt"
	"testing"
)

func TestFrame(t *testing.T) {

	messageFrameWhoami()
}

func messageFrameWhoami() {
	body := `{"statue":1,"Name":"stompserver"}`
	messageFrame, _ := NewMessageFrame("/broadcast/whoami", "90246c88-70a0-45dd-9cbf-cf948f56639d-1", "sub-0")
	messageFrame.SetBody(body)
	fmt.Println(messageFrame.Serialize())
}
