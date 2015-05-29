package flickr

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const retries = 10

const baseURL = "https://api.flickr.com/services/rest/?"

// Flickr represents a connection
type Flickr struct {
	apiKey    string
	apiSecret string
}

// Image represents an image
type Image struct {
	URL    string
	Width  int
	Height int
}

// Rsp API Wrapper
type Rsp struct {
	Stat   string `xml:"stat,attr"`
	Photos Photos `xml:"photos"`
}

// Photos API Wrapper
type Photos struct {
	Page    int     `xml:"page,attr"`
	Pages   int     `xml:"pages,attr"`
	PerPage int     `xml:"perpage,attr"`
	Total   int     `xml:"total,attr"`
	Photos  []Photo `xml:"photo"`
}

// Photo API Wrapper
type Photo struct {
	ID       int    `xml:"id,attr"`
	Owner    string `xml:"owner,attr"`
	Secret   string `xml:"secret,attr"`
	Server   string `xml:"server,attr"`
	Farm     int    `xml:"farm,attr"`
	Title    string `xml:"title,attr"`
	IsPublic bool   `xml:"ispublic,attr"`
	IsFriend bool   `xml:"isfriend,attr"`
	IsFamily bool   `xml:"isfamily,attr"`
	URLL     string `xml:"url_l,attr"`
	HeightL  int    `xml:"height_l,attr"`
	WidthL   int    `xml:"width_l,attr"`
	URLO     string `xml:"url_o,attr"`
	HeightO  int    `xml:"height_o,attr"`
	WidthO   int    `xml:"width_o,attr"`
}

// New creates a new instance of Flickr
func New(apiKey string, apiSecret string) *Flickr {
	return &Flickr{apiKey, apiSecret}
}

// SearchRandom searches for an image and returns a random one
func (f *Flickr) SearchRandom(query string) (*Image, error) {
	var image Image
	url := f.buildRequest("flickr.photos.search",
		[2]string{"text", query},
		[2]string{"per_page", "200"},
		[2]string{"is_getty", "true"},
		[2]string{"extras", "url_l,url_o,o_dims"})
	var rsp Rsp
	for i := 0; i < retries; i++ {
		err := f.request(url, &rsp)
		if nil != err {
			return nil, err
		}
		if len(rsp.Photos.Photos) > 0 {
			break
		} else {
			fmt.Printf("no images found. retry %v ...\n", i)
		}
	}
	if len(rsp.Photos.Photos) <= 0 {
		return nil, errors.New("no images found after " +
			strconv.Itoa(retries) + " retries\n")
	}
	for i := 0; i < 100; i++ {
		photo := rsp.Photos.Photos[random(0, len(rsp.Photos.Photos))]
		if "" != photo.URLO {
			image.URL = photo.URLO
			image.Width = photo.WidthO
			image.Height = photo.HeightO
			break
		} else if "" != photo.URLL {
			image.URL = photo.URLL
			image.Width = photo.WidthL
			image.Height = photo.HeightL
			break
		}
	}
	return &image, nil
}

func random(min, max int) int {
	rand.Seed(time.Now().Unix())
	max = max - min
	if max <= min {
		return min
	}
	return rand.Intn(max-min) + min
}

func (f *Flickr) request(url string, rsp *Rsp) error {
	res, err := http.Get(url)
	if nil != err {
		return err
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if nil != err {
		return err
	}
	xml.Unmarshal(b, rsp)
	return nil
}

func (f *Flickr) buildRequest(method string, parameter ...[2]string) string {
	p := url.Values{}
	p.Add("method", method)
	p.Add("api_key", f.apiKey)
	for _, param := range parameter {
		p.Add(param[0], param[1])
	}
	return baseURL + p.Encode()
}
