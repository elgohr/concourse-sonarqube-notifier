package shared

type Source struct {
	Target     string `json:"target"`
	SonarToken string `json:"sonartoken"`
	Component  string `json:"component"`
	Metrics    string `json:"metrics"`
}

func (s *Source) Valid() bool {
	return len(s.Component) != 0 &&
		len(s.Metrics) != 0 &&
		len(s.Target) != 0 &&
		len(s.SonarToken) != 0
}

type Version map[string]string
