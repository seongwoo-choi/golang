package myutils

import (
	"fmt"
	"net/http"
	"time"

	"github.com/robfig/cron"
)

func EightCron(c *cron.Cron) {
	_ = c.AddFunc("00 00 08 * * *", func() {
		t := time.Now()
		now := t.Format("20060102")
		base_time := "0800"
		numOfRows := "230"
		url := fmt.Sprintf("http://localhost:3000/weather?base_date=%v&base_time=%v&numOfRows=%v", now, base_time, numOfRows)
		response, err := http.Get(url)
		c.Stop()
		if err != nil {
			fmt.Println(err)
		}
		defer response.Body.Close()
	})
	c.Start()
}

func SeventeentCron(c *cron.Cron) {
	_ = c.AddFunc("00 * * * * *", func() {
		t := time.Now()
		now := t.Format("20060102")
		base_time := "1700"
		numOfRows := "84"
		url := fmt.Sprintf("http://localhost:3000/weather?base_date=%v&base_time=%v&numOfRows=%v", now, base_time, numOfRows)
		response, err := http.Get(url)
		c.Stop()
		if err != nil {
			fmt.Println(err)
		}
		defer response.Body.Close()
	})
	c.Start()
}
