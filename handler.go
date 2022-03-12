package urlshort

import (
	"bytes"
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		url, hasPath := pathsToUrls[r.RequestURI]
		if hasPath {
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
	return fn
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	yml = bytes.ReplaceAll(yml, []byte("\t"), []byte(" "))
	parsedYaml, err := parseYAML([]byte(yml))
	if err != nil {
		return nil, err
	}

	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

func parseYAML(yml []byte) ([]map[string]string, error) {
	var t []map[string]string
	err := yaml.Unmarshal(yml, &t)
	return t, err
}

func JSONHandler(json []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJSON, err := parseJSON([]byte(json))
	if err != nil {
		return nil, err
	}

	pathMap := buildMap(parsedJSON)
	return MapHandler(pathMap, fallback), nil
}

func parseJSON(jsn []byte) ([]map[string]string, error) {
	var t []map[string]string
	err := json.Unmarshal(jsn, &t)
	return t, err
}

func buildMap(parsed []map[string]string) map[string]string {
	pathMap := make(map[string]string)
	for _, value := range parsed {
		pathMap[value["path"]] = value["url"]
	}
	return pathMap
}
