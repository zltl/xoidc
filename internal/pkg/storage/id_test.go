package storage

import (
	"log"
	"testing"

	"github.com/bwmarrin/snowflake"
)

func TestSnowFlakeID(t *testing.T) {
	var id int64 = 0
	x := snowflake.ID(id)
	log.Printf("0: %s", x.Base64())
}
