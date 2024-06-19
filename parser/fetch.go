package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const nl = "\r\n"

// * \d+ FETCH \(BODY[] \)`

var regexFetchResponse = regexp.MustCompile(`^* \d+ FETCH \(`)

func ParseFetchResponse(resp string) (map[string]string, error) {
	lineIdx := -1
	lines := strings.Split(resp, nl)
	for i, l := range lines {
		if regexFetchResponse.MatchString(l) {
			lineIdx = i
		}
	}

	if lineIdx == -1 {
		return nil, fmt.Errorf("fetch response not found. response: %q", resp)
	}

	p := New(strings.Join(lines[lineIdx:], nl))
	s := p.ConsumeUntil(' ')
	if s != "*" {
		return nil, fmt.Errorf("could not find '*'. found: %q. remaining: %q", s, p.Remaining())
	}
	s = p.ConsumeUntil(' ')
	if _, err := strconv.ParseUint(s, 10, 32); err != nil {
		return nil, fmt.Errorf("not valid sequence number: %q", s)
	}
	consumed := p.Consume("FETCH (")
	if !consumed {
		return nil, fmt.Errorf("could not find '('. remaining: %q", p.Remaining())
	}

	m := make(map[string]string)
	for {
		consumed = p.Consume(nl)
		if consumed {
			break
		}

		k := p.ConsumeUntil(' ')
		var content string
		consumed = p.Consume("{")
		if consumed {
			sizeStr := p.ConsumeUntil('}')
			size, err := strconv.Atoi(sizeStr)
			if err != nil {
				return nil, fmt.Errorf("expected integer. parsed size: %q. "+
					"remaining: %q. error: %w", sizeStr, p.Remaining(), err)
			}
			content = p.ConsumeN(size)
			consumed = p.Consume(nl)
			if !consumed {
				return nil, fmt.Errorf("expecting \\r\\n. remaining=%q", p.Remaining())

			}
			consumed = p.Consume(" ") || p.Consume(")")
			if !consumed {
				return nil, fmt.Errorf("could not find ' ' or ')'. "+
					"remaining=%q", p.Remaining())
			}
		} else {
			content = p.ConsumeUntil(' ', ')')
		}
		m[k] = content
	}
	return m, nil
}
