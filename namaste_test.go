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
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"testing"
)

const (
	testDir = "namaste-test"
)

func copyDir(srcName string, destName string) (err error) {
	srcInfo, err := os.Stat(srcName)
	if err != nil {
		return err
	}

	err = os.MkdirAll(destName, srcInfo.Mode())
	if err != nil {
		return err
	}

	dp, err := os.Open(srcName)
	if err != nil {
		return err
	}
	defer dp.Close()

	items, err := dp.Readdir(-1)
	for _, item := range items {

		srcPath := path.Join(srcName, item.Name())
		destPath := path.Join(destName, item.Name())

		if item.IsDir() {
			err = copyDir(srcPath, destPath)
			if err != nil {
				return err
			}
		} else {
			err = copyFile(srcPath, destPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func copyFile(srcName string, destName string) (err error) {
	srcFile, err := os.Open(srcName)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destName)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err == nil {
		srcInfo, err := os.Stat(srcName)
		if err != nil {
			err = os.Chmod(destName, srcInfo.Mode())
		}
	}
	return err
}

func TestType(t *testing.T) {
	// Setup this test
	data := "bagit_0.1"
	dataName := fmt.Sprintf("0=%s", data)
	fName := path.Join(testDir, dataName)
	os.RemoveAll(fName)

	// Test method
	_, err := DirType(testDir, data)
	if err != nil {
		t.Errorf("DirType(%q, %q), %q", testDir, dataName, err)
	}
	//log.Printf("DEBUG DirType(%q, %q) -> %q", testDir, dataName, s)

	// Validate the result
	expected := []byte(data + "\n")
	src, err := ioutil.ReadFile(fName)
	if err != nil {
		t.Errorf("Missing %q, %s", fName, err)
		t.FailNow()
	}
	if bytes.Equal(src, expected) == false {
		t.Errorf("expected %q, got %q", expected, src)
	} else {
		os.RemoveAll(fName)
	}
}

func TestWho(t *testing.T) {
	// Setup this test
	data := "Fayman,R."
	dataName := fmt.Sprintf("1=%s", data)
	fName := path.Join(testDir, dataName)
	os.RemoveAll(fName)

	// Test Method
	Who(testDir, data)

	// Validate the result
	expected := []byte(data + "\n")
	src, err := ioutil.ReadFile(fName)
	if err != nil {
		t.Errorf("Missing %q, %s", fName, err)
		t.FailNow()
	}
	if bytes.Equal(src, expected) == false {
		t.Errorf("expected %q, got %q", expected, src)
	} else {
		os.RemoveAll(fName)
	}
}

func TestWhat(t *testing.T) {
	// Setup this test
	data := "Particles"
	dataName := fmt.Sprintf("2=%s", data)
	fName := path.Join(testDir, dataName)
	os.RemoveAll(fName)

	// Test Method
	What(testDir, data)

	// Validate the result
	expected := []byte(data + "\n")
	src, err := ioutil.ReadFile(fName)
	if err != nil {
		t.Errorf("Missing %q, %s", fName, err)
		t.FailNow()
	}
	if bytes.Equal(src, expected) == false {
		t.Errorf("expected %q, got %q", expected, src)
	} else {
		os.RemoveAll(fName)
	}
}

func TestWhen(t *testing.T) {
	// Setup this test
	data := "2018"
	dataName := fmt.Sprintf("3=%s", data)
	fName := path.Join(testDir, dataName)
	os.RemoveAll(fName)

	// Test Method
	When(testDir, data)

	// Validate the result
	expected := []byte(data + "\n")
	src, err := ioutil.ReadFile(fName)
	if err != nil {
		t.Errorf("Missing %q, %s", fName, err)
		t.FailNow()
	}
	if bytes.Equal(src, expected) == false {
		t.Errorf("expected %q, got %q", expected, src)
	} else {
		os.RemoveAll(fName)
	}
}

func TestWhere(t *testing.T) {
	// Setup this test
	data := "Pasadena"
	dataName := fmt.Sprintf("4=%s", data)
	fName := path.Join(testDir, dataName)
	os.RemoveAll(fName)

	// Test Method
	Where(testDir, data)

	// Validate the result
	expected := []byte(data + "\n")
	src, err := ioutil.ReadFile(fName)
	if err != nil {
		t.Errorf("Missing %q, %s", fName, err)
		t.FailNow()
	}
	if bytes.Equal(src, expected) == false {
		t.Errorf("expected %q, got %q", expected, src)
	} else {
		os.RemoveAll(fName)
	}
}

func TestGet(t *testing.T) {
	// Setup some test data to "Get"
	cleanOK := true
	DirType(testDir, "bagit_0.1")
	Who(testDir, "Feynman,R.")
	What(testDir, "Particles")
	When(testDir, "2018")
	Where(testDir, "Pasadena")
	expected := map[string]string{
		"type":  "0=bagit_0.1",
		"who":   "1=Feynman,R.",
		"what":  "2=Particles",
		"when":  "3=2018",
		"where": "4=Pasadena",
	}
	tags, err := Get(testDir, nil)
	if err != nil {
		t.Errorf("Get(%q) failed, %s", testDir, err)
		cleanOK = false
	}
	// Check for expected values
	for key, val := range expected {
		foundIt := false
		for _, tagVal := range tags {
			if tagVal == val {
				foundIt = true
				break
			}
		}
		if foundIt == false {
			t.Errorf("%s, missing %s in tags %+v", val, key, tags)
			cleanOK = false
		}
	}
	// Check for unexpected values
	if len(tags) != len(expected) {
		t.Errorf("Unexpected value(s) %+v", tags)
	}

	// Check for specific type request
	l, err := Get(testDir, []string{"who"})
	if err != nil {
		t.Errorf("Get(%q) failed, %s", testDir, err)
		cleanOK = false
	}
	if len(l) != 1 {
		t.Errorf("expected length 1, got %d - %+v", len(l), l)
		cleanOK = false
	}
	if l[0] != "1=Feynman,R." {
		t.Errorf("expected '1=Feynman,R.', got %q", l[0])
		cleanOK = false
	}

	// Cleanup after test
	if cleanOK {
		for key, _ := range expected {
			os.RemoveAll(path.Join(testDir, key))
		}
	}
}

func TestGetTypes(t *testing.T) {
	// Setup Test data
	expected := map[string]map[string]string{
		"bagit_2323.1": map[string]string{
			"name":  "bagit",
			"major": "2323",
			"minor": "1",
		},
		"redd_0.1333": map[string]string{
			"name":  "redd",
			"major": "0",
			"minor": "1333",
		},
		"dflat_34.22": map[string]string{
			"name":  "dflat",
			"major": "34",
			"minor": "22",
		},
	}
	for key, _ := range expected {
		os.Remove(path.Join(testDir, key))
		DirType(testDir, key)
	}

	// Test method
	types, err := GetTypes(testDir)
	if err != nil {
		t.Errorf("GetTypes(%q) failed, %s", testDir, err)
		t.FailNow()
	}

	// Validate results
	for name, values := range expected {
		parts := strings.SplitN(name, "_", 2)
		key := parts[0]
		if results, ok := types[key]; ok == true {
			if len(results) != 3 {
				t.Errorf("Expected only three fields, got %+v", results)
			} else {
				// Check for matching values
				for k, v := range values {
					if r2, found := results[k]; found == true {
						if r2 != v {
							t.Errorf("expected %q, got %q for %s", k, r2, k)
						}
					} else {
						t.Errorf("missing %q in %+v", k, r2)
					}
				}
			}
		} else {
			t.Errorf("Missing %s in %+v", key, values)
		}
	}
}

func TestMain(m *testing.M) {
	os.RemoveAll(testDir)
	err := os.MkdirAll(testDir, 0775)
	if err != nil {
		log.Fatalf("Can't create dir %s, %s", testDir, err)
	}
	err = copyDir("docs", testDir)
	if err != nil {
		log.Fatalf("Can't setup %q, %q", testDir, err)
	}
	exitCode := m.Run()
	if exitCode == 0 {
		os.RemoveAll(testDir)
	}
	os.Exit(exitCode)
}
