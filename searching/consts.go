package searching

import "time"

const (
	interval = 3
	timeout  = 15 * time.Second
	mapping  = `{
		"%s": {
			"dynamic_templates": [{
				"long2float": {
					"match_mapping_type": "long",
					"mapping": {
						"type": "float"
					}
				}
			}]
		}
	}`
)
