package sanitize

import(
	"strconv"
	"fmt"
	"github.com/opesun/numcon"
	"reflect"
)

func Fast(sch, data map[string]interface{}) (map[string]interface{}, error) {
	e, err := New(sch)
	if err != nil {
		return nil, err
	}
	return e.Extract(data)
}

type Scheme struct {
	Must			bool
	Type			string
	Slice			bool
	SliceMin		int
	SliceMax		int
	AllOrNothing	bool
	Min				int64
	Max				int64
	Regexp			string
	// Below things are for supporting hackish things.
	Ignore			bool
	Specific		map[string]interface{}
	Key				string		// Only field not suppliable, it's just metainformation for validation handlers.
}

type SchemeMap map[string]Scheme

func toScheme(a interface{}) (Scheme, error) {
	empty := Scheme{}
	s := Scheme{}
	s.Specific = map[string]interface{}{}
	ai, err := numcon.Int(a)
	if err == nil && ai == 1 {
		s.Type = "string"
		return s, nil
	}
	am, ok := a.(map[string]interface{})
	if !ok {
		return empty, fmt.Errorf("Can't interpret scheme.")
	}
	s.SliceMin = -1		// 0 is not good as a default value since it can be a legal one too.
	s.SliceMax = -1
	for i, v := range am {
		switch i {
		case "must":
			s.Must, ok = v.(bool)
			if !ok {
				return empty, fmt.Errorf("Bad scheme: must should be a slice.")
			}
		case "type":
			s.Type, ok = v.(string)
			if !ok {
				return empty, fmt.Errorf("Bad scheme: type should be a string.")
			}
		case "slice":
			s.Slice, ok = v.(bool)
			if !ok {
				return empty, fmt.Errorf("Bad scheme: slice should be a bool.")
			}
		case "sliceMin":
			s.SliceMin, err = numcon.Int(v)
			if err != nil {
				return empty, fmt.Errorf("Bad scheme: sliceMin should be int compatible.")
			}
		case "sliceMax":
			s.SliceMax, err = numcon.Int(v)
			if err != nil {
				return empty, fmt.Errorf("Bad scheme: sliceMax should be int compatible.")
			}
		case "allOrNothing":
			s.AllOrNothing, ok = v.(bool)
			if !ok {
				return empty, fmt.Errorf("Bad scheme: allOrNothing should be a bool.")
			}
		case "min":
			s.Min, err = numcon.Int64(v)
			if err != nil {
				return empty, fmt.Errorf("Bad scheme: min should be int64 compatible.")
			}
		case "max":
			s.Max, err = numcon.Int64(v)
			if err != nil {
				return empty, fmt.Errorf("Bad scheme: max should be int64 compatible.")
			}
		case "regexp":
			s.Regexp, ok = v.(string)
			if !ok {
				return empty, fmt.Errorf("Bad scheme: regexp should be a string.")
			}
		case "ignore":
			s.Ignore, ok = v.(bool)
			if !ok {
				return empty, fmt.Errorf("Bad scheme: ignore should be a string.")
			}
		default:
			s.Specific[i] = v
		}
	}
	return s, nil
}

func toSchemeMap(a map[string]interface{}) (SchemeMap, error) {
	s := SchemeMap{}
	for i, v := range a {
		val, err := toScheme(v)
		if err != nil {
			return nil, err
		}
		val.Key = i
		s[i] = val
	}
	return s, nil
}

type FuncMap map[string]func(interface{}, Scheme) (interface{}, error)

type Extractor struct {
	SchemeMap 	SchemeMap
	FuncMap		FuncMap
}

func booler(dat interface{}, s Scheme) (interface{}, error) {
	switch v := dat.(type) {
	case string:
		if v == "false" {
			return false, nil
		}
		if v == "true" {
			return true, nil
		}
	case bool:
		return dat, nil
	}
	return nil, fmt.Errorf("Can't interpret.")
}

