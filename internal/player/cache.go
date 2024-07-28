package player

func NewCache(update func() string) *ViewCache {
	return &ViewCache{
		update: update,
	}
}

type ViewCache struct {
	update func() string
	cache  string
	valid  bool
}

func (v *ViewCache) String() string {
	if !v.valid {
		v.valid = true
		v.cache = v.update()
	}
	return v.cache
}

func (v *ViewCache) Invalidate() {
	v.valid = false
}
