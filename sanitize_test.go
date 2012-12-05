package sanitize_test

import(
	"testing"
	"github.com/opesun/sanitize"
	"fmt"
)

func TestStringMinMust(t *testing.T) {
	dat := map[string]interface{}{
		"a":	"Hey there!",
	}
	scheme1 := map[string]interface{}{
		"a": map[string]interface{}{
			"must": true,
			"type": "string",
			"min":	100,
		},
	}
	ex, err := sanitize.New(scheme1)
	if err != nil {
		t.Fatal()
	}
	_, err = ex.Extract(dat)
	if err == nil {
		t.Fatal()
	}
}

func TestStringMaxMust(t *testing.T) {
	dat := map[string]interface{}{
		"a":	"Hey there!",
	}
	scheme1 := map[string]interface{}{
		"a": map[string]interface{}{
			"must": true,
			"type": "string",
			"max":	3,
		},
	}
	ex, err := sanitize.New(scheme1)
	if err != nil {
		t.Fatal()
	}
	_, err = ex.Extract(dat)
	if err == nil {
		t.Fatal()
	}
}

func TestStringMin(t *testing.T) {
	dat := map[string]interface{}{
		"a":	"Hey there!",
	}
	scheme1 := map[string]interface{}{
		"a": map[string]interface{}{
			"type": "string",
			"min":	100,
		},
	}
	ex, err := sanitize.New(scheme1)
	if err != nil {
		t.Fatal()
	}
	val, err := ex.Extract(dat)
	if err != nil || len(val) != 0 {
		t.Fatal()
	}
}

func TestStringMax(t *testing.T) {
	dat := map[string]interface{}{
		"a":	"Hey there!",
	}
	scheme1 := map[string]interface{}{
		"a": map[string]interface{}{
			"type": "string",
			"max":	3,
		},
	}
	ex, err := sanitize.New(scheme1)
	if err != nil {
		t.Fatal()
	}
	val, err := ex.Extract(dat)
	if err != nil || len(val) != 0 {
		t.Fatal()
	}
}

func TestIntMust(t *testing.T) {
	dat := map[string]interface{}{
		"a":	"Hey there!",
	}
	scheme1 := map[string]interface{}{
		"a": map[string]interface{}{
			"type": "int",
			"max":	3,
			"must": true,
		},
	}
	ex, err := sanitize.New(scheme1)
	if err != nil {
		t.Fatal()
	}
	_, err = ex.Extract(dat)
	if err == nil {
		t.Fatal()
	}
}

func TestInt(t *testing.T) {
	dat := map[string]interface{}{
		"a":	"Hey there!",
	}
	scheme1 := map[string]interface{}{
		"a": map[string]interface{}{
			"type": "int",
			"max":	3,
		},
	}
	ex, err := sanitize.New(scheme1)
	if err != nil {
		t.Fatal()
	}
	val, err := ex.Extract(dat)
	if err != nil || len(val) != 0 {
		t.Fatal()
	}
}

func TestIntMax(t *testing.T) {
	dat := map[string]interface{}{
		"a":	900,
	}
	scheme1 := map[string]interface{}{
		"a": map[string]interface{}{
			"type": "int",
			"max":	700,
			"must": true,
		},
	}
	ex, err := sanitize.New(scheme1)
	if err != nil {
		t.Fatal()
	}
	_, err = ex.Extract(dat)
	if err == nil {
		t.Fatal()
	}
}

func TestMin(t *testing.T) {
	dat := map[string]interface{}{
		"a":	500,
	}
	scheme1 := map[string]interface{}{
		"a": map[string]interface{}{
			"type": "int",
			"min":	600,
			"must": true,
		},
	}
	ex, err := sanitize.New(scheme1)
	if err != nil {
		t.Fatal()
	}
	_, err = ex.Extract(dat)
	if err == nil {
		t.Fatal()
	}
}

func TestFloatMust(t *testing.T) {
	dat := map[string]interface{}{
		"a":	"adsad",
	}
	scheme1 := map[string]interface{}{
		"a": map[string]interface{}{
			"type": "float",
			"must": true,
		},
	}
	ex, err := sanitize.New(scheme1)
	if err != nil {
		t.Fatal()
	}
	_, err = ex.Extract(dat)
	if err == nil {
		t.Fatal()
	}
}

func TestBoolMust(t *testing.T) {
	dat := map[string]interface{}{
		"a":	"asdasd",
	}
	scheme1 := map[string]interface{}{
		"a": map[string]interface{}{
			"type": "bool",
			"must": true,
		},
	}
	ex, err := sanitize.New(scheme1)
	if err != nil {
		t.Fatal()
	}
	_, err = ex.Extract(dat)
	if err == nil {
		t.Fatal()
	}
}

