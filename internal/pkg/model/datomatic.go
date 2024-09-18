package model

type DatoMaticDataFile struct {
	Header Header `xml:"header"`
	Games  []Game `xml:"game"`
}

type Header struct {
	Id          string     `xml:"id"`
	Name        string     `xml:"name"`
	Description string     `xml:"description"`
	Version     string     `xml:"version"`
	Author      string     `xml:"author"`
	Homepage    string     `xml:"Homepage"`
	URL         string     `xml:"url"`
	Clrmamepro  Clrmamepro `xml:"clrmamepro"`
}

type Clrmamepro struct {
	forcenodump string `xml:"forcenodump,attr"`
}

type Game struct {
	Id          string `xml:"id,attr"`
	Name        string `xml:"name,attr"`
	CloneOfId   string `xml:"cloneofid,attr"`
	Description string `xml:"description"`
	Rom         []Rom  `xml:"rom"`
}

type Rom struct {
	Name   string `xml:"name,attr"`
	Size   string `xml:"size,attr"`
	CRC    string `xml:"crc,attr"`
	Md5    string `xml:"md5,attr"`
	Sha1   string `xml:"sha1,attr"`
	Sha256 string `xml:"sha256,attr"`
}
