package badfuncs

import (
	"fmt"
	"go/ast"
	"path"
	"strconv"

	"golang.org/x/tools/go/analysis"
)

const Doc = `detect disallowed functions`

type BadFuncKey struct {
	pkg      string
	funcName string
}

var ReasonHttpMock = "httpmock.Register functions use a global variable which can cause random failures when unit tests from different packages are run in parallel. Instead, declare a local transport using trans=httpmock.NewMockTransport() and use trans.RegisterResponder()"

var badFuncs = map[BadFuncKey]string{}

var Analyzer = &analysis.Analyzer{
	Name:             "badfuncs",
	Doc:              Doc,
	Run:              run,
	RunDespiteErrors: true,
	//FactTypes:        []analysis.Fact{new(foundFact)},
}

func addBadFunc(pkg, name, reason string) {
	key := BadFuncKey{
		pkg:      pkg,
		funcName: name,
	}
	badFuncs[key] = reason
}

func run(pass *analysis.Pass) (interface{}, error) {
	addBadFunc("httpmock", "RegisterResponder", ReasonHttpMock)
	addBadFunc("httpmock", "RegisterRegexpResponder", ReasonHttpMock)
	addBadFunc("httpmock", "RegisterResponderWithQuery", ReasonHttpMock)
	addBadFunc("httpmock", "RegisterNoResponder", ReasonHttpMock)
	addBadFunc("httpmock", "DeactivateAndReset", ReasonHttpMock)

	for _, f := range pass.Files {
		imports := make(map[string]string)
		ast.Inspect(f, func(n ast.Node) bool {
			switch stmt := n.(type) {
			case *ast.ImportSpec:
				// collect import aliases
				pkg, err := strconv.Unquote(stmt.Path.Value)
				if err != nil {
					pkg = stmt.Path.Value
				}
				name := path.Base(pkg)
				if stmt.Name != nil {
					name = stmt.Name.Name
				}
				imports[name] = pkg
			case *ast.CallExpr:
				sel, ok := stmt.Fun.(*ast.SelectorExpr)
				if !ok {
					return true
				}
				key := BadFuncKey{
					pkg:      getIdent(sel.X),
					funcName: getIdent(sel.Sel),
				}
				if _, found := badFuncs[key]; found {
					pass.Reportf(n.Pos(), fmt.Sprintf("Unsafe func %s.%s", key.pkg, key.funcName))
				}
			}
			return true
		})
	}
	return nil, nil
}

func getIdent(expr ast.Expr) string {
	id, ok := expr.(*ast.Ident)
	if !ok {
		return ""
	}
	return id.Name
}
