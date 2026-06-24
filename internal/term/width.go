package term

import "strings"

func PadRight(s string, width int) string {
	currentWidth := DisplayWidth(s)
	if currentWidth >= width {
		return s
	}
	return s + strings.Repeat(" ", width-currentWidth)
}

func DisplayWidth(s string) int {
	width := 0
	for _, r := range s {
		width += runeWidth(r)
	}
	return width
}

func runeWidth(r rune) int {
	switch {
	case r >= 0x1100 && r <= 0x115F:
		return 2
	case r >= 0x2E80 && r <= 0xA4CF:
		return 2
	case r >= 0xAC00 && r <= 0xD7A3:
		return 2
	case r >= 0xF900 && r <= 0xFAFF:
		return 2
	case r >= 0xFE10 && r <= 0xFE6F:
		return 2
	case r >= 0xFF00 && r <= 0xFF60:
		return 2
	case r >= 0xFFE0 && r <= 0xFFE6:
		return 2
	default:
		return 1
	}
}
