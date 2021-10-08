package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/build"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
	"text/template"

	"github.com/juju/errors"
	"github.com/serenize/snaker"
	"github.com/sirupsen/logrus"
)

// usage: gen
func main() {
	pkg, err := build.Default.ImportDir(".", 0)
	if err != nil {
		logrus.Fatalf("process directory failed: %s", err)
	}

	typeMap := BaseStructTypes()
	fs := token.NewFileSet()
	for _, filename := range pkg.GoFiles {
		if strings.HasSuffix(filename, "_gorp.go") {
			continue
		}
		if !strings.HasSuffix(filename, ".go") {
			continue
		}
		file, err := parser.ParseFile(fs, filename, nil, parser.ParseComments)
		if err != nil {
			logrus.Fatalf(errors.ErrorStack(err))
		}
		for _, decl := range file.Decls {
			genDecl, ok := decl.(*ast.GenDecl)
			if !ok {
				continue
			}

			for _, node := range genDecl.Specs {
				typeSpec, ok := node.(*ast.TypeSpec)
				if !ok {
					continue
				}

				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}
				typeMap[typeSpec.Name.Name] = structType
			}
		}
	}

	metas := make(map[string]*entityMeta)
	for _, filename := range pkg.GoFiles {
		if strings.HasSuffix(filename, "_gorp.go") {
			continue
		}
		if !strings.HasSuffix(filename, ".go") {
			continue
		}
		pkgName := pkg.Name
		//fs := token.NewFileSet()
		file, err := parser.ParseFile(fs, filename, nil, parser.ParseComments)
		if err != nil {
			logrus.Fatalf(errors.ErrorStack(err))
		}

		// ast.Inspect(file, traverses)
		cmap := ast.NewCommentMap(fs, file, file.Comments)
		for node, comments := range cmap {
			switch node.(type) {
			case *ast.GenDecl:
				for _, spec := range node.(*ast.GenDecl).Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}
					structType, ok := typeSpec.Type.(*ast.StructType)
					if !ok {
						continue
					}
					metas[typeSpec.Name.Name] = parseEntityMeta(pkgName, typeSpec.Name.Name, structType, comments, typeMap)
				}
			case *ast.TypeSpec:
				typeSpec := node.(*ast.TypeSpec)
				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok {
					continue
				}

				metas[typeSpec.Name.Name] = parseEntityMeta(pkgName, typeSpec.Name.Name, structType, comments, typeMap)
			}
		}
	}
	bs, _ := json.MarshalIndent(metas, "", "    ")
	logrus.Debugf(string(bs))
	// convert rel and mul Name, field to Table Name, field
	for _, meta := range metas {
		if len(meta.Rels) > 0 {
			for _, ref := range meta.Rels {
				if _, ok := metas[ref.Name]; !ok {
					logrus.Fatalf("%v ref not found: %v", meta.Name, ref.Name)
				}
				ref.TableName = metas[ref.Name].Table
			}
		}
		if len(meta.Muls) > 0 {
			for _, mul := range meta.Muls {
				if _, ok := metas[mul.Name]; !ok {
					logrus.Fatalf("%v mul not found: %v", meta.Name, mul.Name)
				}
				mul.TableName = metas[mul.Name].Table
			}
		}
	}
	bs, _ = json.MarshalIndent(metas, "", "    ")
	logrus.Debugf(string(bs))

	var isBool = func(rm interface{}) bool {
		_, ok := rm.(bool)
		return ok
	}

	var isInt64 = func(rm interface{}) bool {
		_, ok := rm.(int64)
		return ok
	}

	funcMap := template.FuncMap{
		"ToUpper": strings.ToUpper,
		"IsBool":  isBool,
		"IsInt64": isInt64,
	}

	tpl, err := template.New("gorp.tpl").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		panic(err)
	}
	// tpl := template.Must(template.ParseFiles(tempfile)).Funcs(template.FuncMap{
	// 	"ToUpper": strings.ToUpper,
	// })
	for _, meta := range metas {
		meta.generate(tpl)
	}
}

