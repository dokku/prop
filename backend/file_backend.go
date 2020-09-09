package backend

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strconv"
	"strings"

	"github.com/xo/dburl"
)

type UnstructuredFileBackend struct {
	NamespaceRoot string
	Namespace     string
	SystemUser    string
	SystemGroup   string
}

// NewUnstructuredFileBackend create new instance of UnstructuredFileBackend
func NewUnstructuredFileBackend(namespace string, url *dburl.URL) (UnstructuredFileBackend, error) {
	systemUser := url.Query().Get("system-user")
	systemGroup := url.Query().Get("system-group")
	backend := UnstructuredFileBackend{}
	backend.NamespaceRoot = path.Join(url.Opaque, namespace)
	backend.Namespace = namespace
	backend.SystemUser = systemUser
	backend.SystemGroup = systemGroup
	return backend, nil
}

func (backend UnstructuredFileBackend) BackendExport() (PropertyCollection, error) {
	return PropertyCollection{}, fmt.Errorf("Not implemented")
}

func (backend UnstructuredFileBackend) BackendImport(p PropertyCollection, clear bool) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (backend UnstructuredFileBackend) BackendReset() (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (backend UnstructuredFileBackend) Del(key string) (bool, error) {
	keyPath := backend.getKeyPath(key)
	if err := os.Remove(keyPath); err != nil {
		if exists, _ := backend.Exists(key); !exists {
			return true, nil
		}

		return false, fmt.Errorf("Unable to remove key %s.%s", backend.Namespace, key)
	}

	return true, nil
}

func (backend UnstructuredFileBackend) Exists(key string) (bool, error) {
	keyPath := backend.getKeyPath(key)
	_, err := os.Stat(keyPath)
	if err != nil {
		return false, err
	}

	return !os.IsNotExist(err), nil
}

func (backend UnstructuredFileBackend) NamespaceExists(namespace string) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (backend UnstructuredFileBackend) NamespaceClear(namespace string) (bool, error) {
	return false, fmt.Errorf("Not implemented")
}

func (backend UnstructuredFileBackend) Get(key string, defaultValue string) (string, error) {
	if exists, _ := backend.Exists(key); !exists {
		if defaultValue != "" {
			return defaultValue, nil
		}

		return "", fmt.Errorf("Key does not exist in namespace")
	}

	keyPath := backend.getKeyPath(key)
	b, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return "", fmt.Errorf("Unable to read key %s.%s", backend.Namespace, key)
	}

	return string(b), nil
}

func (backend UnstructuredFileBackend) GetAll() (map[string]string, error) {
	keyValuePairs := make(map[string]string)
	files, err := ioutil.ReadDir(backend.NamespaceRoot)
	if err != nil {
		return keyValuePairs, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		key := file.Name()
		keyValuePairs[key], _ = backend.Get(key, "")
	}

	return keyValuePairs, nil
}

func (backend UnstructuredFileBackend) GetAllByPrefix(prefix string) (map[string]string, error) {
	keyValuePairs, err := backend.GetAll()
	if err != nil {
		return map[string]string{}, err
	}

	response := make(map[string]string)
	for key, value := range keyValuePairs {
		if strings.HasPrefix(key, prefix) {
			response[key] = value
		}
	}

	return response, nil
}

func (backend UnstructuredFileBackend) Set(key string, value string) (bool, error) {
	if err := backend.touchKey(key); err != nil {
		return false, err
	}

	keyPath := backend.getKeyPath(key)
	file, err := os.Create(keyPath)
	if err != nil {
		return false, fmt.Errorf("Unable to write config value %s.%s: %s", backend.Namespace, key, err.Error())
	}
	defer file.Close()

	fmt.Fprintf(file, value)
	file.Chmod(0600)
	backend.setPermissions(keyPath, 0600)

	return true, nil
}

func (backend UnstructuredFileBackend) Lindex(key string, index int) (string, error) {
	lines, err := backend.Lrange(key)
	if err != nil {
		return "", err
	}

	for i, line := range lines {
		if i == index {
			return line, nil
		}
	}

	return "", fmt.Errorf("Index not found in key: %s.%s", backend.Namespace, key)
}

func (backend UnstructuredFileBackend) Lismember(key string, element string) (bool, error) {
	lines, err := backend.Lrange(key)
	if err != nil {
		return false, err
	}

	for _, line := range lines {
		if line == element {
			return true, nil
		}
	}

	return false, fmt.Errorf("Value not found in list: %s.%s", backend.Namespace, key)
}

