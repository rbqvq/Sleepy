package main

type Config struct {
	Web struct {
		Type   string // network type
		Listen string

		Cert, Key string
	}

	Security struct {
		AllowCORS bool
		Secret    string
	}

	System struct {
		DebugMode bool
	}
}
