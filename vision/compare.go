package vision

import (
	"image"
	"image/png"
	"os"
	"sync"
	"time"
)

type Comparer struct{ cache *ReferenceCache }

func NewComparer() *Comparer { return &Comparer{cache: NewReferenceCache()} }

type cachedHash struct {
	hash  [64]uint8
	mtime time.Time
}

type ReferenceCache struct {
	mu    sync.RWMutex
	cache map[string]cachedHash
}

func NewReferenceCache() *ReferenceCache { return &ReferenceCache{cache: map[string]cachedHash{}} }

func hashImg(img image.Image) [64]uint8 {
	n := Normalize(img)
	var out [64]uint8
	k := 0
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			r, _, _, _ := n.At(x, y).RGBA()
			if r>>8 > 127 {
				out[k] = 1
			}
			k++
		}
	}
	return out
}

func dist(a, b [64]uint8) int {
	d := 0
	for i := 0; i < 64; i++ {
		if a[i] != b[i] {
			d++
		}
	}
	return d
}

func (r *ReferenceCache) get(path string) ([64]uint8, error) {
	st, err := os.Stat(path)
	if err != nil {
		return [64]uint8{}, err
	}
	r.mu.RLock()
	v, ok := r.cache[path]
	r.mu.RUnlock()
	if ok && v.mtime.Equal(st.ModTime()) {
		return v.hash, nil
	}
	f, err := os.Open(path)
	if err != nil {
		return [64]uint8{}, err
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		return [64]uint8{}, err
	}
	h := hashImg(img)
	r.mu.Lock()
	r.cache[path] = cachedHash{hash: h, mtime: st.ModTime()}
	r.mu.Unlock()
	return h, nil
}

func (c *Comparer) MatchesAny(live image.Image, refs []string, threshold int) (bool, error) {
	lh := hashImg(live)
	for _, ref := range refs {
		rh, err := c.cache.get(ref)
		if err != nil {
			return false, err
		}
		if dist(lh, rh) <= threshold {
			return true, nil
		}
	}
	return false, nil
}
