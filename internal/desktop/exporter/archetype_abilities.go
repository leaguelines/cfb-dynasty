package exporter

// Archetype physical ability slot order from https://collegefootball.gg/abilities/
// Kept in sync with github.com/leaguelines/cfb-dynasty/dynasty/archetype_physical_abilities.go.
var archetypePhysicalAbilities = map[string][5]string{
	"Agile":                  {"Screen Enforcer", "Quick Step", "Option Shield", "Outside Shield", "Quick Drop"},
	"Backfield Creator":      {"Off Platform", "Pull Down", "On Time", "Magician", "Mobile Deadeye"},
	"Backfield Threat":       {"360", "Safety Valve", "Takeoff", "Side Step", "Recoup"},
	"Blocking":               {"Strong Grip", "Second Level", "Pocket Shield", "Sidekick", "Screen Enforcer"},
	"Boundary":               {"Jammer", "Blanket Coverage", "Lay Out", "Wrap Up", "Quick Jump"},
	"Box Specialist":         {"Aftershock", "Wrap Up", "Hammer", "Blow Up", "Workhorse"},
	"Bump and Run":           {"Blanket Coverage", "Jammer", "House Call", "Ball Hawk", "Knockout"},
	"Contact Seeker":         {"Downhill", "Workhorse", "Battering Ram", "Ball Security", "Balanced"},
	"Contested Specialist":   {"50/50", "Workhorse", "Balanced", "Headfirst", "Downhill"},
	"Coverage Specialist":    {"Ball Hawk", "Lay Out", "House Call", "Robber", "Knockout"},
	"Dual Threat":            {"Downhill", "Extender", "Option King", "Dot!", "Mobile Resistance"},
	"East/West Playmaker":    {"Recoup", "Shifty", "Side Step", "360", "Arm Bar"},
	"Edge Setter":            {"Grip Breaker", "Inside Disruptor", "Outside Disruptor", "Option Disruptor", "Workhorse"},
	"Elusive Bruiser":        {"Shifty", "Headfirst", "Side Step", "Downhill", "Arm Bar"},
	"Elusive Route Runner":   {"360", "Cutter", "Double Dip", "Recoup", "Side Step"},
	"Field":                  {"Wrap Up", "Robber", "Knockout", "Blanket Coverage", "Ball Hawk"},
	"Gadget":                 {"Side Step", "Shifty", "Dot!", "Cutter", "Extender"},
	"Gap Specialist":         {"Grip Breaker", "Inside Disruptor", "Outside Disruptor", "Option Disruptor", "Workhorse"},
	"Hybrid":                 {"Wrap Up", "Hammer", "Knockout", "Aftershock", "Blow Up"},
	"KP|Accurate":            {"Chip Shot", "Deep Range", "Mega Leg", "", ""},
	"KP|Power":               {"Deep Range", "Mega Leg", "Coffin Corner", "", ""},
	"Lurker":                 {"House Call", "Knockout", "Bouncer", "Robber", "Wrap Up"},
	"North/South Blocker":    {"Headfirst", "Balanced", "Sidekick", "Ball Security", "Strong Grip"},
	"North/South Receiver":   {"Balanced", "Arm Bar", "Safety Valve", "Headfirst", "Downhill"},
	"Pass Protector":         {"Pocket Shield", "Quick Drop", "PA Shield", "Strong Grip", "Wear Down"},
	"Pocket Passer":          {"Resistance", "Step Up", "Sleight Of Hand", "Dot!", "On Time"},
	"Power Rusher":           {"Pocket Disruptor", "Duress", "Grip Breaker", "Workhorse", "Take Down"},
	"Pure Blocker":           {"Strong Grip", "Quick Drop", "Outside Shield", "Pocket Shield", "Second Level"},
	"Pure Possession":        {"Sure Hands", "Wear Down", "Strong Grip", "Outside Shield", "Balanced"},
	"Pure Power":             {"Grip Breaker", "Pocket Disruptor", "Inside Disruptor", "Workhorse", "Hammer"},
	"Pure Runner":            {"Downhill", "Option King", "Shifty", "Side Step", "Workhorse"},
	"Raw Strength":           {"Strong Grip", "Workhorse", "Second Level", "Inside Shield", "Ground N Pound"},
	"Route Artist":           {"Cutter", "Lay Out", "Recoup", "Double Dip", "Sure Hands"},
	"Signal Caller":          {"Take Down", "Workhorse", "Blow Up", "Wrap Up", "Hammer"},
	"Speed Rusher":           {"Quick Jump", "Duress", "Take Down", "Pocket Disruptor", "Recoup"},
	"Speedster":              {"Side Step", "Double Dip", "Take Off", "Recoup", "Shifty"},
	"TE|Gritty Possession":   {"Workhorse", "Strong Grip", "Sure Hands", "Outside Shield", "Battering Ram"},
	"TE|Physical Route Runner": {"Balanced", "50/50", "Cutter", "Downhill", "Sure Hands"},
	"Thumper":                {"Grip Breaker", "Wrap Up", "Aftershock", "Blow Up", "Hammer"},
	"Utility":                {"Safety Valve", "Balanced", "Screen Enforcer", "Sidekick", "Recoup"},
	"Vertical Threat":        {"Workhorse", "Balanced", "Take Off", "Recoup", "50/50"},
	"WR|Gritty Possession":   {"Second Level", "Outside Shield", "Strong Grip", "Workhorse", "Sure Hands"},
	"WR|Physical Route Runner": {"Downhill", "Press Pro", "Sure Hands", "50/50", "Cutter"},
	"Well-Rounded":           {"Pocket Shield", "Outside Shield", "Strong Grip", "Option Shield", "Inside Shield"},
	"Zone":                   {"Knockout", "Lay Out", "House Call", "Ball Hawk", "Bouncer"},
}

func physicalAbilityPositionKey(position string) string {
	switch position {
	case "HB", "RB":
		return "HB"
	case "LT", "LG", "C", "RG", "RT", "OL":
		return "OL"
	case "LE", "RE", "DT", "DL":
		return "DL"
	case "SAM", "MIKE", "WILL", "LB":
		return "LB"
	case "FS", "SS", "S":
		return "S"
	case "K", "P", "PK", "KP":
		return "KP"
	default:
		return position
	}
}

// ArchetypePhysicalAbilityName returns the ability name for an archetype slot (1–5).
func ArchetypePhysicalAbilityName(position, archetypeLabel string, slot int) string {
	if archetypeLabel == "" || slot < 1 || slot > 5 {
		return ""
	}
	if posKey := physicalAbilityPositionKey(position); posKey != "" {
		if abilities, ok := archetypePhysicalAbilities[posKey+"|"+archetypeLabel]; ok {
			return abilities[slot-1]
		}
	}
	if abilities, ok := archetypePhysicalAbilities[archetypeLabel]; ok {
		return abilities[slot-1]
	}
	return ""
}
