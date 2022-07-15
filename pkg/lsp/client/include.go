package client

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/lazzer64/ttsls/pkg/uri"
)

var (
	includePattern       = regexp.MustCompile(`^#include\s*(.*?)\s*$`)
	includeMarkerPattern = regexp.MustCompile(`^----#include\s*(.*?)\s*$`)
)

func libURI(source SourceFile, name string) (uri.URI, error) {
	if strings.HasPrefix(name, "<") && strings.HasSuffix(name, ">") {
		name = strings.TrimSuffix(strings.TrimPrefix(name, "<"), ">")
	}
	if !strings.HasSuffix(name, ".ttslua") {
		name = fmt.Sprintf("%s.ttslua", name)
	}
	if strings.HasPrefix(name, "!") {
		// TODO
	} else {
		name = filepath.Join(filepath.Dir(source.Path()), name)
	}
	return uri.Parse(name)
}

func Expand(client Client, source SourceFile) (string, error) {
	s := ""
	for _, line := range strings.SplitAfter(source.Content(), "\n") {
		if matches := includePattern.FindStringSubmatch(line); len(matches) == 2 {
			libName := matches[1]
			uri, err := libURI(source, libName)
			if err != nil {
				return "", err
			}
			lib, err := client.Files.Get(uri)
			if err != nil {
				return "", err
			}
			content, err := Expand(client, lib)
			if err != nil {
				return "", err
			}
			lineEnd := "\n"
			if strings.HasSuffix(line, "\r\n") {
				lineEnd = "\r\n"
			}
			s += fmt.Sprintf("----#include %s%s", libName, lineEnd)
			if strings.HasPrefix(libName, "<") && strings.HasSuffix(libName, ">") {
				s += fmt.Sprintf("do%s%send%s", lineEnd, content, lineEnd)
			} else {
				s += content
			}
			s += fmt.Sprintf("----#include %s%s", libName, lineEnd)
		} else {
			s += line
		}
	}
	return s, nil
}

func Collapse(client Client, source SourceFile) (string, error) {
	s := ""
	stack := []string{}
	libcontent := ""

	for _, line := range strings.SplitAfter(source.Content(), "\n") {
		if matches := includeMarkerPattern.FindStringSubmatch(line); len(matches) == 2 {
			libName := matches[1]
			lineEnd := "\n"
			if strings.HasSuffix(line, "\r\n") {
				lineEnd = "\r\n"
			}

			// push if empty stack, this will be our #include statement
			if len(stack) == 0 {
				stack = append(stack, libName)
				s += fmt.Sprintf("#include %s%s", libName, lineEnd)
				libcontent = ""
				continue
			}

			// if a new include does not match the most recent scope open a new one
			if stack[len(stack)-1] != libName {
				stack = append(stack, libName)
				libcontent += line
				continue
			}

			// since the scope does match the most recent, close that scope
			stack = stack[:len(stack)-1]

			// if there are no remaining scopes we must have finished the #include
			if len(stack) == 0 {
				uri, err := libURI(source, libName)
				if err != nil {
					return "", err
				}
				if strings.HasPrefix(libName, "<") && strings.HasSuffix(libName, ">") {
					libcontent = strings.TrimPrefix(libcontent, fmt.Sprintf("do%s", lineEnd))
					libcontent = strings.TrimSuffix(libcontent, fmt.Sprintf("end%s", lineEnd))
				}
				lib, err := ReadSourceFile(uri.Path())
				if err != nil {
					return "", err
				}
				libcontent, err = Collapse(client, lib)
				if err != nil {
					return "", err
				}
				client.Files.Write(uri, libcontent)
				continue
			}
		}

		if len(stack) > 0 {
			libcontent += line
			continue
		}

		s += line

	}
	return s, nil
}