func (backend UnstructuredFileBackend) Llen(key string) (int, error) {
	elements, err := backend.Lrange(key)
	if err != nil {
		return 0, err
	}

	return len(elements), nil
}

func (backend UnstructuredFileBackend) Lrange(key string) ([]string, error) {
	if exists, _ := backend.Exists(key); !exists {
		return []string{}, fmt.Errorf("Key does not exist in namespace")
	}

	keyPath := backend.getKeyPath(key)
	file, err := os.Open(keyPath)
	if err != nil {
		return []string{}, err
	}
	defer file.Close()

	var values []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		values = append(values, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		return values, fmt.Errorf("Unable to read config value for %s.%s: %s", backend.Namespace, key, err.Error())
	}

	return values, nil
}

func (backend UnstructuredFileBackend) Lrangefrom(key string, start int) ([]string, error) {
	elements, err := backend.Lrange(key)
	if err != nil {
		return []string{}, err
	}

	var values []string
	if start > len(elements) {
		return values, nil
	}

	for i, element := range elements {
		if i >= start {
			values = append(values, element)
		}
	}

	return values, nil
}

func (backend UnstructuredFileBackend) Lrangefromto(key string, start int, stop int) ([]string, error) {
	elements, err := backend.Lrange(key)
	if err != nil {
		return []string{}, err
	}

	var values []string
	if start > len(elements) {
		return values, nil
	}

	for i, element := range elements {
		if i >= start && i <= stop {
			values = append(values, element)
		}
	}

	return values, nil
}

func (backend UnstructuredFileBackend) Lrem(key string, countToRemove int, element string) (int, error) {
	elements, err := backend.Lrange(key)
	if err != nil {
		return 0, err
	}

	var newElements []string
	removed := 0
	if countToRemove == 0 {
		for _, e := range elements {
			if e == element {
				removed++
				continue
			}
			newElements = append(newElements, e)
		}
	} else {
		if countToRemove < 0 {
			reverse(elements)
		}
		for _, e := range elements {
			if e == element {
				if removed != countToRemove {
					removed++
					continue
				}
			}
			newElements = append(newElements, e)
		}

		if countToRemove < 0 {
			reverse(elements)
		}
	}

	if err = backend.writeList(key, newElements); err != nil {
		return 0, err
	}

	return removed, nil
}

func (backend UnstructuredFileBackend) Lset(key string, index int, element string) (bool, error) {
	if err := backend.touchKey(key); err != nil {
		return false, err
	}

	elements, err := backend.Lrange(key)
	if err != nil {
		return false, err
	}

	absIndex := index
	if index < 0 {
		absIndex = -index - 1
	}

	if absIndex >= len(elements) {
		return false, fmt.Errorf("Index out of range")
	}

	var newElements []string
	element = strings.TrimSpace(element)
	if index < 0 {
		reverse(elements)
	}

	for i, line := range elements {
		if i == absIndex {
			newElements = append(newElements, element)
		} else {
			newElements = append(newElements, line)
		}
	}

	if index < 0 {
		reverse(newElements)
	}

	if err = backend.writeList(key, newElements); err != nil {
		return false, err
	}

	return true, nil
}

func (backend UnstructuredFileBackend) Rpush(key string, newElements ...string) (int, error) {
	if err := backend.touchKey(key); err != nil {
		return 0, err
	}

	elements, err := backend.Lrange(key)
	if err != nil {
		return 0, err
	}

	elements = append(elements, newElements...)

	if err = backend.writeList(key, elements); err != nil {
		return 0, err
	}

	return len(elements), nil
}

func (backend UnstructuredFileBackend) Sadd(key string, newMembers ...string) (int, error) {
	if err := backend.touchKey(key); err != nil {
		return 0, err
	}

	members, err := backend.Smembers(key)
	if err != nil {
		return 0, err
	}

	addedCount := 0
	for _, member := range newMembers {
		if _, ok := members[member]; !ok {
			members[member] = true
			addedCount++
		}
	}

	if err = backend.writeSet(key, members); err != nil {
		return 0, err
	}

	return addedCount, nil
}

func (backend UnstructuredFileBackend) Sismember(key string, member string) (bool, error) {
	if exists, _ := backend.Exists(key); !exists {
		return false, fmt.Errorf("Set does not exist: %s.%s", backend.Namespace, key)
	}

	members, err := backend.Smembers(key)
	if err != nil {
		return false, err
	}

	_, ok := members[member]
	return ok, nil
}

