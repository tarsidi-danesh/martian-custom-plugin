package main

import (
	"context"
	"errors"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
)

// HandlerRegisterer is the symbol the plugin loader will try to load. It must implement the Registerer interface
var HandlerRegisterer = handlerRegisterer("handler-plugin")

// ClientRegisterer is the symbol the plugin loader will try to load. It must implement the RegisterClient interface
var ClientRegisterer = proxyRegisterer("proxy-plugin")

type handlerRegisterer string
type proxyRegisterer string

func (r handlerRegisterer) RegisterHandlers(f func(
	name string,
	handler func(context.Context, map[string]interface{}, http.Handler) (http.Handler, error),
)) {
	f(string(r), r.registerHandlers)
}

func (r proxyRegisterer) RegisterClients(f func(
	name string,
	handler func(context.Context, map[string]interface{}) (http.Handler, error),
)) {
	f(string(r), r.registerClients)
}

func (r handlerRegisterer) registerHandlers(ctx context.Context, extra map[string]interface{}, handler http.Handler) (http.Handler, error) {
	// check the passed configuration and initialize the plugin
	name, ok := extra["name"].(string)
	if !ok {
		return nil, errors.New("wrong config")
	}
	if name != string(r) {
		return nil, fmt.Errorf("unknown register %s", name)
	}
	// return the actual handler wrapping or your custom logic so it can be used as a replacement for the default http handler
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("HostName:", req.Host)
		fmt.Printf("Url: %q \n", html.EscapeString(req.URL.Path))

		handler.ServeHTTP(w, req)
	}), nil
}

func (r proxyRegisterer) registerClients(ctx context.Context, extra map[string]interface{}) (h http.Handler, e error) {
	// check the passed configuration and initialize the plugin
	name, ok := extra["name"].(string)
	if !ok {
		return nil, errors.New("wrong config")
	}
	if name != string(r) {
		return nil, fmt.Errorf("unknown register %s", name)
	}
	// return the actual handler wrapping or your custom logic so it can be used as a replacement for the default http client
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Printf("before proxy handler %s\n", html.EscapeString(req.URL.Path))
		makeOriginalRequest(w, req)
		fmt.Println("after proxy-plugin called")
	}), nil
}

func makeOriginalRequest(w http.ResponseWriter, req *http.Request) {
	client := &http.Client{}
	req.Header["channelId"] = []string{"DESKTOP"}
	req.Header["storeId"] = []string{"TIKETCOM"}
	req.Header["username"] = []string{"username"}
	req.Header["requestId"] = []string{"requestId"}
	req.Header["serviceId"] = []string{"GATEWAY"}
	for h, v := range req.Header {
		fmt.Println(h, "=", v[0])
	}
	// Send an HTTP request and returns an HTTP response object.
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer resp.Body.Close()

	// headers
	for name, values := range resp.Header {
		w.Header()[name] = values
	}

	// status (must come after setting headers and before copying body)
	w.WriteHeader(resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	w.Write(body)

	fmt.Println("request completed")
}

func init() {
	fmt.Println("krakend-example handler plugin loaded!!!")
}

func main() {}
