package wiremelt

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"

	"wiremelt/helpers"
	"wiremelt/utils"
	"wiremelt/worker"

	"github.com/joho/godotenv"
	"github.com/manifoldco/promptui"
)

func WiremeltAscii() (interface{}, error) {
	fmt.Println(`
	:::       ::: ::::::::::: :::::::::  :::::::::: ::::    ::::  :::::::::: :::    ::::::::::: 
	:+:       :+:     :+:     :+:    :+: :+:        +:+:+: :+:+:+ :+:        :+:        :+:     
	+:+       +:+     +:+     +:+    +:+ +:+        +:+ +:+:+ +:+ +:+        +:+        +:+     
	+$+  +:+  +$+     +$+     +$++:++$:  +$++:++$   +$+  +:+  +$+ +$++:++$   +$+        +$+     
	+$+ +$+$+ +$+     +$+     +$+    +$+ +$+        +$+       +$+ +$+        +$+        +$+     
	 $+$+$ $+$+$      $+$     $+$    $+$ $+$        $+$       $+$ $+$        $+$        $+$     
	  $$$   $$$   $$$$$$$$$$$ $$$    $$$ $$$$$$$$$$ $$$       $$$ $$$$$$$$$$ $$$$$$$$$$ $$$     
	`)
	return nil, nil
}

// ClientConfiguration references setup details prior to session configuration, connects external infrastructure to client
type ClientConfiguration struct {
	DBConn    helpers.DatabaseConnection   `json:"dbConn"`
	ReCaptcha helpers.TwoCaptchaConnection `json:"reCaptcha"`
	Proxy     helpers.ProxyConnection      `json:"proxy"`
}

// Initiate New Client confiuration
func newClientConfig(connDB helpers.DatabaseConnection, connReCaptcha helpers.TwoCaptchaConnection, proxy helpers.ProxyConnection) *ClientConfiguration {
	config := &ClientConfiguration{
		DBConn:    connDB,
		ReCaptcha: connReCaptcha,
		Proxy:     proxy,
	}

	return config
}

// View contents of ClientConfiguration
func (config *ClientConfiguration) Read() {
	fmt.Printf("\n\tClient-Configuration: [DBConn] %v | [2Captcha] %v | [Proxy] %v", config.DBConn, config.ReCaptcha, config.Proxy)
	fmt.Println("")
}

// Parse takes ClientConfiguration as input, validates upon initiating actual client
func (config *ClientConfiguration) Parse() bool {
	if config != nil {
		return true
	} else {
		return false
	}
}

// SessionConfiguration references setup details prior to client initialization, can save or load from JSON
type SessionConfiguration struct {
	RepeatCycle     int              `json:"repeatCycle"`
	CPUCores        int              `json:"cpuCores"`
	FactoryQuantity int              `json:"factoryQuantity"`
	WorkerQuantity  int              `json:"workerQuantity"`
	JobsPerMacro    int              `json:"jobsPerMacro"`
	FactoryFocus    map[int]string   `json:"factoryFocus"`
	WorkerRoles     map[int]string   `json:"workerRoles"`
	MacroLibrary    worker.MacroSpec `json:"macroLibrary"`
	ShellCycle      int              `json:"shellCycle"`
	NeuralEnabled   int              `json:"neuralEnabled"`
	TrainLimit      int              `json:"trainLimit"`
}

// Initiate New Session confiuration
func NewSessionConfig(repeatCycle, cpuCores, factoryQuantity, workerQuantity int, jobsPerMacro int, factoryFocus map[int]string, workerRoles map[int]string, macroSpec worker.MacroSpec, shellCycle int, neuralEnabled int, trainLimit int) *SessionConfiguration {
	config := SessionConfiguration{
		RepeatCycle:     repeatCycle,
		CPUCores:        cpuCores,
		FactoryQuantity: factoryQuantity,
		WorkerQuantity:  workerQuantity,
		JobsPerMacro:    jobsPerMacro,
		FactoryFocus:    factoryFocus,
		WorkerRoles:     workerRoles,
		MacroLibrary:    macroSpec,
		ShellCycle:      shellCycle,
		NeuralEnabled:   neuralEnabled,
		TrainLimit:      trainLimit,
	}

	return &config
}

