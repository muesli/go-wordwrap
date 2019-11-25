package wordwrap

import (
	"unicode"
)

// WrapString wraps the given string within lim width in characters.
//
// Wrapping is currently naive and only happens at white-space. A future
// version of the library will implement smarter wrapping. This means that
// pathological cases can dramatically reach past the limit, such as a very
// long word.
func WrapString(s string, lim uint) string {
	// Initialize a buffer with a slightly larger size to account for breaks
	var buf string

	var current uint
	var wordBuf, spaceBuf string
	var ansi bool

	for _, char := range s {
		if char == '\x1B' {
			if len(spaceBuf) > 0 || len(wordBuf) > 0 {
				current += uint(len(spaceBuf) + len(wordBuf))
				buf += spaceBuf
				spaceBuf = ""
				buf += wordBuf
				wordBuf = ""
			}
			buf += string(char)
			ansi = true
		} else if ansi {
			buf += string(char)
			if (char >= 0x40 && char <= 0x5a) || (char >= 0x61 && char <= 0x7a) {
				ansi = false
			}
		} else if char == '\n' {
			if len(wordBuf) == 0 {
				if current+uint(len(spaceBuf)) > lim {
					current = 0
				} else {
					current += uint(len(spaceBuf))
					buf += spaceBuf
				}
				spaceBuf = ""
			} else {
				current += uint(len(spaceBuf) + len(wordBuf))
				buf += spaceBuf
				spaceBuf = ""
				buf += wordBuf
				wordBuf = ""
			}
			buf += string(char)
			current = 0
		} else if unicode.IsSpace(char) {
			if len(spaceBuf) == 0 || len(wordBuf) > 0 {
				current += uint(len(spaceBuf) + len(wordBuf))
				buf += spaceBuf
				spaceBuf = ""
				buf += wordBuf
				wordBuf = ""
			}

			spaceBuf += string(char)
		} else {
			wordBuf += string(char)

			if current+uint(len(spaceBuf)+len(wordBuf)) > lim && uint(len(wordBuf)) < lim {
				buf += "\n"
				current = 0
				spaceBuf = ""
			}
		}
	}

	if len(wordBuf) == 0 {
		if current+uint(len(spaceBuf)) <= lim {
			buf += spaceBuf
		}
	} else {
		buf += spaceBuf
		buf += wordBuf
	}

	return buf
}
