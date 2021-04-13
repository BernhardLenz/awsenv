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

//TODO: test other keys than aws_access_key_id such as metadata_service_timeout in activate
//TODO: test printing of default configs as part of a Profile line
//TODO: test export of e.g. AWS_CCESS_KEY_ID
//TODO: Add versioning and printing of version
//TODO: comment methods

type Profile struct {
	profileName           string
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

//method pointer which can be changed during test case execution
var logFatalf = log.Fatalf
var osExit = os.Exit

func main() {

	listCommand := flag.NewFlagSet("list", flag.ExitOnError) //since list doesn't require parameters, not sure if a FlagSet is needed
	activateCommand := flag.NewFlagSet("activate", flag.ExitOnError)

	if len(os.Args) > 3 {
		fmt.Fprintf(os.Stderr, "ERROR: Too many arguments supplied.\n")
		printUsage()
		osExit(1)
		return //During test case execution osExit may not actually exit
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "list":
			//TODO: handle error
			_ = listCommand.Parse(make([]string, 0))
		case "activate":
			_ = activateCommand.Parse(os.Args[2:])
		case "help", "-help", "--help":
			printUsage()
			osExit(0)
			return //During test case execution osExit may not actually exit
		default:
			fmt.Fprintf(os.Stderr, "ERROR: Unknown command!\n")
			printUsage()
			osExit(1)
			return //During test case execution osExit may not actually exit
		}
	}

	if listCommand.Parsed() || len(os.Args) == 1 {
		parse()
		listProfiles(profiles)

		fmt.Printf("\nTo activate a different Profile run '%s activate <Profile>'", filepath.Base(os.Args[0]))
	} else if activateCommand.Parsed() {
		if len(os.Args) != 3 {
			fmt.Fprintf(os.Stderr, "ERROR: Required parameter <Profile> missing for activate command!\n")
			printUsage()
			osExit(1)
			return //During test case execution osExit may not actually exit
		}
		activateProfileName := os.Args[2]
		if activateProfileName == "default" {
			logFatalf("ERROR: Cannot activate the 'default' Profile as it is already active!")
		}
		parse()

		profile, ok := profiles[activateProfileName]
		if !ok {
			fmt.Fprintf(os.Stderr, "ERROR: Profile '%s' does not exist! Available profiles are: \n\n", activateProfileName)
			listProfiles(profiles)
			osExit(1)
			return //During test case execution osExit may not actually exit
		}

		if profile.isActive {
			fmt.Fprintf(os.Stderr, "Profile '%s' is already active! No changes applied. \n\n", activateProfileName)
			listProfiles(profiles)
			osExit(0)
			return //During test case execution osExit may not actually exit
		}

		setDefaultProfile(activateProfileName)

		listProfiles(profiles)

		setEnvironmentVariables()
	}
} //main

func init() {
	ini.DefaultSection = "default"
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Printf("%s [list]\n", filepath.Base(os.Args[0]))
	fmt.Println(" Lists all available profiles.")
	fmt.Printf("%s activate <Profile>\n", filepath.Base(os.Args[0]))
	fmt.Println(" Activates a given Profile.")
	fmt.Println("")
	fmt.Println("To create a new Profile use 'aws configure''")
} //printUsage

func parse() {

	parseCredentials()
	parseConfig()

} //parse

func parseCredentials() {
	credentialsFile = loadIni(getCredentialFilePath())
	//During test case execution loadIni may actually return a nil if it doesn't find a file
	if credentialsFile == nil {
		return
	}

	//creates new empty defaultCredentialsSection if it doesn't exist already
	defaultCredentialsSection := credentialsFile.Section(ini.DefaultSection)
	defaultProfile.profileName = ini.DefaultSection
	defaultProfile.aws_access_key_id = defaultCredentialsSection.Key("aws_access_key_id").Value()
	defaultProfile.aws_secret_access_key = defaultCredentialsSection.Key("aws_secret_access_key").Value()

	for _, credentialsSection := range credentialsFile.Sections() {
		sectionName := credentialsSection.Name()

		var profile Profile
		profile.profileName = sectionName

		for _, key := range credentialsSection.Keys() {
			keyName := key.Name()
			value := key.Value()
			if "aws_access_key_id" == keyName {
				profile.aws_access_key_id = value
			} else if "aws_secret_access_key" == keyName {
				profile.aws_secret_access_key = value
			}
		}
		if "default" != sectionName {
			profiles[sectionName] = profile
		}
	}

	//Now mark the profiles that match the default as active
	activeProfileFound := false
	for profileName, profile := range profiles {
		if profile.aws_access_key_id == defaultProfile.aws_access_key_id && profile.aws_access_key_id != "" {
			profile.isActive = true
			profiles[profileName] = profile
			activeProfileFound = true
		}
	}

	//if there is no matching Profile then add the default Profile to the profiles map and make it active
	if !activeProfileFound && defaultProfile.aws_access_key_id != "" {
		defaultProfile.isActive = true
		profiles["default"] = defaultProfile
	}

} //parseCredentials

