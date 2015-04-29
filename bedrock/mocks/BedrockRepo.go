package mocks

import "github.com/stretchr/testify/mock"

type BedrockRepo struct {
	mock.Mock
}

func (m *BedrockRepo) UpdateWordPressVersion(version string) {
	m.Called(version)
}
