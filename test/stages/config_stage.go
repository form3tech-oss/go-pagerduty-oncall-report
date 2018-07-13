package stages

import "testing"

type ConfigStage struct {
	t *testing.T
}

func ConfigTest(t *testing.T) (*ConfigStage, *ConfigStage, *ConfigStage) {
	stage := &ConfigStage{
		t: t,
	}

	return stage, stage, stage
}

func (s *ConfigStage) And() *ConfigStage {
	return s
}

//func (s *ConfigStage) XX() *ConfigStage {
//	return s
//}
