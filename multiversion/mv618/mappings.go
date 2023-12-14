package mv618

import (
	_ "embed"

	"github.com/oomph-ac/mv/multiversion/latest"
	"github.com/oomph-ac/mv/multiversion/mappings"
)

var (
	//go:embed mappings/block_states.nbt
	blockStates []byte
	//go:embed mappings/item_runtime_ids.nbt
	itemRuntimeIDData []byte

	Mapping mappings.MVMapping
)

func init() {
	Mapping = mappings.Mapping(blockStates, latest.ItemRuntimeIDData, false)
}
