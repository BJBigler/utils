package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const (
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36"
)

//GetGoQueryDocumentFromURL returns goquery doc from URL
func GetGoQueryDocumentFromURL(url string) *goquery.Document {

	timeout := time.Duration(60 * time.Second)
	cookieJar, _ := cookiejar.New(nil)

	client := http.Client{
		Timeout: timeout,
		Jar:     cookieJar,
	}

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/html")
	req.Header.Set("Accept-Language", "en-US;q=0.8")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Close = true

	resp, err := client.Do(req)

	if err != nil {
		Log(err)
		return nil
	}
	defer resp.Body.Close()

	doc, err3 := goquery.NewDocumentFromResponse(resp)

	if err3 != nil {
		Log(err3)
		return nil
	}

	return doc
}

//GetGoQueryDocumentFromFile returns goquery doc from URL
func GetGoQueryDocumentFromFile(path string) (*goquery.Document, error) {

	r, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(r)

	if err != nil {
		return nil, err
	}

	return doc, nil
}

//GetContent returns contents of a URL
func GetContent(sourceURL string) ([]byte, error) {

	parsedURL, err := url.Parse(sourceURL)
	if err != nil {
		return nil, err
	}

	parsedURL.RawQuery = parsedURL.Query().Encode()

	resp, err := http.Get(parsedURL.String())

	if err != nil {
		return nil, fmt.Errorf("GET error: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GetContent() Status error: %v from URL: %v", resp.StatusCode, sourceURL)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Read body: %v", err)
	}

	return data, nil
}

//GetImagesFromURL returns a slice of image urls, e.g.
//http://www.remotesite.com/img.jpg
func GetImagesFromURL(source string) []string {
	result := []string{}
	//Check URL for errors
	parseURL, err := url.Parse(source)
	if err != nil {
		return result
	}

	doc := GetGoQueryDocumentFromURL(parseURL.RawPath)

	doc.Find("img").Each(func(i int, img *goquery.Selection) {
		val, exists := img.Attr("src")

		if exists {
			imgURL, err := url.Parse(val)
			if err == nil {
				u := parseURL.ResolveReference(imgURL)
				result = append(result, u.RawPath)
			}
		}
	})

	return result
}

//PrepRedirectURL takes an http.Request to create a redirect
//URL back to it
func PrepRedirectURL(r *http.Request) *url.URL {

	result := r.URL
	if result.Host == "" {
		result.Host = r.Host
	}

	result.Scheme = "http"

	if scheme := r.Header.Get("X-Forwarded-Proto"); scheme != "" {
		result.Scheme = scheme
	} else if r.TLS != nil {
		result.Scheme = "https"
	}

	return result
}

//Join takes *basePath* ("http...") and joins it with
//paths to form a URL.
func Join(basePath string, paths ...string) (*url.URL, error) {
	u, err := url.Parse(basePath)
	if err != nil {
		return nil, fmt.Errorf("invalid url")
	}

	p2 := append([]string{u.Path}, paths...)

	result := path.Join(p2...)

	u.Path = result
	return u, nil
}

//Secure returns false only if the Host is localhost
func Secure(r *http.Request) bool {
	if strings.Contains(r.Host, "localhost") {
		return false
	}
	return true

}

//GetIP returns UserIP address
func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}
