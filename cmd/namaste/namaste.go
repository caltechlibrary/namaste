package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	// Caltech Library packages
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/namaste"
)

var (
	// Standard options
	showHelp             bool
	showVersion          bool
	showLicense          bool
	verbose              bool
	generateMarkdownDocs bool

	// App specific options
	dName  string
	asJSON bool
)

func main() {
	appName := path.Base(os.Args[0])
	app := cli.NewCli(namaste.Version)
	// We require an "ACTION" or verb for command to work.
	app.ActionsRequired = true

	// Add help assets
	for k, v := range Help {
		app.AddHelp(k, v)
	}

	// Document our verbs
	app.AddVerb("type", "returns the type of a directory if known")
	app.AddVerb("who", "returns the who value of a directory if known")
	app.AddVerb("what", "returns the what value of a directory if known")
	app.AddVerb("when", "returns the when value of a directory if known")
	app.AddVerb("where", "returns the where value of a directory if known")
	app.AddVerb("get", "returns all the namaste metadata of a directory if known")
	app.AddVerb("gettypes", "returns the types of a directory if known")

	// Standard Options
	app.BoolVar(&showHelp, "h,help", false, "display help")
	app.BoolVar(&showVersion, "v,version", false, "display program version")
	app.BoolVar(&showLicense, "l,license", false, "display license")
	app.BoolVar(&verbose, "V,verbose", false, "verbose output")
	app.BoolVar(&generateMarkdownDocs, "generate-markdown-docs", false, "output documentation in Markdown")

	// App Options
	app.StringVar(&dName, "d,directory", ".", "directory")
	app.BoolVar(&asJSON, "json", false, "output in JSON format")

	app.Parse()

	args := app.Args()

	if generateMarkdownDocs {
		app.GenerateMarkdownDocs(app.Out)
		os.Exit(0)
	}

	if showHelp {
		if showHelp {
			if len(args) > 0 {
				fmt.Fprintf(os.Stdout, app.Help(args...))
			} else {
				app.Usage(os.Stdout)
			}
			os.Exit(0)
		}
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

	var (
		s   string
		err error
	)

	// Read functions
	switch strings.ToLower(args[0]) {
	case "get":
		l, err := namaste.Get(dName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		if verbose {
			fmt.Fprintf(os.Stdout, "namastes: %s\n", strings.Join(l, ", "))
		}
		os.Exit(0)
	case "gettypes":
		m, err := namaste.GetTypes(dName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
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
			fmt.Fprintf(os.Stdout, "namaste - directory type %s - version %s %s\n", name, major, minor)
		}
		os.Exit(0)
	}

	// Write functions
	for _, arg := range args[1:] {
		switch strings.ToLower(args[0]) {
		case "type":
			s, err = namaste.DirType(dName, arg)
		case "who":
			s, err = namaste.Who(dName, arg)
		case "what":
			s, err = namaste.What(dName, arg)
		case "when":
			s, err = namaste.When(dName, arg)
		case "where":
			s, err = namaste.Where(dName, arg)
		default:
			args := os.Args[:]
			args[0] = appName
			fmt.Fprintf(os.Stderr, "Do not understand %q, type %s -help", strings.Join(args, " "), appName)
			os.Exit(1)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			os.Exit(1)
		}
		if verbose && s != "" {
			fmt.Fprintf(os.Stdout, "%s\n", s)
		}
	}
}
