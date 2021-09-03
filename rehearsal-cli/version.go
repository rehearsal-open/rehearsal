// rehearsal-cli/version.go
// Copyright (C) 2021 Kasai Koji

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import "fmt"

const InfomationPlane = `
 /        ____     ____    _  _     ____    __     ____     ___     __     _    
 \(o)    /____\   /____\  //  \\   /____\  /__\   /____\   /___\   /__\   //    
   \\_  //____\\ //_____ //____\\ //_____ //  \\ //____\\ //___   //  \\ //     
   /\  // \ ___//______//______ //______///___//// \ ___//____ \ //___////      
  /  \ \\_ \\__ \\_____ \\  _ ////_____ / ____ \\\_ \\__ ____/ // ____ \\\______ 
_/\  /\_\_\ \__\ \____/  \\ \_/ \_____//_/    \_\\_\ \__\\____//_/    \_\\_____/ 


rehearsal v1.202109 / process-conecting test tool
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

const Abstruct = `
rehearsal v1.202109 Copyright (C) 2021  Kasai Koji

This program comes with ABSOLUTELY NO WARRANTY; for details type 'rehearsal-cli version'.
This is free software, and you are welcome to redistribute it
under certain conditions; type 'rehearsal-cli version' for details.
`

func Version() {
	fmt.Println(InfomationPlane)
}
