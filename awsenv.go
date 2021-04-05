package main

import (
	"flag"
	"fmt"
	"github.com/BernhardLenz/ini"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"time"
)

var usr *user.User

type Profile struct {
	aws_access_key_id     string
	aws_secret_access_key string
	output                string
	region                string
	isActive              bool
}

var defaultProfile Profile
var profiles = make(map[string]Profile)

type Config struct {
	output string
	region string
}

var defaultConfig Config
var configs = make(map[string]Config)

var credentialsFile *ini.File
var configFile *ini.File

func main() {

	var err error
	usr, err = user.Current()
	if err != nil {
		log.Fatal(err)
	}

	listCommand := flag.NewFlagSet("list", flag.ExitOnError)
	activateCommand := flag.NewFlagSet("activate", flag.ExitOnError)

	if len(os.Args) > 3 {
		fmt.Println("ERROR: Too many arguments supplied.")
		printUsage()
		os.Exit(1)
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "list":
			listCommand.Parse(os.Args[2:])
		case "activate":
			activateCommand.Parse(os.Args[2:])
		default:
			fmt.Println("ERROR: Unknown command!")
			printUsage()
			os.Exit(1)
		}
	}

	if listCommand.Parsed() || len(os.Args) == 1 {
		parse()
		listProfiles(profiles)

		fmt.Printf("\nTo activate a different profile run '%s activate <profile>'", filepath.Base(os.Args[0]))
	} else if activateCommand.Parsed() {
		if len(os.Args) != 3 {
			fmt.Println("ERROR: Required parameter <profile> missing for activate command!")
			printUsage()
			os.Exit(1)
		}
		activateProfileName := os.Args[2]
		if activateProfileName == "default" {
			fmt.Println("ERROR: Cannot activate the 'default' profile as it is already active!")
			os.Exit(1)
		}
		parse()

		profile, ok := profiles[activateProfileName]
		if !ok {
			fmt.Printf("ERROR: Profile '%s' does not exist! Available profiles are: \n\n", activateProfileName)
			listProfiles(profiles)
			os.Exit(1)
		}

		if profile.isActive {
			fmt.Printf("Profile '%s' is already active! No changes applied. \n\n", activateProfileName)
			listProfiles(profiles)
			os.Exit(0)
		}

		setDefaultProfile(activateProfileName)

		listProfiles(profiles)
	}
} //main

func printUsage() {
	fmt.Println("Usage:")
	fmt.Printf("%s [list]\n", filepath.Base(os.Args[0]))
	fmt.Println(" Lists all available profiles.")
	fmt.Printf("%s activate <profile>\n", filepath.Base(os.Args[0]))
	fmt.Println(" Activates a given profile.")
} //printUsage

func parse() {

	ini.DefaultSection = "default"

	var err error
	credentialsFile, err = ini.Load(usr.HomeDir + "/.aws/credentials")
	if err != nil {
		fmt.Printf("Failed to read file: %v", err)
		os.Exit(1)
	}

	defaultCredentialsSection := credentialsFile.Section(ini.DefaultSection)

	for _, credentialsSection := range credentialsFile.Sections() {
		//fmt.Printf("Section sectionName: %s\n", sectionName)
		var profile Profile
		sectionName := credentialsSection.Name()
		for _, key := range credentialsSection.Keys() {
			keyName := key.Name()
			value := key.Value()
			//fmt.Printf("%s => %s\n", key, value)
			if "aws_access_key_id" == keyName {
				profile.aws_access_key_id = value
			} else if "aws_secret_access_key" == keyName {
				profile.aws_secret_access_key = value
			}
			//fmt.Printf("Default Profile: %s\n", defaultProfile.aws_access_key_id)
			//fmt.Printf("Default Profile: %s\n", defaultProfile.aws_secret_access_key)
		}
		if "default" != sectionName {
			if profile.aws_access_key_id == defaultCredentialsSection.Key("aws_access_key_id").Value() {
				profile.isActive = true
				defaultProfile = profile
				//fmt.Printf("In default\n")
				//fmt.Printf("profile: %s\n", profile)
				//fmt.Printf("defaultProfile: %s\n", defaultProfile)
			}
			profiles[sectionName] = profile
		} else {
			defaultProfile = profile
		}
	}
	//Only default section in ini file
	if defaultProfile.aws_access_key_id != "" {
		if len(profiles) == 0 {
			defaultProfile.isActive = true
			profiles["default"] = defaultProfile
		} else {
			foundActive := false
			for _, profile := range profiles {
				if profile.isActive == true {
					foundActive = true
					break
				}
			}
			if !foundActive {
				defaultProfile.isActive = true
				profiles["default"] = defaultProfile
			}
		}
	}

	configFile, err := ini.Load(usr.HomeDir + "/.aws/config")
	if err != nil {
		fmt.Printf("Failed to read file: %v", err)
		os.Exit(1)
	}

	//defaultConfigSection := configFile.Section(ini.DefaultSection)

	for _, configSection := range configFile.Sections() {
		//fmt.Printf("Section sectionName: %s\n", sectionName)
		var config Config
		sectionName := configSection.Name()
		for _, key := range configSection.Keys() {
			keyName := key.Name()
			value := key.Value()
			if "DEFAULT" == sectionName {
				if "output" == keyName {
					defaultConfig.output = value
					defaultProfile.output = value
				} else if "region" == keyName {
					defaultConfig.region = value
					defaultProfile.region = value
				}
			} else {
				profile, ok := profiles[sectionName]
				//fmt.Printf("sectionName: %s\n", sectionName)
				//fmt.Printf("profile: %s\n", profile)
				//fmt.Printf("ok: %s\n", ok)
				if ok {
					if "output" == keyName {
						profile.output = value
					} else if "region" == keyName {
						profile.region = value
					}
					//fmt.Printf("profile: %s\n", profile)
					profiles[sectionName] = profile
				}
			}
		}
		configs[sectionName] = config
	}

	foundProfileForDefault := false
	for sectionName, profile := range profiles {
		if profile.aws_access_key_id == defaultProfile.aws_access_key_id {
			foundProfileForDefault = true
			profile.isActive = true
			profiles[sectionName] = profile
			//fmt.Printf("In foundProfileForDefault\n")
			//fmt.Printf("profile: %s\n", profile)
			break
		}
	}

	if !foundProfileForDefault && defaultProfile.isActive {
		profiles["default"] = defaultProfile
	}
} //parse

