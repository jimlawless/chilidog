# chilidog
( a.k.a. Jim's personal static site generator )

I don't expect anyone else to use this but if you're feeling adventurous, here's a brief rundown of how chilidog works.

It passes the input file looking for items enclosed in square brackets.  When it finds one, it tries to look up a function name in the first word ( delimited by a space ) in those brackets.  If any other parameters are present before the closing bracket, they should be delimited by the pipe | symbol.

Text outside of brackets is passed through as HTML with the exception of blank lines.  One or more blank lines is reduced to a single paragraph tag.

In most cases, the chilidog tags are almost identical to their HTML counterparts, but I find it easier to use them as I'm typing. You'll find that some of these ... like "img" have opinionated styling that stems from the program and not from the stylesheet alone. ( Right now, all images are centered. )

To kick the tires, copy the items in the src folder somewhere, compile chilidog.go,  and run the chili_all.sh script.  It should produce an HTML file called index.htm that has all of the output from index.txt, frag1.txt, and frag2.txt.  The output file is dependent on the chili-prefixed styles in the provided style.css.

This is really specific to the simple approach to my web sites.  You can see how I use this on my sites below.  

https://jimlawless.net

https://jumptable.net

https://straypointers.com

