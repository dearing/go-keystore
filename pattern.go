package keystore

import (
	"fmt"
	"regexp"
)

func (c *Collection[T]) MatchKeys(pattern string) ([]string, error) {

	var matches []string

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("error compiling regex: %w", err)
	}

	for key := range c.Keys {
		if regex.MatchString(key) {
			matches = append(matches, key)
		}
	}

	return matches, nil

}

func (c *Collection[T]) Prefix(prefix string) ([]string, error) {

	var matches []string

	for key := range c.Keys {
		if matchWildcard(key, prefix) {
			matches = append(matches, key)
		}
	}

	return matches, nil

}

func (c *Collection[T]) PrefixChan(prefix string) chan string {
	out := make(chan string, 2)
	go func() {
		for key := range c.Keys {
			if matchWildcard(key, prefix) {
				out <- key
			}
		}
		close(out)
	}()

	return out
}

func matchWildcard(key, pattern string) bool {
	keyIndex := 0
	patternIndex := 0

	for keyIndex < len(key) || patternIndex < len(pattern) {
		if patternIndex < len(pattern) {
			if pattern[patternIndex] == '*' {
				patternIndex++
				for keyIndex < len(key) {
					if matchWildcard(key[keyIndex:], pattern[patternIndex:]) {
						return true
					}
					keyIndex++
				}
				return matchWildcard(key[keyIndex:], pattern[patternIndex:])
			} else if pattern[patternIndex] == '?' {
				if keyIndex >= len(key) {
					return false
				}
				keyIndex++
				patternIndex++
			} else if keyIndex < len(key) && key[keyIndex] == pattern[patternIndex] {
				keyIndex++
				patternIndex++
			} else {
				return false
			}
		} else if keyIndex < len(key) {
			return false
		} else {
			return true
		}
	}
	return true
}
