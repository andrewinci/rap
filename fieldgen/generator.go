package fieldgen

import (
	"fmt"
	"math/rand"
	"strconv"
)

// generic data generator
type FieldGen func() (interface{}, error)

func NewFieldGenerator(rawPattern string, seed int64) FieldGen {
	pattern := parsePattern(rawPattern)
	if pattern == nil {
		return nil
	}
	random := rand.New(rand.NewSource(seed))
	return func() (interface{}, error) {
		var res string
		for _, c := range pattern.content {
			for i := 0; i < c.count; i++ {
				j := random.Intn(len(c.options))
				k := random.Intn(len(c.options[j]))
				res += c.options[j][k]
			}
		}
		switch pattern.type_ {
		case avro_boolean:
			return strconv.ParseBool(res)
		case avro_int:
			return strconv.Atoi(res)
		case avro_long:
			return strconv.ParseInt(res, 10, 64)
		case avro_float:
			return strconv.ParseFloat(res, 32)
		case avro_double:
			return strconv.ParseFloat(res, 32)
		case avro_string:
			return res, nil
		default:
			return nil, fmt.Errorf("unsupported avro type %s", pattern.type_)
		}
	}
}
