---
name: mochi
description: Use when the user asks to manage Mochi.cards flashcards, decks, templates, or due cards via the mochi CLI. Supports create, read, update, delete, search, import/export, and LLM-optimized workflows.
---

# Mochi CLI v1.0.0

## Installation

```bash
# One-line install (macOS/Linux)
curl -fsSL https://raw.githubusercontent.com/nerveband/mochi-cli/main/install.sh | bash

# macOS (Apple Silicon)
curl -L https://github.com/nerveband/mochi-cli/releases/latest/download/mochi-cli_Darwin_arm64.tar.gz | tar xz
sudo mv mochi /usr/local/bin/

# macOS (Intel)
curl -L https://github.com/nerveband/mochi-cli/releases/latest/download/mochi-cli_Darwin_x86_64.tar.gz | tar xz
sudo mv mochi /usr/local/bin/

# Linux (x64)
curl -L https://github.com/nerveband/mochi-cli/releases/latest/download/mochi-cli_Linux_x86_64.tar.gz | tar xz
sudo mv mochi /usr/local/bin/

# Windows (x64)
# Download mochi-cli_Windows_x86_64.zip from releases and add to PATH
```

**Repository:** https://github.com/nerveband/mochi-cli
**Releases:** https://github.com/nerveband/mochi-cli/releases
**Mochi API Docs:** https://mochi.cards/docs/api/

## Overview

Use the `mochi` CLI to manage Mochi.cards flashcards, decks, templates, and spaced repetition. Supports multiple output formats (json, table, compact, markdown) and multiple API key profiles for different accounts.

## When to Use

- User asks to create, read, update, or delete Mochi cards or decks
- User wants to search flashcards or manage vocabulary/study materials
- User asks to import or export .mochi files
- User wants to check due cards for review
- User needs to manage templates or attachments
- User wants to automate flashcard workflows or integrate with LLMs

## Setup

```bash
# Add your API key (get from Account Settings → API Keys in Mochi)
mochi config add my-profile YOUR_API_KEY

# Verify connection
mochi deck list

# Multiple profiles
mochi config add work API_KEY_1
mochi config add personal API_KEY_2
mochi config use personal
```

## Instructions

1. **Always verify profile**: Confirm which API key/account to use if multiple profiles exist.
2. **Use quiet mode for scripts**: Add `--quiet` or `-q` for programmatic use.
3. **Default to JSON format**: JSON is default and LLM-friendly. Use `--format table` for human-readable output.
4. **Prefer dry-run**: Use `--dry-run` to preview mutations before executing.
5. **Extract specific fields**: Use `--output-only <field>` or `--id-only` to get only what's needed.
6. **Pipe-friendly**: Read from stdin with `--stdin` for card content.
7. **Handle errors gracefully**: Use `--json-errors` for structured error output.

## Quick Reference

### Configuration
| Task | Command |
|---|---|
| Add profile | `mochi config add <name> <api-key>` |
| List profiles | `mochi config list` |
| Switch profile | `mochi config use <name>` |
| Remove profile | `mochi config remove <name>` |

### Decks
| Task | Command |
|---|---|
| List all decks | `mochi deck list` |
| List (table format) | `mochi deck list --format table` |
| Get deck by ID | `mochi deck get DECK_ID` |
| Create deck | `mochi deck create "Spanish Vocabulary"` |
| Create nested deck | `mochi deck create "Advanced" --parent PARENT_ID` |
| Update deck name | `mochi deck update DECK_ID --name "New Name"` |
| Archive deck | `mochi deck update DECK_ID --archive` |
| Delete deck | `mochi deck delete DECK_ID` |
| Import .mochi file | `mochi deck import vocab.mochi` |
| Export to .mochi | `mochi deck export DECK_ID vocab.mochi` |

### Cards
| Task | Command |
|---|---|
| List all cards | `mochi card list` |
| List cards in deck | `mochi card list --deck DECK_ID` |
| List (table format) | `mochi card list --format table` |
| Get card by ID | `mochi card get CARD_ID` |
| Create card | `mochi card create --deck DECK_ID --content "# Front\n\nBack"` |
| Create from file | `mochi card create --deck DECK_ID --file card.md` |
| Create from stdin | `echo "# Question\n\nAnswer" \| mochi card create --deck DECK_ID --stdin` |
| Update card content | `mochi card update CARD_ID --content "New content"` |
| Update card name | `mochi card update CARD_ID --name "Updated Name"` |
| Archive card | `mochi card update CARD_ID --archive` |
| Delete card | `mochi card delete CARD_ID` |
| Search cards | `mochi card search "keyword"` |
| Search in deck | `mochi card search "keyword" --deck DECK_ID` |

### Templates
| Task | Command |
|---|---|
| List templates | `mochi template list` |
| Get template | `mochi template get TEMPLATE_ID` |

### Due Cards
| Task | Command |
|---|---|
| List due today | `mochi due list` |
| List due on date | `mochi due list --date 2025-12-25` |
| Count due cards | `mochi due count` |
| Due in specific deck | `mochi due list --deck DECK_ID` |

