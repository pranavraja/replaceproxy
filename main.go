package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/elazarl/goproxy"
	"github.com/pranavraja/replaceproxy/termcolor"
)

func ProxyThatModifiesResponsesFromURL(url string, modifier func(*http.Response, *goproxy.ProxyCtx) (newResponse *http.Response, newUrl string)) *goproxy.ProxyHttpServer {
	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest().DoFunc(func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		dump, _ := httputil.DumpRequest(ctx.Req, true)
		println(string(dump))
		return r, nil
	})
	proxy.OnResponse().DoFunc(func(r *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		var coloredStatusCode string
		status := fmt.Sprintf("[%d]", r.StatusCode)
		if url != "" && strings.HasPrefix(ctx.Req.URL.String(), url) {
			r, newUrl := modifier(r, ctx)
			coloredStatusCode = termcolor.ColoredWithBackground(status, termcolor.White, termcolor.BgMagenta, termcolor.Bold)
			fmt.Printf("%s %s %v => %s\n", coloredStatusCode, termcolor.Colored(ctx.Req.Method, termcolor.White, termcolor.Bold), ctx.Req.URL, newUrl)
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
	proxyThatReplacesResponses := ProxyThatModifiesResponsesFromURL(*repl, func(resp *http.Response, ctx *goproxy.ProxyCtx) (*http.Response, string) {
		resp.Body.Close()
		var newBody io.ReadCloser
		if strings.HasPrefix(*with, "http://") {
			replacedResp, err := http.Get(*with)
			if err != nil {
				panic(err)
			}
			resp.StatusCode = replacedResp.StatusCode
			resp.Header = replacedResp.Header
			newBody = replacedResp.Body
		} else {
			var err error
			newBody, err = os.Open(*with)
			if err != nil {
				panic(err)
			}
		}
		resp.Body = newBody
		return resp, *with
	})
	log.Fatal(http.ListenAndServe(":8080", proxyThatReplacesResponses))
}
