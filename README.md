# prop

A golang-based cli interface for manipulating config backed by various datastores.

## Summary

This describes a utility that can be used to store, retrieve, and manipulate various data types against a configured backing store. It can be used as both a library and wrapped via command-line calls.

## Motivation

When writing a command-line utility, it may be necessary to store configuration values on disk. This is typically implemented in a one-off manner, and state is both not always persisted in the same locations nor is it easily distributed for high-availability.

This project seeks to make it easy to manipulate various basic data types across a wide variety of backends via a single entrypoint, making it easier for applications to consider configuration state without needing to reimplement the wheel.

## Terminology

- Set: A mathematical, well-defined collection of distinct objects. See [wikipedia](https://en.wikipedia.org/wiki/Set_(mathematics)) for more details.
- List: An enumerated collection of objects in which repititions are allowed. See [wikipedia](https://en.wikipedia.org/wiki/Sequence) for more details. Lists are zero-indexed.
- Key-Value: A 2-tuple collection of data, addressable by the key name. See [wikipedia](https://en.wikipedia.org/wiki/Attribute–value_pair) for more details.
- Namespace: A distinct collection of symbols that are related to each other. All identifiers within a namespace are unique. See [wikipedia](https://en.wikipedia.org/wiki/Namespace) for more details.
- Backend: A data-access layer that can contain configuration state.

## `prop` configuration state

The `prop` tool stores it's own state locally in the configuration directory for a given OS, as per the [configdir](https://github.com/shibukawa/configdir) project. This directory will contain a `config.json` file that maintains prop's configuration as key-value pairs.

## Data Types

The following data types are implemented within the `prop` tool:

- sets
- lists
- key-value

Inside of `prop`, a bit of data is called a `Property` consists of the following interface:

```go
type Property struct {
  DataType  string
  Namespace string
  Key       string
  Value     string
}
```

A set of data is called a `PropertyCollection`. A `PropertyCollection` can be imported/exported from one backend to another.

```go
type PropertyCollection struct {
   Properties []Property
}
```

## Key and Value specification

Keys may follow the following regex:

```
[\w-/]{1,200}[^/]
```

Values may contain 0 or more utf8 characters and may be a maximum of 65535 characters in length.

## Commands

The following commands are supported.

### `backend` commands

#### `backend export path/to/file`

- Description: Exports a backend to a json file
- Method Signature: `func (b Backend) BackendExport() (p PropertyCollection, exported bool, err error)`

When export a backend, it is assumed that there are is no concurrent access to the backend. In other words, if another process is changing values of the backend, then the export may result in an invalid state.

#### `backend import path/to/file`

- Description: Import a backend to a json file
- Method Signature: `func (b Backend) BackendImport(p PropertyCollection, clear bool) (imported bool, err error)`
- Flags: `--clear`

When importing a backend, properties are merged into the existing backend unless the `--clear` flag is specified.

When migrating a backend, it is assumed that there are is no concurrent access to the backend. In other words, if another process is changing values of the backend, then the import may result in an invalid state.

#### `backend reset`

- Description: Clear all values in a backend
- Method Signature: `func (b Backend) BackendReset() (success bool, err error)`

### `config` commands

Used for configuring `prop`.

Current configuration values that may be manipulated:

- `url`:
  - Type: string
  - Default: `file:/etc/prop.d`
  - environment Variable: `PROP_BACKEND_URL`
  - Description: A configured backend for prop, specified in [DSN](https://en.wikipedia.org/wiki/Data_source_name) form. Backends are built into the prop project. Currently supported backends are `file` and `postgres`
- `namespace`:
  - Type: string
  - Default: `default`
  - Environment Variable: `PROP_NAMESPACE`
  - Description: The default namespace. Commands that allow namespace usage will note as such.

All properties may be specified as environment variables. The `config.json` holding the config will only be read if any of the environment variables are missing.

#### `config get key`

- Description: Get a configuration value

```shell
prop config get url
```

#### `config set key value`

- Description: Set a configuration value

```shell
prop config set url postgres://user:password@host:port/database
```

#### `config del key`

- Description: Delete a configuration value

```
prop config del url
```

### `namespace` commands

#### `namespace exists namespace`

- Description: Checks if there are any keys in a given namespace
- Method Signature: `func (b Backend) NamespaceExists(namespace string) (exists bool, err error)`

#### `namespace clear namespace`

- Description: Delete all keys from a given namespace
- Method Signature: `func (b Backend) NamespaceClear(namespace string) (success bool, err error)`

### global commands

#### `del key`

- Description: Delete a key
- Data Type: `key-value`, `list`, `set`
- Supported Flags: `--namespace`
- Method Signature: `func (b Backend) Del(key string) (success bool, err error)`

### `key-value` commands

### `exists key`

- Description: Check if a exists
- Data Type: `key-value`, `list`, `set`
- Supported Flags: `--namespace`
- Method Signature: `func (b Backend) Exists(key string) (exists bool, err error)`

#### `get key [default]`

- Description: Get the value of a key
- Data Type: `key-value`
- Supported Flags: `--namespace`
- Method Signature: `func (b Backend) Get(key string, defaultValue string) (value string, err error)`

#### `get-all [prefix]`

- Description: Get all key-value tuples
- Data Type: `[(key-value tuple)]`
- Supported Flags: `--namespace`
- Method Signature: `func (b Backend) GetAll() (keyValuePairs map[string]string, err error)`
- Method Signature: `func (b Backend) GetAllByPrefix(prefix string) (keyValuePairs map[string]string, err error)`

#### `set key value`

- Description: Set the string value of a key
- Data Type: `key-value`
- Supported Flags: `--namespace`
- Method Signature: `func (b Backend) Set(key string, value string) (success bool, err error)`

### `list` commands

#### `lindex key index`

- Description: Get an element from a list by its index
- Data Type: `list`
- Supported Flags: `--namespace`
- Method Signature: `func (b Backend) Lindex(key string, index int) (element string, err error)`

#### `lismember key element`

- Description: Determine if a given value is an element in the list
- Data Type: `list`
- Supported Flags: `--namespace`
- Method Signature: `func (b Backend) Lismember(key string, element string) (isMember bool, err error)`

#### `llen key`

- Description: Get the length of a list
- Data Type: `list`
- Supported Flags: `--namespace`
- Method Signature: `func (b Backend) Llen(key string) (length int, err error)`

#### `lrange key [start [stop]]`

- Description: Get a range of elements from a list
- Data Type: `list`
- Supported Flags: `--namespace`
- Method Signature: `func (b Backend) Lrange(key string) ([]string, err error)`
- Method Signature: `func (b Backend) Lrangefrom(key string, start int) ([]string, err error)`
- Method Signature: `func (b Backend) Lrangefromto(key string, start int, stop int) ([]string, err error)`

#### `lrem key count element`

- Description: Remove elements from a list
- Data Type: `list`
- Supported Flags: `--namespace`
- Method Signature: `func (b Backend) Lrem(key string, countToRemove int, element string) (removedCount int, err error)`

#### `lset key index element`

- Description: Set the value of an element in a list by its index
- Data Type: `list`
- Supported Flags: `--namespace`
- IntMethod Signatureerface: `func (b Backend) Lset(key string, index int, element string) (success bool, err error)`

#### `rpush key element [element...]`

- Description: Append one or more elements to a list
- Data Type: `list`
- Supported Flags: `--namespace`
- Method Signature: `func (b Backend) Rpush(key string, newElements ...string) (listLength int, err error)`

### `set` commands

#### `sadd key member [member ..]`

- Description: Add one or more members to a set
- Data Type: `set`
- Supported Flags: `--namespace`
- Method Signature: `func (b Backend) Sadd(key string, newMembers ...string) (addedCount int, err error)`

#### `sismember key member`

- Description: Determine if a given value is a member of a set
- Data Type: `set`
- Supported Flags: `--namespace`
- Method Signature: `func (b Backend) Sismember(key string, member string) (isMember bool, err error)`

#### `smembers key`

- Description: Get all the members in a set
- Data Type: `set`
- Supported Flags: `--namespace`
- Method Signature: `func (b Backend) Smembers(key string) (member map[string]bool, err error)`

#### `srem key member [member ...]`

- Description: Remove one or more members from a set
- Data Type: `set`
- Supported Flags: `--namespace`
- Method Signature: `func (b Backend) Srem(key string, membersToRemove...string) (removedCount int, err error)`

## Backends

Backends should implement the method signatures specified for each command. The following is the base interface:

```go
type Backend interface {
  BackendExport() (PropertyCollection, error)
  BackendImport(p PropertyCollection, clear bool) (bool, error)
  BackendReset() (bool, error)
  Del(key string) (bool, error)
  Exists(key string) (bool, error)
  NamespaceExists(namespace string) (bool, error)
  NamespaceClear(namespace string) (bool, error)
  Get(key string, defaultValue string) (string, error)
  GetAll() (map[string]string, error)
  GetAllByPrefix(prefix string) (map[string]string, error)
  Set(key string, value string) (bool, error)
  Lindex(key string, index int) (string, error)
  Lismember(key string, element string) (bool, error)
  Llen(key string) (int, error)
  Lrange(key string) ([]string, error)
  Lrangefrom(key string, start int) ([]string, error)
  Lrangefromto(key string, start int, stop int) ([]string, error)
  Lrem(key string, countToRemove int, element string) (int, error)
  Lset(key string, index int, element string) (bool, error)
  Rpush(key string, newElements ...string) (int, error)
  Sadd(key string, newMembers ...string) (int, error)
  Sismember(key string, member string) (bool, error)
  Smembers(key string) (map[string]bool, error)
  Srem(key string, membersToRemove ...string) (int, error)
}
```

The following backends are supported.

### File

To configure, run:

```shell
prop config set url file:/etc/prop.d
```

The directory structure is as follows:

```shell
# returns the contents of the file
cat $NAMESPACE/$KEY
```

Key names can include forward slashes, which will be interpreted as a directory structure.

Values are stored in the following json format:

```json
{
   "type": "$data_type",
   "value": "$value"
}
```

When querying for a property, if the type of the value does not match the type specified by the executed command, an error should be raised where possible.

### Redis

To configure, run:

```shell
prop config set url redis://user:password@host:port/database
```

With the redis backend, commands map to their redis equivalents where appropriate. If there is no equivalent redis command, a redis script may be used instead to implement the functionality.

When querying for a property, if the type of the value does not match the type specified by the executed command, an error should be raised where possible.

Namespaces are implemented via key prefixes, with the namespace being prepended to the key name with the delimiter `:`. For instance, a key name of `bar` with a namespace of `foo` would be written as `foo:bar`.

### Postgres

To configure, run:

```shell
prop config set url postgres://user:password@host:port/database
```

The following is the SQL schema:

```
CREATE TYPE "data_types" AS ENUM (
  'key_value',
  'list',
  'set'
);

CREATE TABLE "properties" (
  "id" SERIAL PRIMARY KEY,
  "namespace" varchar NOT NULL DEFAULT 'default',
  "data_type" data_types NOT NULL,
  "key" varchar NOT NULL,
  "value" text NOT NULL,
  "created_at" timestamp
);

CREATE INDEX "namespace_by_data_type" ON "properties" ("namespace", "data_type");

CREATE INDEX "namespace_by_key" ON "properties" ("namespace", "key");

CREATE UNIQUE INDEX ON "properties" ("id");
```

The encoding should be as follows:

- encoding: `pg_char_to_encoding('utf8')`
- datcollate: `en_US.utf8`
- datctype: `en_US.utf8`

When querying for a property, the type of the command should be compared to the type of the retrieved record. If they do not match, then the command should return an error.
