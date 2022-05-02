package elevation

import (
	"fmt"
	"net/url"
	"strings"

	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/HydrologicEngineeringCenter/shape-sql-loader/internal/global"
)

type QueryBuilder struct {
	u *url.URL
}

func NewNationalMapQuery() QueryBuilder {
	u := url.URL{
		Scheme: global.NATIONAL_MAP_SCHEME,
		Host:   global.NATIONAL_MAP_HOST,
		Path:   global.NATIONAL_MAP_PATH,
	}
	qb := QueryBuilder{
		u: &u,
	}
	qb.setParam("dataset", global.NATIONAL_MAP_DATASET)
	return qb
}

func (q *QueryBuilder) setParam(k string, v string) {
	q.u.Query().Set(k, v)
}

func (q *QueryBuilder) delParam(k string) {
	q.u.Query().Del(k)
}

func (q *QueryBuilder) setParams(kv map[string]string) {
	for k, v := range kv {
		q.setParam(k, v)
	}
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

func newQueryResult(path string) (QueryResult, error) {
	f, err := os.Open(path)
	if err != nil {
		return QueryResult{}, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return QueryResult{}, err
	}
	var q QueryResult
	err = json.Unmarshal(b, &q)
	if err != nil {
		return QueryResult{}, err
	}
	return q, nil
}

// cacheKey generates a key to the data file within the key/value store
func (i Item) cacheKey() string {
	urlTokens := strings.Split(i.DownloadURL, "/")
	return fmt.Sprintf(`%s/%s`, global.NATIONAL_MAP_CACHE_BASEPATH, urlTokens[len(urlTokens)-1])
}
