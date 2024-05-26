package metadata

import (
	"archive/zip"
	"encoding/xml"
	"strings"
)

type OfficeCoreProperty struct {
	XMLName        xml.Name `xml:"coreProperties"`
	Creator        string   `xml:"creator"`
	LastModifiedBy string   `xml:"lastModifiedBy"`
}

type OfficeAppProperty struct {
	XMLName     xml.Name `xml:"Properties"`
	Application string   `xml:"Application"`
	Company     string   `xml:"Company"`
	Version     string   `xml:"AppVersion"`
}

var OfficeVersion = map[string]string{
	"16": "2016",
	"15": "2013",
	"14": "2010",
	"12": "2007",
	"11": "2003",
}

func (o *OfficeAppProperty) GetMajorVersion() string {
	tokens := strings.Split(o.Version, ".")
	if len(tokens) < 2 {
		return "Unknown"
	}

	v, ok := OfficeVersion[tokens[0]]
	if !ok {
		return "Unknown"
	}

	return v
}

func process(r *zip.File, prop interface{}) error {
	file, err := r.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	if err := xml.NewDecoder(file).Decode(&prop); err != nil {
		return err
	}

	return nil
}

func NewProperties(z *zip.Reader) (*OfficeAppProperty, *OfficeCoreProperty, error) {
	var appProp OfficeAppProperty
	var coreProp OfficeCoreProperty

	for _, f := range z.File {
		switch f.Name {
		case "docProps.core.xml":
			if err := process(f, &coreProp); err != nil {
				return nil, nil, err
			}
		case "docProps.app.xml":
			if err := process(f, &appProp); err != nil {
				return nil, nil, err
			}
		default:
			continue
		}
	}

	return &appProp, &coreProp, nil
}
