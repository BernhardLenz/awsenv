package main

import (
	"fmt"
	"github.com/BernhardLenz/ini"
	"os"
	"testing"
)

func TestInvalidFile(t *testing.T) {
	t.Cleanup(resetState)

	origLogFatalf := logFatalf

	// After this test, replace the original fatal function
	defer func() { logFatalf = origLogFatalf }()

	errors := []string{}
	logFatalf = func(format string, args ...interface{}) {
		if len(args) > 0 {
			fmt.Fprintf(os.Stderr, format, args)
			fmt.Fprintln(os.Stderr)
			errors = append(errors, fmt.Sprintf(format, args))
		} else {
			fmt.Fprintf(os.Stderr, format)
			fmt.Fprintln(os.Stderr)
			errors = append(errors, format)
		}
	}

	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", "./testdata/nonexisting_file")
	parseCredentials()

	if len(errors) != 1 {
		t.Errorf("TestInvalidFile: Excepted one error for invalid AWS_SHARED_CREDENTIALS_FILE, actual %v", len(errors))
	}

	errors = []string{}

	os.Setenv("AWS_CONFIG_FILE", "./testdata/nonexisting_file")
	parseConfig()

	if len(errors) != 1 {
		t.Errorf("TestInvalidFile: Excepted one error for invalid AWS_CONFIG_FILE, actual %v", len(errors))
	}
} //TestInvalidFile

func TestEmptyFile(t *testing.T) {
	t.Cleanup(resetState)

	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", "./testdata/empty_file")
	parseCredentials()

	if defaultProfile.profileName != "default" {
		t.Errorf("TestEmptyFile: defaultProfile.profileName is not 'default': %s ", defaultProfile.profileName)
		t.Fail()
	}

	if defaultProfile.aws_access_key_id != "" {
		t.Errorf("TestEmptyFile: defaultProfile.aws_access_key_id is not '': %s ", defaultProfile.aws_access_key_id)
		t.Fail()
	}

	if len(profiles) != 0 {
		t.Errorf("TestEmptyFile: len(profiles) is not 0: %d ", len(profiles))
		t.Fail()
	}

	os.Setenv("AWS_CONFIG_FILE", "./testdata/empty_file")
	parseConfig()

	if defaultConfig.region != "" {
		t.Errorf("TestEmptyFile: defaultConfig.region is not '': %s ", defaultConfig.region)
		t.Fail()
	}

	if len(configs) != 0 {
		t.Errorf("TestEmptyFile: len(configs) is not 0: %d ", len(configs))
		t.Fail()
	}
} //TestEmptyFile

func TestOnlyDefaultSection(t *testing.T) {
	t.Cleanup(resetState)

	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", "./testdata/only_default_credentials")

	parseCredentials()

	if defaultProfile.profileName != "default" {
		t.Errorf("TestOnlyDefaultSection: defaultProfile.profileName is not 'one_profile_matching_default_credentials': %s ", defaultProfile.profileName)
		t.Fail()
	}

	if defaultProfile.aws_access_key_id != "12345678901234567890" {
		t.Errorf("TestOnlyDefaultSection: defaultProfile.aws_access_key_id is not '12345678901234567890': %s ", defaultProfile.aws_access_key_id)
		t.Fail()
	}

	if !defaultProfile.isActive {
		t.Errorf("TestOnlyDefaultSection: defaultProfile.isActive is not true: %t ", defaultProfile.isActive)
		t.Fail()
	}

	if len(profiles) != 1 {
		t.Errorf("TestOnlyDefaultSection: len(profiles) is not 1: %d ", len(profiles))
		t.Fail()
	}

	profile := profiles[ini.DefaultSection]

	if profile.profileName != "default" {
		t.Errorf("TestOnlyDefaultSection: Profile.profileName is not 'default': %s ", defaultProfile.profileName)
		t.Fail()
	}

	if profile.aws_access_key_id != "12345678901234567890" {
		t.Errorf("TestOnlyDefaultSection: Profile.aws_access_key_id is not '12345678901234567890': %s ", defaultProfile.aws_access_key_id)
		t.Fail()
	}

	if !profile.isActive {
		t.Errorf("TestOnlyDefaultSection: Profile.isActive is not true: %t ", defaultProfile.isActive)
		t.Fail()
	}

	os.Setenv("AWS_CONFIG_FILE", "./testdata/only_default_config")
	parseConfig()
	//TODO: Test Cases
} //TestOnlyDefaultSection

func TestProfileMatchingDefault(t *testing.T) {
	t.Cleanup(resetState)

	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", "./testdata/one_profile_matching_default_credentials")

	parseCredentials()

	if defaultProfile.profileName != "default" {
		t.Errorf("TestProfileMatchingDefault: defaultProfile.profileName is not 'default': %s ", defaultProfile.profileName)
		t.Fail()
	}

	if defaultProfile.aws_access_key_id != "12345678901234567890" {
		t.Errorf("TestProfileMatchingDefault: defaultProfile.aws_access_key_id is not '12345678901234567890': %s ", defaultProfile.aws_access_key_id)
		t.Fail()
	}

	if defaultProfile.isActive {
		t.Errorf("TestProfileMatchingDefault: defaultProfile.isActive is not false: %t ", defaultProfile.isActive)
		t.Fail()
	}

	if len(profiles) != 1 {
		//mString :=
		//msg := "TestProfileMatchingDefault: len(profiles) is not 1: %d \n " + profiles
		t.Errorf("TestProfileMatchingDefault: len(profiles) is not 1: %d \n %v", len(profiles), profiles)
		t.Fail()
	}

	profile := profiles["profile_matching_default_credentials"]

	if profile.profileName != "profile_matching_default_credentials" {
		t.Errorf("TestProfileMatchingDefault: Profile.profileName is not 'profile_matching_default_credentials': %s ", profile.profileName)
		t.Fail()
	}

	if profile.aws_access_key_id != "12345678901234567890" {
		t.Errorf("TestProfileMatchingDefault: Profile.aws_access_key_id is not '12345678901234567890': %s ", profile.aws_access_key_id)
		t.Fail()
	}

	if !profile.isActive {
		t.Errorf("TestProfileMatchingDefault: Profile.isActive is not true: %t ", profile.isActive)
		t.Fail()
	}
	//
	//
	//os.Setenv("AWS_SHARED_CREDENTIALS_FILE", "./testdata/only_default_config")
	//parseConfig()
} //TestProfileMatchingDefault

func resetState() {
	defaultProfile = Profile{}
	profiles = make(map[string]Profile)

	defaultConfig = Config{}
	configs = make(map[string]Config)
}