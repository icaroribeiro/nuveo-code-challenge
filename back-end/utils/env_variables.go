package utils

import (
    "fmt"
    "github.com/joho/godotenv"
    "os"
)

func GetEnvVariables(filenames []string, envVariablesMap map[string]string) error {
    var err error
    var key string

    err = godotenv.Load(filenames...)

    if err != nil {
        return fmt.Errorf("Failed to load the env file(s): %s", err.Error())
    }

    for key, _ = range envVariablesMap {
        envVariablesMap[key] = os.Getenv(key)

        if envVariablesMap[key] == "" {
            return fmt.Errorf("Failed to read the %s environment variable: it isn't set", key)
        }
    }

    return nil
}
