package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Unrecoverable error: %s\n", r)
		}
	}()

	mainTarget := extractMainTarget()
	targets := extractTargets()

	if mainTarget == nil && len(targets) == 0 {
		log.Println("No targets found, exiting...")
		return
	}

	label := "no"
	if mainTarget != nil {
		label = "a"
	}

	log.Printf("Proxy starting, got %s main target and %d other targets\n", label, len(targets))

	for i, t := range targets {
		log.Printf("[%d] %s\n", i, t.url)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		body, _ := io.ReadAll(r.Body)
		header := r.Header
		uri := r.RequestURI
		clientIP, _, _ := net.SplitHostPort(r.RemoteAddr)

		xff := header.Get("X-Forwarded-For")

		if xff != "" {
			xff = xff + ", "
		}

		xff = xff + clientIP

		header.Set("X-Forwarded-For", xff)
		header.Set("X-Real-IP", clientIP)

		log.Printf("%s %s\n", method, uri)

		forwardRequest(w, mainTarget, targets, method, uri, header, body)
	})

	listenerOpts := []string{
		extractListener(),
		"0.0.0.0:0",
		"127.0.0.1:0",
	}

	for _, listener := range listenerOpts {
		l, err := net.Listen("tcp", listener)

		if err != nil {
			if opErr, ok := err.(*net.OpError); ok {
				if syscallErr, ok := opErr.Err.(*os.SyscallError); ok {
					if syscallErr.Syscall == "bind" {
						log.Printf("Address %s already in use, trying next...\n", listener)

						continue
					}
				}
			}
			log.Println("An error occurred:", err)
			continue
		}

		defer l.Close()

		log.Printf("Listening on %s\n", l.Addr().String())

		if err := http.Serve(l, mux); err != nil {
			log.Println("An error occurred:", err)
		}

		break
	}
}

func forwardRequest(w http.ResponseWriter, mainTarget *target, targets []target, method string, uri string, header http.Header, body []byte) {
	if mainTarget != nil {
		forwardRequestToTarget(w, mainTarget, method, uri, header, body)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	for _, t := range targets {
		t := t

		go forwardRequest(nil, &t, nil, method, uri, header, body)
	}
}

func forwardRequestToTarget(w http.ResponseWriter, target *target, method string, uri string, header http.Header, body []byte) {
	req, err := http.NewRequest(method, target.url+"/"+strings.TrimPrefix(uri, "/"), bytes.NewReader(body))

	if err != nil {
		log.Println("An error occurred:", err)
	}

	req.Header = header

	res, err := target.client.Do(req)

	if w == nil {
		return
	}

	resBody, _ := io.ReadAll(res.Body)

	w.WriteHeader(res.StatusCode)

	for k, v := range res.Header {
		w.Header().Set(k, v[0])
	}

	w.Header()
	w.Write(resBody)
}

func extractMainTarget() *target {
	t := os.Getenv("TARGET_MAIN")

	if t == "" {
		return nil
	}

	u, err := url.Parse(t)

	if err != nil {
		log.Printf("Main target %s not a valid URL\n", t)

		return nil
	}

	return &target{url: strings.TrimSuffix(u.String(), "/")}

}

func extractTargets() []target {
	i := -1

	targets := make([]target, 0)

	for {
		i++

		t := os.Getenv("TARGET_" + strconv.Itoa(i))

		if t == "" {
			if i == 0 {
				// Allow 0 or 1-indexed targets
				continue
			}
			break
		}

		u, err := url.Parse(t)

		if err != nil {
			panic(err)
		}

		targets = append(targets, target{url: strings.TrimSuffix(u.String(), "/")})
	}

	return targets
}

func extractListener() string {
	p := os.Getenv("PORT")
	b := os.Getenv("BIND")

	if p == "" {
		p = "80"
	}

	return b + ":" + p
}

type target struct {
	url    string
	client http.Client
}
