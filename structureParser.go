package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type node struct {
	Name     string  `json:"name"`
	Children []*node `json:"children"`
	Filename *string `json:"filename"`
	ID       *string `json:"id"`
}

func parseStructure(reader io.Reader) (interface{}, error) {
	lines := make([]string, 0)
	r := bufio.NewReader(reader)
	var err error
	line := ""
	for err == nil {
		line, err = r.ReadString('\n')
		if err != nil && err != io.EOF {
			return nil, err
		}
		lines = append(lines, line)
	}
	parsed, err := parseList(lines)
	return parsed, err
}

func parseList(lines []string) ([]*node, error) {
	var err error

	ret := make([]*node, 0)
	if len(lines) == 0 {
		return ret, nil
	}

	if _, second := getTokens(lines[0]); second == "" { // is a category
		baseIndent := getIndent(lines[0])
		for i := 0; i < len(lines); {
			ret = append(ret, &node{Name: strings.Trim(lines[i], "\r\n ")})
			j := i + 1
			for ; j < len(lines) && getIndent(lines[j]) > baseIndent; j++ {
			}
			ret[len(ret)-1].Children, err = parseList(lines[i+1 : j])
			if err != nil {
				return nil, err
			}
			i = j
		}
	} else {
		for i := 0; i < len(lines); i++ {
			first, second := getTokens(lines[i])
			id := encodeName(first)
			ret = append(ret, &node{Name: first, ID: &id, Filename: &second})
		}
	}

	return ret, nil
}

func getIndent(str string) int {
	return len(str) - len(strings.TrimLeft(str, "\r\n "))
}

func getTokens(str string) (first, second string) {
	tokens := strings.Split(str, "{|}")
	if len(tokens) > 1 {
		return strings.Trim(tokens[0], "\r\n "), strings.Trim(tokens[1], "\r\n ")
	}
	return strings.Trim(tokens[0], "\r\n "), ""
}

func encodeName(name string) string {
	return fmt.Sprintf("id%x", name)
}
