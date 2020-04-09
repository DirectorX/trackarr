package config

import (
	"regexp"

	"github.com/antchfx/xmlquery"
	irc "github.com/thoj/go-ircevent"
)

type TrackerInstance struct {
	Name   string
	Config *TrackerConfig
	Info   *TrackerInfo
	IRC    *irc.Connection
}

type TrackerConfig struct {
	Enabled   bool
	ForceHTTP bool
	Bencode   *TrackerBencodeConfig
	Settings  map[string]string
	IRC       *TrackerIrcConfig
}

type TrackerBencodeConfig struct {
	Name bool
	Size bool
}

type TrackerIrcConfig struct {
	Nickname   string
	Channels   []string
	Announcers []string
	Commands   []string
	Host       *string
	Port       *string
	Sasl       TrackerIrcSaslConfig
	Verbose    bool
}

type TrackerIrcSaslConfig struct {
	User string
	Pass string
}

type TrackerInfo struct {
	Name       string
	ShortName  *string
	LongName   string
	Settings   []string
	Servers    []string
	Channels   []string
	Announcers []string

	IgnoreLines       []TrackerIgnore
	LinePatterns      []TrackerPattern
	MultiLinePatterns []TrackerPattern

	LineMatchedRules *xmlquery.Node
}

type TrackerIgnore struct {
	Rxp      *regexp.Regexp
	Expected bool
}

type TrackerPattern struct {
	PatternType MessagePatternType
	Rxp         *regexp.Regexp
	Vars        []string
	Optional    bool
}

type MessagePatternType int

/* Enum */

const (
	LinePattern MessagePatternType = iota + 1
	MultiLinePattern
)
