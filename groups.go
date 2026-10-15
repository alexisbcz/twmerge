package twmerge

import (
	"strings"
)

// classify returns the conflict group of a utility (without variants), or ""
// when the class is not recognized.
func classify(base string) string {
	base = strings.TrimPrefix(base, "-") // negative values: -mt-2
	if base == "" {
		return ""
	}
	// Arbitrary property: [mask-type:luminance] conflicts with the same
	// property only.
	if base[0] == '[' && strings.HasSuffix(base, "]") {
		if i := strings.IndexByte(base, ':'); i > 1 {
			return "arbitrary.." + base[1:i]
		}
		return ""
	}
	if g, ok := exact[base]; ok {
		return g
	}
	// Try the longest prefix first: "border-x-2" is "border-x" + "2",
	// not "border" + "x-2".
	for i := len(base) - 1; i > 0; i-- {
		if base[i] != '-' {
			continue
		}
		if r, ok := prefixes[base[:i]]; ok {
			if g := r(base[i+1:]); g != "" {
				return g
			}
		}
	}
	return ""
}

// resolver maps the value part of a utility to its group, or "" if the
// value does not fit the utility.
type resolver func(v string) string

// always returns a resolver accepting any value.
func always(group string) resolver { return func(string) string { return group } }

// choose returns a resolver that picks group a for values in set and b (if
// non-empty) for everything else.
func choose(set []string, a, b string) resolver {
	m := make(map[string]bool, len(set))
	for _, s := range set {
		m[s] = true
	}
	return func(v string) string {
		if m[stripPostfix(v)] {
			return a
		}
		return b
	}
}

// widthOrColor picks group w for numeric or length values and c otherwise.
func widthOrColor(w, c string) resolver {
	return func(v string) string {
		if isNumber(v) || isArbitraryLength(v) {
			return w
		}
		return c
	}
}

var (
	tshirt       = []string{"3xs", "2xs", "xs", "sm", "md", "lg", "xl", "2xl", "3xl", "4xl", "5xl", "6xl", "7xl", "8xl", "9xl"}
	fontSizes    = append([]string{"base"}, tshirt...)
	borderStyles = []string{"solid", "dashed", "dotted", "double", "hidden", "none"}
	sides        = []string{"x", "y", "t", "r", "b", "l", "s", "e"}
)

var prefixes = map[string]resolver{}

var exact = map[string]string{}

// conflicts lists, for a group, the groups it overrides: p-2 overrides an
// earlier px-4, but px-4 does not override an earlier p-2.
var conflicts = map[string][]string{
	"p":              {"px", "py", "ps", "pe", "pt", "pr", "pb", "pl"},
	"px":             {"pr", "pl"},
	"py":             {"pt", "pb"},
	"m":              {"mx", "my", "ms", "me", "mt", "mr", "mb", "ml"},
	"mx":             {"mr", "ml"},
	"my":             {"mt", "mb"},
	"size":           {"w", "h"},
	"inset":          {"inset-x", "inset-y", "start", "end", "top", "right", "bottom", "left"},
	"inset-x":        {"right", "left"},
	"inset-y":        {"top", "bottom"},
	"gap":            {"gap-x", "gap-y"},
	"overflow":       {"overflow-x", "overflow-y"},
	"overscroll":     {"overscroll-x", "overscroll-y"},
	"rounded":        {"rounded-s", "rounded-e", "rounded-t", "rounded-r", "rounded-b", "rounded-l", "rounded-ss", "rounded-se", "rounded-ee", "rounded-es", "rounded-tl", "rounded-tr", "rounded-br", "rounded-bl"},
	"rounded-s":      {"rounded-ss", "rounded-es"},
	"rounded-e":      {"rounded-se", "rounded-ee"},
	"rounded-t":      {"rounded-tl", "rounded-tr"},
	"rounded-r":      {"rounded-tr", "rounded-br"},
	"rounded-b":      {"rounded-br", "rounded-bl"},
	"rounded-l":      {"rounded-tl", "rounded-bl"},
	"border-w":       {"border-w-x", "border-w-y", "border-w-t", "border-w-r", "border-w-b", "border-w-l", "border-w-s", "border-w-e"},
	"border-w-x":     {"border-w-r", "border-w-l"},
	"border-w-y":     {"border-w-t", "border-w-b"},
	"border-color":   {"border-color-x", "border-color-y", "border-color-t", "border-color-r", "border-color-b", "border-color-l", "border-color-s", "border-color-e"},
	"border-color-x": {"border-color-r", "border-color-l"},
	"border-color-y": {"border-color-t", "border-color-b"},
	"font-size":      {"leading"},
	"translate":      {"translate-x", "translate-y"},
	"scale":          {"scale-x", "scale-y"},
	"line-clamp":     {"display", "overflow"},
}

