package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"wiremelt/pilot"
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
		Description: "Extendible Automation Utility; powers concurrent yet parallel worker-pool operations at scale.",
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
			&cli.StringFlag{
				Name:    "pilot",
				Aliases: []string{"p"},
				Usage:   "Launch Pilot for Rod",
			},
			&cli.StringFlag{
				Name:        "file",
				Aliases:     []string{"f"},
				Destination: &file,
				Usage:       "Load File",
			},
		},
		Action: func(context *cli.Context) error {
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

					if !utils.CheckStringForEmptiness(importName) || !utils.CheckStringForEmptiness(importURL) {
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

							dlSuccess, customScript, err := utils.DownloadTarget(importURL, macroDest) // Download
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
										conf, err := json.Marshal(newConf)                                                                                                                                                                                                                                                                                // Convert SessionConfiguration to JSON object
										if err != nil {
											log.Println(err)
										}

										fmt.Println("\n+ MACRO LIBRARY:", newConf.MacroLibrary)

										strConf := string(conf)                                          // Convert SessionConfiguration JSON object to string
										baseConf64 := base64.StdEncoding.EncodeToString([]byte(strConf)) // Encode `strConf` to base64
										envKeyValue := fmt.Sprintf("SESSION_CONFIG=%s\n", baseConf64)    // Convert base64 to string for .env file

										utils.WriteToEnv("SESSION_CONFIG", envKeyValue)
										fmt.Println("\n\t[✓] SessionConf saved!")
									} else {
										fmt.Println("\n+ TEMP MACRO LIBRARY:", worker.MacroSpecs)
									}
								}
							}
						} else {
							log.Fatalln("[x] macro import must be a valid JavaScript URL route")
						}
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
