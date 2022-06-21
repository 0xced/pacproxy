package pac

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// SmartLoader attempt to detect if we are using js, a url, or a file path
func SmartLoader(thing string) Loader {
	return func() (string, error) {
		if proxies, err := ParseFindProxyString(thing); err == nil {
			log.Print("loading pac as a static string result")
			return fmt.Sprintf("function FindProxyForURL(url, host){ return %q; }", proxies), nil
		}
		if strings.Contains(thing, "FindProxyForURL") && strings.Contains(thing, "{") {
			log.Print("loading pac as string")
			return thing, nil
		}
		if parseURL, parseErr := url.Parse(thing); parseErr == nil {
			switch strings.ToLower(parseURL.Scheme) {
			case "http", "https":
				return HTTPLoader(parseURL)()
			}
		}
		if _, registryErr := getRegistryValue(thing); registryErr == nil {
			return RegistryLoader(thing)()
		}
		return FileLoader(thing)()
	}
}

func FileLoader(file string) Loader {
	return func() (string, error) {
		log.Printf("loading pac from file %q", file)
		buf, err := ioutil.ReadFile(file)
		if err != nil {
			return "", err
		}
		return string(buf), nil
	}
}

func HTTPLoader(u *url.URL) Loader {
	return func() (string, error) {
		log.Printf("loading pac from URL %q", u)
		res, err := http.Get(u.String())
		if err != nil {
			return "", err
		}
		defer res.Body.Close()
		pac, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", err
		}
		return string(pac), nil
	}
}

func RegistryLoader(registryPath string) Loader {
	return func() (string, error) {
		log.Printf("loading pac URL from registry \"%s\"", registryPath)
		if registryURL, registryErr := getRegistryValue(registryPath); registryErr == nil {
			if parseURL, parseErr := url.Parse(registryURL); parseErr == nil {
				return HTTPLoader(parseURL)()
			} else {
				return "", parseErr
			}
		} else {
			return "", registryErr
		}
	}
}
