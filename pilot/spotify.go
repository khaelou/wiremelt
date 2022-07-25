package pilot

import (
	"crypto/md5"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"wiremelt/twocaptcha"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/launcher/flags"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/stealth"
)

func SpotifySignin() {
	/*
		pinkcake@gmail.com:Mid0riba | Plan : Family | Country : US
		jon.chetrit@gmail.com:Statmotion87 | Plan : Family | Country : US
		lewis.albon@yahoo.co.nz:4Lancewood | Plan : Family | Country : NZ
		jo.adamson@xtra.co.nz:pumpkin1 | Plan : Family | Country : NZ
		jacquie@worklife.co.nz:reina699 | Plan : Family | Country : NZ
	*/
	screenshots := "screenshots/"

	loginURL := "https://accounts.spotify.com/en/login?continue=https%3A%2F%2Fopen.spotify.com%2Fartist%2F5TDJKVd91S3WHlgmPGibDy"

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

	page.MustNavigate(loginURL)
	fmt.Printf("stealth.JS: %x\n", md5.Sum([]byte(stealth.JS))) // You can also use stealth.JS directly without rod

	// ** SIGN IN @ Spotify.com **
	time.Sleep(3 * time.Second)

	page.MustElement("#login-username").MustInput("jon.chetrit@gmail.com") // Email
	page.MustElement("#login-password").MustInput("Statmotion87")          // Password

	time.Sleep(3 * time.Second)
	page.MustElement("#login-button").MustClick()

	time.Sleep(3 * time.Second)
	page.MustElement("#onetrust-close-btn-container > button").MustClick() // Hide onetrust-consent-sdk

	follow := page.MustElementX("/html/body/div[4]/div/div[2]/div[3]/div[1]/div[2]/div[2]/div/div/div[2]/main/section/div/div[2]/div[2]/div[4]/div/div/div/div/button[1]")
	followText := strings.ToLower(follow.MustText())
	if strings.Contains(followText, "follow") && !strings.Contains(followText, "following") {
		follow.MustClick()
		fmt.Println("Following Artist!")
	} else {
		fmt.Println("Bypass Follow!", followText)
	}

	time.Sleep(3 * time.Second)
	play := page.MustElementX("/html/body/div[4]/div/div[2]/div[3]/div[1]/div[2]/div[2]/div/div/div[2]/main/section/div/div[2]/div[2]/div[4]/div/div/div/div/div/button")
	play.MustClick()
	fmt.Println("Play!")

	time.Sleep(4 * time.Second)
	page.MustWaitLoad().MustScreenshot(fmt.Sprintf("%sSpotify.png", screenshots))
	fmt.Println("OK")

	time.Sleep(time.Hour)
}