func init() {
	group := func(g string, classes ...string) {
		for _, c := range classes {
			exact[c] = g
		}
	}
	group("display", "block", "inline-block", "inline", "flex", "inline-flex", "table", "inline-table",
		"table-caption", "table-cell", "table-column", "table-column-group", "table-footer-group",
		"table-header-group", "table-row-group", "table-row", "flow-root", "grid", "inline-grid",
		"contents", "list-item", "hidden")
	group("position", "static", "fixed", "absolute", "relative", "sticky")
	group("visibility", "visible", "invisible", "collapse")
	group("isolation", "isolate", "isolation-auto")
	group("text-decoration-line", "underline", "overline", "line-through", "no-underline")
	group("text-transform", "uppercase", "lowercase", "capitalize", "normal-case")
	group("font-style", "italic", "not-italic")
	group("font-smoothing", "antialiased", "subpixel-antialiased")
	group("text-overflow", "truncate")
	group("sr", "sr-only", "not-sr-only")
	group("box-sizing", "box-border", "box-content")
	group("border-collapse", "border-collapse", "border-separate")
	group("container", "container")
	group("grow", "grow")
	group("shrink", "shrink")
	group("border-w", "border")
	group("rounded", "rounded")
	group("shadow", "shadow")
	group("ring-w", "ring")
	group("ring-inset", "ring-inset")
	group("outline-w", "outline")
	group("transition", "transition")
	group("resize", "resize")
	group("space-x-reverse", "space-x-reverse")
	group("space-y-reverse", "space-y-reverse")
	group("divide-x", "divide-x")
	group("divide-y", "divide-y")
	group("blur", "blur")
	for _, s := range sides {
		group("border-w-"+s, "border-"+s)
	}
	for _, r := range []string{"s", "e", "t", "r", "b", "l", "ss", "se", "ee", "es", "tl", "tr", "br", "bl"} {
		group("rounded-"+r, "rounded-"+r)
		prefixes["rounded-"+r] = always("rounded-" + r)
	}

	anyValue := func(groups ...string) {
		for _, g := range groups {
			prefixes[g] = always(g)
		}
	}
	anyValue("p", "px", "py", "ps", "pe", "pt", "pr", "pb", "pl",
		"m", "mx", "my", "ms", "me", "mt", "mr", "mb", "ml",
		"gap", "gap-x", "gap-y", "space-x", "space-y",
		"w", "min-w", "max-w", "h", "min-h", "max-h", "size",
		"inset", "inset-x", "inset-y", "top", "right", "bottom", "left", "start", "end",
		"z", "order", "basis", "grow", "shrink",
		"grid-cols", "grid-rows", "col", "row", "col-span", "col-start", "col-end",
		"row-span", "row-start", "row-end", "auto-cols", "auto-rows", "grid-flow",
		"justify-items", "justify-self", "items", "self", "place-content", "place-items", "place-self",
		"rounded", "leading", "tracking", "line-clamp", "opacity", "cursor", "pointer-events",
		"select", "resize", "overflow", "overflow-x", "overflow-y", "overscroll", "overscroll-x",
		"overscroll-y", "whitespace", "break-before", "break-after", "break-inside",
		"underline-offset", "duration", "ease", "delay", "animate", "aspect", "translate",
		"translate-x", "translate-y", "scale", "scale-x", "scale-y", "rotate", "skew-x", "skew-y",
		"origin", "fill", "accent", "caret", "blur", "backdrop-blur", "brightness", "contrast",
		"grayscale", "invert", "saturate", "sepia", "drop-shadow", "will-change", "appearance",
		"columns", "float", "clear", "align", "indent", "outline-offset", "ring-offset-w",
		"border-spacing", "bg-blend", "mix-blend", "bg-clip", "bg-origin", "scroll-m", "scroll-p",
		"snap", "touch", "content-visibility", "field-sizing", "inset-shadow", "inset-ring",
		"text-shadow", "mask", "perspective", "backface", "transform-style", "box-decoration",
		"table", "caption", "list-image", "hyphens", "object-position")

	prefixes["flex"] = func(v string) string {
		switch v {
		case "row", "row-reverse", "col", "col-reverse":
			return "flex-direction"
		case "wrap", "wrap-reverse", "nowrap":
			return "flex-wrap"
		}
		return "flex"
	}
	prefixes["justify"] = choose([]string{"start", "end", "center", "between", "around", "evenly", "stretch", "baseline", "normal", "center-safe", "end-safe"}, "justify-content", "")
	prefixes["content"] = choose([]string{"normal", "center", "start", "end", "between", "around", "evenly", "baseline", "stretch", "center-safe", "end-safe"}, "align-content", "content")

	prefixes["text"] = func(v string) string {
		switch stripPostfix(v) {
		case "left", "center", "right", "justify", "start", "end":
			return "text-align"
		case "ellipsis", "clip":
			return "text-overflow"
		case "wrap", "nowrap", "balance", "pretty":
			return "text-wrap"
		}
		if contains(fontSizes, stripPostfix(v)) || isArbitraryLength(v) {
			return "font-size"
		}
		return "text-color"
	}
	prefixes["font"] = func(v string) string {
		switch v {
		case "thin", "extralight", "light", "normal", "medium", "semibold", "bold", "extrabold", "black":
			return "font-weight"
		}
		if isNumber(v) || isArbitraryNumber(v) {
			return "font-weight"
		}
		if strings.HasPrefix(v, "stretch-") {
			return "font-stretch"
		}
		return "font-family"
	}
	prefixes["list"] = func(v string) string {
		if v == "inside" || v == "outside" {
			return "list-position"
		}
		return "list-style-type"
	}
	prefixes["break"] = choose([]string{"normal", "words", "all", "keep"}, "word-break", "")

	prefixes["border"] = func(v string) string {
		switch {
		case contains(borderStyles, v):
			return "border-style"
		case isNumber(v) || isArbitraryLength(v):
			return "border-w"
		}
		return "border-color"
	}
	for _, s := range sides {
		prefixes["border-"+s] = widthOrColor("border-w-"+s, "border-color-"+s)
	}
	prefixes["divide"] = choose(borderStyles, "divide-style", "divide-color")
	prefixes["divide-x"] = always("divide-x")
	prefixes["divide-y"] = always("divide-y")

	prefixes["bg"] = func(v string) string {
		switch v {
		case "fixed", "local", "scroll":
			return "bg-attachment"
		case "repeat", "no-repeat", "repeat-x", "repeat-y", "repeat-round", "repeat-space":
			return "bg-repeat"
		case "auto", "cover", "contain":
			return "bg-size"
		case "bottom", "center", "left", "left-bottom", "left-top", "right", "right-bottom",
			"right-top", "top", "top-left", "top-right", "bottom-left", "bottom-right":
			return "bg-position"
		case "none":
			return "bg-image"
		}
		for _, p := range []string{"linear-", "radial", "conic", "gradient-"} {
			if strings.HasPrefix(v, p) {
				return "bg-image"
			}
		}
		if strings.HasPrefix(v, "[url(") || strings.HasPrefix(v, "[image:") {
			return "bg-image"
		}
		return "bg-color"
	}
	for _, stop := range []string{"from", "via", "to"} {
		prefixes[stop] = func(v string) string {
			if strings.HasSuffix(v, "%") || isArbitraryLength(v) {
				return stop + "-position"
			}
			return stop + "-color"
		}
	}

	prefixes["shadow"] = choose(append([]string{"none", "inner"}, tshirt...), "shadow", "shadow-color")
	prefixes["ring"] = widthOrColor("ring-w", "ring-color")
	prefixes["ring-offset"] = widthOrColor("ring-offset-w", "ring-offset-color")
	prefixes["outline"] = func(v string) string {
		switch {
		case contains(borderStyles, v):
			return "outline-style"
		case isNumber(v) || isArbitraryLength(v):
			return "outline-w"
		}
		return "outline-color"
	}
	prefixes["decoration"] = func(v string) string {
		switch {
		case v == "solid" || v == "double" || v == "dotted" || v == "dashed" || v == "wavy":
			return "decoration-style"
		case v == "auto" || v == "from-font" || isNumber(v) || isArbitraryLength(v):
			return "decoration-thickness"
		}
		return "decoration-color"
	}
	prefixes["stroke"] = widthOrColor("stroke-w", "stroke-color")
	prefixes["transition"] = always("transition")
	prefixes["object"] = choose([]string{"contain", "cover", "fill", "none", "scale-down"}, "object-fit", "object-position")
	prefixes["box"] = choose([]string{"border", "content"}, "box-sizing", "")
}

