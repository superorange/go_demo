package mit

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"github.com/elazarl/goproxy"
	"io"
	"log"
	"net"
	"net/http"
)

func orPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func Run() {
	//verbose := flag.Bool("v", true, "should every proxy request be logged to stdout")
	httpAddr := flag.String("httpaddr", ":8080", "proxy http listen address")
	flag.Parse()

	proxy := goproxy.NewProxyHttpServer()
	//proxy.Verbose = *verbose
	if proxy.Verbose {
		log.Printf("Server starting up! - configured to listen on http interface %s ", *httpAddr)
	}

	proxy.NonproxyHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Host == "" {
			log.Println(w, "Cannot handle requests without Host header, e.g., HTTP 1.0")
			return
		}
		req.URL.Scheme = "http"
		req.URL.Host = req.Host
		proxy.ServeHTTP(w, req)
	})
	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		all, _ := io.ReadAll(resp.Body)
		fmt.Printf("接收到数据：%v\n", string(all))
		resp.Body = io.NopCloser(bytes.NewReader(all))
		return resp
	})
	//proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	//
	//	return req, nil
	//})
	//proxy.OnRequest(goproxy.ReqHostMatches(regexp.MustCompile("^.*:80$"))).
	//	HijackConnect(func(req *http.Request, client net.Conn, ctx *goproxy.ProxyCtx) {
	//		defer func() {
	//			if e := recover(); e != nil {
	//				ctx.Logf("error connecting to remote: %v", e)
	//				client.Write([]byte("HTTP/1.1 500 Cannot reach destination\r\n\r\n"))
	//			}
	//			client.Close()
	//		}()
	//		clientBuf := bufio.NewReadWriter(bufio.NewReader(client), bufio.NewWriter(client))
	//
	//		remote, err := connectDial(req.Context(), proxy, "tcp", req.URL.Host)
	//		orPanic(err)
	//		remoteBuf := bufio.NewReadWriter(bufio.NewReader(remote), bufio.NewWriter(remote))
	//		for {
	//			req, err := http.ReadRequest(clientBuf.Reader)
	//			orPanic(err)
	//			orPanic(req.Write(remoteBuf))
	//			orPanic(remoteBuf.Flush())
	//			resp, err := http.ReadResponse(remoteBuf.Reader, req)
	//			orPanic(err)
	//			orPanic(resp.Write(clientBuf.Writer))
	//			orPanic(clientBuf.Flush())
	//		}
	//	})

	log.Fatalln(http.ListenAndServe(*httpAddr, proxy))
}

// copied/converted from https.go
func dial(ctx context.Context, proxy *goproxy.ProxyHttpServer, network, addr string) (c net.Conn, err error) {
	if proxy.Tr.DialContext != nil {
		return proxy.Tr.DialContext(ctx, network, addr)
	}
	var d net.Dialer
	return d.DialContext(ctx, network, addr)
}

// copied/converted from https.go
func connectDial(ctx context.Context, proxy *goproxy.ProxyHttpServer, network, addr string) (c net.Conn, err error) {
	if proxy.ConnectDial == nil {
		return dial(ctx, proxy, network, addr)
	}
	return proxy.ConnectDial(network, addr)
}

type dumbResponseWriter struct {
	net.Conn
}

func (dumb dumbResponseWriter) Header() http.Header {
	panic("Header() should not be called on this ResponseWriter")
}

func (dumb dumbResponseWriter) Write(buf []byte) (int, error) {
	if bytes.Equal(buf, []byte("HTTP/1.0 200 OK\r\n\r\n")) {
		return len(buf), nil // throw away the HTTP OK response from the faux CONNECT request
	}
	return dumb.Conn.Write(buf)
}

func (dumb dumbResponseWriter) WriteHeader(code int) {
	panic("WriteHeader() should not be called on this ResponseWriter")
}

func (dumb dumbResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return dumb, bufio.NewReadWriter(bufio.NewReader(dumb), bufio.NewWriter(dumb)), nil
}