// View contents of SessionConfiguration
func (config *SessionConfiguration) Read() {
	fmt.Printf("\n\tSession-Configuration: %v", *config)
	fmt.Println("")
}

// Parse takes SessionConfiguration as input, validates upon initiating into a Session
func (config *SessionConfiguration) Parse() bool {
	if config != nil {
		return true
	} else {
		return false
	}
}

// Session Initialization
func (session *SessionConfiguration) StartSession(client *ClientConfiguration) {
	defer session.SessionCleanup()

	runtime.GOMAXPROCS(session.CPUCores) // Increase CPU Core processes to SessionConf value

	WiremeltAscii()
	session.Read()
	InitClient(session)
}

// Session Cleanup
func (session *SessionConfiguration) SessionCleanup() {
	fmt.Println()
}

// CLI-based setup of Client Configuration
func PromptClientConfInit() *ClientConfiguration {
	fmt.Println()

	// Strings
	validateString := func(input string) error {
		isString := utils.CheckStringForEmptiness(input)
		if !isString {
			return errors.New("invalid Input")
		}
		return nil
	}

	// Numbers
	validateFloat := func(input string) error {
		_, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return errors.New("invalid Input")
		}
		return nil
	}

	// ** DATABASE **
	// DB Host
	promptDBHost := promptui.Prompt{
		Label:    "DB Host",
		Validate: validateString,
	}
	resultDBHost, err := promptDBHost.Run()
	if err != nil {
		fmt.Printf("resultDBHost Error: %v\n", err)
	}

	// DB Port
	promptDBPort := promptui.Prompt{
		Label:    "DB Port",
		Validate: validateFloat,
	}
	resultDBPort, err := promptDBPort.Run()
	if err != nil {
		fmt.Printf("resultDBPort Error: %v\n", err)
	}

	// DB User
	promptDBUser := promptui.Prompt{
		Label:    "DB User",
		Validate: validateString,
	}
	resultDBUser, err := promptDBUser.Run()
	if err != nil {
		fmt.Printf("resultDBUser Error: %v\n", err)
	}

	// DB Password
	promptDBPass := promptui.Prompt{
		Label:    "DB Password",
		Validate: validateString,
	}
	resultDBPass, err := promptDBPass.Run()
	if err != nil {
		fmt.Printf("resultDBPass Error: %v\n", err)
	}

	// DB Name
	promptDBName := promptui.Prompt{
		Label:    "DB Name",
		Validate: validateString,
	}
	resultDBName, err := promptDBName.Run()
	if err != nil {
		fmt.Printf("resultDBName Error: %v\n", err)
	}

	// ** 2CAPTCHA API **
	// 2Captcha.com Access Key
	prompt2CapAccess := promptui.Prompt{
		Label:    "2Captcha.com Access Key",
		Validate: validateString,
	}
	result2CapAccess, err := prompt2CapAccess.Run()
	if err != nil {
		fmt.Printf("result2CapAccess Error: %v\n", err)
	}

	// ** PROXY **
	// Proxy Host
	promptProxyHost := promptui.Prompt{
		Label:    "Proxy Host",
		Validate: validateString,
	}
	resultProxyHost, err := promptProxyHost.Run()
	if err != nil {
		fmt.Printf("resultProxyHost Error: %v\n", err)
	}

	// Proxy Port
	promptProxyPort := promptui.Prompt{
		Label:    "Proxy Port",
		Validate: validateString,
	}
	resultProxyPort, err := promptProxyPort.Run()

	if err != nil {
		fmt.Printf("Prompt Failed: %v\n", err)
	}

	dbConn := helpers.DatabaseConnection{Host: resultDBHost, Port: resultDBPort, User: resultDBUser, Pass: resultDBPass, Name: resultDBName} // Initialize database connection
	reCapConn := helpers.TwoCaptchaConnection{AccessKey: result2CapAccess}                                                                   // Initialize 2captcha.com connection
	proxyConn := helpers.ProxyConnection{Host: resultProxyHost, Port: resultProxyPort}                                                       // Initialize proxy connection

	newConf := *newClientConfig(dbConn, reCapConn, proxyConn) // Initialize ClientConfiguration with input values
	newConf.UpdateClientConfiguration()

	return &newConf
}

