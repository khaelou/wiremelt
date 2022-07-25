package helpers

// TwoCaptchaConnection references 2Captcha.com credentials
type TwoCaptchaConnection struct {
	AccessKey string
}

func ExtractReCaptchaV2SiteKey(targetURL string) (string, error) {
	// Download Target URL Source to /.temp
	// Extract V2 sitekey from downloaded source file
	// Remove /.temp and it's contents

	return "", nil
}

func ExtractReCaptchaV3SiteKey(targetURL string) (string, error) {
	// Download Target URL Source to /.temp
	// Extract V3 sitekey from downloaded source file
	// Remove /.temp and it's contents

	return "", nil
}