func parseConfig() {

	configFile = loadIni(getConfigFilePath())
	//During test case execution loadIni may actually return a nil if it doesn't find a file
	if configFile == nil {
		return
	}

	defaultConfigSection := configFile.Section(ini.DefaultSection)
	defaultConfig.region = defaultConfigSection.Key("region").Value()
	defaultConfig.output = defaultConfigSection.Key("output").Value()

	for _, configSection := range configFile.Sections() {
		var config Config
		sectionName := configSection.Name()
		for _, key := range configSection.Keys() {
			keyName := key.Name()
			value := key.Value()

			if "output" == keyName {
				config.output = value
			} else if "region" == keyName {
				config.region = value
			}

			//	if "default" == sectionName {
			//		if "output" == keyName {
			//			defaultConfig.output = value
			//			defaultProfile.output = value
			//			if Profile, ok := profiles[ini.DefaultSection]; ok {
			//				Profile.output = value
			//				profiles[ini.DefaultSection] = Profile
			//			}
			//		} else if "region" == keyName {
			//			defaultConfig.region = value
			//			defaultProfile.region = value
			//			if Profile, ok := profiles[ini.DefaultSection]; ok {
			//				Profile.region = value
			//				profiles[ini.DefaultSection] = Profile
			//			}
			//		}
			//	} else {
			//		Profile, ok := profiles[sectionName]
			//		if ok {
			//			if "output" == keyName {
			//				Profile.output = value
			//			} else if "region" == keyName {
			//				Profile.region = value
			//			}
			//			profiles[sectionName] = Profile
			//		}
			//	}
			//}
			configs[sectionName] = config
		}
	}

	//Now set the corresponding fields in profiles
	for sectionName, profile := range profiles {
		if config, ok := configs[sectionName]; ok {
			profile.output = config.output
			profile.region = config.region
			profiles[sectionName] = profile
		}
	}
	defaultProfile.output = defaultConfig.output
	defaultProfile.region = defaultConfig.region

	////Mark the profiles matching the default Profile including the default Profile itself
	////as active
	//foundProfileMatchingDefault := false
	//for sectionName, Profile := range profiles {
	//	if Profile.aws_access_key_id == defaultProfile.aws_access_key_id {
	//		foundProfileMatchingDefault = true
	//		Profile.isActive = true
	//		//set region and output to "" as the default output will be used instead of the
	//		//values of this Profile
	//		if sectionName != ini.DefaultSection {
	//			Profile.region = ""
	//			Profile.output = ""
	//		}
	//		profiles[sectionName] = Profile
	//		//break
	//	}
	//}
	//
	////Only a default Profile exists but no matching named Profile
	////Previously the default Profile wasn't added to the map
	////but since the default Profile in this case is a "stand alone" Profile
	////it needs to be added
	//if !foundProfileMatchingDefault && defaultProfile.isActive {
	//	profiles["default"] = defaultProfile
	//}

} //parseConfig

func loadIni(fileName string) *ini.File {
	file, err := ini.Load(fileName)
	if err != nil {
		errMsg := "Failed to find or read file: " + fileName + ". %v"
		logFatalf(errMsg, err)
	}
	return file
} //loadIni

func getCredentialFilePath() string {
	return getAwsCliFilePath("AWS_SHARED_CREDENTIALS_FILE", "credentials")
} //getCredentialFilePath

func getConfigFilePath() string {
	return getAwsCliFilePath("AWS_CONFIG_FILE", "config")
} //getConfigFilePath

func getAwsCliFilePath(ENV string, fileName string) string {
	ENVVAL := os.Getenv(ENV)
	if ENVVAL != "" {
		return ENVVAL
	}

	return getUser().HomeDir + "/.aws/" + fileName
} //getAwsCliFilePath

func getUser() *user.User {
	usr, err := user.Current()
	if err != nil {
		logFatalf("%v", err)
	}
	return usr
} //getUser

