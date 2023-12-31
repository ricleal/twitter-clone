// Package openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.15.0 DO NOT EDIT.
package openapi

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/9RWT2/bPgz9KgZ/v6MQp83Nt60bigA79JCeimBQbTpVYUuqRK0zAn/3QZLt/HOatiuG",
	"9BTFJMXHR/LZa8hVrZVESRayNdj8AWsejt+NUcYftFEaDQkMj3NVoP8t0OZGaBJKQhadk2BjUCpTc4IM",
	"hKTZJTCgRmP8iys00DKo0Vq+OnpRbx5CLRkhV9C2DAw+OWGwgOwOuoS9+7JlsHhGpDHYklAGw96VDETh",
	"Hw+onRPFYWYGzqL5+SrfA5Qx9+aK5RCi7h8xJ3/9rcURurHmotpJGZ+wd5cheY2jNHhwR4x7BQ2erENz",
	"WI8PEbJU4TJBlbctngURmuTLzRwY/EJjY88vJtPJ1CNQGiXXAjKYTaaTGTDQnB4CDyn5vobjKvbXs8T9",
	"1MwLyOAaaRE9PFSrlbSRwMvpdK//XOtK5CE0fbQeQT/3/iQI6xD4v8ESMvgv3WxI2q1HGoesHarmxvAm",
	"Fr07zj+EpUSVSYc+2EvuKnoTppegxDUdSe0k/taYExYJdj4MrKtrbpoeGa+qLWha2RFmb5TdpvbJoaWv",
	"qmg+rIKOzN0hI+OwPWjlxaFghOgkN8h9pdblOVpbuqpqzontq4Av4ZHtYOwmOl2Loj091vMibIPhNRIa",
	"C9ndOBHzb+AXD7KwOtCvOwQt2KWXbdV8StCWf7lVr+z/WEUFEhfVWe3ONVJsZHLfeMZDO70ovqhPt8Hh",
	"X8hTeJW8QZ0i9HMUpwHZcW3a0Prx0hSJfKcy+eDPI0wultoP8klZCrSfVqVAwicUpU3nR+rRRpWiwnOT",
	"JLeFbVAmLwN8Nf7CiJ9jV5WS2H2UdY3pTNAu2z8BAAD//+uIQUAiDAAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