// --- value helpers -------------------------------------------------------------

func contains(set []string, v string) bool {
	for _, s := range set {
		if s == v {
			return true
		}
	}
	return false
}

// stripPostfix removes a "/modifier" (opacity or line height) outside of
// brackets: "sm/6" -> "sm".
func stripPostfix(v string) string {
	if strings.HasPrefix(v, "[") || strings.HasPrefix(v, "(") {
		return v
	}
	if i := strings.IndexByte(v, '/'); i > 0 {
		return v[:i]
	}
	return v
}

func isNumber(v string) bool {
	if v == "" {
		return false
	}
	dot := false
	for i := 0; i < len(v); i++ {
		switch c := v[i]; {
		case c >= '0' && c <= '9':
		case c == '.' && !dot:
			dot = true
		default:
			return false
		}
	}
	return true
}

// arbitrary returns the inside of "[...]" or "(...)" and its optional type
// label ("length" in "[length:var(--x)]").
func arbitrary(v string) (label, inner string, ok bool) {
	switch {
	case len(v) >= 2 && v[0] == '[' && v[len(v)-1] == ']':
	case len(v) >= 2 && v[0] == '(' && v[len(v)-1] == ')':
	default:
		return "", "", false
	}
	inner = v[1 : len(v)-1]
	if i := strings.IndexByte(inner, ':'); i > 0 && isIdent(inner[:i]) {
		return inner[:i], inner[i+1:], true
	}
	return "", inner, true
}

