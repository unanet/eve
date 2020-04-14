package mux

type Config struct {
	Port        int    `split_words:"true" default:"8080"`
	MetricsPort int    `split_words:"true" default:"3000"`
	ServiceName string `split_words:"true" default:"eve"`
}
