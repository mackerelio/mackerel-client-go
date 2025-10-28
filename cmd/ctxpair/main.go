package main

import (
	"cmp"
	"go/ast"
	"maps"
	"slices"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/unitchecker"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name: "ctxpair",
	Doc:  "ctxpair detects exported functions that have both context.Context and non-context.Context versions.",
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func main() {
	unitchecker.Main(Analyzer)
}

func isClientMethod(fn *ast.FuncDecl) bool {
	if fn.Recv == nil || len(fn.Recv.List) == 0 {
		return false
	}
	recvType := fn.Recv.List[0].Type
	if starExpr, ok := recvType.(*ast.StarExpr); ok {
		recvType = starExpr.X
	}
	if ident, ok := recvType.(*ast.Ident); ok {
		return ident.Name == "Client"
	}
	return false
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	funcs := make(map[string]*ast.FuncDecl)
	inspect.Preorder(nodeFilter, func(n ast.Node) {
		fn, ok := n.(*ast.FuncDecl)
		if ok && isClientMethod(fn) && ast.IsExported(fn.Name.Name) {
			funcs[fn.Name.Name] = fn
		}
	})
	stems := slices.CompactFunc(slices.Sorted(maps.Keys(funcs)), func(s, t string) bool {
		return strings.TrimSuffix(s, "Context") == strings.TrimSuffix(t, "Context")
	})

	for _, stem := range stems {
		fn1, ok1 := funcs[stem]
		fn2, ok2 := funcs[stem+"Context"]
		if !ok1 || !ok2 {
			fn := cmp.Or(fn1, fn2)
			pass.Reportf(fn.Pos(), "exported method %[1]s must have both (%[1]s and %[1]sContext)", stem)
		}
	}
	return nil, nil
}
