package chrome

import (
	"context"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"time"

	"screen-shot-service/config"
	"screen-shot-service/generator"
	"screen-shot-service/logger"

	gover "github.com/mcuadros/go-version"
	"github.com/pkg/errors"
)

// Chrome contains information about a Google Chrome
// instance, with methods to run on it.
type Chrome struct {
	Resolution       string
	ChromeTimeout    int
	ChromeTimeBudget int
	Path             string
	UserAgent        string
	Argvs            []string

	ScreenshotPath string
}

// NewChromeGenerator :
func NewChromeGenerator(c config.Config) generator.IGenerator {
	chrm := &Chrome{
		Resolution: c.Rresolution,

		ChromeTimeout:    c.ChromeTimeOut,
		ChromeTimeBudget: c.ChromeTimeBudget,
	}
	chrm.Setup()

	return chrm
}

// Setup configures a Chrome struct with the path
// specified to what is available on this system.
func (chrome *Chrome) Setup() {
	chrome.chromeLocator()
}

// ChromeLocator looks for an installation of Google Chrome
// and returns the path to where the installation was found
func (chrome *Chrome) chromeLocator() {

	// if we already have a path to chrome (say from a cli flag),
	// check that it exists. If not, continue with the finder logic.
	if _, err := os.Stat(chrome.Path); os.IsNotExist(err) {
		logger.Debugf("Chrome path not set or invalid. Performing search %v", logger.WithFields(map[string]interface{}{"user-path": chrome.Path, "error": err}))
	} else {
		logger.Debug("Chrome path exists, skipping search and version check")
		return
	}

	// Possible paths for Google Chrome or chromium to be at.
	paths := []string{
		"/usr/bin/chromium",
		"/usr/bin/chromium-browser",
		"/usr/bin/google-chrome-stable",
		"/usr/bin/google-chrome",
		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
		"/Applications/Google Chrome Canary.app/Contents/MacOS/Google Chrome Canary",
		"/Applications/Chromium.app/Contents/MacOS/Chromium",
		"C:/Program Files (x86)/Google/Chrome/Application/chrome.exe",
	}

	for _, path := range paths {

		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		logger.Debugf("Google Chrome path %v", logger.WithFields(map[string]interface{}{"chrome-path": chrome.Path}))

		chrome.Path = path

		// check the version for this chrome instance. if the current
		// path is a version that is old enough, use that.
		if chrome.checkVersion("60") {
			break
		}
	}

	// final check to ensure we actually found chrome
	if chrome.Path == "" {
		logger.Fatal("Unable to locate a valid installation of Chrome to use. gowitness needs at least Chrome \n Chrome Canary v60+. Either install Google Chrome or try specifying a valid location with \n  the --chrome-path flag")
	}
}

// checkVersion checks if the version at the chrome.Path is at
// least the lowest version
func (chrome *Chrome) checkVersion(lowestVersion string) bool {
	out, err := exec.Command(chrome.Path, "-version").Output()
	if err != nil {
		logger.Errorf("An error occurred while trying to get the Chrome version %v", logger.WithFields(map[string]interface{}{"chrome-path": chrome.Path, "err": err}))
		return false
	}

	// Convert the output to a simple string
	version := string(out)

	re := regexp.MustCompile(`\d+(\.\d+)+`)
	match := re.FindStringSubmatch(version)
	if len(match) <= 0 {
		logger.Debugf("Unable to determine Chrome version %v", logger.WithFields(map[string]interface{}{"chrome-path": chrome.Path}))
		return false
	}

	// grab the first match in the version extraction
	version = match[0]

	if gover.Compare(version, lowestVersion, "<") {
		logger.Warnf("Chrome version is older than %v %v", lowestVersion, logger.WithFields(map[string]interface{}{"chrome-path": chrome.Path, "chromeversion": version}))

		return false
	}

	logger.Debugf("Chrome version %v", logger.WithFields(map[string]interface{}{"chrome-path": chrome.Path, "chromeversion": version}))
	return true
}

// SetScreenshotPath sets the path for screenshots
func (chrome *Chrome) SetScreenshotPath(p string) error {

	if _, err := os.Stat(p); os.IsNotExist(err) {
		return errors.New("Destination path does not exist")
	}

	logger.Debugf("Screenshot path %v", logger.WithFields(map[string]interface{}{"screenshot-path": p}))

	chrome.ScreenshotPath = p

	return nil
}

