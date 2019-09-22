package maps

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/pkging"
)

// Paths wraps sync.Map and uses the following types:
// key:   string
// value: Path
type Paths struct {
	Current here.Info
	data    *sync.Map
	once    *sync.Once
}

func (m *Paths) Data() *sync.Map {
	if m.once == nil {
		m.once = &sync.Once{}
	}
	m.once.Do(func() {
		if m.data == nil {
			m.data = &sync.Map{}
		}
	})
	return m.data
}

func (m *Paths) MarshalJSON() ([]byte, error) {
	mm := map[string]interface{}{}
	m.Data().Range(func(key, value interface{}) bool {
		mm[fmt.Sprintf("%s", key)] = value
		return true
	})
	return json.Marshal(mm)
}

func (m *Paths) UnmarshalJSON(b []byte) error {
	mm := map[string]pkging.Path{}

	if err := json.Unmarshal(b, &mm); err != nil {
		return err
	}
	for k, v := range mm {
		m.Store(k, v)
	}
	return nil
}

// Delete the key from the map
func (m *Paths) Delete(key string) {
	m.Data().Delete(key)
}

// Load the key from the map.
// Returns Path or bool.
// A false return indicates either the key was not found
// or the value is not of type Path
func (m *Paths) Load(key string) (pkging.Path, bool) {
	i, ok := m.Data().Load(key)
	if !ok {
		return pkging.Path{}, false
	}
	s, ok := i.(pkging.Path)
	return s, ok
}

// Range over the Path values in the map
func (m *Paths) Range(f func(key string, value pkging.Path) bool) {
	m.Data().Range(func(k, v interface{}) bool {
		key, ok := k.(string)
		if !ok {
			return false
		}
		value, ok := v.(pkging.Path)
		if !ok {
			return false
		}
		return f(key, value)
	})
}

// Store a Path in the map
func (m *Paths) Store(key string, value pkging.Path) {
	m.Data().Store(key, value)
}

// Keys returns a list of keys in the map
func (m *Paths) Keys() []string {
	var keys []string
	m.Range(func(key string, value pkging.Path) bool {
		keys = append(keys, key)
		return true
	})
	sort.Strings(keys)
	return keys
}

func (m *Paths) Parse(p string) (pkging.Path, error) {
	p = strings.TrimSpace(p)
	p = filepath.Clean(p)
	p = strings.TrimPrefix(p, m.Current.Dir)

	p = strings.Replace(p, "\\", "/", -1)
	p = strings.TrimSpace(p)

	pt, ok := m.Load(p)
	if ok {
		return pt, nil
	}
	if len(p) == 0 || p == ":" {
		return m.build("", "", "")
	}

	res := pathrx.FindAllStringSubmatch(p, -1)
	if len(res) == 0 {
		return pt, fmt.Errorf("could not parse %q", p)
	}

	matches := res[0]

	if len(matches) != 4 {
		return pt, fmt.Errorf("could not parse %q", p)
	}

	return m.build(p, matches[1], matches[3])
}

var pathrx = regexp.MustCompile("([^:]+)(:(/.+))?")

func (m *Paths) build(p, pkg, name string) (pkging.Path, error) {
	pt := pkging.Path{
		Pkg:  pkg,
		Name: name,
	}

	if strings.HasPrefix(pt.Pkg, "/") || len(pt.Pkg) == 0 {
		pt.Name = pt.Pkg
		pt.Pkg = m.Current.ImportPath
	}

	if len(pt.Name) == 0 {
		pt.Name = "/"
	}

	if pt.Pkg == pt.Name {
		pt.Pkg = m.Current.ImportPath
		pt.Name = "/"
	}

	if !strings.HasPrefix(pt.Name, "/") {
		pt.Name = "/" + pt.Name
	}
	pt.Name = strings.TrimPrefix(pt.Name, m.Current.Dir)
	m.Store(p, pt)
	return pt, nil
}