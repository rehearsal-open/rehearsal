package executer

import (
	"log"
)

func (manager *ExecuteManager) printOut(item *ExecuteItem, obj string) {
	str := item.Exec.GetName()
	for i := len(str); i <= manager.maxNameLen; i++ {
		str += " "
	}
	str += ":"
	str += obj
	if str[len(str)-1] == '\n' || str[len(str)-1] == '\r' {
		str = str[:len(str)-1]
	} else {
		str += "(no carriage return)"
	}
	log.Println(str)
}
