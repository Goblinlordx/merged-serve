package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/spf13/afero"
)

type anyDirs interface {
	len() int
	name(i int) string
	isDir(i int) bool
	size(i int) int64
	modTime(i int) string
}

type fileInfoDirs []fs.FileInfo

func (d fileInfoDirs) len() int             { return len(d) }
func (d fileInfoDirs) isDir(i int) bool     { return d[i].IsDir() }
func (d fileInfoDirs) name(i int) string    { return d[i].Name() }
func (d fileInfoDirs) size(i int) int64     { return d[i].Size() }
func (d fileInfoDirs) modTime(i int) string { return d[i].ModTime().String() }

type dirEntryDirs []fs.DirEntry

func (d dirEntryDirs) len() int          { return len(d) }
func (d dirEntryDirs) isDir(i int) bool  { return d[i].IsDir() }
func (d dirEntryDirs) name(i int) string { return d[i].Name() }
func (d dirEntryDirs) size(i int) int64 {
	s, err := d[i].Info()
	if err != nil {
		return 0
	}

	return s.Size()
}
func (d dirEntryDirs) modTime(i int) string {
	s, err := d[i].Info()
	if err != nil {
		return ""
	}

	return s.ModTime().UTC().Format("2006-01-02 15:04:05Z")
}

type DirEntry struct {
	Name    string
	Path    string
	Size    string
	ModTime string
}

type autoIndexedFS struct {
	fs       http.FileSystem
	template *template.Template
}

func (aifs autoIndexedFS) Open(path string) (http.File, error) {
	idx := strings.HasSuffix(path, "/index.html")

	f, err := aifs.fs.Open(path)
	if !idx {
		if err != nil {
			return nil, err
		}

		return f, nil
	}

	p := strings.TrimSuffix(path, "/index.html")
	if len(p) == 0 {
		p = "."
	}

	f, err = aifs.fs.Open(p)
	if err != nil {
		return nil, err
	}
	if p == "." {
		p = "/"
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if !s.IsDir() {
		return f, nil
	}

	var dirs anyDirs
	if d, ok := f.(fs.ReadDirFile); ok {
		var list dirEntryDirs
		list, err = d.ReadDir(-1)
		dirs = list
	} else {
		var list fileInfoDirs
		list, err = f.Readdir(-1)
		dirs = list
	}

	if err != nil {
		return nil, err
	}
	sort.Slice(dirs, func(i, j int) bool { return dirs.name(i) < dirs.name(j) })

	entries := []DirEntry{}
	if len(p) > 1 {
		split := strings.Split(p, "/")
		up := strings.Join(split[:len(split)-1], "/")
		if len(up) == 0 {
			up = "/"
		}

		entries = append(entries, DirEntry{
			Name: "..",
			Path: up,
		})
	}
	for i, n := 0, dirs.len(); i < n; i++ {
		name := dirs.name(i)
		if dirs.isDir(i) {
			name += "/"
		}
		url := url.URL{Path: name}
		entries = append(entries, DirEntry{
			Path:    url.String(),
			Name:    name,
			Size:    humanize.Bytes(uint64(dirs.size(i))),
			ModTime: dirs.modTime(i),
		})
	}

	mmfs := afero.NewMemMapFs()
	file, err := mmfs.Create("index.html")
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(make([]byte, 0))
	err = aifs.template.Execute(buf, struct {
		Path    string
		Entries []DirEntry
	}{Path: p, Entries: entries})

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	_, err = file.WriteString(buf.String())
	if err != nil {
		return nil, err
	}

	err = file.Close()
	if err != nil {
		return nil, err
	}

	return mmfs.Open("index.html")
}
