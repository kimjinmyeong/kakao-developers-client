package local

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// KeywordSearchResult ...
type KeywordSearchResult struct {
	XMLName xml.Name `xml:"result"`
	Meta    struct {
		TotalCount    int        `json:"total_count" xml:"total_count"`
		PageableCount int        `json:"pageable_count" xml:"pageable_count"`
		IsEnd         bool       `json:"is_end" xml:"is_end"`
		SameName      RegionInfo `json:"same_name" xml:"same_name"`
	} `json:"meta" xml:"meta"`
	Documents []Place `json:"documents" xml:"documents"`
}

// KeywordSearchIterator ...
type KeywordSearchIterator struct {
	Query             string
	CategoryGroupCode string
	Format            string
	AuthKey           string
	X                 string
	Y                 string
	Radius            int
	Rect              string
	Page              int
	Size              int
	Sort              string
}

// KeywordSearch ...
func KeywordSearch(query string) *KeywordSearchIterator {
	return &KeywordSearchIterator{
		Query:             url.QueryEscape(strings.TrimSpace(query)),
		CategoryGroupCode: "",
		Format:            "json",
		AuthKey:           "KakaoAK ",
		X:                 "",
		Y:                 "",
		Radius:            0,
		Rect:              "",
		Page:              1,
		Size:              15,
		Sort:              "accuracy",
	}
}

func (k *KeywordSearchIterator) FormatJSON() *KeywordSearchIterator {
	k.Format = "json"
	return k
}

func (k *KeywordSearchIterator) FormatXML() *KeywordSearchIterator {
	k.Format = "xml"
	return k
}

// AuthorizeWith ...
func (k *KeywordSearchIterator) AuthorizeWith(key string) *KeywordSearchIterator {
	k.AuthKey = "KakaoAK " + strings.TrimSpace(key)
	return k
}

// Category sets the category group code of k.
// There are a few available category group codes:
//
// MT1: Large Supermarket
//
// CS2: Convenience Store
//
// PS3: Daycare Center, Kindergarten
//
// SC4: School
//
// AC5: Academic
//
// PK6: Parking
//
// OL7: Gas Station, Charging Station
//
// SW8: Subway Station
//
// CT1: Culture Facility
//
// AG2: Brokerage
//
// PO3: Public Institution
//
// AT4: Tourist Attractions
//
// FD6: Restaurant
//
// CE7: Cafe
//
// HP8: Hospital
//
// PM9: Pharmacy
//
// BK9: Bank
//
// AD5: Accommodation
func (k *KeywordSearchIterator) Category(groupcode string) *KeywordSearchIterator {
	if groupcode == "MT1" || groupcode == "CS2" || groupcode == "PS3" || groupcode == "SC4" || groupcode == "AC5" || groupcode == "PK6" || groupcode == "OL7" || groupcode == "SW8" || groupcode == "CT1" || groupcode == "AG2" || groupcode == "P03" || groupcode == "AT4" || groupcode == "FD6" || groupcode == "CE7" || groupcode == "HP8" || groupcode == "PM9" || groupcode == "BK9" || groupcode == "AD5" {
		k.CategoryGroupCode = groupcode
	}
	return k
}

// WithRadius ...
func (k *KeywordSearchIterator) WithRadius(x, y float64, radius int) *KeywordSearchIterator {
	k.X = strconv.FormatFloat(x, 'f', -1, 64)
	k.Y = strconv.FormatFloat(y, 'f', -1, 64)
	if 0 <= radius && radius <= 20000 {
		k.Radius = radius
	}
	return k
}

// WithRect ...
func (k *KeywordSearchIterator) WithRect(xMin, yMin, xMax, yMax float64) *KeywordSearchIterator {
	k.Rect = strings.Join([]string{strconv.FormatFloat(xMin, 'f', -1, 64),
		strconv.FormatFloat(yMin, 'f', -1, 64),
		strconv.FormatFloat(xMax, 'f', -1, 64),
		strconv.FormatFloat(yMax, 'f', -1, 64)}, ",")
	return k
}

func (k *KeywordSearchIterator) Result(page int) *KeywordSearchIterator {
	if 1 <= page && page <= 45 {
		k.Page = page
	}
	return k
}

func (k *KeywordSearchIterator) Display(size int) *KeywordSearchIterator {
	if 1 <= size && size <= 45 {
		k.Size = size
	}
	return k
}

// SortBy ...
func (k *KeywordSearchIterator) SortBy(typ string) *KeywordSearchIterator {
	if typ == "accuracy" || typ == "distance" {
		k.Sort = typ
	}
	return k
}

// Next returns the keyword search result and proceeds the iterator to the next page.
func (k *KeywordSearchIterator) Next() (res KeywordSearchResult, err error) {
	// at first, send request to the API server
	client := new(http.Client)
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://dapi.kakao.com/v2/local/search/keyword.%s?query=%s&category_group_code=%s&x=%s&y=%s&radius=%d&rect=%s&page=%d&size=%d&sort=%s", k.
		Format, k.Query, k.CategoryGroupCode, k.X, k.Y, k.Radius, k.Rect, k.Page, k.Size, k.Sort), nil)

	if err != nil {
		return
	}
	// don't forget to close the request for concurrent request
	req.Close = true

	// set authorization header
	req.Header.Set("Authorization", k.AuthKey)

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	// don't forget to close the response body
	defer resp.Body.Close()

	if k.Format == "json" {
		if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
			return
		}
	} else if k.Format == "xml" {
		if err = xml.NewDecoder(resp.Body).Decode(&res); err != nil {
			return
		}
	}

	if res.Meta.IsEnd {
		return res, ErrEndPage
	}

	k.Page++

	return
}