package types

import (
	"strings"

	"go/ast"
)

const (
	suppressPrefix                  = "ginkgo-linter:"
	suppressLengthAssertionWarning  = suppressPrefix + "ignore-len-assert-warning"
	suppressNilAssertionWarning     = suppressPrefix + "ignore-nil-assert-warning"
	suppressErrAssertionWarning     = suppressPrefix + "ignore-err-assert-warning"
	suppressCompareAssertionWarning = suppressPrefix + "ignore-compare-assert-warning"
	suppressAsyncAsertWarning       = suppressPrefix + "ignore-async-assert-warning"
	suppressFocusContainerWarning   = suppressPrefix + "ignore-focus-container-warning"
)

type Config struct {
	SuppressLen     Boolean
	SuppressNil     Boolean
	SuppressErr     Boolean
	SuppressCompare Boolean
	SuppressAsync   Boolean
	ForbidFocus     Boolean
	AllowHaveLen0   Boolean
}

func (s *Config) AllTrue() bool {
	return bool(s.SuppressLen && s.SuppressNil && s.SuppressErr && s.SuppressCompare && s.SuppressAsync && !s.ForbidFocus)
}

func (s *Config) Clone() Config {
	return Config{
		SuppressLen:     s.SuppressLen,
		SuppressNil:     s.SuppressNil,
		SuppressErr:     s.SuppressErr,
		SuppressCompare: s.SuppressCompare,
		SuppressAsync:   s.SuppressAsync,
		ForbidFocus:     s.ForbidFocus,
		AllowHaveLen0:   s.AllowHaveLen0,
	}
}

func (s *Config) UpdateFromComment(commentGroup []*ast.CommentGroup) {
	for _, cmntList := range commentGroup {
		if s.AllTrue() {
			break
		}

		for _, cmnt := range cmntList.List {
			commentLines := strings.Split(cmnt.Text, "\n")
			for _, comment := range commentLines {
				comment = strings.TrimPrefix(comment, "//")
				comment = strings.TrimPrefix(comment, "/*")
				comment = strings.TrimSuffix(comment, "*/")
				comment = strings.TrimSpace(comment)

				switch comment {
				case suppressLengthAssertionWarning:
					s.SuppressLen = true
				case suppressNilAssertionWarning:
					s.SuppressNil = true
				case suppressErrAssertionWarning:
					s.SuppressErr = true
				case suppressCompareAssertionWarning:
					s.SuppressCompare = true
				case suppressAsyncAsertWarning:
					s.SuppressAsync = true
				case suppressFocusContainerWarning:
					s.ForbidFocus = false
				}
			}
		}
	}
}

func (s *Config) UpdateFromFile(cm ast.CommentMap) {

	for key, commentGroup := range cm {
		if s.AllTrue() {
			break
		}

		if _, ok := key.(*ast.GenDecl); ok {
			s.UpdateFromComment(commentGroup)
		}
	}
}
