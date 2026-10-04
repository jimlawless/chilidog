// Copyright 2026 - by Jim Lawless
// License: MIT / X11

package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type TagToFunc struct {
	Params int
	Fn     func([]string)
}

var tags map[string]TagToFunc

var patTag = "[[][^]]*[]]"
var patText = "([^[]|[^[]])"

var reTag, reText, reAll, reCmd *regexp.Regexp

var fout *os.File

var imgPrefix string
var takeALink string
var lastChar byte
var paraWritten bool

var verMajor = 0
var verMinor = 3
var verUpdate = 2

func main() {
	var err error
	fmt.Printf("chilidog template to HTML preprocessor v%d.%d.%d\nby Jim Lawless - https://github.com/jimlawless/chilidog\n\n", verMajor, verMinor, verUpdate)
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, "Syntax:\nchilidog input-template-filename output-HTML-filename\n")
		os.Exit(1)
	}
	initializeGlobals()
	fout, err = os.Create(os.Args[2])
	handleError(err)
	processFile(os.Args[1], true)
	fout.Close()
	fmt.Println("Processing complete.\n")
	os.Exit(0)
}

func initializeGlobals() {
	reTag, _ = regexp.Compile(patTag)
	reAll, _ = regexp.Compile(patTag + "|" + patText)
	imgPrefix = ""
	takeALink = ""
	lastChar = 'n'
	paraWritten = false
	tags = make(map[string]TagToFunc)
	tags["b"] = TagToFunc{0, do_b}
	tags["/b"] = TagToFunc{0, do_s_b}
	tags["a"] = TagToFunc{1, do_a}
	tags["/a"] = TagToFunc{0, do_s_a}
	tags["sq"] = TagToFunc{0, do_sq}
	tags["sq"] = TagToFunc{0, do_s_sq}
	tags["i"] = TagToFunc{0, do_i}
	tags["/i"] = TagToFunc{0, do_s_i}
	tags["google"] = TagToFunc{0, do_google}
	tags["img"] = TagToFunc{1, do_img}
	tags["top"] = TagToFunc{3, do_top}
	tags["btop"] = TagToFunc{3, do_btop}
	tags["atop"] = TagToFunc{2, do_atop}
	tags["libsyn"] = TagToFunc{1, do_libsyn}
	tags["youtube"] = TagToFunc{1, do_youtube}
	tags["apple"] = TagToFunc{1, do_apple}
	tags["amazon"] = TagToFunc{1, do_amazon}
	tags["spotify"] = TagToFunc{1, do_spotify}
	tags["section"] = TagToFunc{1, do_section}
	tags["/section"] = TagToFunc{0, do_s_section}
	tags["linkframe"] = TagToFunc{0, do_linkframe}
	tags["/linkframe"] = TagToFunc{0, do_s_linkframe}
	tags["indent"] = TagToFunc{0, do_indent}
	tags["/indent"] = TagToFunc{0, do_s_indent}
	tags["hr"] = TagToFunc{0, do_hr}
	tags["include"] = TagToFunc{1, do_include}
	tags["imgprefix"] = TagToFunc{1, do_imgprefix}
	tags["takealink"] = TagToFunc{1, do_takealink}
}

func processFile(filename string, mainDocument bool) {
	var waitForEndComment bool
	var waitForEndMain bool

	waitForEndComment = false
	waitForEndMain = false

	content, err := os.ReadFile(filename)
	if err != nil {
		handleError(err)
	}
	matches := reAll.FindAllString(string(content), -1)

	for _, match := range matches {
		if match == "[/main]" {
			waitForEndMain = false
			continue
		}
		if waitForEndMain {
			continue
		}
		if waitForEndComment {
			if match == "[/comment]" {
				waitForEndComment = false
			}
			continue
		}
		if reTag.MatchString(match) {
			paraWritten = false
			lastChar = '*'
			if match == "[comment]" {
				waitForEndComment = true
				continue
			} else if match == "[main]" {
				if !mainDocument {
					waitForEndMain = true
				}
				continue
			} else {
				evalTag(match)
			}
		} else {
			writeEachCharacter(match)
		}
	}
}

func writeEachCharacter(text string) {
	var i int
	for i = 0; i < len(text); i++ {
		c := text[i]
		if c == '\r' {
			continue
		}
		if c != '\n' {
			fout.WriteString(string(c))
			paraWritten = false
			lastChar = c
			continue
		}
		if lastChar == '\n' {
			if !paraWritten {
				fout.WriteString("<p>\n")
				paraWritten = true
				continue
			}
		}
		fout.WriteString(string(c))
		lastChar = c
		paraWritten = false
	}
}

func evalTag(tag string) {
	var i int
	blanks := []string{}
	inner := tag[1 : len(tag)-1]
	for i = 0; i < len(inner); i++ {
		if inner[i] == ' ' {
			break
		}
	}
	if i == len(inner) {
		t, ok := tags[inner]
		if !ok {
			log.Fatal("Tag ", inner, " not found in tag list.")
			os.Exit(1)
		}
		if t.Params == 0 {
			t.Fn(blanks)
		} else {
			log.Fatal("Invalid number of parameters to tag:", inner, " Expecting ", t.Params, " found none.")
			os.Exit(1)
		}
	} else {
		name := inner[0:i]
		args := strings.Split(inner[i+1:], "|")
		t, ok := tags[name]
		if !ok {
			log.Fatal("Tag ", name, " not found in tag list.")
			os.Exit(1)
		}
		if t.Params == len(args) {
			t.Fn(args)
		} else {
			log.Fatal("Invalid number of parameters to tag:", inner, " Expecting ", t.Params, " found ", len(args), ".")
			os.Exit(1)
		}
	}
}

