// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package common

type stringConverter struct {
	c map[string]string
	r map[string]string
}

func NewStringConverter(input map[string]string) stringConverter {
	reversed := make(map[string]string)
	for k, v := range input {
		reversed[v] = k
	}
	return stringConverter{
		c: input,
		r: reversed,
	}
}

func (c stringConverter) Convert(src string) string {
	if dst, ok := c.c[src]; ok {
		return dst
	}
	return src
}

func (c stringConverter) UnConvert(dst string) string {
	if src, ok := c.r[dst]; ok {
		return src
	}
	return dst
}
