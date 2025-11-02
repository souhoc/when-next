package discord

// Snowflake ID Broken Down in Binary:
//
//	111111111111111111111111111111111111111111 11111 11111 111111111111
//	64                                         22    17    12          0
//
//	|FIELD              |BITS    |NUMBER OF BITS|
//	|-------------------|--------|--------------|
//	|Timestamp          |63 to 22|42 bits       |
//	|Internal worker ID |21 to 17|5 bits        |
//	|Internal process ID|16 to 12|5 bits        |
//	|Increment          |11 to 0 |12 bits       |
//
// For more see [Discord docs]
//
// [Discord docs]: https://discord.com/developers/docs/reference#snowflakes
type Snowflake uint64

// Unix returns t as a Unix time, the number of seconds elapsed
// since January 1, 1970 UTC. The result does not depend on the
// location associated with t.
func (s Snowflake) Unix() int64 {
	const epoch = 1420070400000
	return int64((s >> 22) + epoch)
}