func TestComposite(t *testing.T) {
	dat := map[string]interface{}{
		"a":	"asdasd",
		"b":	"20",
		"c":	"Hey there.",
	}
	scheme1 := map[string]interface{}{
		"a": map[string]interface{}{
			"type": "bool",
		},
		"b": map[string]interface{}{
			"type": "float",
			"must":	true,
		},
		"c": map[string]interface{}{
			"type": "string",
			"must": true,
		},
	}
	ex, err := sanitize.New(scheme1)
	if err != nil {
		t.Fatal()
	}
	val, err := ex.Extract(dat)
	if err != nil || len(val) != 2 {
		t.Fatal()
	}
}

func TestUserDefinedType(t *testing.T) {
	dat := map[string]interface{}{
		"a":	"asdasd1",
	}
	scheme1 := map[string]interface{}{
		"a": map[string]interface{}{
			"type": "myType",
			"must": true,
		},
	}
	ex, err := sanitize.New(scheme1)
	if err != nil {
		t.Fatal()
	}
	ex.AddFuncs(sanitize.FuncMap{
		"myType": func(dat interface{}, s sanitize.Scheme) (interface{}, error) {
			val, ok := dat.(string)
			if !ok {
				return nil, fmt.Errorf("Baaad.")
			}
			if val == "asdasd" {
				return val, nil
			}
			return nil, fmt.Errorf("This is baad too.")
		},
	})
	_, err = ex.Extract(dat)
	if err == nil {
		t.Fatal()
	}
}

func TestSlice(t *testing.T) {
	dat := map[string]interface{}{
		"a": []interface{}{30, 20, "xxd"},
	}
	scheme1 := map[string]interface{}{
		"a": map[string]interface{}{
			"type": 	"string",
			"slice": 	true,
		},
	}
	ex, err := sanitize.New(scheme1)
	if err != nil {
		t.Fatal()
	}
	val, err := ex.Extract(dat)
	if err != nil || len(val["a"].([]interface{})) != 1 {
		t.Fatal()
	}
}

func TestSliceMustAllOrNothing(t *testing.T) {
	dat := map[string]interface{}{
		"a": []interface{}{30, 20, "xxd"},
	}
	scheme1 := map[string]interface{}{
		"a": map[string]interface{}{
			"type": 		"string",
			"slice":		true,
			"allOrNothing": true,
			"must":			true,
		},
	}
	ex, err := sanitize.New(scheme1)
	if err != nil {
		t.Fatal()
	}
	_, err = ex.Extract(dat)
	if err == nil {
		t.Fatal()
	}
}

func TestShorthand1(t *testing.T) {
	dat := map[string]interface{}{
		"a": 1,
		"b": "asdsad",
		"c": "Hello.",
	}
	scheme1 := map[string]interface{}{
		"a": 1,
		"b": 1,
		"c": 1,
	}
	ex, err := sanitize.New(scheme1)
	if err != nil {
		t.Fatal()
	}
	val, err := ex.Extract(dat)
	if err != nil || len(val) != 2 {
		t.Fatal()
	}
}

func TestIgnore(t *testing.T) {
	dat := map[string]interface{}{
		"a": "asdsadasd",
		"b": "asdsad",
		"c": "Hello.",
	}
	scheme1 := map[string]interface{}{
		"a": map[string]interface{}{
			"type": 	"string",
			"ignore":	true,
		},
		"b": 1,
		"c": 1,
	}
	ex, err := sanitize.New(scheme1)
	if err != nil {
		t.Fatal()
	}
	val, err := ex.Extract(dat)
	if err != nil || len(val) != 2 {
		t.Fatal()
	}
}

func TestIgnoreMust(t *testing.T) {
	dat := map[string]interface{}{
		"a": "asdsadasd",
		"b": "asdsad",
		"c": "Hello.",
	}
	scheme1 := map[string]interface{}{
		"a": map[string]interface{}{
			"type": 	"string",
			"ignore":	true,
			"min":		1000,
			"must":		true,
		},
		"b": 1,
		"c": 1,
	}
	ex, err := sanitize.New(scheme1)
	if err != nil {
		t.Fatal()
	}
	_, err = ex.Extract(dat)
	if err == nil {
		t.Fatal()
	}
}

func TestConst(t *testing.T) {
	dat := map[string]interface{}{}
	scheme := map[string]interface{}{
		"a": map[string]interface{}{
			"type":		"const",
			"value":	"example",
		},
	}
	ex, err := sanitize.New(scheme)
	if err != nil {
		t.Fatal()
	}
	res, err := ex.Extract(dat)
	if err != nil {
		t.Fatal(err)
	}
	if res["a"] != "example" {
		t.Fatal(res)
	}
}

func TestEq(t *testing.T) {
	dat := map[string]interface{}{
		"a": "example1234",
	}
	scheme := map[string]interface{}{
		"a": map[string]interface{}{
			"type": "eq",
			"value": "example123",
			"must": true,
		},
	}
	res, err := sanitize.Fast(scheme, dat)
	if err == nil {
		t.Fatal()
	}
	dat1 := map[string]interface{}{
		"a": "example123",
	}
	res, err = sanitize.Fast(scheme, dat1)
	if err != nil {
		t.Fatal(err)
	}
	if res["a"] != "example123" {
		t.Fatal(res["a"])
	}
}