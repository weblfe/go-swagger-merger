package swagger

import (
	"encoding/json"
	"github.com/ghodss/yaml"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Merger struct {
	Swagger map[string]any
}

func NewMerger() *Merger {
	merger := new(Merger)
	merger.Swagger = map[string]any{}
	return merger
}

func (m *Merger) AddFile(file string, pattern ...func(string) bool) error {
	stat, err := os.Stat(file)
	if err != nil {
		return err
	}
	if len(pattern) <= 0 {
		pattern = append(pattern, func(s string) bool {
			return strings.HasSuffix(s, "swagger.json")
		})
	}
	if stat.IsDir() {
		return filepath.Walk(file, m.recursive(file, pattern[0]))
	}
	log.Printf("add file %s \n", file)
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err = f.Close()
	}(f)

	content, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	var s1 interface{}
	err = yaml.Unmarshal(content, &s1)
	if err != nil {
		return err
	}

	return m.merge(s1.(map[string]any))
}

func (m *Merger) merge(f map[string]any) error {
	for key, item := range f {
		if i, ok := item.(map[string]interface{}); ok {
			for subKey, subItem := range i {
				if _, ok := m.Swagger[key]; !ok {
					m.Swagger[key] = map[string]any{}
				}

				m.Swagger[key].(map[string]any)[subKey] = subItem
			}
		} else {
			m.Swagger[key] = item
		}
	}
	return nil
}

func (m *Merger) Save(fileName string, beautify ...bool) error {
	var (
		err error
		res []byte
		ext = strings.ToLower(strings.TrimPrefix(filepath.Ext(fileName), "."))
	)
	if len(beautify) <= 0 {
		beautify = append(beautify, false)
	}
	switch ext {
	case `json`, `json5`:
		if !beautify[0] {
			res, err = json.Marshal(m.Swagger)
		} else {
			res, err = json.MarshalIndent(m.Swagger, "", `  `)
		}
	case `yaml`, `yml`:
		res, err = yaml.Marshal(m.Swagger)
	}

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer func(f *os.File) {
		err = f.Close()
	}(f)

	_, err = f.Write(res)
	if err != nil {
		return err
	}

	return nil
}

// 递归处理目录
func (m *Merger) recursive(root string, pattern func(string) bool) filepath.WalkFunc {
	root, _ = filepath.Abs(root)
	return func(path string, info fs.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		if !strings.HasPrefix(path, root) {
			path = filepath.Join(root, path)
		}
		if info.IsDir() {
			if root != path {
				return filepath.Walk(path, m.recursive(path, pattern))
			}
			return nil
		}
		if !pattern(path) {
			return nil
		}
		return m.AddFile(path)
	}
}

func (m *Merger) CreatePatternFilter(suffix []string) func(string) bool {
	return func(s string) bool {
		for _, v := range suffix {
			if strings.HasSuffix(s, v) {
				return true
			}
		}
		return false
	}
}
