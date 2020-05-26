package searching

import (
	sl "code.aliyun.com/bim_backend/zoogoer/gun/singleflight"
)

// StorageWithElastic struct
type StorageWithElastic struct {
	lg sl.Group
	es *ElasticClient
}

// NewStorageWithElastic create StorageWithElastic
func NewStorageWithElastic(es *ElasticClient) *StorageWithElastic {
	o := new(StorageWithElastic)
	o.es = es
	return o
}

// Set function
func (m *StorageWithElastic) Set(item ElasticItem) error {
	return m.es.Update(item)
}

// SetMulti function
func (m *StorageWithElastic) SetMulti(items ...ElasticItem) ([]ElasticItem, error) {
	serr := m.es.UpdateMulti(items...)
	if serr != nil {
		return nil, serr
	}
	return items, nil
}

// Update function
func (m *StorageWithElastic) Update(item ElasticItem) error {
	return m.es.Update(item)
}

// UpdateMulti function
func (m *StorageWithElastic) UpdateMulti(items ...ElasticItem) ([]ElasticItem, error) {
	serr := m.es.UpdateMulti(items...)
	if serr != nil {
		return nil, serr
	}
	return items, nil
}

// SetIfAbsent function
func (m *StorageWithElastic) SetIfAbsent(item ElasticItem) error {
	return m.es.Create(item)
}

// SetMultiIfAbsent function
func (m *StorageWithElastic) SetMultiIfAbsent(items ...ElasticItem) ([]ElasticItem, error) {
	serr := m.es.CreateMulti(items...)
	if serr != nil {
		return nil, serr
	}

	return items, nil
}

// Delete function
func (m *StorageWithElastic) Delete(parrent, id string) error {
	return m.es.Delete(id, parrent)
}

// DeleteMulti function
func (m *StorageWithElastic) DeleteMulti(parrent string, ids ...string) error {
	serr := m.es.DeleteMulti(parrent, ids...)
	if serr != nil {
		return serr
	}

	return nil
}

// Flush function
func (m *StorageWithElastic) Flush(parrent string) error {
	return m.es.DeleteIndex(parrent)
}