func SpotifySignup() {
	screenshots := "screenshots/"

	twoCaptchaAPIKey := "06b1f801f4b0bcc0d1abea45e7306543" // *client.ReCaptcha -> connReCaptcha.AccessKey

	regURL := "https://www.spotify.com/us/signup?forward_url=https%3A%2F%2Fopen.spotify.com%2Fartist%2F5TDJKVd91S3WHlgmPGibDy%23login"
	urlV2CaptchaKey := "6LeO36obAAAAALSBZrY6RYM1hcAY7RLvpDDcJLy3" // Checkbox Key
	urlV3CaptchaKey := "6LfCVLAUAAAAALFwwRnnCJ12DalriUGbj8FW_J39" // Scoring Key

	l := launcher.New().Headless(true)               // Create a browser launcher
	l = l.Set(flags.ProxyServer, "p.webshare.io:80") // Pass '--proxy-server=127.0.0.1:8081' argument to the browser on launch
	controlURL, _ := l.Launch()                      // Launch the browser and get debug URL

	browser := rod.New().ControlURL(controlURL).MustConnect().NoDefaultDevice() // Connect to the newly launched browser
	go browser.MustHandleAuth("vdwgrgse-rotate", "mu99wc82ns38")()              // <-- Notice how HandleAuth returns (mandatory goroutine)
	browser.MustIgnoreCertErrors(true)                                          // Ignore certificate errors since we are using local insecure proxy
	defer browser.MustClose()

	page := stealth.MustPage(browser).MustWindowFullscreen()
	page.SetUserAgent(&proto.NetworkSetUserAgentOverride{UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.3"})
	page.MustNavigate("https://bot.sannysoft.com")
	UserMachineReport(page)

	page.MustNavigate(regURL)

	fmt.Printf("stealth.JS: %x\n", md5.Sum([]byte(stealth.JS))) // You can also use stealth.JS directly without rod

	// ** SIGN UP @ Spotify.com **
	time.Sleep(3 * time.Second)
	page.MustElement("#onetrust-close-btn-container > button").MustClick() // Hide onetrust-consent-sdk

	page.MustElement("#email").MustInput("test@wiremelt.app") // Email
	page.MustElement("#confirm").MustInput("test@wiremelt.app")
	page.MustElement("#password").MustInput("f33t1231")         // Password
	page.MustElement("#displayname").MustInput("test@wiremelt") // Display Name
	page.MustElement("#month").MustSelect("July")               // Birth Month
	page.MustElement("#day").MustInput("1")                     // Birth Day
	page.MustElement("#year").MustInput("2000")                 // Birth Year

	// Random Gender Selection
	minGender := 1
	maxGender := 3
	randomValGender := rand.Intn(maxGender-minGender+1) + minGender

	switch randomValGender {
	case 1:
		page.MustElement("#__next > main > div > div > form > fieldset > div > div:nth-child(1) > label > span.Type__TypeElement-goli3j-0.lfGOlT.TextForLabel-sc-1wen0a8-0.gXvRBb").MustClick() // Male
	case 2:
		page.MustElement("#__next > main > div > div > form > fieldset > div > div:nth-child(2) > label > span.Type__TypeElement-goli3j-0.lfGOlT.TextForLabel-sc-1wen0a8-0.gXvRBb").MustClick() // Female
	default:
		page.MustElement("#__next > main > div > div > form > fieldset > div > div:nth-child(3) > label > span.Type__TypeElement-goli3j-0.lfGOlT.TextForLabel-sc-1wen0a8-0.gXvRBb").MustClick() // Non-binary
	}

	// V2Solver
	log.Println("2Captcha.com")
	twoCaptcha := twocaptcha.New(twoCaptchaAPIKey)

	solvedV2, err := twoCaptcha.SolveRecaptchaV2(regURL, urlV2CaptchaKey)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("V2 Solved:", solvedV2)
	}

	grespV2 := page.MustElement("#g-recaptcha-response")
	grespV2.MustEval(`() => this.style.removeProperty('display');`)
	grespV2.MustEval(`() => this.innerHTML='` + solvedV2 + `';`)

	// V3Solver
	solvedV3, err := twoCaptcha.SolveRecaptchaV3(regURL, urlV3CaptchaKey, "t", "0.3")
	if err != nil {
		panic(err)
	} else {
		fmt.Println("V3 Solved:", solvedV3)
	}

	grespV3 := page.MustElement("#g-recaptcha-response-100000")
	grespV3.MustEval(`() => this.style.removeProperty('display')`)
	grespV3.MustEval(`() => this.innerHTML='` + solvedV3 + `';`)

	page.MustEval("() => window.findRecaptchaClients = function(){return'undefined'!=typeof ___grecaptcha_cfg?Object.entries(___grecaptcha_cfg.clients).map((([e,c])=>{const t={id:e,version:e>=1e4?'V3':'V2'};return Object.entries(c).filter((([e,c])=>c&&'object'==typeof c)).forEach((([c,n])=>{const i=Object.entries(n).find((([e,c])=>c&&'object'==typeof c&&'sitekey'in c&&'size'in c));if('object'==typeof n&&n instanceof HTMLElement&&'DIV'===n.tagName&&(t.pageurl=n.baseURI),i){const[n,a]=i;t.sitekey=a.sitekey;const o='V2'===t.version?'callback':'promise-callback',s=a[o];if(s){t.function=s;const i=[e,c,n,o].map((e=>'['+e+']')).join('');t.callback='___grecaptcha_cfg.clients'+i}else t.callback=null,t.function=null}})),t})):[]}") // Modify the content
	retrieveCallback := page.MustEval(`() => window.findRecaptchaClients();`).Arr()                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        // .Str()                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            // .Str()
	callbackElem := retrieveCallback[0].Get("callback").Str()
	callbackObj := callbackElem[29 : len(callbackElem)-14] // ___grecaptcha_cfg.clients['0']['?']['?']['callback']
	callbackV2 := fmt.Sprintf(".%s.%s.", callbackObj, callbackObj)
	fmt.Println("OBJ:", callbackObj)
	fmt.Println("Callback:", callbackV2)

	page.MustEval(`() => window.execCallback = function(){return ___grecaptcha_cfg.clients[0]` + callbackV2 + `callback('` + solvedV2 + `');}`) // ___grecaptcha_cfg.clients['0']['X']['X']['callback']
	execCallback := page.MustEval(`() => window.execCallback`).Val()
	fmt.Println("Exec Callback:", execCallback)
	if execCallback != nil {
		fmt.Printf("[✓] solved: ___grecaptcha_cfg.clients[0]%scallback(TOKEN)\n", callbackV2)
	} else {
		log.Fatalf("[x] execution error: ___grecaptcha_cfg.clients[0]%scallback(TOKEN)\n", callbackV2)
	}

	time.Sleep(7 * time.Second)
	page.MustElement("#__next > main > div > div > form > div.EmailForm__Center-jwtojv-0.itqiSk > div > button").MustClick()

	time.Sleep(3 * time.Second)
	page.MustWaitLoad().MustScreenshot(fmt.Sprintf("%sSpotify.png", screenshots))
	fmt.Println("OK")

	time.Sleep(time.Hour)
}
