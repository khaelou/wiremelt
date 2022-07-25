package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"unicode"

	v8go "rogchap.com/v8go"
)

func CheckStringForEmptiness(input string) bool {
	if len(input) > 0 {
		return true
	} else {
		return false
	}
}

func CheckStringForIpOrHostname(host string) bool {
	addr := net.ParseIP(host)

	if addr == nil {
		return false
	} else {
		return true
	}
}

func IsStringValidUrl(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}

	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}

func CapitalizeString(s string) string {
	rs := []rune(s)
	inWord := false
	for i, r := range rs {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			if !inWord {
				rs[i] = unicode.ToTitle(r)
			}
			inWord = true
		} else {
			inWord = false
		}
	}
	return string(rs)
}

func LowercaseString(s string) string {
	rs := []rune(s)
	inWord := false
	for i, r := range rs {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			if !inWord {
				rs[i] = unicode.ToLower(r)
			}
			inWord = true
		} else {
			inWord = false
		}
	}
	return string(rs)
}

func YesNoToInt(b string) int {
	var store int

	switch b {
	case "Yes":
		store = 1
	case "No":
		store = 0
	}

	return store
}

func WriteToEnv(key, keyHasValue string) bool {
	envFile, err := ioutil.ReadFile(".env")
	if err != nil {
		// Create ENV file
		f, envInitErr := os.Create(".env")
		if envInitErr != nil {
			log.Fatalln(err, envInitErr)
		}
		defer f.Close()

		envFile, _ = ioutil.ReadFile(".env")
	}
	envLines := string(envFile)

	// Check for existing ENV key
	if strings.Contains(envLines, key) {
		lines := strings.Split(envLines, "\n")

		for i, line := range lines {
			if strings.Contains(line, key) {
				lines[i] = keyHasValue // Replace existing value
			}
		}
		output := strings.Join(lines, "\n")
		err := ioutil.WriteFile(".env", []byte(output), 0644)
		if err != nil {
			log.Fatalln(err)
		}
	} else { // Add key to ENV
		f, envInitErr := os.OpenFile(".env", os.O_APPEND|os.O_WRONLY, 0644)
		if envInitErr != nil {
			log.Fatalln(envInitErr)
		}
		defer f.Close()

		_, writeErr := f.WriteString(keyHasValue)
		if writeErr != nil {
			log.Fatalln(writeErr)
		}
	}

	return true
}

func V8Isolates(script string, showSource bool, isolateOpt ...*v8go.Isolate) (string, error) {
	var isolate *v8go.Isolate
	if len(isolateOpt) > 0 {
		isolate = isolateOpt[0]
	}

	ctx := v8go.NewContext(isolate) // Passing `nil` creates a new Isolate
	defer ctx.Close()
	output, e := ctx.RunScript(script, "macro.js")
	if e != nil {
		log.Fatalf("error V8Isolates: %+v\n", e)
		return "", e
	}

	product := output.String()

	if showSource {
		fmt.Println("JS SOURCE CODE:", script)
		fmt.Println("\nRETURN:", product)
	}

	return product, nil
}

func V8NodeJS(script string, showSource bool) (string, error) {
	_, err := exec.LookPath("node")
	if err != nil {
		log.Fatal(err)
		return fmt.Sprintf("nodeJS LookPath error: %v", err), err
	}

	output, e := exec.Command("node", "-e", script+"\nconsole.log(result)").Output()
	if e != nil {
		log.Fatalf("nodeJS error: %+v\n", e)
	}

	product := strings.TrimSuffix(string(output), "\n")

	if showSource {
		fmt.Println("JS SOURCE CODE:", script)
		fmt.Println("\nRETURN:", product)
	}

	return product, nil
}
