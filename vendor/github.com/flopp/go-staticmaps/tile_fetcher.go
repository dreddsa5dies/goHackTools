// Copyright 2016, 2017 Florian Pigorsch. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package sm

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg" // to be able to decode jpegs
	_ "image/png"  // to be able to decode pngs
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Wessie/appdirs"
)

// TileFetcher downloads map tile images from a TileProvider
type TileFetcher struct {
	tileProvider *TileProvider
	cacheDir     string
	useCaching   bool
	userAgent    string
}

// NewTileFetcher creates a new Tilefetcher struct
func NewTileFetcher(tileProvider *TileProvider) *TileFetcher {
	t := new(TileFetcher)
	t.tileProvider = tileProvider
	app := appdirs.New("go-staticmaps", "flopp.net", "0.1")
	t.cacheDir = fmt.Sprintf("%s/%s", app.UserCache(), tileProvider.Name)
	t.useCaching = true
	t.userAgent = "Mozilla/5.0+(compatible; go-staticmaps/0.1; https://github.com/flopp/go-staticmaps)"
	return t
}

// SetUserAgent sets the HTTP user agent string used when downloading map tiles
func (t *TileFetcher) SetUserAgent(a string) {
	t.userAgent = a
}

func (t *TileFetcher) url(zoom, x, y int) string {
	shard := ""
	ss := len(t.tileProvider.Shards)
	if len(t.tileProvider.Shards) > 0 {
		shard = t.tileProvider.Shards[(x+y)%ss]
	}
	return t.tileProvider.getURL(shard, zoom, x, y)
}

func (t *TileFetcher) cacheFileName(zoom int, x, y int) string {
	return fmt.Sprintf("%s/%d/%d/%d", t.cacheDir, zoom, x, y)
}

// ToggleCaching enables/disables caching
func (t *TileFetcher) ToggleCaching(enabled bool) {
	t.useCaching = enabled
}

// Fetch download (or retrieves from the cache) a tile image for the specified zoom level and tile coordinates
func (t *TileFetcher) Fetch(zoom, x, y int) (image.Image, error) {
	if t.useCaching {
		fileName := t.cacheFileName(zoom, x, y)
		cachedImg, err := t.loadCache(fileName)
		if err == nil {
			return cachedImg, nil
		}
	}

	url := t.url(zoom, x, y)
	data, err := t.download(url)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	if t.useCaching {
		fileName := t.cacheFileName(zoom, x, y)
		if err := t.storeCache(fileName, data); err != nil {
			log.Printf("Failed to store map tile as '%s': %s", fileName, err)
		}
	}

	return img, nil
}

func (t *TileFetcher) download(url string) ([]byte, error) {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", t.userAgent)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GET %s: %s", url, resp.Status)
	}

	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return contents, nil
}

func (t *TileFetcher) loadCache(fileName string) (image.Image, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func (t *TileFetcher) createCacheDir(path string) error {
	src, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(path, 0777)
		}
		return err
	}
	if src.IsDir() {
		return nil
	}
	return fmt.Errorf("File exists but is not a directory: %s", path)
}

func (t *TileFetcher) storeCache(fileName string, data []byte) error {
	dir, _ := filepath.Split(fileName)

	if err := t.createCacheDir(dir); err != nil {
		return err
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = io.Copy(file, bytes.NewBuffer(data)); err != nil {
		return err
	}

	return nil
}
