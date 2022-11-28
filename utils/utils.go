package utils

import (
	"fmt"
	"net"
	"net/url"
	"strconv"

	model "github.com/chumvan/confdb/models"
)

func ParseUDPAddrsFromUsers(users []model.User) (udpAddrs []net.UDPAddr) {
	fmt.Println("received users")
	for _, u := range users {
		fmt.Printf("full string: %v\n", u)
		url, err := url.Parse(u.EntityUrl.String())
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("user: %s\n", url.User.Username())
		host, portStr, err := net.SplitHostPort(url.Host)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("host: %s\n", host)
		fmt.Printf("port: %s\n", portStr)
		port, err := strconv.Atoi(portStr)
		if err != nil {
			fmt.Println(err)
		}
		userAddr := net.UDPAddr{
			IP:   net.ParseIP(host),
			Port: port,
		}
		udpAddrs = append(udpAddrs, userAddr)
	}
	return
}
