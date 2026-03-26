package templating

import (
	"text/template"

	sprig "github.com/Masterminds/sprig/v3"
)

var safeSprigFunctionNames = []string{
	"date", "date_in_zone", "date_modify", "dateInZone", "dateModify", "duration", "durationRound",
	"htmlDate", "htmlDateInZone", "must_date_modify", "mustDateModify", "mustToDate", "toDate", "unixEpoch",
	"abbrev", "abbrevboth", "trunc", "trim", "upper", "lower", "title", "untitle", "substr", "trimall", "trimAll",
	"trimSuffix", "trimPrefix", "nospace", "initials", "snakecase", "camelcase", "kebabcase", "wrap", "wrapWith",
	"contains", "hasPrefix", "hasSuffix", "quote", "squote", "cat", "indent", "nindent", "replace", "plural",
	"sha1sum", "sha256sum", "adler32sum", "toString",
	"atoi", "int64", "int", "float64", "toDecimal",
	"add1", "add", "sub", "div", "mod", "mul", "add1f", "addf", "subf", "divf", "mulf",
	"biggest", "max", "min", "maxf", "minf", "ceil", "floor", "round",
	"split", "splitList", "splitn", "toStrings", "join", "sortAlpha",
	"default", "empty", "coalesce", "all", "any", "compact", "mustCompact",
	"fromJson", "toJson", "toPrettyJson", "toRawJson", "mustFromJson", "mustToJson", "mustToPrettyJson", "mustToRawJson",
	"ternary",
	"typeOf", "typeIs", "typeIsLike", "kindOf", "kindIs", "deepEqual",
	"base", "dir", "clean", "ext", "isAbs", "osBase", "osClean", "osDir", "osExt", "osIsAbs",
	"b64enc", "b64dec", "b32enc", "b32dec",
	"tuple", "list", "dict", "get", "hasKey", "pluck", "keys", "pick", "omit", "values",
	"append", "push", "mustAppend", "mustPush", "prepend", "mustPrepend",
	"first", "mustFirst", "rest", "mustRest", "last", "mustLast", "initial", "mustInitial",
	"reverse", "mustReverse", "uniq", "mustUniq", "without", "mustWithout", "has", "mustHas",
	"slice", "mustSlice", "concat", "dig",
	"semver", "semverCompare",
	"fail",
}

var randomSprigFunctionNames = []string{
	"randAlpha",
	"randAlphaNum",
	"randAscii",
	"randNumeric",
	"randInt",
	"uuidv4",
}

func buildFuncMap(allowNow bool, allowRandom bool) template.FuncMap {
	all := sprig.TxtFuncMap()
	funcs := make(template.FuncMap, len(safeSprigFunctionNames)+8)
	for _, name := range safeSprigFunctionNames {
		if fn, ok := all[name]; ok {
			funcs[name] = fn
		}
	}

	if allowNow {
		if fn, ok := all["now"]; ok {
			funcs["now"] = fn
		}
	}
	if allowRandom {
		for _, name := range randomSprigFunctionNames {
			if fn, ok := all[name]; ok {
				funcs[name] = fn
			}
		}
	}

	funcs["len"] = func(value any) (int, error) {
		return lenValue(value)
	}
	funcs["hash_of"] = func(value any) (string, error) {
		return hashOfValue(value)
	}
	funcs["type_of"] = func(value any) string {
		return typeOfValue(value)
	}
	funcs["exists"] = func(value any) bool {
		return value != nil
	}
	funcs["sample"] = func(value any) (any, error) {
		return sampleValue(value)
	}

	return funcs
}
