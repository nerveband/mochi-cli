# LLM Integration Guide for Mochi CLI

## Best Practices for LLM Workflows

### 1. Configuration Management

**Store API keys securely:**
```bash
# Option 1: Environment variable (recommended for CI/automation)
export MOCHI_API_KEY="your-api-key"

# Option 2: Profile (recommended for interactive use)
mochi config add production YOUR_API_KEY
mochi config use production
```

**Verify configuration:**
```bash
# Check which profile is active
mochi config list

# Test connectivity
mochi deck list --quiet
```

### 2. Output Optimization

**For LLM consumption, always use:**
```bash
mochi <command> --format json --quiet
```

**This provides:**
- Structured, parseable output
- No status messages or progress indicators
- Consistent field names
- Complete data (no truncation)

### 3. Pagination Handling

Mochi API uses bookmarks for pagination:

```bash
# Get first page
result=$(mochi card list --deck DECK_ID --limit 100 --format json)
cards=$(echo $result | jq '.cards')
bookmark=$(echo $result | jq -r '.bookmark')

# Get subsequent pages
while [ "$bookmark" != "null" ] && [ -n "$bookmark" ]; do
  result=$(mochi card list --deck DECK_ID --limit 100 --format json)
  cards=$(echo $cards | jq --argjson new "$(echo $result | jq '.cards')" '. + $new')
  bookmark=$(echo $result | jq -r '.bookmark')
done
```

### 4. Error Handling

**Robust error handling:**
```bash
# Use --json-errors for parseable errors
output=$(mochi card get INVALID_ID --json-errors 2>&1)
if echo "$output" | jq -e '.error' > /dev/null 2>&1; then
  error_msg=$(echo "$output" | jq -r '.error')
  echo "Error: $error_msg"
  exit 1
fi
```

**Check exit codes:**
- `0`: Success
- `1`: User/config error
- `2`: API/network error
- `3`: Internal error

### 5. Content Processing

**Creating cards from LLM output:**
```bash
# Generate content with LLM
llm_output=$(cat <<'EOF'
# Spanish Vocabulary

**Word:** Hola
**Meaning:** Hello
**Example:** ¡Hola! ¿Cómo estás?
EOF
)

# Create card
echo "$llm_output" | mochi card create --deck SPANISH_DECK_ID --stdin --quiet
```

**Processing card content:**
```bash
# Get content and process
content=$(mochi card get CARD_ID --output-only content --quiet)

# Extract front/back for flashcard apps
front=$(echo "$content" | head -1)
back=$(echo "$content" | tail -n +2)
```

### 6. Batch Operations

**Creating multiple cards:**
```bash
# Read cards from JSON array
cat <<'EOF' | jq -c '.[]' | while read card; do
  deck_id=$(echo "$card" | jq -r '.deck_id')
  content=$(echo "$card" | jq -r '.content')
  name=$(echo "$card" | jq -r '.name // empty')
  
  if [ -n "$name" ]; then
    mochi card create --deck "$deck_id" --name "$name" --content "$content" --quiet
  else
    mochi card create --deck "$deck_id" --content "$content" --quiet
  fi
done
EOF
```

**Updating cards in batch:**
```bash
# Archive all cards matching pattern
mochi card search "OLD_TOPIC" --format json | \
  jq -r '.cards[].id' | \
  while read id; do
    mochi card update "$id" --archive --quiet
    sleep 0.5  # Rate limiting
  done
```

### 7. Data Export

**Full backup:**
```bash
mkdir -p mochi-backup-$(date +%Y%m%d)
cd mochi-backup-$(date +%Y%m%d)

# Export decks
mochi deck list --format json > decks.json

# Export cards per deck
mkdir cards
for deck_id in $(mochi deck list --id-only --quiet); do
  mochi card list --deck "$deck_id" --format json > "cards/${deck_id}.json"
done

# Export templates
mochi template list --format json > templates.json
```

**Export to CSV:**
```bash
# Convert cards to CSV
mochi card list --deck DECK_ID --format json | \
  jq -r '.cards[] | [.id, .name, .content] | @csv' > cards.csv
```

### 8. Integration Patterns

