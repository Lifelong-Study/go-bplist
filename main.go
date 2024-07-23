package go_bplist

import (
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
)

var format Format

func main() {
	nodes, err := Parse("info.plist")

	//
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//
	Save(nodes, "sdf.plist")
}

func Parse(path string) ([]PlistNode, error) {
	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	bytes, err := io.ReadAll(fp)
	if err != nil {
		return nil, err
	}

	//
	if !isBinaryPlistFile(bytes) {
		return nil, errors.New("file format is incorrect")
	}

	format = Format{}
	format.Data = bytes
	format = format.fillFooter()

	//
	format.DataStartAddress.Print()
	format.RecordStartAddress.Print()
	fmt.Printf("Data Count: %d\n", format.DataCount)

	//
	nodes := make([]PlistNode, 0)
	for i := 0; i < format.DataCount; i++ {

		//
		keyNode := format.GetData(i)

		//
		valueNode := format.GetData(format.DataCount + i)
		valueNode.Key = keyNode.ValueString

		//
		nodes = append(nodes, valueNode)
	}

	// A -> Z
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Key < nodes[j].Key
	})

	return nodes, nil
}

// Save plist file
func Save(nodes []PlistNode, filename string) {
	fp, _ := os.Create(filename)
	defer fp.Close()

	//
	PrintXML(fp, 0, `<?xml version="1.0" encoding="UTF-8"?>`)
	PrintXML(fp, 0, `<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">`)
	PrintXML(fp, 0, `<plist version="1.0">`)
	PrintXML(fp, 0, `<dict>`)

	for _, node := range nodes {
		node.Print(fp, 1)
	}

	PrintXML(fp, 0, `</dict>`)
	PrintXML(fp, 0, `</plist>`)
}

func PrintXML(fp *os.File, level int, format string, a ...any) {
	spaces := ""
	for range level {
		spaces += "\t"
	}

	text := ""
	text = fmt.Sprintf(format, a...)
	text = fmt.Sprintf("%s%s\n", spaces, text)

	fp.WriteString(text)
}

func isBinaryPlistFile(bytes []byte) bool {
	if bytes[0] == 'b' && bytes[1] == 'p' &&
		bytes[2] == 'l' && bytes[3] == 'i' &&
		bytes[4] == 's' && bytes[5] == 't' {
		if bytes[6] == 0x30 && bytes[7] == 0x30 {
			return true
		}
	}
	return false
}

func (node PlistNode) Print(fp *os.File, level int) {

	// Key
	if len(node.Key) > 0 {
		PrintXML(fp, level, "<key>%s</key>", node.Key)
	}

	// Value
	if node.Type == Dictionary {
		PrintXML(fp, level, "<dict>")

		for _, subNode := range node.Nodes {
			subNode.Print(fp, level+1)
		}

		PrintXML(fp, level, "</dict>")
		return
	} else if node.Type == Array {
		PrintXML(fp, level, "<array>")

		for _, subNode := range node.Nodes {
			subNode.Print(fp, level+1)
		}

		PrintXML(fp, level, "</array>")
		return
	} else if node.Type == Integer {
		PrintXML(fp, level, "<integer>%d</integer>", node.ValueInteger)
	} else if node.Type == String || node.Type == UnicodeString {
		PrintXML(fp, level, "<string>%s</string>", node.ValueString)
	} else if node.Type == NullOrBool {
		if node.ValueBool {
			PrintXML(fp, level, "<true/>")
		} else {
			PrintXML(fp, level, "<false/>")
		}
	}
}
