package bencoder

import (
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

func readInteger(idx *int, str *string, terminator string) (int64, error) {
	out := ""
	lim := strings.Index((*str)[*idx:], terminator)
	if lim == -1 {
		return -1, fmt.Errorf("%s", "error of terminator wasn't found")
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

func readBulkString(idx *int, str *string) (string, error) {
	out := ""
	length, err := readInteger(idx, str, ":")
	if err != nil {
		return "", err
	}
	if length == -1 {
		return "", nil
	}
	(*idx)++
	lim := (*idx) + int(length)
	if lim > len(*str) {
		return "", fmt.Errorf("%s", "error of string lenght isn't enough")
	}
	out = (*str)[*idx:lim]
	if out == "-1" {
		return "", nil
	}
	*idx = lim - 1
	return out, nil
}

func readMap(idx *int, str *string) (map[interface{}]interface{}, error) {
	var out = make(map[interface{}]interface{})
	var prev interface{}
	for i := 0; ; i++ {
		var val interface{}
		if (*str)[*idx] == 'i' {
			(*idx)++
			integer, err := readInteger(idx, str, "e")
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
			integer, err := readInteger(&idx, &str, "e")
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

// Encoder reads an interface, returning an byte array of the bencoded interface
func Encoder(benco interface{}) ([]byte, error) {

	var encoded []byte

	switch ty := benco.(type) {
	case int, int64:
		encoded = append(encoded, 'i')
		encoded = append(encoded, []byte(strconv.FormatInt(reflect.ValueOf(ty).Int(), 10))...)
		encoded = append(encoded, 'e')

	case string:
		encoded = append(encoded, []byte(strconv.Itoa(len(ty)))...)
		encoded = append(encoded, ':')
		encoded = append(encoded, []byte(ty)...)

	case []interface{}:
		encoded = append(encoded, 'l')
		for _, item := range ty {
			encodedItem, err := Encoder(item)
			if err != nil {
				return nil, err
			}
			encoded = append(encoded, encodedItem...)
		}
		encoded = append(encoded, 'e')

	case map[string]interface{}:
		encoded = append(encoded, 'd')
		var keyArr []string
		for key := range ty {
			keyArr = append(keyArr, key)
		}
		slices.Sort(keyArr)
		for _, key := range keyArr {
			encodedKey, err := Encoder(key)
			if err != nil {
				return nil, err
			}
			encodedValue, err := Encoder(ty[key])
			if err != nil {
				return nil, err
			}
			encoded = append(encoded, encodedKey...)
			encoded = append(encoded, encodedValue...)
		}
		encoded = append(encoded, 'e')

	default:
		return nil, fmt.Errorf("unsupported data type")
	}

	return encoded, nil
}
