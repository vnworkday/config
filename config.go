package config

// FromEnv reads environment variables and populates the config.
func FromEnv() *Builder {
	return newBuilder().FromEnv()
}

// FromFile reads a file and populates the config.
func FromFile(file string) *Builder {
	return newBuilder().FromFile(file)
}
