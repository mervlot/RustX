package errors

func FormatSyntaxError(msg string, line int, column int) string {
	return "SyntaxError:\n" + msg + "\nat line " + itoa(line) + " column " + itoa(column)
}

func FormatRuntimeError(msg string, line int, column int) string {
	return "RuntimeError:\n" + msg + "\nat line " + itoa(line) + " column " + itoa(column)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	res := ""
	for n > 0 {
		rem := n % 10
		res = string(rune('0'+rem)) + res
		n = n / 10
	}
	return res
}
