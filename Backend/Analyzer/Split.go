package analyzer

func SplitCommand(cmd string, separator rune) []string {
	i := 0
	l := 0
	state := 0
	splitted := ""
	var res []string
	res = make([]string, 1)

	for i < len(cmd) {
		if cmd[i] == '\n' {
			break
		}

		switch state {
		case 0:
			if cmd[i] == byte(separator) {
				i++
				continue
			}
			state = 1
			splitted = splitted + string(cmd[i])
			break
		case 1:
			if cmd[i] == byte(separator) {
				res = append(res, splitted)
				l++
				state = 0
				i++
				splitted = ""
				continue
			}
			if cmd[i] == byte('"') {
				state = 2
				i++
				continue
			}
			splitted = splitted + string(cmd[i])
			break
		default:
			if cmd[i] == byte('"') {
				state = 1
				i++
				continue
			}
			splitted = splitted + string(cmd[i])
			break
		}

		i++
	}

	if cmd[i-1] != ' ' {
		res = append(res, splitted)
		l++
	}

	return res[1 : l+1]
}