func stringer(dat interface{}, s Scheme) (interface{}, error) {
	dat_str, ok := dat.(string)
	if !ok {
		return nil, fmt.Errorf("Not a string.")
	}
	if len(dat_str) < int(s.Min) {
		return nil, fmt.Errorf("String is too short.")
	}
	if s.Max != 0 && len(dat_str) > int(s.Max) {
		return nil, fmt.Errorf("String is too long.")
	}
	// Insert regexp check here.
	return dat_str, nil
}

func floater(dat interface{}, s Scheme) (interface{}, error) {
	var val float64
	var err error
	switch v := dat.(type) {
	case string:
		val, err = strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, err
		}
	default:
		val, err = numcon.Float64(v)
		if err != nil {
			return nil, err
		}
	}
	if s.Min > int64(val) {
		return nil, fmt.Errorf("Float value is too small.")
	}
	if s.Max != 0 && s.Max < int64(val) {
		return nil, fmt.Errorf("Float value is too large.")
	}
	return val, nil
}

func inter(dat interface{}, s Scheme) (interface{}, error) {
	var val int64
	var err error
	switch v := dat.(type) {
	case string:
		val, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
	default:
		val, err = numcon.Int64(v)
		if err != nil {
			return nil, err
		}
	}
	if s.Min > val {
		return nil, fmt.Errorf("Int value is too small.")
	}
	if s.Max != 0 && s.Max < val {
		return nil, fmt.Errorf("Int value is too large.")
	}
	return val, nil
}

func conster(dat interface{}, s Scheme) (interface{}, error) {
	return s.Specific["value"], nil
}

func eqer(dat interface{}, s Scheme) (interface{}, error) {
	val := s.Specific["value"]
	if reflect.DeepEqual(dat, val) {
		return dat, nil
	}
	return nil, fmt.Errorf("%v is not equal to %v.", s.Key, val)
}

func anier(dat interface{}, s Scheme) (interface{}, error) {
	if dat == nil {
		return nil, fmt.Errorf("any means anything but not nothing")
	}
	return dat, nil
}

func New(scheme_map map[string]interface{}) (*Extractor, error) {
	schemeMap, err := toSchemeMap(scheme_map)
	if err != nil {
		return nil, err
	}
	funcMap := map[string]func(interface{}, Scheme)(interface{}, error) {
		"string": 	stringer,
		"float":	floater,
		"bool":		booler,
		"int":		inter,
		"const":	conster,
		"eq":		eqer,
		"any":		anier,
	}
	return &Extractor{schemeMap, funcMap}, nil
}

func (e *Extractor) AddFuncs(a FuncMap) {
	for i, v := range a {
		e.FuncMap[i] = v
	}
}

func (e *Extractor) Extract(data map[string]interface{}) (map[string]interface{}, error) {
	ret := map[string]interface{}{}
	for i, v := range e.SchemeMap {
		current := data[i]
		c_func, ok := e.FuncMap[v.Type]
		if !ok {
			return nil, fmt.Errorf("No handler for type %v.", i)
		}
		if slice, ok := current.([]interface{}); ok {
			if !v.Slice {
				return nil, fmt.Errorf("This shouldn't be a slice.")
			}
			f := []interface{}{}
			for _, v1 := range slice {
				if v.SliceMax > -1 && len(f) > v.SliceMax {
					break
				}
				val, err := c_func(v1, v)
				if err != nil {
					if v.AllOrNothing {
						return nil, fmt.Errorf("Slice member is no good..")
					} else {
						continue
					}
				}
				f = append(f, val)
			}
			if v.SliceMin > -1 && len(f) < v.SliceMin {
				if v.Must {
					return nil, fmt.Errorf("Slice length is too small.")
				} else {
					continue
				}
			}
			if !v.Ignore {
				ret[i] = f
			}
		} else {
			val, err := c_func(current, v)
			if err != nil {
				if v.Must {
					return nil, err
				} else {
					continue
				}
			}
			if !v.Ignore {
				ret[i] = val
			}
		}
	}
	return ret, nil
}