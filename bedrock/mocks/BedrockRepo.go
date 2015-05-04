package mocks

import "github.com/austinpray/bump-bedrock/Godeps/_workspace/src/github.com/stretchr/testify/mock"

type BedrockRepo struct {
	mock.Mock
}

func (m *BedrockRepo) UpdateWordPressVersion(version string) string {
	ret := m.Called(version)

	r0 := ret.Get(0).(string)

	return r0
}
