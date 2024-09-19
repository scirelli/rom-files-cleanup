package model

// <datafile xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="https://datomatic.no-intro.org/stuff https://datomatic.no-intro.org/stuff/schema_nointro_datfile_v3.xsd">
// <!DOCTYPE datafile PUBLIC "-//Logiqx//DTD ROM Management Datafile//EN" "http://www.logiqx.com/Dats/datafile.dtd">
type DATFile struct {
	Datafile Datafile `xml:"datafile"`
}

type Datafile struct {
	Header            Header `xml:"header"`
	Games             []Game `xml:"game"`
	xmlns             string `xml:"xmlns:xsi,attr"`
	XSISchemaLocation string `xml:"xsi:schemaLocation,attr"`
}

type DATHeader struct {
	Id          string     `xml:"id"`
	Name        string     `xml:"name"`
	Description string     `xml:"description"`
	Version     string     `xml:"version"`
	Author      string     `xml:"author"`
	Homepage    string     `xml:"Homepage"`
	URL         string     `xml:"url"`
	Subset      string     `xml:"subset"`
	Clrmamepro  clrmamepro `xml:"clrmamepro"`
}

type clrmamepro struct {
	forcenodump string `xml:"forcenodump,attr"`
}

type DATGame struct {
	Id          string `xml:"id,attr"`
	Name        string `xml:"name,attr"`
	CloneOfId   string `xml:"cloneofid,attr"`
	Description string `xml:"description"`
	Rom         []Rom  `xml:"rom"`
}

type DATRom struct {
	Name   string `xml:"name,attr"`
	Size   string `xml:"size,attr"`
	CRC    string `xml:"crc,attr"`
	Md5    string `xml:"md5,attr"`
	Sha1   string `xml:"sha1,attr"`
	Sha256 string `xml:"sha256,attr"`
	Status string `xml:"status,attr"`
	Serial string `xml:"serial,attr"`
	Header string `xml:"header,attr"`
}
