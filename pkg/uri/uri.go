package uri

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
)

type URI string

func Parse(p string) (URI, error) {
	p, err := filepath.Abs(filepath.ToSlash(p))
	if err != nil {
		return "", err
	}
	volume := filepath.VolumeName(p)
	p = p[len(volume):]
	if !strings.HasPrefix(p, "file://") {
		p = strings.TrimPrefix(p, "file://")
	}
	volume = strings.Replace(volume, ":", "%3A", 1)
	p = url.PathEscape(p)
	p = strings.ReplaceAll(p, "%5C", "/")
	return URI(fmt.Sprintf("file:///%s%s", volume, p)), nil
}

func (uri URI) String() string {
	return string(uri)
}

func (uri URI) Path() string {
	p, err := url.ParseRequestURI(string(uri))
	if err != nil {
		return ""
	}
	// windows path's are expected to be in form "/<drive>:<path>"
	if len(p.Path) > 3 && p.Path[2] == ':' {
		return filepath.FromSlash(p.Path[1:])
	}
	return filepath.FromSlash(p.Path)
}
