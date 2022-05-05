package elevation

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"encoding/json"
	"io/ioutil"

	"github.com/HydrologicEngineeringCenter/shape-sql-loader/internal/global"
)

type Query struct {
	u      *url.URL
	params map[string]string
}

func NewNationalMapQuery() Query {
	u := url.URL{
		Scheme: global.NATIONAL_MAP_SCHEME,
		Host:   global.NATIONAL_MAP_HOST,
		Path:   global.NATIONAL_MAP_PATH,
	}
	qb := Query{
		u:      &u,
		params: make(map[string]string),
	}
	qb.setParam("dataset", global.NATIONAL_MAP_DATASET)
	qb.setParam("prodFormats", "GeoTIFF")
	return qb
}

func (q *Query) setParam(k string, v string) {
	q.params[k] = v
}

func (q *Query) delParam(k string) {
	delete(q.params, k)
}

func (q *Query) setParams(kv map[string]string) {
	for k, v := range kv {
		q.setParam(k, v)
	}
}

func (q *Query) String() string {
	s := q.u.String() + "?"
	p := url.Values{}
	for k, v := range q.params {
		p.Add(k, v)
	}
	return s + p.Encode()
}

func (q *Query) QueryName(n string) (QueryResult, error) {
	q.setParam("q", n)
	r, err := q.sendRequest()
	return r, err
}

func (q *Query) QueryBoundingBox(b BoundingBox) (QueryResult, error) {
	q.setParam("bbox", fmt.Sprintf("%f,%f,%f,%f", b.MinX, b.MinY, b.MaxX, b.MaxY))
	r, err := q.sendRequest()
	return r, err
}

// newQueryResult deserializes json bytes into a QueryResult struct
func (q *Query) sendRequest() (QueryResult, error) {
	u := q.String()
	resp, err := http.Get(u)
	if err != nil {
		return QueryResult{}, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return QueryResult{}, err
	}
	var r QueryResult
	err = json.Unmarshal(b, &r)
	if err != nil {
		return QueryResult{}, err
	}
	return r, nil
}

type QueryResult struct {
	Total            int      `json:"total"`
	Items            []Item   `json:"items"`
	Error            []string `json:"errors"`
	Messages         []string `json:"messages"`
	SciencebaseQuery string   `json:"sciencebaseQuery"`
	FilteredOut      int      `json:"filteredOut"`
}

type Item struct {
	Title             string      `json:"title"`
	MoreInfo          string      `json:"moreInfo"`
	SourceID          string      `json:"sourceId"`
	SourceName        string      `json:"sourceName"`
	SourceOriginID    string      `json:"sourceOriginId"`
	SourceOriginName  string      `json:"sourceOriginName"`
	MetaURL           string      `json:"metaUrl"`
	VendorMetaURL     string      `json:"vendorMetaUrl"`
	PublicationDate   string      `json:"publicationDate"`
	LastUpdated       string      `json:"lastUpdated"`
	DateCreated       string      `json:"dateCreated"`
	SizeInBytes       int         `json:"sizeInBytes"`
	Extent            string      `json:"extent"`
	Format            string      `json:"format"`
	DownloadURL       string      `json:"downloadURL"`
	DownloadURLRaster string      `json:"downloadURLRaster"`
	PreviewGraphicURL string      `json:"previewGraphicURL"`
	DownloadLazURL    string      `json:"downloadLazURL"`
	Urls              Urls        `json:"urls"`
	Datasets          []string    `json:"datasets"`
	BoundingBox       BoundingBox `json:"boundingBox"`
	BestFitIndex      float32     `json:"bestFitIndex"`
	Body              string      `json:"body"`
	ProcessingURL     string      `json:"processingUrl"`
	ModificationInfo  string      `json:"modificationInfo"`
}

type Urls struct {
	Tiff string `json:"TIFF"`
}

// cacheKey generates a key to the data file within the key/value store
func (i Item) cacheKey() (string, error) {
	urlTokens := strings.Split(i.DownloadURL, "/")
	return url.QueryUnescape(fmt.Sprintf(`%s%s`, global.NATIONAL_MAP_CACHE_BASEPATH, urlTokens[len(urlTokens)-1]))
}
