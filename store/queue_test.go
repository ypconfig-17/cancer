package store

import (
	"fmt"
	"testing"
	"time"
)

func Test_Queue(t *testing.T) {

	queue, _ := openFileQueue("test", "queue_01")
	defer queue.file.Close()

	go func() {
		for i := 0; i < 1000; i++ {
			sprintf := fmt.Sprintf("%s%d", "你好啊", i)
			mesage1 := CreteMesage(queue, []byte(sprintf))
			queue.Add(mesage1)
		}
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			pollMsg := queue.Poll()
			if pollMsg != nil {
				fmt.Printf("取出数据:%s\n", pollMsg.ToString())
			}
		}

	}()

	time.Sleep(time.Second * 10)
}
