package main

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func Test_Rendering(t *testing.T) {
	input, err := ioutil.ReadFile("test/test.md")
	if err != nil {
		t.Fatal("testing... failed to read input(test/test.md): %s", err)
	}

	expectOutput, err := ioutil.ReadFile("test/index.html")
	if err != nil {
		t.Fatal("testing... failed to read expected output(test/index.html): %s", err)
	}

	output, _ := render(input,
		renderOptions{baseCss: []string{impressjsName + "/css/impress-demo.css"}, scriptSrcs: defaultJs, scriptLiterals: defaultJsLiteral})

	if !bytes.Equal(expectOutput, output) {
		t.Fatal("render output does not match expected output")
	}
}
