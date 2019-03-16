package Utils

import (
	"fmt"
	"math/rand"
)

func RandCode() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}
