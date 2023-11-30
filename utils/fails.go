package utils

func Failed(err *error, fails *[]string) {
	if *err != nil {
		*fails = append(*fails, (*err).Error())
		*err = nil
	}
}
