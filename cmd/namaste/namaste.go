package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	// Caltech Library packages
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/namaste"
)

var (
	synopsis = `_%s_ is a tool for adding metadata as directory entry`

	description = `_%s_ is a command line tool for basic metadata 
in "name as text" format with a directory. . The metadata supported 
includes directory type (with major/minor version numbers), who 
created it, what is it, what it is, when was it created, where it 
was created.

You can see "namaste" metadata by looking at the directory
contents without any special software. Namaste fields start with
zero (type), one (who), two (what), three (when) or four (where).
This is followed by an equal sign then the value of the metadata
field. E.g.

` + "```" + `
   0=bagit_0.1
   1=Twain,M.
   2=Hamlet
   3=2008
   4=Seattle
` + "```" + `
`

	examples = `
Here is an example of workflow to add type, author, title, 
year and place to a raw ePub folder named "hamlet-epub".
` + "```" + `
    cd hamlet-epub
	namaste type ePub_3
	namaste who "Twain, Mark"
	namaste what "Hamlet"
	namaste when "2008"
	namaste where "Seattle, Washington, USA"
` + "```" + `
`

	license = `
%s %s

Copyright (c) 2018, Caltech
All rights not granted herein are expressly reserved by Caltech.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
`

	// Standard options
	showHelp         bool
	showVersion      bool
	showLicense      bool
	verbose          bool
	generateMarkdown bool
	generateManPage  bool

	// App specific options
	dName    string
	asValues bool
	asJSON   bool
)

