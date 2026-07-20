package exporter

// Mental ability icon assets from https://collegefootball.gg/abilities/ (mental section).
var mentalAbilityAssets = map[string]struct {
	Slug  string
	Ext   string
	Label string
}{
	"Adrenaline":      {Slug: "adrenaline", Ext: "png", Label: "Adrenaline"},
	"BestFriend":      {Slug: "best-friend", Ext: "png", Label: "Best Friend"},
	"ClearHeaded":     {Slug: "clear-headed", Ext: "png", Label: "Clear Headed"},
	"ClutchKicker":    {Slug: "clutch-kicker", Ext: "png", Label: "Clutch Kicker"},
	"DefensiveRally":  {Slug: "defensive-rally", Ext: "png", Label: "Defensive Rally"},
	"DBRally":         {Slug: "legion", Ext: "png", Label: "Legion"},
	"DLRally":         {Slug: "defensive-rally", Ext: "png", Label: "Defensive Rally"},
	"FanFavorite":     {Slug: "fan-favorite", Ext: "png", Label: "Fan Favorite"},
	"HomeFanFavorite": {Slug: "fan-favorite", Ext: "png", Label: "Fan Favorite"},
	"FieldGeneral":    {Slug: "field-general", Ext: "png", Label: "Field General"},
	"Headstrong":      {Slug: "headstrong", Ext: "png", Label: "Headstrong"},
	"HotHead":         {Slug: "hot-head", Ext: "webp", Label: "Hot Head"},
	"Instinct":        {Slug: "instinct", Ext: "png", Label: "Instinct"},
	"Legion":          {Slug: "legion", Ext: "png", Label: "Legion"},
	"OffensiveRally":  {Slug: "offensive-rally", Ext: "png", Label: "Offensive Rally"},
	"OLRally":         {Slug: "offensive-rally", Ext: "png", Label: "Offensive Rally"},
	"RoadDog":         {Slug: "road-dog", Ext: "png", Label: "Road Dog"},
	"RoadFanFavorite": {Slug: "road-dog", Ext: "png", Label: "Road Dog"},
	"TeamPlayer":      {Slug: "team-player", Ext: "png", Label: "Team Player"},
	"TheNatural":      {Slug: "the-natural", Ext: "png", Label: "The Natural"},
	"WinningTime":     {Slug: "winning-time", Ext: "png", Label: "Winning Time"},
}
