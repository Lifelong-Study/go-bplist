package go_bplist

type Format struct {
	Data               Bytes
	DataCount          int
	RecordStartAddress Bytes
	DataStartAddress   Bytes
	Footer             HexFooter
}

func (format Format) fillFooter() Format {
	byte_size := len(format.Data)
	start := byte_size - 32
	end := byte_size - 0

	footer_bytes := Bytes(format.Data[start:end])

	//
	footer := HexFooter{}
	footer.AlwaysNULL = footer_bytes.GetBytes(0, 6)
	footer.SizeOfTheItemsInTheOffsetTable = footer_bytes.GetBytes(6, 1)
	footer.ObjectReferenceSize = footer_bytes.GetBytes(7, 1)
	footer.NumberOfObjectsInTheOffsetTable = footer_bytes.GetBytes(8, 8)
	footer.TopOfTheTable = footer_bytes.GetBytes(16, 8)
	footer.OffsetTableStart = footer_bytes.GetBytes(24, 8)

	//
	footer.AlwaysNULL.Print()
	footer.SizeOfTheItemsInTheOffsetTable.Print()
	footer.ObjectReferenceSize.Print()
	footer.NumberOfObjectsInTheOffsetTable.Print()
	footer.TopOfTheTable.Print()
	footer.OffsetTableStart.Print()

	//
	format.Footer = footer

	//
	format.DataStartAddress = format.Data.GetBytes(footer.OffsetTableStart.ToInt(), 2)

	//
	DataCount, DataStartAddress := format.Data.GetDataCount(format.DataStartAddress)
	format.DataCount = DataCount
	format.DataStartAddress = DataStartAddress

	//
	format.RecordStartAddress = footer.OffsetTableStart

	return format
}

func (format Format) GetData(i int) PlistNode {

	//
	dataAddress := format.DataStartAddress.ToInt()
	dataOffset := format.Data.GetBytes(dataAddress+i, 1).ToInt()

	//
	recordAddress := format.RecordStartAddress.ToInt() + 2*dataOffset
	recordOffset := format.Data.GetBytes(recordAddress, 2)

	//
	return format.Data.GetValue(recordOffset)
}
