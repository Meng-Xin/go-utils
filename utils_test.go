package utils

import (
	"fmt"
	"golang.org/x/time/rate"
	"testing"
	"time"
)

type MyTaskQueue struct {
	msg string
}

func (m *MyTaskQueue) Consumer(bytes []byte) error {
	var err error
	fmt.Println(bytes)
	return err
}

func (m *MyTaskQueue) MsgContent() string {
	return m.msg
}

func TestTaskQueue(t *testing.T) {
	msg := fmt.Sprintf("这是测试任务")
	task := &MyTaskQueue{
		msg,
	}
	queueExchange := &QueueExchange{
		"test.rabbit",
		"rabbit.key",
		"test.rabbit.mq",
		"direct",
	}
	mq := New(queueExchange)
	mq.RegisterProducer(task)
	mq.RegisterReceiver(task)
	mq.Start()
	time.Sleep(time.Second * 5)
}

func TestRate(t *testing.T) {
	limiter := rate.NewLimiter(100, 200)
	for i := 0; i < 1000; i++ {
		go func() {
			if limiter.Allow() {
				fmt.Printf("Success 第:%d号协程\n", i)
			} else {
				fmt.Printf("Error 第:%d号协程\n", i)
			}
		}()
	}
}
