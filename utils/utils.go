package utils

func ExitOnErr(err error, msg string, args ...interface{}) {
	args = append(args, err)
	if err != nil {
		log.Fatalf(msg+" , error: %v", args...)
	}
}
