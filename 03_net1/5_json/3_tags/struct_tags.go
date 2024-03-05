package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type UserType int

const (
	Unknown UserType = iota
	Regular
	Paying
)

func (a *UserType) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch strings.ToLower(s) {
	default:
		*a = Unknown
	case "regular":
		*a = Regular
	case "paying":
		*a = Paying
	}

	return nil
}

func (a UserType) MarshalJSON() ([]byte, error) {
	var s string
	switch a {
	default:
		s = "unknown"
	case Regular:
		s = "regular"
	case Paying:
		s = "zebra"
	}

	return json.Marshal(s)
}

type User struct {
	ID       int `json:"user_id,omitempty,string"`
	Username string
	Address  string   `json:",omitempty"`
	Comnpany string   `json:"-"`
	Type     UserType `json:","`
}

func main() {
	u := &User{
		ID:       1,
		Username: "rvasily",
		Address:  "",
		Comnpany: "Mail.Ru Group",
		Type:     Regular,
	}
	result, _ := json.Marshal(u)
	fmt.Printf("json string: %s\n", string(result))
}