func isIdent(s string) bool {
	for i := 0; i < len(s); i++ {
		if c := s[i]; !(c >= 'a' && c <= 'z' || c == '-') {
			return false
		}
	}
	return true
}

var lengthUnits = []string{"px", "rem", "em", "%", "vh", "vw", "dvh", "dvw", "svh", "svw", "lvh", "lvw",
	"vmin", "vmax", "ch", "ex", "lh", "rlh", "cm", "mm", "in", "pt", "pc", "cqw", "cqh", "cqi", "cqb", "cqmin", "cqmax"}

func isArbitraryLength(v string) bool {
	label, inner, ok := arbitrary(v)
	if !ok {
		return false
	}
	if label != "" {
		return label == "length" || label == "size" || label == "percentage"
	}
	if inner == "0" {
		return true
	}
	for _, fn := range []string{"calc(", "min(", "max(", "clamp("} {
		if strings.HasPrefix(inner, fn) {
			return true
		}
	}
	for _, u := range lengthUnits {
		if strings.HasSuffix(inner, u) && isNumber(strings.TrimSuffix(inner, u)) {
			return true
		}
	}
	return false
}

func isArbitraryNumber(v string) bool {
	label, inner, ok := arbitrary(v)
	return ok && (label == "number" || label == "" && isNumber(inner))
}
