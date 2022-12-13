package pilot

import (
	"fmt"
	"log"
	"time"

	"github.com/fatih/color"
	"github.com/go-rod/rod"
)

func InitPilot(include interface{}) {
	fmt.Println(color.HiBlueString(fmt.Sprintf("\n~ INIT_PILOT: %v", include)))

	/*
		page := rod.New().NoDefaultDevice().MustConnect().MustPage("https://www.khaelou.com")
		page.MustWindowFullscreen() // debug
		page.MustWaitLoad().MustScreenshot(fmt.Sprintf("/pilot/screenshots/%s.png", "khaelou"))
		time.Sleep(time.Second * 2) // debug
	*/

	BrowserNavigator("https://www.wikipedia.org/", "wikipedia")
	BrowserNavigator("https://www.apple.com/", "apple")
	BrowserNavigator("https://www.google.com/", "google")
}

// BrowserNavigator:
func BrowserNavigator(targetURL string, resultName string) {
	browser := rod.New().MustConnect().NoDefaultDevice()
	page := browser.MustPage(targetURL).MustWindowFullscreen()

	err := rod.Try(func() {
		page.MustElement("#searchInput").MustInput("limewire")
		page.MustElement("#search-form > fieldset > button").MustClick()
	})
	log.Println(err)

	page.MustWaitLoad().MustScreenshot(fmt.Sprintf("pilot/screenshots/%s.png", resultName))
	time.Sleep(time.Second * 4)
}
