# LLM Quick Reference for Mochi CLI

## Getting Started

```bash
# Check version and capabilities
mochi version

# List all available commands
mochi --help

# See help for specific command group
mochi card --help
mochi deck --help
```

## Essential Patterns

### 1. Configuration

```bash
# Set up API key (do this first!)
mochi config add myprofile YOUR_API_KEY

# Verify setup
mochi deck list
```

### 2. Working with Decks

```bash
# Get all deck IDs (for iteration)
mochi deck list --id-only --quiet

# Get full deck info
mochi deck get DECK_ID --format json

# Create a new deck
mochi deck create "Deck Name" --format json

# Create nested deck
mochi deck create "Child Deck" --parent PARENT_ID
```

### 3. Working with Cards

```bash
# List cards in a deck
mochi card list --deck DECK_ID --format json

# Get specific card
mochi card get CARD_ID --format json

# Create card with content
mochi card create --deck DECK_ID --content "# Question\n\nAnswer" --format json

# Create from stdin (great for LLM output)
echo "# Generated Card\n\nContent" | mochi card create --deck DECK_ID --stdin

# Search cards
mochi card search "keyword" --format json

# Get due cards
mochi due list --format json
```

### 4. Bulk Operations

```bash
# Count cards per deck
for deck_id in $(mochi deck list --id-only --quiet); do
  count=$(mochi card list --deck $deck_id --quiet --format json | jq '.cards | length')
  echo "$deck_id: $count cards"
done

# Export all cards
mkdir -p export
for card_id in $(mochi card list --id-only --quiet); do
  mochi card get $card_id --format json > "export/${card_id}.json"
done

# Archive old cards (use with caution)
mochi card list --format json | jq -r '.cards[] | select(.archived == false) | .id' | while read id; do
  mochi card update $id --archive --quiet
done
```

### 5. Content Processing

```bash
# Extract just the content field
mochi card get CARD_ID --output-only content --quiet

# Get multiple fields
mochi card get CARD_ID --format json | jq '{id, name, content}'

# Process card content with external tools
mochi card get CARD_ID --output-only content --quiet | pandoc -f markdown -t html
```

## Output Format Guide

| Format | Use Case | Example |
|--------|----------|---------|
| `json` (default) | Full API response, LLM processing | `mochi deck list --format json` |
| `compact` | Flat array, minimal | `mochi deck list --format compact` |
| `table` | Human readable, terminal | `mochi deck list --format table` |
| `markdown` | Documentation, export | `mochi card get ID --format markdown` |

## Error Handling

```bash
# Get JSON errors for programmatic handling
mochi card get INVALID_ID --json-errors
# Output: {"error": "card not found"}

# Check exit codes
mochi card get VALID_ID --quiet
if [ $? -eq 0 ]; then
  echo "Success"
elif [ $? -eq 2 ]; then
  echo "API error"
fi
```

## Common Workflows

### Import Cards from JSON

```bash
# cards.json contains array of {content, deck_id}
jq -c '.[]' cards.json | while read card; do
  deck_id=$(echo $card | jq -r '.deck_id')
  content=$(echo $card | jq -r '.content')
  mochi card create --deck $deck_id --content "$content" --quiet
done
```

### Sync Due Cards to External System

```bash
mochi due list --format json | jq -c '.cards[]' | while read card; do
  id=$(echo $card | jq -r '.id')
  content=$(echo $card | jq -r '.content')
  # Send to your system
  curl -X POST https://your-api.com/cards \
    -d "{\"id\": \"$id\", \"content\": \"$content\"}"
done
```

### Generate Study Report

```bash
# Get stats for all decks
echo "Deck Study Report"
echo "================="
date

total_cards=0
for deck_id in $(mochi deck list --id-only --quiet); do
  deck_name=$(mochi deck get $deck_id --output-only name --quiet)
  card_count=$(mochi card list --deck $deck_id --quiet --format json | jq '.cards | length')
  due_count=$(mochi due list --deck $deck_id --quiet --format json | jq '.cards | length')
  
  echo "$deck_name: $card_count cards, $due_count due"
  total_cards=$((total_cards + card_count))
done

echo ""
echo "Total: $total_cards cards"
```

## Template Reference

### Card Object Structure

```json
{
  "id": "ABC123",
  "content": "# Card Content",
  "name": "Card Name",
  "deck-id": "DECK123",
  "template-id": "TEMPLATE123",
  "pos": "1",
  "archived?": false,
  "review-reverse?": false,
  "created-at": "2024-01-01T00:00:00Z",
  "updated-at": "2024-01-01T00:00:00Z"
}
```

### Deck Object Structure

```json
{
  "id": "DECK123",
  "name": "Deck Name",
  "parent-id": "PARENT123",
  "sort": 1,
  "archived?": false,
  "sort-by": "lexicographically",
  "cards-view": "list",
  "show-sides?": true
}
```

## Tips for LLMs

1. **Always use `--quiet`** when output is being parsed
2. **Use `--format json`** for structured data
3. **Use `--dry-run`** to test destructive operations
4. **Check version** with `mochi version` if encountering issues
5. **Use `--id-only`** to get simple lists for iteration
6. **Pipe to `jq`** for advanced JSON processing

## Rate Limiting

Mochi API has rate limits and concurrency limits (1 concurrent request per account). The CLI handles this, but for bulk operations, add delays:

```bash
for id in $(mochi card list --id-only --quiet); do
  mochi card get $id --quiet
  sleep 0.5  # Be nice to the API
done
```