func parseEntityMeta(pkg string, name string, structType *ast.StructType, commentGroup []*ast.CommentGroup, typeMap map[string]*ast.StructType) *entityMeta {
	em := &entityMeta{
		Pkg:    pkg,
		Name:   name,
		Fields: FlatFields(structType, typeMap),
		Rels:   make(map[string]*Ref),
	}
	for _, comments := range commentGroup {
		for _, comment := range comments.List {
			annotation := newAnnotation(comment.Text)
			if annotation == nil {
				continue
			}
			em.Init = true
			switch annotation.Name {
			case "@TABLE":
				em.Table = annotation.Key
			case "@SHARD":
				em.ShardKey = annotation.Key
			case "@SHARDING":
				var err error
				em.Sharding, err = strconv.Atoi(annotation.Key)
				if err != nil {
					panic(err)
				}
			case "@PK":
				em.ID = annotation.Key
			case "@VER":
				em.Version = annotation.Key
			case "@REL":
				em.Rels[annotation.Key] = NewRef(annotation.Vals[0])
			case "@MUL":
				em.Muls = append(em.Muls, NewMul(annotation))
			case "@NOINIT":
				em.Init = false
			}
		}
	}
	if em.Version == "" && em.Fields["UpdatedSeq"] != "" {
		em.Version = em.Fields["UpdatedSeq"]
	}

	return em
}

type entityMeta struct {
	Pkg      string
	Name     string
	Fields   map[string]string
	Table    string
	ShardKey string
	Sharding int
	Init     bool
	ID       string
	Version  string
	Rels     map[string]*Ref
	Muls     []*Mul

	// for tpl
	Imports []string
}

func NewRef(s string) *Ref {
	tmp := strings.Split(s, ".")
	return &Ref{
		Name:  tmp[0],
		Field: tmp[1],
	}
}
func NewRefs(list []string) []*Ref {
	var refs []*Ref
	for _, s := range list {
		refs = append(refs, NewRef(s))
	}
	return refs
}

func NewMul(annotation *annotation) *Mul {
	tmp := strings.Split(annotation.Key, ",")
	refs := NewRefs(annotation.Vals)
	return &Mul{
		Edge:  tmp[0],
		Name:  tmp[1],
		Left:  refs[0],
		Right: refs[1],
	}
}

type Mul struct {
	Edge      string
	Name      string
	TableName string
	Left      *Ref
	Right     *Ref
}

type Ref struct {
	Name      string
	Field     string
	TableName string
}

func (e *entityMeta) generate(tpl *template.Template) {
	if len(e.Table)+len(e.ID)+len(e.Rels)+len(e.Muls) == 0 {
		return
	}

	if e.Table == "" {
		e.Table = snaker.SnakeToCamelLower(e.Name)
	}

	if e.ID == "" {
		e.ID = "ID"
	}
	e.Imports = []string{
		"encoding/json",
		"github.com/zkcrescent/chaos/gorp-util",
		"github.com/juju/errors",
		"gopkg.in/gorp.v2",
	}
	if e.Sharding > 0 {
		e.Imports = append([]string{"fmt"}, e.Imports...)
	}
	//if (e.Fields["CreatedTime"] + e.Fields["UpdatedTime"] + e.Fields["RemovedTime"]) != "" {
	//	e.Imports = append(e.Imports, "time")
	//}

	buf := new(bytes.Buffer)
	if err := tpl.Execute(buf, e); err != nil {
		logrus.Fatal(err)
	}

	//filename := fmt.Sprintf("%s_gorp.go", strings.ToLower(e.Name))
	//err := ioutil.WriteFile(filename, buf.Bytes(), 0644)
	//if err != nil {
	//	logrus.Fatalf(errors.ErrorStack(err))
	//}

	src, err := format.Source(buf.Bytes())
	if err != nil {
		logrus.Fatalf(errors.ErrorStack(err))
	}
	filename := fmt.Sprintf("%s_gorp.go", strings.ToLower(e.Name))
	logrus.Info(filename)
	err = ioutil.WriteFile(filename, src, 0644)
	if err != nil {
		logrus.Fatalf(errors.ErrorStack(err))
	}
}

func FlatFields(structType *ast.StructType, typeMap map[string]*ast.StructType) map[string]string {
	m := make(map[string]string)
	for _, f := range structType.Fields.List {
		if f.Tag != nil {
			tag, ok := reflect.StructTag(strings.ReplaceAll(f.Tag.Value, "`", "")).Lookup("db")
			if ok && tag != "-" {
				m[f.Names[0].Name] = strings.TrimSpace(strings.Split(tag, ",")[0])
			}
		}
	}
	for _, f := range structType.Fields.List {
		if len(f.Names) == 0 {
			var em string
			if selector, ok := f.Type.(*ast.SelectorExpr); ok {
				em = selector.Sel.Name
			} else if star, ok := f.Type.(*ast.StarExpr); ok {
				em = star.X.(*ast.Ident).Name
			} else {
				em = f.Type.(*ast.Ident).Name
			}
			emm := FlatFields(typeMap[em], typeMap)
			for k, v := range emm {
				if _, ok := m[k]; !ok {
					m[k] = v
				}
			}
		}
	}
	return m
}
