package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

func main() {
	output := make(map[string]interface{})
	host := os.Args[1]
	hostUrl, err := url.Parse(host)
	if err != nil {
		log.Fatalf("invalid host url %s\n", err.Error())
	}
	output["host_url"] = hostUrl.String()
	client := http.Client{
		Timeout: time.Second * 10,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	req, err := http.NewRequest("GET", hostUrl.String(), nil)
	if err != nil {
		log.Println(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println(err)
	}

	sel := doc.Find("script").First()

	re := regexp.MustCompile(`.*DeviceModel = '(.*?)'`)
	model := re.FindStringSubmatch(sel.Text())
	if len(model) > 1 {
		output["model"] = model[1]
	} else {
		output["model"] = "No model"
	}

	lib := doc.Find(`link[href*="lib/"]`)

	src, ok := lib.Attr("href")
	if ok {
		date := strings.Split(src, "/")[2]
		t, err := strconv.ParseInt(date, 10, 64)
		if err != nil {
			panic(err)
		}
		output["lib_date"] = time.Unix(t, 0).String()
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	connection := fmt.Sprintf("wss://%s/ws/cli", hostUrl.Host)

	dialer := ws.Dialer{
		Timeout: time.Second * 10,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	conn, _, _, err := dialer.Dial(ctx, connection)
	data := make([]byte, 0)
	if err != nil {
		if _, ok := err.(ws.StatusError); ok {
			output["hostname"] = "400 Resp"
			output["websocket"] = false
		}
	} else {
		defer conn.Close()
		data, err = wsutil.ReadServerBinary(conn)
		if err != nil {
			log.Println(err)
		}
	}

	re_login := regexp.MustCompile(`\r(.*?)\slogin`)
	hostname := re_login.FindSubmatch(data)
	if hostname != nil {
		if len(hostname[0]) > 1 {
			output["hostname"] = fmt.Sprintf("Host is \"%s\"", hostname[1])
			output["websocket"] = true
		}
	}

	for k, v := range output {
		fmt.Printf("%s: %v\n", k, v)
	}
}
