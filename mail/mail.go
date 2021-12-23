package mail

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"net/url"
	"os"
	"strconv"
)

func GetDialer() (d *gomail.Dialer, err error) {
	URL := os.Getenv("MAILER_URL")
	if URL == "" {
		err = fmt.Errorf("no MAILER_URL")
		return
	}
	u, err := url.Parse(URL)
	if err != nil {
		return
	}
	port, err := strconv.Atoi(u.Port())
	if err != nil {
		return
	}
	pass, _ := u.User.Password()
	hostname := u.Hostname()
	username := u.User.Username()

	d = gomail.NewDialer(hostname, port, username, pass)
	return
}
