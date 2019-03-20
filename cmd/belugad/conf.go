package main

import (
//	"gopkg.in/yaml.v2"
)

type Config struct {
	Keys []APIKey
}

type APIKey struct {
	Key                 string
	AllowedDomainsRegex string
}
