package utils

import (
	"os"
	"os/signal"
	"time"

	"code.aliyun.com/bim_backend/zoogoer/gun/errors"
)

// Do function to do the thing by interval until success
func Do(interval int, errstr string, fn func() (interface{}, error)) (interface{}, error) {
	d := time.Duration(interval) * time.Second
	timer := time.NewTimer(d)
	c := make(chan os.Signal)
	signal.Notify(c)

	defer timer.Stop()
	defer close(c)

	for {
		select {
		case <-timer.C:
			val, err := fn()
			if err == nil {
				return val, nil
			}
			timer.Reset(d)
		case <-c:
			return nil, errors.NewClientErr(nil, errstr, "", "DataBase.do", nil)
		}
	}
}
