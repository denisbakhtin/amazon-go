package utility

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/microcosm-cc/bluemonday"
)

//AppendIfMissing appends string to slice, if its not there yet
func AppendIfMissing(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

//Parameterize parameterizes string, to make it url friendly
func Parameterize(source string) (result string) {
	result = strings.ToLower(source)
	re := regexp.MustCompile("[^a-z0-9\\-_]+")
	result = re.ReplaceAllString(result, "-")
	re = regexp.MustCompile("\\-{2,}")
	result = re.ReplaceAllString(result, "-")
	re = regexp.MustCompile("^\\-|\\-$")
	result = re.ReplaceAllString(result, "")
	return
}

//Min returns int minimum of a & b
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

//Max returns int maximum of a & b
func Max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

//Pagination stores pagination element
type Pagination struct {
	Class string
	URL   string
	Title string
	Rel   string
}

//CurrentPage retrieves page number query parameter
func CurrentPage(c *gin.Context) int {
	currentPage := 1
	if pageStr := c.Query("page"); pageStr != "" {
		currentPage, _ = strconv.Atoi(pageStr)
	}
	currentPage = int(math.Max(float64(1), float64(currentPage)))
	return currentPage
}

//Paginator creates paginator
func Paginator(currentPage, totalPages int, curURL *url.URL) []Pagination {
	currentPage = int(Max(1, Min(int(currentPage), int(totalPages))))
	queryValues, err := url.ParseQuery(curURL.RawQuery)
	if err != nil {
		queryValues = url.Values{}
	}
	nextID := 5
	lastID := 6
	if totalPages < 3 {
		nextID = 4
		lastID = 5
	}
	//first + last + prev + next + 3 adjusent == 7
	pagination := make([]Pagination, 7)

	if totalPages > 1 {
		//prev links
		if currentPage > 1 {
			newURL := *curURL
			newURL.RawQuery = pageQuery(&queryValues, 1)
			pagination[0] = Pagination{Class: "first_page", URL: newURL.RequestURI(), Title: "First"}
			newURL.RawQuery = pageQuery(&queryValues, currentPage-1)
			pagination[1] = Pagination{Class: "previous_page", URL: newURL.RequestURI(), Title: "Previous", Rel: "prev"}
		} else {
			pagination[0] = Pagination{Class: "first_page disabled", URL: "", Title: "First"}
			pagination[1] = Pagination{Class: "previous_page disabled", URL: "", Title: "Previous"}
		}

		//page numbers
		switch currentPage {
		case 1:
			pagination[2] = Pagination{Class: "active", URL: "", Title: "1"}
			if 2 <= totalPages {
				newURL := *curURL
				newURL.RawQuery = pageQuery(&queryValues, 2)
				pagination[3] = Pagination{Class: "", URL: newURL.RequestURI(), Title: "2"}
			}
			if 3 <= totalPages {
				newURL := *curURL
				newURL.RawQuery = pageQuery(&queryValues, 3)
				pagination[4] = Pagination{Class: "", URL: newURL.RequestURI(), Title: "3"}
			}
		case totalPages:
			if 3 <= totalPages {
				pagination[4] = Pagination{Class: "active", URL: "", Title: fmt.Sprintf("%d", totalPages)}
				newURL := *curURL
				newURL.RawQuery = pageQuery(&queryValues, totalPages-1)
				pagination[3] = Pagination{Class: "", URL: newURL.RequestURI(), Title: fmt.Sprintf("%d", totalPages-1)}

				newURL.RawQuery = pageQuery(&queryValues, totalPages-2)
				pagination[2] = Pagination{Class: "", URL: newURL.RequestURI(), Title: fmt.Sprintf("%d", totalPages-2)}
			}
			if 2 == totalPages {
				pagination[3] = Pagination{Class: "active", URL: "", Title: fmt.Sprintf("%d", totalPages)}

				newURL := *curURL
				newURL.RawQuery = pageQuery(&queryValues, totalPages-1)
				pagination[2] = Pagination{Class: "", URL: newURL.RequestURI(), Title: fmt.Sprintf("%d", totalPages-1)}
			}

		default:
			newURL := *curURL
			newURL.RawQuery = pageQuery(&queryValues, currentPage+1)
			pagination[4] = Pagination{Class: "", URL: newURL.RequestURI(), Title: fmt.Sprintf("%d", currentPage+1)}
			pagination[3] = Pagination{Class: "active", URL: "", Title: fmt.Sprintf("%d", currentPage)}
			newURL.RawQuery = pageQuery(&queryValues, currentPage-1)
			pagination[2] = Pagination{Class: "", URL: newURL.RequestURI(), Title: fmt.Sprintf("%d", currentPage-1)}
		}

		//next links
		if currentPage < totalPages {
			newURL := *curURL
			newURL.RawQuery = pageQuery(&queryValues, currentPage+1)
			pagination[nextID] = Pagination{Class: "next_page", URL: newURL.RequestURI(), Title: "Next", Rel: "next"}

			newURL.RawQuery = pageQuery(&queryValues, totalPages)
			pagination[lastID] = Pagination{Class: "last_page", URL: newURL.RequestURI(), Title: "Last"}

		} else {
			pagination[lastID] = Pagination{Class: "last_page disabled", URL: "", Title: "Last"}
			pagination[nextID] = Pagination{Class: "next_page disabled", URL: "", Title: "Next"}
		}
		return pagination
	}
	return nil
}

func pageQuery(query *url.Values, page int) string {
	if page > 1 {
		query.Set("page", fmt.Sprintf("%d", page))
	} else {
		query.Del("page")
	}
	return query.Encode()
}

//SubtractArray subtracts one array from another, mainly for keywords meta tags
func SubtractArray(minuend, subtrahend []string) (difference []string) {
	difference = make([]string, 0, len(minuend))
	for i := range minuend {
		found := false
		for j := range subtrahend {
			if minuend[i] == subtrahend[j] {
				found = true
				break
			}
		}
		if !found {
			difference = append(difference, minuend[i])
		}
	}
	return
}

//SubtractUint64Array subtracts one uint64 array from another
func SubtractUint64Array(minuend, subtrahend []uint64) (difference []uint64) {
	difference = make([]uint64, 0, len(minuend))
	for i := range minuend {
		found := false
		for j := range subtrahend {
			if minuend[i] == subtrahend[j] {
				found = true
				break
			}
		}
		if !found {
			difference = append(difference, minuend[i])
		}
	}
	return
}

//GetFileNameFromURL retreives file name from url
func GetFileNameFromURL(url string) string {
	re, err := regexp.Compile("^.*/(.+)$")
	if err != nil {
		return ""
	}
	res := re.FindStringSubmatch(url)
	if len(res) > 1 {
		return res[1] //(.+) match
	}
	return ""
}

//FileExists checks file existence
func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//StringSliceContains check if string slice contains needed str
func StringSliceContains(str string, strs []string) bool {
	for i := range strs {
		if strs[i] == str {
			return true
		}
	}
	return false
}

//AcceptsHTML checks if Request accepts html response
func AcceptsHTML(r *http.Request) bool {
	if len(r.Header["Accept"]) > 0 {
		return strings.Contains(r.Header["Accept"][0], "text/html")
	}
	return false
}

//AcceptsJSON checks if Request accepts json response
func AcceptsJSON(r *http.Request) bool {
	if len(r.Header["Accept"]) > 0 {
		return strings.Contains(r.Header["Accept"][0], "application/json")
	}
	return false
}

//AppendToCache appends string to cache
func AppendToCache(src []string, index int, value string) []string {
	if index < len(src) {
		src[index] = value
		return src
	}
	dst := make([]string, index+1)
	copy(dst, src)
	dst[index] = value
	return dst
}

//NilOrRefID returns nil if ID == 0 or its reference
func NilOrRefID(id uint64) *uint64 {
	if id == 0 {
		return nil
	}
	return &id
}

//ParseUint converts string to uint
func ParseUint(s string) (uint, error) {
	res, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(res), nil
}

//EditorialSanitizer returns a basic html sanitizer for editorial reviews
func EditorialSanitizer() *bluemonday.Policy {
	p := bluemonday.NewPolicy()
	p.AllowStandardAttributes()                                                          //title, id, dir, lang
	p.AllowAttrs("class").Matching(regexp.MustCompile(`[a-zA-Z0-9\:\-_\.]+`)).Globally() //allow class attribute
	p.AllowLists()                                                                       //li, ul, ol
	//p.AllowImages()                                                                      //allow img
	p.AllowElements("h1", "h2", "h3", "h4", "h5", "h6")
	p.AllowElements("br", "div", "hr", "p", "span")
	p.AllowElements("abbr", "acronym", "cite", "code", "dfn", "em", "figcaption", "mark", "s", "samp", "strong", "sub", "sup")
	p.AllowElements("b", "i", "pre", "small", "strike", "tt", "u")
	return p
}

//ReviewsSanitizer returns a basic html sanitizer for customer reviews
func ReviewsSanitizer() *bluemonday.Policy {
	p := bluemonday.NewPolicy()
	p.AllowStandardAttributes()                                                          //title, id, dir, lang
	p.AllowAttrs("class").Matching(regexp.MustCompile(`[a-zA-Z0-9\:\-_\.]+`)).Globally() //allow class attribute
	p.AllowLists()                                                                       //li, ul, ol
	p.AllowImages()                                                                      //allow img
	p.AllowElements("h1", "h2", "h3", "h4", "h5", "h6")
	p.AllowElements("br", "div", "hr", "p", "span")
	p.AllowElements("abbr", "acronym", "cite", "code", "dfn", "em", "figcaption", "mark", "s", "samp", "strong", "sub", "sup")
	p.AllowElements("b", "i", "pre", "small", "strike", "tt", "u")
	p.AllowElements("script")
	return p
}

//GetAmazonURL creates a client that simulates a real web browser to avoid captcha gate, returns an unprocessed response and error
func GetAmazonURL(path string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:63.0) Gecko/20100101 Firefox/63.0")

	return client.Do(req)
}

