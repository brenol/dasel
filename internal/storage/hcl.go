package storage

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

func init() {
	registerReadParser([]string{"hcl"}, []string{".hcl"}, &HCLParser{})
	registerWriteParser([]string{"hcl"}, []string{".hcl"}, &HCLParser{})
}

// HCLParser is a Parser implementation to handle hcl2 files.
type HCLParser struct {
}

// FromBytes returns some data that is represented by the given bytes.
func (p *HCLParser) FromBytes(byteData []byte) (interface{}, error) {
	file, diags := hclsyntax.ParseConfig(byteData, "x", hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, fmt.Errorf("could not parse hcl config: %s", diags.Error())
	}

	target := map[string]interface{}{}

	diags = gohcl.DecodeBody(file.Body, nil, &target)
	if diags.HasErrors() {
		return nil, fmt.Errorf("could not decode hcl body: %s", diags.Error())
	}

	return &BasicSingleDocument{Value: hclToMap(target)}, nil
}

func hclToMap(data interface{}) interface{} {
	switch val := data.(type) {
	case map[string]interface{}:
		for k, v := range val {
			val[k] = hclToMap(v)
		}
		return val
	case *hcl.Attribute:
		x, _ := val.Expr.Value(nil)
		switch x.Type() {
		case cty.Bool:
			return x.True()
		case cty.Number:
			floatVal, _ := x.AsBigFloat().Float64()
			return floatVal
		case cty.String:
			return x.AsString()
		}
		fmt.Printf("%s, %s\n", val.Name, x.Type().GoString())
	}
	return data
}

// ToBytes returns a slice of bytes that represents the given value.
func (p *HCLParser) ToBytes(value interface{}, options ...ReadWriteOption) ([]byte, error) {
	// buffer := new(bytes.Buffer)
	// encoder := hcl.NewEncoder(buffer)
	//
	// indent := "  "
	// prettyPrint := true
	//
	// for _, o := range options {
	// 	switch o.Key {
	// 	case OptionIndent:
	// 		if value, ok := o.Value.(string); ok {
	// 			indent = value
	// 		}
	// 	case OptionPrettyPrint:
	// 		if value, ok := o.Value.(bool); ok {
	// 			prettyPrint = value
	// 		}
	// 	}
	// }
	//
	// if !prettyPrint {
	// 	indent = ""
	// }
	// encoder.SetIndent("", indent)
	//
	// switch v := value.(type) {
	// case SingleDocument:
	// 	if err := encoder.Encode(v.Document()); err != nil {
	// 		return nil, fmt.Errorf("could not encode single document: %w", err)
	// 	}
	// case MultiDocument:
	// 	for index, d := range v.Documents() {
	// 		if err := encoder.Encode(d); err != nil {
	// 			return nil, fmt.Errorf("could not encode multi document [%d]: %w", index, err)
	// 		}
	// 	}
	// default:
	// 	if err := encoder.Encode(v); err != nil {
	// 		return nil, fmt.Errorf("could not encode default document type: %w", err)
	// 	}
	// }
	// return buffer.Bytes(), nil
	return nil, nil
}