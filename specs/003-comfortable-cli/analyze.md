# Analysis: Comfortable CLI

## Architectural Refinements

1. **TTY-Aware UI Package**:
   - I will create `internal/ui` to encapsulate all terminal interaction.
   - It will handle symbol mapping (✓, !, ✗) and colorization.
   - It will automatically disable colors if `!isatty.IsTerminal` or if `--no-color` is set.

2. **Doctor Command Implementation**:
   - The `doctor` logic will live in `internal/poof` as a series of check functions.
   - It will reuse the `config.LoadConfig` and `db.Client` logic but strictly in read-only mode.

3. **Plan (Dry-Run) Capability**:
   - The `poof.Engine` will be updated to handle a `DryRun` flag.
   - In `DryRun`, it will fetch rows, generate masked values, but *not* execute the `UPDATE` statements.
   - It will return a slice of "Diff" objects containing `PK`, `Column`, `OldValue`, and `NewValue`.

4. **Self-Documenting Init**:
   - `init` will use a hardcoded template with placeholders.
   - `--explain` will inject more verbose comments into the template before writing.

## Applied Suggestions (Automatic)

- **Cobra Flags**: I will move common flags like `--db` and `--config` to `PersistentFlags` on the root command if appropriate, though some commands (like `init`) might not need all of them. For simplicity and following the SRS, I'll keep them consistent.
- **Safety**: `validate` and `doctor` will be used as internal gates inside `plan` and `apply`.

## Risks & Mitigations

- **Risk**: `plan` diffing large tables could be slow or memory-intensive.
- **Mitigation**: Strictly limit the number of rows processed in `plan` mode (e.g., first 5 rows per table).
- **Risk**: Color dependencies.
- **Mitigation**: Use well-maintained libraries (`fatih/color`) and ensure robust TTY detection.
