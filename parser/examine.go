package parser

import (
	"errors"
	"regexp"
	"strconv"
)

var regexpUIDValidity = regexp.MustCompile(`(?m)^\* OK \[UIDVALIDITY (\d+)\]`)

func ParseExamineResponse(resp string) (uint32, error) {
	submatches := regexpUIDValidity.FindStringSubmatch(resp)
	if len(submatches) != 2 {
		return 0, errors.New("could not find UIDVALIDITY")
	}
	uid, err := strconv.ParseUint(submatches[1], 10, 32)
	return uint32(uid), err
}
