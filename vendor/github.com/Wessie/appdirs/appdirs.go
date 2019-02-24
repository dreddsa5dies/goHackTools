// A port of the excellent python module `appdirs`.
// See https://github.com/ActiveState/appdirs for the python version.
package appdirs

import (
	"os/user"
	"strings"
)

// App is a helper type to create easy access across your program to the appdirs
// functions.
//
// The *App type has 6 methods that map to the 6 functions exported by `appdirs`.
// All methods take no arguments, and supply the function it wraps with arguments
// pre-set in the struct on creation.
type App struct {
	Name    string
	Author  string
	Version string
	Roaming bool
	Opinion bool
}

// New returns a new App helper that has various methods for receiving
// relevant directories for your application.
//
// The following defaults are used for the two fields not settable by New:
// Roaming: false, Opinion: true
//
// If you want to set these, create your own App struct by the usual means.
func New(name, author, version string) *App {
	return &App{
		Name:    name,
		Author:  author,
		Version: version,
		Roaming: false,
		Opinion: true,
	}
}

// UserData returns the full path to the user-specific data directory
func (app *App) UserData() string {
	return UserDataDir(app.Name, app.Author, app.Version, app.Roaming)
}

// SiteData returns the full path to the user-shared data directory
func (app *App) SiteData() string {
	return SiteDataDir(app.Name, app.Author, app.Version)
}

// SiteConfig returns the full path to the user-shared configuration directory
func (app *App) SiteConfig() string {
	return SiteConfigDir(app.Name, app.Author, app.Version)
}

// UserConfig returns the full path to the user-specific configuration directory
func (app *App) UserConfig() string {
	return UserConfigDir(app.Name, app.Author, app.Version, app.Roaming)
}

// UserCache returns the full path to the user-specific cache directory
func (app *App) UserCache() string {
	return UserCacheDir(app.Name, app.Author, app.Version, app.Opinion)
}

// UserLog returns the full path to the user-specific log directory
func (app *App) UserLog() string {
	return UserLogDir(app.Name, app.Author, app.Version, app.Opinion)
}

// ExpandUser is a helper function that expands the first '~' it finds in the
// passed path with the home directory of the current user.
//
// Note: This only works on environments similar to bash.
func ExpandUser(path string) string {
	if u, err := user.Current(); err == nil {
		return strings.Replace(path, "~", u.HomeDir, -1)
	}
	return path
}
