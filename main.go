package mink

import (
	"fmt"
	"strings"
	"time"

	"github.com/snowmerak/mink/gtype"
)

/*
Element ... set of essential method for Mink.
Compare method is comparing I and parameter.
Filter method is
*/
type Element interface{}

//Mink ... struct declared linq method with Element slice.
type Mink struct {
	str   string
	slice []Element
	maps  map[Element]Element
	err   error
	typ   uint8
}

//From ... from slice or string to Mink
func From(t interface{}) Mink {
	g := Mink{str: "", slice: nil, maps: nil, err: nil, typ: 0}

	if v, ok := t.([]Element); ok {
		g.slice = v
		g.typ = gtype.Slice
		return g
	}
	if v, ok := t.(map[Element]Element); ok {
		g.maps = v
		g.typ = gtype.Maps
		return g
	}
	if v, ok := t.(string); ok {
		g.str = v
		g.typ = gtype.Str
		return g
	}
	g.err = fmt.Errorf("%v: From: parameter is invalid", time.Now())
	return g
}

//Filter ... apply criterion
func (g Mink) Filter(c func(e Element) bool) Mink {
	switch g.typ {
	case gtype.Slice:
		ns := make([]Element, 0, len(g.slice)/2)
		for _, v := range g.slice {
			if c(v) {
				ns = append(ns, v)
			}
		}
		g.slice = ns
	case gtype.Str:
		sb := strings.Builder{}
		for _, v := range g.str {
			if c(v) {
				sb.WriteRune(v)
			}
		}
		g.str = sb.String()
	case gtype.Maps:
		nm := map[Element]Element{}
		for k, v := range g.maps {
			if c(v) {
				nm[k] = v
			}
		}
		g.maps = nm
	}
	return g
}

/*
OrderBy ... order by parameter
the paremeter return when choose a true, not false
*/
func (g Mink) OrderBy(f func(a Element, b Element) bool) Mink {
	if g.typ != gtype.Slice {
		return g
	}
	rt := <-asyncMergeSortSlice(g.slice, f)
	g.slice = rt
	return g
}

/*
Unwrap ... unwrap to interface{} can convert to slice, map, string, error
*/
func (g Mink) Unwrap() interface{} {
	if g.err != nil {
		return g.err
	}
	switch g.typ {
	case gtype.Str:
		return g.str
	case gtype.Slice:
		return g.slice
	case gtype.Maps:
		return g.maps
	}
	return nil
}
