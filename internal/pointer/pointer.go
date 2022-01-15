// Package pointer is a helper package to convert concrete types to pointers
package pointer

// IntP takes in an int and returns it's pointer
func IntP(in int) (out *int) {
	out = &in
	return
}

// StringP takes in a string and returns it's pointer
func StringP(in string) (out *string) {
	out = &in
	return
}
