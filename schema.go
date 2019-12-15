/*
Copyright 2019 Konstantin Vasilev (burmuley@gmail.com)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "strings"
    "text/template"
)

// data stucture to store information about Go imports
type templateImports struct {
    Imports map[string]struct{}
    Prefix  string
}

type schemaPrefixList map[string]templateImports

//////////////////////////////////////////////////////////////////////
// JSON schema `type` field wrapper
//////////////////////////////////////////////////////////////////////
type schemaTypeWrapper struct {
    Type string `json:"-"`
}

func (s schemaTypeWrapper) String() string {
    return fmt.Sprintf("%s", s.Type)
}

func (s *schemaTypeWrapper) UnmarshalJSON(b []byte) error {
    var tmp interface{}

    if err := json.Unmarshal(b, &tmp); err != nil {
        return schemaError{string(b), err}
    }

    switch tmp.(type) {
    case string:
        s.Type = fmt.Sprintf("%s", tmp)
    case []interface{}:
        s.Type = schemaTypeNumber
    default:
        s.Type = schemaTypeUnknown
        return schemaError{
            errInfo: string(b),
            err:     fmt.Errorf("%s", "unknown schema type"),
        }
    }

    return nil

}

//////////////////////////////////////////////////////////////////////
// JSON schema `items` field wrapper
//////////////////////////////////////////////////////////////////////
type schemaItemsWrapper struct {
    Items    *schemaJSONProperty   `json:"-"`
    ItemsArr []*schemaJSONProperty `json:"-"`
}

func (s schemaItemsWrapper) GetGoType() string {
    if s.ItemsArr != nil {
        return "interface{}"
    } else if s.Items != nil {
        return s.Items.GetGoType()
    }

    return ""
}

func (s schemaItemsWrapper) GetDescription() string {
    return ""
}

func (s *schemaItemsWrapper) UnmarshalJSON(b []byte) error {
    var tmp interface{}

    if err := json.Unmarshal(b, &tmp); err != nil {
        return schemaError{string(b), err}
    }

    switch tmp.(type) {
    case []interface{}:
        if err := json.Unmarshal(b, &s.ItemsArr); err != nil {
            return schemaError{string(b), err}
        }

        return nil
    default:
        if err := json.Unmarshal(b, &s.Items); err != nil {
            return schemaError{string(b), err}
        }

        return nil
    }
}

//////////////////////////////////////////////////////////////////////
// JSON schema element data structure
//////////////////////////////////////////////////////////////////////
type schemaJSONProperty struct {
    Type        schemaTypeWrapper              `json:"type,omitempty"`
    Descr       string                         `json:"description,omitempty"`
    AllOf       []*schemaJSONProperty          `json:"allOf,omitempty"`
    OneOf       []*schemaJSONProperty          `json:"oneOf,omitempty"`
    Properties  map[string]*schemaJSONProperty `json:"properties,omitempty"`
    Required    []string                       `json:"required,omitempty"`
    Enum        []interface{}                  `json:"enum,omitempty"` // TODO: make a wrapper (can be int or string)
    EnumNames   []string                       `json:"enum_names,omitempty"`
    Items       *schemaItemsWrapper            `json:"items,omitempty"`
    Ref         string                         `json:"$ref,omitempty"`
    stripPrefix bool                           `json:"-"`
    addPrefix   string                         `json:"-"`
}

func (s schemaJSONProperty) GetType() string {
    if len(s.AllOf) > 0 || len(s.OneOf) > 0 {
        return schemaTypeMultiple
    } else if len(s.Ref) > 0 {
        return schemaTypeBuiltin
    }

    if len(fmt.Sprint(s.Type)) > 0 {
        return fmt.Sprint(s.Type)
    } else if len(s.Properties) > 0 {
        return schemaTypeObject
    }

    return schemaTypeUnknown
}

func (s schemaJSONProperty) GetGoType() (goTypes string) {
    if s.AllOf != nil {
        tmpArr := make([]string, 0)
        for _, r := range s.AllOf {
            tmpArr = append(tmpArr, r.GetGoType())
        }
        return strings.Join(tmpArr, "\n")
    } else if s.OneOf != nil {
        tmpArr := make([]string, 0)
        for _, r := range s.OneOf {
            tmpArr = append(tmpArr, r.GetGoType())
        }
        return strings.Join(tmpArr, "\n")
    }

    if len(s.Ref) > 0 {
        var ref string

        if s.stripPrefix {
            stripped := strings.Split(s.Ref, "#")
            ref = strings.Join([]string{"#", stripped[len(stripped)-1]}, "")
        } else if len(s.addPrefix) > 0 {
            stripped := strings.Split(s.Ref, "#")
            if stripped[0] != s.addPrefix {
                ref = strings.Join([]string{s.addPrefix, "#", stripped[len(stripped)-1]}, "")
            } else {
                ref = s.Ref
            }
        } else {
            ref = s.Ref
        }

        return getObjectTypeName(ref)
    }

    if fmt.Sprint(s.Type) == schemaTypeArray {
        return s.Items.GetGoType()
    }

    return detectGoType(fmt.Sprintf("%s", s.Type))
}

func (s schemaJSONProperty) GetDescription() string {
    return s.Descr
}

func (s schemaJSONProperty) GetProperties() (pMap map[string]schemaJSONProperty) {
    if len(s.Properties) > 0 {
        pMap = make(map[string]schemaJSONProperty, len(s.Properties))

        for k, v := range s.Properties {
            pMap[k] = *v
        }

        return pMap
    }

    return nil
}

//////////////////////////////////////////////////////////////////////
// VK API type definition
//////////////////////////////////////////////////////////////////////
type typeDefinition map[string]IType

func (o *typeDefinition) Render(tmpl *template.Template) ([]byte, error) {
    var buf bytes.Buffer

    if err := tmpl.Execute(&buf, o); err != nil {
        return []byte{}, err
    }

    return buf.Bytes(), nil

}
