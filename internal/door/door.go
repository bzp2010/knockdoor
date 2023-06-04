package door

// Door is an interface used to represent the "door" mechanism
type Door interface {
	Open(ip string) error
}
