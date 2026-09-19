/*
Package config offers the functionality of Punch configuration files
and creating unique test IDs. Alongside with the types and structs
that the rest of Punch uses throughout the entire application.
*/
package config

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
)

// ParsePunchConfig reads, parses, and returns the Punch
// configuration file.
func ParsePunchConfig(directory string) (*ConfigFile, error) {
	// Attempt to change directory
	err := os.Chdir(directory)
	if err != nil {
		return &ConfigFile{}, err
	}

	// Attempt to read Punch configuration file
	file, err := os.ReadFile(configFileName)
	if err != nil {
		return &ConfigFile{}, fmt.Errorf(
			"no configuration file provided: %w",
			err,
		)
	}

	var config ConfigFile
	// Attempt to parse Punch Configuration file
	if err := json.Unmarshal(file, &config); err != nil {
		return &ConfigFile{}, fmt.Errorf(
			"failed to unmarshal %s: %w",
			configFileName,
			err,
		)
	}

	return &config, nil
}

// GenerateTestID creates and returns a unique test ID of config.TestIDLen
// characters.
func GenerateTestID() TestID {
	return TestID(rand.Text()[0:testIDLen])
}

// SetupTestMetadata is a part of the initialization process:
// It generates a unique test ID (with included retry-handling),
// persists the ID and the clients configuration to Redis.
// Disclaimer: Currently unimplemented; SetupTestMetadata steps
// are included in the body of the function.
func SetupTestMetadata(d *ClientTestData) (testID TestID, err error) {
	// Temporary use of variable ClientTestData
	fmt.Printf("%s", d.TestID[0:0])
	// Steps:
	// 1. Create test ID w/ retry handling
	// (Test ID format: TEST-12CHARACTERS)
	// 2. Persist test ID and configuration to Redis
	// 3. Return persisted test ID

	// Note: Planned use.
	return "", nil
}
