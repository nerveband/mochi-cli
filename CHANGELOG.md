# Changelog

All notable changes to mochi-cli will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-02-02

### Added
- Complete Mochi.cards API implementation
- Card management (create, read, update, delete, search)
- Deck management (create, read, update, delete, nested hierarchies)
- Template operations (list, get)
- Due cards queries
- Attachment management
- Multi-profile configuration support
- Import/export .mochi files
- Multiple output formats (JSON, Table, Compact, Markdown)
- LLM-optimized features:
  - Quiet mode (`--quiet`)
  - JSON error output (`--json-errors`)
  - Field extraction (`--output-only`, `--id-only`)
  - Stdin support for piping content
- Shell completions (Bash, Zsh, Fish, PowerShell)
- Dry-run mode for previewing changes
- Auto-updater with daily update checks
- Demo video and GIF showcasing features

### Documentation
- Comprehensive README with usage examples
- LLM integration guides
- Installation instructions for all platforms
- API reference and command documentation

[1.0.0]: https://github.com/nerveband/mochi-cli/releases/tag/v1.0.0
