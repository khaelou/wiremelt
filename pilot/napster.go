package pilot

import (
	"crypto/md5"
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/launcher/flags"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/stealth"
)

func NapsterSignin() {
	/*
		pauline-spagnolo@hotmail.fr:villadona
	*/

	screenshots := "screenshots/"

	loginURL := "https://web.napster.com/auth/login"

	l := launcher.New().Headless(false)                // Create a browser launcher
	l = l.Set(flags.ProxyServer, "45.137.60.112:6640") // Pass '--proxy-server=127.0.0.1:8081' argument to the browser on launch
	controlURL, _ := l.Launch()                        // Launch the browser and get debug URL

	browser := rod.New().ControlURL(controlURL).MustConnect().NoDefaultDevice() // Connect to the newly launched browser
	go browser.MustHandleAuth("vdwgrgse", "mu99wc82ns38")()                     // <-- Notice how HandleAuth returns (mandatory goroutine)
	browser.MustIgnoreCertErrors(true)                                          // Ignore certificate errors since we are using local insecure proxy
	defer browser.MustClose()

	page := stealth.MustPage(browser).MustWindowFullscreen()
	page.SetUserAgent(&proto.NetworkSetUserAgentOverride{UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.3"})
	page.MustNavigate("https://bot.sannysoft.com")
	UserMachineReport(page)

	page.MustNavigate(loginURL)

	fmt.Printf("stealth.JS: %x\n", md5.Sum([]byte(stealth.JS))) // You can also use stealth.JS directly without rod

	// ** SIGN IN @ Napster.com **
	time.Sleep(3 * time.Second)

	page.MustElement("#root > div > div > div > div > div > div > div > form > div:nth-child(1) > div.sc-jSgvzq.ggOWiZ.sc-xyEDr.ctOeOM > input").MustInput("pauline-spagnolo@hotmail.fr") // Email
	page.MustElement("#root > div > div > div > div > div > div > div > form > div:nth-child(2) > div.sc-jSgvzq.ggOWiZ.sc-xyEDr.ctOeOM > span > input").MustInput("villadona")            // Password

	time.Sleep(3 * time.Second)
	page.MustElement("#root > div > div > div > div > div > div > div > form > div.sc-jSgvzq.hAcLwY > button.sc-gsTEea.cbWaac").MustClick()

	time.Sleep(3 * time.Second)
	page.MustElement("#root > div > div > div.sc-jSgvzq.kSmGbF > div.sc-jSgvzq.sc-jfJyPD.fpndQW.ibaZra > div > div > div.sc-gyUflj.hGsnQZ > div > a:nth-child(4) > div > div > div").MustClick()

	time.Sleep(3 * time.Second)
	play := page.MustElementX("/html/body/div[1]/div/div/div[1]/div[2]/div/div/div/div[3]/div[1]/div[2]/div[3]/div[1]/button")
	play.MustClick()
	fmt.Println("Play!")

	time.Sleep(4 * time.Second)
	page.MustWaitLoad().MustScreenshot(fmt.Sprintf("%sNapster.png", screenshots))
	fmt.Println("OK")

	time.Sleep(time.Hour)
}
