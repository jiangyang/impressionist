% impressionist foobarx
% css http://jmblog.github.io/color-themes-for-google-code-prettify/css/themes/tomorrow-night-eighties.css
% author foolio
% title no really no
% author fooz
% title Impressionist Rocks!
% css http://fonts.googleapis.com/css?family=Open+Sans:regular,semibold,italic,italicsemibold|PT+Sans:400,700,400italic,700italic|PT+Serif:400,700,400italic,700italic
% css http://jmblog.github.io/color-themes-for-google-code-prettify/css/themes/github.css

<!-- foo bar comment-->
% slide bored
% class step slide
% x -1000
% y -1500
% duck  quux

% author foo
% title bar

Aren't you just **bored** with all those slides-based presentations?

% slide
% class slide
% x 0 data-y -1500
Don't you think that presentations given `in modern browsers` shouldn't ~~copy the limits~~ of 'classic' slide decks?

% slide
% class step slide
% data-x 1000   y   -1500
Would you like to <strong>impress your audience</strong> with <strong>stunning visualization</strong> of your talk?

% slide
% x 0 y -3000
###### Would you also like to show some code?
```go
// [Timers](timers) are for when you want to do
// something once in the future - _tickers_ are for when
// you want to do something repeatedly at regular
// intervals. Here's an example of a ticker that ticks
// periodically until we stop it.

package main

import "time"
import "fmt"

func main() {

    // Tickers use a similar mechanism to timers: a
    // channel that is sent values. Here we'll use the
    // `range` builtin on the channel to iterate over
    // the values as they arrive every 500ms.
    ticker := time.NewTicker(time.Millisecond * 500)
    go func() {
        for t := range ticker.C {
            fmt.Println("Tick at", t)
        }
    }()

    // Tickers can be stopped like timers. Once a ticker
    // is stopped it won't receive any more values on its
    // channel. We'll stop ours after 1500ms.
    time.Sleep(time.Millisecond * 1500)
    ticker.Stop()
    fmt.Println("Ticker stopped")
}
```

% slide
% class step
% data-x 0    data-y 0 scale 4
### then you should try
impress.js  
no rhyme intended


% slide its
% class step
% x 850 y 3000 rotate 90 scale 5
It's a **presentation tool**  
inspired by the idea behind [prezi.com](http://prezi.com)  
and based on the **power of CSS3 transforms and transitions** in modern browsers.

% slide big
% class step
% x 3500 y 2100 rotate 180 scale 6
visualize your **big** thoughts

% slide tiny
% class step
% x 2825 y 2325 z -3000 rotate 300 data-scale 1
and **tiny** ideas

% slide ing
% class step
% x  3500 y -850 rotate 270  scale 6
by **positioning**, **rotating** and **scaling** them on an infinite canvas

% slide imagination
% class step
% x 6700 y -300 scale 6
the only **limit** is your **imagination**

% slide source
% class step
% x 6300 y 2000 rotate 20 scale 4
want to know more?
>[use the source](http://github.com/bartaz/impress.js), Luke!

% slide one-more-thing
% class step
% data-x 6000 data-y 4000 scale 2
one more thing...

% slide its-in-3d
% class step
% x 6200 y 4300 z -100 rotate-x -40 rotate-y 10 scale 2
have you noticed it's in **3D**
* beat that, prezi ;)

% slide overview
% class step
% x 3000 y 1500
% data-scale 10
