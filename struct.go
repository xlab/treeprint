package treeprint

import (
	"fmt"
	"reflect"
	"strings"
)

type StructTreeOption int

const (
	StructNameTree StructTreeOption = iota
	StructValueTree
	StructTagTree
	StructTypeTree
	StructTypeSizeTree
	StructFormattedTree
)

func Repr(v interface{}) string {
	tree := New()
	err := valueTree(tree, v)
	if err != nil {
		return err.Error()
	}
	return tree.String()
}

func FromStruct(v interface{}, opt ...StructTreeOption) (Tree, error) {
	var treeOpt StructTreeOption
	if len(opt) > 0 {
		treeOpt = opt[0]
	}
	switch treeOpt {
	case StructNameTree:
		tree := New()
		err := nameTree(tree, v)
		return tree, err
	case StructValueTree:
		tree := New()
		err := valueTree(tree, v)
		return tree, err
	case StructTagTree:
		tree := New()
		err := tagTree(tree, v)
		return tree, err
	case StructTypeTree:
		tree := New()
		err := typeTree(tree, v)
		return tree, err
	case StructTypeSizeTree:
		tree := New()
		err := typeSizeTree(tree, v)
		return tree, err
	default:
		err := fmt.Errorf("treeprint: invalid StructTreeOption %v", treeOpt)
		return nil, err
	}
}

func nameTree(tree Tree, v interface{}) error {
	typ, val, err := checkType(v)
	if err != nil {
		return err
	}
	fields := typ.NumField()
	for i := 0; i < fields; i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)
		name, omit := getMeta(field.Name, field.Tag)
		if omit && isEmpty(&fieldValue) {
			continue
		}
		typ, val, isStruct := getValue(field.Type, &fieldValue)
		if !isStruct {
			tree.AddNode(name)
			continue
		} else if subNum := typ.NumField(); subNum == 0 {
			tree.AddNode(name)
			continue
		}
		branch := tree.AddBranch(name)
		if err := nameTree(branch, val.Interface()); err != nil {
			err := fmt.Errorf("%v on struct branch %s", name)
			return err
		}
	}
	return nil
}

func getMeta(fieldName string, tag reflect.StructTag) (name string, omit bool) {
	if tagStr := tag.Get("tree"); len(tagStr) > 0 {
		name, omit = tagSpec(tagStr)
	}
	if len(name) == 0 {
		name = fieldName
	} else if trimmed := strings.TrimSpace(name); len(trimmed) == 0 {
		name = fieldName
	}
	return
}

func valueTree(tree Tree, v interface{}) error {
	typ, val, err := checkType(v)
	if err != nil {
		return err
	}
	fields := typ.NumField()
	for i := 0; i < fields; i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)
		name, omit := getMeta(field.Name, field.Tag)
		if omit && isEmpty(&fieldValue) {
			continue
		}
		typ, val, isStruct := getValue(field.Type, &fieldValue)
		if !isStruct {
			tree.AddMetaNode(val.Interface(), name)
			continue
		} else if subNum := typ.NumField(); subNum == 0 {
			tree.AddMetaNode(val.Interface(), name)
			continue
		}
		branch := tree.AddBranch(name)
		if err := valueTree(branch, val.Interface()); err != nil {
			err := fmt.Errorf("%v on struct branch %s", name)
			return err
		}
	}
	return nil
}

func tagTree(tree Tree, v interface{}) error {
	return nil
}

func typeTree(tree Tree, v interface{}) error {
	return nil
}

func typeSizeTree(tree Tree, v interface{}) error {
	return nil
}

func getValue(typ reflect.Type, val *reflect.Value) (reflect.Type, *reflect.Value, bool) {
	switch typ.Kind() {
	case reflect.Ptr:
		typ = typ.Elem()
		if typ.Kind() == reflect.Struct {
			elem := val.Elem()
			return typ, &elem, true
		}
	case reflect.Struct:
		return typ, val, true
	}
	return typ, val, false
}

func checkType(v interface{}) (reflect.Type, *reflect.Value, error) {
	typ := reflect.TypeOf(v)
	val := reflect.ValueOf(v)
	switch typ.Kind() {
	case reflect.Ptr:
		typ = typ.Elem()
		if typ.Kind() != reflect.Struct {
			err := fmt.Errorf("treeprint: %T is not a struct we could work with", v)
			return nil, nil, err
		}
		val = val.Elem()
	case reflect.Struct:
	default:
		err := fmt.Errorf("treeprint: %T is not a struct we could work with", v)
		return nil, nil, err
	}
	return typ, &val, nil
}
