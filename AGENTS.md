# AGENTS.md — VocabBunny Core API

## Project Overview

Go REST API for VocabBunny language-learning content management. Exposes CRUD endpoints for a content hierarchy: **Section → Lesson → Unit → QuestionSet → Question → QuestionChoice/Tag**.

- **Module:** `github.com/namchokGithub/vocabunny-core-api`
- **Go:** 1.24.1
- **Framework:** Echo v4
- **ORM:** GORM + PostgreSQL
- **Auth:** JWT (access + refresh tokens)

---

## Architecture

Hexagonal (ports & adapters). Three concentric layers — never import inward past your layer.

```
cmd/                         ← entrypoint
configs/                     ← env config
infrastructure/              ← DB, Redis, JWT, storage init
protocol/                    ← HTTP server + route registration
internal/
  core/
    domain/                  ← pure structs, no imports from other layers
    port/                    ← interface definitions (repository.go, service.go)
    helper/                  ← shared utilities (pagination, errors, context, validator)
    service/content/         ← business logic implements port.*Service
    service/identity/
  repository/content/        ← GORM implements port.*Repository
                               model.go  — GORM models
                               repository.go — all content repos in one file
  handler/content/           ← Echo handlers, DTOs, transformers
  handler/identity/
  middleware/
```

### Key contracts

| Layer | Talks to |
|---|---|
| handler | port.*Service (interface) |
| service | port.*Repository (interface) |
| repository | GORM + domain types |

---

## Active Refactor: `?include=` Query Parameter Support

Branch: `feat/game-play`. Follow **progress.md** for the full task list. Summary below.

### What is done

| File | Status |
|---|---|
| `internal/core/domain/common.go` | `Includes` type, `ParseIncludes`, `.Has()` |
| `internal/core/domain/content.go` | Pointer relation fields on `Lesson`, `Unit`, `QuestionSet`, `Question`; `Includes Includes` added to all four Query structs |
| `internal/core/port/repository.go` | `FindByID` signatures updated for Lesson, Unit, QuestionSet, Question |
| `internal/core/port/service.go` | `FindByID` signatures updated for Lesson, Unit, QuestionSet, Question |
| `internal/repository/content/model.go` | GORM association fields added to `LessonModel`, `UnitModel`, `QuestionSetModel`, `QuestionModel` |
| `internal/repository/content/repository.go` | All `FindByID` / `FindAll` methods updated with selective `Preload`; internal Create/Update calls pass `nil` |

### What still needs doing

Work through steps A → F in order; each step must compile before starting the next.

---

#### Step A — Fix domain mappers in `repository.go`

File: `internal/repository/content/repository.go`

The four `toDomain*` functions (lines ~1229–1288) map GORM models to domain structs but do not yet populate the new pointer relation fields. Update each one to check the zero-UUID sentinel:

```go
func toDomainLesson(model LessonModel) domain.Lesson {
    lesson := domain.Lesson{
        ID: model.ID, SectionID: model.SectionID, /* ... */
    }
    if model.Section.ID != uuid.Nil {
        s := toDomainSection(model.Section)
        lesson.Section = &s
    }
    return lesson
}

func toDomainUnit(model UnitModel) domain.Unit {
    unit := domain.Unit{ /* scalar fields */ }
    if model.Lesson.ID != uuid.Nil {
        l := toDomainLesson(model.Lesson)
        unit.Lesson = &l
    }
    return unit
}

func toDomainQuestionSet(model QuestionSetModel) domain.QuestionSet {
    qs := domain.QuestionSet{ /* scalar fields */ }
    if model.Unit.ID != uuid.Nil {
        u := domain.Unit{
            ID: model.Unit.ID, LessonID: model.Unit.LessonID,
            Slug: model.Unit.Slug, Title: model.Unit.Title,
            /* other scalar fields — no nested Lesson */
        }
        qs.Unit = &u
    }
    if model.Unit.Lesson.ID != uuid.Nil {
        l := domain.Lesson{
            ID: model.Unit.Lesson.ID, SectionID: model.Unit.Lesson.SectionID,
            Slug: model.Unit.Lesson.Slug, Title: model.Unit.Lesson.Title,
            /* scalar only */
        }
        qs.Lesson = &l
    }
    return qs
}

func toDomainQuestion(model QuestionModel) domain.Question {
    q := domain.Question{ /* scalar fields */ }
    if model.QuestionSet.ID != uuid.Nil {
        qs := toDomainQuestionSet(model.QuestionSet)
        q.QuestionSet = &qs
    }
    return q
}
```

---

#### Step B — Service files (4 files + 1 auxiliary)

Fix all `FindByID` call sites that no longer match the new signature.

**Rules:**
- Internal/validation/write calls (e.g., existence checks inside `RunInTx`) → pass `nil`
- Public `FindByID(ctx, id)` method on the service → accept `includes domain.Includes` and pass through

Files and approximate lines:

| File | Lines to fix |
|---|---|
| `internal/core/service/content/lesson_service.go` | 70 (Update tx), 115 (Delete tx); also update `FindByID` method signature at line 122 |
| `internal/core/service/content/unit_service.go` | symmetric to lesson |
| `internal/core/service/content/question_set_service.go` | symmetric to lesson |
| `internal/core/service/content/question_service.go` | multiple sites — internal ones get `nil`, public method gets `includes` |
| `internal/core/service/content/question_choice_service.go` | line ~44 — calls `s.questionRepository.FindByID`; pass `nil` |

