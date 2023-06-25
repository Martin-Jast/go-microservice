package utils

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// SetupEnvVars gets a file with the declared env vars and set them in the environment
// envs should be declared as KEY="VALUE" with the value always inside " symbols
func SetupEnvVars(path string) error {
	envs, err := SafeGetFileAsLines(path)
	if err != nil {
		return err
	}

	for _, env := range envs {
		 keyRegex := regexp.MustCompile(`(^[A-Za-z])\w+`)
		 keyName := keyRegex.FindStringSubmatch(env)
		 if len(keyName) != 0 {
			// found a key, now get the value and declare as envVar
			valueRegex := regexp.MustCompile(`\"(.*)\"`)
			value := valueRegex.FindStringSubmatch(env)
			if len(value) == 0 {
				fmt.Printf("No Value set for %s env var. \n", keyName[0])
				continue;
			}
			err := os.Setenv(keyName[0], trimQuotationMarks(value[0]))
			if err != nil {
				fmt.Printf("Error setting %s env var.", keyName[0])
				panic(err)
			}
		 }
	}

	return nil
}

func trimQuotationMarks(s string) string {
	t := s
	quotation:= "\""
	if strings.HasSuffix(t, "\"") {
        t = t[:len(t)-len(quotation)]
    }
	if strings.HasPrefix(t, "\"") {
        t = t[len(quotation):]
    }
	return t
}

// CheckIfNeededVarsAreSet returns error in case of missing env var
func CheckIfNeededVarsAreSet(vars []string, verbose bool) error {
	unset := []string{}
	for _, envVar := range vars {
		val, isSet := os.LookupEnv(envVar)
		if verbose {
			fmt.Printf("LookupEnv: %s variable is set: %t, value: %s \n", envVar, isSet, val)
		}
		if !isSet {
			unset = append(unset, envVar)
		}
	}
	if len(unset) > 0 {
		return fmt.Errorf("envs not set: %s", strings.Join(unset, ", "))
	}
	return nil
}