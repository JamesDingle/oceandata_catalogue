package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
	"unicode"
	"strings"
)

type data struct {
	filename string
	link     string
	mod_time time.Time
	size     int
}

func is_all_whitespace(s string) (ok bool) {
	for _, char := range s {
		if !unicode.IsSpace(char) {
			return false
		}
	}
	return true
}

func check_header(n *html.Node) (arr []string, ok bool) {

	for c := n.FirstChild; c != nil; c = c.NextSibling {

		if c.Type == html.TextNode {
			if !is_all_whitespace(c.Data) {
				arr = append(arr, c.Data)
				ok = true
			}
		} else {
			vals, test := check_header(c)
			if test {
				arr = append(arr, vals...)
				ok = test

			}

		}
	}

	return
}

func f_table(n *html.Node, f string) (arr []string, ok bool) {

	for c := n.FirstChild; c != nil; c = c.NextSibling {

		if c.Data == f {

			if c.Type == html.ElementNode {
				arr, ok = check_header(c)
				return
			}

		} else {
			arr, ok = f_table(c, f)
			if ok {
				return
			}
		}

	}

	return
}

func f_link(n *html.Node) (links map[string]string, ok bool) {

	links = make(map[string]string)

	if n.Type == html.ElementNode && n.Data == "a" {

		if len(n.Attr) > 0 {
			if n.Attr[0].Key == "href" {
				link := n.Attr[0].Val
				base := path.Base(link)
				links[base] = link
				ok = true
			}
		}

	} else {

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			n_links, nok := f_link(c)
			if nok {
				for key, val := range n_links {
					links[key] = val
				}
				ok = true
			}
		}
	}

	return
}

func dateFromString(s string) (t time.Time) {
	layout := "2006-01-02 15:04:05"
	t, err := time.Parse(layout, s)

	if err != nil {
		fmt.Println(err)
	}
	return t
}

func printData(d data) {
	fmt.Println("Filename: ", d.filename, " Link: ", d.link, " Last modified: ", d.mod_time, " Size: ", d.size)
}

func handle(url string) {
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("ERROR: Failed to fetch \"" + url + "\"")
		return
	}

	doc, err := html.Parse(resp.Body)

	if err != nil {
		fmt.Println("ERROR: Failed to parse \"" + url + "\"")
		return
	}

	var header []string
	var h_ok bool

	header, h_ok = f_table(doc, "thead")
	fmt.Println("Returned header: ", header, " Checked: ", h_ok)

	h_len := len(header)

	if !h_ok || h_len == 0 {
		fmt.Println("Exiting due to no table header found")
		os.Exit(2)
	}

	vals, v_ok := f_table(doc, "tbody")

	v_len := len(vals)

	if !v_ok || v_len == 0 {
		fmt.Println("Exiting due to no table header found")
		os.Exit(2)
	}
	var results []data

	for i := 0; i < v_len; i += 3 {
		var item data
		item.filename = vals[i]
		item.mod_time = dateFromString(vals[i+1])
		item.size, _ = strconv.Atoi(vals[i+2])
		results = append(results, item)
	}

	linkmap, linkok := f_link(doc)

	for i, val := range results {
		if linkok {
			if link, ok := linkmap[val.filename]; ok {
				if strings.HasPrefix(link, "/") {
					link = path.Join(url, path.Base(link))
				}
				val.link = link
				results[i] = val
			}
		}
		printData(val)
	}

}

func main() {

	urls := os.Args[1:]
	for _, url := range urls {
		fmt.Println("Handling URL: ", url)
		handle(url)
	}

}
