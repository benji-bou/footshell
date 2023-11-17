package loader

import (
	"fmt"
	"path/filepath"
	"plugin"
	"strings"
)

func ParseSymbolName(pluginPath string) string {

	pluginName := filepath.Base(pluginPath)
	extension := filepath.Ext(pluginName)
	pluginName = pluginName[:len(pluginName)-len(extension)]
	symbolName := strings.ToUpper(string(pluginName[0])) + pluginName[1:]
	return symbolName
}

func Symbol[T any](path string, prefix string) (T, string, error) {

	var zeroValue T
	plug, err := plugin.Open(path)

	if err != nil {
		return zeroValue, "", fmt.Errorf("loading plugin failed: %w", err)
	}
	symbolName := ParseSymbolName(path)
	sym, err := plug.Lookup(prefix + symbolName)
	if err != nil {
		return zeroValue, symbolName, fmt.Errorf("loading plugin symbol failed: %w", err)
	}
	loadedPlugin, isOk := sym.(T)
	if !isOk {
		return zeroValue, symbolName, fmt.Errorf("loaded plugin symbol does not respect %T interface", zeroValue)
	}
	return loadedPlugin, symbolName, nil
}
