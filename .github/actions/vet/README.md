# ock/vet

Validates document metadata according to a defined schema.

## Usage

> [!IMPORTANT]
> A schema must be present for validation to occur. See [here](../../../README.md)
> for information on how initialize a schema and the [schema reference](../../../docs/references/schema.md)
> for information on how to define it.

Specify this action in a GitHub Actions workflow:

```yaml
steps:
- uses: slewiskelly/ock/.github/actions/vet@v0
  with:
    version: latest # Default.
```

### Configuration

| Input     | Required? | Default     | Description                                                                                                                  |
|-----------|-----------|----------   |------------------------------------------------------------------------------------------------------------------------------|
| `format`  | No        | `summary`   | Display format (`json` \| `summary`)                                                                                         |
| `glob`    | No        | `**/*.md`   | Pattern to filter files                                                                                                      |
| `level`   | No        | `error`     | Minimum error level to display (`error` \| `warn`)                                                                           |
| `path`    | No        | `.`         | Root directory containing files to validate, all subdirectories within the root directory are traversed                      |
| `schema`  | No        | `./.schema` | Location of the schema file to validate against                                                                              |
| `version` | No        | `latest`    | Version of `ock` to be installed; can be either: a semantic version (`v0`, `v0.1.0`), branch (`main`), or `latest`           |
