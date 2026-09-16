# Planned

Roadmap for Punch: features, config, commands, and infrastructure that are
decided but not yet implemented.

## Features

- [ ] Increase detail of performance summaries:
  - [ ] Report total successful, fatal, and failed requests
  - [ ] Include latencies: P99, P95, and P50
- [ ] Configuration (`punch.json`):
  - [ ] Add global request timeouts and per child overrides
  - [ ] Add request headers and request bodies
  - [ ] See [punch-planned.jsonc](../punch-planned.jsonc) for the planned layout
- [ ] Implement remaining HTTP methods
- [ ] Allow multiple separate tests within configuration file (single object->array)
- [ ] Integrate CD with GitHub Actions
- [ ] Global Punch configuration (`~/.punch`):

  ```bash
  ~/.punch
  ├── config/
  ├── tests/
      ├── bow/
      │   ├── logs/
      │   │   ├── metadata.json    # includes: dates, config, etc.
      │   │   ├── bow-test-id-logs-1.json
      │   │   ├── bow-test-id-logs-1.csv
      │   ├── results/
      │   │   ├── bow-test-id-1-results.json
      │   │   └── bow-test-id-1-results.csv
      └── lilify/
          ├── logs/
          └── results/
  ```

- [ ] Persist test outputs:
  - [ ] Serve entire test output as plain text, log, JSON and CSV files
  - [ ] Optionally persist test output to `~/.punch`:
    - [ ] Logs
    - [ ] Results:
      - [ ] .json (save must be included in config file or via the CLI)
- [ ] Improve test output formatting and structure
- [ ] Improve error reports
- [ ] Containerize builds with Docker (include build step for CI)

## Commands

- [ ] Run a load test and redirect Punch's output to a file. File
      support will include `.json`, `.log`, `.csv`, and `.txt` files:

  ```bash
  # Layout
  punch run --out <file-name-with-ext> <directory>
  ```

- [ ] Run a load test and persist output to `~/.punch`:

  ```bash
  punch run --save [directory]
  ```

- [ ] List all persisted logs with their test IDs, names, and dates:

  ```bash
  punch logs list
  ```

- [ ] View a saved test's output from `~/.punch`. If `[test-name]` is
      unspecified, Punch outputs last test results if they were set
      to be persisted:

  ```bash
  punch logs [test-name]
  ```

  The `logs` command only returns logs from test runs that were
  started with `--save-logs` or specified in a punch.json.

## Infrastructure

- [ ] Containerization: Docker
- [ ] CD: GitHub Actions

## Configuration

See:

- Configuration Sample: [punch-planned.json](../punch.json)
- Configuration Details: [punch-planned.jsonc](../punch.jsonc)

---

- [ ] Move to a top-level `tests` array so one Punch file can hold
      multiple, independently named test definitions (`global_name` /
      `global_timeout_secs` at the file level; `name` / `desc` /
      `base_url` per test)
- [ ] Add a global timeout (`global_timeout_secs`) with per-child
      override (`override_timeout_secs`)
- [ ] Rename `target` to `base_url` per test
- [ ] Rename child `name` to `path`
- [ ] Rename `total_requests` + `base_duration_secs` to a rate-based
      `requests_per_sec` + `duration_secs`
- [ ] Rename `want_status_code` to `expected_status_code`
- [ ] Add request bodies per child (`body`)
- [ ] Add `desc` fields for documenting tests and the overall file
- [ ] Drop `grace_period_percent` (superseded by per-test timeouts)
