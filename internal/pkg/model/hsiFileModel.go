package model

//https://github.com/Klozz/xMAME/tree/master/hash

type HSIFile struct {
	HashFile HashFile `xml:"hashfile"`
}

type HashFile struct {
	Hash []HSIHash `xml:"hash"`
}

type HSIHash struct {
	CRC32        string `xml:"crc32,attr"`
	SHA1         string `xml:"sha1,attr"`
	MD5          string `xml:md5,attr"`
	Name         string `xml:"name,attr"`
	Year         string `xml:year"`
	Manufacturer string `xml:manufacturer"`
	Status       string `xml:"status"`
	ExtraInfo    string `xml:"extrainfo"`
}
