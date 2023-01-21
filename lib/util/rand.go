package util

import (
	"log"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func RandBytes(length int) (bytes []byte) {
	bytes = make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		log.Panic(err)
	}
	return
}
