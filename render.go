package main

import (
	"bytes"
	"errors"
	"github.com/russross/blackfriday"
	"github.com/sourcegraph/syntaxhighlight"
	"sort"
	"strconv"
	"strings"
)

// options passed in from the outside
type renderOptions struct {
	baseCss        []string
	scriptSrcs     []string
	scriptLiterals []string
}

// user defined header parameters
// defined in a block title with special name "impressionist"
type udfHeader struct {
	title  string
	author string
	css    []string
}

// wrap blackfriday's stock html renderer and override accordingly
type renderer struct {
	renderOptions
	udfHeader
	// blackfriday renderer struct
	*blackfriday.Html
	// internal state
	slideCount int
	inSlide    bool
}

func render(input []byte, opts renderOptions) ([]byte, error) {
	htmlFlags := 0 | blackfriday.HTML_HREF_TARGET_BLANK | blackfriday.HTML_USE_XHTML | blackfriday.HTML_USE_SMARTYPANTS | blackfriday.HTML_SMARTYPANTS_FRACTIONS | blackfriday.HTML_SMARTYPANTS_LATEX_DASHES
	extensions := 0 | blackfriday.EXTENSION_NO_INTRA_EMPHASIS | blackfriday.EXTENSION_TABLES | blackfriday.EXTENSION_FENCED_CODE | blackfriday.EXTENSION_AUTOLINK | blackfriday.EXTENSION_STRIKETHROUGH | blackfriday.EXTENSION_SPACE_HEADERS | blackfriday.EXTENSION_HEADER_IDS | blackfriday.EXTENSION_TITLEBLOCK
	// make blackfriday renderer
	renderer := &renderer{Html: blackfriday.HtmlRenderer(htmlFlags, "", "").(*blackfriday.Html)}
	// initialize defaults
	renderer.renderOptions = opts
	renderer.udfHeader = udfHeader{title: "Impressionist", author: "Impressionist", css: []string{}}
	// render
	buf := new(bytes.Buffer)
	bfout := blackfriday.Markdown(input, renderer, extensions)
	if renderer.slideCount == 0 {
		return []byte(""), errors.New("No slide defined")
	}
	renderer._documentHeader(buf)
	buf.Write(bfout)
	renderer._documentFooter(buf)
	return buf.Bytes(), nil
}

func (r *renderer) TitleBlock(out *bytes.Buffer, text []byte) {
	// needs to look at first key of first prop
	// if impressionist, populate udfHeader struct
	// if slide, write slide separator header
	// else let blackfriday handle it
	text = bytes.TrimPrefix(text, []byte("% "))
	keyvals := bytes.Split(text, []byte("\n% "))
	parts := bytes.Split(keyvals[0], []byte(" "))
	if bytes.Equal([]byte("impressionist"), parts[0]) {
		r.updateUdfHeader(keyvals)
	} else if bytes.Equal([]byte("slide"), parts[0]) {
		r.writeSlideHeader(out, keyvals)
	} else {
		r.Html.TitleBlock(out, text)
	}
}

func (r *renderer) updateUdfHeader(keyvals [][]byte) {
	for idx, keyval := range keyvals {
		if idx == 0 {
			continue
		}
		key, val := splitKeyVal(keyval)
		keyStr := string(key)
		if keyStr == "title" {
			r.udfHeader.title = string(val)
		} else if keyStr == "author" {
			r.udfHeader.author = string(val)
		} else if keyStr == "css" {
			r.udfHeader.css = append(r.udfHeader.css, string(val))
		}
	}
}

func (r *renderer) writeSlideHeader(out *bytes.Buffer, keyvals [][]byte) {
	attrsMap := make(map[string]string)
	for i, keyval := range keyvals {
		key, val := splitKeyVal(keyval)
		if i == 0 {
			r.slideCount += 1
			if len(val) == 0 {
				val = []byte("slide_" + strconv.Itoa(r.slideCount))
			}
			attrsMap["id"] = string(val)
		} else if bytes.Equal([]byte("class"), key) {
			attrsMap["class"] = string(val)
		} else {
			parseAttrLine(attrsMap, keyval)
		}
	}
	l.debug("attributes map %v", attrsMap)

	if r.inSlide {
		out.WriteString("\n</div>\n")
	}
	if c, ok := attrsMap["class"]; !ok {
		attrsMap["class"] = "step"
	} else if c[:4] != "step" {
		attrsMap["class"] = "step " + c
	}
	out.WriteString("<div ")
	// sort map keys
	idx := 0
	keys := make([]string, len(attrsMap))
	for k_, _ := range attrsMap {
		keys[idx] = k_
		idx++
	}
	sort.Strings(keys)
	for _, k := range keys {
		attrEscape(out, []byte(k))
		out.WriteString("=\"")
		attrEscape(out, []byte(attrsMap[k]))
		out.WriteString("\" ")
	}
	out.WriteString(">\n")
	r.inSlide = true
}

