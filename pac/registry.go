package pac

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"strings"
)

func getRegistryValue(path string) (string, error) {
	var (
		key registry.Key
		err error
	)

	const HKCR = `HKEY_CLASSES_ROOT\`
	const HKCU = `HKEY_CURRENT_USER\`
	const HKLM = `HKEY_LOCAL_MACHINE\`
	const HKU = `HKEY_USERS\`
	const HKCC = `HKEY_CURRENT_CONFIG\`

	i := strings.LastIndex(path, `\`)
	if strings.HasPrefix(path, HKCR) && i > len(HKCR) {
		key, err = registry.OpenKey(registry.CLASSES_ROOT, path[len(HKCR):i], registry.QUERY_VALUE)
	} else if strings.HasPrefix(path, HKCU) && i > len(HKCU) {
		key, err = registry.OpenKey(registry.CURRENT_USER, path[len(HKCU):i], registry.QUERY_VALUE)
	} else if strings.HasPrefix(path, HKLM) && i > len(HKLM) {
		key, err = registry.OpenKey(registry.LOCAL_MACHINE, path[len(HKLM):i], registry.QUERY_VALUE)
	} else if strings.HasPrefix(path, HKU) && i > len(HKU) {
		key, err = registry.OpenKey(registry.USERS, path[len(HKU):i], registry.QUERY_VALUE)
	} else if strings.HasPrefix(path, HKCC) && i > len(HKCC) {
		key, err = registry.OpenKey(registry.CURRENT_CONFIG, path[len(HKCC):i], registry.QUERY_VALUE)
	} else {
		return "", fmt.Errorf("path '%s' is not a valid registry path", path)
	}

	if err != nil {
		return "", err
	}
	defer key.Close()

	value, z, err := key.GetStringValue(path[i+1:])
	if err != nil {
		return "", err
	}
	_ = z

	return value, nil
}