//UniqueStrings only unique strings from data slice
func UniqueStrings(data []string) []string {
	resmap := make(map[string]bool)
	for i := range data {
		resmap[data[i]] = true
	}
	keys := make([]string, len(resmap))
	i := 0
	for k := range resmap {
		keys[i] = k
		i++
	}
	return keys
}

//SplitCamelWords splits string into camel cased words
func SplitCamelWords(source string) string {
	re := regexp.MustCompile("([A-Z]+)")
	return strings.TrimLeft(re.ReplaceAllString(source, " $1"), " ")
}

//DelURLTag removes tag query param from url
func DelURLTag(amazonURL string) string {
	eventURL, err := url.Parse(amazonURL)
	if err != nil {
		log.Printf("Error parsing feed URL: %s\n", amazonURL)
		return amazonURL
	}
	//delete credentials
	queryValues := eventURL.Query()
	queryValues.Del("tag") //delete tag before requesting page
	eventURL.RawQuery = queryValues.Encode()
	return eventURL.String()
}

//UTCTime converts time t to UTC local time
func UTCTime(t time.Time) time.Time {
	_, offset := time.Now().Zone()
	return t.Add(-time.Duration(offset) * time.Second).UTC()
}

//AmazonTime converts time t to Amazon.com local time (-4 UTC)
func AmazonTime(t time.Time) time.Time {
	_, offset := time.Now().Zone()
	return t.Add(-time.Duration(offset) * time.Second).Add(-time.Duration(4) * time.Hour).UTC()
}

//BeginningOfDay returns beginning of the day
func BeginningOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}
