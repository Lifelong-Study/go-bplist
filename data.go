package go_bplist

import (
	"fmt"
	"sort"
)

type Bytes []byte

func (data Bytes) GetBytes(index int, length int) Bytes {
	values := make(Bytes, length)

	for i := 0; i < length; i += 1 {
		values[i] = data[index+i]
	}

	return values
}

func (data Bytes) GetValue(address Bytes) PlistNode {

	//
	addressValue := address.ToInt()
	value := data[addressValue]
	moreOffset := 1
	count := 0
	node := PlistNode{}

	//
	if (value & 0xD0) == 0xD0 {
		count = int(value - 0xD0)

		if count == 0xF {
			count = data.GetBytes(addressValue+1, 2).ToInt()
			count -= 0x1000
			moreOffset = 3
		}

		nodes := []PlistNode{}

		for i := 0; i < count; i++ {

			//
			keyIndex := format.Data.GetBytes(addressValue+i+1, 1)
			keyOffset := format.RecordStartAddress.ToInt() + 2*keyIndex.ToInt()
			keyBytes := format.Data.GetBytes(keyOffset, 2)
			keyNode := format.Data.GetValue(keyBytes)

			//
			valueIndex := format.Data.GetBytes(addressValue+i+1+count, 1)
			valueOffset := format.RecordStartAddress.ToInt() + 2*valueIndex.ToInt()
			valueBytes := format.Data.GetBytes(valueOffset, 2)
			valueNode := format.Data.GetValue(valueBytes)

			//
			valueNode.Key = keyNode.ValueString
			nodes = append(nodes, valueNode)
		}

		// A -> Z
		sort.Slice(nodes, func(i, j int) bool {
			return nodes[i].Key < nodes[j].Key
		})

		//
		node.Type = Dictionary
		node.Nodes = nodes

		return node
	} else if (value & 0xA0) == 0xA0 {
		count = int(value - 0xA0)

		if count == 0xF {
			count = data.GetBytes(addressValue+1, 2).ToInt()
			count -= 0x1000
			moreOffset = 3
		}

		node.Type = Array

		for i := 0; i < count; i++ {
			offset := data[addressValue+moreOffset+i]
			keyOffset := format.RecordStartAddress.ToInt() + 2*int(offset)
			keyBytes := format.Data.GetBytes(keyOffset, 2)
			subNode := format.Data.GetValue(keyBytes)

			node.Nodes = append(node.Nodes, subNode)
		}

		return node
	} else if (value & 0x80) == 0x80 { // UID Value
		fmt.Printf("UID Value 待處理\n")
		return node
	} else if (value & 0x60) == 0x60 {
		count = int(value - 0x60)

		if count == 0xF {
			count = data.GetBytes(addressValue+1, 2).ToInt()
			count -= 0x1000
			moreOffset = 3
		}

		data := data.GetBytes(addressValue+moreOffset, 2*count)

		text := UInt16ArrayToString(data)

		node.Type = UnicodeString
		node.ValueString = text

		return node
	} else if value&0x50 == 0x50 { // String
		count = int(value - 0x50)

		if count == 0xF {
			count = data.GetBytes(addressValue+1, 2).ToInt()
			count -= 0x1000
			moreOffset = 3
		}

		data := data.GetBytes(addressValue+moreOffset, count)

		node.Type = String
		node.ValueString = string(data)

		return node
	} else if value&0x40 == 0x40 { //
		fmt.Printf("Data 待處理\n")
		node.Type = Dataa

		return node
	} else if value&0x30 == 0x30 { //
		fmt.Printf("Date 待處理\n")
		node.Type = Date
		return node
	} else if value&0x20 == 0x20 { // Real
		fmt.Printf("Real 待處理\n")
		node.Type = Real
		return node
	} else if value&0x10 == 0x10 { // Integer
		count = int(value - 0x10)
		valueInteger := format.Data.GetBytes(addressValue+count, 2).ToInt()
		valueInteger -= 0x1000

		node.Type = Integer
		node.ValueInteger = valueInteger

		return node
	} else {
		node.Type = NullOrBool

		if value == 0x09 {
			node.ValueBool = true
			return node
		} else if value == 0x08 {
			node.ValueBool = false
			return node
		}

		address.Print()
		fmt.Printf("無法辨識: 0x%.2X\n", int(value))
		fmt.Printf("NullOrBool 待處理\n")
		return node
	}
}

func (data Bytes) GetDataCount(addr Bytes) (int, Bytes) {
	index := addr.ToInt()
	value := int(data[index])

	// 超過 15
	if value&0xD0 == 0xD0 {
		value = data.GetBytes(index+1, 2).ToInt()
		value -= 0x1000

		return value, IntToByte(addr.ToInt() + 3)
	}

	return value, IntToByte(addr.ToInt() + 1)
}

func (data Bytes) Print() {
	for i := 0; i < len(data); i++ {
		fmt.Printf("0x%.2X ", data[i])
	}
	fmt.Println()
}

func (data Bytes) ToInt() int {
	value := 0
	for i := 0; i < len(data); i++ {
		value = (value << 8) + int(data[i])
	}
	return value
}

func IntToByte(value int) Bytes {
	count := 0
	if value <= 0xFF {
		count = 1
	} else if value <= 0xFFFF {
		count = 2
	} else if value <= 0xFFFFFF {
		count = 3
	} else {
		count = 4
	}
	bytes := make([]byte, count)

	for i := 0; i < count; i++ {
		bytes[i] = uint8(value >> (8 * (count - i - 1)))
	}

	return bytes
}
