package location

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type LocData struct {
	IANAName     string
	OriginalName string
	Population   int
}

var (
	// map of normalized name -> LocData
	index map[string]LocData
	//go:embed _data/cities.tsv
	citiesData []byte

	// offsetRegex matches common timezone shorthands and offsets
	offsetRegex = regexp.MustCompile(`^(?i)(?:UTC|GMT)?\s*([+-]\d{1,4}(?::?\d{2})?)$`)

	// tzAbbreviations maps common timezone abbreviations to their UTC offsets
	// in seconds, sourced from
	// https://en.wikipedia.org/wiki/List_of_time_zone_abbreviations.
	tzAbbreviations = map[string]int{
		"ACDT":  10*3600 + 30*60,
		"ACST":  9*3600 + 30*60,
		"ACT":   -5 * 3600,
		"ACWST": 8*3600 + 45*60,
		"ADT":   -3 * 3600,
		"AEDT":  11 * 3600,
		"AEST":  10 * 3600,
		"AET":   10 * 3600,
		"AFT":   4*3600 + 30*60,
		"AKDT":  -8 * 3600,
		"AKST":  -9 * 3600,
		"ALMT":  6 * 3600,
		"AMST":  -3 * 3600,
		"AMT":   -4 * 3600,
		"ANAT":  12 * 3600,
		"AQTT":  5 * 3600,
		"ART":   -3 * 3600,
		"AST":   -4 * 3600,
		"AWST":  8 * 3600,
		"AZOST": 0,
		"AZOT":  -1 * 3600,
		"AZT":   4 * 3600,
		"BNT":   8 * 3600,
		"BIOT":  6 * 3600,
		"BIT":   -12 * 3600,
		"BOT":   -4 * 3600,
		"BRST":  -2 * 3600,
		"BRT":   -3 * 3600,
		"BST":   1 * 3600,
		"BTT":   6 * 3600,
		"CAT":   2 * 3600,
		"CCT":   6*3600 + 30*60,
		"CDT":   -5 * 3600,
		"CEST":  2 * 3600,
		"CET":   1 * 3600,
		"CHADT": 13*3600 + 45*60,
		"CHAST": 12*3600 + 45*60,
		"CHOT":  8 * 3600,
		"CHOST": 9 * 3600,
		"CHST":  10 * 3600,
		"CHUT":  10 * 3600,
		"CIST":  -8 * 3600,
		"CKT":   -10 * 3600,
		"CLST":  -3 * 3600,
		"CLT":   -4 * 3600,
		"COST":  -4 * 3600,
		"COT":   -5 * 3600,
		"CST":   -6 * 3600,
		"CT":    -6 * 3600,
		"CVT":   -1 * 3600,
		"CWST":  8*3600 + 45*60,
		"CXT":   7 * 3600,
		"DAVT":  7 * 3600,
		"DDUT":  10 * 3600,
		"DFT":   1 * 3600,
		"EASST": -5 * 3600,
		"EAST":  -6 * 3600,
		"EAT":   3 * 3600,
		"ECT":   -5 * 3600,
		"EDT":   -4 * 3600,
		"EEST":  3 * 3600,
		"EET":   2 * 3600,
		"EGST":  0,
		"EGT":   -1 * 3600,
		"EST":   -5 * 3600,
		"ET":    -5 * 3600,
		"FET":   3 * 3600,
		"FJT":   12 * 3600,
		"FKST":  -3 * 3600,
		"FKT":   -4 * 3600,
		"FNT":   -2 * 3600,
		"GALT":  -6 * 3600,
		"GAMT":  -9 * 3600,
		"GET":   4 * 3600,
		"GFT":   -3 * 3600,
		"GILT":  12 * 3600,
		"GIT":   -9 * 3600,
		"GST":   4 * 3600,
		"GYT":   -4 * 3600,
		"HDT":   -9 * 3600,
		"HAEC":  2 * 3600,
		"HST":   -10 * 3600,
		"HKT":   8 * 3600,
		"HMT":   5 * 3600,
		"HOVST": 8 * 3600,
		"HOVT":  7 * 3600,
		"ICT":   7 * 3600,
		"IDLW":  -12 * 3600,
		"IDT":   3 * 3600,
		"IOT":   6 * 3600,
		"IRDT":  4*3600 + 30*60,
		"IRKT":  8 * 3600,
		"IRST":  3*3600 + 30*60,
		"IST":   5*3600 + 30*60,
		"JST":   9 * 3600,
		"KALT":  2 * 3600,
		"KGT":   6 * 3600,
		"KOST":  11 * 3600,
		"KRAT":  7 * 3600,
		"KST":   9 * 3600,
		"LHST":  10*3600 + 30*60,
		"LINT":  14 * 3600,
		"MAGT":  12 * 3600,
		"MART":  -(9*3600 + 30*60),
		"MAWT":  5 * 3600,
		"MDT":   -6 * 3600,
		"MET":   1 * 3600,
		"MEST":  2 * 3600,
		"MHT":   12 * 3600,
		"MIST":  11 * 3600,
		"MIT":   -(9*3600 + 30*60),
		"MMT":   6*3600 + 30*60,
		"MSK":   3 * 3600,
		"MST":   -7 * 3600,
		"MT":    -7 * 3600,
		"MUT":   4 * 3600,
		"MVT":   5 * 3600,
		"MYT":   8 * 3600,
		"NCT":   11 * 3600,
		"NDT":   -(2*3600 + 30*60),
		"NFT":   11 * 3600,
		"NOVT":  7 * 3600,
		"NPT":   5*3600 + 45*60,
		"NST":   -(3*3600 + 30*60),
		"NT":    -(3*3600 + 30*60),
		"NUT":   -11 * 3600,
		"NZDT":  13 * 3600,
		"NZDST": 13 * 3600,
		"NZST":  12 * 3600,
		"OMST":  6 * 3600,
		"ORAT":  5 * 3600,
		"PDT":   -7 * 3600,
		"PET":   -5 * 3600,
		"PETT":  12 * 3600,
		"PGT":   10 * 3600,
		"PHOT":  13 * 3600,
		"PHT":   8 * 3600,
		"PHST":  8 * 3600,
		"PKT":   5 * 3600,
		"PMDT":  -2 * 3600,
		"PMST":  -3 * 3600,
		"PONT":  11 * 3600,
		"PST":   -8 * 3600,
		"PT":    -8 * 3600,
		"PWT":   9 * 3600,
		"PYST":  -3 * 3600,
		"PYT":   -4 * 3600,
		"RET":   4 * 3600,
		"ROTT":  -3 * 3600,
		"SAKT":  11 * 3600,
		"SAMT":  4 * 3600,
		"SAST":  2 * 3600,
		"SBT":   11 * 3600,
		"SCT":   4 * 3600,
		"SDT":   -10 * 3600,
		"SGT":   8 * 3600,
		"SLST":  5*3600 + 30*60,
		"SRET":  11 * 3600,
		"SRT":   -3 * 3600,
		"SST":   -11 * 3600,
		"SYOT":  3 * 3600,
		"TAHT":  -10 * 3600,
		"THA":   7 * 3600,
		"TFT":   5 * 3600,
		"TJT":   5 * 3600,
		"TKT":   13 * 3600,
		"TLT":   9 * 3600,
		"TMT":   5 * 3600,
		"TRT":   3 * 3600,
		"TOT":   13 * 3600,
		"TST":   8 * 3600,
		"TVT":   12 * 3600,
		"ULAST": 9 * 3600,
		"ULAT":  8 * 3600,
		"UYST":  -2 * 3600,
		"UYT":   -3 * 3600,
		"UZT":   5 * 3600,
		"VET":   -4 * 3600,
		"VLAT":  10 * 3600,
		"VOLT":  3 * 3600,
		"VOST":  6 * 3600,
		"VUT":   11 * 3600,
		"WAKT":  12 * 3600,
		"WAST":  2 * 3600,
		"WAT":   1 * 3600,
		"WEST":  1 * 3600,
		"WET":   0,
		"WIB":   7 * 3600,
		"WIT":   9 * 3600,
		"WITA":  8 * 3600,
		"WGST":  -2 * 3600,
		"WGT":   -3 * 3600,
		"WST":   8 * 3600,
		"YAKT":  9 * 3600,
		"YEKT":  5 * 3600,
	}
)

