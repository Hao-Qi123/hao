package pkg

import (
	"fmt"
	"math/rand"
	"time"
)

func OrderNo() string {
	return time.Now().Format("20060102150405") + fmt.Sprintf("%06d", rand.Intn(1000000))
}
