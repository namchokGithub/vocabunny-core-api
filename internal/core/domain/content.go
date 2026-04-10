package domain

import "github.com/google/uuid"

type Section struct {
	ID          uuid.UUID
	Slug        string
	Title       string
	Description *string
	OrderNo     int
	IsPublished bool
	AuditFields
}

type SectionUpdateInput struct {
	ID          uuid.UUID
	Slug        EntityField[string]
	Title       EntityField[string]
	Description EntityField[*string]
	OrderNo     EntityField[int]
	IsPublished EntityField[bool]
	ActorID     string
}

type SectionQuery struct {
	Paging      Paging
	Search      string
	IsPublished *bool
	SortBy      string
	SortOrder   string
}

type Lesson struct {
	ID          uuid.UUID
	SectionID   uuid.UUID
	Slug        string
	Title       string
	Description *string
	OrderNo     int
	IsPublished bool
	AuditFields
}

type LessonUpdateInput struct {
	ID          uuid.UUID
	SectionID   EntityField[uuid.UUID]
	Slug        EntityField[string]
	Title       EntityField[string]
	Description EntityField[*string]
	OrderNo     EntityField[int]
	IsPublished EntityField[bool]
	ActorID     string
}

type LessonQuery struct {
	Paging      Paging
	SectionID   *uuid.UUID
	Search      string
	IsPublished *bool
	SortBy      string
	SortOrder   string
}

type Unit struct {
	ID          uuid.UUID
	LessonID    uuid.UUID
	Slug        string
	Title       string
	Description *string
	OrderNo     int
	IsPublished bool
	AuditFields
}

type UnitUpdateInput struct {
	ID          uuid.UUID
	LessonID    EntityField[uuid.UUID]
	Slug        EntityField[string]
	Title       EntityField[string]
	Description EntityField[*string]
	OrderNo     EntityField[int]
	IsPublished EntityField[bool]
	ActorID     string
}

type UnitQuery struct {
	Paging      Paging
	LessonID    *uuid.UUID
	Search      string
	IsPublished *bool
	SortBy      string
	SortOrder   string
}

type QuestionSet struct {
	ID               uuid.UUID
	UnitID           uuid.UUID
	Slug             string
	Title            string
	Description      *string
	OrderNo          int
	EstimatedSeconds *int
	IsPublished      bool
	Version          int
	AuditFields
}

type QuestionSetUpdateInput struct {
	ID               uuid.UUID
	UnitID           EntityField[uuid.UUID]
	Slug             EntityField[string]
	Title            EntityField[string]
	Description      EntityField[*string]
	OrderNo          EntityField[int]
	EstimatedSeconds EntityField[*int]
	IsPublished      EntityField[bool]
	Version          EntityField[int]
	ActorID          string
}

type QuestionSetQuery struct {
	Paging      Paging
	UnitID      *uuid.UUID
	Search      string
	IsPublished *bool
	Version     *int
	SortBy      string
	SortOrder   string
}

type Question struct {
	ID            uuid.UUID
	QuestionSetID uuid.UUID
	Type          string
	QuestionText  string
	BlankPosition *int
	Explanation   *string
	ImageURL      *string
	Difficulty    int
	OrderNo       int
	IsActive      bool
	Choices       []QuestionChoice
	Tags          []Tag
	AuditFields
}

type QuestionUpdateInput struct {
	ID            uuid.UUID
	QuestionSetID EntityField[uuid.UUID]
	Type          EntityField[string]
	QuestionText  EntityField[string]
	BlankPosition EntityField[*int]
	Explanation   EntityField[*string]
	ImageURL      EntityField[*string]
	Difficulty    EntityField[int]
	OrderNo       EntityField[int]
	IsActive      EntityField[bool]
	ActorID       string
}

type QuestionQuery struct {
	Paging         Paging
	QuestionSetID  *uuid.UUID
	Type           *string
	IsActive       *bool
	Search         string
	SortBy         string
	SortOrder      string
	IncludeChoices bool
	IncludeTags    bool
}

type QuestionChoice struct {
	ID          uuid.UUID
	QuestionID  uuid.UUID
	ChoiceText  string
	ChoiceOrder int
	IsCorrect   bool
	AuditFields
}

type QuestionChoiceInput struct {
	ID          uuid.UUID
	ChoiceText  string
	ChoiceOrder int
	IsCorrect   bool
}

type QuestionChoiceUpdateInput struct {
	ID          uuid.UUID
	ChoiceText  EntityField[string]
	ChoiceOrder EntityField[int]
	IsCorrect   EntityField[bool]
	ActorID     string
}

type QuestionChoiceQuery struct {
	Paging     Paging
	QuestionID *uuid.UUID
	SortBy     string
	SortOrder  string
}

type Tag struct {
	ID   uuid.UUID
	Name string
	AuditFields
}

type TagUpdateInput struct {
	ID      uuid.UUID
	Name    EntityField[string]
	ActorID string
}

type TagQuery struct {
	Paging    Paging
	Search    string
	SortBy    string
	SortOrder string
}
