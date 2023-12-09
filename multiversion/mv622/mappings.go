package mv622

import (
	_ "embed"

	"github.com/oomph-ac/mv/multiversion/latest"
	"github.com/oomph-ac/mv/multiversion/mappings"
)

var (
	//go:embed mappings/block_states.nbt
	blockStates []byte

	Mapping mappings.MVMapping
)

func init() {
	Mapping = mappings.Mapping(blockStates, latest.ItemRuntimeIDData, false)
}
