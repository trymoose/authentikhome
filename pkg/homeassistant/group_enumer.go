// Code generated by "enumer -type Group -linecomment"; DO NOT EDIT.

package homeassistant

import (
	"fmt"
	"strings"
)

const _GroupName = "system-userssystem-admin"

var _GroupIndex = [...]uint8{0, 12, 24}

const _GroupLowerName = "system-userssystem-admin"

func (i Group) String() string {
	if i < 0 || i >= Group(len(_GroupIndex)-1) {
		return fmt.Sprintf("Group(%d)", i)
	}
	return _GroupName[_GroupIndex[i]:_GroupIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _GroupNoOp() {
	var x [1]struct{}
	_ = x[GroupUsers-(0)]
	_ = x[GroupAdmin-(1)]
}

var _GroupValues = []Group{GroupUsers, GroupAdmin}

var _GroupNameToValueMap = map[string]Group{
	_GroupName[0:12]:       GroupUsers,
	_GroupLowerName[0:12]:  GroupUsers,
	_GroupName[12:24]:      GroupAdmin,
	_GroupLowerName[12:24]: GroupAdmin,
}

var _GroupNames = []string{
	_GroupName[0:12],
	_GroupName[12:24],
}

// GroupString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func GroupString(s string) (Group, error) {
	if val, ok := _GroupNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _GroupNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to Group values", s)
}

// GroupValues returns all values of the enum
func GroupValues() []Group {
	return _GroupValues
}

// GroupStrings returns a slice of all String values of the enum
func GroupStrings() []string {
	strs := make([]string, len(_GroupNames))
	copy(strs, _GroupNames)
	return strs
}

// IsAGroup returns "true" if the value is listed in the enum definition. "false" otherwise
func (i Group) IsAGroup() bool {
	for _, v := range _GroupValues {
		if i == v {
			return true
		}
	}
	return false
}
