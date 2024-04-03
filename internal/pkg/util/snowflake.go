package util

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/snowflake"
)

var (
	snowflakeNode *snowflake.Node
	mu            sync.Mutex
	once          sync.Once
)

func GetSnowflakeID() string {
	once.Do(func() {
		snowflakeNode, _ = snowflake.NewNode(1)
	})
	mu.Lock()
	defer mu.Unlock()
	return snowflakeNode.Generate().String()
}

func GenValidateCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())
	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}
