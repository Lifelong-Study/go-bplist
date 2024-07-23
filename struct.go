package go_bplist

type PlistNode struct {
	Type         int
	Key          string
	ValueBool    bool
	ValueInteger int
	ValueString  string
	Nodes        []PlistNode
}

const (
	NullOrBool    int = 0x00
	Integer       int = 0x01
	Real          int = 0x02
	Date          int = 0x03
	Dataa         int = 0x04
	String        int = 0x05
	UnicodeString int = 0x06
	UID           int = 0x08
	Array         int = 0x0A
	Dictionary    int = 0x0D
)

type HexFooter struct {
	AlwaysNULL                      Bytes
	SizeOfTheItemsInTheOffsetTable  Bytes
	ObjectReferenceSize             Bytes
	NumberOfObjectsInTheOffsetTable Bytes
	TopOfTheTable                   Bytes
	OffsetTableStart                Bytes
}
