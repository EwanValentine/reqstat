package analyzer

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type FieldType string

const (
	TypeString  FieldType = "string"
	TypeNumber  FieldType = "number"
	TypeBool    FieldType = "bool"
	TypeNull    FieldType = "null"
	TypeObject  FieldType = "object"
	TypeArray   FieldType = "array"
	TypeUnknown FieldType = "unknown"
)

type Field struct {
	Name     string
	Type     FieldType
	Children []*Field
	Example  any
	Count    int
}

type Schema struct {
	Root       *Field
	TotalKeys  int
	MaxDepth   int
	ArrayItems int
}

func AnalyzeJSON(data []byte) (*Schema, error) {
	var parsed any
	if err := json.Unmarshal(data, &parsed); err != nil {
		return nil, err
	}

	schema := &Schema{}
	schema.Root = analyzeValue("root", parsed, 0, schema)

	return schema, nil
}

func analyzeValue(name string, value any, depth int, schema *Schema) *Field {
	if depth > schema.MaxDepth {
		schema.MaxDepth = depth
	}

	field := &Field{
		Name:  name,
		Count: 1,
	}

	switch v := value.(type) {
	case nil:
		field.Type = TypeNull
		field.Example = nil

	case bool:
		field.Type = TypeBool
		field.Example = v

	case float64:
		field.Type = TypeNumber
		field.Example = v

	case string:
		field.Type = TypeString
		if len(v) > 50 {
			field.Example = v[:50] + "..."
		} else {
			field.Example = v
		}

	case []any:
		field.Type = TypeArray
		schema.ArrayItems += len(v)
		if len(v) > 0 {
			child := analyzeValue("[item]", v[0], depth+1, schema)
			field.Children = append(field.Children, child)
			field.Count = len(v)
		}

	case map[string]any:
		field.Type = TypeObject
		schema.TotalKeys += len(v)

		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			child := analyzeValue(k, v[k], depth+1, schema)
			field.Children = append(field.Children, child)
		}

	default:
		field.Type = TypeUnknown
	}

	return field
}

func (s *Schema) String() string {
	if s.Root == nil {
		return "empty"
	}
	var sb strings.Builder
	s.Root.format(&sb, 0)
	return sb.String()
}

func (f *Field) format(sb *strings.Builder, indent int) {
	prefix := strings.Repeat("  ", indent)

	switch f.Type {
	case TypeObject:
		if f.Name == "root" {
			sb.WriteString("{\n")
		} else {
			sb.WriteString(fmt.Sprintf("%s%s: {\n", prefix, f.Name))
		}
		for _, child := range f.Children {
			child.format(sb, indent+1)
		}
		if f.Name == "root" {
			sb.WriteString("}")
		} else {
			sb.WriteString(fmt.Sprintf("%s}\n", prefix))
		}

	case TypeArray:
		countInfo := ""
		if f.Count > 1 {
			countInfo = fmt.Sprintf(" (%d items)", f.Count)
		}
		if len(f.Children) > 0 {
			child := f.Children[0]
			if child.Type == TypeObject {
				sb.WriteString(fmt.Sprintf("%s%s: [%s\n", prefix, f.Name, countInfo))
				sb.WriteString(fmt.Sprintf("%s  {\n", prefix))
				for _, grandchild := range child.Children {
					grandchild.format(sb, indent+2)
				}
				sb.WriteString(fmt.Sprintf("%s  }\n", prefix))
				sb.WriteString(fmt.Sprintf("%s]\n", prefix))
			} else {
				sb.WriteString(fmt.Sprintf("%s%s: [%s]%s\n", prefix, f.Name, child.Type, countInfo))
			}
		} else {
			sb.WriteString(fmt.Sprintf("%s%s: []%s\n", prefix, f.Name, countInfo))
		}

	default:
		example := ""
		if f.Example != nil {
			switch e := f.Example.(type) {
			case string:
				example = fmt.Sprintf(" // e.g. %q", e)
			case float64:
				if e == float64(int(e)) {
					example = fmt.Sprintf(" // e.g. %d", int(e))
				} else {
					example = fmt.Sprintf(" // e.g. %.2f", e)
				}
			case bool:
				example = fmt.Sprintf(" // e.g. %t", e)
			}
		}
		sb.WriteString(fmt.Sprintf("%s%s: %s%s\n", prefix, f.Name, f.Type, example))
	}
}

func (s *Schema) Summary() string {
	return fmt.Sprintf("Keys: %d | Depth: %d | Array items: %d",
		s.TotalKeys, s.MaxDepth, s.ArrayItems)
}
