package model

// https://github.com/mamedev/mame/tree/master/hash
// <!DOCTYPE softwarelist SYSTEM "softwarelist.dtd">
// https://github.com/mamedev/mame/blob/master/hash/softwarelist.dtd
type SoftwareList struct {
	Name        string     `xml:"name,attr"`
	Description string     `xml:"description"`
	Software    []Software `xml:"software"`
}

type Software struct {
	Name        string     `xml:"name,attr"`
	Description string     `xml:"description"`
	Year        string     `xml:"year"`
	Publisher   string     `xml:"publisher"`
	Sharedfeat  Sharedfeat `xml:"sharedfeat"`
	Part        Part       `xml:"part"`
	Info        []Info     `xml:"info"`
}

type Info struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type Sharedfeat struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type Part struct {
	Name      string    `xml:"name,attr"`
	Interface string    `xml:"interface,attr"`
	Feature   []Feature `xml:"feature"`
	DataArea  DataArea  `xml:"dataarea"`
}

type Feature struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type DataArea struct {
	Name string    `xml:"name,attr"`
	Size string    `xml:"size,attr"`
	Rom  []MAMERom `xml:"rom"`
	Disk []Disk    `xml:"disk"`
}

type MAMERom struct {
	Name     string `xml:"name,attr"`
	Size     string `xml:"size,attr"`
	CRC      string `xml:"crc,attr"`
	SHA1     string `xml:"sha1,attr"`
	Offset   string `xml:"offset,attr"`
	LoadFlag string `xml:"loadflag,attr"`
	Status   string `xml:"status"`
}

type Disk struct {
	Name   string `xml:"name,attr"`
	SHA1   string `xml:"sha1,attr"`
	Status string `xml:"status"`
}
