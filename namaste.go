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
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func makeNamaste(tag, value string) string {
	return fmt.Sprintf("%s=%s", tag, value)
}

func getNamaste(dName, tag string) ([]string, error) {
	results := []string{}
	dInfo, err := os.Stat(dName)
	if err != nil {
		return nil, err
	}
	if dInfo.IsDir() == false {
		return nil, fmt.Errorf("expected %q to be a directory", dName)
	}

	dir, err := os.Open(dName)
	if err != nil {
		return nil, err
	}
	defer dir.Close()
	items, err := dir.Readdir(-1)
	prefix := fmt.Sprintf("%s=", tag)
	for _, item := range items {
		if strings.HasPrefix(item.Name(), prefix) {
			results = append(results, item.Name())
		}
	}
	//	sort.Strings(results)
	return results, nil
}

func setNamaste(dName, tag, value string) (string, error) {
	dInfo, err := os.Stat(dName)
	if err != nil {
		return "", err
	}
	if dInfo.IsDir() == false {
		return "", fmt.Errorf("%q is not a directory", dName)
	}
	namaste := makeNamaste(tag, value)
	return namaste, ioutil.WriteFile(path.Join(dName, namaste), []byte(value+"\n"), 0664)
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

func Get(dName string) ([]string, error) {
	dInfo, err := os.Stat(dName)
	if err != nil {
		return nil, err
	}
	if dInfo.IsDir() == false {
		return nil, fmt.Errorf("%q is not a directory", dName)
	}
	dir, err := os.Open(dName)
	if err != nil {
		return nil, err
	}
	defer dir.Close()
	items, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}
	results := []string{}
	for _, item := range items {
		name := item.Name()
		if strings.HasPrefix(name, "0=") || strings.HasPrefix(name, "1=") || strings.HasPrefix(name, "2=") || strings.HasPrefix(name, "3=") || strings.HasPrefix(name, "4=") {
			results = append(results, name)
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
	for _, t := range typeTags {
		s := strings.SplitN(strings.TrimPrefix(t, "0="), "_", 2)
		key := s[0]
		version := strings.SplitN(s[1], ".", 2)
		types[key] = map[string]string{
			"name":  key,
			"major": version[0],
			"minor": version[1],
		}
	}
	return types, nil
}
