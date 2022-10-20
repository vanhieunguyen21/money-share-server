package util

import "math/rand"

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const numericBytes = "0123456789"

func RandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = letterBytes[rand.Int63() % int64(len(letterBytes))]
	}
	return string(b)
}

func RandomStringRange(minLength int, maxLength int) string {
	if maxLength < minLength {
		panic("maxLength cannot be smaller than minLength")
	}
	return RandomString(rand.Intn(maxLength-minLength) + minLength)
}

func RandomNumericString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = numericBytes[rand.Int63() % int64(len(numericBytes))]
	}
	return string(b)
}

func RandomNumericStringRange(minLength int, maxLength int) string {
	if maxLength < minLength {
		panic("maxLength cannot be smaller than minLength")
	}
	return RandomNumericString(rand.Intn(maxLength-minLength) + minLength)
}