package templating

import (
	"sort"
	"testing"

	sprig "github.com/Masterminds/sprig/v3"
)

var reviewedDeniedSprigFunctionNames = map[string]string{
	"ago":                        "time-dependent helper",
	"bcrypt":                     "credential/hash helper not needed in templates",
	"buildCustomCert":            "certificate material helper",
	"chunk":                      "not needed in supported template surface",
	"deepCopy":                   "object graph helper not needed in templates",
	"decryptAES":                 "cryptographic helper",
	"derivePassword":             "credential derivation helper",
	"encryptAES":                 "cryptographic helper",
	"env":                        "host environment access",
	"expandenv":                  "host environment access",
	"genCA":                      "certificate material helper",
	"genCAWithKey":               "certificate material helper",
	"genPrivateKey":              "key material helper",
	"genSelfSignedCert":          "certificate material helper",
	"genSelfSignedCertWithKey":   "certificate material helper",
	"genSignedCert":              "certificate material helper",
	"genSignedCertWithKey":       "certificate material helper",
	"getHostByName":              "network lookup helper",
	"hello":                      "diagnostic helper",
	"htpasswd":                   "credential/hash helper not needed in templates",
	"merge":                      "map merge helper can hide unexpected shape changes",
	"mergeOverwrite":             "map merge helper can hide unexpected shape changes",
	"mustChunk":                  "not needed in supported template surface",
	"mustDeepCopy":               "object graph helper not needed in templates",
	"mustMerge":                  "map merge helper can hide unexpected shape changes",
	"mustMergeOverwrite":         "map merge helper can hide unexpected shape changes",
	"mustRegexFind":              "regex helpers are not part of supported template surface",
	"mustRegexFindAll":           "regex helpers are not part of supported template surface",
	"mustRegexMatch":             "regex helpers are not part of supported template surface",
	"mustRegexReplaceAll":        "regex helpers are not part of supported template surface",
	"mustRegexReplaceAllLiteral": "regex helpers are not part of supported template surface",
	"mustRegexSplit":             "regex helpers are not part of supported template surface",
	"randBytes":                  "random helper is not part of current random opt-in surface",
	"regexFind":                  "regex helpers are not part of supported template surface",
	"regexFindAll":               "regex helpers are not part of supported template surface",
	"regexMatch":                 "regex helpers are not part of supported template surface",
	"regexQuoteMeta":             "regex helpers are not part of supported template surface",
	"regexReplaceAll":            "regex helpers are not part of supported template surface",
	"regexReplaceAllLiteral":     "regex helpers are not part of supported template surface",
	"regexSplit":                 "regex helpers are not part of supported template surface",
	"repeat":                     "unbounded output helper",
	"seq":                        "unbounded sequence helper",
	"set":                        "mutating map helper",
	"sha512sum":                  "not needed in supported template surface",
	"shuffle":                    "randomized helper",
	"swapcase":                   "not needed in supported template surface",
	"unset":                      "mutating map helper",
	"until":                      "unbounded sequence helper",
	"untilStep":                  "unbounded sequence helper",
	"urlJoin":                    "URL helper not needed in templates",
	"urlParse":                   "URL helper not needed in templates",
}

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

func TestAllSprigFunctionsExplicitlyReviewed(t *testing.T) {
	reviewed := map[string]string{}
	for _, name := range safeSprigFunctionNames {
		reviewed[name] = "safe"
	}
	for _, name := range randomSprigFunctionNames {
		reviewed[name] = "random opt-in"
	}
	reviewed["now"] = "time opt-in"
	for name, reason := range reviewedDeniedSprigFunctionNames {
		reviewed[name] = "denied: " + reason
	}

	var unreviewed []string
	for name := range sprig.TxtFuncMap() {
		if _, ok := reviewed[name]; !ok {
			unreviewed = append(unreviewed, name)
		}
	}
	sort.Strings(unreviewed)
	for _, name := range unreviewed {
		t.Errorf(
			"unreviewed Sprig function %q: add to safeSprigFunctionNames, "+
				"randomSprigFunctionNames, or reviewedDeniedSprigFunctionNames",
			name,
		)
	}
}

func TestSafeSprigOSPathFunctionsAreStringOnlyAndAllowed(t *testing.T) {
	funcs := buildFuncMap(false, false)

	for _, name := range []string{"osBase", "osClean", "osDir", "osExt", "osIsAbs"} {
		if _, ok := funcs[name]; !ok {
			t.Fatalf("expected reviewed string-only path helper %s to remain available", name)
		}
	}
}
