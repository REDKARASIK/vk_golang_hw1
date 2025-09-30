package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type Config struct {
	Comments  bool `json:"comments"`
	Limit     int
	Type      Enum
	Servers   []string
	CustomCfg struct {
		Flags int
	} `json:"tmp"`
}

const (
	Unknown Enum = iota
	Prod
	Staging
)

type Enum int

var _ json.Unmarshaler = new(Enum)

func (a *Enum) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch strings.ToLower(s) {
	default:
		*a = Unknown
	case "prod":
		*a = Prod
	case "staging":
		*a = Staging
	}

	return nil
}

var (
	config = &Config{}
)

func main() {
	data, err := os.ReadFile("./config.json")
	if err != nil {
		log.Fatalln("cant read config file:", err)
	}

	err = json.Unmarshal(data, config)
	if err != nil {
		log.Fatalln("cant parse config:", err)
	}

	if config.Comments {
		fmt.Println("Comments per page", config.Limit)
		fmt.Println("Comments services", config.Servers)
	} else {
		fmt.Println("Comments disabled")
	}
	fmt.Printf("Cfg: %+v", config)
}
