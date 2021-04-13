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
			fmt.Fprint(os.Stderr, format)
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

	if len(configs) != 1 {
		t.Errorf("TestEmptyFile: len(configs) is not 1: %d ", len(configs))
		t.Fail()
	}

	config := configs[ini.DefaultSection]

	if config.output != "" {
		t.Errorf("TestOnlyDefaultSection: default config.output is not '': %s ", config.output)
		t.Fail()
	}

	if config.region != "" {
		t.Errorf("TestOnlyDefaultSection: default config.region is not '': %s ", config.region)
		t.Fail()
	}

} //TestEmptyFile

func TestEmptyConfigFile(t *testing.T) {
	t.Cleanup(resetState)

	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", "./testdata/one_profile_matching_default_credentials")

	parseCredentials()

	os.Setenv("AWS_CONFIG_FILE", "./testdata/empty_file")
	parseConfig()

	if defaultConfig.output != "" {
		t.Errorf("TestOnlyDefaultSection: defaultConfig.output is not '': %s ", defaultConfig.output)
		t.Fail()
	}

	if defaultConfig.region != "" {
		t.Errorf("TestOnlyDefaultSection: defaultConfig.region is not 'us-east-1': %s ", defaultConfig.region)
		t.Fail()
	}

	config := configs[ini.DefaultSection]
	profile := profiles[ini.DefaultSection]

	if config.output != "" {
		t.Errorf("TestOnlyDefaultSection: default config.output is not '': %s ", config.output)
		t.Fail()
	}

	if config.region != "" {
		t.Errorf("TestOnlyDefaultSection: default config.region is not '': %s ", config.region)
		t.Fail()
	}

	if defaultProfile.region != "" {
		t.Errorf("TestOnlyDefaultSection: defaultProfile.region is not '': %s ", defaultProfile.region)
		t.Fail()
	}

	if defaultProfile.output != "" {
		t.Errorf("TestOnlyDefaultSection: defaultProfile.output is not '': %s ", defaultProfile.output)
		t.Fail()
	}

	if profile.region != "" {
		t.Errorf("TestOnlyDefaultSection: default Profile.region is not '': %s ", profile.region)
		t.Fail()
	}

	if profile.output != "" {
		t.Errorf("TestOnlyDefaultSection: default Profile.output is not '': %s ", profile.output)
		t.Fail()
	}
} //TestOnlyDefaultSection

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

	if defaultProfile.region != "" {
		t.Errorf("TestOnlyDefaultSection: defaultProfile.region is not '': %s ", defaultProfile.region)
		t.Fail()
	}

	if defaultProfile.output != "" {
		t.Errorf("TestOnlyDefaultSection: defaultProfile.output is not '': %s ", defaultProfile.output)
		t.Fail()
	}

	if len(profiles) != 1 {
		t.Errorf("TestOnlyDefaultSection: len(profiles) is not 1: %d ", len(profiles))
		t.Fail()
	}

	profile := profiles[ini.DefaultSection]

	if profile.profileName != "default" {
		t.Errorf("TestOnlyDefaultSection: Profile.profileName is not 'default': %s ", profile.profileName)
		t.Fail()
	}

	if profile.aws_access_key_id != "12345678901234567890" {
		t.Errorf("TestOnlyDefaultSection: Profile.aws_access_key_id is not '12345678901234567890': %s ", profile.aws_access_key_id)
		t.Fail()
	}

	if !profile.isActive {
		t.Errorf("TestOnlyDefaultSection: default Profile.isActive is not true: %t ", profile.isActive)
		t.Fail()
	}

	if profile.region != "" {
		t.Errorf("TestOnlyDefaultSection: default Profile.region is not '': %s ", profile.region)
		t.Fail()
	}

	if profile.output != "" {
		t.Errorf("TestOnlyDefaultSection: default Profile.output is not '': %s ", profile.output)
		t.Fail()
	}

	os.Setenv("AWS_CONFIG_FILE", "./testdata/only_default_config")
	parseConfig()

	if defaultConfig.output != "json" {
		t.Errorf("TestOnlyDefaultSection: defaultConfig.output is not 'json': %s ", defaultConfig.output)
		t.Fail()
	}

	if defaultConfig.region != "us-east-1" {
		t.Errorf("TestOnlyDefaultSection: defaultConfig.region is not 'us-east-1': %s ", defaultConfig.region)
		t.Fail()
	}

	if len(configs) != 1 {
		t.Errorf("TestOnlyDefaultSection: len(configs) is not 1: %d ", len(configs))
		t.Fail()
	}

	config := configs[ini.DefaultSection]
	profile = profiles[ini.DefaultSection]

	if config.output != "json" {
		t.Errorf("TestOnlyDefaultSection: default config.output is not 'json': %s ", config.output)
		t.Fail()
	}

	if config.region != "us-east-1" {
		t.Errorf("TestOnlyDefaultSection: default config.region is not 'us-east-1': %s ", config.region)
		t.Fail()
	}

	if defaultProfile.region != "us-east-1" {
		t.Errorf("TestOnlyDefaultSection: defaultProfile.region is not 'us-east-1': %s ", defaultProfile.region)
		t.Fail()
	}

	if defaultProfile.output != "json" {
		t.Errorf("TestOnlyDefaultSection: defaultProfile.output is not 'json': %s ", defaultProfile.output)
		t.Fail()
	}

	if profile.region != "us-east-1" {
		t.Errorf("TestOnlyDefaultSection: default Profile.region is not 'us-east-1': %s ", profile.region)
		t.Fail()
	}

	if profile.output != "json" {
		t.Errorf("TestOnlyDefaultSection: default Profile.output is not 'json': %s ", profile.output)
		t.Fail()
	}
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
		t.Errorf("TestProfileMatchingDefault: len(profiles) is not 1: %d \n %v", len(profiles), profiles)
		t.Fail()
	}

	profile := profiles["profile_matching_default_credentials"]

	if profile.profileName != "profile_matching_default_credentials" {
		t.Errorf("TestProfileMatchingDefault: profile_matching_default_credentials Profile.profileName is not 'profile_matching_default_credentials': %s ", profile.profileName)
		t.Fail()
	}

	if profile.aws_access_key_id != "12345678901234567890" {
		t.Errorf("TestProfileMatchingDefault: profile_matching_default_credentials Profile.aws_access_key_id is not '12345678901234567890': %s ", profile.aws_access_key_id)
		t.Fail()
	}

	if !profile.isActive {
		t.Errorf("TestProfileMatchingDefault: profile_matching_default_credentials Profile.isActive is not true: %t ", profile.isActive)
		t.Fail()
	}
} //TestProfileMatchingDefault

