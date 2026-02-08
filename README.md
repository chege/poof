# ✨ poof

[![Go Version](https://img.shields.io/github/go-mod/go-version/christopher/masker?color=00ADD8&logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Safe by Default](https://img.shields.io/badge/Safety-First-blueviolet)](https://github.com/christopher/masker)

**Deterministic, parallel-safe, and declarative data masking for PostgreSQL.**

`poof` is a trust-first anonymization tool designed to help developers and operators mask sensitive production data for use in staging, testing, and local environments—without the risk of breaking constraints or leaking data.

---

## 🚀 Try poof in 2 Minutes

Experience `poof` immediately using our self-contained demo environment.

**Requirements:** [Docker](https://www.docker.com/) & [Go](https://go.dev/)

```bash
# 1. Clone the repo
git clone https://github.com/christopher/masker.git && cd masker

# 2. Run the automated demo
task demo
```

**What happens next?**
1. A local PostgreSQL container spins up.
2. Sample PII (Personally Identifiable Information) is loaded.
3. `poof` analyzes the schema and generates a **Masking Plan**.
4. You'll see a side-by-side diff of how the data *would* be transformed.

---

## 💎 Key Features

*   **🔒 Safe-by-Default**: Refuses to touch any database not explicitly named in your `allowed_db_names` allowlist.
*   **🧠 Deterministic**: Masking values are seeded by `MD5(table_name + primary_key)`, ensuring consistent results across multiple runs and environments.
*   **⚡ Parallel-Safe**: High-performance worker pool architecture scales to millions of rows while maintaining deterministic integrity.
*   **🛠 Declarative**: Define your rules in human-readable TOML. No complex SQL scripts or brittle ETL pipelines.
*   **🔐 Secret Management**: Supports `${ENV_VAR}` expansion in configuration files to keep credentials out of version control.
*   **🚀 High-Performance Batching**: Uses PostgreSQL's `UPDATE ... FROM (VALUES ...)` to apply changes in bulk (default 500 rows/batch), providing 10x-50x speedups.
*   **♻️ Smart Retries**: Automatically detects `UNIQUE` constraint violations and retries generation with a deterministic incrementer until successful.
*   **🛡 No-Schema-Mod**: `poof` never runs `ALTER`, `DROP`, or `DISABLE CONSTRAINTS`. It works within your existing schema rules.

---

## 🕹 Commands

| Command | Purpose |
| :--- | :--- |
| `poof init` | Scaffold a new `poof.toml` configuration file. |
| `poof analyze` | **Advisory Mode.** Inspects your DB and suggests columns that need masking. Supports `--json`. |
| `poof validate` | Performs deep semantic and schema validation of your configuration. |
| `poof plan` | Dry-run preview. Shows exactly what will happen without changing a single byte. |
| `poof apply` | Executes the masking. Supports `--dry-run` for a full simulation. |
| `poof doctor` | Checks your environment, DB connections, and config health. |

---

## 📝 Configuration

Define your masking rules in `poof.toml`:

```toml
[databases.staging]
dsn = "postgres://poof_user:${STAGING_DB_PASS}@staging-db:5432/app"

[databases.local]
dsn = "postgres://localhost:5432/dev_db"
default = true

[safety]
allowed_db_names = ["staging_db", "dev_db"]
salt = "your-secret-global-salt"
```

[[tables]]
name = "users"
pk = "id"

  [[tables.columns]]
  name = "email"
  [tables.columns.gen]
  type = "faker"
  provider = "email"

  [[tables.columns]]
  name = "full_name"
  [tables.columns.gen]
  type = "faker"
  provider = "full_name"

  [[tables.columns]]
  name = "secret_key"
  [tables.columns.gen]
  type = "hash"
```

---

## 🧩 Built-in Generators

`poof` comes packed with specialized generators:

*   **`faker`**: Names, Emails, Phones, Addresses, Company names, IPs.
*   **`hash`**: Deterministic MD5 hashes for tokens or IDs.
*   **`counter`**: Incremental numbers for sequencing.
*   **`template`**: Dynamic strings using Go templates.
*   **`constant`**: Fixed values for status flags or static text.
*   **`null`**: For clearing out sensitive optional fields.

---

## 🛠 Development

We use [Task](https://taskfile.dev/) for a "boring" and predictable development workflow:

*   `task build`: Compile the binary.
*   `task fmt`: Format code and sort imports.
*   `task lint`: Run comprehensive static analysis.
*   `task test`: Run unit and integration tests (uses Testcontainers).
*   `task ready`: The full quality gate (fmt + lint + test).

---

## ⚖️ License

Distributed under the MIT License. See `LICENSE` for more information.

---
*Built with ❤️ for the PostgreSQL community.*
