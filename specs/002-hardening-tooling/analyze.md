# Analysis: hardening-tooling

## Architectural Refinements

1. **Strict HCL Validation**:
   - `hclsimple.DecodeFile` is convenient but doesn't easily allow for "fail hard on unknown fields" without additional logic.
   - I will transition to using the lower-level `hcl` package. I'll parse the file into an `hcl.File`, then use `hcl.Body.Content` with a defined schema. Any remaining attributes or blocks in the body will trigger a diagnostic error.

2. **Deterministic Provider Implementations**:
   - `ipv4_address`: Use the row-seeded `rand.Rand` to generate 4 bytes.
   - `short_text`: Generate a fixed-length string of random characters from a restricted set (a-z, A-Z, 0-9) using the row-seeded `rand.Rand`.
   - `username`: Combine fake first and last names or use a set of deterministic patterns.
   - `company_name`: Use a deterministic list of prefixes/suffixes or predefined names.

3. **Taskfile Orchestration**:
   - The `Taskfile.yml` will be the single source of truth for all workflows.
   - `task check` will be a prerequisite for all commits.
   - `init`, `plan`, `apply` will map to `./poof` commands with standard flags to ensure consistency.

4. **Error Hardening**:
   - Wrap all `pgx` and `hcl` errors with context-specific messages (e.g., `"failed to update column %s in table %s: %w"`).
   - Ensure the tool provides a non-zero exit code on *any* diagnostic failure or runtime error.

## Applied Suggestions (Automatic)

- **Dead Code**: I will scan for unused exported functions or internal helpers and remove them.
- **Package Boundaries**: Ensure `internal/config` does not depend on `internal/generator` to maintain a clean dependency graph.
- **Boilerplate**: Standardize all Go files with proper headers if missing.

## Risks & Mitigations

- **Risk**: `Taskfile.yml` might be complex to maintain.
- **Mitigation**: Keep tasks simple and atomic. Avoid logic inside tasks; prefer calling Go commands or simple shell scripts if absolutely necessary (though SRS says avoid scripts).
