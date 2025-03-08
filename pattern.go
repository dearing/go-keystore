package keystore

import (
	"fmt"
	"regexp"
)

// String returns a string representation of the collection
func (c *Collection[T]) String() string {
	return fmt.Sprintf("Collection: '%s', size=%d", c.Description, c.Len())
}

// MatchValues returns a slice of values that match the given pattern
// The pattern is a regular expression, of which the syntax is described here: https://golang.org/pkg/regexp/
//
// ex: values, err := db.MatchValues("c[0-9]")
func (c *Collection[T]) MatchValues(pattern string) ([]T, error) {

	var matches []T

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("error compiling regex: %w", err)
	}

	for key, record := range c.Keys {
		if regex.MatchString(key) {
			matches = append(matches, record.Value)
		}
	}

	return matches, nil

}

// MatchKeys returns a slice of keys that match the given pattern
// The pattern is a regular expression, of which the syntax is described here: https://golang.org/pkg/regexp/
//
// ex: keys, err := db.MatchKeys("c[0-9]")
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

// Prefix returns a slice of keys that match the given prefix
// Supported wildcards are '*' to match zero or more characters and '?' to match exactly one character.
//
// ex: users, err := db.Prefix("user:")
func (c *Collection[T]) Prefix(prefix string) ([]string, error) {

	var matches []string

	for key := range c.Keys {
		if matchWildcard(key, prefix) {
			matches = append(matches, key)
		}
	}

	return matches, nil

}

// PrefixChan returns a channel of keys that match the given prefix
// Supported wildcards are '*' to match zero or more characters and '?' to match exactly one character.
//
// ex: for key := range db.PrefixChan("user:") {}
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

// matchWildcard returns true if the key matches the pattern
// Supported wildcards are '*' to match zero or more characters and '?' to match exactly one character.
// TODO: feels clunky, could be improved
// ex: matchWildcard("c1", "c*")
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
