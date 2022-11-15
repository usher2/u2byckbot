package main

const (
	TBLOCK_UNKNOWN = iota
	TBLOCK_URL
	TBLOCK_HTTPS
	TBLOCK_DOMAIN
	TBLOCK_IP
)

type TMinContent struct {
	Id                 int32
	BlockType          int32 // for protobuf
	RegistryUpdateTime int64
	Url                []TUrl
	Ip4                []TIp4
	Domain             []TDomain
	Pack               []byte
	U2Hash             uint64
}

type TContent struct {
	Id           int32     `json:"id"`
	Decision     string    `json:"d"`
	DecisionInfo string    `json:"di"`
	IncludeTime  int64     `json:"it"`
	Url          []TUrl    `json:"url,omitempty"`
	Ip4          []TIp4    `json:"ip4,omitempty"`
	Domain       []TDomain `json:"dm,omitempty"`
	HttpsBlock   int       `json:"hb"`
	U2Hash       uint64    `json:"u2h"`
}

type TDomain struct {
	Domain string `json:"dm"`
}

type TUrl struct {
	Url string `json:"u"`
}

type TIp4 struct {
	Ip4 uint32 `json:"ip4"`
}
