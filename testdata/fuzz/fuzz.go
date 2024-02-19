package fuzzing

import "github.com/PratikBhusal/glamour"

func Fuzz(data []byte) int {
	_, err := glamour.RenderBytes(data, glamour.DarkStyle)
	if err != nil {
		return 0
	}
	return 1
}
