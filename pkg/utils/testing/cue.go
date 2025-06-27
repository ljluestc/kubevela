package testing

import (
	"os"
	"path/filepath"
	"testing"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/load"
	"github.com/stretchr/testify/assert"
)

// LoadCUEFilesIntoValue loads CUE files into a cue.Value for testing
func LoadCUEFilesIntoValue(t *testing.T, files []string) cue.Value {
	assert := assert.New(t)
	
	cfg := &load.Config{
		ModuleRoot: t.TempDir(),
	}
	
	instances := load.Instances(files, cfg)
	assert.NotEmpty(instances, "expected at least one CUE instance to be loaded")
	
	r := cue.Runtime{}
	instance, err := r.Build(instances[0])
	assert.NoError(err, "should compile CUE instance")
	
	return instance.Value()
}

// CreateTempCUEFile creates a temporary CUE file with the given content
func CreateTempCUEFile(t *testing.T, content string) string {
	assert := assert.New(t)
	
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.cue")
	
	err := os.WriteFile(filePath, []byte(content), 0644)
	assert.NoError(err, "should write temporary CUE file")
	
	return filePath
}

// SetupCUETestEnvironment sets up the necessary environment variables for CUE tests
func SetupCUETestEnvironment() {
	// Ensure CUE_REGISTRY is set to a valid value for tests
	os.Setenv("CUE_REGISTRY", "memory://")
}
