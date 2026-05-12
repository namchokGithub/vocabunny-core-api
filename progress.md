# Refactor Progress: Content API Include Support

## Status: ✅ COMPLETE

`go build ./...` passes with zero errors. All steps A–F implemented and verified.

---

## Goal

Add optional `?include=` query parameter to GET list and GET by ID endpoints
across the content hierarchy (Lesson → Section, Unit → Lesson,
QuestionSet → Unit/Lesson, Question → QuestionSet/Choices/Tags).

Selective GORM preload only — no `clause.Associations`, no deep recursive
nesting, no N+1 queries.

---

## What Was Done (all complete)

### ✅ 1. `internal/core/domain/common.go`

Added `Includes` type (a `map[string]struct{}`), `ParseIncludes(raw string)`,
and `(Includes).Has(key string) bool`. Lives in domain so every layer can
import it without creating cycles.

### ✅ 2. `internal/core/domain/content.go`

Added optional pointer relation fields to domain structs:

| Struct | New Field |
|---|---|
| `Lesson` | `Section *Section` |
| `Unit` | `Lesson *Lesson` |
| `QuestionSet` | `Unit *Unit`, `Lesson *Lesson` |
| `Question` | `QuestionSet *QuestionSet` |

Added `Includes Includes` to `LessonQuery`, `UnitQuery`, `QuestionSetQuery`,
`QuestionQuery`. Also **removed** `IncludeChoices bool` and `IncludeTags bool`
from `QuestionQuery` — replaced by `Includes`.

### ✅ 3. `internal/core/port/repository.go` & `service.go`

Updated `FindByID` signatures for `Lesson`, `Unit`, `QuestionSet`, `Question`
in both `*Repository` and `*Service` interfaces:

```go
FindByID(ctx context.Context, id uuid.UUID, includes domain.Includes) (..., error)
```

`Section`, `QuestionChoice`, and `Tag` interfaces are **unchanged** — they
have no parent relations to include.

### ✅ 4. `internal/repository/content/model.go`

Added GORM association fields (zero-value by default; populated only when
explicitly preloaded):

| Model | New Field |
|---|---|
| `LessonModel` | `Section SectionModel \`gorm:"foreignKey:SectionID"\`` |
| `UnitModel` | `Lesson LessonModel \`gorm:"foreignKey:LessonID"\`` |
| `QuestionSetModel` | `Unit UnitModel \`gorm:"foreignKey:UnitID"\`` |
| `QuestionModel` | `QuestionSet QuestionSetModel \`gorm:"foreignKey:QuestionSetID"\`` |

`UnitModel.Lesson` is reused by GORM when preloading `"Unit.Lesson"` for
`QuestionSet` — no extra model field needed.

### ✅ 5. `internal/repository/content/repository.go`

**Preload wiring (FindByID + FindAll):**
- `lessonRepository` — `Preload("Section")` when `includes.Has("section")`
- `unitRepository` — `Preload("Lesson")` when `includes.Has("lesson")`
- `questionSetRepository` — `Preload("Unit.Lesson")` for `include=lesson`, `Preload("Unit")` for `include=unit` only
- `questionRepository.FindByID` — `Preload("QuestionSet")` for `include=question_set`; always loads choices+tags (backward compat)
- `questionRepository.FindAll` — `Preload("QuestionSet")` + choices/tags now opt-in via `Includes.Has("choices")` / `Includes.Has("tags")`
- All `Create`/`Update` internal calls pass `nil`

**Domain mappers updated:**
- `toDomainLesson` — sets `lesson.Section` when `model.Section.ID != uuid.Nil`
- `toDomainUnit` — sets `unit.Lesson` when `model.Lesson.ID != uuid.Nil`
- `toDomainQuestionSet` — sets `qs.Unit` from `model.Unit` (scalar fields only); sets `qs.Lesson` from `model.Unit.Lesson` (Lesson travels via the Unit association)
- `toDomainQuestion` — sets `question.QuestionSet` when `model.QuestionSet.ID != uuid.Nil`

### ✅ 6. Service files (5 files)

All `FindByID` call sites updated:

| File | Change |
|---|---|
| `lesson_service.go` | Internal calls → `nil`; public `FindByID` accepts + passes `includes` |
| `unit_service.go` | Internal calls → `nil`; public `FindByID` accepts + passes `includes` |
| `question_set_service.go` | Internal calls → `nil`; public `FindByID` accepts + passes `includes` |
| `question_service.go` | Internal calls → `nil`; public `FindByID` accepts + passes `includes` |
| `question_choice_service.go` | `s.questionRepository.FindByID` → passes `nil` |