func listProfiles(profiles map[string]Profile) {

	activeLength := 2
	nameLength := 20
	awsAccessKeyIdLength := 20
	regionLength := 10
	outputLength := 10

	fmt.Printf(fs(activeLength), " ")
	fmt.Printf(fs(nameLength), "PROFILE")
	fmt.Printf("    ")
	fmt.Printf(fs(awsAccessKeyIdLength), "AWS_ACCESS_KEY_ID")
	fmt.Printf("    ")
	fmt.Printf(fs(regionLength), "REGION")
	fmt.Printf("    ")
	fmt.Printf(fs(outputLength), "OUTPUT")
	fmt.Printf("\n")

	for sectionName, profile := range profiles {
		if profile.isActive {
			fmt.Printf("* ")
		} else {
			fmt.Printf("  ")
		}
		fmt.Printf(fs(nameLength), sectionName)
		trunc(sectionName, nameLength)

		fmt.Printf(fs(awsAccessKeyIdLength), profile.aws_access_key_id)
		trunc(profile.aws_access_key_id, awsAccessKeyIdLength)

		fmt.Printf(fs(regionLength), profile.region)
		trunc(profile.region, regionLength)

		fmt.Printf(fs(outputLength), profile.output)
		trunc(profile.output, outputLength)

		fmt.Printf("\n")
	}
} //listProfiles

//format string pattern
func fs(l int) string {
	return "%-" + strconv.Itoa(l) + "." + strconv.Itoa(l) + "s"
} //fs

//truncate string and pad with ...
func trunc(s string, l int) {
	if len(s) > l {
		fmt.Printf("... ")
	} else {
		fmt.Printf("    ")
	}
} //trunc

func setDefaultProfile(fromSectionName string) {
	defaultSection := credentialsFile.Section(ini.DefaultSection)

	//make a backup of the current default section so it doesn't get lost
	//the default section is only active if there is no matching profile present
	if defaultProfile.isActive {
		defaultBackupSectionName := "default-" + time.Now().Format("20060102150405")
		credentialsFile.NewSection(defaultBackupSectionName)
		defaultBackupSection := credentialsFile.Section(defaultBackupSectionName)
		for _, key := range defaultSection.Keys() {
			keyName := key.Name()
			value := key.Value()
			defaultBackupSection.NewKey(keyName, value)
		}
	}

	for _, key := range defaultSection.Keys() {
		keyName := key.Name()
		defaultSection.DeleteKey(keyName)
	}

	fromSection := credentialsFile.Section(fromSectionName)
	for _, key := range fromSection.Keys() {
		keyName := key.Name()
		value := key.Value()
		defaultSection.NewKey(keyName, value)
	}

	credentialsFile.SaveTo(usr.HomeDir + "/.aws/credentials")

	fmt.Printf("Activated Profile '%s'\n\n", fromSectionName)

	parse()
} //setDefaultProfile
