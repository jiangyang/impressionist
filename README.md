Impressionist
=====

write a markdown file and generate an [impressjs](https://github.com/bartaz/impress.js/) presentation

## usage

`go get github.com/jiangyang/impressionist`

then

`impressionist -outdir presentation input.md`
The input markdown file will be rendered into `index.html` and then put into the output dir. A version of impressjs will be download and put in the output dir so your presentation can be run.

options:
>  -democss=false: specify to include the impressjs demo css  
>  -outdir=".": output directory, default to pwd, will attempt creation if specified  

## example markdown

`````markdown
<!--
    start a file with a pandoc block specifying the global attributes

    impressionist: specifies this block should be seen as global attributes

    author: the name of the author
    title: the title
    both will be put in the generated html, and title is also used as the page title

    css: specify links to external css to use
    common use case like: specify one for the fonts, specify another for code syntax highlighting
-->
% impressionist
% author fooz
% title Impressionist Rocks!
% css http://fonts.googleapis.com/css?family=Open+Sans
% css http://jmblog.github.io/color-themes-for-google-code-prettify/css/themes/github.css

<!-- 
    start every slide with a pandoc block setting the attributes of the slide
    slide: indicates this is a new slide, when a value is specified, the value will be used as the element id
    class: the css classes on the slide element, impressjs "step" class can be omitted and will be added by default
    
    x,y...etc are deemed as impressjs attributes they can be written in full form like data-x 1000 or "data-" part can be omitted

    more conveniently they can be written in one line, e.g.
    x -1000 y -1500 rotate-x 10 rotate-y 20 z -100 scale 1
-->
% slide bored
% class step slide
% x -1000
% y -1500
Aren't you just **bored** with all those slides-based presentations?


<!-- 
    fenced code blocks are supported and will be annotated according to http://google-code-prettify.googlecode.com
-->
% slide
% x 0 y -3000 rotate 180 scale 10
###### Would you also like to show some code?
```go
package main

import "time"
import "fmt"

func main() {

    go func() {
        for t := range ticker.C {
            fmt.Println("Tick at", t)
        }
    }()

    time.Sleep(time.Millisecond * 1500)
    ticker.Stop()
    fmt.Println("Ticker stopped")
}
```
`````