func normalizeString(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, _ := transform.String(t, s)
	return strings.ToLower(result)
}

// init initializes the timezone index from the embedded TSV data.
func init() {
	index = map[string]LocData{}
	// read embedded TSV line by line
	lines := bytes.SplitSeq(citiesData, []byte{'\n'})
	for line := range lines {
		if len(line) == 0 {
			continue
		}
		// extract TSV fields
		parts := bytes.Split(line, []byte{'\t'})
		if len(parts) < 5 {
			continue
		}
		cityName := string(parts[0])
		countryCode := string(parts[1])
		stateCode := string(parts[2])
		population, _ := strconv.Atoi(string(parts[3]))
		ianaName := string(parts[4])
		// normalize city name for consistent searching
		normalized := normalizeString(cityName)
		// build searchable keys (city name only, city+country, city+state)
		keys := []string{
			normalized,
		}
		if countryCode != "" {
			keys = append(keys, normalized+","+strings.ToLower(countryCode))
		}
		if stateCode != "" {
			keys = append(keys, normalized+","+strings.ToLower(stateCode))
		}
		// populate the index, mapping each key to the IANA timezone data
		for _, k := range keys {
			existing, exists := index[k]
			if !exists || population > existing.Population {
				index[k] = LocData{
					IANAName:     ianaName,
					OriginalName: cityName,
					Population:   population,
				}
			}
		}
	}
}

