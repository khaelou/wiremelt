package pilot

import (
	"fmt"
	"log"
	"time"

	"wiremelt/twocaptcha"

	"github.com/go-rod/rod"
)

func InitGoogleDemo() {
	screenshots := "screenshots/"

	twoCaptchaAPIKey := "06b1f801f4b0bcc0d1abea45e7306543"

	originURL := "https://www.google.com/recaptcha/api2/demo"
	urlV2CaptchaKey := "6Le-wvkSAAAAAPBMRTvw0Q4Muexq9bi0DJwx_mJ-"

	browser := rod.New().MustConnect().NoDefaultDevice()
	page := browser.MustPage(originURL).MustWindowFullscreen()

	// ** DEMO @ Google.com **
	time.Sleep(3 * time.Second)

	// V2Solver
	log.Println("2Captcha.com")
	twoCaptcha := twocaptcha.New(twoCaptchaAPIKey)

	solvedV2, err := twoCaptcha.SolveRecaptchaV2(originURL, urlV2CaptchaKey)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("V2 Solved:", solvedV2)
	}

	grespV2 := page.MustElement("#g-recaptcha-response")
	grespV2.MustEval(`() => this.style.removeProperty('display');`)
	grespV2.MustEval(`() => this.innerHTML='` + solvedV2 + `';`)

	page.MustEval("() => window.findRecaptchaClients = function(){return'undefined'!=typeof ___grecaptcha_cfg?Object.entries(___grecaptcha_cfg.clients).map((([e,c])=>{const t={id:e,version:e>=1e4?'V3':'V2'};return Object.entries(c).filter((([e,c])=>c&&'object'==typeof c)).forEach((([c,n])=>{const i=Object.entries(n).find((([e,c])=>c&&'object'==typeof c&&'sitekey'in c&&'size'in c));if('object'==typeof n&&n instanceof HTMLElement&&'DIV'===n.tagName&&(t.pageurl=n.baseURI),i){const[n,a]=i;t.sitekey=a.sitekey;const o='V2'===t.version?'callback':'promise-callback',s=a[o];if(s){t.function=s;const i=[e,c,n,o].map((e=>'['+e+']')).join('');t.callback='___grecaptcha_cfg.clients'+i}else t.callback=null,t.function=null}})),t})):[]}") // Modify the content
	retrieveCallback := page.MustEval(`() => window.findRecaptchaClients();`).Arr()                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                        // .Str()                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                            // .Str()
	callbackElem := retrieveCallback[0].Get("callback").Str()
	callbackObj := callbackElem[29 : len(callbackElem)-14]
	callbackV2 := fmt.Sprintf(".%s.%s.", callbackObj, callbackObj)
	fmt.Println("OBJ:", callbackObj)
	fmt.Println("Callback:", callbackV2)

	page.MustEval(`() => window.execCallback = function(){return ___grecaptcha_cfg.clients[0]` + callbackV2 + `callback('` + solvedV2 + `');}`)
	execCallback := page.MustEval(`() => window.execCallback`).Val()
	fmt.Println("Exec Callback:", execCallback)
	if execCallback != nil {
		fmt.Printf("[✓] solved: ___grecaptcha_cfg.clients[0]%scallback(TOKEN)\n", callbackV2)
	} else {
		log.Fatalf("[x] execution error: ___grecaptcha_cfg.clients[0]%scallback(TOKEN)\n", callbackV2)
	}

	time.Sleep(7 * time.Second)
	page.MustElement("#recaptcha-demo-submit").MustClick()

	time.Sleep(3 * time.Second)
	page.MustWaitLoad().MustScreenshot(fmt.Sprintf("%sGoogle.png", screenshots))
	fmt.Println("OK")

	time.Sleep(3 * time.Second)
}
