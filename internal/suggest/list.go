package suggest

import "math/rand/v2"

// List is a Suggester backed by a fixed list of ingredients. It ignores the
// requested ingredient and returns one item picked at random from the list.
//
// This is the first-iteration stub: it proves the end-to-end path (HTTP →
// Suggester → JSON) returns a real suggestion without yet matching the input.
type List struct {
	items []Alternative
}

// ingredients is the fixed set of food & beverage items the API can suggest.
// Per AGENTS.md, every entry is a food or beverage ingredient.
var ingredients = []Alternative{
	{Name: "butter"},
	{Name: "olive oil"},
	{Name: "coconut oil"},
	{Name: "ghee"},
	{Name: "margarine"},
	{Name: "sunflower oil"},
	{Name: "honey"},
	{Name: "maple syrup"},
	{Name: "agave nectar"},
	{Name: "molasses"},
	{Name: "brown sugar"},
	{Name: "oat milk"},
	{Name: "almond milk"},
	{Name: "soy milk"},
	{Name: "coconut milk"},
	{Name: "buttermilk"},
	{Name: "greek yogurt"},
	{Name: "cashew cream"},
	{Name: "applesauce"},
	{Name: "flaxseed meal"},
}

// NewList returns a List seeded with the fixed ingredient set.
func NewList() *List {
	return &List{items: ingredients}
}

// Suggest returns a single ingredient chosen at random from the fixed list.
// The requested ingredient is ignored in this iteration. It returns
// ErrNotFound only if the list is empty.
func (l *List) Suggest(_ string) ([]Alternative, error) {
	if len(l.items) == 0 {
		return nil, ErrNotFound
	}
	pick := l.items[rand.IntN(len(l.items))]
	return []Alternative{pick}, nil
}
