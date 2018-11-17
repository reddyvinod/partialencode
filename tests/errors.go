package tests

//partialencode:json
type ErrorIntSlice []int

//partialencode:json
type ErrorBoolSlice []bool

//partialencode:json
type ErrorUintSlice []uint

//partialencode:json
type ErrorStruct struct {
	Int      int    `json:"int"`
	String   string `json:"string"`
	Slice    []int  `json:"slice"`
	IntSlice []int  `json:"int_slice"`
}

type ErrorNestedStruct struct {
	ErrorStruct ErrorStruct `json:"error_struct"`
	Int         int         `json:"int"`
}

//partialencode:json
type ErrorIntMap map[uint32]string
