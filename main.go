package main

import (
	"flag"
	"fmt"
	"github.com/elazarl/goproxy"
	"github.com/vrischmann/termcolor"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func ProxyThatModifiesResponsesFromURL(url string, modifier func(*http.Response, *goproxy.ProxyCtx) *http.Response) *goproxy.ProxyHttpServer {
	proxy := goproxy.NewProxyHttpServer()
	proxy.OnResponse().DoFunc(func(r *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		var coloredStatusCode string
		status := fmt.Sprintf("[%d]", r.StatusCode)
		if strings.HasPrefix(ctx.Req.URL.String(), url) {
			r = modifier(r, ctx)
			coloredStatusCode = termcolor.ColoredWithBackground(status, termcolor.White, termcolor.BgMagenta, termcolor.Bold)
			fmt.Printf("%s %s %v => %s\n", coloredStatusCode, termcolor.Colored(ctx.Req.Method, termcolor.White, termcolor.Bold), ctx.Req.URL, url)
			return r
		}
		switch r.StatusCode {
		case 301, 302:
			coloredStatusCode = termcolor.Colored(status, termcolor.Cyan)
		case 200:
			coloredStatusCode = termcolor.Colored(status, termcolor.White)
		case 404, 500:
			coloredStatusCode = termcolor.Colored(status, termcolor.Red)
		default:
			coloredStatusCode = termcolor.Colored(status, termcolor.Magenta)
		}
		fmt.Printf("%s %s %v\n", coloredStatusCode, termcolor.Colored(ctx.Req.Method, termcolor.White, termcolor.Bold), ctx.Req.URL)
		return r
	})
	return proxy
}

func main() {
	repl := flag.String("repl", "", "URL to replace")
	with := flag.String("with", "", "URL to substitute")
	flag.Parse()
	proxyThatReplacesResponses := ProxyThatModifiesResponsesFromURL(*repl, func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		resp.Body.Close()
		var newBody io.ReadCloser
		if strings.HasPrefix(*with, "http://") {
			replacedResp, err := http.Get(*with)
			if err != nil {
				panic(err)
			}
			newBody = replacedResp.Body
		} else {
			var err error
			newBody, err = os.Open(*with)
			if err != nil {
				panic(err)
			}
		}
		resp.Body = newBody
		return resp
	})
	log.Fatal(http.ListenAndServe(":8080", proxyThatReplacesResponses))
}
