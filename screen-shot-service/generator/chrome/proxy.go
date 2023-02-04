package chrome

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"screen-shot-service/logger"
	"strings"
)

const listeningURL string = "127.0.0.1"

type forwardingProxy struct {
	targetURL *url.URL
	server    *httputil.ReverseProxy
	listener  net.Listener
	port      int
}

func (proxy *forwardingProxy) start() error {
	logger.Debugf("Initializing shitty forwarding proxy %v", logger.WithFields(map[string]interface{}{"target-url": proxy.targetURL}))

	// *Dont* verify remote certificates.
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Start the proxy and assign our custom Transport
	proxy.targetURL.Path = "/" // set the path to / as this becomes the base path
	proxy.server = httputil.NewSingleHostReverseProxy(proxy.targetURL)
	proxy.server.Transport = transport

	// Get an open port for this proxy instance to run on.
	var err error
	proxy.listener, err = net.Listen("tcp", listeningURL+":0")
	if err != nil {
		return err
	}

	// Set the port we used so that the caller of this method
	// can discover where to find this proxy instance.
	proxy.port = proxy.listener.Addr().(*net.TCPAddr).Port
	logger.Debugf("forwarding proxy listening port %v", logger.WithFields(map[string]interface{}{"target-url": proxy.targetURL, "listen-port": proxy.port}))

	// Finally, the goroutine for the proxy service.
	go func() {
		logger.Debugf("Starting shitty forwarding proxy goroutine %v", logger.WithFields(map[string]interface{}{"target-url": proxy.targetURL, "listen-address": proxy.listener.Addr()}))

		// Create an isolated ServeMux
		//  ref: https://golang.org/pkg/net/http/#ServeMux
		httpServer := http.NewServeMux()
		httpServer.HandleFunc("/", proxy.handle)

		if err := http.Serve(proxy.listener, httpServer); err != nil {

			// Probably a better way to handle these cases. Meh.
			if strings.Contains(err.Error(), "use of closed network connection") {
				return
			}

			// Looks like something is actually wrong
			logger.Errorf("Shitty forwarding proxy broke %v", logger.WithFields(map[string]interface{}{"err": err}))
		}
	}()

	return nil
}

// handle gets called on each request. We use this to update the host header.
func (proxy *forwardingProxy) handle(w http.ResponseWriter, r *http.Request) {
	logger.Debugf("Making proxied request %v", logger.WithFields(map[string]interface{}{"target-url": proxy.targetURL, "request": r.URL}))

	// Replace the host so that the Host: header is correct
	r.Host = proxy.targetURL.Host

	proxy.server.ServeHTTP(w, r)
}

// Stops the proxy
func (proxy *forwardingProxy) stop() {
	logger.Debugf("Stopping shitty forwarding proxy %v", logger.WithFields(map[string]interface{}{"target-url": proxy.targetURL, "port": proxy.port}))

	proxy.listener.Close()
}