func handleError(e error) {
	if e != nil {
		log.Fatal(e)
		os.Exit(1)
	}
}

func do_b(args []string) {
	fout.WriteString("<strong>")
}
func do_s_b(args []string) {
	fout.WriteString("</strong>")
}
func do_a(args []string) {
	var pos int
	pos = 8
	if args[0][pos] == '/' {
		pos = 9
	}
	if takeALink == "" {
		fout.WriteString("<a href=\"" + args[0] + "\">")
	} else {
		fout.WriteString("<a href=\"" + takeALink + args[0][pos:] + "\">")
	}
}
func do_s_a(args []string) {
	fout.WriteString("</a>")
}
func do_sq(args []string) {
	fout.WriteString("[")
}
func do_s_sq(args []string) {
	fout.WriteString("]")
}
func do_i(args []string) {
	fout.WriteString("<em>")
}
func do_s_i(args []string) {
	fout.WriteString("</em>")
}
func do_google(args []string) {
	fout.WriteString("<span class=\"chili-google\"><span style=\"color: blue; \">G</span><span style=\"color: red;\">o</span><span style=\"color: yellow;\">o</span><span style=\"color: blue;\">g</span><span style=\"color: green;\">l</span><span style=\"color: red;\">e</span></span>")

}
func do_img(args []string) {
	fout.WriteString("<div style=\"text-align: center;\"><img src=\"" + imgPrefix + args[0] + "\" class=\"chili-img\"><p></div>\n")
}
func do_top(args []string) {
	fout.WriteString("<div class=\"chili-marquis\"><div style=\"margin-left: 1em; margin-right: 1em;\"><p>")
	fout.WriteString("<pre>\n\n</pre><h1 class=\"chili-h1-inverse\">" + args[0] + "</h1><p><pre>\n\n</pre><img src=\"" + imgPrefix + args[1] + "\" class=\"chili-img\"><p><pre>\n\n</pre>")
	fout.WriteString("<span style=\"font-family: Courier New; font-size: 14px;\">Posted on " + args[2] + "</span><p><pre>\n\n</pre></div></div>\n")
}

func do_btop(args []string) {
	fout.WriteString("<div class=\"chili-marquis\"><div style=\"margin-left: 1em; margin-right: 1em;\"><p>")
	fout.WriteString("<pre>\n\n</pre><h1 class=\"chili-h1-inverse\">" + args[0] + "</h1><p><img src=\"" + imgPrefix + args[1] + "\" class=\"chili-img\"><p>")
	fout.WriteString("<span style=\"font-family: Courier New; font-size: 14px;\">" + args[2] + "</span><p><pre>\n\n</pre></div></div>\n")
}

func do_atop(args []string) {
	fout.WriteString("<div class=\"chili-marquis\"><div style=\"margin-left: 1em; margin-right: 1em;\"><p>")
	fout.WriteString("<pre>\n\n</pre><h1 class=\"chili-h1-inverse\">" + args[0] + "</h1><p>")
	if args[1] != "*" {
		fout.WriteString("<span style=\"font-family: Courier New; font-size: 14px;\">Posted on " + args[1] + "</span>\n")
	}
	fout.WriteString("<p><pre>\n\n</pre></div></div>\n")
}

func do_libsyn(args []string) {
	fout.WriteString("<a href=\"" + args[0] + "\">LibSyn Link</a><p>\n")
}
func do_youtube(args []string) {
	fout.WriteString("<a href=\"" + args[0] + "\">YouTube Link</a><p>\n")
}
func do_apple(args []string) {
	fout.WriteString("<a href=\"" + args[0] + "\">Apple Podcasts Link</a><p>\n")
}
func do_amazon(args []string) {
	fout.WriteString("<a href=\"" + args[0] + "\">Amazon Music Link</a><p>\n")
}
func do_spotify(args []string) {
	fout.WriteString("<a href=\"" + args[0] + "\">Spotify Link</a><p>\n")
}
func do_section(args []string) {
	fout.WriteString("<p><div class=\"chili-inset\">")
	fout.WriteString("<div style=\"margin-left: 1em; margin-right: 1em;\"><p>")
	fout.WriteString("<h2>" + args[0] + "</h2><p>")
}
func do_s_section(args []string) {
	fout.WriteString("<p><pre>\n\n</pre></div></div>\n")
}
func do_linkframe(args []string) {
	fout.WriteString("<div style=\"border-radius: 5px; font-family: border-width: thin; border-style: solid;\"><p>")
	fout.WriteString("<div style=\"margin-left: 1EM;\">")
}
func do_s_linkframe(args []string) {
	fout.WriteString("</div></div>")
}
func do_indent(args []string) {
	fout.WriteString("<div style=\"margin-left: 1em;\">\n")
}
func do_s_indent(args []string) {
	fout.WriteString("</div>\n")
}
func do_hr(args []string) {
	fout.WriteString("<hr>\n")
}
func do_include(args []string) {
	processFile(args[0], false)
}
func do_imgprefix(args []string) {
	imgPrefix = args[0]
}
func do_takealink(args []string) {
	takeALink = args[0]
}
