package mv671

import (
	_ "embed"
	"github.com/oomph-ac/mv/multiversion/mappings"
)

var (
	//go:embed mappings/block_states.nbt
	blockStates []byte
	//go:embed mappings/item_runtime_ids.nbt
	itemRuntimeIDs []byte

	Mapping mappings.MVMapping
)

func init() {
	Mapping = mappings.Mapping(blockStates, itemRuntimeIDs)
}
