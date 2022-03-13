package avrogen

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/hamba/avro"
)

// generic data generator
type fieldGen func() (interface{}, error)

func newFieldGen(rawPattern string, seed int64) fieldGen {
	pattern := parsePattern(rawPattern)
	if pattern == nil {
		return nil
	}
	random := rand.New(rand.NewSource(seed))
	return func() (interface{}, error) {
		var res string
		for _, c := range pattern.content {
			for i := 0; i < c.count; i++ {
				// pick a random patternIdx
				patternIdx := random.Intn(len(c.options))
				// generate the list of options for the selected
				// pattern
				options := c.options[patternIdx]()
				k := random.Intn(len(options))
				res += options[k]
			}
		}
		switch pattern.type_ {
		case string(avro.Boolean):
			return strconv.ParseBool(res)
		case string(avro.Int):
			return strconv.Atoi(res)
		case string(avro.Long):
			return strconv.ParseInt(res, 10, 64)
		case string(avro.Float):
			res, err := strconv.ParseFloat(res, 32)
			if err != nil {
				return nil, err
			}
			return float32(res), nil
		case string(avro.Double):
			return strconv.ParseFloat(res, 64)
		case string(avro.String):
			return res, nil
		default:
			return nil, fmt.Errorf("unsupported avro type %s", pattern.type_)
		}
	}
}

func defaultIntFieldGen(seed int64) fieldGen    { return newFieldGen("{int}[0-9]{4}", seed) }
func defaultLongFieldGen(seed int64) fieldGen   { return newFieldGen("{long}[0-9]{7}", seed) }
func defaultStringFieldGen(seed int64) fieldGen { return newFieldGen("{string}[a-Z|0-9]{10}", seed) }
func defaultFloatFieldGen(seed int64) fieldGen {
	return newFieldGen("{float}[0]{1}[.]{1}[0-9]{3}", seed)
}
func defaultDoubleFieldGen(seed int64) fieldGen {
	return newFieldGen("{double}[0]{1}[.]{1}[0-9]{3}", seed)
}
func defaultBooleanFieldGen(seed int64) fieldGen {
	return newFieldGen("{boolean}[true|false]{1}", seed)
}
