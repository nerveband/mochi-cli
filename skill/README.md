# Mochi CLI Skill for LLM Agents

This skill enables LLM agents (Claude, ChatGPT, etc.) to use mochi-cli for managing Mochi.cards flashcards and decks.

## Installation

### For Claude Code / Claude Desktop

1. **Copy the skill to your skills directory:**
   ```bash
   cp -r skill ~/.claude/skills/mochi
   ```

2. **Or use skillshare (if installed):**
   ```bash
   skillshare install https://github.com/nerveband/mochi-cli
   ```

3. **Verify installation:**
   ```bash
   ls ~/.claude/skills/mochi
   # Should show: SKILL.md
   ```

### For Other LLM Tools

Copy `SKILL.md` to your LLM's skills/prompts directory according to your tool's documentation.

## Usage

Once installed, your LLM agent will automatically use this skill when you:
- Ask to create, read, update, or delete Mochi cards or decks
- Want to manage flashcards or study materials
- Need to import/export .mochi files
- Ask about due cards for review
- Want to automate flashcard workflows

## Examples

**Create a vocabulary deck:**
```
"Create a Spanish vocabulary deck and add cards for common greetings"
```

**Export flashcards:**
```
"Export my Spanish deck to a .mochi file"
```

**Check due cards:**
```
"Show me what flashcards are due for review today"
```

**Bulk operations:**
```
"Create flashcards from this list of French vocabulary words"
```

## Requirements

- mochi-cli must be installed (see main README.md)
- Mochi API key configured (`mochi config add ...`)

## Features

The skill provides:
- ✅ Full Mochi.cards API coverage
- ✅ LLM-optimized output formats
- ✅ Batch operations and automation
- ✅ Error handling and validation
- ✅ Quiet mode for programmatic use
- ✅ Multiple account support

## Links

- **Main Repository:** https://github.com/nerveband/mochi-cli
- **Releases:** https://github.com/nerveband/mochi-cli/releases
- **Mochi.cards:** https://mochi.cards
