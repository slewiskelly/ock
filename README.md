# Ock

A tool for querying and validating document metadata.

> [!WARNING]
> ___This is a work in progress and is considered highly experimental.___

## Overview

> [!NOTE]
> Ock only targets YAML frontmatter in markdown files.

### Querying

Documents can be listed as well as metadata retrieved.

Expressions can be used to filter files or fields, as necessary. These
expressions are expressed using [CUE](https://cuelang.org/) syntax.

### Validation

Document metadata is validated against a user-defined schema, expressed in
[CUE](https://cuelang.org/).

By default the definition used for validation is `#Metadata` within the
`.schema.cue` file.

See the usage section for initializing a default schema.

## Usage

> [!TIP]
> For full usage and examples use `ock help` or `ock <command> --help`.

### Initialization

To initialize a default schema file (`.schema.cue`) in the current working
directory:

```shell
ock init
```

With the default contents being:

```cue
import "time"

#Metadata: {
        title:    string
        status:   #Status
        owner:    #Owner
        reviewed: #Date
        tags: [...string]
}

#Status: "archived" | "draft" | "published"

#Date: time.Format(time.RFC3339Date)

#Owner: string
```

This schema can be changed according to specific needs.

> [!NOTE]
> The `#Metadata` definition is (___always___) used to validate files against.

### Querying

To retrieve the metadata of a single file:

```shell
ock get [flags] <path>
```

To list all (markdown) files rooted at the given path:

```shell
ock list [flags] <path>
```

### Validation

To validate file(s), rooted at the given path, against the schema:

```shell
ock vet [flags] <path>
```
