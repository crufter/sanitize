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
		t.Fail()
	}
	_, err = ex.Extract(dat)
	if err == nil {
		t.Fail()
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
		t.Fail()
	}
	_, err = ex.Extract(dat)
	if err == nil {
		t.Fail()
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
		t.Fail()
	}
	val, err := ex.Extract(dat)
	if err != nil || len(val) != 0 {
		t.Fail()
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
		t.Fail()
	}
	val, err := ex.Extract(dat)
	if err != nil || len(val) != 0 {
		t.Fail()
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
		t.Fail()
	}
	_, err = ex.Extract(dat)
	if err == nil {
		t.Fail()
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
		t.Fail()
	}
	val, err := ex.Extract(dat)
	if err != nil || len(val) != 0 {
		t.Fail()
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
		t.Fail()
	}
	_, err = ex.Extract(dat)
	if err == nil {
		t.Fail()
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
		t.Fail()
	}
	_, err = ex.Extract(dat)
	if err == nil {
		t.Fail()
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
		t.Fail()
	}
	_, err = ex.Extract(dat)
	if err == nil {
		t.Fail()
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
		t.Fail()
	}
	_, err = ex.Extract(dat)
	if err == nil {
		t.Fail()
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
		t.Fail()
	}
	val, err := ex.Extract(dat)
	if err != nil || len(val) != 2 {
		t.Fail()
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
		t.Fail()
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
		t.Fail()
	} 
}