package entity

import "strings"

const (
	TaskNameRegexp       string = "[a-zA-Z0-9][a-zA-Z_0-9]*"
	AboutInfomationPlane string = `
 ______________________
[_/_\_/_\_/_\_/_\_/_\_/]
 \ /  /    _     \  \ /    ____        ____        __  _        ____        ____        ____        ____         ____        _        
 \ / /    |o|     \ \ /   /____\      /____\      /_/  \\      /____\      /____\      /____\      /____\       /____\      //        
  \ /   __/\\__    \ /   //____\\    //____      //_____\\    //____      //    \\    //____\\    //______     //    \\    //         
  /_/      /\      \_\  // \ ___/   /______\  _ /________ \_ /______\  _ //______\\_ // \ ___/   /_______ \__ //______\\_ //          
 //_\     /  \     /_\\ \\_ \ \____ \\_______// \\   ___/  / \\_______// \ ______  / \\_ \ \____ ________\  / \ ______  / \\_________ 
 //_\  __/\  /\__  /_\\  \_\ \____/  \_______/   \\  \____/   \_______/   \\ \____/   \_\ \____/ \_________/   \\ \____/   \________/ 

Rehearsal v0.9.0 / process-connecting test tool
Copyright (C) 2021  Kasai Koji

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
`
	GeneralPulicLicenseAbstruct string = `
rehearsal v0.9.0  Copyright (C) 2021  Kasai Koji
This program comes with ABSOLUTELY NO WARRANTY; for details type 'rehearsal about'.
This is free software, and you are welcome to redistribute it
under certain conditions; type 'rehearsal about' for details.
	`
)

func AboutInfomation() string {
	strs := strings.Split(AboutInfomationPlane, "\n")
	for i, _ := range strs {
		l := len(strs[i])
		if i < 3 {
			strs[i] = "\x1b[33m" + strs[i]
		} else if i < 9 {
			strs[i] = "\x1b[31m" + strs[i][0:7] + "\x1b[0m" + strs[i][7:17] + "\x1b[31m" + strs[i][17:24] + "\x1b[36m" + strs[i][24:l]
		} else {
			strs[i] = "\x1b[0m" + strs[i]
		}
	}
	return strings.Join(strs, "\n")
}