func TestMultipleProfileMatchingDefault(t *testing.T) {
	t.Cleanup(resetState)

	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", "./testdata/two_profiles_matching_default_credentials")

	parseCredentials()

	if defaultProfile.profileName != "default" {
		t.Errorf("TestMultipleProfileMatchingDefault: defaultProfile.profileName is not 'default': %s ", defaultProfile.profileName)
		t.Fail()
	}

	if defaultProfile.isActive {
		t.Errorf("TestMultipleProfileMatchingDefault: defaultProfile.isActive is not false: %t ", defaultProfile.isActive)
		t.Fail()
	}

	if len(profiles) != 3 {
		t.Errorf("TestMultipleProfileMatchingDefault: len(profiles) is not 3: %d \n %v", len(profiles), profiles)
		t.Fail()
	}

	profile := profiles["profile1_matching_default_credentials"]

	if profile.profileName != "profile1_matching_default_credentials" {
		t.Errorf("TestMultipleProfileMatchingDefault: profile1_matching_default_credentials Profile.profileName is not 'profile_matching_default_credentials': %s ", profile.profileName)
		t.Fail()
	}

	if !profile.isActive {
		t.Errorf("TestMultipleProfileMatchingDefault: profile1_matching_default_credentials Profile.isActive is not true: %t ", profile.isActive)
		t.Fail()
	}

	profile = profiles["profile2_matching_default_credentials"]

	if profile.profileName != "profile2_matching_default_credentials" {
		t.Errorf("TestMultipleProfileMatchingDefault: profile2_matching_default_credentials Profile.profileName is not 'profile2_matching_default_credentials': %s ", profile.profileName)
		t.Fail()
	}

	if !profile.isActive {
		t.Errorf("TestMultipleProfileMatchingDefault: profile2_matching_default_credentials Profile.isActive is not true: %t ", profile.isActive)
		t.Fail()
	}

	profile = profiles["profile3_not_matching_default_credentials"]

	if profile.profileName != "profile3_not_matching_default_credentials" {
		t.Errorf("TestMultipleProfileMatchingDefault: profile3_not_matching_default_credentials Profile.profileName is not 'profile3_matching_default_credentials': %s ", profile.profileName)
		t.Fail()
	}

	if profile.isActive {
		t.Errorf("TestMultipleProfileMatchingDefault: profile3_not_matching_default_credentials Profile.isActive is true: %t ", profile.isActive)
		t.Fail()
	}
} //TestMultipleProfileMatchingDefault

