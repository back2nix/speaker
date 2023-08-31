package keylogger

const (
	L_CTRL  = "L_CTRL"
	L_SHIFT = "L_SHIFT"
	L_ALT   = "L_ALT"
	C       = "C"
	Z       = "Z"
	F       = "F"
	P       = "P"
)

var KeyCombos = [][]string{
	{L_CTRL, L_SHIFT, C},
	{L_CTRL, C},
	{L_ALT, C},
	{L_CTRL, Z},
	{L_ALT, Z},
	{L_ALT, F},
	{L_CTRL, L_ALT, P},
}

// KeyMap = map[string]bool{}

// func init() {
// 	for _, it := range KeyCombos {
// 		combined := strings.Join(it, " + ")
// 		KeyMap[combined] = true
// 	}
// }
