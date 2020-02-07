package gapt

var assets = &AssetManager{}

type AssetManager struct {
}

// Assets returns the singleton application wide asset manager instance.
// All modules are required to register t??
func Assets() *AssetManager {
	return assets
}
