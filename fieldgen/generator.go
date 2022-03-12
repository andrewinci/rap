package fieldgen

// generic data generator
type Gen func() interface{}

func NewGenerator(pattern string) Gen {
	return func() interface{} {
		return 1
	}
}
