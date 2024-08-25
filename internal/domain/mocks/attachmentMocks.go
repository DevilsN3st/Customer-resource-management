// Code generated by MockGen. DO NOT EDIT.
// Source: attachments.go

// Package mock_domain is a generated GoMock package.
package mock_domain

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/icrxz/crm-api-core/internal/domain"
)

// MockAttachmentRepository is a mock of AttachmentRepository interface.
type MockAttachmentRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAttachmentRepositoryMockRecorder
}

// MockAttachmentRepositoryMockRecorder is the mock recorder for MockAttachmentRepository.
type MockAttachmentRepositoryMockRecorder struct {
	mock *MockAttachmentRepository
}

// NewMockAttachmentRepository creates a new mock instance.
func NewMockAttachmentRepository(ctrl *gomock.Controller) *MockAttachmentRepository {
	mock := &MockAttachmentRepository{ctrl: ctrl}
	mock.recorder = &MockAttachmentRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAttachmentRepository) EXPECT() *MockAttachmentRepositoryMockRecorder {
	return m.recorder
}

// GetByCommentID mocks base method.
func (m *MockAttachmentRepository) GetByCommentID(ctx context.Context, commentID string) ([]domain.Attachment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByCommentID", ctx, commentID)
	ret0, _ := ret[0].([]domain.Attachment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByCommentID indicates an expected call of GetByCommentID.
func (mr *MockAttachmentRepositoryMockRecorder) GetByCommentID(ctx, commentID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByCommentID", reflect.TypeOf((*MockAttachmentRepository)(nil).GetByCommentID), ctx, commentID)
}

// GetByID mocks base method.
func (m *MockAttachmentRepository) GetByID(ctx context.Context, attachmentID string) (domain.Attachment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, attachmentID)
	ret0, _ := ret[0].(domain.Attachment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockAttachmentRepositoryMockRecorder) GetByID(ctx, attachmentID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockAttachmentRepository)(nil).GetByID), ctx, attachmentID)
}

// Save mocks base method.
func (m *MockAttachmentRepository) Save(ctx context.Context, attachment domain.Attachment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, attachment)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockAttachmentRepositoryMockRecorder) Save(ctx, attachment interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockAttachmentRepository)(nil).Save), ctx, attachment)
}

// SaveBatch mocks base method.
func (m *MockAttachmentRepository) SaveBatch(ctx context.Context, attachments []domain.Attachment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveBatch", ctx, attachments)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveBatch indicates an expected call of SaveBatch.
func (mr *MockAttachmentRepositoryMockRecorder) SaveBatch(ctx, attachments interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveBatch", reflect.TypeOf((*MockAttachmentRepository)(nil).SaveBatch), ctx, attachments)
}

// MockAttachmentBucket is a mock of AttachmentBucket interface.
type MockAttachmentBucket struct {
	ctrl     *gomock.Controller
	recorder *MockAttachmentBucketMockRecorder
}

// MockAttachmentBucketMockRecorder is the mock recorder for MockAttachmentBucket.
type MockAttachmentBucketMockRecorder struct {
	mock *MockAttachmentBucket
}

// NewMockAttachmentBucket creates a new mock instance.
func NewMockAttachmentBucket(ctrl *gomock.Controller) *MockAttachmentBucket {
	mock := &MockAttachmentBucket{ctrl: ctrl}
	mock.recorder = &MockAttachmentBucketMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAttachmentBucket) EXPECT() *MockAttachmentBucketMockRecorder {
	return m.recorder
}

// Download mocks base method.
func (m *MockAttachmentBucket) Download(ctx context.Context, attachmentID string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Download", ctx, attachmentID)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Download indicates an expected call of Download.
func (mr *MockAttachmentBucketMockRecorder) Download(ctx, attachmentID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Download", reflect.TypeOf((*MockAttachmentBucket)(nil).Download), ctx, attachmentID)
}
