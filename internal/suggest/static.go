package suggest

import "strings"

// Static is an in-memory Suggester backed by a hardcoded map.
//
// It is the stub implementation used to bootstrap the project. Replace it
// with a database- or LLM-backed Suggester later without changing callers.
type Static struct {
	data map[string][]Alternative
}

// NewStatic returns a Static seeded with a small set of common substitutions.
func NewStatic() *Static {
	return &Static{
		data: map[string][]Alternative{
			"butter": {
				{Name: "olive oil", Notes: "use ~3/4 the amount for baking"},
				{Name: "coconut oil", Notes: "adds a mild coconut flavour"},
			},
			"milk": {
				{Name: "oat milk", Notes: "neutral, good for baking"},
				{Name: "almond milk", Notes: "lower calorie, slightly nutty"},
			},
			"egg": {
				{Name: "flaxseed meal", Notes: "1 tbsp + 3 tbsp water per egg"},
				{Name: "applesauce", Notes: "1/4 cup per egg, adds moisture"},
			},
			"sugar": {
				{Name: "honey", Notes: "use ~3/4 the amount, reduce liquids"},
				{Name: "maple syrup", Notes: "use ~3/4 the amount"},
			},
		},
	}
}

// Suggest returns alternatives for the given ingredient, matching
// case-insensitively. It returns ErrNotFound when the ingredient is unknown.
func (s *Static) Suggest(ingredient string) ([]Alternative, error) {
	key := strings.ToLower(strings.TrimSpace(ingredient))
	alts, ok := s.data[key]
	if !ok {
		return nil, ErrNotFound
	}
	return alts, nil
}
