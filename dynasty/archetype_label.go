package dynasty

// archetypePTEntry maps a Player.PT_* trait flag to the CFB 26/27 UI archetype name.
// Order matters: the first matching flag wins when multiple are set.
var archetypePTEntries = []archetypePTEntry{
	// Quarterback
	{field: "PT_QBPOCKETPASSER", label: "Pocket Passer"},
	{field: "PT_QBBACKFIELDCREATOR", label: "Backfield Creator"},
	{field: "PT_QBDUALTHREAT", label: "Dual Threat"},
	{field: "PT_QBPURERUNNER", label: "Pure Runner"},

	// Halfback (RB)
	{field: "PT_HBPOWERBACK", label: "Contact Seeker"},
	{field: "PT_HBELUSIVEBACK", label: "East/West Playmaker"},
	{field: "PT_HBRECEIVINGBACK", label: "Backfield Threat"},
	{field: "PT_HBELUSIVEPOWER", label: "Elusive Bruiser"},
	{field: "PT_HBPOWERRECEIVING", label: "North/South Receiver"},
	{field: "PT_HBPOWERBLOCKING", label: "North/South Blocker"},

	// Wide receiver
	{field: "PT_WRPHYSICALRECEIVER", label: "Contested Specialist"},
	{field: "PT_WRPOWERBLOCKING", label: "Speedster"},
	{field: "PT_WRGADGET", label: "Gadget"},
	{field: "PT_WRPHYSICALROUTERUNNER", label: "Physical Route Runner"},
	{field: "PT_WRPLAYMAKER", label: "Route Artist"},
	{field: "PT_WRELUSIVEROUTERUNNER", label: "Elusive Route Runner"},
	{field: "PT_WRPHYSICALBLOCKER", label: "Gritty Possession"},

	// Tight end
	{field: "PT_TEVERTICALTHREAT", label: "Vertical Threat"},
	{field: "PT_TEPHYSICALROUTERUNNER", label: "Physical Route Runner"},
	{field: "PT_TEBLOCKING", label: "Pure Blocker"},
	{field: "PT_TEPHYSICALBLOCKER", label: "Gritty Possession"},

	// Defensive line
	{field: "PT_DLPUREPOWER", label: "Pure Power"},
	{field: "PT_DLPOWERRUSHER", label: "Power Rusher"},
	{field: "PT_DLSPEEDRUSHER", label: "Speed Rusher"},
	// Shared by DE_RunStopper (Edge Setter) and DT_NoseTackle (Gap Specialist);
	// resolved in archetypeLabelFromRecord via PlayerType/position.
	{field: "PT_DLRUNSTOPPER", label: "Gap Specialist"},

	// Linebacker
	{field: "PT_LBFIELDGENERAL", label: "Signal Caller"},
	{field: "PT_LBPASSCOVERAGE", label: "Lurker"},
	{field: "PT_LBRUNSUPPORT", label: "Thumper"},

	// Safety
	{field: "PT_SRUNSUPPORT", label: "Box Specialist"},
}

