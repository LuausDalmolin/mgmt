// Mgmt
// Copyright (C) 2013-2018+ James Shubin and the project contributors
// Written by James Shubin <james@shubin.ca> and the project contributors
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package resources

// Edge is a struct that represents a graph's edge.
type Edge struct {
	Name   string
	Notify bool // should we send a refresh notification along this edge?

	refresh bool // is there a notify pending for the dest vertex ?
}

// String is a required method of the Edge interface that we must fulfill.
func (obj *Edge) String() string {
	return obj.Name
}

// Compare returns true if two edges are equivalent. Otherwise it returns false.
func (obj *Edge) Compare(edge *Edge) bool {
	if obj.Name != edge.Name {
		return false
	}
	if obj.Notify != edge.Notify {
		return false
	}
	// FIXME: should we compare this as well?
	//if obj.refresh != edge.refresh {
	//	return false
	//}
	return true
}

// Refresh returns the pending refresh status of this edge.
func (obj *Edge) Refresh() bool {
	return obj.refresh
}

// SetRefresh sets the pending refresh status of this edge.
func (obj *Edge) SetRefresh(b bool) {
	obj.refresh = b
}