Pattern for the public method (example for Lesson):

```go
func (s *lessonService) FindByID(ctx context.Context, id uuid.UUID, includes domain.Includes) (domain.Lesson, error) {
    return s.lessonRepository.FindByID(ctx, id, includes)
}
```

---

#### Step C — DTOs (`internal/handler/content/`)

Add summary DTOs (flat, no nesting) to the relevant `*_dto.go` files, and add the optional relation field to each `*Response`:

```go
// section_dto.go
type SectionSummaryDTO struct {
    ID        string `json:"id"`
    Slug      string `json:"slug"`
    Title     string `json:"title"`
}

// lesson_dto.go
type LessonSummaryDTO struct {
    ID        string `json:"id"`
    SectionID string `json:"section_id"`
    Slug      string `json:"slug"`
    Title     string `json:"title"`
}
// Add to LessonResponse:
Section *SectionSummaryDTO `json:"section,omitempty"`

// unit_dto.go
type UnitSummaryDTO struct {
    ID       string `json:"id"`
    LessonID string `json:"lesson_id"`
    Slug     string `json:"slug"`
    Title    string `json:"title"`
}
// Add to UnitResponse:
Lesson *LessonSummaryDTO `json:"lesson,omitempty"`

// question_set_dto.go
type QuestionSetSummaryDTO struct {
    ID      string `json:"id"`
    UnitID  string `json:"unit_id"`
    Slug    string `json:"slug"`
    Title   string `json:"title"`
    Version int    `json:"version"`
}
// Add to QuestionSetResponse:
Unit    *UnitSummaryDTO    `json:"unit,omitempty"`
Lesson  *LessonSummaryDTO  `json:"lesson,omitempty"`

// question_dto.go
// Add to QuestionResponse:
QuestionSet *QuestionSetSummaryDTO `json:"question_set,omitempty"`
```

---

#### Step D — Transformers (`internal/handler/content/`)

Update `toLessonResponse`, `toUnitResponse`, `toQuestionSetResponse`, `toQuestionResponse` in the corresponding `*_transformer.go` files to map non-nil domain relation pointers to summary DTOs:

```go
// example in lesson_transformer.go
func toLessonResponse(item domain.Lesson) LessonResponse {
    resp := LessonResponse{ /* existing fields */ }
    if item.Section != nil {
        resp.Section = &SectionSummaryDTO{
            ID:    item.Section.ID.String(),
            Slug:  item.Section.Slug,
            Title: item.Section.Title,
        }
    }
    return resp
}
```

---

#### Step E — Handlers

Parse `?include=` in handlers for Lesson, Unit, QuestionSet, Question.

**FindByID:**
```go
includes := domain.ParseIncludes(c.QueryParam("include"))
item, err := h.service.FindByID(c.Request().Context(), id, includes)
```

**FindAll:**
```go
query.Includes = domain.ParseIncludes(c.QueryParam("include"))
```

For the Question `FindAll` handler: **remove** the old `include_choices` and `include_tags` boolean query params — they are now replaced by `?include=choices,tags`.

---

#### Step F — API.md

Add an **Include Query Parameter** section documenting:
- Supported includes per entity (see table below)
- Syntax: `?include=unit,lesson`
- Performance note: selective GORM preload, summary-only on list responses, no deep nesting
- Breaking change: question list no longer includes choices/tags by default; use `?include=choices,tags`

| Entity | Supported includes |
|---|---|
| Lesson | `section` |
| Unit | `lesson` |
| QuestionSet | `unit`, `lesson` |
| Question | `question_set`, `choices`, `tags` |

---

## Key Architecture Decisions

| Decision | Rationale |
|---|---|
| `Includes` lives in `domain` | Avoids import cycles; domain is the shared nucleus |
| `FindByID` signature extended (not a new method) | Clean interface; internal callers just pass `nil` |
| Internal service calls pass `nil` includes | Validation/write paths don't need relations; avoids wasteful joins |
| `Question.FindByID` always loads choices+tags | Single-item read was always full; backward compat preserved |
| `Question.FindAll` choices/tags are now opt-in | Breaking but intentional — consistent with `?include=` pattern |
| `include=lesson` on QuestionSet uses `Preload("Unit.Lesson")` | GORM resolves via WHERE IN — 3 queries total for a list page, not N+1 |
| Zero-UUID check to detect loaded associations | `if model.Section.ID != uuid.Nil` — no extra flags needed |
| Summary DTOs are flat | Prevents recursive/circular DTO structures |
| Preload added **after** COUNT query | COUNT must not have joins that inflate row count |

---

## Common Patterns

### Error helpers (`internal/core/helper/errors.go`)
```go
helper.BadRequest("code", "message", err)
helper.NotFound("code", "message", err)
helper.Conflict("code", "message", err)
helper.Internal("code", "message", err)
```

### Respond helpers
```go
helper.RespondSuccess(c, http.StatusOK, payload)
helper.RespondError(c, err)
```

### EntityField (partial update)
```go
domain.NewEntityField(value)  // marks a field as "explicitly set"
if input.Slug.Set { model.Slug = input.Slug.Value }
```

### Transaction
```go
s.txManager.RunInTx(ctx, func(txCtx context.Context) error { ... })
```
Pass `txCtx` to all repository calls inside the closure so they share the transaction.

---

## Build & Test

```bash
go build ./...          # verify no compile errors
go test ./...           # run all tests
```

Ensure `go build ./...` passes with zero errors before marking any step complete.
