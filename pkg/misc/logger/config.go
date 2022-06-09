package logger

// Config config
type Config struct {
	Level           int      `yaml:"level"`
	Development     bool     `yaml:"development"`
	Sampling        Sampling `yaml:"sampling"`
	OutputPath      []string `yaml:"outputPath"`
	ErrorOutputPath []string `yaml:"errorOutputPath"`
}

// Sampling Sampling
type Sampling struct {
	Initial    int `yaml:"initial"`
	Thereafter int `yaml:"thereafter"`
}
