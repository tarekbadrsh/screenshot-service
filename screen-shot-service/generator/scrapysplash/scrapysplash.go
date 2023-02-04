package scrapysplash

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"screen-shot-service/config"
	"screen-shot-service/generator"
	"screen-shot-service/logger"
)

const defaultLuaScript = ` 
		-- main script
		function main(splash)
			-- set screen size
			splash:set_viewport_size(%v)
			assert(splash:go(splash.args.url))
			return splash:png{}
		end
`

// Splash :
type Splash struct {
	Resolution       string
	ScrapySplashHost string
	ScrapyURL        url.URL
	/* lua is the language used inside scrapy/splash to render the page

	- The default script
	```
		-- main script
		function main(splash)
			-- set screen size
			splash:set_viewport_size(${Resolution})
			assert(splash:go(splash.args.url))
			return splash:png{render_all=true}
		end
	```
	*/
	LuaSource string
}

// NewSplashGenerator :
func NewSplashGenerator(c config.Config) generator.IGenerator {
	splsh := &Splash{
		Resolution:       c.Rresolution,
		ScrapySplashHost: c.ScrapySplashHost,
		LuaSource:        fmt.Sprintf(defaultLuaScript, c.Rresolution),
	}

	return splsh
}

// Setup configures Splash struct
func createSplashURL(splsh *Splash, targetURL string) string {

	url := url.URL{Scheme: "http", Host: splsh.ScrapySplashHost, Path: "execute"}
	q := url.Query()
	q.Add("timeout", "90.0")
	q.Add("lua_source", splsh.LuaSource)
	q.Add("url", targetURL)
	url.RawQuery = q.Encode()

	return url.String()

}

// ScreenshotURL : takes a screenshot of a URL
func (splsh *Splash) ScreenshotURL(targetURL string, destination string) error {
	generatorURL := createSplashURL(splsh, targetURL)

	response, err := http.Get(generatorURL)
	if err != nil {
		logger.Error(err)
		return err
	}
	defer response.Body.Close()

	out, err := os.Create(destination)
	if err != nil {
		logger.Error(err)
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, response.Body)
	if err != nil {
		logger.Error(err)
		return err
	}

	logger.Infof("Generate screenshot succeed %v", destination)
	return nil
}
