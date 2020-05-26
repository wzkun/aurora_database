package searching

import "encoding/json"

// ElasticItem A generic value can store in the cache.
type ElasticItem interface {
	Idx() string
	Parrent() string
	JSONData() json.RawMessage
}
