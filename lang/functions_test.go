package lang

import (
	"testing"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcl/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

// TestFunctions tests that functions are callable through the functionality
// in the langs package, via HCL.
//
// These tests are primarily here to assert that the functions are properly
// registered in the functions table, rather than to test all of the details
// of the functions. Each function should only have one or two tests here,
// since the main set of unit tests for a function should live alongside that
// function either in the "funcs" subdirectory here or over in the cty
// function/stdlib package.
//
// One exception to that is we can use this test mechanism to assert common
// patterns that are used in real-world configurations which rely on behaviors
// implemented either in this lang package or in HCL itself, such as automatic
// type conversions. The function unit tests don't cover those things because
// they call directly into the functions.
//
// With that said then, this test function should contain at least one simple
// test case per function registered in the functions table (just to prove
// it really is registered correctly) and possibly a small set of additional
// functions showing real-world use-cases that rely on type conversion
// behaviors.
func TestFunctions(t *testing.T) {
	tests := []struct {
		src  string
		want cty.Value
	}{
		// Please maintain this list in alphabetical order by function, with
		// a blank line between the group of tests for each function.

		{
			`abs(-1)`,
			cty.NumberIntVal(1),
		},

		{
			`ceil(1.2)`,
			cty.NumberIntVal(2),
		},

		{
			`chunklist(["a", "b", "c"], 1)`,
			cty.ListVal([]cty.Value{
				cty.ListVal([]cty.Value{
					cty.StringVal("a"),
				}),
				cty.ListVal([]cty.Value{
					cty.StringVal("b"),
				}),
				cty.ListVal([]cty.Value{
					cty.StringVal("c"),
				}),
			}),
		},

		{
			`cidrhost("192.168.1.0/24", 5)`,
			cty.StringVal("192.168.1.5"),
		},

		{
			`cidrnetmask("192.168.1.0/24")`,
			cty.StringVal("255.255.255.0"),
		},

		{
			`cidrsubnet("192.168.2.0/20", 4, 6)`,
			cty.StringVal("192.168.6.0/24"),
		},

		{
			`coalesce("first", "second", "third")`,
			cty.StringVal("first"),
		},

		{
			`coalescelist(["first", "second"], ["third", "fourth"])`,
			cty.ListVal([]cty.Value{
				cty.StringVal("first"), cty.StringVal("second"),
			}),
		},

		{
			`compact(["test", "", "test"])`,
			cty.ListVal([]cty.Value{
				cty.StringVal("test"), cty.StringVal("test"),
			}),
		},

		{
			`contains(["a", "b"], "a")`,
			cty.True,
		},
		{ // Should also work with sets, due to automatic conversion
			`contains(toset(["a", "b"]), "a")`,
			cty.True,
		},

		{
			`distinct(["a", "b", "a", "b"])`,
			cty.ListVal([]cty.Value{
				cty.StringVal("a"), cty.StringVal("b"),
			}),
		},

		{
			`element(["hello"], 0)`,
			cty.StringVal("hello"),
		},

		{
			`file("hello.txt")`,
			cty.StringVal("hello!"),
		},

		{
			`flatten([tolist(["a", "b"]), tolist(["c", "d"])])`,
			cty.ListVal([]cty.Value{
				cty.StringVal("a"),
				cty.StringVal("b"),
				cty.StringVal("c"),
				cty.StringVal("d"),
			}),
		},

		{
			`index(["a", "b", "c"], "a")`,
			cty.NumberIntVal(0),
		},

		{
			`keys({"hello"=1, "goodbye"=42})`,
			cty.TupleVal([]cty.Value{
				cty.StringVal("goodbye"),
				cty.StringVal("hello"),
			}),
		},

		{
			`length(["the", "quick", "brown", "bear"])`,
			cty.NumberIntVal(4),
		},

		{
			`list("hello")`,
			cty.ListVal([]cty.Value{
				cty.StringVal("hello"),
			}),
		},

		{
			`lookup({hello=1, goodbye=42}, "goodbye")`,
			cty.NumberIntVal(42),
		},

		{
			`map("hello", "world")`,
			cty.MapVal(map[string]cty.Value{
				"hello": cty.StringVal("world"),
			}),
		},

		{
			`matchkeys(["a", "b", "c"], ["ref1", "ref2", "ref3"], ["ref1"])`,
			cty.ListVal([]cty.Value{
				cty.StringVal("a"),
			}),
		},

		{
			`merge({"a"="b"}, {"c"="d"})`,
			cty.ObjectVal(map[string]cty.Value{
				"a": cty.StringVal("b"),
				"c": cty.StringVal("d"),
			}),
		},

		{
			`reverse(["a", true, 0])`,
			cty.TupleVal([]cty.Value{cty.Zero, cty.True, cty.StringVal("a")}),
		},

		{
			`setproduct(["development", "staging", "production"], ["app1", "app2"])`,
			cty.ListVal([]cty.Value{
				cty.TupleVal([]cty.Value{cty.StringVal("development"), cty.StringVal("app1")}),
				cty.TupleVal([]cty.Value{cty.StringVal("development"), cty.StringVal("app2")}),
				cty.TupleVal([]cty.Value{cty.StringVal("staging"), cty.StringVal("app1")}),
				cty.TupleVal([]cty.Value{cty.StringVal("staging"), cty.StringVal("app2")}),
				cty.TupleVal([]cty.Value{cty.StringVal("production"), cty.StringVal("app1")}),
				cty.TupleVal([]cty.Value{cty.StringVal("production"), cty.StringVal("app2")}),
			}),
		},

		{
			`slice(["a", "b", "c", "d"], 1, 3)`,
			cty.ListVal([]cty.Value{
				cty.StringVal("b"), cty.StringVal("c"),
			}),
		},

		{
			`transpose({"a" = ["1", "2"], "b" = ["2", "3"]})`,
			cty.MapVal(map[string]cty.Value{
				"1": cty.ListVal([]cty.Value{cty.StringVal("a")}),
				"2": cty.ListVal([]cty.Value{cty.StringVal("a"), cty.StringVal("b")}),
				"3": cty.ListVal([]cty.Value{cty.StringVal("b")}),
			}),
		},

		{
			`values({"hello"="world", "what's"="up"})`,
			cty.TupleVal([]cty.Value{
				cty.StringVal("world"),
				cty.StringVal("up"),
			}),
		},

		{
			`zipmap(["hello", "bar"], ["world", "baz"])`,
			cty.ObjectVal(map[string]cty.Value{
				"hello": cty.StringVal("world"),
				"bar":   cty.StringVal("baz"),
			}),
		},
	}

	for _, test := range tests {
		t.Run(test.src, func(t *testing.T) {
			expr, parseDiags := hclsyntax.ParseExpression([]byte(test.src), "test.hcl", hcl.Pos{Line: 1, Column: 1})
			if parseDiags.HasErrors() {
				for _, diag := range parseDiags {
					t.Error(diag.Error())
				}
				return
			}

			data := &dataForTests{} // no variables available; we only need literals here
			scope := &Scope{
				Data:    data,
				BaseDir: "./testdata/functions-test", // for the functions that read from the filesystem
			}

			got, diags := scope.EvalExpr(expr, cty.DynamicPseudoType)
			if diags.HasErrors() {
				for _, diag := range diags {
					t.Errorf("%s: %s", diag.Description().Summary, diag.Description().Detail)
				}
				return
			}

			if !test.want.RawEquals(got) {
				t.Errorf("wrong result\nexpr: %s\ngot:  %#v\nwant: %#v", test.src, got, test.want)
			}
		})
	}
}
