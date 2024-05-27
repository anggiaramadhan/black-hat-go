package main

import (
	"archive/zip"
	"bing-metadata/metadata"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func handler(i int, s *goquery.Selection) {
	url, ok := s.Find("a").Attr("href")
	if !ok {
		log.Fatalln("error find href")
		return
	}

	fmt.Printf("%d: %s", i, url)
	res, err := http.Get(url)
	if err != nil {
		log.Fatalln("error get url")
		return
	}

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln("error when read all")
		return
	}
	defer res.Body.Close()

	r, err := zip.NewReader(bytes.NewReader(buf), int64(len(buf)))
	if err != nil {
		log.Fatalln("error when read zip")
		return
	}

	cp, er, err := metadata.NewProperties(r)
	if err != nil {
		log.Fatalln("error reading metadata")
		return
	}

	log.Printf(
		"%25s %25s - %s %s\n",
		er.Creator,
		er.LastModifiedBy,
		cp.Application,
		cp.GetMajorVersion())
}

func main() {
	if len(os.Args) < 3 {
		log.Fatalln("missing command. Usage main.go domain ext")
	}

	domain := os.Args[1]
	filetype := os.Args[2]

	q := fmt.Sprintf(
		"site:%s && filetype:%s && instreamset:(url title):%s",
		domain,
		filetype,
		filetype,
	)

	search := fmt.Sprintf("http://www.bing.com/search?q=%s", url.QueryEscape(q))
	doc, err := goquery.NewDocument(search)
	if err != nil {
		log.Fatalln("document not found")
	}

	s := "html body div#b_content ol#b_results li.b_algo div.b_title h2"
	fmt.Printf("doc size: %d\n", doc.Size())
	html, err := doc.Html()
	if err != nil {
		log.Fatalln("html error")
	}

	f, err := os.OpenFile("result.html", os.O_CREATE|os.O_RDWR, 0664)
	if err != nil {
		log.Fatalln("error open file")
	}
	defer f.Close()

	n, err := f.Write([]byte(html))
	if n == 0 {
		log.Fatalln("error write file")
	}

	if err != nil {
		log.Fatalln("error write html file")
	}

	doc.Find(s).Each(handler)
}
