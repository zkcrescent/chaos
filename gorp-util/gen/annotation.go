package main

import "strings"

func newAnnotation(comment string) *annotation {
	annotation := new(annotation)
	s := strings.TrimPrefix(comment, "//")
	s = strings.TrimSpace(s)
	if !strings.HasPrefix(s, "@") {
		return nil
	}

	tmp := strings.Split(s, "=")
	nameKey := strings.Split(tmp[0], "(")
	annotation.Name = strings.TrimSpace(nameKey[0])
	if len(nameKey) > 1 {
		nameKey[1] = strings.Replace(nameKey[1], ")", "", 1)
		annotation.Key = strings.TrimSpace(nameKey[1])
	}
	var vals []string
	if len(tmp) > 1 {
		for _, f := range strings.Split(tmp[1], ",") {
			vals = append(vals, strings.TrimSpace(f))
		}
	}
	annotation.Vals = vals

	return annotation
}

type annotation struct {
	Name string
	Key  string
	Vals []string
}
