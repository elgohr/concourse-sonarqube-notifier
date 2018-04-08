package shared

type Source struct {
	Target     string `json:"target"`
	SonarToken string `json:"sonartoken"`
	Component  string `json:"component"`
	Metrics    string `json:"metrics"`
}

func (s *Source) Valid() bool {
	if len(s.Component) == 0 ||
		len(s.Metrics) == 0 ||
		len(s.Target) == 0 ||
		len(s.SonarToken) == 0 {
		return false
	}
	return true
}

type Version map[string]string

func HasError(err error) bool {
	return err != nil
}
