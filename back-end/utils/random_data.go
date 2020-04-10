package utils

import (
    "math/rand"
    "time"
)

// It seeds the random data generator once in the program.
func InitializeRandomization() {
    rand.Seed(time.Now().UnixNano())
}

// It generates a random string with a specific length based on chars in the range a-z-A-Z.
func GenerateRandomString(length int) string {
    var charset string
    var bytes []byte
    var i int

    charset = "abcdefghijklmnopqrstuvwxyz" +
        "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

    bytes = make([]byte, length)

    for i = 0; i < length; i++ {
        bytes[i] = charset[rand.Intn(len(charset))]
    }

    return string(bytes)
}

// It generates an integer higher or equal to minValue and smaller than maxValue.
func GenerateRandomInteger(minValue int, maxValue int) int {
    return rand.Intn(maxValue-minValue) + minValue
}