// CLI-based setup of Session Configuration
func PromptSessionConfInit() *SessionConfiguration {
	fmt.Println()

	// Strings
	validateString := func(input string) error {
		isString := utils.CheckStringForEmptiness(input)
		if !isString {
			return errors.New("invalid Input")
		}
		return nil
	}

	// Numbers
	validateFloat := func(input string) error {
		_, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return errors.New("invalid Input")
		}
		return nil
	}

	// ** SETUP **
	// Repeat Cycle
	promptRepeatCycle := promptui.Select{
		Label: "Repeat Cycle?",
		Items: []string{"No", "Yes"},
	}
	_, resultRepeatCycle, err := promptRepeatCycle.Run()
	if err != nil {
		log.Fatalln("resultRepeatCycle Error:", err)
	}

	// CPU Cores
	promptCPUCores := promptui.Prompt{
		Label:    "CPU Cores",
		Validate: validateFloat,
	}
	resultCPUCores, err := promptCPUCores.Run()
	if err != nil {
		fmt.Printf("resultCPUCores Error: %v\n", err)
	}

	// Factory Quantity
	promptFactoryQtn := promptui.Prompt{
		Label:    "Factory Quantity",
		Validate: validateFloat,
	}
	resultFactoryQtn, err := promptFactoryQtn.Run()
	if err != nil {
		fmt.Printf("resultFactoryQtn Error: %v\n", err)
	}

	// Worker Quantity
	promptWorkerQtn := promptui.Prompt{
		Label:    "Workers Per Factory",
		Validate: validateFloat,
	}
	resultWorkerQtn, err := promptWorkerQtn.Run()
	if err != nil {
		fmt.Printf("resultWorkerQtn Error: %v\n", err)
	}

	// Jobs Per Macro
	promptJobsPerMacro := promptui.Prompt{
		Label:    fmt.Sprintf("Jobs Per Macro (x%s)", resultWorkerQtn),
		Validate: validateFloat,
	}
	resultJobsPerMacro, err := promptJobsPerMacro.Run()
	if err != nil {
		fmt.Printf("resultJobsPerMacro Error: %v\n", err)
	}

	// Factory Focus
	promptFactoryFocus := promptui.Prompt{
		Label:    fmt.Sprintf("Factory Focus (x%s)", resultFactoryQtn),
		Validate: validateString,
	}
	resultFactoryFocus, err := promptFactoryFocus.Run()
	if err != nil {
		fmt.Printf("resultFactoryFocus Error: %v\n", err)
	}

	// Worker Roles
	promptWorkerRoles := promptui.Prompt{
		Label:    fmt.Sprintf("Worker Role / Label (x%s)", resultWorkerQtn),
		Validate: validateString,
	}
	resultWorkerRoles, err := promptWorkerRoles.Run()
	if err != nil {
		fmt.Printf("resultWorkerRoles Error: %v\n", err)
	}

	// Macro Spec
	promptMacroSpec := promptui.Prompt{
		Label:    "Macro Specification",
		Validate: validateString,
	}
	resultMacroSpec, err := promptMacroSpec.Run()
	if err != nil {
		fmt.Printf("resultMacroSpec Error: %v\n", err)
	}

	// Neural Enabled
	promptNeuralEnabled := promptui.Select{
		Label: "Neural Network?",
		Items: []string{"No", "Yes"},
	}
	_, resultNeuralEnabled, err := promptNeuralEnabled.Run()
	if err != nil {
		log.Fatalln("resultNeuralEnabled Error:", err)
	}

	// Train Limit
	promptTrainLimit := promptui.Prompt{
		Label:    "Train Limit",
		Validate: validateFloat,
	}
	resultTrainLimit, err := promptTrainLimit.Run()
	if err != nil {
		fmt.Printf("resultTrainLimit Error: %v\n", err)
	}

	// Conv RepeatCycle
	repeatCycle := utils.YesNoToInt(resultRepeatCycle)

	// Conv CPUCores
	cpuCores, err := strconv.Atoi(resultCPUCores)
	if err != nil {
		log.Println("cpuCores Error:", err)
	}
	if cpuCores > runtime.NumCPU() { // Ensure cpuCores value is less than or equal to that of the hardware's
		cpuCores = runtime.NumCPU()
	}

	// Conv FactoryQtn
	factoryQtn, err := strconv.Atoi(resultFactoryQtn)
	if err != nil {
		log.Fatalln("factoryQtn Error:", err)
	}

	// Conv WorkerQtn
	workerQtn, err := strconv.Atoi(resultWorkerQtn)
	if err != nil {
		log.Fatalln("workerQtn Error:", err)
	}

	// Conv JobsPerMacro
	jobsPerMacro, err := strconv.Atoi(resultJobsPerMacro)
	if err != nil {
		log.Fatalln("jobsPerMacro Error:", err)
	}

	// Conv FactoryFocus
	focalResults := strings.ReplaceAll(resultFactoryFocus, " ", "")
	focalPoints := strings.Split(focalResults, ",")
	focalMap := make(map[int]string)
	for i := 0; i < len(focalPoints); i++ {
		focalMap[i] = focalPoints[i]
	}
	factoryFocus := focalMap

	// Conv WorkerRoles
	posResults := strings.ReplaceAll(resultWorkerRoles, " ", "")
	workerPositions := strings.Split(posResults, ",")
	roleMap := make(map[int]string)
	for i := 0; i < len(workerPositions); i++ {
		roleMap[i] = strings.TrimSpace(workerPositions[i])
	}
	workerRoles := roleMap

	// Conv MacroSpecs
	specResults := strings.ReplaceAll(resultMacroSpec, " ", "")
	macroSpecs := strings.Split(specResults, ",")
	specMap := make(map[int]string)
	for i := 0; i < len(macroSpecs); i++ {
		if strings.Contains(macroSpecs[i], "macro.") {
			ignMacroOpr := strings.TrimPrefix(macroSpecs[i], "macro")
			ignDotNot := strings.TrimPrefix(ignMacroOpr, ".")
			genMacroName := utils.CapitalizeString(ignDotNot)

			specMap[i] = strings.TrimSpace(genMacroName) // job.Macro
		} else {
			specMap[i] = strings.TrimSpace(macroSpecs[i]) // job.Macro
		}
	}

	// Assign ParamArg to Macro
	for _, macro := range specMap {
		if !strings.Contains(macro, "*") { // * denotes to ignore passed paramArg
			promptMacroParam := promptui.Prompt{
				Label:    fmt.Sprintf("Macro `%s` Param / Script", macro),
				Validate: validateString,
			}
			resultMacroParam, err := promptMacroParam.Run()
			if err != nil {
				fmt.Printf("resultMacroParam Error: %v\n", err)
			}

			worker.MacroSpecs[macro] = strings.TrimSpace(resultMacroParam) // job.paramArg
		} else {
			worker.MacroSpecs[macro] = "" // job.paramArg
		}
	}

	// Conv NeuralEnabled
	neuralEnabled := utils.YesNoToInt(resultNeuralEnabled)

	// Conv TrainLimit
	trainLimit, err := strconv.Atoi(resultTrainLimit)
	if err != nil {
		log.Fatalln("trainLimit Error:", err)
	}

	// Shell
	var shellCycle int
	if repeatCycle != 1 { // Ensure shell can't initialize with repeatCycle
		promptShellCycle := promptui.Select{
			Label: "Shell after Cycle?",
			Items: []string{"Yes", "No"},
		}
		_, resultShellCycle, err := promptShellCycle.Run()
		if err != nil {
			log.Fatalln("resultShellCycle Error:", err)
		}

		// Conv Shell
		shellCycle = utils.YesNoToInt(resultShellCycle)
	} else {
		shellCycle = 0
	}

	macroSpec := worker.MacroSpecs

	newConf := *NewSessionConfig(repeatCycle, cpuCores, factoryQtn, workerQtn, jobsPerMacro, factoryFocus, workerRoles, macroSpec, shellCycle, neuralEnabled, trainLimit) // Initialize SessionConfiguration with input values
	conf, err := json.Marshal(newConf)                                                                                                                                    // Convert SessionConfiguration to JSON object
	if err != nil {
		log.Println(err)
	}

	strConf := string(conf)                                          // Convert SessionConfiguration JSON object to string
	baseConf64 := base64.StdEncoding.EncodeToString([]byte(strConf)) // Encode `strConf` to base64
	envKeyValue := fmt.Sprintf("SESSION_CONFIG=%s\n", baseConf64)    // Convert base64 to string for .env file

	utils.WriteToEnv("SESSION_CONFIG", envKeyValue)
	fmt.Println("\n\t[✓] SessionConf saved!")

	return &newConf
}

