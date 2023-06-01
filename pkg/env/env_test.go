package env_test

// import (
// 	"github.com/stretchr/testify/assert"
// 	// "gitlab.neoway.com.br/malibu/go-logger-lib/nlog"
// 	// "gitlab.neoway.com.br/sebrae-sp-ddm-datacapture/pkg/env"
// 	"os"
// 	"testing"
// )

// var (
// 	testLog = nlog.NewLogger("Test-API-test", "", nlog.LevelError)
// )

// func TestGetString(t *testing.T) {
// 	t.Run("ReturnValue", func(t *testing.T) {
// 		err := os.Setenv("VAR_TEST", "custom_value")
// 		assert.NoError(t, err)

// 		defer os.Unsetenv("VAR_TEST")
// 		env.CheckRequired(testLog, "VAR_TEST")

// 		value := env.GetString("VAR_TEST", "default_value")
// 		assert.Equal(t, "custom_value", value)
// 	})

// 	t.Run("ReturnDefaultValue", func(t *testing.T) {
// 		value := env.GetString("VAR_TEST", "default_value")
// 		assert.Equal(t, "default_value", value)
// 	})
// }

// func TestGetInt(t *testing.T) {
// 	t.Run("ReturnValue", func(t *testing.T) {
// 		err := os.Setenv("VAR_TEST", "123")
// 		assert.NoError(t, err)

// 		defer os.Unsetenv("VAR_TEST")

// 		env.CheckRequired(testLog, "VAR_TEST")
// 		value := env.GetInt("VAR_TEST", 999)
// 		assert.Equal(t, 123, value)
// 	})

// 	t.Run("ValueIsNotInt", func(t *testing.T) {
// 		err := os.Setenv("VAR_TEST", "not_integer_value")
// 		assert.NoError(t, err)

// 		defer os.Unsetenv("VAR_TEST")

// 		env.CheckRequired(testLog, "VAR_TEST")
// 		value := env.GetInt("VAR_TEST", 999)
// 		assert.Equal(t, 999, value)
// 	})

// 	t.Run("ReturnDefaultValue", func(t *testing.T) {
// 		value := env.GetInt("VAR_TEST", 222)
// 		assert.Equal(t, 222, value)
// 	})
// }
