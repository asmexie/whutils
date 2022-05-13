package whutils

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/asmexie/gopub/common"
	"github.com/patrickmn/go-cache"
)

func unescaped(x string) interface{} {
	return template.HTML(x)
}

func httpGet(reqURL string) (string, error) {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return "", common.ERR(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", common.ERR(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", common.ERR(err)
	}
	return string(data), nil
}

func (tm *tmplMGR) newTemplate(name string, parent *template.Template) *template.Template {
	var tmpl *template.Template
	if parent == nil {
		tmpl = template.New(name)
	} else {
		tmpl = parent.New(name)
	}

	tmpl.Funcs(template.FuncMap{
		"unescaped": unescaped,
		"runTmplURL": func(cacheMinutes int, reqURL string) interface{} {
			if cacheMinutes > 0 {
				if v, found := tm.cc.Get(reqURL); found {
					return template.HTML(v.(string))
				}
			}
			if data, err := httpGet(reqURL); err == nil {
				data = strings.TrimSpace(data)
				if cacheMinutes > 0 {
					tm.cc.Set(reqURL, data, time.Minute*time.Duration(cacheMinutes))
				}
				return template.HTML(data)
			}
			return ""
		},
	})

	return tmpl
}

type tmplMGR struct {
	cc *cache.Cache
}

var _tmplMGR = tmplMGR{
	cc: cache.New(time.Second*10, time.Second*20),
}

func (tm *tmplMGR) parseFiles(cacheTime time.Duration, filenames ...string) (*template.Template, error) {
	var t *template.Template
	var cacheName string
	if len(filenames) == 0 {
		// Not really a problem, but be consistent.
		return nil, fmt.Errorf("html/template: no files named in call to ParseFiles")
	}
	if cacheTime > 0 {
		v, found := tm.cc.Get(filenames[0])
		if found {
			return v.(*template.Template), nil
		}
	}

	for _, filename := range filenames {
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		s := string(b)
		var tmpl *template.Template

		name := filepath.Base(filename)
		if t == nil {
			t = tm.newTemplate(name, nil)
			cacheName = filename
		}

		if name == t.Name() {
			tmpl = t
		} else {
			tmpl = tm.newTemplate(name, t)
		}
		_, err = tmpl.Parse(s)
		if err != nil {
			return nil, err
		}

	}
	if cacheTime > 0 {
		tm.cc.Set(cacheName, t, cacheTime)
	}
	return t, nil
}

// GetTmpl ...
func (tm *tmplMGR) GetTmpl(cachetime time.Duration, webDir string, tmplNames ...string) *template.Template {
	tmplPaths := []string{}
	for _, tmplName := range tmplNames {
		tmplPath := filepath.Join(webDir, "/tmpl/", tmplName)
		tmplPaths = append(tmplPaths, tmplPath)
	}
	t, err := tm.parseFiles(cachetime, tmplPaths...)
	if err != nil {
		panic(err)
	}

	return t
}

// GetTmpl ...
func GetTmpl(cachetime time.Duration, webDir string, tmplNames ...string) *template.Template {
	return _tmplMGR.GetTmpl(cachetime, webDir, tmplNames...)
}