func parseAttrLine(m map[string]string, keyval []byte) {
	var key, val string
	var keyset, valset bool
	for spaceloc := bytes.IndexRune(keyval, ' '); spaceloc != -1; spaceloc = bytes.IndexRune(keyval, ' ') {
		if spaceloc != 0 {
			if !keyset {
				key = string(keyval[0:spaceloc])
				keyset = true
			} else if !valset {
				val = string(keyval[0:spaceloc])
				valset = true
			}
			if keyset && valset {
				m[prefixKey(key)] = val
				keyset = false
				valset = false
			}
		}
		keyval = keyval[spaceloc+1:]
	}
	if !keyset {
		key = string(keyval)
	} else if !valset {
		val = string(keyval)
	}
	m[prefixKey(key)] = val
}

func prefixKey(k string) string {
	if strings.TrimPrefix(k, "data-") == k {
		return "data-" + k
	}
	return k
}

func splitKeyVal(keyval []byte) (key, val []byte) {
	spaceloc := bytes.IndexRune(keyval, ' ')
	if spaceloc != -1 {
		key = keyval[0:spaceloc]
		val = keyval[spaceloc:]
	}
	key = bytes.TrimSpace(key)
	val = bytes.TrimSpace(val)
	return
}

func (r *renderer) BlockCode(out *bytes.Buffer, text []byte, lang string) {
	count := 0
	for _, elt := range strings.Fields(lang) {
		if elt[0] == '.' {
			elt = elt[1:]
		}
		if len(elt) == 0 {
			continue
		}
		if count == 0 {
			out.WriteString("<pre class=\"prettyprint\"><code class=\"language-")
		} else {
			out.WriteByte(' ')
		}
		attrEscape(out, []byte(elt))
		count++
	}

	if count == 0 {
		out.WriteString("<pre class=\"prettyprint\"><code>")
	} else {
		out.WriteString("\">")
	}

	highlighted, err := syntaxhighlight.AsHTML(text)
	if err != nil {
		l.error("syntax highlight failure: %s", err)
		attrEscape(out, text)
	} else {
		out.Write(highlighted)
	}
	out.WriteString("</code></pre>\n")
}

func (r *renderer) DocumentHeader(_ *bytes.Buffer) {
	// pass
}

func (r *renderer) _documentHeader(out *bytes.Buffer) {
	out.WriteString("<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" ")
	out.WriteString("\"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\">\n")
	out.WriteString("<html xmlns=\"http://www.w3.org/1999/xhtml\">\n")
	out.WriteString("<head>\n")
	out.WriteString("  <title>")
	attrEscape(out, []byte(r.udfHeader.title))
	out.WriteString("  </title>\n")
	out.WriteString("  <meta charset=\"utf-8\"/>\n")
	out.WriteString("  <meta title=\"")
	attrEscape(out, []byte(r.udfHeader.title))
	out.WriteString("\"/>\n")
	out.WriteString("  <meta author=\"")
	attrEscape(out, []byte(r.udfHeader.author))
	out.WriteString("\"/>\n")
	// base css provided by impressionist
	for _, s := range r.renderOptions.baseCss {
		out.WriteString("  <link rel=\"stylesheet\" href=\"")
		attrEscape(out, []byte(s))
		out.WriteString("\"/>\n")
	}
	// user defined css in udf header
	for _, s := range r.udfHeader.css {
		out.WriteString("  <link rel=\"stylesheet\" href=\"")
		attrEscape(out, []byte(s))
		out.WriteString("\"/>\n")
	}
	out.WriteString("</head>\n")
	// impressjs
	out.WriteString("<body>\n")
	out.WriteString("<div id=\"impress\">\n")
}

func (r *renderer) DocumentFooter(out *bytes.Buffer) {
	if r.inSlide {
		out.WriteString("\n</div>\n")
	}
	out.WriteString("\n</div>\n")
}

func (r *renderer) _documentFooter(out *bytes.Buffer) {
	for _, s := range r.renderOptions.scriptSrcs {
		out.WriteString("<script src=\"")
		attrEscape(out, []byte(s))
		out.WriteString("\"></script>\n")
	}

	for _, s := range r.renderOptions.scriptLiterals {
		out.WriteString("<script>")
		out.WriteString(s)
		out.WriteString("</script>\n")
	}

	out.WriteString("\n</body>\n</html>\n")
}

// copied from russross/blackfriday

// Using if statements is a bit faster than a switch statement. As the compiler
// improves, this should be unnecessary this is only worthwhile because
// attrEscape is the single largest CPU user in normal use.
// Also tried using map, but that gave a ~3x slowdown.
func escapeSingleChar(char byte) (string, bool) {
	if char == '"' {
		return "&quot;", true
	}
	if char == '&' {
		return "&amp;", true
	}
	if char == '<' {
		return "&lt;", true
	}
	if char == '>' {
		return "&gt;", true
	}
	return "", false
}

func attrEscape(out *bytes.Buffer, src []byte) {
	org := 0
	for i, ch := range src {
		if entity, ok := escapeSingleChar(ch); ok {
			if i > org {
				// copy all the normal characters since the last escape
				out.Write(src[org:i])
			}
			org = i + 1
			out.WriteString(entity)
		}
	}
	if org < len(src) {
		out.Write(src[org:])
	}
}