// playerTypeArchetypeLabels maps legacy PlayerType enum values to CFB UI labels
// when no PT_* trait flag is set (common for CB, OL, K/P, and some edge cases).
var playerTypeArchetypeLabels = map[string]string{
	// Quarterback (fallback when PT_* flags are unset)
	"QB_FieldGeneral":  "Pocket Passer",
	"QB_Improviser":    "Backfield Creator",
	"QB_Scrambler":     "Dual Threat",
	"QB_PureScrambler": "Pure Runner",

	// Halfback
	"HB_PowerBack":      "Contact Seeker",
	"HB_ElusiveBack":    "East/West Playmaker",
	"HB_ReceivingBack":  "Backfield Threat",
	"HB_ElusivePower":   "Elusive Bruiser",
	"HB_PowerReceiving": "North/South Receiver",
	"HB_PowerBlocking":  "North/South Blocker",

	// Wide receiver
	"WR_DeepThreat":          "Speedster",
	"WR_Playmaker":           "Route Artist",
	"WR_PhysicalRouteRunner": "Physical Route Runner",
	"WR_ShiftyRouteRunner":   "Elusive Route Runner",
	"WR_GadgetReceiver":      "Gadget",
	"WR_Physical":            "Contested Specialist",
	"WR_PhysicalBlocker":     "Gritty Possession",

	// Tight end
	"TE_VerticalThreat":      "Vertical Threat",
	"TE_PhysicalRouteRunner": "Physical Route Runner",
	"TE_Blocking":            "Pure Blocker",
	"TE_Possession":          "Pure Possession",

	// Cornerback
	"CB_Zone":         "Zone",
	"CB_MantoMan":     "Bump and Run",
	"CB_Slot":         "Boundary",
	"CB_HybridCorner": "Field",

	// Safety
	"S_Zone":       "Coverage Specialist",
	"S_Hybrid":     "Hybrid",
	"S_RunSupport": "Box Specialist",

	// Offensive line (shared naming across C/G/OT)
	"OT_PassProtector": "Pass Protector",
	"OT_Power":         "Raw Strength",
	"OT_Agile":         "Agile",
	"OT_WellRounded":   "Well-Rounded",
	"G_PassProtector":  "Pass Protector",
	"G_Power":          "Raw Strength",
	"G_Agile":          "Agile",
	"G_WellRounded":    "Well-Rounded",
	"C_PassProtector":  "Pass Protector",
	"C_Power":          "Raw Strength",
	"C_Agile":          "Agile",
	"C_WellRounded":    "Well-Rounded",

	// Defensive line fallbacks
	"DE_SmallerSpeedRusher": "Speed Rusher",
	"DE_PowerRusher":        "Power Rusher",
	"DE_PurePower":          "Pure Power",
	"DE_RunStopper":         "Edge Setter",
	"DT_NoseTackle":         "Gap Specialist",
	"DT_PurePower":          "Pure Power",
	"DT_SpeedRusher":        "Speed Rusher",
	"DT_PowerRusher":        "Power Rusher",
	"OLB_SpeedRusher":       "Speed Rusher",
	"OLB_PowerRusher":       "Power Rusher",
	"OLB_PassCoverage":      "Lurker",
	"OLB_RunStopper":        "Thumper",
	"MLB_FieldGeneral":      "Signal Caller",
	"MLB_PassCoverage":      "Lurker",
	"MLB_RunStopper":        "Thumper",

	// Fullback
	"FB_Blocking": "Blocking",
	"FB_Utility":  "Utility",

	// Kicker / punter / specialists
	"KP_Accurate": "Accurate",
	"KP_Power":    "Power",
	"KR_Balanced": "Balanced",
	"PR_Balanced": "Balanced",
	"LS_Power":    "Power",
	"LS_Accurate": "Accurate",
}

type archetypePTEntry struct {
	field string
	label string
}

func archetypeLabelFromRecord(record Record) string {
	for _, entry := range archetypePTEntries {
		value, ok := record.Get(entry.field)
		if ok && value.Bool {
			if entry.field == "PT_DLRUNSTOPPER" {
				return dlRunStopperLabel(record)
			}
			return entry.label
		}
	}
	if label, ok := playerTypeArchetypeLabels[stringField(record, "PlayerType")]; ok {
		return label
	}
	return ""
}

// dlRunStopperLabel disambiguates the shared PT_DLRUNSTOPPER flag. Edges
// (DE_RunStopper / LE / RE) display as "Edge Setter"; interior run-stoppers
// (DT_NoseTackle and other DTs) stay "Gap Specialist".
func dlRunStopperLabel(record Record) string {
	switch stringField(record, "PlayerType") {
	case "DE_RunStopper":
		return "Edge Setter"
	case "DT_NoseTackle":
		return "Gap Specialist"
	}
	switch stringField(record, "Position") {
	case "LE", "RE":
		return "Edge Setter"
	}
	return "Gap Specialist"
}
