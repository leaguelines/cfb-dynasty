package dynasty

import "strconv"

// EnumSchema describes enum members for a schema field.
type EnumSchema struct {
	Name      string
	MaxLength int
	Members   []EnumMember
}

// EnumMember is one value in an enum schema.
type EnumMember struct {
	Name              string
	Value             int
	UnformattedValue  string
}

// NameForBits returns the member name for raw bits read from a record.
func (e *EnumSchema) NameForBits(bits uint32, bitLength int) string {
	if e == nil {
		return ""
	}
	unformatted := strconv.FormatUint(uint64(bits), 2)
	if e.MaxLength > 0 {
		if len(unformatted) > e.MaxLength {
			trimmed := unformatted[:len(unformatted)-e.MaxLength]
			for _, ch := range trimmed {
				if ch != '0' {
					return ""
				}
			}
			unformatted = unformatted[len(unformatted)-e.MaxLength:]
		} else {
			for len(unformatted) < e.MaxLength {
				unformatted = "0" + unformatted
			}
		}
	}
	return e.nameForUnformatted(unformatted)
}

func (e *EnumSchema) nameForUnformatted(unformatted string) string {
	var matches []EnumMember
	for _, member := range e.Members {
		if member.Name == "First_" || member.Name == "Last_" {
			continue
		}
		if member.UnformattedValue == unformatted {
			matches = append(matches, member)
		}
	}
	if len(matches) == 0 {
		return ""
	}
	for _, member := range matches {
		if member.Name == "" || member.Name[len(member.Name)-1] != '_' {
			return member.Name
		}
	}
	return matches[0].Name
}
