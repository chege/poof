# Analysis: TOML Configuration, Autonomous & Safe Inline Masking

## Architectural Refinements

1. **Strict TOML Parsing**:
   - I will use `toml.NewDecoder(r).DisallowUnknownFields()` from `github.com/BurntSushi/toml` to ensure the tool fails hard on any unrecognized configuration keys.
   - The configuration model will be flattened where possible to map cleanly to the SRS TOML structure.

2. **Single Database Focus**:
   - The `Config` struct will be updated to include a `Database` section with `DSN`.
   - This simplifies the CLI as users no longer *need* to provide `--db` if it's in the config, though the flag may still act as an override.

3. **In-Place Mutation Safety**:
   - The `masker.Engine` already uses transactions. I will ensure that in `DryRun` mode, the transaction is *always* rolled back, even if no changes were attempted, as a double-safety measure.
   - The engine will be updated to fetch the database name *before* starting work to verify against `allowed_db_names`.

4. **Autonomous Mode**:
   - All `ui.Error` calls followed by `os.Exit(1)` are already non-interactive. I will ensure that no new code adds `fmt.Scan` or similar.
   - The `--yes` flag will be documented as the way to skip the plan review in non-TTY environments.

## Applied Suggestions (Automatic)

- **Dependency Cleanup**: I will remove `github.com/hashicorp/hcl/v2` from `go.mod`.
- **Default Path**: `root.go` will look for `dbmask.toml` if no path is provided.

## Risks & Mitigations

- **Risk**: Users losing data due to in-place mutation.
- **Mitigation**: Mandatory `dbmask plan` output before apply (unless `--yes` is used) and highly visible dry-run mode.
- **Risk**: Migrating existing users from HCL.
- **Mitigation**: This is a breaking change (Version 6.0). I will provide a clear `init` command to help users recreate their configs in TOML.
