package pointer

func IntP(in int) (out *int) {
	out = &in
	return
}

func StringP(in string) (out *string) {
	out = &in
	return
}
