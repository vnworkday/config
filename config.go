package config

import (
	"flag"
	"os"
)

const configFile = ".env"

const (
	profileEnvVar  = "profile"
	profileDefault = "local"
)

// LoadConfig reads environment variables and adds them to the config map.
//
// Parameters:
//   - in: A pointer to struct to populate with the config values loaded from the environment and the config file.
//
// Returns:
//   - A pointer to the struct with populated config values, if successful, otherwise nil.
//   - An error if the config values could not be loaded, otherwise nil.
func LoadConfig[T any](in *T) (*T, error) {
	loadProfile()

	err := newBuilder().
		FromEnv().
		FromFile(configFile).
		MapTo(in)
	if err != nil {
		return nil, err
	}

	return in, nil
}

// IsLocal returns true if the profile is set to "local".
func IsLocal() bool {
	profile := GetProfile()

	return profile == profileDefault
}

// GetProfile returns the current profile.
func GetProfile() string {
	profile, found := os.LookupEnv(profileEnvVar)

	if !found {
		return profileDefault
	}

	return profile
}

// LoadProfile reads the profile from the command line arguments and sets it as an environment variable.
func loadProfile() {
	var profile string

	flag.StringVar(&profile, profileEnvVar, profileDefault, "Profile to use for configuration")
	flag.Parse()

	if profile == "" {
		profile = profileDefault
	}

	if err := os.Setenv(profileEnvVar, profile); err != nil {
		panic(err)
	}
}
