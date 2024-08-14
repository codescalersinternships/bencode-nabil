package bencoder

import (
	"fmt"
	"strconv"
	"strings"
)

func readInteger(idx *int, str *string, terminator string) (int64, error) {
	out := ""
	lim := strings.Index((*str)[*idx:], terminator)
	if lim == -1 {
		return -1, fmt.Errorf("%s","error of terminator wasn't found")
	}
	lim = lim + (*idx)
	out = (*str)[*idx:lim]
	*idx = lim 
	integer, err := strconv.ParseInt(out, 10, 64)
	if err != nil {
		return -1, err
	}
	return integer, nil
}

func readBulkString(idx *int, str *string) (interface{}, error) {
	out := ""
	length, err := readInteger(idx, str,":")
	if err != nil {
		return "", err
	}
	if length == -1 {
		return nil, nil
	}
	(*idx)++
	lim := (*idx) + int(length)
	if lim > len(*str) {
		return "", fmt.Errorf("%s","error of string lenght isn't enough")
	}
	out = (*str)[*idx:lim]
	if out == "-1" {
		return "", nil
	}
	*idx = lim -1
	return out, nil
}


func readMap(idx *int, str *string) (map[interface{}]interface{}, error) {
	var out = make(map[interface{}]interface{})
	var prev interface{}
	for i := 0; ; i++ {
		var val interface{}
		if (*str)[*idx] == 'i' {
			(*idx)++
			integer, err := readInteger(idx, str,"e")
			if err != nil {
				return nil, err
			}
			val = integer
			(*idx)++
			if i%2 == 1 {
				out[prev] = val
			}
	
			prev = val
			continue
		}
		if (*str)[*idx] >= '0' && (*str)[*idx] <= '9' {
			bulkString, err := readBulkString(idx, str)
			if err != nil {
				return nil, err
			}
			val = bulkString
			(*idx)++
			if i%2 == 1 {
				out[prev] = val
			}
	
			prev = val
			continue
		}
		if (*str)[*idx] == 'l' {
			(*idx)++
			array, err := Decoder(*str, idx)
			if err != nil {
				return nil, err
			}
			val = array
			(*idx)++
			if i%2 == 1 {
				out[prev] = val
			}
	
			prev = val
			continue
		}
		
		if (*str)[*idx] == 'd' {
			(*idx)++
			respMap, err := readMap(idx, str)
			if err != nil {
				return nil, err
			}
			val = respMap
			(*idx)++
			if i%2 == 1 {
				out[prev] = val
			}
	
			prev = val
			continue
		}
		if (*str)[*idx] == 'e' {
			break
		}
	}
	return out, nil
}

// Decoder reads an bencoded string, returning an array of interfaces of items
func Decoder(str string, start ...*int) (interface{}, error) {
	var out []interface{}
	var idx int = 0
	if len(start) > 0 {
		idx = *start[0]
	}
	for ; idx < len(str); idx++ {
		if str[idx] == 'i' {
			idx++
			integer, err := readInteger(&idx, &str,"e")
			if err != nil {
				return nil, err
			}
			out = append(out, integer)
			continue
		}
		if str[idx] >= '0' && str[idx] <= '9' {
			bulkString, err := readBulkString(&idx, &str)
			if err != nil {
				return nil, err
			}
			out = append(out, bulkString)
			continue
		}
		if str[idx] == 'l' {
			idx++
			array, err := Decoder(str, &idx)
			if err != nil {
				return nil, err
			}
			out = append(out, array)
			continue
		}
		
		if str[idx] == 'd' {
			idx++
			respMap, err := readMap(&idx, &str)
			if err != nil {
				return nil, err
			}
			out = append(out, respMap)
			continue
		}
		if str[idx] == 'e' {
			break
		}
	}
	if len(start) > 0 {
		*start[0] = idx
	}
	if len(out) == 1 {
		return out[0], nil
	}
	return out, nil
}

