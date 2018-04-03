package main

import (
	"net/http"
	"fmt"
	"golang.org/x/net/html"
	"unicode"
	"os"
	"time"
	//strftime "github.com/jehiah/go-strftime"
)

type data struct {
	filename string
	link string
	mod_time time.Time
	size int
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
				return arr, ok
			}
		} else {
			vals, test := check_header(c)
			if test {
				arr = append(arr, vals...)
				ok = test
				return arr, ok

			}

		}
	}

	return
}

func f_table(n *html.Node, f string) (arr []string, ok bool) {


	//fmt.Println("Searching table under: ", f)

	for c := n.FirstChild; c != nil; c = c.NextSibling {

		if c.Data == f {
			fmt.Println("Found node: ", f)

			if c.Type == html.ElementNode {
				fmt.Println("Is an element node...")
			} else {
				fmt.Println(c.Type)
			}

		}

		if c.Type == html.ElementNode && c.Data == f {
			h_arr, h_ok := check_header(c)
			if h_ok && f == "tbody" {
				fmt.Println("----", arr[0])
				link, check := f_link(c)
				if check {
					fmt.Println("Found link: ", link)
					h_arr = append(h_arr, link)
				}
				arr = append(arr, h_arr...)
				ok = true
			}
			//if ok { return }
		} else {
			arr, ok = f_table(c,f)
			if ok {
				//fmt.Println(f_link(c.FirstChild))
				return
			}
		}
	}

	return
}

func f_link(n *html.Node) (s string, ok bool) {
	if n.Type == html.ElementNode && n.Data == "a" {
		s = n.Attr[0].Val
		ok = true
		return
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		s, ok = f_link(c)
		if ok { return }
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



	header, h_ok := f_table(doc, "thead")

	h_len := len(header)

	fmt.Println("Header elements: ", h_len)

	if !h_ok || h_len == 0  {
		fmt.Println("Exiting due to no table header found")
		os.Exit(2)
	}

	for i, val := range(header) {
		fmt.Println(i, val)
	}


	//vals , v_ok := f_table(doc, "tbody")
	//fmt.Println(vals)
	//
	////fmt.Println(vals)
	//v_len := len(vals)
	//
	//if !v_ok || v_len == 0  {
	//	fmt.Println("Exiting due to no table header found")
	//	os.Exit(2)
	//}
	//var results []data
	//
	//for i := 0; i < v_len; i+=3 {
	//	var item data
	//	fmt.Println(i, v_len)
	//	item.filename = vals[i]
	//	item.mod_time = dateFromString(vals[i+1])
	//	item.size, _ = strconv.Atoi(vals[i+2])
	//	//item.link = vals[i+3]
	//	results = append(results, item)
	//}
	//
	//for _, val := range(results) {
	//	printData(val)
	//}


}

func main() {

	url := "https://oceandata.sci.gsfc.nasa.gov/Ancillary/LUTs/modis/"
	handle(url)
	//
	////  //Kick off the handle process (concurrently)
	//urls := os.Args[1:]
	//for _, url := range urls {
	//	fmt.Println("Handling URL: ", url)
	//	handle(url)
	//}

}