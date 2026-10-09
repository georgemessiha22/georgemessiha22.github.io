# Agent rules

- `resume.schema.json` is the handwritten JSON Schema (draft-07) for `resume.yaml`,
  used by the Neovim yaml language server via the modeline on line 1 of `resume.yaml`.
- Whenever you add, rename, or remove a field in `resume.yaml`, `internal/config/load.go`
  (yaml DTOs) or `internal/model/resume.go`, update `resume.schema.json` in the same change.
  The schema uses `additionalProperties: false`, so missing fields show up as errors.
