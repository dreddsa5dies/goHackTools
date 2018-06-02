go-crypt (`crypt`)
==================

[![Build Status](https://travis-ci.org/amoghe/go-crypt.svg)](https://travis-ci.org/amoghe/go-crypt)

Package `crypt` provides go language wrappers around crypt(3). For further information on crypt see the
[man page](http://man7.org/linux/man-pages/man3/crypt.3.html)

If you have questions about how to use crypt (the C function), it is likely this is not the package you
are looking for.

**NOTE** Depending on the platform, this package provides a `Crypt` function that is backed by different
flavors of the libc crypt. This is done by detecting the GOOS and trying to build using `crypt_r` (the GNU
extension) when on linux, and wrapping around plain 'ol `crypt` (guarded by a global lock) otherwise.

Example
-------
```go
import (
	"fmt"
	"github.com/amoghe/go-crypt"
)

func main() {
	md5, err := crypt.Crypt("password", "in")
	if err != nil {
		fmt.Errorf("error:", err)
		return
	}

	sha512, err := crypt.Crypt("password", "$6$SomeSaltSomePepper$")
	if err != nil {
		fmt.Errorf("error:", err)
		return
	}

	fmt.Println("MD5:", md5)
	fmt.Println("SHA512:", sha512)
}
```

A Note On "Salt"
----------------

You can find out more about salt [here](https://en.wikipedia.org/wiki/Salt_(cryptography))

The hash algorithm can be selected via the salt string. Here is how to do it (relevant
section from the man page):

```
   If salt is a character string starting with the characters
   "$id$" followed by a string terminated by "$":

       $id$salt$encrypted

   then instead of using the DES machine, id identifies the
   encryption method used and this then determines how the rest
   of the password string is interpreted.  The following values
   of id are supported:

          ID  | Method
          ─────────────────────────────────────────────────────────
          1   | MD5
          2a  | Blowfish (not in mainline glibc; added in some
              | Linux distributions)
          5   | SHA-256 (since glibc 2.7)
          6   | SHA-512 (since glibc 2.7)

   So $5$salt$encrypted is an SHA-256 encoded password and
   $6$salt$encrypted is an SHA-512 encoded one.

   "salt" stands for the up to 16 characters following "$id$" in
   the salt.  The encrypted part of the password string is the
   actual computed password.  The size of this string is fixed:

   MD5     | 22 characters
   SHA-256 | 43 characters
   SHA-512 | 86 characters
```

Platforms
---------

This package has been tested on the following platforms:
- ubuntu 14.04.2 (libc 2.19)
- ubuntu 12.04.5 (libc 2.15)
- centos         (libc 2.17)
- fedora 22      (libc 2.21)

All the platforms tested on have GNU libc (with extensions) so that the GOOS=linux always
compiles the reentrant versions of the crypt function (`crypt_r`), and exposes it to go land.

Other platforms (freebsd, netbsd) should also work (in theory) since their libc expose at least
a posix compliant crypt function. On these platforms the fallback should compile and expose the
'plain' (non reentrant, thus globally locked) crypt function.

Unfortunately, I do not have access to machines that run anything other than Linux, hence the other
platforms have not been tested, however I believe they should work just fine. If you can verify this
(or provide a patch that fixes this), I would be grateful.

TODO
----
* Find someone with access to *BSD system(s)

License
-------

Released under the [MIT License](LICENSE)
