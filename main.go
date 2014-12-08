package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
)

const impressjsURL string = "http://github.com/bartaz/impress.js/zipball/0.5.3"
const impressjsName string = "bartaz-impress.js-e8fbd0c"

var defaultCss = []string{}

var defaultJs = []string{
	impressjsName + "/js/impress.js"}

var defaultJsLiteral = []string{
	"impress().init()"}

func main() {
	var err error = nil
	var outdir, infile string
	var nobasestyle bool
	flag.StringVar(&outdir, "outdir", ".", "output directory, default to pwd, will attempt creation if specified")
	flag.BoolVar(&nobasestyle, "nobasestyle", false, "specify to NOT include the impressjs demo css")
	flag.Parse()
	if len(flag.Args()) > 1 {
		l.error("multiple input files specified but not supported %v", flag.Args())
		os.Exit(1)
	}
	infile = flag.Args()[0]
	if !nobasestyle {
		defaultCss = append(defaultCss, impressjsName+"/css/impress-demo.css")
	}

	// input md file
	infile_abs, err := filepath.Abs(infile)
	if err != nil {
		l.error("failed to locate input file %s", infile)
		os.Exit(1)
	}

	if !exists(infile_abs) {
		l.error("input markdown file does NOT exist %s", infile_abs)
		os.Exit(1)
	}
	md, err := ioutil.ReadFile(infile_abs)
	if err != nil {
		l.error("something is wrong reading the markdown file: %s", err)
		os.Exit(1)
	}
	// render
	out, err := render(md, renderOptions{baseCss: defaultCss, scriptSrcs: defaultJs, scriptLiterals: defaultJsLiteral})
	if err != nil {
		l.error("rendering failed: %s", err)
	}
	// output folder has to exist otherwise create
	outdir_abs, err := filepath.Abs(outdir)
	if err != nil {
		l.error("failed to locate path %s", err)
		os.Exit(1)
	}
	l.debug("output folder abs: %s", outdir_abs)
	if !exists(outdir_abs) {
		err = os.MkdirAll(outdir_abs, os.ModeDir|0777)
		if err != nil {
			l.error("error making directory %s; caused by: %s", outdir_abs, err)
			os.Exit(1)
		}
	}
	// lets run the actual thing
	err = ioutil.WriteFile(filepath.Join(outdir_abs, "index.html"), out, 0775)
	if err != nil {
		l.error("error writing html file: %s", err)
		os.Exit(1)
	}

	// download impressjs
	impressjs_folder := filepath.Join(outdir_abs, impressjsName)
	if !exists(impressjs_folder) {
		// download
		err = getFile(impressjsURL, filepath.Join(outdir_abs, impressjsName+".zip"))
		if err != nil {
			l.error("error downloading impress js at %s; caused by: %s", impressjsURL, err)
			os.Exit(1)
		}
		// unzip
		err = unzip(filepath.Join(outdir_abs, impressjsName)+".zip", outdir_abs)
		if err != nil {
			l.error("error extracting impress js; caused by: %s", err)
			os.Exit(1)
		}
		// delete zip
		err = os.Remove(filepath.Join(outdir_abs, impressjsName) + ".zip")
		if err != nil {
			l.error("error removing impressjs zip, please remove manually if necessary; caused by: %s", err)
		}
	}

	l.info("run your presentation at %s", filepath.Join(outdir_abs, "index.html"))
	os.Exit(0)
}

func exists(p string) bool {
	var err error = nil
	if _, err = os.Stat(p); err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			l.error("error accessing directory %s; caused by: %s", p, err)
			os.Exit(1)
		}
	}
	return true
}
