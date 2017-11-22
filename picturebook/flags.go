package picturebook

import (
	"strings"
)

// please make these precompile golang regexp thingies
// (20171122/thisisaaronland)

type IncludeFlag []string

func (i *IncludeFlag) String() string {
	return strings.Join(*i, "\n")
}

func (i *IncludeFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type ExcludeFlag []string

func (e *ExcludeFlag) String() string {
	return strings.Join(*e, "\n")
}

func (e *ExcludeFlag) Set(value string) error {
	*e = append(*e, value)
	return nil
}
