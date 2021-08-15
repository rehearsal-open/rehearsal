// entities/task.go
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

package entities

func (t *Task) Name() string       { return t.name }
func (t *Task) Fullname() string   { return t.phase + "::" + t.name }
func (t *Task) Detail() TaskDetail { return t.detail }
func (t *Task) CheckFormat() error {
	if err := taskkindFilter(t.Kind); err != nil {
		return err
	} else if err := encodingFilter(t.Encoding); err != nil {
		return err
	} else if err := t.detail.CheckFormat(); err != nil {
		return err
	}
	return nil
}
