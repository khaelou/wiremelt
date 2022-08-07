package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"wiremelt/macro"
	"wiremelt/pilot"
	//"wiremelt/pilot"
	"wiremelt/shell"
	"wiremelt/utils"
	"wiremelt/wiremelt"
	"wiremelt/worker"

	"github.com/manifoldco/promptui"
	"github.com/urfave/cli/v2"
	"golang.org/x/exp/maps"
)

func main() {
	var client *wiremelt.ClientConfiguration
	var session *wiremelt.SessionConfiguration

	var file string

	// CLI Initialization
	app := &cli.App{
		Name:        "Wiremelt",
		Usage:       "Extendible Automation Utility",
		Description: "Extendible automation utility for parallel concurrent worker-pool operations at scale.",
		Version:     "1.0.0",

		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "client",
				Aliases: []string{"c"},
				Usage:   "Setup new Client Configuration",
			},
			&cli.StringFlag{
				Name:    "session",
				Aliases: []string{"s"},
				Usage:   "Setup new Session Configuration",
			},
			&cli.StringFlag{
				Name:    "macro",
				Aliases: []string{"m"},
				Usage:   "Macro Library / Import Custom Macro",
			},
			&cli.StringFlag{
				Name:    "shell",
				Aliases: []string{"sh"},
				Usage:   "Launch Shell",
			},
			&cli.StringFlag{
				Name:    "web",
				Aliases: []string{"w"},
				Usage:   "Launch Web UI",
			},
			/*
				&cli.StringFlag{
					Name:    "pilot",
					Aliases: []string{"p"},
					Usage:   "Launch Pilot for Rod",
				},
			*/
			&cli.StringFlag{
				Name:        "file",
				Aliases:     []string{"f"},
				Destination: &file,
				Usage:       "Load File",
			},
		},
		Action: func(context *cli.Context) error {
			// Strings
			validateString := func(input string) error {
				isString := utils.CheckStringForEmptiness(input)
				if !isString {
					return errors.New("invalid Input")
				}
				return nil
			}

			fmt.Println("\n\tWIREMELT")

			if context.Args().Len() > 0 {
				fmt.Println()

				flag := context.Args().Get(0)

				switch flag {
				case "client":
					promptNewClientConf := promptui.Select{
						Label: "New Client Configuration?",
						Items: []string{"Yes", "No"},
					}

					_, reply, err := promptNewClientConf.Run()
					if err != nil {
						log.Fatalf("promptNewClientConf error: %v\n", err)
					}

					switch reply {
					case "Yes":
						wiremelt.PromptClientConfInit()
					case "No":
						loadClientConf := wiremelt.LoadClientConfiguration()
						loadClientConf.Read()
					}
				case "session":
					promptNewSessionConf := promptui.Select{
						Label: "New Session Configuration?",
						Items: []string{"Yes", "No"},
					}

					_, reply, err := promptNewSessionConf.Run()
					if err != nil {
						log.Fatalf("promptNewSessionConf error: %v\n", err)
					}

					switch reply {
					case "Yes":
						wiremelt.PromptSessionConfInit()
					case "No":
						loadSessConf := wiremelt.LoadSessionConfiguration()
						loadSessConf.Read()
					}
				case "macro":
					importName := context.Args().Get(1)
					importURL := context.Args().Get(2)

					if !utils.CheckStringForEmptiness(importName) {
						loadSess := wiremelt.LoadSessionConfiguration()
						fmt.Println("~ MACRO LIBRARY:", loadSess.MacroLibrary)
					} else {
						ignMacroOpr := strings.TrimPrefix(importName, "macro")
						ignDotNot := strings.TrimPrefix(ignMacroOpr, ".")
						genMacroName := utils.CapitalizeString(ignDotNot)
						specID := fmt.Sprintf("macro.%s", genMacroName)

						// Ensure URL is a valid route and contains JS document
						if utils.IsStringValidUrl(importURL) && strings.Contains(importURL, ".js") {
							macroDest := fmt.Sprintf("%s.js", utils.LowercaseString(genMacroName))

							dlSuccess, customScript, err := utils.DownloadTarget(importURL, macroDest, true) // Download external JS script
							if err != nil {
								fmt.Println("\t[x]", specID, "IMPORT", genMacroName, "@", importURL)
							}

							if dlSuccess {
								fmt.Println("\t[✓]", specID, "IMPORT", genMacroName, "@", customScript)

								fileDir := fmt.Sprintf("%s%s", worker.MacroFileDir, customScript)
								if _, err := os.Stat(fileDir); !os.IsNotExist(err) {
									fmt.Println()

									// Parse file
									fileContent, err := ioutil.ReadFile(fileDir)
									if err != nil {
										log.Fatalln(err)
									}

									convScript := string(fileContent)                // Convert []byte to string
									execJS, err := utils.V8NodeJS(convScript, false) // Test script
									if err != nil {
										log.Println(err, execJS)
									}

									// Add to MacroSpecs
									ignParamName := genMacroName
									worker.MacroSpecs[ignParamName] = customScript // "MacroName": "macro.js"

									if wiremelt.DoesEnvFileExist() {
										existingSession := wiremelt.LoadSessionConfiguration()
										maps.Copy(existingSession.MacroLibrary, worker.MacroSpecs)                                                                                                                                                                                                                                                        // Copy local macroSpecs into saved sessionConf
										newConf := wiremelt.NewSessionConfig(existingSession.RepeatCycle, existingSession.CPUCores, existingSession.FactoryQuantity, existingSession.WorkerQuantity, existingSession.JobsPerFactory, existingSession.FactoryFocus, existingSession.WorkerRoles, existingSession.MacroLibrary, existingSession.ShellCycle) // Initialize SessionConfiguration with input values
										newConf.UpdateSessionConfiguration()
									} else {
										fmt.Println("\n+ TEMP MACRO LIBRARY:", worker.MacroSpecs)
									}
								}
							}
						} else {
							existingSession := wiremelt.LoadSessionConfiguration()
							existingMacroLibrary := existingSession.MacroLibrary

							addNewMacro := func() {
								if !utils.CheckStringForEmptiness(importURL) {
									importName = fmt.Sprintf("%s*", importName)
									importURL = ""
								}

								fmt.Println("\t> ADD DEFAULT MACRO:", importName)

								worker.MacroSpecs[importName] = importURL

								if wiremelt.DoesEnvFileExist() {
									maps.Copy(existingMacroLibrary, worker.MacroSpecs)                                                                                                                                                                                                                                                                // Copy local macroSpecs into saved sessionConf
									newConf := wiremelt.NewSessionConfig(existingSession.RepeatCycle, existingSession.CPUCores, existingSession.FactoryQuantity, existingSession.WorkerQuantity, existingSession.JobsPerFactory, existingSession.FactoryFocus, existingSession.WorkerRoles, existingSession.MacroLibrary, existingSession.ShellCycle) // Initialize SessionConfiguration with input values
									newConf.UpdateSessionConfiguration()
								} else {
									fmt.Println("\n+ TEMP MACRO LIBRARY:", worker.MacroSpecs)
								}
							}

							if _, ok := macro.MacroLibrary[importName]; ok { // Default Macro
								addNewMacro()
							} else if strings.Contains(importURL, ".js") { // Custom Macro
								addNewMacro()
							} else {
								log.Fatalln("[x] macro import must reference a default macro or (.js) JavaScript file.")
							}
						}
					}
				case "del":
					if wiremelt.DoesEnvFileExist() {
						sessConf := wiremelt.LoadSessionConfiguration()
						macroSpec := sessConf.MacroLibrary

						fmt.Println("\n~ SESSION MACROS:", macroSpec)
						fmt.Println()

						promptTargetMacro := promptui.Prompt{
							Label:    "Delete Macro",
							Validate: validateString,
						}
						resultTargetMacro, err := promptTargetMacro.Run()
						if err != nil {
							fmt.Printf("resultTargetMacro Error: %v\n", err)
						}

						parseTarget := strings.TrimSpace(utils.CapitalizeString(resultTargetMacro))

						if _, ok := macroSpec[parseTarget]; ok {
							delete(macroSpec, parseTarget)

							newConf := *wiremelt.NewSessionConfig(sessConf.RepeatCycle, sessConf.CPUCores, sessConf.FactoryQuantity, sessConf.WorkerQuantity, sessConf.JobsPerFactory, sessConf.FactoryFocus, sessConf.WorkerRoles, macroSpec, sessConf.ShellCycle) // Initialize SessionConfiguration with input values
							newConf.UpdateSessionConfiguration()

							fmt.Println("\n- MACRO LIBRARY:", macroSpec)
						}
					} else {
						log.Fatalln("[x] macro del requires a session configuration with specified macros.")
					}
				case "shell":
					macroSpec := wiremelt.LoadSessionConfiguration().MacroLibrary
					shell.InitShell(macroSpec)
				case "web":
					// WebAssembly
				case "pilot":
					pilot.InitPilot()
				case "file":
					fmt.Println("LOAD FILE")
				default:
					fmt.Printf("Flag: `%v`\n", flag) // .Get(i) obtains element by index from cli.Context.Args()
				}

				fmt.Println()
			} else {
				// Check for an existing ClientConfig in .env file
				if wiremelt.DoesEnvFileExist() {
					client = wiremelt.LoadClientConfiguration()
				} else {
					client = wiremelt.PromptClientConfInit()
				}

				if client.Parse() {
					client.Read()

					// Check for an existing SessionConfig in .env file
					if wiremelt.DoesSessionConfExist() {
						session = wiremelt.LoadSessionConfiguration()
					} else {
						session = wiremelt.PromptSessionConfInit()
					}

					session.StartSession(client)
				}
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
