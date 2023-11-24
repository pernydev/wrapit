package basics

import (
	"strconv"
	"time"
)

type Snowflake struct {
	Timestamp time.Time
	WorkerID  int64
	ProcessID int64
	Increment int64
}

func ParseSF(snowflake string) Snowflake {
	var sf Snowflake

	snowflakeInt, _ := strconv.ParseInt(snowflake, 10, 64)
	sf.Increment = snowflakeInt & 0xFFF
	sf.ProcessID = (snowflakeInt >> 12) & 0x3F
	sf.WorkerID = (snowflakeInt >> 17) & 0x1F
	sf.Timestamp = time.Unix(((snowflakeInt>>22)+1420070400000)/1000, 0)

	return sf
}
