// Package twmerge merges Tailwind CSS class lists so that later classes
// override conflicting earlier ones, like the tailwind-merge JavaScript
// library that shadcn/ui's cn() is built on:
//
//	Merge("px-2 py-1 bg-red-500", "p-3 bg-blue-500") // "p-3 bg-blue-500"
//	Merge("rounded-md", "rounded-full")              // "rounded-full"
//	Merge("hover:bg-accent", "bg-primary")           // both kept: different variants
//
// It understands the Tailwind v4 default utilities: variants (hover:,
// dark:, has-[>svg]:, [&_svg]:, ...), the important modifier (!p-2 or
// p-2!), negative values, arbitrary values ([3px], (--var)) and
// arbitrary properties ([mask-type:luminance]). Classes it does not
// recognize are always kept.
package twmerge

import (
	"slices"
	"strings"
	"sync"
)

// Merge joins the class lists and removes classes overridden by a later
// class in the same group with the same variants.
func Merge(classLists ...string) string {
	// The cache key keeps the first two lists apart instead of joining them,
	// so the common cn(base, className) call allocates nothing on a hit.
	var k cacheKey
	switch len(classLists) {
	case 0:
		return ""
	case 1:
		k = cacheKey{classLists[0], ""}
	case 2:
		k = cacheKey{classLists[0], classLists[1]}
	default:
		k = cacheKey{strings.Join(classLists, " "), ""}
	}
	if v, ok := cache.get(k); ok {
		return v
	}
	out := merge(k.a + " " + k.b)
	cache.put(k, out)
	return out
}

type cacheKey struct{ a, b string }

func merge(classes string) string {
	fields := strings.Fields(classes)
	seen := make(map[string]struct{}, len(fields))
	keep := make([]bool, len(fields))
	for i := len(fields) - 1; i >= 0; i-- {
		cls := fields[i]
		mods, important, base := parse(cls)
		group := classify(base)
		if group == "" {
			keep[i] = true
			continue
		}
		prefix := modifierKey(mods)
		if important {
			prefix += "!"
		}
		key := prefix + group
		if _, dup := seen[key]; dup {
			continue
		}
		keep[i] = true
		seen[key] = struct{}{}
		for _, c := range conflicts[group] {
			seen[prefix+c] = struct{}{}
		}
	}
	var b strings.Builder
	for i, cls := range fields {
		if !keep[i] {
			continue
		}
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(cls)
	}
	return b.String()
}

// parse splits "dark:hover:!bg-red-500" into its variants, the important
// flag and the base utility. Colons inside brackets or parentheses are not
// separators.
func parse(cls string) (mods []string, important bool, base string) {
	depth, start := 0, 0
	for i := 0; i < len(cls); i++ {
		switch cls[i] {
		case '[', '(':
			depth++
		case ']', ')':
			depth--
		case ':':
			if depth == 0 {
				mods = append(mods, cls[start:i])
				start = i + 1
			}
		}
	}
	base = cls[start:]
	if strings.HasPrefix(base, "!") {
		important, base = true, base[1:]
	} else if strings.HasSuffix(base, "!") {
		important, base = true, base[:len(base)-1]
	}
	return mods, important, base
}

// orderSensitive variants change meaning when reordered, so they (and
// arbitrary variants) are kept in place; the variants between them are
// sorted so "hover:focus:" and "focus:hover:" compare equal.
var orderSensitive = map[string]bool{
	"before": true, "after": true, "placeholder": true, "file": true,
	"marker": true, "selection": true, "first-line": true, "first-letter": true,
	"backdrop": true, "*": true, "**": true,
}

func modifierKey(mods []string) string {
	if len(mods) == 0 {
		return ""
	}
	out := make([]string, 0, len(mods))
	var run []string
	flush := func() {
		slices.Sort(run)
		out = append(out, run...)
		run = run[:0]
	}
	for _, m := range mods {
		if strings.HasPrefix(m, "[") || orderSensitive[m] {
			flush()
			out = append(out, m)
			continue
		}
		run = append(run, m)
	}
	flush()
	return strings.Join(out, ":") + ":"
}

// --- cache -------------------------------------------------------------------

// Component class lists repeat constantly, so results are memoized. The cache
// is bounded: when full it is simply cleared.
type boundedCache struct {
	mu sync.RWMutex
	m  map[cacheKey]string
}

const cacheSize = 4096

var cache = &boundedCache{m: make(map[cacheKey]string, 256)}

func (c *boundedCache) get(k cacheKey) (string, bool) {
	c.mu.RLock()
	v, ok := c.m[k]
	c.mu.RUnlock()
	return v, ok
}

func (c *boundedCache) put(k cacheKey, v string) {
	c.mu.Lock()
	if len(c.m) >= cacheSize {
		clear(c.m)
	}
	c.m[k] = v
	c.mu.Unlock()
}
