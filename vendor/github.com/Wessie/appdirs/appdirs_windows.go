package appdirs

import (
	"path/filepath"
	"syscall"
	"unsafe"
)

var (
	shell32, _            = syscall.LoadLibrary("shell32.dll")
	getKnownFolderPath, _ = syscall.GetProcAddress(shell32, "SHGetKnownFolderPath")

	ole32, _         = syscall.LoadLibrary("Ole32.dll")
	coTaskMemFree, _ = syscall.GetProcAddress(ole32, "CoTaskMemFree")
)

// These are KNOWNFOLDERID constants that are passed to GetKnownFolderPath
var (
	rfidLocalAppData = syscall.GUID{
		0xf1b32785,
		0x6fba,
		0x4fcf,
		[8]byte{0x9d, 0x55, 0x7b, 0x8e, 0x7f, 0x15, 0x70, 0x91},
	}
	rfidRoamingAppData = syscall.GUID{
		0x3eb685db,
		0x65f9,
		0x4cf6,
		[8]byte{0xa0, 0x3a, 0xe3, 0xef, 0x65, 0x72, 0x9f, 0x3d},
	}
	rfidProgramData = syscall.GUID{
		0x62ab5d82,
		0xfdc1,
		0x4dc3,
		[8]byte{0xa9, 0xdd, 0x07, 0x0d, 0x1d, 0x49, 0x5d, 0x97},
	}
)

func userDataDir(name, author, version string, roaming bool) (path string) {
	if author == "" {
		author = name
	}

	var rfid syscall.GUID
	if roaming {
		rfid = rfidRoamingAppData
	} else {
		rfid = rfidLocalAppData
	}

	path, err := getFolderPath(rfid)

	if err != nil {
		return ""
	}

	if path, err = filepath.Abs(path); err != nil {
		return ""
	}

	if name != "" {
		path = filepath.Join(path, author, name)
	}

	if name != "" && version != "" {
		path = filepath.Join(path, version)
	}

	return path
}

func siteDataDir(name, author, version string) (path string) {
	path, err := getFolderPath(rfidProgramData)

	if err != nil {
		return ""
	}

	if path, err = filepath.Abs(path); err != nil {
		return ""
	}

	if author == "" {
		author = name
	}

	if name != "" {
		path = filepath.Join(path, author, name)
	}

	if name != "" && version != "" {
		path = filepath.Join(path, version)
	}

	return path
}

func userConfigDir(name, author, version string, roaming bool) string {
	return UserDataDir(name, author, version, roaming)
}

func siteConfigDir(name, author, version string) (path string) {
	return SiteDataDir(name, author, version)
}

func userCacheDir(name, author, version string, opinion bool) (path string) {
	if author == "" {
		author = name
	}

	path, err := getFolderPath(rfidLocalAppData)

	if err != nil {
		return ""
	}

	if path, err = filepath.Abs(path); err != nil {
		return ""
	}

	if name != "" {
		path = filepath.Join(path, author, name)
		if opinion {
			path = filepath.Join(path, "Cache")
		}
	}

	if name != "" && version != "" {
		path = filepath.Join(path, version)
	}

	return path
}

func userLogDir(name, author, version string, opinion bool) (path string) {
	path = UserDataDir(name, author, version, false)

	if opinion {
		path = filepath.Join(path, "Logs")
	}

	return path
}

func getFolderPath(rfid syscall.GUID) (string, error) {
	var res uintptr

	ret, _, callErr := syscall.Syscall6(
		uintptr(getKnownFolderPath),
		4,
		uintptr(unsafe.Pointer(&rfid)),
		0,
		0,
		uintptr(unsafe.Pointer(&res)),
		0,
		0,
	)

	if callErr != 0 && ret != 0 {
		return "", callErr
	}

	defer syscall.Syscall(uintptr(coTaskMemFree), 1, res, 0, 0)
	return ucs2PtrToString(res), nil
}

func ucs2PtrToString(p uintptr) string {
	ptr := (*[4096]uint16)(unsafe.Pointer(p))

	return syscall.UTF16ToString((*ptr)[:])
}
