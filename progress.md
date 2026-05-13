# Branch Progress: feat/game-play

All tasks on this branch are complete. `go build ./...` passes with zero errors.

---

## Task 1: Content API `?include=` Query Parameter Support

**Status: ✅ COMPLETE**

### Goal

Add optional `?include=` query parameter to GET list and GET by ID endpoints across the content hierarchy (Lesson → Section, Unit → Lesson, QuestionSet → Unit/Lesson, Question → QuestionSet/Choices/Tags). Selective GORM preload only — no `clause.Associations`, no deep nesting, no N+1 queries.

### What Was Done

| Step | File(s) | Change |
|---|---|---|
| A | `domain/common.go` | `Includes` type, `ParseIncludes`, `.Has()` |
| A | `domain/content.go` | Pointer relation fields on `Lesson`, `Unit`, `QuestionSet`, `Question`; `Includes` added to all four Query structs; removed `IncludeChoices`/`IncludeTags` from `QuestionQuery` |
| A | `port/repository.go`, `port/service.go` | `FindByID` signatures extended with `includes domain.Includes` for Lesson, Unit, QuestionSet, Question |
| B | `repository/content/model.go` | GORM association fields on `LessonModel`, `UnitModel`, `QuestionSetModel`, `QuestionModel` |
| C | `repository/content/repository.go` | Selective `Preload` in all `FindByID` / `FindAll` methods; domain mappers updated to set pointer relation fields on non-nil associations; internal/write callers pass `nil` |
| D | 5 service files | Public `FindByID` methods accept + pass `includes`; internal/write callers pass `nil` |
| E | 5 handler `*_dto.go` files | Summary DTOs added (`SectionSummaryDTO`, `LessonSummaryDTO`, `UnitSummaryDTO`, `QuestionSetSummaryDTO`); optional relation fields on response types |
| F | 4 `*_transformer.go` files | Map non-nil domain relation pointers to summary DTOs |
| G | 4 handler files | Parse `?include=` via `domain.ParseIncludes`; pass to service |
| H | `API.md` | Include query parameter reference section; updated query param tables; breaking-change note for Question list |

### Key Decisions

| Decision | Rationale |
|---|---|
| `Includes` lives in `domain` package | Avoids import cycles — domain is the shared nucleus all layers can import |
| `FindByID` signature extended, not a new method | Clean interface; internal callers just pass `nil` |
| Internal service/write callers pass `nil` | Validation and write paths don't need relations; avoids wasteful joins |
| `Question.FindByID` always loads choices+tags | Single-item read was always full; backward compat preserved |
| `Question.FindAll` choices/tags now opt-in | Breaking but intentional — consistent with `?include=` pattern |
| `include=lesson` on QuestionSet uses `Preload("Unit.Lesson")` | GORM resolves via WHERE IN — 3 queries total for a list page, not N+1 |
| Zero-UUID sentinel `!= uuid.Nil` to detect loaded associations | No extra flags needed |
| Summary DTOs are flat | Prevents recursive/circular DTO structures |
| Preload added **after** COUNT query | COUNT must not have joins that inflate the row count |
| Lesson travels via `model.Unit.Lesson` in QuestionSet mapper | Lesson has no direct FK on QuestionSet; it arrives through the already-loaded Unit association |

### Supported Includes

| Entity | `?include=` values |
|---|---|
| Lesson | `section` |
| Unit | `lesson` |
| QuestionSet | `unit`, `lesson` |
| Question (list) | `question_set`, `choices`, `tags` |
| Question (by ID) | `question_set` (choices+tags always loaded) |

---

## Task 2: API Response Standardization

**Status: ✅ COMPLETE**

### Goal

Standardize response envelopes across all content APIs:
- Consistent `{ success, data, meta: { code } }` for success
- Consistent `{ success, error: { code, message } }` for errors
- Environment-prefixed business codes (`vocab-{env}-{code}`)
- Numeric constants for all code categories

### What Was Done

| Step | File(s) | Change |
|---|---|---|
| 1 | `internal/constants/error_codes.go` (new) | Numeric constants for all categories: 2000–2005 (success), 4001–4004 (validation), 4101–4104/4301 (auth), 4201–4208 (content), 4901–4902 (conflict), 5001–5004 (system) |
| 2 | `internal/core/helper/presenter.go` | Exported `SuccessResponse`, `ErrorResponse`, `Meta`, `APIError`; added `BuildCode`; `RespondSuccess` now variadic for backward compat; `RespondError` wraps code with `BuildCode` |
| 3 | `internal/core/helper/validator.go` | `validation_error` → `4001`, `invalid_request` → `4002`, `validator_unavailable` → `5001` |
| 4 | `internal/repository/content/repository.go` | `mapGormNotFound` for Section/Lesson/Unit/QuestionSet/Question → numeric not-found codes (4201–4205) + `5001` fallback |
| 5 | 4 content service files | Duplicate-slug codes (`duplicate_*_slug`) → `constants.CodeDuplicateSlug` (4207) |
| 6 | `common_transformer.go` + 9 handler files | `parseUUID` → `4004`; all direct UUID parse errors → `4004`; `RespondSuccess` calls tagged: Create→2001, Update→2002, Delete→2003, FindByID/FindAll→2000 |
| 7 | `API.md` | Updated Common Structures (new envelopes); added Response Codes section (format, tables, examples) |

### Key Decisions

| Decision | Rationale |
|---|---|
| `BuildCode` reads `os.Getenv("APP_ENV")` directly | Avoids import cycle from `internal/core/helper` → `configs`; env var is available at runtime without a config struct |
| `RespondSuccess` variadic code parameter | Backward compatible — identity handler callers unchanged; `meta` field is simply absent when no code is passed (`omitempty`) |
| Numeric codes scoped to content not-found and duplicate slug | These are the most user-visible and explicitly in the spec; other string codes still render correctly via `BuildCode` wrapping |
| Identity handler codes left as descriptive strings | Out of scope; they still render as `vocab-{env}-{string}` — valid format, just not numeric |
| `CodeInternalError` (5001) for `mapGormNotFound` fallback | All DB-level find failures are internal errors; the message field still distinguishes them for debugging |
| Question choice and tag not-found codes left unchanged | Not listed in the content code spec (4201–4208 range covers Section→Question only) |

### Response Code Reference

| Constant | Value | Typical HTTP |
|---|---|---|
| `CodeSuccess` | `2000` | 200 |
| `CodeCreated` | `2001` | 201 |
| `CodeUpdated` | `2002` | 200 |
| `CodeDeleted` | `2003` | 200 |
| `CodeValidationFailed` | `4001` | 400 |
| `CodeInvalidRequestBody` | `4002` | 400 |
| `CodeInvalidQueryParam` | `4004` | 400 |
| `CodeSectionNotFound` | `4201` | 404 |
| `CodeLessonNotFound` | `4202` | 404 |
| `CodeUnitNotFound` | `4203` | 404 |
| `CodeQuestionSetNotFound` | `4204` | 404 |
| `CodeQuestionNotFound` | `4205` | 404 |
| `CodeDuplicateSlug` | `4207` | 409 |
| `CodeInternalError` | `5001` | 500 |
