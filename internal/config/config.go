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

	"github.com/yuriongit/punch/internal/domain"
)

// LoadConfig reads, parses, and returns the Punch
// configuration file.
func LoadConfig(dir string) (*domain.ConfigFile, error) {
	// Attempt to change directory
	err := os.Chdir(dir)
	if err != nil {
		return &domain.ConfigFile{}, err
	}

	// Attempt to read Punch configuration cnfFile
	cnfFile, err := os.ReadFile(domain.ConfigFileName)
	if err != nil {
		return &domain.ConfigFile{}, fmt.Errorf(
			"no configuration file provided: %w",
			err,
		)
	}

	var cnf domain.ConfigFile
	// Attempt to parse Punch Configuration file
	if err := json.Unmarshal(cnfFile, &cnf); err != nil {
		return &domain.ConfigFile{}, fmt.Errorf(
			"failed to unmarshal %s: %w",
			domain.ConfigFileName,
			err,
		)
	}

	return &cnf, nil
}

// CreateTestID creates and returns a unique test ID of config.TestIDLen
// characters.
func CreateTestID() domain.TestID {
	return domain.TestID(rand.Text()[0:domain.TestIDLen])
}

// SetupTestMetadata is a part of the initialization process:
// It generates a unique test ID (with included retry-handling),
// persists the ID and the clients configuration to SQLite maybe
// Disclaimer: Currently unimplemented; SetupTestMetadata steps
// are included in the body of the function.
func SetupTestMetadata(d *domain.ClientTestData) (testID domain.TestID, err error) {
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
