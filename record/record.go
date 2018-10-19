package record

type Record struct {
	size           int64 // the size of the record in bytes when serialized
	SequenceNumber int64
	Body           []byte
}
