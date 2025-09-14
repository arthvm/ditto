package git

import "fmt"

type gitArg interface {
	String() string
	isDiffArg()
	isLogArg()
}

type GitOption string

func (o GitOption) String() string { return string(o) }
func (o GitOption) isDiffArg()     {}
func (o GitOption) isLogArg()      {}

func Branches(base string, head string) gitArg {
	return GitOption(fmt.Sprintf("%s..%s", base, head))
}
