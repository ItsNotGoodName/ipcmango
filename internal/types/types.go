package types

type Key struct {
	ID   int64
	UUID string
}

type Timestamp struct {
	Created_At Time
	Updated_At Time
}
