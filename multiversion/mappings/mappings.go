package mappings

// MVMapping holds all data blocks, items related.
type MVMapping struct {
	MVBlockMapping
	MVItemMapping
}

// Mapping returns MVMapping instance of all block and item entries and values in the maps from the resource JSON.
func Mapping(blockStateData, itemRuntimeIDData []byte) MVMapping {
	return MVMapping{
		MVBlockMapping: blockMapping(blockStateData),
		MVItemMapping:  itemMapping(itemRuntimeIDData),
	}
}
