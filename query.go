package elasticquery

import (
	"github.com/mtlkn/json"
)

type Query interface {
	JSON() *json.Object
}