func TestDuplicateProfiles(t *testing.T) {
	t.Cleanup(resetState)

	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", "./testdata/duplicate_profiles_credentials")

	parseCredentials()

	if defaultProfile.profileName != "default" {
		t.Errorf("TestDuplicateProfiles: defaultProfile.profileName is not 'default': %s ", defaultProfile.profileName)
		t.Fail()
	}

	//ini picks the 2nd default profile
	if defaultProfile.aws_access_key_id != "12345678901234567891" {
		t.Errorf("TestDuplicateProfiles: defaultProfile.aws_access_key_id is not '12345678901234567891': %s ", defaultProfile.aws_access_key_id)
		t.Fail()
	}

	if !defaultProfile.isActive {
		t.Errorf("TestDuplicateProfiles: defaultProfile.isActive is not active: %t ", defaultProfile.isActive)
		t.Fail()
	}

	if len(profiles) != 2 {
		t.Errorf("TestDuplicateProfiles: len(profiles) is not 2: %d \n %v", len(profiles), profiles)
		t.Fail()
	}

	profile := profiles["duplicate"]

	if profile.profileName != "duplicate" {
		t.Errorf("TestDuplicateProfiles: duplicate Profile.profileName is not 'duplicate': %s ", profile.profileName)
		t.Fail()
	}

	//ini picks the 2nd duplicate profile
	if profile.aws_access_key_id != "12345678901234567893" {
		t.Errorf("TestDuplicateProfiles: duplicate Profile.aws_access_key_id is not '12345678901234567893': %s ", defaultProfile.aws_access_key_id)
		t.Fail()
	}

	if profile.isActive {
		t.Errorf("TestDuplicateProfiles: duplicate Profile.isActive is active: %t ", profile.isActive)
		t.Fail()
	}

} //TestDuplicateProfiles

func TestDuplicateMixedCaseProfiles(t *testing.T) {
	t.Cleanup(resetState)

	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", "./testdata/duplicate_mixed_case_profiles_credentials")

	parseCredentials()

	if defaultProfile.profileName != "default" {
		t.Errorf("TestDuplicateMixedCaseProfiles: defaultProfile.profileName is not 'default': %s ", defaultProfile.profileName)
		t.Fail()
	}

	//awsenv configures ini picks the lower case default profile
	if defaultProfile.aws_access_key_id != "12345678901234567890" {
		t.Errorf("TestDuplicateMixedCaseProfiles: defaultProfile.aws_access_key_id is not '12345678901234567891': %s ", defaultProfile.aws_access_key_id)
		t.Fail()
	}

	if !defaultProfile.isActive {
		t.Errorf("TestDuplicateMixedCaseProfiles: defaultProfile.isActive is not active: %t ", defaultProfile.isActive)
		t.Fail()
	}

	if len(profiles) != 4 {
		t.Errorf("TestDuplicateMixedCaseProfiles: len(profiles) is not 4 %d \n %v", len(profiles), profiles)
		t.Fail()
	}

	profile := profiles["DEFAULT"]

	if profile.profileName != "DEFAULT" {
		t.Errorf("TestDuplicateMixedCaseProfiles: DEFAULT Profile.profileName is not 'DEFAULT': %s ", profile.profileName)
		t.Fail()
	}

	if profile.aws_access_key_id != "12345678901234567891" {
		t.Errorf("TestDuplicateMixedCaseProfiles: DEFAULT Profile.aws_access_key_id is not '12345678901234567891': %s ", defaultProfile.aws_access_key_id)
		t.Fail()
	}

	if profile.isActive {
		t.Errorf("TestDuplicateMixedCaseProfiles: DEFAULT Profile.isActive is active: %t ", profile.isActive)
		t.Fail()
	}

	profile = profiles["duplicate"]

	if profile.profileName != "duplicate" {
		t.Errorf("TestDuplicateMixedCaseProfiles: duplicate Profile.profileName is not 'duplicate': %s ", profile.profileName)
		t.Fail()
	}

	if profile.aws_access_key_id != "12345678901234567892" {
		t.Errorf("TestDuplicateMixedCaseProfiles: duplicate Profile.aws_access_key_id is not '12345678901234567892': %s ", profile.aws_access_key_id)
		t.Fail()
	}

	if profile.isActive {
		t.Errorf("TestDuplicateMixedCaseProfiles: duplicate Profile.isActive is active: %t ", profile.isActive)
		t.Fail()
	}

	profile = profiles["DUPLICATE"]

	if profile.profileName != "DUPLICATE" {
		t.Errorf("TestDuplicateMixedCaseProfiles: DUPLICATE Profile.profileName is not 'DUPLICATE': %s ", profile.profileName)
		t.Fail()
	}

	if profile.aws_access_key_id != "12345678901234567893" {
		t.Errorf("TestDuplicateMixedCaseProfiles: DUPLICATE Profile.aws_access_key_id is not '12345678901234567892': %s ", profile.aws_access_key_id)
		t.Fail()
	}

	if profile.isActive {
		t.Errorf("TestDuplicateMixedCaseProfiles: DUPLICATE Profile.isActive is active: %t ", profile.isActive)
		t.Fail()
	}

} //TestDuplicateMixedCaseProfiles

func resetState() {
	defaultProfile = Profile{}
	profiles = make(map[string]Profile)

	defaultConfig = Config{}
	configs = make(map[string]Config)
}