func listProfiles(profiles map[string]Profile) {

	activeLength := 2
	nameLength := 20
	awsAccessKeyIdLength := 20
	regionLength := 15
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
		truncPrintf(sectionName, nameLength)

		fmt.Printf(fs(awsAccessKeyIdLength), maskAccessKey(profile.aws_access_key_id, awsAccessKeyIdLength))
		fmt.Printf("    ") //perhaps only one blank here

		if profile.region == "" {
			if defaultConfig.region != "" {
				fmt.Printf("[%."+strconv.Itoa(regionLength-2)+"s", defaultConfig.region)
				if len(defaultConfig.region) > regionLength-2 {
					fmt.Printf("...")
				}
				fmt.Printf("] ")
				fmt.Printf(fs(regionLength-len(defaultConfig.region)-3+4), "                  ")
			} else {
				fmt.Printf(fs(regionLength), profile.region)
				truncPrintf(profile.region, regionLength)
			}
		} else {
			fmt.Printf(fs(regionLength), profile.region)
			truncPrintf(profile.region, regionLength)
		}

		if profile.output == "" {
			if defaultConfig.output != "" {
				fmt.Printf("[%."+strconv.Itoa(outputLength-2)+"s", defaultConfig.output)
				if len(defaultConfig.output) > outputLength-2 {
					fmt.Printf("...")
				}
				fmt.Printf("] ")
				fmt.Printf(fs(outputLength-len(defaultConfig.output)-3+4), "                  ")
			} else {
				fmt.Printf(fs(outputLength), profile.output)
				truncPrintf(profile.output, outputLength)
			}
		} else {
			fmt.Printf(fs(outputLength), profile.output)
			truncPrintf(profile.output, outputLength)
		}

		fmt.Printf("\n")
	}

	if len(profiles) > 0 {
		fmt.Printf("\n")
		fmt.Println("Profiles with * are active profiles. Profiles with region or output in [] are using the default config.")
	}
} //listProfiles

//format string pattern to eg %-10.10s
func fs(l int) string {
	//- for left justify
	//cut off after first number
	//pad to last number
	return "%-" + strconv.Itoa(l) + "." + strconv.Itoa(l) + "s"
} //fs

//TODO: rename func
//truncate string longer than l and if longer pad with "... " otherwise pad with "    "
func truncPrintf(s string, l int) {
	if len(s) > l {
		fmt.Printf("... ")
	} else {
		fmt.Printf("    ")
	}
} //truncPrintf

func maskAccessKey(s string, l int) string {
	var r string
	if len(s) <= 4 {
		r = s
	} else if len(s) <= l {
		last4 := s[len(s)-4:]
		prefix := ""
		for i := 0; i < len(s)-4; i++ {
			prefix += "*"
		}
		r = prefix + last4
	} else {
		last4 := s[len(s)-4:]
		prefix := ""
		for i := 0; i < l-4-3; i++ {
			prefix += "*"
		}
		r = "..." + prefix + last4
	}
	return r
} //maskAccessKey

func setDefaultProfile(fromSectionName string) {
	defaultSection := credentialsFile.Section(ini.DefaultSection)

	//make a backup of the current default section so it doesn't get lost
	//the default section is only active if there is no matching Profile present
	if defaultProfile.isActive {
		defaultBackupSectionName := "default-" + time.Now().Format("20060102150405")
		_, _ = credentialsFile.NewSection(defaultBackupSectionName)
		defaultBackupSection := credentialsFile.Section(defaultBackupSectionName)
		for _, key := range defaultSection.Keys() {
			keyName := key.Name()
			value := key.Value()
			_, _ = defaultBackupSection.NewKey(keyName, value)
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
		_, _ = defaultSection.NewKey(keyName, value)
	}

	//TODO: handle error
	_ = credentialsFile.SaveTo(getUser().HomeDir + "/.aws/credentials")

	fmt.Printf("Activated Profile '%s'\n\n", fromSectionName)

	parse()
} //setDefaultProfile

func setEnvironmentVariables() {
	if defaultProfile.aws_access_key_id != "" {
		os.Setenv("AWS_ACCESS_KEY_ID", defaultProfile.aws_access_key_id)
		os.Setenv("AWS_SECRET_ACCESS_KEY", defaultProfile.aws_secret_access_key)
		//TODO:
		//set AWS_DEFAULT_REGION
	}
} //setEnvironmentVariables
