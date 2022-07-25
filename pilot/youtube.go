package pilot

import (
	"crypto/md5"
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/launcher/flags"
	"github.com/go-rod/stealth"
)

func YouTubeCast() {
	screenshots := "screenshots/"

	targetURL := "https://www.youtube.com/watch?v=N_L6ZLrS2Dc&autoplay=1"

	l := launcher.New().Headless(false)              // Create a browser launcher
	l = l.Set(flags.ProxyServer, "p.webshare.io:80") // Pass '--proxy-server=127.0.0.1:8081' argument to the browser on launch
	controlURL, _ := l.Launch()                      // Launch the browser and get debug URL

	browser := rod.New().ControlURL(controlURL).MustConnect().NoDefaultDevice() // Connect to the newly launched browser
	go browser.MustHandleAuth("vdwgrgse-rotate", "mu99wc82ns38")()              // <-- Notice how HandleAuth returns (mandatory goroutine)
	browser.MustIgnoreCertErrors(true)                                          // Ignore certificate errors since we are using local insecure proxy
	defer browser.MustClose()

	page := stealth.MustPage(browser).MustWindowFullscreen()
	//page.SetUserAgent(&proto.NetworkSetUserAgentOverride{UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.3"})
	//page.MustNavigate("https://bot.sannysoft.com")
	//UserMachineReport(page)

	page.MustNavigate(targetURL)

	fmt.Printf("stealth.JS: %x\n", md5.Sum([]byte(stealth.JS))) // You can also use stealth.JS directly without rod

	// ** PLAY IN @ Youtube.com **
	time.Sleep(7 * time.Second)
	page.MustElementX("/html/body/ytd-app/ytd-consent-bump-v2-lightbox/tp-yt-paper-dialog/div[4]/div[2]/div[6]/div[1]/ytd-button-renderer[2]/a/tp-yt-paper-button").MustClick()

	time.Sleep(3 * time.Second)
	page.MustElementX("/html/body/ytd-app/div[1]/ytd-page-manager/ytd-watch-flexy/div[5]/div[1]/div/div[1]/div/div/div/ytd-player/div/div/div[4]/button").MustClick()
	fmt.Println("Play!")

	time.Sleep(3 * time.Second)
	page.Eval(`() => document.querySelector("#movie_player > div.ytp-chrome-bottom > div.ytp-chrome-controls > div.ytp-left-controls > button").click()`)
	fmt.Println("Play2!")

	time.Sleep(3 * time.Second)
	page.MustElement("#movie_player > div.ytp-chrome-bottom > div.ytp-chrome-controls > div.ytp-left-controls > button").MustClick()
	fmt.Println("Play3!")

	time.Sleep(4 * time.Second)
	page.MustWaitLoad().MustScreenshot(fmt.Sprintf("%sYouTube.png", screenshots))
	fmt.Println("OK")

	time.Sleep(time.Hour)
}
