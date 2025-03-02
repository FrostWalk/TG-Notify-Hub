package formatters

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"plugin"
	"strings"
)

// normal maps are safe for concurrent read access
var loadedPlugins = make(map[string]Formatter)

func loadPlugin(path string) (Formatter, error) {
	p, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}

	// Look up the exported symbol "Formatter".
	symPlugin, err := p.Lookup("Formatter")
	if err != nil {
		return nil, err
	}

	// Assert that the symbol implements the plugin.Formatter interface.
	plug, ok := symPlugin.(Formatter)
	if !ok {
		return nil, fmt.Errorf("unexpected type for Formatter in %s", path)
	}

	return plug, nil
}

// LoadPluginsFromFolder load all the .so file present in the folder and map each of them to his slug
func LoadPluginsFromFolder(folder string) error {
	err := ensureFolderExists(folder)
	if err != nil {
		return err
	}

	entries, err := os.ReadDir(folder)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		// Only process files (skip directories) that have the .so extension.
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".so") {
			continue
		}

		fullPath := filepath.Join(folder, entry.Name())
		plug, err := loadPlugin(fullPath)
		if err != nil {
			log.Printf("Error loading plugin %s: %v", fullPath, err)
			continue
		}

		loadedPlugins[plug.Slug()] = plug
		log.Printf("Loaded plugin: %s from %s", plug.Slug(), fullPath)
	}

	return nil
}

// GetPluginFromSlug retrieves a loaded plugin by its slug.
func GetPluginFromSlug(slug string) (bool, Formatter) {
	plug, ok := loadedPlugins[slug]
	return ok, plug
}

func ensureFolderExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}
