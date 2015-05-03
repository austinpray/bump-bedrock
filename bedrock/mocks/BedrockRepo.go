package mocks

import "github.com/stretchr/testify/mock"

type BedrockRepo struct {
	mock.Mock
}

func (m *BedrockRepo) UpdateWordPressVersion(version string) string {
	ret := m.Called(version)

	r0 := ret.Get(0).(string)

	return r0
}
