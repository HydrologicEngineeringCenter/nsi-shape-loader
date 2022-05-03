package elevation

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/HydrologicEngineeringCenter/shape-sql-loader/internal/global"
	"github.com/usace/filestore"
)

type cacheItem filestore.FileStoreResultObject

// TODO maybe spin this out into a general lib, seems useful

// ElevationAccessor acts as a caching service around the National Map API and
// the local filestore. The accessor query the National Map service for a list
// of relevant files based on a BoundingBox generated from a set of Points.
// If the files are already available in the localCache, then it uses that data.
// Otherwise, the accessor downloads to localCache before loading.
type ElevationAccessor struct {
	queryResult QueryResult
	localCache  filestore.FileStore
	cacheObjs   *[]cacheItem // this wrangling is so convoluted TODO maybe refactor
}

func NewElevationAccessor(p Points) (ElevationAccessor, error) {
	b := p.BoundingBox()
	q, err := b.QueryNationalMap()
	if err != nil {
		return ElevationAccessor{}, err
	}
	localFS, err := filestore.NewFileStore(filestore.BlockFSConfig{})
	if err != nil {
		return ElevationAccessor{}, err
	}
	e := ElevationAccessor{
		queryResult: q,
		localCache:  localFS,
	}
	err = e.refreshCacheObjs()
	if err != nil {
		return ElevationAccessor{}, err
	}
	return e, nil
}

// GetElevation fills the nil Elevation field for each point
func (e *ElevationAccessor) GetElevation(p Points) error {
	errs := make(chan error, 1)
	// loop through all available items and download relevant file to local cache
	for _, i := range e.queryResult.Items {
		// filter for only USGS 1/3 arc-second dataset
		if strings.Contains(i.Title, "USGS 13 arc-second") {
			existsInCache, err := e.cacheContains(i)
			if err != nil {
				return err
			}
			if p.IsIntersecting(i) && !existsInCache {
				// TODO might be something blocking the multithreading here, not seeing the mutiple threads spawning, could be API rate/concurrency limit
				go func() {
					errs <- e.downloadData(i)
				}()
				if err := <-errs; err != nil {
					return err
				}
			}
		}
	}
	// loop through all cachedItem TIFF
	for _, cacheItem := range *e.cacheObjs {
		i, err := e.getItemFromCacheItem(cacheItem)
		if err != nil {
			return err
		}
		// intersect points relevant for each cacheItem TIFF
		intersectedPoints := i.BoundingBox.Intersect(p)
		cachedKey, err := i.cacheKey()
		if err != nil {
			return err
		}
		g, err := newGDALAccessor(cachedKey)
		if err != nil {
			return err
		}
		// populate elevation data for each point
		for _, point := range intersectedPoints {
			if point.NilElevation() {
				// boxed pointer - a trick from rust
				err = g.calculateElevation(i.BoundingBox, *point)
				if err != nil {
					return err
				}
				boxed := float64(0)
				point.Elevation = &boxed // TODO setting to 0 for testing
			}
		}
	}
	return nil
}

// getItemFromCacheItem finds the corresponding Item obj from cacheItem
func (e *ElevationAccessor) getItemFromCacheItem(c cacheItem) (Item, error) {
	for _, i := range e.queryResult.Items {
		cachedKey, err := i.cacheKey()
		if err != nil {
			return Item{}, err
		}
		if cachedKey == (c.Path + c.Name) {
			return i, nil
		}
	}
	return Item{}, errors.New(fmt.Sprintf("No QueryResult Item located at: %s", c.Path))
}

// refreshCacheObjs keeps an index cache in memory during app lifetime
// rather than querying everytime there's a need and slowing down performance
func (e *ElevationAccessor) refreshCacheObjs() error {
	o, err := e.localCache.GetDir(global.NATIONAL_MAP_CACHE_BASEPATH, false)
	if err != nil {
		return err
	}
	var flush []cacheItem
	for _, i := range *o {
		// coerce to the new alias type
		coerced := cacheItem(i)
		flush = append(flush, coerced)
	}
	e.cacheObjs = &flush
	if err != nil {
		return err
	}
	return err
}

// downloadData sends out a get request to the National Map API
// and download data to localCache
func (e *ElevationAccessor) downloadData(i Item) error {
	cachedKey, err := i.cacheKey()
	if err != nil {
		return err
	}
	out, err := os.Create(cachedKey)
	if err != nil {
		return err
	}
	defer out.Close()
	resp, err := http.Get(i.DownloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	err = e.refreshCacheObjs()
	if err != nil {
		return err
	}
	return nil
}

// cacheContains returns true if item is already downloaded and in local cache
func (e *ElevationAccessor) cacheContains(i Item) (bool, error) {
	cachedKey, err := i.cacheKey()
	if err != nil {
		return false, err
	}
	if _, err := os.Stat(cachedKey); errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return true, nil
}
