package gitprompt

import (
	"testing"
)

var all = &GitStatus{
	Branch:    "master",
	Sha:       "0455b83f923a40f0b485665c44aa068bc25029f5",
	Untracked: 0,
	Modified:  1,
	Staged:    2,
	Conflicts: 3,
	Ahead:     4,
	Behind:    5,
}

func TestPrinterEmpty(t *testing.T) {
	actual := Print(nil, "%h")
	assertOutput(t, "", actual)
}

func TestPrinterData(t *testing.T) {
	actual := Print(all, "%h %u %m %s %c %a %b")
	assertOutput(t, "master 0 1 2 3 4 5", actual)
}

func TestPrinterUnicode(t *testing.T) {
	actual := Print(all, "%h ✋%u ⚡️%m 🚚%s ❗️%c ⬆%a ⬇%b")
	assertOutput(t, "master ✋0 ⚡️1 🚚2 ❗️3 ⬆4 ⬇5", actual)
}

func TestShortSHA(t *testing.T) {
	actual := Print(&GitStatus{Sha: "858828b5e153f24644bc867598298b50f8223f9b"}, "%h")
	assertOutput(t, "858828b", actual)
}

func TestPrinterColorAttributes(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		expected string
	}{
		{
			name:     "red",
			format:   "#r%h",
			expected: "\x1b[31mmaster\x1b[0m",
		},
		{
			name:     "bold",
			format:   "@b%h",
			expected: "\x1b[1mmaster\x1b[0m",
		},
		{
			name:     "color & attribute",
			format:   "#r@bA",
			expected: "\x1b[1;31mA\x1b[0m",
		},
		{
			name:     "color & attribute reversed",
			format:   "@b#rA",
			expected: "\x1b[1;31mA\x1b[0m",
		},
		{
			name:     "ignore format until non-whitespace",
			format:   "A#r#g#b     B@i\tC",
			expected: "A     \x1b[34mB\t\x1b[3mC\x1b[0m",
		},
		{
			name:     "reset color",
			format:   "#rA#_B",
			expected: "\x1b[31mA\x1b[0mB",
		},
		{
			name:     "reset attributes",
			format:   "@bA@_B",
			expected: "\x1b[1mA\x1b[0mB",
		},
		{
			name:     "reset attribute",
			format:   "#ggreen @b@igreen_bold_italic @Bgreen_italic",
			expected: "\x1b[32mgreen \x1b[1;3mgreen_bold_italic \x1b[0;3;32mgreen_italic\x1b[0m",
		},
		{
			name:     "ending with #",
			format:   "%h#",
			expected: "master#",
		},
		{
			name:     "ending with !",
			format:   "%h!",
			expected: "master!",
		},
		{
			name:     "ending with @",
			format:   "%h@",
			expected: "master@",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := Print(all, test.format)
			assertOutput(t, test.expected, actual)
		})
	}
}

func TestPrinterGroups(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		expected string
	}{
		{
			name:     "groups",
			format:   "<[%h][ B%b A%a][ U%u][ C%c]>",
			expected: "<master B5 A4 C3>",
		},
		{
			name:     "group color",
			format:   "<[#r%h]-[#g%u]%a[-#b%b]>",
			expected: "<\x1b[31mmaster\x1b[0m-4-\x1b[34m5\x1b[0m>",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := Print(all, test.format)
			assertOutput(t, test.expected, actual)
		})
	}
}

func TestPrinterNonMatching(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		expected string
	}{
		{
			name:     "data valid odd",
			format:   "%%%h",
			expected: "%%master",
		},
		{
			name:     "data valid even",
			format:   "%%%%h",
			expected: "%%%%h",
		},
		{
			name:     "data invalid odd",
			format:   "%%%z",
			expected: "%%%z",
		},
		{
			name:     "data invalid even",
			format:   "%%%%z",
			expected: "%%%%z",
		},
		{
			name:     "color valid odd",
			format:   "###rA",
			expected: "##\x1b[31mA\x1b[0m",
		},
		{
			name:     "color valid even",
			format:   "####rA",
			expected: "####rA",
		},
		{
			name:     "color invalid odd",
			format:   "###zA",
			expected: "###zA",
		},
		{
			name:     "color invalid even",
			format:   "####zA",
			expected: "####zA",
		},
		{
			name:     "attribute valid odd",
			format:   "@@@bA",
			expected: "@@\x1b[1mA\x1b[0m",
		},
		{
			name:     "attribute valid even",
			format:   "@@@@bA",
			expected: "@@@@bA",
		},
		{
			name:     "attribute invalid odd",
			format:   "@@@zA",
			expected: "@@@zA",
		},
		{
			name:     "attribute invalid even",
			format:   "@@@@zA",
			expected: "@@@@zA",
		},
		{
			name:     "trailing %",
			format:   "A%",
			expected: "A%",
		},
		{
			name:     "trailing #",
			format:   "A#",
			expected: "A#",
		},
		{
			name:     "trailing @",
			format:   "A@",
			expected: "A@",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := Print(all, test.format)
			assertOutput(t, test.expected, actual)
		})
	}
}

func TestPrinterEscape(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		expected string
	}{
		{
			name:     "data",
			format:   "A\\%h",
			expected: "A%h",
		},
		{
			name:     "color",
			format:   "A\\#rB",
			expected: "A#rB",
		},
		{
			name:     "attribute",
			format:   "A\\!bB",
			expected: "A!bB",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := Print(all, test.format)
			assertOutput(t, test.expected, actual)
		})
	}
}

func assertOutput(t *testing.T, expected, actual string) {
	if actual == expected {
		return
	}
	actualEscaped := actual + "\x1b[0m"
	t.Errorf(`
Expected:    %s
            %q
Actual:      %s
            %q`,
		expected,
		expected,
		actualEscaped,
		actual,
	)
}
