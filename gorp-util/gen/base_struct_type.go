package main

import (
	"go/ast"
	"reflect"

	gu "github.com/zkcrescent/chaos/gorp-util"
)

func BaseStructTypes() map[string]*ast.StructType {
	m := make(map[string]*ast.StructType)
	for _, i := range []interface{}{
		gu.Base{},
		gu.EnableBase{},
		gu.TaskBase{},
		gu.TaskEnableBase{},
		gu.FullTaskBase{},
	} {
		m = typeToStructType(reflect.TypeOf(i), m)
	}
	return m
}

func typeToStructType(ty reflect.Type, m map[string]*ast.StructType) map[string]*ast.StructType {
	if ty.Kind() != reflect.Struct {
		return m
	}
	if _, ok := m[ty.Name()]; ok {
		return m
	}

	st := &ast.StructType{
		Fields: &ast.FieldList{},
	}

	for i := 0; i < ty.NumField(); i++ {
		f := ty.Field(i)
		sf := &ast.Field{
			Tag: &ast.BasicLit{},
		}
		if v, ok := f.Tag.Lookup("db"); ok {
			sf.Names = append(sf.Names, &ast.Ident{
				Name: f.Name,
			})
			sf.Tag.Value = "`db:" + `"` + v + `"` + "`"
		} else {
			sf.Type = &ast.Ident{
				Name: f.Name,
			}
		}
		st.Fields.List = append(st.Fields.List, sf)
	}

	if st.Fields.List != nil {
		m[ty.Name()] = st
	}
	return m
}
