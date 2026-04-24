package vault

import (
	"encoding/json"
	"fmt"
	"os"
)

// policyFile is the on-disk representation of a Policy.
type policyFile struct {
	Name  string       `json:"name"`
	Rules []policyRule `json:"rules"`
}

type policyRule struct {
	KeyPattern  string `json:"key_pattern"`
	AllowRead   bool   `json:"allow_read"`
	AllowWrite  bool   `json:"allow_write"`
	AllowDelete bool   `json:"allow_delete"`
}

// LoadPolicy reads a JSON policy file from path and returns a Policy.
func LoadPolicy(path string) (Policy, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Policy{}, fmt.Errorf("loading policy file %q: %w", path, err)
	}

	var pf policyFile
	if err := json.Unmarshal(data, &pf); err != nil {
		return Policy{}, fmt.Errorf("parsing policy file %q: %w", path, err)
	}

	if pf.Name == "" {
		return Policy{}, fmt.Errorf("policy file %q missing required field \"name\"", path)
	}

	p := Policy{Name: pf.Name}
	for _, r := range pf.Rules {
		if r.KeyPattern == "" {
			return Policy{}, fmt.Errorf("policy %q has a rule with empty key_pattern", pf.Name)
		}
		p.Rules = append(p.Rules, PolicyRule{
			KeyPattern:  r.KeyPattern,
			AllowRead:   r.AllowRead,
			AllowWrite:  r.AllowWrite,
			AllowDelete: r.AllowDelete,
		})
	}
	return p, nil
}