// distance calculates the Levenshtein distance between two strings using a
// dynamic programming approach. It computes the minimum number of
// single-character edits (insertions, deletions, or substitutions)
// required to change string s1 into string s2.
//
// The algorithm constructs an (n+1) x (m+1) matrix where each cell (i, j)
// represents the cost to transform the prefix of s1 of length i into the
// prefix of s2 of length j.
//
// Example: transforming "cat" to "cars"
// Base matrix initialization (costs to transform to/from empty strings):
//
//	    ""  c  a  r  s
//	""   0  1  2  3  4
//	 c   1
//	 a   2
//	 t   3
//
// Matrix calculation step:
// For each cell (i, j), the value is the minimum of:
// 1. Deletion:       matrix[i-1][j] + 1
// 2. Insertion:      matrix[i][j-1] + 1
// 3. Substitution:   matrix[i-1][j-1] + cost (0 if chars match, 1 otherwise)
//
//	(j-1)      (j)
//
// (i-1) [Sub/Match] [Del]
// (i)   [Ins]       [Target] -> Target = min(Del+1, Ins+1, Sub+cost)
//
// Completed matrix for "cat" -> "cars":
//
//	    ""  c  a  r  s
//	""   0  1  2  3  4
//	 c   1  0  1  2  3
//	 a   2  1  0  1  2
//	 t   3  2  1  1  2  <- Final Distance = 2
func distance(s1, s2 string) int {
	// convert strings to rune slices to correctly handle multi-byte characters
	r1, r2 := []rune(s1), []rune(s2)
	n, m := len(r1), len(r2)
	// fast path: if either string is empty, the distance is the length of the
	// other string
	if n == 0 {
		return m
	}
	if m == 0 {
		return n
	}
	// initialize the matrix with base costs for transforming prefixes into empty
	// strings
	d := make([][]int, n+1)
	for i := range d {
		d[i] = make([]int, m+1)
		d[i][0] = i
	}
	for j := range d[0] {
		d[0][j] = j
	}
	// compute the Levenshtein matrix by finding the minimum of deletion,
	// insertion, and substitution costs
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			cost := 1
			if r1[i-1] == r2[j-1] {
				cost = 0
			}
			minVal := min(d[i-1][j-1]+cost, min(d[i][j-1]+1, d[i-1][j]+1))
			d[i][j] = minVal
		}
	}

	return d[n][m]
}

// find resolves a user-provided string to a *time.Location and its original name.
func Find(name string) (*time.Location, string, error) {
	if name == "" || name == "@" {
		return time.Local, "Local", nil
	}
	if loc, zoneName, ok := parseShorthand(name); ok {
		return loc, zoneName, nil
	}
	normalized := normalizeString(name)
	normalized = strings.ReplaceAll(normalized, ", ", ",")
	// exact match
	if data, ok := index[normalized]; ok {
		loc, err := time.LoadLocation(data.IANAName)
		return loc, data.OriginalName, err
	}
	// fuzzy match
	bestMatch := ""
	bestDist := 100
	for key := range index {
		dist := distance(normalized, key)
		if dist < bestDist {
			bestDist = dist
			bestMatch = key
		} else if dist == bestDist {
			// break ties by preferring the higher population
			if index[key].Population > index[bestMatch].Population {
				bestMatch = key
			} else if index[key].Population == index[bestMatch].Population {
				// break ties deterministically by choosing the alphabetically smaller key
				if key < bestMatch {
					bestMatch = key
				}
			}
		}
	}
	// Use a distance threshold to avoid returning completely unrelated
	// locations. It allows up to half the length of the string in edits (plus 1
	// to be lenient for very short strings).
	if bestMatch != "" && bestDist <= len(normalized)/2+1 {
		data := index[bestMatch]
		loc, err := time.LoadLocation(data.IANAName)
		return loc, data.OriginalName, err
	}
	return nil, "", errors.New("location not found")
}

// parseShorthand parses common timezone shorthands like UTC, UTC+8, +0800, -07:00
func parseShorthand(name string) (*time.Location, string, bool) {
	upper := strings.ToUpper(name)
	if upper == "UTC" || upper == "GMT" {
		return time.UTC, upper, true
	}
	if offsetSecs, ok := tzAbbreviations[upper]; ok {
		return time.FixedZone(upper, offsetSecs), upper, true
	}

	m := offsetRegex.FindStringSubmatch(name)
	if m == nil {
		return nil, "", false
	}

	offsetStr := m[1]
	sign := 1
	if offsetStr[0] == '-' {
		sign = -1
	}
	offsetStr = offsetStr[1:]
	offsetStr = strings.ReplaceAll(offsetStr, ":", "")

	var hours, mins int
	switch {
	case len(offsetStr) <= 2:
		hours, _ = strconv.Atoi(offsetStr)
	case len(offsetStr) == 3:
		hours, _ = strconv.Atoi(offsetStr[:1])
		mins, _ = strconv.Atoi(offsetStr[1:])
	case len(offsetStr) == 4:
		hours, _ = strconv.Atoi(offsetStr[:2])
		mins, _ = strconv.Atoi(offsetStr[2:])
	default:
		return nil, "", false
	}

	offsetSecs := sign * ((hours * 3600) + (mins * 60))
	zoneName := fmt.Sprintf("UTC%+03d:%02d", sign*hours, mins)
	if mins == 0 {
		zoneName = fmt.Sprintf("UTC%+d", sign*hours)
	}

	return time.FixedZone(zoneName, offsetSecs), zoneName, true
}
