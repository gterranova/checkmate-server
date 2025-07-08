package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type Router struct {
	Handlers map[string]func(ctx *fiber.Ctx) (err error)
}

var (
	ErrorNotFound             = fmt.Errorf("not found")
	ErrorConsecutiveWildcards = fmt.Errorf("consecutive wildcards")
	ErrorNestedParam          = fmt.Errorf("nested named param")
	ErrorUnopenedParam        = fmt.Errorf("not opened named param")
	ErrorUnclosedParam        = fmt.Errorf("unclosed named param")
	ErrorUnmatched            = fmt.Errorf("pattern does not match")
)

func extractWildcard(str, prefix, suffix string) (string, error) {
	//fmt.Printf("\"%s\" \"%s\" \"%s\" -> ", str, prefix, suffix)
	if strings.HasPrefix(str, prefix) && strings.HasSuffix(str, suffix) {
		// If the element contains a wildcard, extract the text between the wildcards.
		result := strings.TrimPrefix(strings.TrimSuffix(str, suffix), prefix)
		return result, nil
	}
	return "", ErrorNotFound
}

func splitPattern(pattern string) ([]string, error) {
	segments := make([]string, 0)
	var str string
	var opened bool
	for _, c := range pattern {
		if string(c) == "*" {
			if len(str) == 0 {
				return segments, ErrorConsecutiveWildcards
			}
			segments = append(segments, str)
			segments = append(segments, "*")
			str = ""
		} else if string(c) == "{" {
			if opened {
				return segments, ErrorNestedParam
			}
			if len(str) == 0 {
				return segments, ErrorConsecutiveWildcards
			}
			segments = append(segments, str)
			str = ""
			opened = true
		} else if string(c) == "}" {
			if !opened {
				return segments, ErrorUnopenedParam
			}
			segments = append(segments, "{"+str+"}")
			str = ""
			opened = false
		} else {
			str += string(c)
		}
	}
	if opened {
		return segments, ErrorUnclosedParam
	}
	if len(str) > 0 {
		segments = append(segments, str)
	}
	return segments, nil
}

func isWildcard(s string) bool {
	return s == "*" || (strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}"))
}

func parseWildcard(str, pattern string) (map[string]string, error) {
	var suffix, currentString string

	// Create a map to store the extracted text.
	result := make(map[string]string)

	// Split the string into a slice of strings using the wildcard as a delimiter.
	s, err := splitPattern(pattern)
	if err != nil {
		return result, err
	}

	// Loop through each slice element and check if it contains a wildcard.
	//fmt.Println("parts", s)
	for i, elem := range s {
		if elem == "" {
			continue
		}
		if isWildcard(elem) {
			chunk := str
			if len(s) <= i+1 {
				suffix = ""
			} else {
				suffix = s[i+1]
				if !strings.Contains(chunk, suffix) {
					return result, ErrorUnmatched
				}
				chunk = chunk[:strings.Index(chunk, suffix)+len(suffix)]
			}
			//fmt.Printf("Elem: %v, Current: %v, Chunk is \"%s\"\nsuffix=%v\n", elem, currentString, chunk, suffix)
			extractedText, _ := extractWildcard(chunk, currentString, suffix)

			var key string
			if elem == "*" {
				key = fmt.Sprintf("*%d", i)
			} else {
				key = elem[1 : len(elem)-1]
			}
			result[key] = extractedText
			currentString += extractedText
		} else {
			currentString += elem
		}
	}

	return result, nil
}

func (r *Router) Handle(ctx *fiber.Ctx) (err error) {

	path := ctx.Path()
	if path == "/" {
		// Handle root path
	} else {
		path = strings.ReplaceAll(filepath.Join("", path), "\\", "/")
		for pattern, handle := range r.Handlers {
			if params, err := parseWildcard(path, pattern); err != nil {
				ctx.Locals("params", params)
				return handle(ctx)
			}
		}
	}
	return ctx.Next()
}
