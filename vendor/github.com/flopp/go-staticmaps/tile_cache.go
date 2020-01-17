package sm

import (
	"fmt"
	"os"

	"github.com/Wessie/appdirs"
)

// TileCache provides cache information to the tile fetcher
type TileCache interface {
	// Root path to store cached tiles in with no trailing slash.
	Path() string
	// Permission to set when creating missing cache directories.
	Perm() os.FileMode
}

// TileCacheStaticPath provides a static path to the tile fetcher.
type TileCacheStaticPath struct {
	path string
	perm os.FileMode
}

// Path to the cache.
func (c *TileCacheStaticPath) Path() string {
	return c.path
}

// Perm instructs the permission to set when creating missing cache directories.
func (c *TileCacheStaticPath) Perm() os.FileMode {
	return c.perm
}

// NewTileCache stores cache files in a static path.
func NewTileCache(rootPath string, perm os.FileMode) *TileCacheStaticPath {
	return &TileCacheStaticPath{
		path: rootPath,
		perm: perm,
	}
}

// NewTileCacheFromUserCache stores cache files in a user-specific cache directory.
func NewTileCacheFromUserCache(name string, perm os.FileMode) *TileCacheStaticPath {
	app := appdirs.New("go-staticmaps", "flopp.net", "0.1")
	return &TileCacheStaticPath{
		path: fmt.Sprintf("%s/%s", app.UserCache(), name),
		perm: perm,
	}
}
