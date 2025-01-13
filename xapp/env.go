package xapp

type Environment int

const (
	EnvTest Environment = iota
	EnvStaging
	EnvProd
)
