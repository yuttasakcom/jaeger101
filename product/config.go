package main

import "os"

func configEnv() {
	os.Setenv("PROJECT_NAME", "product")
	os.Setenv("LOGGING_SERVICE", "product")
}
