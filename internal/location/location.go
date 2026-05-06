package location

import (
	"bytes"
	_ "embed"
	"errors"
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
