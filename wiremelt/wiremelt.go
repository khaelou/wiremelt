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

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/manifoldco/promptui"
)

func WiremeltAscii() (interface{}, error) {
	color.HiYellow(`
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

// Initiate New Client configuration
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
	color.HiMagenta("\n\tClient-Configuration: [DBConn] %v | [2Captcha] %v | [Proxy] %v", config.DBConn, config.ReCaptcha, config.Proxy)
}

// Parse takes ClientConfiguration as input, validates upon initiating actual client
func (config *ClientConfiguration) Parse() bool {
	if config == nil {
		return false
	} else {
		return true
	}
}

// SessionConfiguration references setup details prior to workload execution, can save or load from JSON
type SessionConfiguration struct {
	SessionName     string           `json:"sessionName"`
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

// Initiate New Session configuration
func NewSessionConfig(sessionName string, repeatCycle, cpuCores, factoryQuantity, workerQuantity int, jobsPerMacro int, factoryFocus map[int]string, workerRoles map[int]string, macroSpec worker.MacroSpec, shellCycle int, neuralEnabled int, trainLimit int) *SessionConfiguration {
	config := SessionConfiguration{
		SessionName:     sessionName,
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
	color.HiBlue("\n\tSession-Configuration: %v", *config)
}

// Parse takes SessionConfiguration as input, validates upon initiating into a Session
func (config *SessionConfiguration) Parse() bool {
	if config == nil {
		return false
	} else {
		return true
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

// Strings (Validation)
var validateString = func(input string) error {
	isString := utils.CheckStringForEmptiness(input)
	if !isString {
		return errors.New("invalid Input")
	}
	return nil
}

// Numbers (Validation)
var validateFloat = func(input string) error {
	_, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return errors.New("invalid Input")
	}
	return nil
}

// Shell prompts for SessionConfiguration
var sessionPrompts = func() (int, int, int, int, int, map[int]string, map[int]string, worker.MacroSpec, int, int, int) {
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
		Label:    fmt.Sprintf("Worker Role (x%s)", resultWorkerQtn),
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

	specResults := strings.ReplaceAll(resultMacroSpec, " ", "")
	macroSpecs := strings.Split(specResults, ",")
	macroCount := len(macroSpecs)

	// Jobs Per Macro
	promptJobsPerMacro := promptui.Prompt{
		Label:    fmt.Sprintf("Jobs Per Macro (x%v)", macroCount),
		Validate: validateFloat,
	}
	resultJobsPerMacro, err := promptJobsPerMacro.Run()
	if err != nil {
		fmt.Printf("resultJobsPerMacro Error: %v\n", err)
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
	specResults = strings.ReplaceAll(resultMacroSpec, " ", "")
	macroSpecs = strings.Split(specResults, ",")
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

			// CHECK: If imported macro is declared w/o script; auto-locate script using camelCaseNaming.js

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
			Items: []string{"No", "Yes"},
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

	return repeatCycle, cpuCores, factoryQtn, workerQtn, jobsPerMacro, factoryFocus, workerRoles, macroSpec, shellCycle, neuralEnabled, trainLimit
}

// CLI-based setup of Client Configuration
func PromptClientConfInit() *ClientConfiguration {
	fmt.Println()

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

	// ** SETUP **
	// Session Name
	promptSessionName := promptui.Prompt{
		Label:    "Session Name",
		Validate: validateString,
	}
	resultSessionName, err := promptSessionName.Run()
	if err != nil {
		fmt.Printf("resultSessionName Error: %v\n", err)
	}

	// Ensure pre-existing sessions do not share target name
	_, _, savedSessionFound, savedSessionConf := DoesSessionsExist(resultSessionName, true)
	if savedSessionFound {
		color.HiCyan("\t- SESSION_LOCATED: The previous configuration will be overwritten!\n")
		savedSessionConf.AlterSessionConfiguration()

		return &savedSessionConf
	} else {
		repeatCycle, cpuCores, factoryQtn, workerQtn, jobsPerMacro, factoryFocus, workerRoles, macroSpec, shellCycle, neuralEnabled, trainLimit := sessionPrompts()
		newConf := *NewSessionConfig(resultSessionName, repeatCycle, cpuCores, factoryQtn, workerQtn, jobsPerMacro, factoryFocus, workerRoles, macroSpec, shellCycle, neuralEnabled, trainLimit) // Initialize SessionConfiguration with input values
		conf, err := json.Marshal(newConf)                                                                                                                                                       // Convert SessionConfiguration to JSON object
		if err != nil {
			log.Println(err)
		}

		strConf := string(conf)                                           // Convert SessionConfiguration JSON object to string
		baseConf64 := base64.StdEncoding.EncodeToString([]byte(strConf))  // Encode `strConf` to base64
		envKeyValue := fmt.Sprintf("SESSION_CONFIG=%s\n", baseConf64)     // Convert base64 to string
		envKeyValueNoLine := fmt.Sprintf("SESSION_CONFIG=%s", baseConf64) // Convert base64 to string

		sessions, sessionsFound, sessionFound, _ := DoesSessionsExist(resultSessionName, true)
		if sessions != nil && sessionsFound && sessionFound {
			utils.WriteToEnv("SESSION_CONFIG", envKeyValueNoLine)
			color.Blue("\n\t[✓] Session '%s' updated!", sess.SessionName)
		} else {
			if sessions == nil && !sessionsFound {
				sessionConfigs := make([]string, 0)                 // Create slice containing existing base64 JSON results
				sessionConfigs = append(sessionConfigs, baseConf64) // Add newly created sessionConf to the slice
				//fmt.Printf("~ INIT_SESSIONS_B4_ENCODE: %v\n", sessionConfigs)

				// +-
				sessConf, err := json.Marshal(sessionConfigs) // Convert SessionConfigs to JSON object
				if err != nil {
					log.Println("sessConf", err)
				}

				// Encode slice to string
				baseSessConf64 := base64.StdEncoding.EncodeToString([]byte(sessConf)) // Encode `sessConf` to base64
				sessEnvKeyValue := fmt.Sprintf("SESSIONS=%s", baseSessConf64)         // Convert base64 to string
				//fmt.Printf("~ INIT_SESSIONS_ENCODED_SLICE: %s\n", baseSessConf64)

				/*
					// Decode slice from string
					rawSavedSessConf64, err := base64.StdEncoding.DecodeString(baseSessConf64) // Decode SessionConf string
					if err != nil {
						log.Fatalln("101", err)
					}
					sessEnvDecode := make([]string, 0)
					json.Unmarshal(rawSavedSessConf64, &sessEnvDecode)
					fmt.Printf("~ INIT_SESSIONS_DECODED_SLICE: %v\n", sessEnvDecode)
				*/
				// --

				utils.WriteToEnv("SESSION_CONFIG", envKeyValue)
				utils.WriteToEnv("SESSIONS", sessEnvKeyValue)
				color.Magenta("\n\t[✓] Added '%s' to Saved Sessions!", resultSessionName)
			} else {
				if !sessionFound {
					if DoesEnvFileExist() {
						_, present := os.LookupEnv("SESSIONS")
						if present {
							sessions := os.Getenv("SESSIONS")                           // Get Sessions value from .env SESSIONS key
							rawSess64, err := base64.StdEncoding.DecodeString(sessions) // Decode Sessions string
							if err != nil {
								log.Fatalln("exist", err)
							}

							savedSessions := make([]string, 0)
							json.Unmarshal(rawSess64, &savedSessions) // Convert SessionConfiguration to JSON object
							//fmt.Printf("\n~ INIT_PREV_SESSIONS_DECODED: %v\n", savedSessions)

							updateSessions := make([]string, 0)
							updateSessions = append(updateSessions, savedSessions...) // Get existing sessions
							updateSessions = append(updateSessions, baseConf64)       // Add newly created session
							//fmt.Printf("~ INIT_UPDATE_SESSIONS_B4_ENCODE: %v\n", updateSessions)

							// +-
							sessConf, err := json.Marshal(updateSessions) // Convert updateSessions to JSON object
							if err != nil {
								log.Println("sessConf", err)
							}

							baseSessConf64 := base64.StdEncoding.EncodeToString([]byte(sessConf)) // Encode `sessConf` to base64
							sessEnvKeyValue := fmt.Sprintf("SESSIONS=%s", baseSessConf64)         // Convert base64 to string
							//fmt.Printf("~ INIT_UPDATE_SESSIONS_ENCODED_SLICE: %s\n", baseSessConf64)
							// --

							utils.WriteToEnv("SESSION_CONFIG", envKeyValueNoLine)
							utils.WriteToEnv("SESSIONS", sessEnvKeyValue)
							color.Magenta("\n\t[✓✓] Added '%s' to Saved Sessions!", resultSessionName)
						}
					} else {
						fmt.Println(color.HiRedString("\n\t~ INIT_UPDATE_SESSIONS :. Env file misplaced.\n"))
					}
				}
			}
		}

		return &newConf
	}
}

// Return if ENV file exists
func DoesEnvFileExist() bool {
	err := godotenv.Load(".env")
	return err == nil
}

// Return if ClientConf exists in ENV file
func DoesClientConfExist() bool {
	if DoesEnvFileExist() {
		_, present := os.LookupEnv("CLIENT_CONFIG")
		if present {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

// Return if SessionConf exists in ENV file
func DoesSessionConfExist() bool {
	if DoesEnvFileExist() {
		_, present := os.LookupEnv("SESSION_CONFIG")
		if present {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

// Return if Sessions exists in ENV file
func DoesSessionsExist(lookup string, stateLocated bool) ([]string, bool, bool, SessionConfiguration) { // savedSessions[], sessionsExist?, sessionFound?
	if DoesEnvFileExist() {
		_, present := os.LookupEnv("SESSIONS")
		if present {
			sessions := os.Getenv("SESSIONS")                           // Get Sessions value from .env SESSIONS key
			rawSess64, err := base64.StdEncoding.DecodeString(sessions) // Decode Sessions string
			if err != nil {
				log.Fatalln("exist", err)
			}

			savedSessions := make([]string, 0)
			json.Unmarshal(rawSess64, &savedSessions) // Convert SessionConfiguration to JSON object
			//fmt.Printf("\n~ EXISTING_SESSIONS_DECODED: %v\n", savedSessions)

			var existingConf SessionConfiguration

			sessionConfigs := savedSessions
			for _, session := range sessionConfigs {
				//fmt.Println("SAVED_SESSION:", session)

				rawSessionConf64, err := base64.StdEncoding.DecodeString(session) // Decode SessionConf string
				if err != nil {
					log.Fatalln("rawSessionConf64", err)
				}
				loadConf := &SessionConfiguration{}
				json.Unmarshal(rawSessionConf64, &loadConf) // Convert SessionConfiguration to JSON object
				//fmt.Printf("\n- EXISTING_SESSION_DECODED: %s | %v\n", rawSessionConf64, loadConf)

				if strings.Contains(strings.ToLower(lookup), strings.ToLower(loadConf.SessionName)) {
					if stateLocated {
						fmt.Println(color.HiMagentaString(fmt.Sprintf("\n\t~ SAVED_SESSION_CONF :. '%s' located! \n\t%v", loadConf.SessionName, loadConf)))
					}
					existingConf = *loadConf

					return sessionConfigs, true, true, existingConf
				} else {
					if lookup == "" {
						fmt.Println(color.HiBlueString("\n\t~ EXISTING_SESSION_CONF :. '%s' discovered. \n\t%v", loadConf.SessionName, loadConf))
					}
				}
			}

			return sessionConfigs, true, false, existingConf
		} else {
			return nil, false, false, SessionConfiguration{}
		}
	} else {
		return nil, false, false, SessionConfiguration{}
	}
}

// Return JSON of sessions stored
func GetSessions() []interface{} {
	if DoesEnvFileExist() {
		_, present := os.LookupEnv("SESSIONS")
		if present {
			sessions := os.Getenv("SESSIONS")                           // Get Sessions value from .env SESSIONS key
			rawSess64, err := base64.StdEncoding.DecodeString(sessions) // Decode Sessions string
			if err != nil {
				log.Fatalln("exist", err)
			}

			savedSessions := make([]string, 0)
			json.Unmarshal(rawSess64, &savedSessions) // Convert SavedSessions to JSON object

			decodedSessions := make([]interface{}, 0)
			for _, session := range savedSessions {
				rawSessionConf64, err := base64.StdEncoding.DecodeString(session) // Decode SessionConf string
				if err != nil {
					log.Fatalln("rawSessionConf64", err)
				}
				loadConf := &SessionConfiguration{}
				json.Unmarshal(rawSessionConf64, &loadConf) // Convert SessionConfiguration to JSON object

				decodedSessions = append(decodedSessions, loadConf)
			}

			return decodedSessions
		} else {
			fmt.Println(color.HiRedString("\n\t~ GET_SESSIONS : Env file misplaced.\n"))
		}
	}

	return nil
}

// Commit `SESSIONS` via env file with updated configuration value
func UpdateSessions(targetConf SessionConfiguration, stateLocated bool) {
	if DoesEnvFileExist() {
		_, present := os.LookupEnv("SESSIONS")
		if present {
			sessions := os.Getenv("SESSIONS")                           // Get Sessions value from .env SESSIONS key
			rawSess64, err := base64.StdEncoding.DecodeString(sessions) // Decode Sessions string
			if err != nil {
				log.Fatalln("exist", err)
			}

			savedSessions := make([]string, 0)
			json.Unmarshal(rawSess64, &savedSessions) // Convert SessionConfiguration to JSON object
			//fmt.Printf("~ PREV_UPDATE_SESSIONS_DECODED: %v\n", savedSessions)

			sessionConfigs := savedSessions
			for i, session := range sessionConfigs {
				//fmt.Println("SAVED_SESSION:", session)

				rawSessionConf64, err := base64.StdEncoding.DecodeString(session) // Decode SessionConf string
				if err != nil {
					log.Fatalln("rawSessionConf64", err)
				}
				loadConf := &SessionConfiguration{}
				json.Unmarshal(rawSessionConf64, &loadConf) // Convert SessionConfiguration to JSON object
				//fmt.Printf("- UPDATE_SESSION_DECODED: %s | %v\n", rawSessionConf64, loadConf)

				updateTargetConf, err := json.Marshal(targetConf) // Convert SessionConfiguration to JSON object
				if err != nil {
					log.Println(err)
				}
				strTarConf := string(updateTargetConf)                                 // Convert SessionConfiguration JSON object to string
				baseTarConf64 := base64.StdEncoding.EncodeToString([]byte(strTarConf)) // Encode `strConf` to base64

				if strings.Contains(strings.ToLower(targetConf.SessionName), strings.ToLower(loadConf.SessionName)) {
					if stateLocated {
						fmt.Println(color.HiMagentaString(fmt.Sprintf("\n\t~ TAR_SAVED_SESSION_CONF :. '%s' located!", loadConf.SessionName)))
					}
					savedSessions[i] = baseTarConf64
				} else {
					if loadConf != nil {
						if stateLocated {
							fmt.Println(color.HiBlueString("\n\t~ SAVED_SESSION_CONF :. '%s' exists.", loadConf.SessionName))
						}
					} else {
						sessionConfigs = append(sessionConfigs, baseTarConf64) // Add saved sessionConf to the updated slice
						fmt.Println(color.HiMagentaString("\n\t~ SAVED_SESSION_CONF :. Added '%s' to Saved Sessions!", loadConf.SessionName))
					}
				}
			}

			sessConf, err := json.Marshal(sessionConfigs) // Convert SessionConfigs to JSON object
			if err != nil {
				log.Println("sessConf", err)
			}

			baseSessConf64 := base64.StdEncoding.EncodeToString([]byte(sessConf)) // Encode `sessConf` to base64
			sessEnvKeyValue := fmt.Sprintf("SESSIONS=%s", baseSessConf64)         // Convert base64 to string
			//fmt.Printf("~ UPDATE_SESSIONS_ENCODED_SLICE: %s\n", baseSessConf64)

			utils.WriteToEnv("SESSIONS", sessEnvKeyValue)
		} else {
			fmt.Println(color.HiRedString("\n\t~ UPDATE_SESSIONS : Env file misplaced.\n"))
		}
	}
}

// Obtain ClientConf from ENV file for usage
func LoadClientConfiguration() *ClientConfiguration {
	if DoesEnvFileExist() {
		clientConf := os.Getenv("CLIENT_CONFIG")                      // Get ClientConf value from .env CLIENT_CONFIG key
		rawConf64, err := base64.StdEncoding.DecodeString(clientConf) // Decode ClientConf string
		if err != nil {
			log.Fatalln(err)
		}
		//fmt.Printf("DECODED: %s\n", rawConf64)

		loadConf := &ClientConfiguration{}
		json.Unmarshal(rawConf64, &loadConf) // Convert ClientConfiguration to JSON object
		//fmt.Println("ClientConf:", loadConf)

		DoesSessionsExist("", false)

		return loadConf
	} else {
		return nil
	}
}

// Commit `CLIENT_CONFIG` via env file with updated configuration changes
func (newConfig *ClientConfiguration) UpdateClientConfiguration() {
	conf, err := json.Marshal(newConfig) // Convert ClientConfiguration to JSON object
	if err != nil {
		log.Println(err)
	}

	strConf := string(conf)                                          // Convert ClientConfiguration JSON object to string
	baseConf64 := base64.StdEncoding.EncodeToString([]byte(strConf)) // Encode `strConf` to base64
	envKeyValue := fmt.Sprintf("CLIENT_CONFIG=%s\n", baseConf64)     // Convert base64 to string for .env file

	utils.WriteToEnv("CLIENT_CONFIG", envKeyValue)
	color.HiGreen("\n\t[✓] ClientConf saved!")
}

// Obtain SessionConf from ENV file for usage
func LoadSessionConfiguration() *SessionConfiguration {
	if DoesEnvFileExist() {
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
	} else {
		return nil
	}
}

// Alter exisiting session with a fresh configuration
func (existingConf *SessionConfiguration) AlterSessionConfiguration() {
	fmt.Println()

	repeatCycle, cpuCores, factoryQtn, workerQtn, jobsPerMacro, factoryFocus, workerRoles, macroSpec, shellCycle, neuralEnabled, trainLimit := sessionPrompts()
	alteredConf := *NewSessionConfig(existingConf.SessionName, repeatCycle, cpuCores, factoryQtn, workerQtn, jobsPerMacro, factoryFocus, workerRoles, macroSpec, shellCycle, neuralEnabled, trainLimit)
	alteredConf.UpdateSessionConfiguration()

	UpdateSessions(alteredConf, true)
}

// Commit `SESSION_CONFIG` via env file with updated configuration changes
func (newConfig *SessionConfiguration) UpdateSessionConfiguration() {
	conf, err := json.Marshal(newConfig) // Convert SessionConfiguration to JSON object
	if err != nil {
		log.Println(err)
	}

	strConf := string(conf)                                          // Convert SessionConfiguration JSON object to string
	baseConf64 := base64.StdEncoding.EncodeToString([]byte(strConf)) // Encode `strConf` to base64
	envKeyValue := fmt.Sprintf("SESSION_CONFIG=%s", baseConf64)      // Convert base64 to string for .env file

	sessions, sessionsFound, sessionFound, _ := DoesSessionsExist(newConfig.SessionName, true)
	if sessions != nil && sessionsFound && sessionFound {
		utils.WriteToEnv("SESSION_CONFIG", envKeyValue)
		color.Green("\n\t[✓] Session '%s' is now active!", newConfig.SessionName)

		UpdateSessions(*newConfig, false)
	}
}
