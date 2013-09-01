package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	cWhite = "\r\n\t "
	cSep   = "{|}"
)

type node struct {
	Name     string  `json:"name"`
	Children []*node `json:"children"`
	Filename string  `json:"filename"`
	ID       string  `json:"id"`
}

func parseStructure(reader io.Reader) (*node, error) {
	lines := make([]string, 0)
	r := bufio.NewReader(reader)
	var err error
	line := ""
	for err == nil {
		line, err = r.ReadString('\n')
		if err != nil && err != io.EOF {
			return nil, err
		}
		if len(strings.Trim(line, cWhite)) > 0 {
			lines = append(lines, line)
		}
	}
	if len(lines) < 2 {
		return nil, errors.New("Not enough lines in structure file.")
	}
	parsed, err := parseList(lines)
	if parsed != nil && len(parsed) > 0 {
		return parsed[0], err
	} else {
		return nil, err
	}
}

func parseList(lines []string) ([]*node, error) {
	var err error

	ret := make([]*node, 0)
	if len(lines) == 0 {
		return ret, nil
	}

	baseIndent := getIndent(lines[0])
	for i := 0; i < len(lines); {
		if baseIndent > getIndent(lines[i]) {
			return nil, errors.New("Indents are inconsistent.")
		}
		if baseIndent == getIndent(lines[i]) {
			first, second := getTokens(lines[i])
			n := &node{ID: getID()}
			if first != "" {
				n.Name = first
			}
			if second != "" {
				n.Filename = second
			}
			ret = append(ret, n)
			i++
		} else {
			childrenIndent := getIndent(lines[i])
			childrenStart := i
			childrenEnd := i + 1
			for ; childrenEnd < len(lines) && getIndent(lines[childrenEnd]) >= childrenIndent; childrenEnd++ {
			}
			ret[len(ret)-1].Children, err = parseList(lines[childrenStart:childrenEnd])
			if err != nil {
				return nil, err
			}
			i = childrenEnd
		}
	}

	return ret, nil
}

func getIndent(str string) int {
	return len(str) - len(strings.TrimLeft(str, cWhite))
}

func getTokens(str string) (first, second string) {
	tokens := strings.Split(str, cSep)
	if len(tokens) > 1 {
		return strings.Trim(tokens[0], cWhite), strings.Trim(tokens[1], cWhite)
	}
	return strings.Trim(tokens[0], cWhite), ""
}

var counter = 0

func getID() string {
	counter++
	return fmt.Sprintf("id%x", counter)
}