// Return if ENV file exists
func DoesEnvFileExist() bool {
	err := godotenv.Load(".env")
	return err == nil
}

// Return if ClientConf exists in ENV file
func DoesClientConfExist() bool {
	_, present := os.LookupEnv("CLIENT_CONFIG")
	if present {
		return true
	} else {
		return false
	}
}

// Return if SessionConf exists in ENV file
func DoesSessionConfExist() bool {
	_, present := os.LookupEnv("SESSION_CONFIG")
	if present {
		return true
	} else {
		return false
	}
}

// Obtain ClientConf from ENV file for usage
func LoadClientConfiguration() *ClientConfiguration {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("error loading .env file")
	}

	clientConf := os.Getenv("CLIENT_CONFIG")                      // Get ClientConf value from .env CLIENT_CONFIG key
	rawConf64, err := base64.StdEncoding.DecodeString(clientConf) // Decode ClientConf string
	if err != nil {
		log.Fatalln(err)
	}
	//fmt.Printf("DECODED: %s\n", rawConf64)

	loadConf := &ClientConfiguration{}
	json.Unmarshal(rawConf64, &loadConf) // Convert ClientConfiguration to JSON object
	//fmt.Println("ClientConf:", loadConf)

	return loadConf
}

func (newConfig *ClientConfiguration) UpdateClientConfiguration() {
	conf, err := json.Marshal(newConfig) // Convert ClientConfiguration to JSON object
	if err != nil {
		log.Println(err)
	}

	strConf := string(conf)                                          // Convert ClientConfiguration JSON object to string
	baseConf64 := base64.StdEncoding.EncodeToString([]byte(strConf)) // Encode `strConf` to base64
	envKeyValue := fmt.Sprintf("CLIENT_CONFIG=%s\n", baseConf64)     // Convert base64 to string for .env file

	utils.WriteToEnv("CLIENT_CONFIG", envKeyValue)
	fmt.Println("\n\t[✓] ClientConf saved!")
}

