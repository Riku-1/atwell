// Code generated by mockery v2.4.0. DO NOT EDIT.

package mocks

import (
	domain "golang-api/domain"

	mock "github.com/stretchr/testify/mock"
)

// ArticleUsecase is an autogenerated mock type for the ArticleUsecase type
type ArticleUsecase struct {
	mock.Mock
}

// GetAll provides a mock function with given fields:
func (_m *ArticleUsecase) GetAll() ([]domain.Article, error) {
	ret := _m.Called()

	var r0 []domain.Article
	if rf, ok := ret.Get(0).(func() []domain.Article); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Article)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
