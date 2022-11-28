package application

import (
	"log"
	"os"
	"time"
)

var (
	Hostname = ""
	Name     = ""
)

var start = time.Now()

func init() {
	var err error
	if Hostname, err = os.Hostname(); err != nil {
		log.Fatal(err)
	}
}

type HealthzRespone struct {
	Uptime   string `json:"uptime"`
	Time     string `json:"time"`
	Hostname string `json:"hostname"`
	Name     string `json:"name"`
}

func Healthz() *HealthzRespone {
	return &HealthzRespone{
		Uptime:   time.Since(start).String(),
		Time:     time.Now().Format(time.RFC3339Nano),
		Hostname: Hostname,
		Name:     Name,
	}
}
