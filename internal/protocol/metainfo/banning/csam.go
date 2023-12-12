package banning

import (
	"errors"

	"github.com/bitmagnet-io/bitmagnet/internal/protocol/metainfo"
	"github.com/bitmagnet-io/bitmagnet/internal/regex"
)

var csamWords = []string{
	"pedo",
	"pedofile",
	"pedofilia",
	"preteen",
	"pthc",
	"ptsc",
	"lsbar",
	"lsm",
	"underage",
	"1yo",
	"2yo",
	"3yo",
	"4yo",
	"5yo",
	"6yo",
	"7yo",
	"8yo",
	"9yo",
	"10yo",
	"11yo",
	"12yo",
	"13yo",
	"14yo",
	"15yo",
	"16yo",
	"17yo",
	"hebefilia",
	"opva",
}

var csamRegex = regex.NewRegexFromNames(csamWords...)

type csamChecker struct{}

func (c csamChecker) Check(info metainfo.Info) error {
	if csamRegex.MatchString(info.BestName()) {
		return errors.New("torrent appears to contain CSAM")
	}
	return nil
}
