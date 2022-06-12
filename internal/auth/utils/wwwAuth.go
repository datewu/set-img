package utils

import (
	"strings"
)

// ConsumeParams parses the authorization header and returns the parameters.
func ConsumeParams(v string) (params map[string]string, newv string) {
	for {
		var name, value string
		name, value, v = consumeParam(v)
		if name == "" {
			break
		}
		// Use only the first occurrence of each param name.
		// This is required in some other places that don't use consumeParams,
		// but it seems like reasonable behavior in general.
		if _, seen := params[name]; seen {
			continue
		}
		if params == nil {
			params = make(map[string]string)
		}
		params[name] = value
	}
	return params, v
}

func consumeParam(v string) (name, value, newv string) {
	v = skipWSAnd(v, ';')
	name, v = consumeItem(v)
	if name == "" {
		return "", "", v
	}
	name = strings.ToLower(name)
	v = skipWS(v)
	if peek(v) == '=' {
		v = skipWS(v[1:])
		value, v = consumeItemOrQuoted(v)
	}
	return name, value, v
}

func skipWSAnd(v string, and byte) string {
	for v != "" && (v[0] == ' ' || v[0] == '\t' || v[0] == and) {
		v = v[1:]
	}
	return v
}

func skipWS(v string) string {
	for v != "" && (v[0] == ' ' || v[0] == '\t') {
		v = v[1:]
	}
	return v
}

// consumeItem returns the item from the beginning of v, and the rest of v.
// An item is a run of text up to whitespace, comma, semicolon, or equal sign.
// Callers should check that the item is non-empty if they need to make progress.
func consumeItem(v string) (item, newv string) {
	for i := 0; i < len(v); i++ {
		switch v[i] {
		case ' ', '\t', ',', ';', '=':
			return v[:i], v[i:]
		}
	}
	return v, ""
}

func peek(v string) byte {
	if v == "" {
		return 0
	}
	return v[0]
}

func consumeItemOrQuoted(v string) (text, newv string) {
	if peek(v) == '"' {
		text, newv = consumeQuoted(v)
		return
	}
	return consumeItem(v)
}

func consumeQuoted(v string) (text, newv string) {
	return consumeDelimited(v, '"', '"')
}

func consumeDelimited(v string, opener, closer byte) (text, newv string) {
	if peek(v) != opener {
		return "", v
	}
	v = v[1:]

	// In the common case, when there are no quoted pairs,
	// we can simply slice the string between the outermost delimiters.
	nesting := 1
	i := 0
	for ; i < len(v); i++ {
		switch v[i] {
		case closer:
			nesting--
			if nesting == 0 {
				return v[:i], v[i+1:]
			}
		case opener:
			nesting++
		case '\\': // start of a quoted pair
			goto buffered
		}
	}
	// We've reached the end of v, but nesting is still > 0.
	// This is an unterminated string.
	return v, ""

buffered:
	// Once we have encountered a quoted pair, we have to unquote into a buffer.
	b := &strings.Builder{}
	b.WriteString(v[:i])
	quoted := false
	for ; i < len(v); i++ {
		switch {
		case quoted:
			b.WriteByte(v[i])
			quoted = false
		case v[i] == closer:
			nesting--
			if nesting == 0 {
				return b.String(), v[i+1:]
			}
			b.WriteByte(v[i])
		case v[i] == opener:
			nesting++
			b.WriteByte(v[i])
		case v[i] == '\\':
			quoted = true
		default:
			b.WriteByte(v[i])
		}
	}
	// We've reached the end of v, but nesting is still > 0.
	// This is an unterminated string.
	return b.String(), ""
}