**With other CLI tools:**
```bash
# Sync with Anki (using anki-cli)
mochi card list --deck DECK_ID --format json | \
  jq -r '.cards[] | "\(.id)\t\(.content)"' | \
  while IFS=$'\t' read -r id content; do
    anki add --deck "Mochi Import" --front "$id" --back "$content"
  done

# Process with Python
mochi deck list --format json | python3 -c "
import sys, json
decks = json.load(sys.stdin)['decks']
for deck in decks:
    print(f'{deck[\"id\"]}: {deck[\"name\"]} ({deck.get(\"sort\", 0)})')
"
```

**Webhook integration:**
```bash
# Send due cards to webhook
mochi due list --format json | \
  jq -c '.cards[]' | \
  while read card; do
    curl -X POST https://hooks.example.com/mochi \
      -H "Content-Type: application/json" \
      -d "$card"
  done
```

### 9. Testing and Validation

**Dry-run mode:**
```bash
# Test card creation
mochi card create --deck TEST_DECK --content "Test" --dry-run
# Shows what would happen without making changes

# Validate JSON before sending
echo "Invalid JSON" | mochi card create --deck DECK --stdin --dry-run
```

**Health checks:**
```bash
#!/bin/bash
# health-check.sh

if ! mochi deck list --quiet > /dev/null 2>&1; then
  echo "ERROR: Cannot connect to Mochi API"
  exit 1
fi

due_count=$(mochi due count --quiet --format json | jq '.count')
if [ "$due_count" -gt 100 ]; then
  echo "WARNING: $due_count cards due for review"
fi

echo "OK: Mochi CLI is healthy"
```

### 10. Performance Optimization

**Minimize API calls:**
```bash
# BAD: Multiple calls
for id in $(mochi deck list --id-only --quiet); do
  mochi deck get $id  # N calls for N decks
  mochi card list --deck $id  # N more calls
  mochi due list --deck $id  # N more calls
done

# GOOD: Bulk operations where possible
mochi deck list --format json  # 1 call for all decks
mochi card list --limit 100  # Paginate efficiently
```

**Use appropriate limits:**
```bash
# Get exactly what you need
mochi card list --deck DECK --limit 10  # Small sample
mochi card list --deck DECK --limit 100  # Full page
```

### 11. Template Workflows

**Inspect template structure:**
```bash
# Get template details
mochi template get TEMPLATE_ID --format json | jq '{
  name: .name,
  fields: [.fields[] | {id, name, type}]
}'
```

**Create cards with templates:**
```bash
# Create vocabulary card from template
mochi card create \
  --deck VOCAB_DECK \
  --template TEMPLATE_ID \
  --content "Word: ${word}\nDefinition: ${definition}" \
  --format json
```

### 12. Automation Examples

**Daily study report:**
```bash
#!/bin/bash
# daily-report.sh

echo "📚 Mochi Study Report - $(date)"
echo ""

total_due=$(mochi due count --quiet)
echo "Total due today: $total_due cards"

if [ "$total_due" -gt 0 ]; then
  echo ""
  echo "Breakdown by deck:"
  mochi deck list --format json | \
    jq -r '.decks[].id' | \
    while read deck_id; do
      count=$(mochi due count --deck $deck_id --quiet)
      if [ "$count" -gt 0 ]; then
        name=$(mochi deck get $deck_id --output-only name --quiet)
        echo "  - $name: $count cards"
      fi
    done
fi
```

**Card creation from todo list:**
```bash
#!/bin/bash
# todo-to-cards.sh

TODO_DECK="YOUR_TODO_DECK_ID"

# Convert each todo to a card
cat todos.txt | while read todo; do
  content="# Todo\n\n${todo}\n\n- [ ] Complete"
  mochi card create --deck "$TODO_DECK" --content "$content" --quiet
  echo "Created card for: $todo"
done
```

## Troubleshooting

### Common Issues

**1. Authentication errors:**
```bash
# Check API key
mochi config list

# Test with explicit key
mochi deck list --api-key YOUR_KEY
```

**2. Rate limiting:**
```bash
# Add delays between requests
sleep 1

# Use --quiet to reduce output overhead
```

**3. Large content:**
```bash
# Split large content into chunks
# Use --file instead of --content for long text
```

**4. Unicode issues:**
```bash
# Ensure UTF-8 locale
export LANG=en_US.UTF-8
```
