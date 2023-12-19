package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"golang.org/x/exp/slices"
)

var (
	pluginName        = "krakend-plugin-jwe"
	HandlerRegisterer = registerer(pluginName)
	secret            = os.Getenv("CORE_JWE_SHARED_SECRET")
	symmetricKey      = os.Getenv("CORE_JWE_SYMMETRIC_KEY")
)

type registerer string

func (r registerer) RegisterHandlers(f func(
	name string,
	handler func(context.Context, map[string]interface{}, http.Handler) (http.Handler, error),
)) {
	f(string(r), r.registerHandlers)
}

func (r registerer) registerHandlers(_ context.Context, extra map[string]interface{}, h http.Handler) (http.Handler, error) {
	config, ok := extra[pluginName].(map[string]interface{})
	if !ok {
		return h, errors.New("configuration not found")
	}

	paths := getProtectedPaths(config)
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if !slices.Contains(paths, req.URL.Path) {
			h.ServeHTTP(w, req)
			return
		}

		token := strings.Split(req.Header.Get("Authorization"), " ")[1]
		decrypted, err := decryptAndDecode(&token, symmetricKey, secret)
		if err != nil {
			logger.Debug("JWT validation failed: %q", err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		logger.Debug("Allowed request from user uuid ", decrypted["uuid"])
		req.Header.Set("X-User-Id", decrypted["uuid"].(string))
		h.ServeHTTP(w, req)
	}), nil
}

func getProtectedPaths(config map[string]interface{}) []string {
	pathsIn := config["paths"].([]interface{})
	paths := make([]string, len(pathsIn))
	for i, path := range pathsIn {
		paths[i] = path.(string)
		logger.Debug(fmt.Sprintf("[PLUGIN: %s] Path %s will be JWE-protected", HandlerRegisterer, path))
	}
	return paths
}

func main() {}
