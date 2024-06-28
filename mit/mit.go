package mit

import (
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/elazarl/goproxy"
)

func Run() {
	//verbose := flag.Bool("v", true, "should every proxy request be logged to stdout")
	httpAddr := flag.String("httpaddr", ":8089", "proxy http listen address")
	flag.Parse()

	proxy := goproxy.NewProxyHttpServer()
	//proxy.Verbose = *verbose
	log.Printf("Server starting up! - configured to listen on http interface %s ", *httpAddr)

	proxy.NonproxyHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.Host == "" {
			log.Println(w, "Cannot handle requests without Host header, e.g., HTTP 1.0")
			return
		}
		req.URL.Scheme = "http"
		req.URL.Host = req.Host
		proxy.ServeHTTP(w, req)
	})
	proxy.OnRequest().HandleConnect(&MyMitm{})
	proxy.OnResponse().DoFunc(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		all, _ := io.ReadAll(resp.Body)
		fmt.Printf("接收到数据：%v\n", string(all))
		resp.Body = io.NopCloser(bytes.NewReader(all))
		return resp
	})
	log.Fatalln(http.ListenAndServe(*httpAddr, proxy))
}

type MyMitm struct {
}

func (m MyMitm) HandleConnect(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {

	var ca = []byte(`-----BEGIN CERTIFICATE REQUEST-----
MIIC5jCCAc4CAQAwgaIxEjAQBgNVBAMMCWNoaW5hLmNvbTESMBAGA1UECgwJ55Sw
6ZSm5bKXMRIwEAYDVQQLDAnnlLDplKblspcxDDAKBgNVBAYTAz8/PzESMBAGA1UE
CAwJ55Sw6ZSm5bKXMRIwEAYDVQQHDAnnlLDplKblspcxEjAQBgNVBAkMCeeUsOmU
puWylzEaMBgGCSqGSIb3DQEJARYLPz8/QD8/Py5jb20wggEiMA0GCSqGSIb3DQEB
AQUAA4IBDwAwggEKAoIBAQC06cVErQKFK5tJIXWZSvAADDZcHWSvXq5WlZPf0Rlq
zIdN5+pGsIhA6h7jM/n4h5I16BofVdrFLBMEmgBlBVhMmpEaQZ2YM80VN1pIVJqr
+4qqZLHHXUi0K8eMmyj9OPODBBnQUO1+Pea5St03wkomdbFBrmp0I8j96/Ughpcg
KZCQZ9ObACSD6LCaYFBGjZWwC35evUg/z19a8iQN4CFi0qLMeVslkPw+NvD5xNlo
uQDu7ZYI00Ukq3RDnPsk5um3USx9IQGRqfYj4InF8CbGOOUGgnFYccLX3DZtCcRa
jjvGLQW0sr3p6NydiqnHh6c5B6z9pLyAdGwAqbRbAiSnAgMBAAEwDQYJKoZIhvcN
AQELBQADggEBAHd6rGoYtmWMyBM2JKzTkMYjB3/lSCavp8DlKfm5wvnK5WANp7A7
+KvFyoamHvkGWWHoJzozsxfQkSe8ZmJNw9QmXFwECtJNo89p7+mKuqac9DZPmX9N
I+NcvVGvUE1I4qdT1Ik06Ipoy6nE8MbLZuj7T7CokbfL/AzrVxkFRTeL/VJCngNe
yesKssGwFq6yaJp2VRu9pQzQyMAD5sShSQ0ckQ2PIpYDVMRGWAtsXpXCV4EVTXrJ
oHkjcbjh7kxq2xAlhaHfGySMdewhFTApaOPZr9fOrYHkO43IFZ5rQNStZSUXz18V
OfA4EubuqOcRFFuejcvVIqr7jGkmLNmTrio=
-----END CERTIFICATE REQUEST-----`)
	var key = []byte(`-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC06cVErQKFK5tJ
IXWZSvAADDZcHWSvXq5WlZPf0RlqzIdN5+pGsIhA6h7jM/n4h5I16BofVdrFLBME
mgBlBVhMmpEaQZ2YM80VN1pIVJqr+4qqZLHHXUi0K8eMmyj9OPODBBnQUO1+Pea5
St03wkomdbFBrmp0I8j96/UghpcgKZCQZ9ObACSD6LCaYFBGjZWwC35evUg/z19a
8iQN4CFi0qLMeVslkPw+NvD5xNlouQDu7ZYI00Ukq3RDnPsk5um3USx9IQGRqfYj
4InF8CbGOOUGgnFYccLX3DZtCcRajjvGLQW0sr3p6NydiqnHh6c5B6z9pLyAdGwA
qbRbAiSnAgMBAAECggEAUVmLNInth69nmNdcV418ZTEYoowvEbKsB1AkWfDfEoic
0PnXfWj1I+eC8xyUq15e7zGKyYtkH+RlOtz7D3H7Vhrs45ccw+uLLR6iUMMGeA5Y
uK0lyeWXAzlqdj3xDQi2azQFXYh+epVgMfLJjcCmcivbBJNm2AmdYBzhsXoD69gx
57bMme59I+VEXMLkdXc671dsSDIHy1fiRBfwq0fGSxXCpKFea2agJDT5XZFjZB4u
dKv0M7Qvo/OiZBud7YBsiUTVNt9/NnBSxFU4yx/yFWbVCFzdYpcdkDcAXXK2Jc1m
azccqLGTIBJ9Ym8M4yqFGm1DqcrTjOqTosUFUcV3LQKBgQDzFchibnV34IpIukMw
301YIKGAvcwC9FGhJVcXpMZKCdWUHYQHZkqeTDwsi9Qt6clR5p8tXAqP9c7vAMmo
i9QK6Mz1olMm7LcUe0G/E0k8yV3nB8d+p8GO9Z8ApoGwnLEuynXOih7f6R8sI10c
+EnzgWJX4M2PrIZN2jX++TeZNQKBgQC+hmIXlV//xfPK/2OPEigDcIdgquQBfFT9
2IgrYcdA9dx87jDEY3L4QMB1uszOYwzvtl8tAL+pYXFwGzDzRC+nelUXkjtFEHic
ZFYtLuCr9budkFrUQX/scrcrdc/QJxnhD6FXxSav1uBthfU/RJkgHi7xDk/LdCjn
/0JsH1ed6wKBgGCV7yDtMs+G8GslVz07/Mdfb8xvnXgvC6Az7f7/Aaq5bZuEfslR
46QyNlac2JnForBgKi8juy6oRKjCb14A8SfEiGuxK8jzlWsV7nG1gAwfFqiNdr7k
eQwMnDjt7+n02JH28Ag46TuerTwkcQLpxLh0WFcCg6rqqhKU+Y9uBqFhAoGAHx4z
x0ZOd6gNOYqc6DE+99DZS6CdvOBhwVQsaWl+8c02RfFhZbIYhROOW6w25z6mTkCr
Kt6Eb1XLAVRvmkv4vJHuc/seUxltmZ2JtbeCWpO4IPQC4cgQ7L2PzTlgx86bG3dC
EuPQfcfKwBixBbRejjBf2l9MCR7fz4SRhLdZyCECgYEA1Aw/qK1pcDDDtSekqo12
KvZYWIAME3i6cSrA/f39N313G0eRSML/RLaQ7tggXJWfWoqOjVuxo4+VJEhrY3qo
OrrRTVBnwzakmhf+88F9bpy3mzwPz3z6J8u0BWbLmCc949YmIS3zqI5P/ESJMWra
ZmET0SLFm0f5pfDstzEYmkw=
-----END PRIVATE KEY-----`)

	pair, err := tls.X509KeyPair(ca, key)
	if err != nil {
		fmt.Printf("error %v\n", err)
		return nil, ""
	}

	return &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(&pair)}, host
}
