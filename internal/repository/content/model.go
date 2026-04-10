package content

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SectionModel struct {
	ID          uuid.UUID      `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	Slug        string         `gorm:"column:slug;size:255;uniqueIndex;not null"`
	Title       string         `gorm:"column:title;size:255;not null"`
	Description *string        `gorm:"column:description"`
	OrderNo     int            `gorm:"column:order_no;not null;default:0"`
	IsPublished bool           `gorm:"column:is_published;not null;default:false;index"`
	CreatedBy   string         `gorm:"column:created_by;size:255"`
	UpdatedBy   string         `gorm:"column:updated_by;size:255"`
	CreatedAt   time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (SectionModel) TableName() string {
	return "tbl_sections"
}

type LessonModel struct {
	ID          uuid.UUID      `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	SectionID   uuid.UUID      `gorm:"column:section_id;type:uuid;not null;index"`
	Slug        string         `gorm:"column:slug;size:255;not null"`
	Title       string         `gorm:"column:title;size:255;not null"`
	Description *string        `gorm:"column:description"`
	OrderNo     int            `gorm:"column:order_no;not null;default:0"`
	IsPublished bool           `gorm:"column:is_published;not null;default:false;index"`
	CreatedBy   string         `gorm:"column:created_by;size:255"`
	UpdatedBy   string         `gorm:"column:updated_by;size:255"`
	CreatedAt   time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (LessonModel) TableName() string {
	return "tbl_lessons"
}

type UnitModel struct {
	ID          uuid.UUID      `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	LessonID    uuid.UUID      `gorm:"column:lesson_id;type:uuid;not null;index"`
	Slug        string         `gorm:"column:slug;size:255;not null"`
	Title       string         `gorm:"column:title;size:255;not null"`
	Description *string        `gorm:"column:description"`
	OrderNo     int            `gorm:"column:order_no;not null;default:0"`
	IsPublished bool           `gorm:"column:is_published;not null;default:false;index"`
	CreatedBy   string         `gorm:"column:created_by;size:255"`
	UpdatedBy   string         `gorm:"column:updated_by;size:255"`
	CreatedAt   time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (UnitModel) TableName() string {
	return "tbl_units"
}

type QuestionSetModel struct {
	ID               uuid.UUID      `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	UnitID           uuid.UUID      `gorm:"column:unit_id;type:uuid;not null;index"`
	Slug             string         `gorm:"column:slug;size:255;not null"`
	Title            string         `gorm:"column:title;size:255;not null"`
	Description      *string        `gorm:"column:description"`
	OrderNo          int            `gorm:"column:order_no;not null;default:0"`
	EstimatedSeconds *int           `gorm:"column:estimated_seconds"`
	IsPublished      bool           `gorm:"column:is_published;not null;default:false;index"`
	Version          int            `gorm:"column:version;not null;default:1"`
	CreatedBy        string         `gorm:"column:created_by;size:255"`
	UpdatedBy        string         `gorm:"column:updated_by;size:255"`
	CreatedAt        time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt        time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt        gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (QuestionSetModel) TableName() string {
	return "tbl_question_sets"
}

type QuestionModel struct {
	ID            uuid.UUID      `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	QuestionSetID uuid.UUID      `gorm:"column:question_set_id;type:uuid;not null;index"`
	Type          string         `gorm:"column:type;size:64;not null;index"`
	QuestionText  string         `gorm:"column:question_text;not null"`
	BlankPosition *int           `gorm:"column:blank_position"`
	Explanation   *string        `gorm:"column:explanation"`
	ImageURL      *string        `gorm:"column:image_url"`
	Difficulty    int            `gorm:"column:difficulty;not null;default:1"`
	OrderNo       int            `gorm:"column:order_no;not null;default:0"`
	IsActive      bool           `gorm:"column:is_active;not null;default:true;index"`
	CreatedBy     string         `gorm:"column:created_by;size:255"`
	UpdatedBy     string         `gorm:"column:updated_by;size:255"`
	CreatedAt     time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (QuestionModel) TableName() string {
	return "tbl_questions"
}

type QuestionChoiceModel struct {
	ID          uuid.UUID      `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	QuestionID  uuid.UUID      `gorm:"column:question_id;type:uuid;not null;index"`
	ChoiceText  string         `gorm:"column:choice_text;not null"`
	ChoiceOrder int            `gorm:"column:choice_order;not null;default:0"`
	IsCorrect   bool           `gorm:"column:is_correct;not null;default:false"`
	CreatedBy   string         `gorm:"column:created_by;size:255"`
	UpdatedBy   string         `gorm:"column:updated_by;size:255"`
	CreatedAt   time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (QuestionChoiceModel) TableName() string {
	return "tbl_question_choices"
}

type TagModel struct {
	ID        uuid.UUID      `gorm:"column:id;type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string         `gorm:"column:name;size:255;uniqueIndex;not null"`
	CreatedBy string         `gorm:"column:created_by;size:255"`
	UpdatedBy string         `gorm:"column:updated_by;size:255"`
	CreatedAt time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (TagModel) TableName() string {
	return "tbl_tags"
}

type QuestionTagModel struct {
	QuestionID uuid.UUID      `gorm:"column:question_id;type:uuid;primaryKey"`
	TagID      uuid.UUID      `gorm:"column:tag_id;type:uuid;primaryKey"`
	CreatedBy  string         `gorm:"column:created_by;size:255"`
	UpdatedBy  string         `gorm:"column:updated_by;size:255"`
	CreatedAt  time.Time      `gorm:"column:created_at;not null"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;not null"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (QuestionTagModel) TableName() string {
	return "tbl_question_tags"
}
