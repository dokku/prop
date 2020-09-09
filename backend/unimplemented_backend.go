package backend

import (
	"fmt"
)

type UnimplementedBackend struct {
}

// NewUnimplementedBackend create new instance of UnimplementedBackend
func NewUnimplementedBackend() (UnimplementedBackend, error) {
	return UnimplementedBackend{}, nil
}

func (backend UnimplementedBackend) BackendExport() (PropertyCollection, bool, error) {
	return PropertyCollection{}, false, fmt.Errorf("Not implemented")
}

func (backend UnimplementedBackend) BackendImport(clear bool) (PropertyCollection, bool, error) {
	return PropertyCollection{}, false, fmt.Errorf("Not implemented")
}

func (backend UnimplementedBackend) BackendClear() (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (backend UnimplementedBackend) Del(key string) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (backend UnimplementedBackend) Exists(key string) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (backend UnimplementedBackend) NamespaceExists(namespace string) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (backend UnimplementedBackend) NamespaceClear(namespace string) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (backend UnimplementedBackend) Get(key string, defaultValue string) (string, error) {
	return "", fmt.Errorf("Not implemented")
}

func (backend UnimplementedBackend) GetAll() (map[string]string, error) {
	keyValuePairs := make(map[string]string)
	return keyValuePairs, fmt.Errorf("Not implemented")
}

func (backend UnimplementedBackend) GetAllByPrefix(prefix string) (map[string]string, error) {
	keyValuePairs := make(map[string]string)
	return keyValuePairs, fmt.Errorf("Not implemented")
}

func (backend UnimplementedBackend) Set(key string, value string) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (backend UnimplementedBackend) Lindex(key string, index int) (string, error) {
	return "", fmt.Errorf("Not implemented")
}

func (backend UnimplementedBackend) Lismember(key string, element string) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (backend UnimplementedBackend) Llen(key string) (int, error) {
	return 0, fmt.Errorf("Not implemented")
}

func (backend UnimplementedBackend) Lrange(key string) ([]string, error) {
	return []string{}, fmt.Errorf("Not implemented")
}

func (backend UnimplementedBackend) Lrangefrom(key string, start int) ([]string, error) {
	return []string{}, fmt.Errorf("Not implemented")
}

func (backend UnimplementedBackend) Lrangefromto(key string, start int, stop int) ([]string, error) {
	return []string{}, fmt.Errorf("Not implemented")
}

func (backend UnimplementedBackend) Lrem(key string, countToRemove int, element string) (int, error) {
	return 0, fmt.Errorf("Not implemented")
}

func (backend UnimplementedBackend) Lset(key string, index int, element string) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (backend UnimplementedBackend) Rpush(key string, newElements ...string) (int, error) {
	return 0, fmt.Errorf("Not implemented")
}

func (backend UnimplementedBackend) Sadd(key string, newMembers ...string) (int, error) {
	return 0, fmt.Errorf("Not implemented")
}

func (backend UnimplementedBackend) Sismember(key string, member string) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (backend UnimplementedBackend) Smembers(key string) (map[string]bool, error) {
	return map[string]bool{}, fmt.Errorf("Not implemented")
}

func (backend UnimplementedBackend) Srem(key string, membersToRemove ...string) (int, error) {
	return 0, fmt.Errorf("Not implemented")
}