// Obtain SessionConf from ENV file for usage
func LoadSessionConfiguration() *SessionConfiguration {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln("error loading .env file")
	}

	sessionConf := os.Getenv("SESSION_CONFIG")                     // Get SessionConf value from .env SESSION_CONFIG key
	rawConf64, err := base64.StdEncoding.DecodeString(sessionConf) // Decode SessionConf string
	if err != nil {
		log.Fatalln(err)
	}
	//fmt.Printf("DECODED: %s\n", rawConf64)

	loadConf := &SessionConfiguration{}
	json.Unmarshal(rawConf64, &loadConf) // Convert SessionConfiguration to JSON object
	//fmt.Println("SessionConf:", loadConf)

	worker.MacroSpecs = loadConf.MacroLibrary

	return loadConf
}

func (newConfig *SessionConfiguration) UpdateSessionConfiguration() {
	conf, err := json.Marshal(newConfig) // Convert SessionConfiguration to JSON object
	if err != nil {
		log.Println(err)
	}

	fmt.Println("\n+ MACRO LIBRARY:", newConfig.MacroLibrary)

	strConf := string(conf)                                          // Convert SessionConfiguration JSON object to string
	baseConf64 := base64.StdEncoding.EncodeToString([]byte(strConf)) // Encode `strConf` to base64
	envKeyValue := fmt.Sprintf("SESSION_CONFIG=%s\n", baseConf64)    // Convert base64 to string for .env file

	utils.WriteToEnv("SESSION_CONFIG", envKeyValue)
	fmt.Println("\n\t[✓] SessionConf updated!")
}
