package xshell

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"strings"
)

var cache *db

type db struct {
	Path      string
	Resources []*resource
}

func (d *db) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("Path: %v\r\n", d.Path))
	if len(d.Resources) == 0 {
		buf.WriteString("\tresource: nil")
	} else {
		for i, resource := range d.Resources {
			buf.WriteString(fmt.Sprintf("\tresource: %3d)%v\r\n", i, resource.String()))
		}
	}

	return buf.String()
}

func (d *db) load() error {
	file, err := os.OpenFile(d.Path, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(d); err != nil {
		return err
	}
	return nil
}

func (d *db) dump() error {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(d); err != nil {
		return err
	}

	file, err := os.OpenFile(d.Path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(buf.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func (d *db) get(index int) *resource {
	if index >= d.size() || index < 0 {
		return nil
	}
	return d.Resources[index]
}

func (d *db) query() []*resource {
	return d.Resources
}

func (d *db) match(key string) []*resource {
	var result = make([]*resource, 0, len(d.Resources))
	for _, resource := range d.Resources {
		if strings.Contains(resource.String(), key) {
			result = append(result, resource)
		}
	}
	return result
}

func (d *db) size() int {
	return len(d.Resources)
}

func (d *db) insert(resource *resource) int {
	resource.No = d.size()
	d.Resources = append(d.Resources, resource)
	return 1
}

func (d *db) delete(index int) int {
	if index >= d.size() || index < 0 {
		return 0
	}

	for i := index + 1; i < d.size(); i++ {
		d.Resources[i].No = i - 1
	}

	if index == d.size()-1 {
		d.Resources = d.Resources[:index]
	} else {
		d.Resources = append(d.Resources[:index], d.Resources[index+1:]...)
	}

	return 1
}

func (d *db) update(index int, resource *resource) int {
	if index >= d.size() || index < 0 {
		return 0
	}
	d.Resources[index] = resource
	return 1
}
