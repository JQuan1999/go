package main

func entityParser(text string) string {
	var ret string
	idx := 0
	htmlMap := map[string]byte{
		"&quot;":  '"',
		"&apos;":  '\'',
		"&amp;":   '&',
		"&gt;":    '>',
		"&lt;":    '<',
		"&frasl;": '/',
	}
	textLen := len(text)
	for idx < textLen {
		if text[idx] == '&' {
			for key, value := range htmlMap {
				keyLen := len(key)
				if idx+keyLen <= textLen && text[idx:idx+keyLen] == key {
					idx += keyLen
					ret += string(value)
					break
				}
			}
			ret += text[idx : idx+1]
			idx++
		} else {
			ret += text[idx : idx+1]
			idx++
		}
	}
	return ret
}
