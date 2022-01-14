package env

import (
	"github.com/hill-daniel/drizzle"
	"github.com/pkg/errors"
	"log"
	"os"
	"strings"
)

// ExpandVars replaces ${var} or $var in the string according to the values
// of the current environment variables or, if not found, with a secret manager.
func ExpandVars(variables map[string]string, secretRetriever drizzle.SecretRetriever) (map[string]string, error) {
	keyToValueMap := make(map[string]string)
	for key, value := range variables {
		if !strings.HasPrefix(value, "$") {
			keyToValueMap[key] = value
			continue
		}
		envValue := os.ExpandEnv(value)
		if envValue != "" {
			keyToValueMap[key] = envValue
			continue
		}
		secret := os.Expand(value, func(s string) string {
			secretValue, err := secretRetriever.RetrieveSecret(s)
			if err != nil {
				log.Println(errors.Wrapf(err, "failed to retrieve value for key %q", value))
				return ""
			}
			return secretValue
		})
		keyToValueMap[key] = secret
	}
	return keyToValueMap, nil
}
