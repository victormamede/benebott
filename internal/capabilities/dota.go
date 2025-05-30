package capabilities

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"google.golang.org/genai"
)

var DotaPlayerAccountDeclaration genai.FunctionDeclaration = genai.FunctionDeclaration{
	Name:        "dota_player_account",
	Description: "Gets dota player account information",
	Parameters: &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"playerId": &genai.Schema{Type: genai.TypeString},
		},
		Required: []string{"playerId"},
	},
}

type DotaHero struct {
	Name          string
	LocalizedName string
}

func DotaPlayerAccount(playerId string) CallResponse {
	fmt.Println("Getting dota account for player", playerId)

	response, err := http.Get(fmt.Sprintf("https://api.opendota.com/api/players/%s", playerId))
	callResponse := map[string]any{}

	if err != nil {
		callResponse["error"] = err.Error()
		return callResponse
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		callResponse["error"] = err.Error()
		return callResponse
	}

	err = json.Unmarshal(responseData, &callResponse)
	if err != nil {
		callResponse["error"] = err.Error()
		return callResponse
	}

	return callResponse
}

var DotaPlayerMatchesDeclaration genai.FunctionDeclaration = genai.FunctionDeclaration{
	Name:        "dota_player_matches",
	Description: "Gets match information for a dota player by id.",
	Parameters: &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"playerId": &genai.Schema{Type: genai.TypeString, Description: "The player ID"},
			"limit":    &genai.Schema{Type: genai.TypeInteger, Description: "The number of matches to fetch, not higher than 10"},
		},
		Required: []string{"playerId", "limit"},
	},
}

type DotaPlayerMatchResponse struct {
	MatchId      float64 `json:"match_id"`
	PlayerSlot   float64 `json:"player_slot"`
	RadiantWin   bool    `json:"radiant_win"`
	Duration     float64 `json:"duration"`
	HeroId       float64 `json:"hero_id"`
	StartTime    float64 `json:"start_time"`
	Kills        float64 `json:"kills"`
	Deaths       float64 `json:"deaths"`
	Assists      float64 `json:"assists"`
	Skill        float64 `json:"skill"`
	AverageRank  float64 `json:"average_rank"`
	LeaverStatus float64 `json:"leaver_status"`
	PartySize    float64 `json:"party_size"`
	HeroVariant  float64 `json:"hero_variant"`
}

func DotaPlayerMatches(playerId string, limit int) CallResponse {
	fmt.Println("Getting dota matches for player", playerId, "limit", limit)

	response, err := http.Get(fmt.Sprintf("https://api.opendota.com/api/players/%s/matches?limit=%d", playerId, limit))
	callResponse := map[string]any{}

	if err != nil {
		callResponse["error"] = err.Error()
		return callResponse
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		callResponse["error"] = err.Error()
		return callResponse
	}

	items := []DotaPlayerMatchResponse{}
	err = json.Unmarshal(responseData, &items)
	if err != nil {
		callResponse["error"] = err.Error()
		return callResponse
	}

	parsedItems := []any{}
	for _, item := range items {
		team := "radiant"
		won := item.RadiantWin
		if item.PlayerSlot >= 128 {
			team = "dire"
			won = !item.RadiantWin
		}
		startTime := time.Unix(int64(item.StartTime), 0)
		playerLeft := item.LeaverStatus != 0

		parsedItems = append(parsedItems, map[string]any{
			"match_id":         item.MatchId,
			"team":             team,
			"won":              won,
			"duration_minutes": item.Duration / 60.0,
			"hero":             heroes[uint8(item.HeroId)].Name,
			"start_time":       startTime.Format(time.RFC3339),
			"kills":            item.Kills,
			"deaths":           item.Deaths,
			"assists":          item.Assists,
			"did_player_leave": playerLeft,
			"party_size":       item.PartySize,
		})
	}

	callResponse["items"] = parsedItems
	return callResponse
}

var DotaHeroesDeclaration genai.FunctionDeclaration = genai.FunctionDeclaration{
	Name:        "dota_heroes",
	Description: `Gets information for each dota hero and their id. `,
	Parameters: &genai.Schema{
		Type:       genai.TypeObject,
		Properties: map[string]*genai.Schema{},
	},
}

