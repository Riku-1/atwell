// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// YahooJapanAuthInfrastructure is an autogenerated mock type for the YahooJapanAuthInfrastructure type
type YahooJapanAuthInfrastructure struct {
	mock.Mock
}

// GetPublicKeyList provides a mock function with given fields:
func (_m *YahooJapanAuthInfrastructure) GetPublicKeyList() (*http.Response, error) {
	ret := _m.Called()

	var r0 *http.Response
	if rf, ok := ret.Get(0).(func() *http.Response); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*http.Response)
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

// Token provides a mock function with given fields: code
func (_m *YahooJapanAuthInfrastructure) Token(code string) (*http.Response, error) {
	ret := _m.Called(code)

	var r0 *http.Response
	if rf, ok := ret.Get(0).(func(string) *http.Response); ok {
		r0 = rf(code)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*http.Response)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(code)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UserInfo provides a mock function with given fields: accessToken
func (_m *YahooJapanAuthInfrastructure) UserInfo(accessToken string) (*http.Response, error) {
	ret := _m.Called(accessToken)

	var r0 *http.Response
	if rf, ok := ret.Get(0).(func(string) *http.Response); ok {
		r0 = rf(accessToken)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*http.Response)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(accessToken)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
