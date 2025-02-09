// Code generated by mockery v2.51.1. DO NOT EDIT.

package mocks

import (
	domain "go-banking/domain"
	dto "go-banking/dto"

	mock "github.com/stretchr/testify/mock"
)

// AccountRepository is an autogenerated mock type for the AccountRepository type
type AccountRepository struct {
	mock.Mock
}

// CreateAccount provides a mock function with given fields: account
func (_m *AccountRepository) CreateAccount(account domain.Account) (*domain.Account, error) {
	ret := _m.Called(account)

	if len(ret) == 0 {
		panic("no return value specified for CreateAccount")
	}

	var r0 *domain.Account
	var r1 error
	if rf, ok := ret.Get(0).(func(domain.Account) (*domain.Account, error)); ok {
		return rf(account)
	}
	if rf, ok := ret.Get(0).(func(domain.Account) *domain.Account); ok {
		r0 = rf(account)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Account)
		}
	}

	if rf, ok := ret.Get(1).(func(domain.Account) error); ok {
		r1 = rf(account)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAccountByCustomerID provides a mock function with given fields: id
func (_m *AccountRepository) GetAccountByCustomerID(id string) ([]domain.Account, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetAccountByCustomerID")
	}

	var r0 []domain.Account
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]domain.Account, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(string) []domain.Account); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Account)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAccountByID provides a mock function with given fields: id
func (_m *AccountRepository) GetAccountByID(id string) (*domain.Account, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetAccountByID")
	}

	var r0 *domain.Account
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*domain.Account, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(string) *domain.Account); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Account)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAccounts provides a mock function with no fields
func (_m *AccountRepository) GetAccounts() ([]dto.AccountWithCustomer, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetAccounts")
	}

	var r0 []dto.AccountWithCustomer
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]dto.AccountWithCustomer, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []dto.AccountWithCustomer); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]dto.AccountWithCustomer)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SoftDeleteAccount provides a mock function with given fields: account
func (_m *AccountRepository) SoftDeleteAccount(account domain.Account) (*domain.Account, error) {
	ret := _m.Called(account)

	if len(ret) == 0 {
		panic("no return value specified for SoftDeleteAccount")
	}

	var r0 *domain.Account
	var r1 error
	if rf, ok := ret.Get(0).(func(domain.Account) (*domain.Account, error)); ok {
		return rf(account)
	}
	if rf, ok := ret.Get(0).(func(domain.Account) *domain.Account); ok {
		r0 = rf(account)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Account)
		}
	}

	if rf, ok := ret.Get(1).(func(domain.Account) error); ok {
		r1 = rf(account)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateAccount provides a mock function with given fields: account
func (_m *AccountRepository) UpdateAccount(account domain.Account) (*domain.Account, error) {
	ret := _m.Called(account)

	if len(ret) == 0 {
		panic("no return value specified for UpdateAccount")
	}

	var r0 *domain.Account
	var r1 error
	if rf, ok := ret.Get(0).(func(domain.Account) (*domain.Account, error)); ok {
		return rf(account)
	}
	if rf, ok := ret.Get(0).(func(domain.Account) *domain.Account); ok {
		r0 = rf(account)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Account)
		}
	}

	if rf, ok := ret.Get(1).(func(domain.Account) error); ok {
		r1 = rf(account)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewAccountRepository creates a new instance of AccountRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAccountRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *AccountRepository {
	mock := &AccountRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
