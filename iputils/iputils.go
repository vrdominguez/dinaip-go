package iputils

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
)

var ipObtainServices = []string{
	"http://dinadns01.dinaserver.com",
	"http://dinadns02.dinaserver.com",
	"https://api.ipify.org",
	"https://checkip.amazonaws.com",
	"http://ipinfo.io/ip",
	"http://ifconfig.me/ip",
}

func GetPublicIp() (string, error) {
	for _, service := range ipObtainServices {
		parsedURL, err := url.Parse(service)
		if err != nil {
			continue // We cannot parse de url we want to use
		}

		response, err := http.Get(parsedURL.String())
		if err != nil {
			continue // We cannot get ip from service, try next
		}
		defer response.Body.Close()

		ip, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Error reading response...")
			continue // We cannot get ip from service, try next
		}

		if trimedIP := strings.TrimRight(string(ip), "\r\n"); trimedIP != "" {
			return trimedIP, nil
		}
	}

	return "", fmt.Errorf("unable to retrieve public IP from any service")
}

func FqdnToIP(domain string) (string, error) {
	ip := net.ParseIP(domain)
	if ip != nil {
		return ip.String(), nil
	}

	ips, err := net.LookupIP(domain)
	if err != nil {
		return "", err
	}

	for _, ip := range ips {
		if ip.To4() != nil {
			return ip.String(), nil
		}
	}

	cname, err := net.LookupCNAME(domain)
	if err != nil {
		return "", err
	}

	return FqdnToIP(cname)
}
