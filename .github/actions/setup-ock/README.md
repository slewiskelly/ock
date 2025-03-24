# ock/setup-ock

Installs `ock` and adds it to the `PATH`.

## Usage

Specify this action in a GitHub Actions workflow:

```yaml
steps:
- uses: slewiskelly/ock/.github/actions/setup-ock@v0
  with:
    version: latest # Default.
```

### Configuration

| Input     | Required? | Default  | Description                                                                                                                  |
|-----------|-----------|----------|------------------------------------------------------------------------------------------------------------------------------|
| `version` | No        | `latest` | Version of `ock` to be installed; can be either: a semantic version (`v0`, `v0.1.0`), branch (`main`), or `latest`           |
