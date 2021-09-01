// tasks/buffer/packet_base.go
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

package buffer

import "github.com/rehearsal-open/rehearsal/entities"

func (pb *packetBase) Close() error {
	pb.nClosed++
	if pb.nClosed >= pb.nSend {
		pb.bytes = nil
	}
	return nil
}

func (p *packetBase) Sender() (*entities.Task, entities.TaskElement) {
	return p.buffer.task, p.buffer.element
}
