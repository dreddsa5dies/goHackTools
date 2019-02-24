// Copyright 2016, 2017 Florian Pigorsch. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package sm

import (
	"strings"
)

// hasPrefix checks if 's' has prefix 'prefix'; returns 'true' and the remainder on success, and 'false', 's' otherwise.
func hasPrefix(s string, prefix string) (bool, string) {
	if strings.HasPrefix(s, prefix) {
		return true, strings.TrimPrefix(s, prefix)
	}
	return false, s
}