### ✅ 7. DTOs — `internal/handler/content/`

Added flat summary types (no nested objects):

| File | Added |
|---|---|
| `section_dto.go` | `SectionSummaryDTO{ID, Slug, Title}` |
| `lesson_dto.go` | `LessonSummaryDTO{ID, SectionID, Slug, Title}`; `LessonResponse.Section *SectionSummaryDTO` |
| `unit_dto.go` | `UnitSummaryDTO{ID, LessonID, Slug, Title}`; `UnitResponse.Lesson *LessonSummaryDTO` |
| `question_set_dto.go` | `QuestionSetSummaryDTO{ID, UnitID, Slug, Title, Version}`; `QuestionSetResponse.Unit`, `.Lesson` |
| `question_dto.go` | `QuestionResponse.QuestionSet *QuestionSetSummaryDTO` |

### ✅ 8. Transformers — `internal/handler/content/`

All four transformers updated to map non-nil domain relation pointers to their
summary DTOs:

- `lesson_transformer.go` → `toLessonResponse` maps `item.Section`
- `unit_transformer.go` → `toUnitResponse` maps `item.Lesson`
- `question_set_transformer.go` → `toQuestionSetResponse` maps `item.Unit` and `item.Lesson`
- `question_transformer.go` → `toQuestionResponse` maps `item.QuestionSet`

### ✅ 9. Handlers — `internal/handler/content/`

All four handlers updated to parse `?include=` and pass it through:

```go
// FindByID
includes := domain.ParseIncludes(c.QueryParam("include"))
item, err := h.service.FindByID(c.Request().Context(), id, includes)

// FindAll
query.Includes = domain.ParseIncludes(c.QueryParam("include"))
```

`question_handler.go`: removed old `include_choices` / `include_tags` boolean
query params — replaced by `?include=choices,tags`.

### ✅ 10. `API.md`

- Added **Content — Include Query Parameter** reference section: syntax,
  supported values per entity, response shape example, performance notes.
- Added GET by ID response examples (plain + with `?include=`) for Lesson,
  Unit, QuestionSet, and Question.
- Updated query params tables for all four entities: added `include`, `search`,
  and `page / limit / sort_by / sort_order` rows.
- Replaced `include_choices` / `include_tags` rows in Questions table with
  single `include` row plus breaking-change note.

### ✅ 11. `AGENTS.md` (new file)

Created a comprehensive agent guide covering project overview, architecture,
the full refactor spec (steps A–F with code patterns), key architecture
decisions, common patterns, and build/test commands.

---

## Key Architecture Decisions

| Decision | Rationale |
|---|---|
| `Includes` lives in `domain` package | Avoids import cycles; domain is the shared nucleus |
| `FindByID` signature extended (not a new method) | Clean interface; internal callers just pass `nil` |
| Internal service calls pass `nil` includes | Validation and write paths don't need relations; avoids wasteful joins |
| `Question.FindByID` always loads choices+tags | Single-item read was always full; backward compat preserved |
| `Question.FindAll` choices/tags now opt-in | Consistent with `include=` pattern; was default-true before (intentional breaking change) |
| `include=lesson` on QuestionSet uses `Preload("Unit.Lesson")` | GORM resolves via WHERE IN — 3 queries total for a list page, not N+1 |
| Zero-UUID check to detect loaded associations | `if model.Section.ID != uuid.Nil` — no extra flags needed |
| Summary DTOs are flat | Prevents recursive/circular DTO structures; only id+slug+title-level data on relations |
| Preload added **after** COUNT query | COUNT must not have joins that inflate the row count |
| `toDomainQuestionSet` reads Lesson via `model.Unit.Lesson` | Lesson has no direct FK on QuestionSet; it travels through the Unit association already loaded by `Preload("Unit.Lesson")` |

---

## Supported Includes Reference

| Entity | `?include=` values | Notes |
|---|---|---|
| Lesson | `section` | Embeds parent Section summary |
| Unit | `lesson` | Embeds parent Lesson summary |
| QuestionSet | `unit`, `lesson` | `lesson` alone resolves via `Unit.Lesson` |
| Question (list) | `question_set`, `choices`, `tags` | choices + tags opt-in on list |
| Question (by ID) | `question_set` | choices + tags always loaded on single GET |
