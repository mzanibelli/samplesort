package cache

import "path/filepath"

func Path(root, key, format string) (string, error) {
	file := key + format
	if filepath.IsAbs(file) {
		return file, nil
	}
	if filepath.Base(file) == file {
		return filepath.Join(root, file), nil
	}
	rel, err := filepath.Rel(root, file)
	if err != nil {
		return "", err
	}
	return filepath.Join(root, rel), nil
}
