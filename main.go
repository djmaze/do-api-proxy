package main

import (
  "fmt"
  "log"
  "net/http"
  "net/http/httputil"
  "net/url"
  "strings"
)

type Message struct {
  Droplet_ID int64
  Type string
}

var config EnvConfig

// Serve a reverse proxy for a given url
func serveReverseProxy(target string, res http.ResponseWriter, req *http.Request) {
	// parse the url
	url, _ := url.Parse(target)

	// create the reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Update the headers to allow for SSL redirection
	req.URL.Host = url.Host
	req.URL.Scheme = url.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = url.Host

	// Note that ServeHttp is non blocking and uses a go routine under the hood
	proxy.ServeHTTP(res, req)
}

func isReqAllowed(req *http.Request) bool {
  path := req.URL.Path

  if !config.IsRequestAuthCorrect(req) {
    log.Println("Request auth is incorrect")
    return false
  }

  if strings.HasPrefix(path, "/v2/volumes") {
    parts := strings.Split(path, "/")

    // allow volume list
    if len(parts) == 3 {
      log.Printf("Allowing volume list")
      return true
    }

    // allow volume action if volume is in environment
    volume := parts[3]
    for _, v := range config.Volumes {
      if v.Id == volume {
        log.Printf("Volume %s allowed\n", v.Name)
        return true
      }
    }

    log.Printf("Volume %s forbidden\n", volume)
    return false
  }

  // Allow viewing all actions for now
  if strings.HasPrefix(path, "/v2/actions") {
    log.Printf("Allow viewing action(s)")
    return true
  }

  log.Printf("Request does not match\n")
  return false
}

func handleRequestAndRedirect(res http.ResponseWriter, req *http.Request) {
  if !isReqAllowed(req) {
    res.WriteHeader(405)
    fmt.Fprintf(res, "not allowed\n")
    return
  }

  log.Printf("Proxying request %s\n", req.URL.Path)
  req.Header.Set("Authorization", "Bearer " + config.TargetApiToken)
  serveReverseProxy(config.TargetApiURL, res, req)
}

func main() {
  config = GetEnvConfig()

	// start server
	http.HandleFunc("/", handleRequestAndRedirect)
  log.Printf("Listening on port %s\n", config.ListenPort)
  if err := http.ListenAndServe(":" + config.ListenPort, nil); err != nil {
		panic(err)
	}
}
