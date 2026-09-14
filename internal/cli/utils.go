// Package cli is the most important package for Punch's
// CLI-tool as it implements the core functionality it
// relies on. Disclaimer: Some of the types in
// ./cli/types.go can seem to hold duplication. Despite,
// package server and package cli sharing type and struct
// names, they're completely different.
package cli

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
)

// GenerateTestID creates and returns a unique test ID of 12
// characters.
func GenerateTestID() string {
	return rand.Text()[0:12]
}

// SetupTestMetadata is a part of the initialization process:
// It generates a unique test ID (with included retry-handling),
// persists the ID and the clients configuration to Redis.
// Disclaimer: Currently unimplemented; SetupTestMetadata steps
// are included in the body of the function.
func SetupTestMetadata(d *ClientTestData) (testID string, err error) {
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

// ParsePunchConfig reads, parses, and returns the Punch 
// configuration file.
func ParsePunchConfig(directory string) (PunchConfig, error) {
  // Attempt to change directory
  err := os.Chdir(directory)
  if err != nil {
    return PunchConfig{}, fmt.Errorf(
      "failed to change directory to /%s: %s", 
      directory, 
      err.Error(),
    )
  }
  // Attempt to read Punch configuration file
  file, err := os.ReadFile(PunchConfigFileName)
  if err != nil {
    return PunchConfig{}, fmt.Errorf(
      "failed to read Punch configuration: %s", 
      err.Error(),
    )
  }
  
  var config PunchConfig
  // Attempt to parse Punch Configuration file
  if err := json.Unmarshal(file, &config); err != nil {
    return PunchConfig{}, fmt.Errorf( 
      "failed to unmarshal %s: %s", 
      PunchConfigFileName, 
      err.Error(),
    )
  }
  
  return config, nil
}
