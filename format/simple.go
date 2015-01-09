package format

import (
	"os"
	"encoding/xml"
	"io"
	"bytes"
	"strings"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"bufio"
)

type Simple struct {
	Src string
}

func (f Simple) Ext() string {
	return "xml"
}
func (f Simple) FromJson(to string, jsonMap map[string]string) error {
	file, err := os.Create(to)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("Error creating file %q due to: %v\n", to, err)
	}

	bufWriter := bufio.NewWriter(file)
	defer bufWriter.Flush()

	bufWriter.WriteString("<strings>\n")
	for key,translation := range jsonMap {
		bufWriter.WriteString("    <")
		bufWriter.WriteString(key)
		bufWriter.WriteString(">")
		bufWriter.WriteString(translation)
		bufWriter.WriteString("</")
		bufWriter.WriteString(key)
		bufWriter.WriteString(">\n")
	}

	bufWriter.WriteString("</strings>\n")
	return nil
}
func (f Simple) ToJson(to string) error {

	sourceFile, err := os.Open(f.Src)
	defer sourceFile.Close()
	if err != nil {
		return fmt.Errorf("Unable to open file %q due to %v\n", f.Src, err)
	}

	parser := xml.NewDecoder(sourceFile)

	jsonMap := map[string]string{}
	var current string
	var text bytes.Buffer
	for {
		token, err := parser.Token()

		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("Error processing XML: %v", err)
		}

		switch t := token.(type) {
		case xml.StartElement:
			current = xml.StartElement(t).Name.Local
			text = bytes.Buffer{}
		case xml.EndElement:
			jsonMap[current] = text.String()
		case xml.CharData:
			text.WriteString(strings.TrimSpace(string(xml.CharData(t))))
		}
	}


	jsonBytes, err := json.MarshalIndent(jsonMap, "", "  ")

	if err != nil {
		return fmt.Errorf("Failed to convert string map to json due to: %v", err)
	}

	err = ioutil.WriteFile(to, jsonBytes, os.FileMode(0))
	if err != nil {
		return fmt.Errorf("Failed to write json file %q due to: %v", to, err)
	}

	return nil
}
