package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday/v2"
)

const (
	defaultTemplate = `<!DOCTYPE html>
		<html>
			<head>
				<meta http-equiv="content-type" content="text/html; charset=utf-8">
				<title>{{ .Title }}</title>
			</head>
			<body>
				{{ .Body }}
			</body>
		</html>
		`
)

type content struct {
	Title string
	Body  template.HTML
}

func main() {
	filename := flag.String("file", "", "Markdown file to preview")
	skipPreview := flag.Bool("s", false, "Skip auto-preview")
	tFname := flag.String("t", "", "Alternate template name")
	flag.Parse()

	// If user did not provide input file, show usage
	if *filename == "" {
		flag.Usage()
		os.Exit(1)
	}
	if err := run(*filename, *tFname, os.Stdout, *skipPreview); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(filename string, tFname string, out io.Writer, skipPreview bool) error {
	input, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	htmlData, err := parseContent(input, tFname)
	if err != nil {
		return err
	}

	temp, err := os.CreateTemp("", "mdp*.html")
	if err != nil {
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}
	outName := temp.Name()

	fmt.Fprintln(out, outName)

	if err := saveHTML(outName, htmlData); err != nil {
		return err
	}

	if skipPreview {
		return nil
	}

	defer os.Remove(outName)

	return preview(outName)
}

func parseContent(input []byte, tFname string) ([]byte, error) {
	output := blackfriday.Run(input)
	body := bluemonday.UGCPolicy().SanitizeBytes(output)

	// Parse the contents of the defaultTemplate const into a new Template
	t, err := template.New("mdp").Parse(defaultTemplate)
	if err != nil {
		return nil, err
	}

	// If user provided alternate template file, replace template
	if tFname != "" {
		t, err = template.ParseFiles(tFname)
		if err != nil {
			return nil, err
		}
	}

	c := content{
		Title: "Markdown Preview Tool",
		Body:  template.HTML(body),
	}

	var buffer bytes.Buffer

	if err := t.Execute(&buffer, c); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func saveHTML(outFname string, data []byte) error {
	// Write the bytes to the file
	// Note: the third parameter represents file permissions.
	// We're using 0644 for creating a file that's both readable and writable
	// by the owner but only readable by anyone else
	return os.WriteFile(outFname, data, 0644)
}

func preview(fname string) error {
	cName := ""
	cParams := []string{}

	switch runtime.GOOS {
	case "linux":
		cName = "xdg-open"
	case "windows":
		cName = "cmd.exe"
		cParams = []string{"/C", "start"}
	case "darwin":
		cName = "open"
	default:
		return fmt.Errorf("OS not supported")
	}

	cParams = append(cParams, fname)

	// Locate executable in PATH
	cPath, err := exec.LookPath(cName)
	if err != nil {
		return err
	}

	// Open the file using the default program
	err = exec.Command(cPath, cParams...).Run()

	// Give the browser some time to open the file before deleting it
	time.Sleep(2 * time.Second)
	return err
}