func DotaHeroes() CallResponse {
	fmt.Println("Getting dota heroes")

	response, err := http.Get("https://api.opendota.com/api/heroes")
	callResponse := map[string]any{}

	if err != nil {
		callResponse["error"] = err.Error()
		return callResponse
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		callResponse["error"] = err.Error()
		return callResponse
	}

	items := []any{}
	err = json.Unmarshal(responseData, &items)
	if err != nil {
		callResponse["error"] = err.Error()
		return callResponse
	}

	callResponse["heroes"] = items
	return callResponse
}

var heroes map[uint8]DotaHero = map[uint8]DotaHero{
	1: DotaHero{
		Name:          "npc_dota_hero_antimage",
		LocalizedName: "Anti-Mage",
	},
	2: DotaHero{
		Name:          "npc_dota_hero_axe",
		LocalizedName: "Axe",
	},
	3: DotaHero{
		Name:          "npc_dota_hero_bane",
		LocalizedName: "Bane",
	},
	4: DotaHero{
		Name:          "npc_dota_hero_bloodseeker",
		LocalizedName: "Bloodseeker",
	},
	5: DotaHero{
		Name:          "npc_dota_hero_crystal_maiden",
		LocalizedName: "Crystal Maiden",
	},
	6: DotaHero{
		Name:          "npc_dota_hero_drow_ranger",
		LocalizedName: "Drow Ranger",
	},
	7: DotaHero{
		Name:          "npc_dota_hero_earthshaker",
		LocalizedName: "Earthshaker",
	},
	8: DotaHero{
		Name:          "npc_dota_hero_juggernaut",
		LocalizedName: "Juggernaut",
	},
	9: DotaHero{
		Name:          "npc_dota_hero_mirana",
		LocalizedName: "Mirana",
	},
	10: DotaHero{
		Name:          "npc_dota_hero_morphling",
		LocalizedName: "Morphling",
	},
	11: DotaHero{
		Name:          "npc_dota_hero_nevermore",
		LocalizedName: "Shadow Fiend",
	},
	12: DotaHero{
		Name:          "npc_dota_hero_phantom_lancer",
		LocalizedName: "Phantom Lancer",
	},
	13: DotaHero{
		Name:          "npc_dota_hero_puck",
		LocalizedName: "Puck",
	},
	14: DotaHero{
		Name:          "npc_dota_hero_pudge",
		LocalizedName: "Pudge",
	},
	15: DotaHero{
		Name:          "npc_dota_hero_razor",
		LocalizedName: "Razor",
	},
	16: DotaHero{
		Name:          "npc_dota_hero_sand_king",
		LocalizedName: "Sand King",
	},
	17: DotaHero{
		Name:          "npc_dota_hero_storm_spirit",
		LocalizedName: "Storm Spirit",
	},
	18: DotaHero{
		Name:          "npc_dota_hero_sven",
		LocalizedName: "Sven",
	},
	19: DotaHero{
		Name:          "npc_dota_hero_tiny",
		LocalizedName: "Tiny",
	},
	20: DotaHero{
		Name:          "npc_dota_hero_vengefulspirit",
		LocalizedName: "Vengeful Spirit",
	},
	21: DotaHero{
		Name:          "npc_dota_hero_windrunner",
		LocalizedName: "Windranger",
	},
	22: DotaHero{
		Name:          "npc_dota_hero_zuus",
		LocalizedName: "Zeus",
	},
	23: DotaHero{
		Name:          "npc_dota_hero_kunkka",
		LocalizedName: "Kunkka",
	},
	25: DotaHero{
		Name:          "npc_dota_hero_lina",
		LocalizedName: "Lina",
	},
	26: DotaHero{
		Name:          "npc_dota_hero_lion",
		LocalizedName: "Lion",
	},
	27: DotaHero{
		Name:          "npc_dota_hero_shadow_shaman",
		LocalizedName: "Shadow Shaman",
	},
	28: DotaHero{
		Name:          "npc_dota_hero_slardar",
		LocalizedName: "Slardar",
	},
	29: DotaHero{
		Name:          "npc_dota_hero_tidehunter",
		LocalizedName: "Tidehunter",
	},
	30: DotaHero{
		Name:          "npc_dota_hero_witch_doctor",
		LocalizedName: "Witch Doctor",
	},
	31: DotaHero{
		Name:          "npc_dota_hero_lich",
		LocalizedName: "Lich",
	},
	32: DotaHero{
		Name:          "npc_dota_hero_riki",
		LocalizedName: "Riki",
	},
	33: DotaHero{
		Name:          "npc_dota_hero_enigma",
		LocalizedName: "Enigma",
	},
	34: DotaHero{
		Name:          "npc_dota_hero_tinker",
		LocalizedName: "Tinker",
	},
	35: DotaHero{
		Name:          "npc_dota_hero_sniper",
		LocalizedName: "Sniper",
	},
	36: DotaHero{
		Name:          "npc_dota_hero_necrolyte",
		LocalizedName: "Necrophos",
	},
	37: DotaHero{
		Name:          "npc_dota_hero_warlock",
		LocalizedName: "Warlock",
	},
	38: DotaHero{
		Name:          "npc_dota_hero_beastmaster",
		LocalizedName: "Beastmaster",
	},
	39: DotaHero{
		Name:          "npc_dota_hero_queenofpain",
		LocalizedName: "Queen of Pain",
	},
	40: DotaHero{
		Name:          "npc_dota_hero_venomancer",
		LocalizedName: "Venomancer",
	},
	41: DotaHero{
		Name:          "npc_dota_hero_faceless_void",
		LocalizedName: "Faceless Void",
	},
	42: DotaHero{
		Name:          "npc_dota_hero_skeleton_king",
		LocalizedName: "Wraith King",
	},
	43: DotaHero{
		Name:          "npc_dota_hero_death_prophet",
		LocalizedName: "Death Prophet",
	},
	44: DotaHero{
		Name:          "npc_dota_hero_phantom_assassin",
		LocalizedName: "Phantom Assassin",
	},
	45: DotaHero{
		Name:          "npc_dota_hero_pugna",
		LocalizedName: "Pugna",
	},
	46: DotaHero{
		Name:          "npc_dota_hero_templar_assassin",
		LocalizedName: "Templar Assassin",
	},
	47: DotaHero{
		Name:          "npc_dota_hero_viper",
		LocalizedName: "Viper",
	},
	48: DotaHero{
		Name:          "npc_dota_hero_luna",
		LocalizedName: "Luna",
	},
	49: DotaHero{
		Name:          "npc_dota_hero_dragon_knight",
		LocalizedName: "Dragon Knight",
	},
	50: DotaHero{
		Name:          "npc_dota_hero_dazzle",
		LocalizedName: "Dazzle",
	},
	51: DotaHero{
		Name:          "npc_dota_hero_rattletrap",
		LocalizedName: "Clockwerk",
	},
	52: DotaHero{
		Name:          "npc_dota_hero_leshrac",
		LocalizedName: "Leshrac",
	},
	53: DotaHero{
		Name:          "npc_dota_hero_furion",
		LocalizedName: "Nature's Prophet",
	},
	54: DotaHero{
		Name:          "npc_dota_hero_life_stealer",
		LocalizedName: "Lifestealer",
	},
	55: DotaHero{
		Name:          "npc_dota_hero_dark_seer",
		LocalizedName: "Dark Seer",
	},
	56: DotaHero{
		Name:          "npc_dota_hero_clinkz",
		LocalizedName: "Clinkz",
	},
	57: DotaHero{
		Name:          "npc_dota_hero_omniknight",
		LocalizedName: "Omniknight",
	},
	58: DotaHero{
		Name:          "npc_dota_hero_enchantress",
		LocalizedName: "Enchantress",
	},
	59: DotaHero{
		Name:          "npc_dota_hero_huskar",
		LocalizedName: "Huskar",
	},
	60: DotaHero{
		Name:          "npc_dota_hero_night_stalker",
		LocalizedName: "Night Stalker",
	},
	61: DotaHero{
		Name:          "npc_dota_hero_broodmother",
		LocalizedName: "Broodmother",
	},
	62: DotaHero{
		Name:          "npc_dota_hero_bounty_hunter",
		LocalizedName: "Bounty Hunter",
	},
	63: DotaHero{
		Name:          "npc_dota_hero_weaver",
		LocalizedName: "Weaver",
	},
	64: DotaHero{
		Name:          "npc_dota_hero_jakiro",
		LocalizedName: "Jakiro",
	},
	65: DotaHero{
		Name:          "npc_dota_hero_batrider",
		LocalizedName: "Batrider",
	},
	66: DotaHero{
		Name:          "npc_dota_hero_chen",
		LocalizedName: "Chen",
	},
	67: DotaHero{
		Name:          "npc_dota_hero_spectre",
		LocalizedName: "Spectre",
	},
	68: DotaHero{
		Name:          "npc_dota_hero_ancient_apparition",
		LocalizedName: "Ancient Apparition",
	},
	69: DotaHero{
		Name:          "npc_dota_hero_doom_bringer",
		LocalizedName: "Doom",
	},
	70: DotaHero{
		Name:          "npc_dota_hero_ursa",
		LocalizedName: "Ursa",
	},
	71: DotaHero{
		Name:          "npc_dota_hero_spirit_breaker",
		LocalizedName: "Spirit Breaker",
	},
	72: DotaHero{
		Name:          "npc_dota_hero_gyrocopter",
		LocalizedName: "Gyrocopter",
	},
	73: DotaHero{
		Name:          "npc_dota_hero_alchemist",
		LocalizedName: "Alchemist",
	},
	74: DotaHero{
		Name:          "npc_dota_hero_invoker",
		LocalizedName: "Invoker",
	},
	75: DotaHero{
		Name:          "npc_dota_hero_silencer",
		LocalizedName: "Silencer",
	},
	76: DotaHero{
		Name:          "npc_dota_hero_obsidian_destroyer",
		LocalizedName: "Outworld Destroyer",
	},
	77: DotaHero{
		Name:          "npc_dota_hero_lycan",
		LocalizedName: "Lycan",
	},
	78: DotaHero{
		Name:          "npc_dota_hero_brewmaster",
		LocalizedName: "Brewmaster",
	},
	79: DotaHero{
		Name:          "npc_dota_hero_shadow_demon",
		LocalizedName: "Shadow Demon",
	},
	80: DotaHero{
		Name:          "npc_dota_hero_lone_druid",
		LocalizedName: "Lone Druid",
	},
	81: DotaHero{
		Name:          "npc_dota_hero_chaos_knight",
		LocalizedName: "Chaos Knight",
	},
	82: DotaHero{
		Name:          "npc_dota_hero_meepo",
		LocalizedName: "Meepo",
	},
	83: DotaHero{
		Name:          "npc_dota_hero_treant",
		LocalizedName: "Treant Protector",
	},
	84: DotaHero{
		Name:          "npc_dota_hero_ogre_magi",
		LocalizedName: "Ogre Magi",
	},
	85: DotaHero{
		Name:          "npc_dota_hero_undying",
		LocalizedName: "Undying",
	},
	86: DotaHero{
		Name:          "npc_dota_hero_rubick",
		LocalizedName: "Rubick",
	},
	87: DotaHero{
		Name:          "npc_dota_hero_disruptor",
		LocalizedName: "Disruptor",
	},
	88: DotaHero{
		Name:          "npc_dota_hero_nyx_assassin",
		LocalizedName: "Nyx Assassin",
	},
	89: DotaHero{
		Name:          "npc_dota_hero_naga_siren",
		LocalizedName: "Naga Siren",
	},
	90: DotaHero{
		Name:          "npc_dota_hero_keeper_of_the_light",
		LocalizedName: "Keeper of the Light",
	},
	91: DotaHero{
		Name:          "npc_dota_hero_wisp",
		LocalizedName: "Io",
	},
	92: DotaHero{
		Name:          "npc_dota_hero_visage",
		LocalizedName: "Visage",
	},
	93: DotaHero{
		Name:          "npc_dota_hero_slark",
		LocalizedName: "Slark",
	},
	94: DotaHero{
		Name:          "npc_dota_hero_medusa",
		LocalizedName: "Medusa",
	},
	95: DotaHero{
		Name:          "npc_dota_hero_troll_warlord",
		LocalizedName: "Troll Warlord",
	},
	96: DotaHero{
		Name:          "npc_dota_hero_centaur",
		LocalizedName: "Centaur Warrunner",
	},
	97: DotaHero{
		Name:          "npc_dota_hero_magnataur",
		LocalizedName: "Magnus",
	},
	98: DotaHero{
		Name:          "npc_dota_hero_shredder",
		LocalizedName: "Timbersaw",
	},
	99: DotaHero{
		Name:          "npc_dota_hero_bristleback",
		LocalizedName: "Bristleback",
	},
	100: DotaHero{
		Name:          "npc_dota_hero_tusk",
		LocalizedName: "Tusk",
	},
	101: DotaHero{
		Name:          "npc_dota_hero_skywrath_mage",
		LocalizedName: "Skywrath Mage",
	},
	102: DotaHero{
		Name:          "npc_dota_hero_abaddon",
		LocalizedName: "Abaddon",
	},
	103: DotaHero{
		Name:          "npc_dota_hero_elder_titan",
		LocalizedName: "Elder Titan",
	},
	104: DotaHero{
		Name:          "npc_dota_hero_legion_commander",
		LocalizedName: "Legion Commander",
	},
	105: DotaHero{
		Name:          "npc_dota_hero_techies",
		LocalizedName: "Techies",
	},
	106: DotaHero{
		Name:          "npc_dota_hero_ember_spirit",
		LocalizedName: "Ember Spirit",
	},
	107: DotaHero{
		Name:          "npc_dota_hero_earth_spirit",
		LocalizedName: "Earth Spirit",
	},
	108: DotaHero{
		Name:          "npc_dota_hero_abyssal_underlord",
		LocalizedName: "Underlord",
	},
	109: DotaHero{
		Name:          "npc_dota_hero_terrorblade",
		LocalizedName: "Terrorblade",
	},
	110: DotaHero{
		Name:          "npc_dota_hero_phoenix",
		LocalizedName: "Phoenix",
	},
	111: DotaHero{
		Name:          "npc_dota_hero_oracle",
		LocalizedName: "Oracle",
	},
	112: DotaHero{
		Name:          "npc_dota_hero_winter_wyvern",
		LocalizedName: "Winter Wyvern",
	},
	113: DotaHero{
		Name:          "npc_dota_hero_arc_warden",
		LocalizedName: "Arc Warden",
	},
	114: DotaHero{
		Name:          "npc_dota_hero_monkey_king",
		LocalizedName: "Monkey King",
	},
	119: DotaHero{
		Name:          "npc_dota_hero_dark_willow",
		LocalizedName: "Dark Willow",
	},
	120: DotaHero{
		Name:          "npc_dota_hero_pangolier",
		LocalizedName: "Pangolier",
	},
	121: DotaHero{
		Name:          "npc_dota_hero_grimstroke",
		LocalizedName: "Grimstroke",
	},
	123: DotaHero{
		Name:          "npc_dota_hero_hoodwink",
		LocalizedName: "Hoodwink",
	},
	126: DotaHero{
		Name:          "npc_dota_hero_void_spirit",
		LocalizedName: "Void Spirit",
	},
	128: DotaHero{
		Name:          "npc_dota_hero_snapfire",
		LocalizedName: "Snapfire",
	},
	129: DotaHero{
		Name:          "npc_dota_hero_mars",
		LocalizedName: "Mars",
	},
	131: DotaHero{
		Name:          "npc_dota_hero_ringmaster",
		LocalizedName: "Ringmaster",
	},
	135: DotaHero{
		Name:          "npc_dota_hero_dawnbreaker",
		LocalizedName: "Dawnbreaker",
	},
	136: DotaHero{
		Name:          "npc_dota_hero_marci",
		LocalizedName: "Marci",
	},
	137: DotaHero{
		Name:          "npc_dota_hero_primal_beast",
		LocalizedName: "Primal Beast",
	},
	138: DotaHero{
		Name:          "npc_dota_hero_muerta",
		LocalizedName: "Muerta",
	},
	145: DotaHero{
		Name:          "npc_dota_hero_kez",
		LocalizedName: "Kez",
	},
}