// ScreenshotURL takes a screenshot of a URL
func (chrome *Chrome) ScreenshotURL(u string, destination string) error {
	targetURL, err := url.ParseRequestURI(u)
	if err != nil {
		logger.Errorf("Parse Url Error %v", err)
		return err
	}

	logger.Debugf("Full path to screenshot save using Chrome %v", logger.WithFields(map[string]interface{}{"url": targetURL, "full-destination": destination}))

	// Start with the basic headless arguments
	var chromeArguments = []string{
		"--headless", "--disable-gpu", "--hide-scrollbars", "-no-sandbox",
		"--disable-crash-reporter",
		"--user-agent=" + chrome.UserAgent,
		"--window-size=" + chrome.Resolution, "--screenshot=" + destination,
		"--virtual-time-budget=" + strconv.Itoa(chrome.ChromeTimeBudget*1000),
	}

	// Append extra arguments
	if len(chrome.Argvs) > 0 {
		for _, a := range chrome.Argvs {
			chromeArguments = append(chromeArguments, a)
		}
	}

	logger.Info(logger.WithFields(map[string]interface{}{"Chrome Arguments": chromeArguments}))

	// When we are running as root, chromiun will flag the 'cant
	// run as root' thing. Handle that case.
	if os.Geteuid() == 0 {
		logger.Debugf("Running as root, adding --no-sandbox %v", logger.WithFields(map[string]interface{}{"euid": os.Geteuid()}))
		chromeArguments = append(chromeArguments, "--no-sandbox")
	}

	// Check if we need to add a proxy hack for Chrome headless to
	// stfu about certificates :>
	if targetURL.Scheme == "https" {

		// Chrome headless... you suck. Proxy to the target
		// so that we can ignore SSL certificate issues.
		// proxy := shittyProxy{targetURL: targetURL}
		originalPath := targetURL.Path
		proxy := forwardingProxy{targetURL: targetURL}

		// Give the shitty proxy a few moments to start up.
		time.Sleep(500 * time.Millisecond)

		// Start the proxy and grab the listening port we should tell
		// Chrome to connect to.
		if err := proxy.start(); err != nil {
			logger.Errorf("Failed to start proxy for HTTPS request %v", logger.WithFields(map[string]interface{}{"error": err}))
			return err
		}

		// Update the URL scheme back to http, the proxy will handle the SSL
		proxyURL, _ := url.Parse("http://localhost:" + strconv.Itoa(proxy.port) + "/")
		proxyURL.Path = originalPath

		// I am not 100% sure if this does anything, but lets add --allow-insecure-localhost
		// anyways.
		chromeArguments = append(chromeArguments, "--allow-insecure-localhost")

		// set the URL to call to the proxy we are starting up
		chromeArguments = append(chromeArguments, proxyURL.String())

		// when we are done, stop the hack :|
		defer proxy.stop()

	} else {
		// Finally add the url to screenshot
		chromeArguments = append(chromeArguments, targetURL.String())
	}

	logger.Debugf("Google Chrome arguments %v", logger.WithFields(map[string]interface{}{"arguments": chromeArguments}))

	// get a context to run the command in
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(chrome.ChromeTimeout)*time.Second)
	defer cancel()

	// Prepare the command to run...
	cmd := exec.CommandContext(ctx, chrome.Path, chromeArguments...)

	logger.Infof("Taking screenshot %v", logger.WithFields(map[string]interface{}{"url": targetURL, "destination": destination}))

	// ... and run it!
	startTime := time.Now()
	if err := cmd.Start(); err != nil {
		logger.Fatal(err)
	}

	// Wait for the screenshot to finish and handle the error that may occur.
	if err := cmd.Wait(); err != nil {

		// If if this error was as a result of a timeout
		if ctx.Err() == context.DeadlineExceeded {
			logger.Errorf("Timeout reached while waiting for screenshot to finish %v",
				logger.WithFields(map[string]interface{}{"url": targetURL, "destination": destination, "err": err}))
			return ctx.Err()
		}

		logger.Errorf("Screenshot failed %v", logger.WithFields(map[string]interface{}{"url": targetURL, "destination": destination, "err": err}))
		return err
	}

	logger.Infof("Screenshot taken %v", logger.WithFields(map[string]interface{}{"url": targetURL, "destination": destination, "duration": time.Since(startTime)}))

	return nil
}
