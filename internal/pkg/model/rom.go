package model

type Release struct {
	Name   string `json:"name"`
	Region string `json:"region"`
}

type ROM struct {
	Name string `json:"name"`
	CRC  string `json:"crc"`
	MD5  string `json:"md5"`
	SHA1 string `json:"sha1"`
}

type Game struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Release     []Release `json:"release"`
	ROM         ROM       `json:"rom"`
}
