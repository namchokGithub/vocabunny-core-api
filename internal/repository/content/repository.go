package content

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/domain"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/helper"
	"github.com/namchokGithub/vocabunny-core-api/internal/core/port"
	"gorm.io/gorm"
)

type Dependencies struct {
	DB *gorm.DB
}

type Repository struct {
	Section        port.SectionRepository
	Lesson         port.LessonRepository
	Unit           port.UnitRepository
	QuestionSet    port.QuestionSetRepository
	Question       port.QuestionRepository
	QuestionChoice port.QuestionChoiceRepository
	Tag            port.TagRepository
}

func NewRepository(deps Dependencies) *Repository {
	base := &baseRepository{db: deps.DB}
	return &Repository{
		Section:        &sectionRepository{baseRepository: base},
		Lesson:         &lessonRepository{baseRepository: base},
		Unit:           &unitRepository{baseRepository: base},
		QuestionSet:    &questionSetRepository{baseRepository: base},
		Question:       &questionRepository{baseRepository: base},
		QuestionChoice: &questionChoiceRepository{baseRepository: base},
		Tag:            &tagRepository{baseRepository: base},
	}
}

type baseRepository struct {
	db *gorm.DB
}

// dbWithContext keeps repository methods transaction-safe.
// When a transaction is attached to the context, every query in the same request
// uses that transaction automatically.
func (r *baseRepository) dbWithContext(ctx context.Context) *gorm.DB {
	tx, ok := helper.TxFromContext(ctx).(*gorm.DB)
	if ok && tx != nil {
		return tx.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}

type sectionRepository struct {
	*baseRepository
}

func (r *sectionRepository) Create(ctx context.Context, section domain.Section) (domain.Section, error) {
	now := time.Now()
	model := SectionModel{
		ID:          section.ID,
		Slug:        section.Slug,
		Title:       section.Title,
		Description: section.Description,
		OrderNo:     section.OrderNo,
		IsPublished: section.IsPublished,
		CreatedBy:   section.CreatedBy,
		UpdatedBy:   section.UpdatedBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if model.ID == uuid.Nil {
		model.ID = uuid.New()
	}
	if err := r.dbWithContext(ctx).Create(&model).Error; err != nil {
		return domain.Section{}, helper.Internal("create_section_failed", "failed to create section", err)
	}
	return r.FindByID(ctx, model.ID)
}

func (r *sectionRepository) Update(ctx context.Context, input domain.SectionUpdateInput) (domain.Section, error) {
	var model SectionModel
	if err := r.dbWithContext(ctx).First(&model, "id = ?", input.ID).Error; err != nil {
		return domain.Section{}, mapGormNotFound(err, "section_not_found", "section not found", "find_section_failed", "failed to load section")
	}
	if input.Slug.Set {
		model.Slug = input.Slug.Value
	}
	if input.Title.Set {
		model.Title = input.Title.Value
	}
	if input.Description.Set {
		model.Description = input.Description.Value
	}
	if input.OrderNo.Set {
		model.OrderNo = input.OrderNo.Value
	}
	if input.IsPublished.Set {
		model.IsPublished = input.IsPublished.Value
	}
	model.UpdatedBy = input.ActorID
	model.UpdatedAt = time.Now()

	if err := r.dbWithContext(ctx).Save(&model).Error; err != nil {
		return domain.Section{}, helper.Internal("update_section_failed", "failed to update section", err)
	}
	return r.FindByID(ctx, model.ID)
}

func (r *sectionRepository) Delete(ctx context.Context, id uuid.UUID, actorID string) error {
	if err := r.dbWithContext(ctx).Model(&SectionModel{}).Where("id = ?", id).Updates(map[string]any{
		"updated_by": actorID,
		"updated_at": time.Now(),
	}).Delete(&SectionModel{}, "id = ?", id).Error; err != nil {
		return helper.Internal("delete_section_failed", "failed to delete section", err)
	}
	return nil
}

func (r *sectionRepository) FindByID(ctx context.Context, id uuid.UUID) (domain.Section, error) {
	var model SectionModel
	if err := r.dbWithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return domain.Section{}, mapGormNotFound(err, "section_not_found", "section not found", "find_section_failed", "failed to find section")
	}
	return toDomainSection(model), nil
}

func (r *sectionRepository) FindAll(ctx context.Context, query domain.SectionQuery) (domain.PageResult[domain.Section], error) {
	db := r.dbWithContext(ctx).Model(&SectionModel{})
	if query.Search != "" {
		search := "%" + strings.ToLower(query.Search) + "%"
		db = db.Where("LOWER(slug) LIKE ? OR LOWER(title) LIKE ? OR LOWER(description) LIKE ?", search, search, search)
	}
	if query.IsPublished != nil {
		db = db.Where("is_published = ?", *query.IsPublished)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return domain.PageResult[domain.Section]{}, helper.Internal("count_sections_failed", "failed to count sections", err)
	}

	sortBy := safeSort(query.SortBy, []string{"slug", "title", "order_no", "is_published", "created_at", "updated_at"}, "order_no")
	sortOrder := safeSortOrder(query.SortOrder)

	var models []SectionModel
	if err := db.Order("tbl_sections." + sortBy + " " + sortOrder).
		Order("tbl_sections.id ASC").
		Offset(query.Paging.Offset()).
		Limit(query.Paging.Limit).
		Find(&models).Error; err != nil {
		return domain.PageResult[domain.Section]{}, helper.Internal("list_sections_failed", "failed to list sections", err)
	}

	items := make([]domain.Section, 0, len(models))
	for _, model := range models {
		items = append(items, toDomainSection(model))
	}

	return domain.PageResult[domain.Section]{
		Items: items,
		Paging: domain.Paging{
			Page:  query.Paging.Page,
			Limit: query.Paging.Limit,
			Total: total,
		},
	}, nil
}

func (r *sectionRepository) ExistsBySlug(ctx context.Context, slug string, excludeID *uuid.UUID) (bool, error) {
	return existsByStringField(ctx, r.dbWithContext, &SectionModel{}, "slug", slug, excludeID, "exists_section_failed", "failed to check section uniqueness")
}

type lessonRepository struct {
	*baseRepository
}

func (r *lessonRepository) Create(ctx context.Context, lesson domain.Lesson) (domain.Lesson, error) {
	now := time.Now()
	model := LessonModel{
		ID:          lesson.ID,
		SectionID:   lesson.SectionID,
		Slug:        lesson.Slug,
		Title:       lesson.Title,
		Description: lesson.Description,
		OrderNo:     lesson.OrderNo,
		IsPublished: lesson.IsPublished,
		CreatedBy:   lesson.CreatedBy,
		UpdatedBy:   lesson.UpdatedBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if model.ID == uuid.Nil {
		model.ID = uuid.New()
	}
	if err := r.dbWithContext(ctx).Create(&model).Error; err != nil {
		return domain.Lesson{}, helper.Internal("create_lesson_failed", "failed to create lesson", err)
	}
	return r.FindByID(ctx, model.ID)
}

func (r *lessonRepository) Update(ctx context.Context, input domain.LessonUpdateInput) (domain.Lesson, error) {
	var model LessonModel
	if err := r.dbWithContext(ctx).First(&model, "id = ?", input.ID).Error; err != nil {
		return domain.Lesson{}, mapGormNotFound(err, "lesson_not_found", "lesson not found", "find_lesson_failed", "failed to load lesson")
	}
	if input.SectionID.Set {
		model.SectionID = input.SectionID.Value
	}
	if input.Slug.Set {
		model.Slug = input.Slug.Value
	}
	if input.Title.Set {
		model.Title = input.Title.Value
	}
	if input.Description.Set {
		model.Description = input.Description.Value
	}
	if input.OrderNo.Set {
		model.OrderNo = input.OrderNo.Value
	}
	if input.IsPublished.Set {
		model.IsPublished = input.IsPublished.Value
	}
	model.UpdatedBy = input.ActorID
	model.UpdatedAt = time.Now()

	if err := r.dbWithContext(ctx).Save(&model).Error; err != nil {
		return domain.Lesson{}, helper.Internal("update_lesson_failed", "failed to update lesson", err)
	}
	return r.FindByID(ctx, model.ID)
}

func (r *lessonRepository) Delete(ctx context.Context, id uuid.UUID, actorID string) error {
	if err := r.dbWithContext(ctx).Model(&LessonModel{}).Where("id = ?", id).Updates(map[string]any{
		"updated_by": actorID,
		"updated_at": time.Now(),
	}).Delete(&LessonModel{}, "id = ?", id).Error; err != nil {
		return helper.Internal("delete_lesson_failed", "failed to delete lesson", err)
	}
	return nil
}

func (r *lessonRepository) FindByID(ctx context.Context, id uuid.UUID) (domain.Lesson, error) {
	var model LessonModel
	if err := r.dbWithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return domain.Lesson{}, mapGormNotFound(err, "lesson_not_found", "lesson not found", "find_lesson_failed", "failed to find lesson")
	}
	return toDomainLesson(model), nil
}

func (r *lessonRepository) FindAll(ctx context.Context, query domain.LessonQuery) (domain.PageResult[domain.Lesson], error) {
	db := r.dbWithContext(ctx).Model(&LessonModel{})
	if query.SectionID != nil {
		db = db.Where("section_id = ?", *query.SectionID)
	}
	if query.Search != "" {
		search := "%" + strings.ToLower(query.Search) + "%"
		db = db.Where("LOWER(slug) LIKE ? OR LOWER(title) LIKE ? OR LOWER(description) LIKE ?", search, search, search)
	}
	if query.IsPublished != nil {
		db = db.Where("is_published = ?", *query.IsPublished)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return domain.PageResult[domain.Lesson]{}, helper.Internal("count_lessons_failed", "failed to count lessons", err)
	}

	sortBy := safeSort(query.SortBy, []string{"slug", "title", "order_no", "is_published", "created_at", "updated_at"}, "order_no")
	sortOrder := safeSortOrder(query.SortOrder)

	var models []LessonModel
	if err := db.Order("tbl_lessons." + sortBy + " " + sortOrder).
		Order("tbl_lessons.id ASC").
		Offset(query.Paging.Offset()).
		Limit(query.Paging.Limit).
		Find(&models).Error; err != nil {
		return domain.PageResult[domain.Lesson]{}, helper.Internal("list_lessons_failed", "failed to list lessons", err)
	}

	items := make([]domain.Lesson, 0, len(models))
	for _, model := range models {
		items = append(items, toDomainLesson(model))
	}

	return domain.PageResult[domain.Lesson]{
		Items: items,
		Paging: domain.Paging{
			Page:  query.Paging.Page,
			Limit: query.Paging.Limit,
			Total: total,
		},
	}, nil
}

func (r *lessonRepository) ExistsBySlug(ctx context.Context, sectionID uuid.UUID, slug string, excludeID *uuid.UUID) (bool, error) {
	db := r.dbWithContext(ctx).Model(&LessonModel{}).Where("section_id = ? AND LOWER(slug) = ?", sectionID, strings.ToLower(strings.TrimSpace(slug)))
	if excludeID != nil {
		db = db.Where("id <> ?", *excludeID)
	}
	var count int64
	if err := db.Count(&count).Error; err != nil {
		return false, helper.Internal("exists_lesson_failed", "failed to check lesson uniqueness", err)
	}
	return count > 0, nil
}

type unitRepository struct {
	*baseRepository
}

func (r *unitRepository) Create(ctx context.Context, unit domain.Unit) (domain.Unit, error) {
	now := time.Now()
	model := UnitModel{
		ID:          unit.ID,
		LessonID:    unit.LessonID,
		Slug:        unit.Slug,
		Title:       unit.Title,
		Description: unit.Description,
		OrderNo:     unit.OrderNo,
		IsPublished: unit.IsPublished,
		CreatedBy:   unit.CreatedBy,
		UpdatedBy:   unit.UpdatedBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if model.ID == uuid.Nil {
		model.ID = uuid.New()
	}
	if err := r.dbWithContext(ctx).Create(&model).Error; err != nil {
		return domain.Unit{}, helper.Internal("create_unit_failed", "failed to create unit", err)
	}
	return r.FindByID(ctx, model.ID)
}

func (r *unitRepository) Update(ctx context.Context, input domain.UnitUpdateInput) (domain.Unit, error) {
	var model UnitModel
	if err := r.dbWithContext(ctx).First(&model, "id = ?", input.ID).Error; err != nil {
		return domain.Unit{}, mapGormNotFound(err, "unit_not_found", "unit not found", "find_unit_failed", "failed to load unit")
	}
	if input.LessonID.Set {
		model.LessonID = input.LessonID.Value
	}
	if input.Slug.Set {
		model.Slug = input.Slug.Value
	}
	if input.Title.Set {
		model.Title = input.Title.Value
	}
	if input.Description.Set {
		model.Description = input.Description.Value
	}
	if input.OrderNo.Set {
		model.OrderNo = input.OrderNo.Value
	}
	if input.IsPublished.Set {
		model.IsPublished = input.IsPublished.Value
	}
	model.UpdatedBy = input.ActorID
	model.UpdatedAt = time.Now()

	if err := r.dbWithContext(ctx).Save(&model).Error; err != nil {
		return domain.Unit{}, helper.Internal("update_unit_failed", "failed to update unit", err)
	}
	return r.FindByID(ctx, model.ID)
}

func (r *unitRepository) Delete(ctx context.Context, id uuid.UUID, actorID string) error {
	if err := r.dbWithContext(ctx).Model(&UnitModel{}).Where("id = ?", id).Updates(map[string]any{
		"updated_by": actorID,
		"updated_at": time.Now(),
	}).Delete(&UnitModel{}, "id = ?", id).Error; err != nil {
		return helper.Internal("delete_unit_failed", "failed to delete unit", err)
	}
	return nil
}

func (r *unitRepository) FindByID(ctx context.Context, id uuid.UUID) (domain.Unit, error) {
	var model UnitModel
	if err := r.dbWithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return domain.Unit{}, mapGormNotFound(err, "unit_not_found", "unit not found", "find_unit_failed", "failed to find unit")
	}
	return toDomainUnit(model), nil
}

func (r *unitRepository) FindAll(ctx context.Context, query domain.UnitQuery) (domain.PageResult[domain.Unit], error) {
	db := r.dbWithContext(ctx).Model(&UnitModel{})
	if query.LessonID != nil {
		db = db.Where("lesson_id = ?", *query.LessonID)
	}
	if query.Search != "" {
		search := "%" + strings.ToLower(query.Search) + "%"
		db = db.Where("LOWER(slug) LIKE ? OR LOWER(title) LIKE ? OR LOWER(description) LIKE ?", search, search, search)
	}
	if query.IsPublished != nil {
		db = db.Where("is_published = ?", *query.IsPublished)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return domain.PageResult[domain.Unit]{}, helper.Internal("count_units_failed", "failed to count units", err)
	}

	sortBy := safeSort(query.SortBy, []string{"slug", "title", "order_no", "is_published", "created_at", "updated_at"}, "order_no")
	sortOrder := safeSortOrder(query.SortOrder)

	var models []UnitModel
	if err := db.Order("tbl_units." + sortBy + " " + sortOrder).
		Order("tbl_units.id ASC").
		Offset(query.Paging.Offset()).
		Limit(query.Paging.Limit).
		Find(&models).Error; err != nil {
		return domain.PageResult[domain.Unit]{}, helper.Internal("list_units_failed", "failed to list units", err)
	}

	items := make([]domain.Unit, 0, len(models))
	for _, model := range models {
		items = append(items, toDomainUnit(model))
	}

	return domain.PageResult[domain.Unit]{
		Items: items,
		Paging: domain.Paging{
			Page:  query.Paging.Page,
			Limit: query.Paging.Limit,
			Total: total,
		},
	}, nil
}

func (r *unitRepository) ExistsBySlug(ctx context.Context, lessonID uuid.UUID, slug string, excludeID *uuid.UUID) (bool, error) {
	db := r.dbWithContext(ctx).Model(&UnitModel{}).Where("lesson_id = ? AND LOWER(slug) = ?", lessonID, strings.ToLower(strings.TrimSpace(slug)))
	if excludeID != nil {
		db = db.Where("id <> ?", *excludeID)
	}
	var count int64
	if err := db.Count(&count).Error; err != nil {
		return false, helper.Internal("exists_unit_failed", "failed to check unit uniqueness", err)
	}
	return count > 0, nil
}

type questionSetRepository struct {
	*baseRepository
}

func (r *questionSetRepository) Create(ctx context.Context, questionSet domain.QuestionSet) (domain.QuestionSet, error) {
	now := time.Now()
	model := QuestionSetModel{
		ID:               questionSet.ID,
		UnitID:           questionSet.UnitID,
		Slug:             questionSet.Slug,
		Title:            questionSet.Title,
		Description:      questionSet.Description,
		OrderNo:          questionSet.OrderNo,
		EstimatedSeconds: questionSet.EstimatedSeconds,
		IsPublished:      questionSet.IsPublished,
		Version:          questionSet.Version,
		CreatedBy:        questionSet.CreatedBy,
		UpdatedBy:        questionSet.UpdatedBy,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	if model.ID == uuid.Nil {
		model.ID = uuid.New()
	}
	if model.Version == 0 {
		model.Version = 1
	}
	if err := r.dbWithContext(ctx).Create(&model).Error; err != nil {
		return domain.QuestionSet{}, helper.Internal("create_question_set_failed", "failed to create question set", err)
	}
	return r.FindByID(ctx, model.ID)
}

func (r *questionSetRepository) Update(ctx context.Context, input domain.QuestionSetUpdateInput) (domain.QuestionSet, error) {
	var model QuestionSetModel
	if err := r.dbWithContext(ctx).First(&model, "id = ?", input.ID).Error; err != nil {
		return domain.QuestionSet{}, mapGormNotFound(err, "question_set_not_found", "question set not found", "find_question_set_failed", "failed to load question set")
	}
	if input.UnitID.Set {
		model.UnitID = input.UnitID.Value
	}
	if input.Slug.Set {
		model.Slug = input.Slug.Value
	}
	if input.Title.Set {
		model.Title = input.Title.Value
	}
	if input.Description.Set {
		model.Description = input.Description.Value
	}
	if input.OrderNo.Set {
		model.OrderNo = input.OrderNo.Value
	}
	if input.EstimatedSeconds.Set {
		model.EstimatedSeconds = input.EstimatedSeconds.Value
	}
	if input.IsPublished.Set {
		model.IsPublished = input.IsPublished.Value
	}
	if input.Version.Set {
		model.Version = input.Version.Value
	}
	model.UpdatedBy = input.ActorID
	model.UpdatedAt = time.Now()

	if err := r.dbWithContext(ctx).Save(&model).Error; err != nil {
		return domain.QuestionSet{}, helper.Internal("update_question_set_failed", "failed to update question set", err)
	}
	return r.FindByID(ctx, model.ID)
}

func (r *questionSetRepository) Delete(ctx context.Context, id uuid.UUID, actorID string) error {
	if err := r.dbWithContext(ctx).Model(&QuestionSetModel{}).Where("id = ?", id).Updates(map[string]any{
		"updated_by": actorID,
		"updated_at": time.Now(),
	}).Delete(&QuestionSetModel{}, "id = ?", id).Error; err != nil {
		return helper.Internal("delete_question_set_failed", "failed to delete question set", err)
	}
	return nil
}

func (r *questionSetRepository) FindByID(ctx context.Context, id uuid.UUID) (domain.QuestionSet, error) {
	var model QuestionSetModel
	if err := r.dbWithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return domain.QuestionSet{}, mapGormNotFound(err, "question_set_not_found", "question set not found", "find_question_set_failed", "failed to find question set")
	}
	return toDomainQuestionSet(model), nil
}

func (r *questionSetRepository) FindAll(ctx context.Context, query domain.QuestionSetQuery) (domain.PageResult[domain.QuestionSet], error) {
	db := r.dbWithContext(ctx).Model(&QuestionSetModel{})
	if query.UnitID != nil {
		db = db.Where("unit_id = ?", *query.UnitID)
	}
	if query.Search != "" {
		search := "%" + strings.ToLower(query.Search) + "%"
		db = db.Where("LOWER(slug) LIKE ? OR LOWER(title) LIKE ? OR LOWER(description) LIKE ?", search, search, search)
	}
	if query.IsPublished != nil {
		db = db.Where("is_published = ?", *query.IsPublished)
	}
	if query.Version != nil {
		db = db.Where("version = ?", *query.Version)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return domain.PageResult[domain.QuestionSet]{}, helper.Internal("count_question_sets_failed", "failed to count question sets", err)
	}

	sortBy := safeSort(query.SortBy, []string{"slug", "title", "order_no", "version", "is_published", "created_at", "updated_at"}, "order_no")
	sortOrder := safeSortOrder(query.SortOrder)

	var models []QuestionSetModel
	if err := db.Order("tbl_question_sets." + sortBy + " " + sortOrder).
		Order("tbl_question_sets.id ASC").
		Offset(query.Paging.Offset()).
		Limit(query.Paging.Limit).
		Find(&models).Error; err != nil {
		return domain.PageResult[domain.QuestionSet]{}, helper.Internal("list_question_sets_failed", "failed to list question sets", err)
	}

	items := make([]domain.QuestionSet, 0, len(models))
	for _, model := range models {
		items = append(items, toDomainQuestionSet(model))
	}

	return domain.PageResult[domain.QuestionSet]{
		Items: items,
		Paging: domain.Paging{
			Page:  query.Paging.Page,
			Limit: query.Paging.Limit,
			Total: total,
		},
	}, nil
}

func (r *questionSetRepository) ExistsBySlugVersion(ctx context.Context, unitID uuid.UUID, slug string, version int, excludeID *uuid.UUID) (bool, error) {
	db := r.dbWithContext(ctx).Model(&QuestionSetModel{}).
		Where("unit_id = ? AND LOWER(slug) = ? AND version = ?", unitID, strings.ToLower(strings.TrimSpace(slug)), version)
	if excludeID != nil {
		db = db.Where("id <> ?", *excludeID)
	}
	var count int64
	if err := db.Count(&count).Error; err != nil {
		return false, helper.Internal("exists_question_set_failed", "failed to check question set uniqueness", err)
	}
	return count > 0, nil
}

type questionRepository struct {
	*baseRepository
}

func (r *questionRepository) Create(ctx context.Context, question domain.Question) (domain.Question, error) {
	now := time.Now()
	model := QuestionModel{
		ID:            question.ID,
		QuestionSetID: question.QuestionSetID,
		Type:          question.Type,
		QuestionText:  question.QuestionText,
		BlankPosition: question.BlankPosition,
		Explanation:   question.Explanation,
		ImageURL:      question.ImageURL,
		Difficulty:    question.Difficulty,
		OrderNo:       question.OrderNo,
		IsActive:      question.IsActive,
		CreatedBy:     question.CreatedBy,
		UpdatedBy:     question.UpdatedBy,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if model.ID == uuid.Nil {
		model.ID = uuid.New()
	}
	if model.Difficulty == 0 {
		model.Difficulty = 1
	}
	if err := r.dbWithContext(ctx).Create(&model).Error; err != nil {
		return domain.Question{}, helper.Internal("create_question_failed", "failed to create question", err)
	}
	return r.FindByID(ctx, model.ID)
}

func (r *questionRepository) Update(ctx context.Context, input domain.QuestionUpdateInput) (domain.Question, error) {
	var model QuestionModel
	if err := r.dbWithContext(ctx).First(&model, "id = ?", input.ID).Error; err != nil {
		return domain.Question{}, mapGormNotFound(err, "question_not_found", "question not found", "find_question_failed", "failed to load question")
	}
	if input.QuestionSetID.Set {
		model.QuestionSetID = input.QuestionSetID.Value
	}
	if input.Type.Set {
		model.Type = input.Type.Value
	}
	if input.QuestionText.Set {
		model.QuestionText = input.QuestionText.Value
	}
	if input.BlankPosition.Set {
		model.BlankPosition = input.BlankPosition.Value
	}
	if input.Explanation.Set {
		model.Explanation = input.Explanation.Value
	}
	if input.ImageURL.Set {
		model.ImageURL = input.ImageURL.Value
	}
	if input.Difficulty.Set {
		model.Difficulty = input.Difficulty.Value
	}
	if input.OrderNo.Set {
		model.OrderNo = input.OrderNo.Value
	}
	if input.IsActive.Set {
		model.IsActive = input.IsActive.Value
	}
	model.UpdatedBy = input.ActorID
	model.UpdatedAt = time.Now()

	if err := r.dbWithContext(ctx).Save(&model).Error; err != nil {
		return domain.Question{}, helper.Internal("update_question_failed", "failed to update question", err)
	}
	return r.FindByID(ctx, model.ID)
}

func (r *questionRepository) Delete(ctx context.Context, id uuid.UUID, actorID string) error {
	if err := r.dbWithContext(ctx).Model(&QuestionModel{}).Where("id = ?", id).Updates(map[string]any{
		"updated_by": actorID,
		"updated_at": time.Now(),
		"is_active":  false,
	}).Delete(&QuestionModel{}, "id = ?", id).Error; err != nil {
		return helper.Internal("delete_question_failed", "failed to delete question", err)
	}
	return nil
}

func (r *questionRepository) FindByID(ctx context.Context, id uuid.UUID) (domain.Question, error) {
	var model QuestionModel
	if err := r.dbWithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return domain.Question{}, mapGormNotFound(err, "question_not_found", "question not found", "find_question_failed", "failed to find question")
	}
	return r.loadQuestionRelations(ctx, model, true, true)
}

func (r *questionRepository) FindAll(ctx context.Context, query domain.QuestionQuery) (domain.PageResult[domain.Question], error) {
	db := r.dbWithContext(ctx).Model(&QuestionModel{})
	if query.QuestionSetID != nil {
		db = db.Where("question_set_id = ?", *query.QuestionSetID)
	}
	if query.Type != nil {
		db = db.Where("type = ?", *query.Type)
	}
	if query.IsActive != nil {
		db = db.Where("is_active = ?", *query.IsActive)
	}
	if query.Search != "" {
		search := "%" + strings.ToLower(query.Search) + "%"
		db = db.Where("LOWER(question_text) LIKE ? OR LOWER(explanation) LIKE ?", search, search)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return domain.PageResult[domain.Question]{}, helper.Internal("count_questions_failed", "failed to count questions", err)
	}

	sortBy := safeSort(query.SortBy, []string{"type", "difficulty", "order_no", "is_active", "created_at", "updated_at"}, "order_no")
	sortOrder := safeSortOrder(query.SortOrder)

	var models []QuestionModel
	if err := db.Order("tbl_questions." + sortBy + " " + sortOrder).
		Order("tbl_questions.id ASC").
		Offset(query.Paging.Offset()).
		Limit(query.Paging.Limit).
		Find(&models).Error; err != nil {
		return domain.PageResult[domain.Question]{}, helper.Internal("list_questions_failed", "failed to list questions", err)
	}

	items := make([]domain.Question, 0, len(models))
	for _, model := range models {
		question, err := r.loadQuestionRelations(ctx, model, query.IncludeChoices, query.IncludeTags)
		if err != nil {
			return domain.PageResult[domain.Question]{}, err
		}
		items = append(items, question)
	}

	return domain.PageResult[domain.Question]{
		Items: items,
		Paging: domain.Paging{
			Page:  query.Paging.Page,
			Limit: query.Paging.Limit,
			Total: total,
		},
	}, nil
}

// ReplaceChoices applies "edit mode" semantics for question choices.
// Existing rows with a known ID are updated in place, new rows are created,
// and rows omitted from the payload are soft-deleted.
func (r *questionRepository) ReplaceChoices(ctx context.Context, questionID uuid.UUID, choices []domain.QuestionChoiceInput, actorID string) error {
	db := r.dbWithContext(ctx)

	var existing []QuestionChoiceModel
	if err := db.Where("question_id = ?", questionID).Find(&existing).Error; err != nil {
		return helper.Internal("load_question_choices_failed", "failed to load question choices", err)
	}

	existingByID := make(map[uuid.UUID]QuestionChoiceModel, len(existing))
	for _, item := range existing {
		existingByID[item.ID] = item
	}

	now := time.Now()
	keepIDs := make(map[uuid.UUID]struct{}, len(choices))
	for _, choice := range choices {
		if choice.ID != uuid.Nil {
			if current, ok := existingByID[choice.ID]; ok {
				current.ChoiceText = choice.ChoiceText
				current.ChoiceOrder = choice.ChoiceOrder
				current.IsCorrect = choice.IsCorrect
				current.UpdatedBy = actorID
				current.UpdatedAt = now
				if err := db.Save(&current).Error; err != nil {
					return helper.Internal("update_question_choice_failed", "failed to update question choice", err)
				}
				keepIDs[current.ID] = struct{}{}
				continue
			}
		}

		model := QuestionChoiceModel{
			ID:          choice.ID,
			QuestionID:  questionID,
			ChoiceText:  choice.ChoiceText,
			ChoiceOrder: choice.ChoiceOrder,
			IsCorrect:   choice.IsCorrect,
			CreatedBy:   actorID,
			UpdatedBy:   actorID,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		if model.ID == uuid.Nil {
			model.ID = uuid.New()
		}
		if err := db.Create(&model).Error; err != nil {
			return helper.Internal("create_question_choice_failed", "failed to create question choice", err)
		}
		keepIDs[model.ID] = struct{}{}
	}

	for _, item := range existing {
		if _, ok := keepIDs[item.ID]; ok {
			continue
		}
		if err := db.Model(&QuestionChoiceModel{}).Where("id = ?", item.ID).Updates(map[string]any{
			"updated_by": actorID,
			"updated_at": now,
		}).Delete(&QuestionChoiceModel{}, "id = ?", item.ID).Error; err != nil {
			return helper.Internal("delete_question_choice_failed", "failed to delete question choice", err)
		}
	}

	return nil
}

// ReplaceTags handles the composite primary key table carefully.
// If a tag relation was soft-deleted before, we restore it instead of creating
// a duplicate row that would fail on the same primary key.
func (r *questionRepository) ReplaceTags(ctx context.Context, questionID uuid.UUID, tagIDs []uuid.UUID, actorID string) error {
	db := r.dbWithContext(ctx)

	var existing []QuestionTagModel
	if err := db.Unscoped().Where("question_id = ?", questionID).Find(&existing).Error; err != nil {
		return helper.Internal("load_question_tags_failed", "failed to load question tags", err)
	}

	desired := make(map[uuid.UUID]struct{}, len(tagIDs))
	for _, tagID := range tagIDs {
		if tagID == uuid.Nil {
			continue
		}
		desired[tagID] = struct{}{}
	}

	now := time.Now()
	existingByTagID := make(map[uuid.UUID]QuestionTagModel, len(existing))
	for _, item := range existing {
		existingByTagID[item.TagID] = item
	}

	for tagID := range desired {
		if current, ok := existingByTagID[tagID]; ok {
			updates := map[string]any{
				"updated_by": actorID,
				"updated_at": now,
			}
			if current.DeletedAt.Valid {
				updates["deleted_at"] = nil
			}
			if err := db.Unscoped().Model(&QuestionTagModel{}).
				Where("question_id = ? AND tag_id = ?", questionID, tagID).
				Updates(updates).Error; err != nil {
				return helper.Internal("restore_question_tag_failed", "failed to restore question tag", err)
			}
			continue
		}

		model := QuestionTagModel{
			QuestionID: questionID,
			TagID:      tagID,
			CreatedBy:  actorID,
			UpdatedBy:  actorID,
			CreatedAt:  now,
			UpdatedAt:  now,
		}
		if err := db.Create(&model).Error; err != nil {
			return helper.Internal("create_question_tag_failed", "failed to create question tag", err)
		}
	}

	for _, item := range existing {
		if item.DeletedAt.Valid {
			continue
		}
		if _, ok := desired[item.TagID]; ok {
			continue
		}
		if err := db.Model(&QuestionTagModel{}).
			Where("question_id = ? AND tag_id = ?", questionID, item.TagID).
			Updates(map[string]any{
				"updated_by": actorID,
				"updated_at": now,
			}).
			Delete(&QuestionTagModel{}, "question_id = ? AND tag_id = ?", questionID, item.TagID).Error; err != nil {
			return helper.Internal("delete_question_tag_failed", "failed to delete question tag", err)
		}
	}

	return nil
}

func (r *questionRepository) loadQuestionRelations(ctx context.Context, model QuestionModel, includeChoices bool, includeTags bool) (domain.Question, error) {
	question := toDomainQuestion(model)

	if includeChoices {
		choices, err := r.loadChoicesByQuestionID(ctx, model.ID)
		if err != nil {
			return domain.Question{}, err
		}
		question.Choices = choices
	}
	if includeTags {
		tags, err := r.loadTagsByQuestionID(ctx, model.ID)
		if err != nil {
			return domain.Question{}, err
		}
		question.Tags = tags
	}

	return question, nil
}

func (r *questionRepository) loadChoicesByQuestionID(ctx context.Context, questionID uuid.UUID) ([]domain.QuestionChoice, error) {
	var models []QuestionChoiceModel
	if err := r.dbWithContext(ctx).Where("question_id = ?", questionID).Order("choice_order ASC, id ASC").Find(&models).Error; err != nil {
		return nil, helper.Internal("load_question_choices_failed", "failed to load question choices", err)
	}

	items := make([]domain.QuestionChoice, 0, len(models))
	for _, model := range models {
		items = append(items, toDomainQuestionChoice(model))
	}
	return items, nil
}

func (r *questionRepository) loadTagsByQuestionID(ctx context.Context, questionID uuid.UUID) ([]domain.Tag, error) {
	var models []TagModel
	if err := r.dbWithContext(ctx).
		Table("tbl_tags").
		Joins("JOIN tbl_question_tags ON tbl_question_tags.tag_id = tbl_tags.id AND tbl_question_tags.deleted_at IS NULL").
		Where("tbl_question_tags.question_id = ?", questionID).
		Order("tbl_tags.name ASC").
		Find(&models).Error; err != nil {
		return nil, helper.Internal("load_question_tags_failed", "failed to load question tags", err)
	}

	items := make([]domain.Tag, 0, len(models))
	for _, model := range models {
		items = append(items, toDomainTag(model))
	}
	return items, nil
}

type questionChoiceRepository struct {
	*baseRepository
}

func (r *questionChoiceRepository) Create(ctx context.Context, choice domain.QuestionChoice) (domain.QuestionChoice, error) {
	now := time.Now()
	model := QuestionChoiceModel{
		ID:          choice.ID,
		QuestionID:  choice.QuestionID,
		ChoiceText:  choice.ChoiceText,
		ChoiceOrder: choice.ChoiceOrder,
		IsCorrect:   choice.IsCorrect,
		CreatedBy:   choice.CreatedBy,
		UpdatedBy:   choice.UpdatedBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if model.ID == uuid.Nil {
		model.ID = uuid.New()
	}
	if err := r.dbWithContext(ctx).Create(&model).Error; err != nil {
		return domain.QuestionChoice{}, helper.Internal("create_question_choice_failed", "failed to create question choice", err)
	}
	return r.FindByID(ctx, model.ID)
}

func (r *questionChoiceRepository) Update(ctx context.Context, input domain.QuestionChoiceUpdateInput) (domain.QuestionChoice, error) {
	var model QuestionChoiceModel
	if err := r.dbWithContext(ctx).First(&model, "id = ?", input.ID).Error; err != nil {
		return domain.QuestionChoice{}, mapGormNotFound(err, "question_choice_not_found", "question choice not found", "find_question_choice_failed", "failed to load question choice")
	}
	if input.ChoiceText.Set {
		model.ChoiceText = input.ChoiceText.Value
	}
	if input.ChoiceOrder.Set {
		model.ChoiceOrder = input.ChoiceOrder.Value
	}
	if input.IsCorrect.Set {
		model.IsCorrect = input.IsCorrect.Value
	}
	model.UpdatedBy = input.ActorID
	model.UpdatedAt = time.Now()

	if err := r.dbWithContext(ctx).Save(&model).Error; err != nil {
		return domain.QuestionChoice{}, helper.Internal("update_question_choice_failed", "failed to update question choice", err)
	}
	return r.FindByID(ctx, model.ID)
}

func (r *questionChoiceRepository) Delete(ctx context.Context, id uuid.UUID, actorID string) error {
	if err := r.dbWithContext(ctx).Model(&QuestionChoiceModel{}).Where("id = ?", id).Updates(map[string]any{
		"updated_by": actorID,
		"updated_at": time.Now(),
	}).Delete(&QuestionChoiceModel{}, "id = ?", id).Error; err != nil {
		return helper.Internal("delete_question_choice_failed", "failed to delete question choice", err)
	}
	return nil
}

func (r *questionChoiceRepository) FindByID(ctx context.Context, id uuid.UUID) (domain.QuestionChoice, error) {
	var model QuestionChoiceModel
	if err := r.dbWithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return domain.QuestionChoice{}, mapGormNotFound(err, "question_choice_not_found", "question choice not found", "find_question_choice_failed", "failed to find question choice")
	}
	return toDomainQuestionChoice(model), nil
}

func (r *questionChoiceRepository) FindAll(ctx context.Context, query domain.QuestionChoiceQuery) (domain.PageResult[domain.QuestionChoice], error) {
	db := r.dbWithContext(ctx).Model(&QuestionChoiceModel{})
	if query.QuestionID != nil {
		db = db.Where("question_id = ?", *query.QuestionID)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return domain.PageResult[domain.QuestionChoice]{}, helper.Internal("count_question_choices_failed", "failed to count question choices", err)
	}

	sortBy := safeSort(query.SortBy, []string{"choice_order", "created_at", "updated_at"}, "choice_order")
	sortOrder := safeSortOrder(query.SortOrder)

	var models []QuestionChoiceModel
	if err := db.Order("tbl_question_choices." + sortBy + " " + sortOrder).
		Order("tbl_question_choices.id ASC").
		Offset(query.Paging.Offset()).
		Limit(query.Paging.Limit).
		Find(&models).Error; err != nil {
		return domain.PageResult[domain.QuestionChoice]{}, helper.Internal("list_question_choices_failed", "failed to list question choices", err)
	}

	items := make([]domain.QuestionChoice, 0, len(models))
	for _, model := range models {
		items = append(items, toDomainQuestionChoice(model))
	}

	return domain.PageResult[domain.QuestionChoice]{
		Items: items,
		Paging: domain.Paging{
			Page:  query.Paging.Page,
			Limit: query.Paging.Limit,
			Total: total,
		},
	}, nil
}

type tagRepository struct {
	*baseRepository
}

func (r *tagRepository) Create(ctx context.Context, tag domain.Tag) (domain.Tag, error) {
	now := time.Now()
	model := TagModel{
		ID:        tag.ID,
		Name:      tag.Name,
		CreatedBy: tag.CreatedBy,
		UpdatedBy: tag.UpdatedBy,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if model.ID == uuid.Nil {
		model.ID = uuid.New()
	}
	if err := r.dbWithContext(ctx).Create(&model).Error; err != nil {
		return domain.Tag{}, helper.Internal("create_tag_failed", "failed to create tag", err)
	}
	return r.FindByID(ctx, model.ID)
}

func (r *tagRepository) Update(ctx context.Context, input domain.TagUpdateInput) (domain.Tag, error) {
	var model TagModel
	if err := r.dbWithContext(ctx).First(&model, "id = ?", input.ID).Error; err != nil {
		return domain.Tag{}, mapGormNotFound(err, "tag_not_found", "tag not found", "find_tag_failed", "failed to load tag")
	}
	if input.Name.Set {
		model.Name = input.Name.Value
	}
	model.UpdatedBy = input.ActorID
	model.UpdatedAt = time.Now()

	if err := r.dbWithContext(ctx).Save(&model).Error; err != nil {
		return domain.Tag{}, helper.Internal("update_tag_failed", "failed to update tag", err)
	}
	return r.FindByID(ctx, model.ID)
}

func (r *tagRepository) Delete(ctx context.Context, id uuid.UUID, actorID string) error {
	if err := r.dbWithContext(ctx).Model(&TagModel{}).Where("id = ?", id).Updates(map[string]any{
		"updated_by": actorID,
		"updated_at": time.Now(),
	}).Delete(&TagModel{}, "id = ?", id).Error; err != nil {
		return helper.Internal("delete_tag_failed", "failed to delete tag", err)
	}
	return nil
}

func (r *tagRepository) FindByID(ctx context.Context, id uuid.UUID) (domain.Tag, error) {
	var model TagModel
	if err := r.dbWithContext(ctx).First(&model, "id = ?", id).Error; err != nil {
		return domain.Tag{}, mapGormNotFound(err, "tag_not_found", "tag not found", "find_tag_failed", "failed to find tag")
	}
	return toDomainTag(model), nil
}

func (r *tagRepository) FindAll(ctx context.Context, query domain.TagQuery) (domain.PageResult[domain.Tag], error) {
	db := r.dbWithContext(ctx).Model(&TagModel{})
	if query.Search != "" {
		search := "%" + strings.ToLower(query.Search) + "%"
		db = db.Where("LOWER(name) LIKE ?", search)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return domain.PageResult[domain.Tag]{}, helper.Internal("count_tags_failed", "failed to count tags", err)
	}

	sortBy := safeSort(query.SortBy, []string{"name", "created_at", "updated_at"}, "name")
	sortOrder := safeSortOrder(query.SortOrder)

	var models []TagModel
	if err := db.Order("tbl_tags." + sortBy + " " + sortOrder).
		Order("tbl_tags.id ASC").
		Offset(query.Paging.Offset()).
		Limit(query.Paging.Limit).
		Find(&models).Error; err != nil {
		return domain.PageResult[domain.Tag]{}, helper.Internal("list_tags_failed", "failed to list tags", err)
	}

	items := make([]domain.Tag, 0, len(models))
	for _, model := range models {
		items = append(items, toDomainTag(model))
	}

	return domain.PageResult[domain.Tag]{
		Items: items,
		Paging: domain.Paging{
			Page:  query.Paging.Page,
			Limit: query.Paging.Limit,
			Total: total,
		},
	}, nil
}

func (r *tagRepository) ExistsByName(ctx context.Context, name string, excludeID *uuid.UUID) (bool, error) {
	return existsByStringField(ctx, r.dbWithContext, &TagModel{}, "name", name, excludeID, "exists_tag_failed", "failed to check tag uniqueness")
}

func toDomainSection(model SectionModel) domain.Section {
	item := domain.Section{
		ID:          model.ID,
		Slug:        model.Slug,
		Title:       model.Title,
		Description: model.Description,
		OrderNo:     model.OrderNo,
		IsPublished: model.IsPublished,
		AuditFields: toAuditFields(model.CreatedAt, model.UpdatedAt, model.DeletedAt, model.CreatedBy, model.UpdatedBy),
	}
	return item
}

func toDomainLesson(model LessonModel) domain.Lesson {
	item := domain.Lesson{
		ID:          model.ID,
		SectionID:   model.SectionID,
		Slug:        model.Slug,
		Title:       model.Title,
		Description: model.Description,
		OrderNo:     model.OrderNo,
		IsPublished: model.IsPublished,
		AuditFields: toAuditFields(model.CreatedAt, model.UpdatedAt, model.DeletedAt, model.CreatedBy, model.UpdatedBy),
	}
	return item
}

func toDomainUnit(model UnitModel) domain.Unit {
	item := domain.Unit{
		ID:          model.ID,
		LessonID:    model.LessonID,
		Slug:        model.Slug,
		Title:       model.Title,
		Description: model.Description,
		OrderNo:     model.OrderNo,
		IsPublished: model.IsPublished,
		AuditFields: toAuditFields(model.CreatedAt, model.UpdatedAt, model.DeletedAt, model.CreatedBy, model.UpdatedBy),
	}
	return item
}

func toDomainQuestionSet(model QuestionSetModel) domain.QuestionSet {
	item := domain.QuestionSet{
		ID:               model.ID,
		UnitID:           model.UnitID,
		Slug:             model.Slug,
		Title:            model.Title,
		Description:      model.Description,
		OrderNo:          model.OrderNo,
		EstimatedSeconds: model.EstimatedSeconds,
		IsPublished:      model.IsPublished,
		Version:          model.Version,
		AuditFields:      toAuditFields(model.CreatedAt, model.UpdatedAt, model.DeletedAt, model.CreatedBy, model.UpdatedBy),
	}
	return item
}

func toDomainQuestion(model QuestionModel) domain.Question {
	item := domain.Question{
		ID:            model.ID,
		QuestionSetID: model.QuestionSetID,
		Type:          model.Type,
		QuestionText:  model.QuestionText,
		BlankPosition: model.BlankPosition,
		Explanation:   model.Explanation,
		ImageURL:      model.ImageURL,
		Difficulty:    model.Difficulty,
		OrderNo:       model.OrderNo,
		IsActive:      model.IsActive,
		AuditFields:   toAuditFields(model.CreatedAt, model.UpdatedAt, model.DeletedAt, model.CreatedBy, model.UpdatedBy),
	}
	return item
}

func toDomainQuestionChoice(model QuestionChoiceModel) domain.QuestionChoice {
	item := domain.QuestionChoice{
		ID:          model.ID,
		QuestionID:  model.QuestionID,
		ChoiceText:  model.ChoiceText,
		ChoiceOrder: model.ChoiceOrder,
		IsCorrect:   model.IsCorrect,
		AuditFields: toAuditFields(model.CreatedAt, model.UpdatedAt, model.DeletedAt, model.CreatedBy, model.UpdatedBy),
	}
	return item
}

func toDomainTag(model TagModel) domain.Tag {
	item := domain.Tag{
		ID:          model.ID,
		Name:        model.Name,
		AuditFields: toAuditFields(model.CreatedAt, model.UpdatedAt, model.DeletedAt, model.CreatedBy, model.UpdatedBy),
	}
	return item
}

func toAuditFields(createdAt, updatedAt time.Time, deletedAt gorm.DeletedAt, createdBy, updatedBy string) domain.AuditFields {
	fields := domain.AuditFields{
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		CreatedBy: createdBy,
		UpdatedBy: updatedBy,
	}
	if deletedAt.Valid {
		value := deletedAt.Time
		fields.DeletedAt = &value
	}
	return fields
}

func mapGormNotFound(err error, notFoundCode, notFoundMessage, internalCode, internalMessage string) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return helper.NotFound(notFoundCode, notFoundMessage, err)
	}
	return helper.Internal(internalCode, internalMessage, err)
}

func safeSort(value string, allowed []string, fallback string) string {
	for _, item := range allowed {
		if value == item {
			return value
		}
	}
	return fallback
}

func safeSortOrder(value string) string {
	if strings.ToLower(value) == "desc" {
		return "desc"
	}
	return "asc"
}

func existsByStringField(
	ctx context.Context,
	dbWithContext func(context.Context) *gorm.DB,
	model any,
	field string,
	value string,
	excludeID *uuid.UUID,
	internalCode string,
	internalMessage string,
) (bool, error) {
	db := dbWithContext(ctx).Model(model).Where("LOWER("+field+") = ?", strings.ToLower(strings.TrimSpace(value)))
	if excludeID != nil {
		db = db.Where("id <> ?", *excludeID)
	}
	var count int64
	if err := db.Count(&count).Error; err != nil {
		return false, helper.Internal(internalCode, internalMessage, err)
	}
	return count > 0, nil
}
