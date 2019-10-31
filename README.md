# prop

A golang-based cli interface for manipulating config backed by various datastores.

## Summary

This describes a utility that can be used to store, retrieve, and manipulate various data types against a configured backing store. It can be used as both a library and wrapped via command-line calls.

## Motivation

When writing a command-line utility, it may be necessary to store configuration values on disk. This is typically implemented in a one-off manner, and state is both not always persisted in the same locations nor is it easily distributed for high-availability.

This project seeks to make it easy to manipulate various basic data types across a wide variety of backends via a single entrypoint, making it easier for applications to consider configuration state without needing to reimplement the wheel.

## Terminology

- Set: A mathematical, well-defined collection of distinct objects. See [wikipedia](https://en.wikipedia.org/wiki/Set_(mathematics)) for more details.
- List: An enumerated collection of objects in which repititions are allowed. See [wikipedia](https://en.wikipedia.org/wiki/Sequence) for more details.
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

## Key and Value specification

Keys may follow the following regex:

```
[\w-/]{1,200}[^/]
```

Values may contain 0 or more utf8 characters and may be a maximum of 65535 characters in length.

## Commands

The following commands are supported.

### `backend` commands

#### `backend migrate backend_dsn`

- Description: Migrate from one backend to another

When migrating a backend, it is assumed that there are is no concurrent access to the backend. In other words, if another process is changing values of either backend, then the migration may result in an invalid state.

#### `backend reset`

- Description: Clear all values in a backend

### `del key`

- Description: Delete a key
- Data Type: `key-value`, `list`, `set`
- Supported Flags: `--namespace`

### `config` commands

Used for configuring `prop`.

Current configuration values that may be manipulated:

- `backend`: (type: string, default: `file:///etc/prop.d`) A configured backend for prop, specified in [DSN](https://en.wikipedia.org/wiki/Data_source_name) form. Backends are built into the prop project. Currently supported backends are `file` and `postgres`
- `namespace`: (type: string, default: `default`) The default namespace. Commands that allow namespace usage will note as such.

#### `config get key`

- Description: Get a configuration value

```shell
prop config get backend
```

#### `config set key value`

- Description: Set a configuration value

```shell
prop config set backend postgres://user:password@host:port/database
```

#### `config del key`

- Description: Delete a configuration value

```
prop config del backend
```

### `namespace` commands

#### `namespace exists namespace`

- Description: Checks if there are any keys in a given namespace

#### `namespace clear namespace`

- Description: Delete all keys from a given namespace

### `key-value` commands

#### `get key [default]`

- Description: Get the value of a key
- Data Type: `key-value`
- Supported Flags: `--namespace`

#### `get-all [prefix]`

- Description: Get all key-value tuples
- Data Type: `[(key-value tuple)]`
- Supported Flags: `--namespace`

#### `set key`

- Description: Set the string value of a key
- Data Type: `key-value`
- Supported Flags: `--namespace`

### `list` commands

#### `lindex key index`

- Description: Get an element from a list by its index
- Data Type: `list`
- Supported Flags: `--namespace`

#### `lismember key element`

- Description: Determine if a given value is an element in the list
- Data Type: `list`
- Supported Flags: `--namespace`

#### `llen key`

- Description: Get the length of a list
- Data Type: `list`
- Supported Flags: `--namespace`

#### `lrange key [start [stop]]`

- Description: Get a range of elements from a list
- Data Type: `list`
- Supported Flags: `--namespace`

#### `lrem key count element`

- Description: Remove elements from a list
- Data Type: `list`
- Supported Flags: `--namespace`

#### `lset key index element`

- Description: Set the value of an element in a list by its index
- Data Type: `list`
- Supported Flags: `--namespace`

#### `rpush key element`

- Description: Append one or more members to a list
- Data Type: `list`
- Supported Flags: `--namespace`

### `set` commands

#### `sadd key member [member ..]`

- Description: Add one or more members to a set
- Data Type: `set`
- Supported Flags: `--namespace`

#### `sismember key member`

- Description: Determine if a given value is a member of a set
- Data Type: `set`
- Supported Flags: `--namespace`

#### `smembers key`

- Description: Get all the members in a set
- Data Type: `set`
- Supported Flags: `--namespace`

#### `srem key member [member ...]`

- Description: Remove one or more members from a set
- Data Type: `set`
- Supported Flags: `--namespace`

## Backends

The following backends are supported.

### File

To configure, run:

```shell
prop config set backend file:///etc/prop.d
```

The directory structure is as follows:

```shell
# returns the contents of the file
cat $NAMESPACE/$KEY
```

Key names can include forward slashes, which will be interpreted as a directory structure.

When querying for a property, there is no guarantee that the value will be of the type expected. As such, users should take care to always interact with a given key using the correct commands.

> For consideration: Should we serialize the value into json, such that we have something like: `{"type": "[key-value,list,set]" "value": "value here"}`? This would allow the interface to introspect on the type correctly, though at the cost of complicating the backend a bit more.

### Redis

To configure, run:

```shell
prop config set backend redis://user:password@host:port/database
```

With the redis backend, commands map to their redis equivalents where appropriate. If there is no equivalent redis command, a redis script may be used instead to implement the functionality.

When querying for a property, if the type of the value does not match the type specified by the executed command, an error should be raised where possible.

Namespaces are implemented via key prefixes, with the namespace being prepended to the key name with the delimiter `:`. For instance, a key name of `bar` with a namespace of `foo` would be written as `foo:bar`.

### Postgres

To configure, run:

```shell
prop config set backend postgres://user:password@host:port/database
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
