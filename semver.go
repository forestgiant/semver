package semver

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	numerals     string = "0123456789"
	alphabet            = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	alphaNumeric        = alphabet + numerals
)

// Version struct represents a semantic version
type Version struct {
	Major uint64
	Minor uint64
	Patch uint64
}

func (v *Version) String() string {
	b := make([]byte, 0, 5)
	b = strconv.AppendUint(b, v.Major, 10)
	b = append(b, '.')
	b = strconv.AppendUint(b, v.Minor, 10)
	b = append(b, '.')
	b = strconv.AppendUint(b, v.Patch, 10)

	return string(b)
}

// Equal accepts to semver strings and compares them
func (v *Version) Equal(v2 *Version) bool {
	if v.Major != v2.Major {
		return false
	}

	if v.Minor != v2.Minor {
		return false
	}
	if v.Patch != v2.Patch {
		return false
	}

	return true
}

// CreateFlag takes a semver Version struct and sets
// a -version flag to return a semver string
func (v *Version) CreateFlag() *bool {
	var (
		versionUsage = fmt.Sprintf("Prints current version: v. %v", v)
		versionPtr   = flag.Bool("version", false, versionUsage)
	)

	// Set up short hand flags
	flag.BoolVar(versionPtr, "v", false, versionUsage+" (shorthand)")

	return versionPtr
}

// CreateFlagAndParse takes a semver string and creates a version flag
// It will parse all flags (flag.Parse())
func CreateFlagAndParse(s string) error {
	// Error if flags are already parsed
	if flag.Parsed() {
		return errors.New("Flags have been parsed")
	}

	// Create Version struct for supplied strings
	v, err := NewVersion(s)
	if err != nil {
		return err
	}

	// Create version flag
	versionPtr := v.CreateFlag()

	// Parse all flags
	flag.Parse()

	if *versionPtr {
		fmt.Println(v)
		os.Exit(0)
	}

	return nil
}

// Equal accepts to semver strings and compares them
func Equal(s1 string, s2 string) bool {
	// Create Version struct for supplied strings
	v, err := NewVersion(s1)
	if err != nil {
		return false
	}

	v2, err := NewVersion(s2)
	if err != nil {
		return false
	}

	return v.Equal(v2)
}

// NewVersion parses a version string
func NewVersion(s string) (*Version, error) {
	if len(s) == 0 {
		return nil, errors.New("Version string empty")
	}

	// Split into major.minor.patch
	parts := strings.SplitN(s, ".", 3)
	if len(parts) != 3 {
		return nil, errors.New("Major.Minor.Patch elements not found")
	}

	// Major
	if !containsOnly(parts[0], numerals) {
		return nil, fmt.Errorf("Invalid character(s) found in major number %q", parts[0])
	}
	if hasLeadingZeroes(parts[0]) {
		return nil, fmt.Errorf("Major number must not contain leading zeroes %q", parts[0])
	}
	major, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return nil, err
	}

	// Minor
	if !containsOnly(parts[1], numerals) {
		return nil, fmt.Errorf("Invalid character(s) found in minor number %q", parts[1])
	}
	if hasLeadingZeroes(parts[1]) {
		return nil, fmt.Errorf("Minor number must not contain leading zeroes %q", parts[1])
	}
	minor, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return nil, err
	}

	// Patch
	if !containsOnly(parts[2], numerals) {
		return nil, fmt.Errorf("Invalid character(s) found in patch number %q", parts[2])
	}
	if hasLeadingZeroes(parts[2]) {
		return nil, fmt.Errorf("Patch number must not contain leading zeroes %q", parts[2])
	}
	patch, err := strconv.ParseUint(parts[2], 10, 64)
	if err != nil {
		return nil, err
	}

	v := new(Version)
	v.Major = major
	v.Minor = minor
	v.Patch = patch

	return v, nil
}

func containsOnly(s string, compare string) bool {
	return strings.IndexFunc(s, func(r rune) bool {
		return !strings.ContainsRune(compare, r)
	}) == -1
}

func hasLeadingZeroes(s string) bool {
	return len(s) > 1 && s[0] == '0'
}
