//
// namaste.go implements the "Name as Text" encoding/decoding for embedding metadata signatures in a directory.
//
// Authors R. S. Doiel, <rsdoiel@library.caltech.edu>
//
// Copyright (c) 2018, Caltech
// All rights not granted herein are expressly reserved by Caltech.
//
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
package namaste

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	// Caltech Library Package
	"github.com/caltechlibrary/storage"
)

const (
	Version = `v0.0.2`
)

var (
	normalizeFieldName = map[string]string{
		"type":  "0",
		"who":   "1",
		"what":  "2",
		"when":  "3",
		"where": "4",
	}
)

func Encode(tag, value string) string {
	if s, ok := normalizeFieldName[strings.ToLower(tag)]; ok == true {
		tag = s
	}
	return fmt.Sprintf("%s=%s", tag, value)
}

func Decode(value string) string {
	for _, prefix := range []string{"0=", "1=", "2=", "3=", "4="} {
		if strings.HasPrefix(value, prefix) {
			return value[2:]
		}
	}
	return value
}

func getNamaste(dName, tag string) ([]string, error) {
	options := map[string]interface{}{}
	sType := storage.FS
	prefix := Encode(tag, "")
	if strings.HasPrefix(dName, "s3://") {
		sType = storage.S3
		u, err := url.Parse(dName)
		if err != nil {
			return nil, err
		}
		if os.Getenv("AWS_SDK_LOAD_CONFIG") == "" {
			os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
		}
		options = storage.EnvToOptions(os.Environ())
		options["AwsBucket"] = u.Host
		if u.Path != "" {
			dName = path.Join(u.Path, prefix)
		} else {
			dName = prefix
		}
		if strings.HasPrefix(dName, "/") {
			dName = dName[1:]
		}
		store, err := storage.Init(sType, options)
		if err != nil {
			return nil, err
		}
		results := []string{}
		items, err := store.ReadDir(dName)
		if err != nil {
			return nil, err
		}
		for _, item := range items {
			name := item.Name()
			if strings.HasPrefix(name, prefix) {
				results = append(results, name)
			}
		}
		return results, nil
	}
	if strings.HasPrefix(dName, "gs://") {
		sType = storage.GS
		u, err := url.Parse(dName)
		if err != nil {
			return nil, err
		}
		options := storage.EnvToOptions(os.Environ())
		options["GoogleBucket"] = u.Host
		if u.Path != "" {
			dName = path.Join(u.Path, prefix)
		} else {
			dName = prefix
		}
		if strings.HasPrefix(dName, "/") {
			dName = dName[1:]
		}
		store, err := storage.Init(sType, options)
		if err != nil {
			return nil, err
		}
		results := []string{}
		items, err := store.ReadDir(dName)
		if err != nil {
			return nil, err
		}
		for _, item := range items {
			name := item.Name()
			if strings.HasPrefix(name, prefix) {
				results = append(results, name)
			}
		}
		return results, nil
	}

	// NOTE: if we get this far we assume we're working with
	// storage that supports the concept of directory/folder
	store, err := storage.Init(sType, options)
	if err != nil {
		return nil, err
	}
	results := []string{}
	dInfo, err := store.Stat(dName)
	if err != nil {
		return nil, err
	}
	if dInfo.IsDir() == false {
		return nil, fmt.Errorf("expected %q to be a directory", dName)
	}
	items, err := store.ReadDir(dName)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		name := item.Name()
		if strings.HasPrefix(name, prefix) {
			results = append(results, name)
		}
	}
	return results, nil
}

func setNamaste(dName, tag, value string) (string, error) {
	options := map[string]interface{}{}
	sType := storage.FS
	if strings.HasPrefix(dName, "s3://") {
		sType = storage.S3
		u, err := url.Parse(dName)
		if err != nil {
			return "", err
		}
		if os.Getenv("AWS_SDK_LOAD_CONFIG") == "" {
			os.Setenv("AWS_SDK_LOAD_CONFIG", "1")
		}
		options = storage.EnvToOptions(os.Environ())
		options["AwsBucket"] = u.Host
		if u.Path != "" {
			dName = u.Path
		} else {
			dName = ""
		}
	}
	if strings.HasPrefix(dName, "gs://") {
		sType = storage.GS
		u, err := url.Parse(dName)
		if err != nil {
			return "", err
		}
		options = storage.EnvToOptions(os.Environ())
		options["GoogleBucket"] = u.Host
		if u.Path != "" {
			dName = u.Path
		} else {
			dName = ""
		}
	}
	store, err := storage.Init(sType, options)
	if err != nil {
		return "", err
	}
	if sType == storage.FS {
		dInfo, err := store.Stat(dName)
		if err != nil {
			return "", err
		}
		if dInfo.IsDir() == false {
			return "", fmt.Errorf("%q is not a directory", dName)
		}
	}
	sNamaste := Encode(tag, value)
	return sNamaste, store.WriteFile(path.Join(dName, sNamaste), []byte(value+"\n"), 0664)
}

func DirType(dName, val string) (string, error) {
	return setNamaste(dName, "0", val)
}

func Who(dName, val string) (string, error) {
	return setNamaste(dName, "1", val)
}

func What(dName, val string) (string, error) {
	return setNamaste(dName, "2", val)
}

func When(dName, val string) (string, error) {
	return setNamaste(dName, "3", val)
}

func Where(dName, val string) (string, error) {
	return setNamaste(dName, "4", val)
}

func Get(dName string, kinds []string) ([]string, error) {
	if len(kinds) == 0 {
		kinds = []string{"0", "1", "2", "3", "4"}
	} else {
		// Convert to numeric string from human text, e.g. type, who, when
		for i, val := range kinds {
			if s, ok := normalizeFieldName[strings.ToLower(val)]; ok == true {
				kinds[i] = s
			}
		}
	}
	results := []string{}
	for _, kind := range kinds {
		l, err := getNamaste(dName, kind)
		if err != nil {
			return results, err
		}
		if len(l) > 0 {
			results = append(results, l...)
		}
	}
	return results, nil
}

func GetTypes(dName string) (map[string]map[string]string, error) {
	typeTags, err := getNamaste(dName, "0")
	if err != nil {
		return nil, err
	}
	types := map[string]map[string]string{}
	var (
		name    string
		version []string
	)
	for _, t := range typeTags {
		s := strings.SplitN(strings.TrimPrefix(t, "0="), "_", 2)
		name = s[0]
		if len(s) > 1 {
			version = strings.SplitN(s[1], ".", 2)
		} else {
			version = []string{}
		}
		switch len(version) {
		case 2:
			types[name] = map[string]string{
				"name":  name,
				"major": version[0],
				"minor": version[1],
			}
		case 1:
			types[name] = map[string]string{
				"name":  name,
				"major": version[1],
			}
		default:
			types[name] = map[string]string{
				"name": name,
			}
		}
	}
	return types, nil
}