func main() {
	appName := path.Base(os.Args[0])
	app := cli.NewCli(namaste.Version)
	// We require an "ACTION" or verb for command to work.
	app.VerbsRequired = true

	// Add Help Docs
	app.AddHelp("license", []byte(fmt.Sprintf(license, appName, namaste.Version)))
	app.AddHelp("synopsis", []byte(fmt.Sprintf(synopsis, appName)))
	app.AddHelp("description", []byte(fmt.Sprintf(description, appName)))
	app.AddHelp("examples", []byte(examples))

	// Add help assets
	for k, v := range Help {
		app.AddHelp(k, v)
	}

	// Standard Options
	app.BoolVar(&showHelp, "h,help", false, "display help")
	app.BoolVar(&showVersion, "v,version", false, "display program version")
	app.BoolVar(&showLicense, "l,license", false, "display license")
	app.BoolVar(&verbose, "V,verbose", true, "(default true) verbose output")

	app.BoolVar(&generateMarkdown, "generate-markdown", false, "output documentation in Markdown")
	app.BoolVar(&generateManPage, "generate-manpage", false, "output documentation in 'nroff -man' format")

	// App Options
	app.StringVar(&dName, "d,directory", ".", "directory")
	app.BoolVar(&asJSON, "json", false, "output in JSON format")
	app.BoolVar(&asValues, "values", false, "output value only, one per line")

	// Read Verbs
	verb := app.NewVerb("get", "returns namaste metadata of a directory if known", func(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
		err := flagSet.Parse(args)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		args = flagSet.Args()

		l, err := namaste.Get(dName, args)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		if asValues {
			for i, val := range l {
				l[i] = namaste.Decode(val)
			}
		}
		if asJSON {
			src, err := json.Marshal(l)
			if err != nil {
				fmt.Fprintf(eout, "%s\n", err)
				return 1
			}
			fmt.Fprintf(out, "%s\n", src)
			return 0
		}
		if asValues {
			fmt.Fprintf(out, "%s\n", strings.Join(l, "\n"))
		} else {
			fmt.Fprintf(out, "namastes: %s\n", strings.Join(l, ", "))
		}
		return 0
	})
	verb.SetParams("[TYPE]", "[TYPE ...]")
	verb.BoolVar(&asValues, "values", false, "output values only, one per line")
	verb.BoolVar(&asJSON, "j,json", false, "set json output")

	verb = app.NewVerb("gettypes", "returns the types of a directory if known", func(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
		err := flagSet.Parse(args)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		args = flagSet.Args()

		m, err := namaste.GetTypes(dName)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		if asJSON {
			src, err := json.Marshal(m)
			if err != nil {
				fmt.Fprintf(eout, "%s\n", err)
				return 1
			}
			fmt.Fprintf(out, "%s\n", src)
			return 0
		}
		for _, val := range m {
			name, major, minor := "", "", ""
			if s, ok := val["name"]; ok == true {
				name = s
			}
			if s, ok := val["major"]; ok == true {
				major = s
			}
			if s, ok := val["minor"]; ok == true {
				minor = s
			}
			// Insert newline if required
			if asValues {
				fmt.Fprintf(out, "%s_%s.%s\n", name, major, minor)
			} else {
				fmt.Fprintf(out, "namaste - directory type %s - version %s %s\n", name, major, minor)
			}
		}
		return 0
	})
	verb.BoolVar(&asValues, "values", false, "output values only, one per line")
	verb.BoolVar(&asJSON, "j,json", false, "set json output")

	// Write Verbs
	verb = app.NewVerb("type", "set the type of a directory", func(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
		err := flagSet.Parse(args)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		args = flagSet.Args()
		if len(args) == 0 {
			fmt.Fprintf(eout, "Missing value\n")
			return 1
		}

		s, err := namaste.DirType(dName, args[0])
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		if asJSON && s != "" {
			m := map[string]string{
				"type": s,
			}
			src, err := json.Marshal(m)
			if err != nil {
				fmt.Fprintf(eout, "%s\n", err)
				return 1
			}
			fmt.Fprintf(out, "%s\n", src)
		}
		if verbose && s != "" {
			fmt.Fprintf(out, "%s\n", s)
		}
		return 0
	})
	verb.BoolVar(&asJSON, "j,json", false, "set json output")
	verb.BoolVar(&verbose, "V,verbose", false, "set verbose output")

	verb = app.NewVerb("who", "sets the who value of a directory", func(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
		err := flagSet.Parse(args)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		args = flagSet.Args()
		if len(args) == 0 {
			fmt.Fprintf(eout, "Missing value\n")
			return 1
		}

		s, err := namaste.Who(dName, args[0])
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		if asJSON && s != "" {
			m := map[string]string{
				"who": s,
			}
			src, err := json.Marshal(m)
			if err != nil {
				fmt.Fprintf(eout, "%s\n", err)
				return 1
			}
			fmt.Fprintf(out, "%s\n", src)
		}
		if verbose && s != "" {
			fmt.Fprintf(out, "%s\n", s)
		}
		return 0
	})
	verb.BoolVar(&asJSON, "j,json", false, "set json output")
	verb.BoolVar(&verbose, "V,verbose", false, "set verbose output")

	verb = app.NewVerb("what", "sets the what value of a directory", func(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
		err := flagSet.Parse(args)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		args = flagSet.Args()
		if len(args) == 0 {
			fmt.Fprintf(eout, "Missing value\n")
			return 1
		}

		s, err := namaste.What(dName, args[0])
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		if asJSON && s != "" {
			m := map[string]string{
				"what": s,
			}
			src, err := json.Marshal(m)
			if err != nil {
				fmt.Fprintf(eout, "%s\n", err)
				return 1
			}
			fmt.Fprintf(out, "%s\n", src)
		}
		if verbose && s != "" {
			fmt.Fprintf(out, "%s\n", s)
		}
		return 0
	})
	verb.BoolVar(&asJSON, "j,json", false, "set json output")
	verb.BoolVar(&verbose, "V,verbose", false, "set verbose output")

	verb = app.NewVerb("when", "sets the when value of a directory", func(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
		err := flagSet.Parse(args)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		args = flagSet.Args()
		if len(args) == 0 {
			fmt.Fprintf(eout, "Missing value\n")
			return 1
		}

		s, err := namaste.When(dName, args[0])
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		if asJSON && s != "" {
			m := map[string]string{
				"when": s,
			}
			src, err := json.Marshal(m)
			if err != nil {
				fmt.Fprintf(eout, "%s\n", err)
				return 1
			}
			fmt.Fprintf(out, "%s\n", src)
		}
		if verbose && s != "" {
			fmt.Fprintf(out, "%s\n", s)
		}
		return 0
	})
	verb.BoolVar(&asJSON, "j,json", false, "set json output")
	verb.BoolVar(&verbose, "V,verbose", false, "set verbose output")

	verb = app.NewVerb("where", "sets the where value of a directory", func(in io.Reader, out io.Writer, eout io.Writer, args []string, flagSet *flag.FlagSet) int {
		err := flagSet.Parse(args)
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		args = flagSet.Args()
		if len(args) == 0 {
			fmt.Fprintf(eout, "Missing value\n")
			return 1
		}

		s, err := namaste.Where(dName, args[0])
		if err != nil {
			fmt.Fprintf(eout, "%s\n", err)
			return 1
		}
		if asJSON && s != "" {
			m := map[string]string{
				"where": s,
			}
			src, err := json.Marshal(m)
			if err != nil {
				fmt.Fprintf(eout, "%s\n", err)
				return 1
			}
			fmt.Fprintf(out, "%s\n", src)
		}
		if verbose && s != "" {
			fmt.Fprintf(out, "%s\n", s)
		}
		return 0
	})
	verb.BoolVar(&asJSON, "j,json", false, "set json output")
	verb.BoolVar(&verbose, "V,verbose", false, "set verbose output")

	app.Parse()

	args := app.Args()

	if generateMarkdown {
		app.GenerateMarkdown(os.Stdout)
		os.Exit(0)
	}
	if generateManPage {
		app.GenerateManPage(os.Stdout)
		os.Exit(0)
	}

	if showHelp {
		if len(args) > 0 {
			fmt.Fprintf(os.Stdout, app.Help(args...))
		} else {
			app.Usage(os.Stdout)
		}
		os.Exit(0)
	}

	if showLicense {
		fmt.Fprintln(os.Stdout, app.License())
		os.Exit(0)
	}

	if showVersion {
		fmt.Fprintln(os.Stdout, app.Version())
		os.Exit(0)
	}

	if len(args) < 1 {
		app.Usage(os.Stderr)
		os.Exit(1)
	}

	// Application Logic
	exitCode := app.Run(args)
	if exitCode != 0 {
		os.Exit(exitCode)
	}
}
