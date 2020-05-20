package commandproxy

import (
	"bytes"
	"fmt"
	"strings"
)

var (
	errQuoteNotClose = fmt.Errorf("quote is not closed")
	errEmptyCmmand   = fmt.Errorf("cmd is empty")
)

func SplitCommand(cmd string) ([]string, error) {
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return nil, errEmptyCmmand
	}
	var ss []string

	var buf bytes.Buffer
	for len(cmd) != 0 {
		i := strings.IndexAny(cmd, "'\"\\ \n\t")
		if i == -1 {
			buf.WriteString(cmd)
			break
		}
		buf.WriteString(cmd[:i])
		c := cmd[i]
		cmd = cmd[i+1:]
		switch c {
		case '"':
			i := strings.IndexByte(cmd, '"')
			if i == -1 {
				return nil, errQuoteNotClose
			}
			buf.WriteString(cmd[:i])
			cmd = cmd[i+1:]
		case '\'':
			i := strings.IndexByte(cmd, '\'')
			if i == -1 {
				return nil, errQuoteNotClose
			}
			buf.WriteString(cmd[:i])
			cmd = cmd[i+1:]
		case '\\':
			if len(cmd) == 0 {
				break
			}
			if cmd[0] != '\n' {
				buf.WriteByte(cmd[0])
				cmd = cmd[:1]
				continue
			}
			fallthrough
		case '\t', ' ':
			if buf.Len() != 0 {
				ss = append(ss, buf.String())
				buf.Reset()
			}
		}
	}

	if buf.Len() != 0 {
		if len(ss) == 0 {
			return []string{cmd}, nil
		}
		ss = append(ss, buf.String())
	}

	return ss, nil
}

func ReplaceEscape(s string, re map[byte]string) string {
	re['%'] = "%"
	var buf bytes.Buffer
	for len(s) != 0 {
		i := strings.IndexByte(s, '%')
		if i == -1 || i >= len(s)-1 {
			if buf.Len() == 0 {
				return s
			}
			buf.WriteString(s)
			break
		}
		buf.WriteString(s[:i])
		s = s[i:]
		rep, ok := re[s[1]]
		if ok {
			buf.WriteString(rep)
		} else {
			buf.WriteString(s[:2])
		}
		s = s[2:]
	}
	return buf.String()
}
