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

func TestParseExprs(t *testing.T) {
	path := "../../../config/exprs"
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

func TestSetFilterAndByWithPrefixFilter(t *testing.T) {
	expr := NewExpr()
	expr.setFilterAndByWithPrefix("DM")
	cluster := "container_label_DM_VCLUSTER='%s'"
	user := "container_label_DM_USER='%s'"
	app := "container_label_DM_APP_NAME='%s'"
	slot := "container_label_DM_SLOT_ID='%s'"
	task := "container_label_DM_TASK_ID='%s'"
	if expr.Filter.Cluster != cluster {
		t.Errorf("expect %s, got %s", cluster, expr.Filter.Cluster)
	}
	if expr.Filter.User != user {
		t.Errorf("expect %s, got %s", user, expr.Filter.User)
	}
	if expr.Filter.App != app {
		t.Errorf("expect %s, got %s", app, expr.Filter.App)
	}
	if expr.Filter.Slot != slot {
		t.Errorf("expect %s, got %s", slot, expr.Filter.Slot)
	}
	if expr.Filter.Task != task {
		t.Errorf("expect %s, got %s", task, expr.Filter.Task)
	}
}

func TestSetFilterAndByWithPrefixBy(t *testing.T) {
	expr := NewExpr()
	expr.setFilterAndByWithPrefix("DM")
	cluster := "container_label_DM_VCLUSTER"
	user := "container_label_DM_USER"
	app := "container_label_DM_APP_NAME"
	slot := "container_label_DM_SLOT_ID"
	task := "container_label_DM_TASK_ID"
	if expr.By.Cluster != cluster {
		t.Errorf("expect %s, got %s", cluster, expr.By.Cluster)
	}
	if expr.By.User != user {
		t.Errorf("expect %s, got %s", user, expr.By.User)
	}
	if expr.By.App != app {
		t.Errorf("expect %s, got %s", app, expr.By.App)
	}
	if expr.By.Slot != slot {
		t.Errorf("expect %s, got %s", slot, expr.By.Slot)
	}
	if expr.By.Task != task {
		t.Errorf("expect %s, got %s", task, expr.By.Task)
	}
}
