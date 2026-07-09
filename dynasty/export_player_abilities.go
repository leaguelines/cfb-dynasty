package dynasty

import "fmt"

func physicalAbilitiesFromRecord(record Record) []PlayerPhysicalAbilityExport {
	position := stringField(record, "Position")
	archetypeLabel := archetypeLabelFromRecord(record)
	var out []PlayerPhysicalAbilityExport
	for slot := 1; slot <= 5; slot++ {
		tier := normalizeEnum(stringField(record, fmt.Sprintf("PhysicalAbility%d", slot)))
		if tier == "" || tier == "None" || tier == "Invalid" {
			continue
		}
		export := PlayerPhysicalAbilityExport{
			Slot: slot,
			Tier: tier,
		}
		if name := physicalAbilityName(position, archetypeLabel, slot); name != "" {
			export.Name = name
		}
		out = append(out, export)
	}
	return out
}

func mentalAbilitiesFromRecord(record Record) []PlayerMentalAbilityExport {
	var out []PlayerMentalAbilityExport
	for slot := 1; slot <= 3; slot++ {
		name := normalizeEnum(stringField(record, fmt.Sprintf("MentalAbility%d", slot)))
		if name == "" || name == "None" || name == "Invalid" {
			continue
		}
		export := PlayerMentalAbilityExport{Name: name}
		if rank, ok := intFieldOK(record, fmt.Sprintf("MentalAbilityRank%d", slot)); ok && rank > 0 {
			export.Rank = &rank
		}
		out = append(out, export)
	}
	return out
}
