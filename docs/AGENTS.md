# docs/ — AGENTS.md

## Purpose

This folder drives **Phase 1** of the vibecode workflow: generating the foundational documents before any code is written.

## Folder Layout

```
docs/
├── 00_intake/        ← Start here: fill out the questionnaire
├── 01_prompts/       ← Prompt files to paste into any AI tool
└── 02_outputs/       ← Save AI-generated outputs here
```

## Order of Operations

1. **Fill in `00_intake/intake_questionnaire.md`** — this is the source of truth for all `[[PLACEHOLDER]]` tokens
2. **Work through `01_prompts/` in order** (01 → 06):
   - Fill placeholders from the intake answers
   - Paste the prompt into your AI tool
   - Save the output to `02_outputs/` with a matching filename (e.g. `02_outputs/01_prd.md`)
3. **Review all outputs** before starting Phase 2
4. **Seed `adr/` and `memory/`** from the architecture and implementation plan docs

## Rules for AI Tools

- Do **not** modify files in `01_prompts/` — these are prompt templates
- Place all generated documents in `02_outputs/`
- The OpenAPI spec (`04_api_spec`) in `02_outputs/` is the contract for contract tests in `tests/contract/`
- If a document is revised, save the new version to `02_outputs/` (overwrite is fine; git history tracks changes)
