package vault

import "sort"

// CompareOptions controls how secret comparisons are performed.
type CompareOptions struct {
	IgnoreKeys    []string
	OnlyKeys      []string
	CaseSensitive bool
}

// CompareResult holds the result of comparing two sets of secrets.
type CompareResult struct {
	Matching   []string
	Differing  []string
	OnlyInA    []string
	OnlyInB    []string
}

// Summary returns a human-readable summary of the comparison.
func (r *CompareResult) Summary() string {
	return formatCompareSummary(r)
}

func formatCompareSummary(r *CompareResult) string {
	s := ""
	s += formatLine("Matching", len(r.Matching))
	s += formatLine("Differing", len(r.Differing))
	s += formatLine("Only in A", len(r.OnlyInA))
	s += formatLine("Only in B", len(r.OnlyInB))
	return s
}

func formatLine(label string, count int) string {
	if count == 0 {
		return ""
	}
	return label + ": " + itoa(count) + "\n"
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	result := ""
	for n > 0 {
		result = string(rune('0'+n%10)) + result
		n /= 10
	}
	return result
}

// CompareSecrets compares two secret maps and returns a CompareResult.
func CompareSecrets(a, b map[string]string, opts CompareOptions) CompareResult {
	result := CompareResult{}

	skip := func(key string) bool {
		for _, ig := range opts.IgnoreKeys {
			if ig == key {
				return true
			}
		}
		if len(opts.OnlyKeys) > 0 {
			for _, ok := range opts.OnlyKeys {
				if ok == key {
					return false
				}
			}
			return true
		}
		return false
	}

	allKeys := map[string]struct{}{}
	for k := range a {
		allKeys[k] = struct{}{}
	}
	for k := range b {
		allKeys[k] = struct{}{}
	}

	keys := make([]string, 0, len(allKeys))
	for k := range allKeys {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		if skip(k) {
			continue
		}
		va, inA := a[k]
		vb, inB := b[k]
		switch {
		case inA && inB:
			if va == vb {
				result.Matching = append(result.Matching, k)
			} else {
				result.Differing = append(result.Differing, k)
			}
		case inA:
			result.OnlyInA = append(result.OnlyInA, k)
		case inB:
			result.OnlyInB = append(result.OnlyInB, k)
		}
	}
	return result
}
