package templating

import "testing"

func TestBuildFuncMapRestrictsDangerousSprigFunctions(t *testing.T) {
	funcs := buildFuncMap(false, false)

	for _, name := range []string{"env", "expandenv", "repeat", "seq", "until", "untilStep", "getHostByName"} {
		if _, ok := funcs[name]; ok {
			t.Fatalf("expected %s to be excluded from func map", name)
		}
	}
	if _, ok := funcs["now"]; ok {
		t.Fatalf("expected now to be excluded when deterministic helpers are disabled")
	}
	if _, ok := funcs["uuidv4"]; ok {
		t.Fatalf("expected random helpers to be excluded when disabled")
	}
}

func TestBuildFuncMapAllowsDeterministicOptIns(t *testing.T) {
	funcs := buildFuncMap(true, true)

	if _, ok := funcs["now"]; !ok {
		t.Fatalf("expected now helper to be present")
	}
	if _, ok := funcs["uuidv4"]; !ok {
		t.Fatalf("expected random helper to be present")
	}
	if _, ok := funcs["toJson"]; !ok {
		t.Fatalf("expected safe sprig helper to remain available")
	}
}
