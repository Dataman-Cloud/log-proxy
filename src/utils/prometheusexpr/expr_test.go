package prometheusexpr

import "testing"

func TestNewExpr(t *testing.T) {
	if NewExpr() == nil {
		t.Error("faild")
	} else {
		t.Log("success")
	}
}

func TestParseExpr(t *testing.T) {
	path := "./testfiles/"
	err := Exprs(path)
	if err == nil {
		t.Error("failed")
	} else {
		t.Log("success to get the expr files")
	}
}

func TestGetExprFiles(t *testing.T) {
	path := "./testfile/"
	err := Exprs(path)
	if err == nil {
		t.Error("failed")
	} else {
		t.Log("success to get the expr files")
	}
}

func TestExprs(t *testing.T) {
	path := "../../../config/exprs/"
	err := Exprs(path)
	if err != nil {
		t.Error("failed")
	} else {
		t.Log("success to get the expr files")
	}
}

func TestGetExprs(t *testing.T) {
	path := "../../../config/exprs/"
	err := Exprs(path)
	if err != nil {
		t.Error("failed")
	} else {
		t.Log("success to get the expr files")
	}
	exprsList := GetExprs()
	if exprsList == nil {
		t.Error("failed")
	} else {
		t.Log("success to get the expr files")
	}
}