### Attachments
| Task | Command |
|---|---|
| Add attachment | `mochi attachment add CARD_ID /path/to/file.png` |
| Delete attachment | `mochi attachment delete CARD_ID filename.png` |

## Output Formats

### JSON (Default - LLM-Friendly)
```bash
mochi deck list --format json
# {"decks": [...], "bookmark": "..."}
```

### Table (Human-Readable)
```bash
mochi card list --format table
# ID        CONTENT           UPDATED
# abc123    Hello / Hola      2m ago
```

### Compact (Flat Array)
```bash
mochi deck list --format compact
# [{"id":"...","name":"..."},...]
```

### Markdown
```bash
mochi card get CARD_ID --format markdown
# Returns formatted markdown
```

## LLM & Scripting Features

### Quiet Mode
Suppress all status messages, output only data:
```bash
mochi deck list --quiet
```

### Field Extraction
Get only specific fields:
```bash
# Get only IDs
mochi deck list --id-only
mochi card list --output-only id

# Extract any field
mochi card get CARD_ID --output-only content
```

### Pipe-Friendly Workflows
```bash
# Create cards from files
cat card.md | mochi card create --deck DECK_ID --stdin

# Process with jq
mochi card list --quiet | jq '.cards[].id'

# Bulk export
for id in $(mochi card list --deck DECK_ID --id-only --quiet); do
  mochi card get $id > "cards/${id}.json"
done
```

### JSON Errors
```bash
mochi card get INVALID_ID --json-errors
# {"error": "card not found"}
```

### Dry-Run Mode
Preview changes before executing:
```bash
mochi card create --deck DECK_ID --content "Test" --dry-run
# Shows what would be created without actually creating it
```

## Common Workflows

### Create Vocabulary Deck
```bash
# Create deck
DECK_ID=$(mochi deck create "Spanish Vocab" --id-only --quiet)

# Add cards
echo "# Hola\n\nHello" | mochi card create --deck $DECK_ID --stdin
echo "# Gracias\n\nThank you" | mochi card create --deck $DECK_ID --stdin
echo "# Adiós\n\nGoodbye" | mochi card create --deck $DECK_ID --stdin

# Verify
mochi card list --deck $DECK_ID --format table
```

### Export and Backup
```bash
# Export all decks
for deck_id in $(mochi deck list --id-only --quiet); do
  deck_name=$(mochi deck get $deck_id --output-only name --quiet)
  mochi deck export $deck_id "backups/${deck_name}.mochi"
done
```

### Study Due Cards
```bash
# Check what's due today
mochi due count
mochi due list --format table

# Review by deck
mochi due list --deck DECK_ID
```

### Bulk Card Creation
```bash
# From CSV or structured data
cat vocabulary.csv | while IFS=, read front back; do
  echo "# $front\n\n$back" | mochi card create --deck DECK_ID --stdin --quiet
done
```

## Environment Variables

Override API key without config:
```bash
export MOCHI_API_KEY="your-api-key"
mochi deck list  # Uses env var
```

## Error Handling

Exit codes:
- `0`: Success
- `1`: User error (invalid input, missing args)
- `2`: API error (network, auth, rate limit)
- `3`: Config error

Error categories (with `--json-errors`):
- `AUTH_ERROR`: Invalid or missing API key
- `NOT_FOUND`: Resource not found
- `RATE_LIMIT`: Too many requests
- `API_ERROR`: Server-side error
- `CONFIG_ERROR`: Configuration issue

## Configuration

Config file: `~/.mochi-cli/config.json`

Precedence:
1. CLI flags (`--api-key`)
2. Environment variable (`MOCHI_API_KEY`)
3. Profile in config file
4. Error if none found

## Self-Update

```bash
# Check for updates
mochi upgrade

# Check version
mochi version
```

## Tips for LLM Workflows

1. **Always use --quiet for scripts**: Suppresses status messages
2. **Extract only needed fields**: Use `--output-only` to reduce token usage
3. **Use JSON format**: Default format is LLM-friendly structured data
4. **Batch operations**: Loop over IDs with `--id-only --quiet`
5. **Preview mutations**: Use `--dry-run` to verify before executing
6. **Handle errors properly**: Use `--json-errors` for structured error handling

## Examples for Common LLM Tasks

### Create Study Materials from Text
```bash
# Parse text and create flashcards
echo "Parse this text and create flashcards" | llm |
while read -r card; do
  echo "$card" | mochi card create --deck DECK_ID --stdin --quiet
done
```

### Analyze Due Cards
```bash
# Get due cards and summarize
mochi due list --format json | jq '.cards[] | {id, content}' | llm "Summarize these cards"
```

### Bulk Update Cards
```bash
# Get all cards and process with LLM
for id in $(mochi card list --deck DECK_ID --id-only --quiet); do
  content=$(mochi card get $id --output-only content --quiet)
  updated=$(echo "$content" | llm "Improve this flashcard")
  echo "$updated" | mochi card update $id --stdin --quiet
done
```
