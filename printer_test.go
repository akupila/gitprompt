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
	actual, w := Print(nil, "%h")
	assertOutput(t, "", actual)
	assertWidth(t, 0, w)
}

func TestPrinterData(t *testing.T) {
	actual, w := Print(all, "%h %u %m %s %c %a %b")
	assertOutput(t, "master 0 1 2 3 4 5", actual)
	assertWidth(t, 18, w)
}

func TestPrinterUnicode(t *testing.T) {
	actual, w := Print(all, "%h ✋%u ⚡️%m 🚚%s ❗️%c ⬆%a ⬇%b")
	assertOutput(t, "master ✋0 ⚡️1 🚚2 ❗️3 ⬆4 ⬇5", actual)
	assertWidth(t, 26, w)
}

func TestShortSHA(t *testing.T) {
	actual, w := Print(&GitStatus{Sha: "858828b5e153f24644bc867598298b50f8223f9b"}, "%h")
	assertOutput(t, "858828b", actual)
	assertWidth(t, 7, w)
}

func TestPrinterColorAttributes(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		expected string
		width    int
	}{
		{
			name:     "red",
			format:   "#r%h",
			expected: "\x1b[31mmaster\x1b[0m",
			width:    6,
		},
		{
			name:     "bold",
			format:   "@b%h",
			expected: "\x1b[1mmaster\x1b[0m",
			width:    6,
		},
		{
			name:     "color & attribute",
			format:   "#r@bA",
			expected: "\x1b[1;31mA\x1b[0m",
			width:    1,
		},
		{
			name:     "color & attribute reversed",
			format:   "@b#rA",
			expected: "\x1b[1;31mA\x1b[0m",
			width:    1,
		},
		{
			name:     "ignore format until non-whitespace",
			format:   "A#r#g#b     B@i\tC",
			expected: "A     \x1b[34mB\t\x1b[3mC\x1b[0m",
			width:    9,
		},
		{
			name:     "reset color",
			format:   "#rA#_B",
			expected: "\x1b[31mA\x1b[0mB",
			width:    2,
		},
		{
			name:     "reset attributes",
			format:   "@bA@_B",
			expected: "\x1b[1mA\x1b[0mB",
			width:    2,
		},
		{
			name:     "reset attribute",
			format:   "#ggreen @b@igreen_bold_italic @Bgreen_italic",
			expected: "\x1b[32mgreen \x1b[1;3mgreen_bold_italic \x1b[0;3;32mgreen_italic\x1b[0m",
			width:    36,
		},
		{
			name:     "ending with #",
			format:   "%h#",
			expected: "master#",
			width:    7,
		},
		{
			name:     "ending with !",
			format:   "%h!",
			expected: "master!",
			width:    7,
		},
		{
			name:     "ending with @",
			format:   "%h@",
			expected: "master@",
			width:    7,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, w := Print(all, test.format)
			assertOutput(t, test.expected, actual)
			assertWidth(t, test.width, w)
		})
	}
}

func TestPrinterGroups(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		expected string
		width    int
	}{
		{
			name:     "groups",
			format:   "<[%h][ B%b A%a][ U%u][ C%c]>",
			expected: "<master B5 A4 C3>",
			width:    17,
		},
		{
			name:     "group color",
			format:   "<[#r%h]-[#g%u]%a[-#b%b]>",
			expected: "<\x1b[31mmaster\x1b[0m-4-\x1b[34m5\x1b[0m>",
			width:    12,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, w := Print(all, test.format)
			assertOutput(t, test.expected, actual)
			assertWidth(t, test.width, w)
		})
	}
}

func TestPrinterNonMatching(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		expected string
		width    int
	}{
		{
			name:     "data valid odd",
			format:   "%%%h",
			expected: "%%master",
			width:    8,
		},
		{
			name:     "data valid even",
			format:   "%%%%h",
			expected: "%%%%h",
			width:    5,
		},
		{
			name:     "data invalid odd",
			format:   "%%%z",
			expected: "%%%z",
			width:    4,
		},
		{
			name:     "data invalid even",
			format:   "%%%%z",
			expected: "%%%%z",
			width:    5,
		},
		{
			name:     "color valid odd",
			format:   "###rA",
			expected: "##\x1b[31mA\x1b[0m",
			width:    3,
		},
		{
			name:     "color valid even",
			format:   "####rA",
			expected: "####rA",
			width:    6,
		},
		{
			name:     "color invalid odd",
			format:   "###zA",
			expected: "###zA",
			width:    5,
		},
		{
			name:     "color invalid even",
			format:   "####zA",
			expected: "####zA",
			width:    6,
		},
		{
			name:     "attribute valid odd",
			format:   "@@@bA",
			expected: "@@\x1b[1mA\x1b[0m",
			width:    3,
		},
		{
			name:     "attribute valid even",
			format:   "@@@@bA",
			expected: "@@@@bA",
			width:    6,
		},
		{
			name:     "attribute invalid odd",
			format:   "@@@zA",
			expected: "@@@zA",
			width:    5,
		},
		{
			name:     "attribute invalid even",
			format:   "@@@@zA",
			expected: "@@@@zA",
			width:    6,
		},
		{
			name:     "trailing %",
			format:   "A%",
			expected: "A%",
			width:    2,
		},
		{
			name:     "trailing #",
			format:   "A#",
			expected: "A#",
			width:    2,
		},
		{
			name:     "trailing @",
			format:   "A@",
			expected: "A@",
			width:    2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, w := Print(all, test.format)
			assertOutput(t, test.expected, actual)
			assertWidth(t, test.width, w)
		})
	}
}

func TestPrinterEscape(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		expected string
		width    int
	}{
		{
			name:     "data",
			format:   "A\\%h",
			expected: "A%h",
			width:    3,
		},
		{
			name:     "color",
			format:   "A\\#rB",
			expected: "A#rB",
			width:    4,
		},
		{
			name:     "attribute",
			format:   "A\\!bB",
			expected: "A!bB",
			width:    4,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, w := Print(all, test.format)
			assertOutput(t, test.expected, actual)
			assertWidth(t, test.width, w)
		})
	}
}

func assertOutput(t *testing.T, expected, actual string) {
	t.Helper()
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

func assertWidth(t *testing.T, expected, actual int) {
	t.Helper()
	if expected != actual {
		t.Errorf("Width does not match; expected %d, actual %d", expected, actual)
	}
}
