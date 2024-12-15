package arraystore

import "github.com/vadim-ivlev/url-shortener/internal/apptypes"

type Index struct {
	data    map[string]int64
	keyFunc func(apptypes.URLShortener) string
}

// NewIndex initializes a new index with a key function
func NewIndex(keyFunc func(apptypes.URLShortener) string) *Index {
	return &Index{
		data:    make(map[string]int64),
		keyFunc: keyFunc,
	}
}

// Add adds a record to the index
func (i *Index) Add(record apptypes.URLShortener, idx int64) {
	i.data[i.keyFunc(record)] = idx
}

// Get returns the index of the record
func (i *Index) Get(record apptypes.URLShortener) (int64, bool) {
	idx, ok := i.data[i.keyFunc(record)]
	return idx, ok
}

// Delete removes a record from the index
func (i *Index) Delete(record apptypes.URLShortener) {
	delete(i.data, i.keyFunc(record))
}
