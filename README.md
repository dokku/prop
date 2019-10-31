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

## Commands

The following commands are supported.

### `config`

Can be used to configure various components of prop.

```shell
# getting a value
prop config get backend

# setting a value
prop config set backend postgres://user:password@host:port/database

# clearing the configured value
prop config del backend
```

Current configuration values that may be set:

- `backend`: (type: string, default: `file:///etc/prop.d`) A configured backend for prop, specified in [DSN](https://en.wikipedia.org/wiki/Data_source_name) form. Backends are built into the prop project. Currently supported backends are `file` and `postgres`
- `namespace`: (type: string, default: `default`) The default namespace. Commands that allow namespace usage will note as such.

#### `del key`

- Description: Delete a key
- Data Type: `key-value`, `list`, `set`
- Supported Flags: `--namespace`

#### `key-value` commands

##### `get key [default]`

- Description: Get the value of a key
- Data Type: `key-value`
- Supported Flags: `--namespace`

##### `get-all [prefix]`

- Description: Get all key-value tuples
- Data Type: `[(key-value tuple)]`
- Supported Flags: `--namespace`

##### `set key`

- Description: Set the string value of a key
- Data Type: `key-value`
- Supported Flags: `--namespace`

#### `list` commands

##### `lindex key index`

- Description: Get an element from a list by its index
- Data Type: `list`
- Supported Flags: `--namespace`

##### `lismember key element`

- Description: Determine if a given value is an element in the list
- Data Type: `list`
- Supported Flags: `--namespace`

##### `llen key`

- Description: Get the length of a list
- Data Type: `list`
- Supported Flags: `--namespace`

##### `lrange key [start [stop]]`

- Description: Get a range of elements from a list
- Data Type: `list`
- Supported Flags: `--namespace`

##### `lrem key count element`

- Description: Remove elements from a list
- Data Type: `list`
- Supported Flags: `--namespace`

##### `lset key index element`

- Description: Set the value of an element in a list by its index
- Data Type: `list`
- Supported Flags: `--namespace`

##### `rpush key element`

- Description: Append one or more members to a list
- Data Type: `list`
- Supported Flags: `--namespace`

#### `set` commands

##### `sadd key member [member ..]`

- Description: Add one or more members to a set
- Data Type: `set`
- Supported Flags: `--namespace`

##### `sismember key member`

- Description: Determine if a given value is a member of a set
- Data Type: `set`
- Supported Flags: `--namespace`

##### `smembers key`

- Description: Get all the members in a set
- Data Type: `set`
- Supported Flags: `--namespace`

##### `srem key member [member ...]`

- Description: Remove one or more members from a set
- Data Type: `set`
- Supported Flags: `--namespace`

### Backends

The following backends are supported.

#### File

```shell
prop config set backend file:///etc/prop.d
```

#### Postgres

```shell
prop config set backend postgres://user:password@host:port/database
```