func (backend UnstructuredFileBackend) Smembers(key string) (map[string]bool, error) {
	members := make(map[string]bool)
	if exists, _ := backend.Exists(key); !exists {
		return members, fmt.Errorf("Key does not exist in namespace")
	}

	keyPath := backend.getKeyPath(key)
	file, err := os.Open(keyPath)
	if err != nil {
		return members, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		members[scanner.Text()] = true
	}

	if err = scanner.Err(); err != nil {
		return members, fmt.Errorf("Unable to read config value for %s.%s: %s", backend.Namespace, key, err.Error())
	}

	return members, nil
}

func (backend UnstructuredFileBackend) Srem(key string, membersToRemove ...string) (int, error) {
	if exists, _ := backend.Exists(key); !exists {
		return 0, fmt.Errorf("Key does not exist in namespace")
	}

	members, err := backend.Smembers(key)
	if err != nil {
		return 0, err
	}

	removedCount := 0
	for _, member := range membersToRemove {
		if _, ok := members[member]; ok {
			delete(members, member)
			removedCount++
		}
	}
	if err = backend.writeSet(key, members); err != nil {
		return 0, err
	}

	return removedCount, nil
}

func (backend UnstructuredFileBackend) getKeyPath(key string) string {
	return path.Join(backend.NamespaceRoot, key)
}

// propertyTouch ensures a given application property file exists
func (backend UnstructuredFileBackend) touchKey(key string) error {
	if exists, _ := backend.Exists(key); exists {
		return nil
	}

	if err := backend.makeNamespaceDirectory(); err != nil {
		return fmt.Errorf("Unable to create config directory for %s: %s", backend.Namespace, err.Error())
	}

	keyPath := backend.getKeyPath(key)
	file, err := os.Create(keyPath)
	if err != nil {
		return fmt.Errorf("Unable to writeconfig value %s.%s: %s", backend.Namespace, key, err.Error())
	}
	defer file.Close()

	return nil
}

func (backend UnstructuredFileBackend) writeList(key string, elements []string) error {
	keyPath := backend.getKeyPath(key)
	file, err := os.OpenFile(keyPath, os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(file)
	for _, element := range elements {
		fmt.Fprintln(w, element)
	}
	if err = w.Flush(); err != nil {
		return fmt.Errorf("Unable to write config value %s.%s: %s", backend.Namespace, key, err.Error())
	}

	file.Chmod(0600)
	backend.setPermissions(keyPath, 0600)
	return nil
}

func (backend UnstructuredFileBackend) writeSet(key string, members map[string]bool) error {
	keyPath := backend.getKeyPath(key)
	file, err := os.OpenFile(keyPath, os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(file)
	for member := range members {
		fmt.Fprintln(w, member)
	}
	if err = w.Flush(); err != nil {
		return fmt.Errorf("Unable to write config value %s.%s: %s", backend.Namespace, key, err.Error())
	}

	file.Chmod(0600)
	backend.setPermissions(keyPath, 0600)
	return nil
}

// makeNamespaceDirectory ensures that a property path exists
func (backend UnstructuredFileBackend) makeNamespaceDirectory() error {
	if err := os.MkdirAll(backend.NamespaceRoot, 0755); err != nil {
		return err
	}
	return backend.setPermissions(backend.NamespaceRoot, 0755)
}

// setPermissions sets the proper owner and filemode for a given file
func (backend UnstructuredFileBackend) setPermissions(path string, fileMode os.FileMode) error {
	if err := os.Chmod(path, fileMode); err != nil {
		return err
	}

	group, err := user.LookupGroup(backend.SystemGroup)
	if err != nil {
		return err
	}
	user, err := user.Lookup(backend.SystemUser)
	if err != nil {
		return err
	}

	uid, err := strconv.Atoi(user.Uid)
	if err != nil {
		return err
	}

	gid, err := strconv.Atoi(group.Gid)
	if err != nil {
		return err
	}
	return os.Chown(path, uid, gid)
}

func getenvWithDefault(key string, defaultValue string) (val string) {
	val = os.Getenv(key)
	if val == "" {
		val = defaultValue
	}
	return
}

func reverse(ss []string) {
	last := len(ss) - 1
	for i := 0; i < len(ss)/2; i++ {
		ss[i], ss[last-i] = ss[last-i], ss[i]
	}
}
