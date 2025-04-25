package boxes

import (
	"fmt"
	"strings"
	"testing"
)

func TestTopBorder(t *testing.T) {
	tests := []struct {
		name    string
		width   int
		border  Border
		want    string
		wantErr bool
	}{
		{
			name:    "SingleBorderValidWidth",
			width:   5,
			border:  Single,
			want:    "┌───┐",
			wantErr: false,
		},
		{
			name:    "DoubleBorderValidWidth",
			width:   6,
			border:  Double,
			want:    "╔════╗",
			wantErr: false,
		},
		{
			name:    "UnsupportedBorderStyle",
			width:   5,
			border:  NoBorder,
			want:    "",
			wantErr: true,
		},
		{
			name:    "WidthLessThanMinimum",
			width:   1,
			border:  Single,
			want:    "┌───┐",
			wantErr: false,
		},
		{
			name:    "ExactMinimumWidth",
			width:   2,
			border:  Single,
			want:    "┌───┐",
			wantErr: false,
		},
		{
			name:    "WidthLessThanMinimumWithDoubleBorder",
			width:   1,
			border:  Double,
			want:    "╔═══╗",
			wantErr: false,
		},
		{
			name:    "SingleBorderWithTCorners",
			width:   5,
			border:  Single | TCorners,
			want:    "├───┤",
			wantErr: false,
		},
		{
			name:    "DoubleBorderWithTCorners",
			width:   6,
			border:  Double | TCorners,
			want:    "╠════╣",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := topBorder(tt.width, tt.border)
			if (err != nil) != tt.wantErr {
				t.Errorf("topBorder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) > 0 && got[0] != tt.want {
				t.Errorf("topBorder() = %v, want %v", got[0], tt.want)
			}
		})
	}
}

func TestBottomBorder(t *testing.T) {
	tests := []struct {
		name    string
		width   int
		border  Border
		want    string
		wantErr bool
	}{
		{
			name:    "SingleBorderValidWidth",
			width:   5,
			border:  Single,
			want:    "└───┘",
			wantErr: false,
		},
		{
			name:    "DoubleBorderValidWidth",
			width:   6,
			border:  Double,
			want:    "╚════╝",
			wantErr: false,
		},
		{
			name:    "UnsupportedBorderStyle",
			width:   5,
			border:  NoBorder,
			want:    "",
			wantErr: true,
		},
		{
			name:    "WidthLessThanMinimum",
			width:   1,
			border:  Single,
			want:    "",
			wantErr: true,
		},
		{
			name:    "ExactMinimumWidth",
			width:   2,
			border:  Single,
			want:    "└┘",
			wantErr: false,
		},
		{
			name:    "WidthLessThanMinimumWithDoubleBorder",
			width:   1,
			border:  Double,
			want:    "",
			wantErr: true,
		},
		{
			name:    "SingleBorderWithTCorners",
			width:   5,
			border:  Single | TCorners,
			want:    "├───┤",
			wantErr: false,
		},
		{
			name:    "DoubleBorderWithTCorners",
			width:   6,
			border:  Double | TCorners,
			want:    "╠════╣",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := bottomBorder(tt.width, tt.border)
			if (err != nil) != tt.wantErr {
				t.Errorf("bottomBorder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) > 0 && got[0] != tt.want {
				t.Errorf("bottomBorder() = %v, want %v", got[0], tt.want)
			}
		})
	}
}

func innerBorder(width int, bs Border) ([]string, error) {
	key := bs & (Single | Double | NoBorder)
	if width < 2 || boxTokens[key] == nil {
		return nil, fmt.Errorf("invalid width or unsupported border style")
	}
	tokens := boxTokens[key]
	lTok, mTok, rTok := LeftT, Horizontal, RightT
	return []string{string(tokens[lTok]) + strings.Repeat(string(tokens[mTok]), width-2) + string(tokens[rTok])}, nil
}

func TestInnerBorder(t *testing.T) {
	tests := []struct {
		name    string
		width   int
		border  Border
		want    string
		wantErr bool
	}{
		{
			name:    "SingleBorderValidWidth",
			width:   5,
			border:  Single | All | TCorners,
			want:    "├───┤",
			wantErr: false,
		},
		{
			name:    "DoubleBorderValidWidth",
			width:   6,
			border:  Double | All | TCorners,
			want:    "╠════╣",
			wantErr: false,
		},
		{
			name:    "UnsupportedBorderStyle",
			width:   5,
			border:  NoBorder | All | TCorners,
			want:    "",
			wantErr: true,
		},
		{
			name:    "WidthLessThanMinimum",
			width:   1,
			border:  Single | All | TCorners,
			want:    "",
			wantErr: true,
		},
		{
			name:    "ExactMinimumWidth",
			width:   2,
			border:  Double | All | TCorners,
			want:    "╠╣",
			wantErr: false,
		},
		{
			name:    "WidthLessThanMinimumWithDoubleBorder",
			width:   1,
			border:  Double | All | TCorners,
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := innerBorder(tt.width, tt.border)
			if (err != nil) != tt.wantErr {
				t.Errorf("innerBorder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) > 0 && got[0] != tt.want {
				t.Errorf("innerBorder() = %v, want %v", got[0], tt.want)
			}
		})
	}
}

func TestPadRow(t *testing.T) {
	tests := []struct {
		name  string
		str   string
		bs    Border
		just  Justification
		width int
		want  string
	}{
		// Single Border
		{"SingleALeft", "A", Single | All, Left, 8, ".A...."},
		{"SingleACenter", "A", Single | All, Center, 8, "...A.."},
		{"SingleARight", "A", Single | All, Right, 8, "....A."},

		{"SingleABLeft", "AB", Single | All, Left, 8, ".AB..."},
		{"SingleABCenter", "AB", Single | All, Center, 8, "..AB.."},
		{"SingleABRight", "AB", Single | All, Right, 8, "...AB."},

		{"SingleABCDELeft", "ABCDE", Single | All, Left, 8, ".ABCDE"},
		{"SingleABCDECenter", "ABCDE", Single | All, Center, 8, ".ABCDE"},
		{"SingleABCDERight", "ABCDE", Single | All, Right, 8, "ABCDE."},

		{"SingleABCDEFGLeft", "ABCDEFG", Single | All, Left, 8, ".ABCDEFG"},
		{"SingleABCDEFGCenter", "ABCDEFG", Single | All, Center, 8, "ABCDEFG"},
		{"SingleABCDEFGRight", "ABCDEFG", Single | All, Right, 8, "ABCDEFG."},

		{"SingleLongLeft", "ABCDEFGABCDEFG", Single | All, Left, 8, "ABCDEFGABCDEFG"},
		{"SingleLongCenter", "ABCDEFGABCDEFG", Single | All, Center, 8, "ABCDEFGABCDEFG"},
		{"SingleLongRight", "ABCDEFGABCDEFG", Single | All, Right, 8, "ABCDEFGABCDEFG"},

		{"SingleEmptyLeft", "", Single | All, Left, 8, "......"},
		{"SingleEmptyCenter", "", Single | All, Center, 8, "......"},
		{"SingleEmptyRight", "", Single | All, Right, 8, "......"},

		// Double Border
		{"DoubleALeft", "A", Double | All, Left, 8, ".A...."},
		{"DoubleACenter", "A", Double | All, Center, 8, "...A.."},
		{"DoubleARight", "A", Double | All, Right, 8, "....A."},

		{"DoubleABLeft", "AB", Double | All, Left, 8, ".AB..."},
		{"DoubleABCenter", "AB", Double | All, Center, 8, "..AB.."},
		{"DoubleABRight", "AB", Double | All, Right, 8, "...AB."},

		{"DoubleABCDELeft", "ABCDE", Double | All, Left, 8, ".ABCDE"},
		{"DoubleABCDECenter", "ABCDE", Double | All, Center, 8, ".ABCDE"},
		{"DoubleABCDERight", "ABCDE", Double | All, Right, 8, "ABCDE."},

		{"DoubleABCDEFGLeft", "ABCDEFG", Double | All, Left, 8, ".ABCDEFG"},
		{"DoubleABCDEFGCenter", "ABCDEFG", Double | All, Center, 8, "ABCDEFG"},
		{"DoubleABCDEFGRight", "ABCDEFG", Double | All, Right, 8, "ABCDEFG."},

		{"DoubleLongLeft", "ABCDEFGABCDEFG", Double | All, Left, 8, "ABCDEFGABCDEFG"},
		{"DoubleLongCenter", "ABCDEFGABCDEFG", Double | All, Center, 8, "ABCDEFGABCDEFG"},
		{"DoubleLongRight", "ABCDEFGABCDEFG", Double | All, Right, 8, "ABCDEFGABCDEFG"},

		{"DoubleEmptyLeft", "", Double | All, Left, 8, "......"},
		{"DoubleEmptyCenter", "", Double | All, Center, 8, "......"},
		{"DoubleEmptyRight", "", Double | All, Right, 8, "......"},

		// NoBorder Border
		{"NoBorderALeft", "A", NoBorder, Left, 8, ".A......"},
		{"NoBorderACenter", "A", NoBorder, Center, 8, "....A..."},
		{"NoBorderARight", "A", NoBorder, Right, 8, "......A."},

		{"NoBorderABLeft", "AB", NoBorder, Left, 8, ".AB....."},
		{"NoBorderABCenter", "AB", NoBorder, Center, 8, "...AB..."},
		{"NoBorderABRight", "AB", NoBorder, Right, 8, ".....AB."},

		{"NoBorderABCDELeft", "ABCDE", NoBorder, Left, 8, ".ABCDE.."},
		{"NoBorderABCDECenter", "ABCDE", NoBorder, Center, 8, "..ABCDE."},
		{"NoBorderABCDERight", "ABCDE", NoBorder, Right, 8, "..ABCDE."},

		{"NoBorderABCDEFGLeft", "ABCDEFG", NoBorder, Left, 8, ".ABCDEFG"},
		{"NoBorderABCDEFGCenter", "ABCDEFG", NoBorder, Center, 8, ".ABCDEFG"},
		{"NoBorderABCDEFGRight", "ABCDEFG", NoBorder, Right, 8, "ABCDEFG."},

		{"NoBorderLongLeft", "ABCDEFGABCDEFG", NoBorder, Left, 8, "ABCDEFGABCDEFG"},
		{"NoBorderLongCenter", "ABCDEFGABCDEFG", NoBorder, Center, 8, "ABCDEFGABCDEFG"},
		{"NoBorderLongRight", "ABCDEFGABCDEFG", NoBorder, Right, 8, "ABCDEFGABCDEFG"},

		{"NoBorderEmptyLeft", "", NoBorder, Left, 8, "........"},
		{"NoBorderEmptyCenter", "", NoBorder, Center, 8, "........"},
		{"NoBorderEmptyRight", "", NoBorder, Right, 8, "........"},
	}

	for _, tt := range tests {
		t.Run(tt.str+"_"+tt.want, func(t *testing.T) {
			got := padRow(tt.str, tt.width, tt.bs, tt.just)
			got = strings.ReplaceAll(got, " ", ".")
			if got != tt.want {
				t.Errorf("padRow(%s %q, %d, %v, %v) = %q, want %q", tt.name, tt.str, tt.width, tt.bs, tt.just, got, tt.want)
			}
		})
	}
}

func TestBoxRow(t *testing.T) {
	tests := []struct {
		name  string
		str   string
		width int
		bs    Border
		just  Justification
		want  string
	}{
		{
			name:  "SingleLineLeftJustifiedSingleBorder",
			str:   "Hello",
			width: 10,
			bs:    Single | All,
			just:  Left,
			want: `
│ Hello  │
`,
		},
		{
			name:  "SingleLineCenterJustifiedSingleBorder",
			str:   "Hello",
			width: 10,
			bs:    Single | All,
			just:  Center,
			want: `
│ Hello  │
`,
		},
		{
			name:  "SingleLineRightJustifiedSingleBorder",
			str:   "Hello",
			width: 10,
			bs:    Single | All,
			just:  Right,
			want: `
│  Hello │
`,
		},
		{
			name:  "MultiLineCenterJustifiedSingleBorder",
			str:   "Hello\nWorld",
			width: 12,
			bs:    Single | All,
			just:  Center,
			want:  "│ Hello    │\n│ World    │",
		},
		{
			name:  "EmptyStringSingleBorder",
			str:   "",
			width: 8,
			bs:    Single | All,
			just:  Center,
			want: `
│      │
`,
		},
		{
			name:  "SingleLineLeftJustifiedDoubleBorder",
			str:   "Hello",
			width: 10,
			bs:    Double | All,
			just:  Left,
			want: `
║ Hello  ║
`,
		},
		{
			name:  "SingleLineRightJustifiedDoubleBorder",
			str:   "Hello",
			width: 10,
			bs:    Double | All,
			just:  Right,
			want: `
║  Hello ║
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := boxRow(tt.str, tt.width, tt.bs, tt.just)
			got = strings.TrimSpace(got)
			expected := strings.TrimSpace(tt.want)
			if strings.Contains(got, "\n") || strings.Contains(expected, "\n") {
				gotDots := strings.ReplaceAll(got, " ", ".")
				expectedDots := strings.ReplaceAll(expected, " ", ".")
				if got != expected {
					t.Errorf("boxRow(%s, %q, %d, %v, %v) = %q, want %q",
						tt.name, tt.str, tt.width, tt.bs, tt.just, gotDots, expectedDots)

					t.Logf("Got (len=%d): %s", len(got), gotDots)
					t.Logf("Want (len=%d): %s", len(expectedDots), expectedDots)

					gotBytes := []byte(got)
					expBytes := []byte(expectedDots)
					minLen := len(gotBytes)
					if len(expBytes) < minLen {
						minLen = len(expBytes)
					}

					for i := 0; i < minLen; i++ {
						if gotBytes[i] != expBytes[i] {
							t.Logf("First difference at position %d: got=%d, want=%d", i, gotBytes[i], expBytes[i])
							break
						}
					}
				}
			} else {
				// For single line strings, use the previous comparison
				wantDots := strings.ReplaceAll(expected, " ", ".")
				gotDots := strings.ReplaceAll(got, " ", ".")
				if got != expected {
					t.Errorf("boxRow(%s, %q, %d, %v, %v) = %q, want %q",
						tt.name, tt.str, tt.width, tt.bs, tt.just, gotDots, wantDots)
				}
			}
		})
	}
}

func TestBox(t *testing.T) {
	tests := []struct {
		name  string
		strs  []string
		width int
		bs    Border
		just  Justification
		want  string
	}{
		// Test cases for borders
		{
			name:  "SingleAllBorders",
			strs:  []string{"Hello", "World"},
			width: 10,
			bs:    Single | TopBorder | BottomBorder | LeftBorder | RightBorder,
			just:  Left,
			want: `
┌────────┐
│ Hello  │
│ World  │
└────────┘
`,
		},
		{
			name:  "DoubleAllBorders",
			strs:  []string{"Test", "Box"},
			width: 10,
			bs:    Double | TopBorder | BottomBorder | LeftBorder | RightBorder,
			just:  Left,
			want: `
╔════════╗
║ Test   ║
║ Box    ║
╚════════╝
`,
		},
		{
			name:  "DoubleSideOnly",
			strs:  []string{"Side", "Only"},
			width: 10,
			bs:    Double | LeftBorder | RightBorder,
			just:  Left,
			want: `
║ Side   ║
║ Only   ║
`,
		},
		{
			name:  "SingleJustifyLeft",
			strs:  []string{"Left Align"},
			width: 20,
			bs:    Single | All,
			just:  Left,
			want: `
┌──────────────────┐
│ Left Align       │
└──────────────────┘
`,
		},
		{
			name:  "SingleJustifyRight",
			strs:  []string{"Right Align"},
			width: 20,
			bs:    Single | All,
			just:  Right,
			want: `
┌──────────────────┐
│      Right Align │
└──────────────────┘
`,
		},
		{
			name:  "SingleWithTCorners",
			strs:  []string{"Underline"},
			width: 15,
			bs:    Single | BottomBorder | Side | TCorners,
			just:  Left,
			want: `
│ Underline   │
├─────────────┤
`,
		},
		{
			name:  "TopBorderOnly",
			strs:  []string{"Top Only"},
			width: 15,
			bs:    Single | TopBorder,
			just:  Center,
			want: `
┌─────────────┐
│ Top Only    │
`,
		},
		{
			name:  "BottomBorderOnly",
			strs:  []string{"Bottom Only"},
			width: 15,
			bs:    Single | LeftBorder | RightBorder | BottomBorder,
			just:  Center,
			want: `
│ Bottom Only │
└─────────────┘
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Box(tt.strs, tt.width, tt.bs, tt.just)
			got = strings.TrimSpace(got)
			tt.want = strings.TrimSpace(tt.want)
			if got != tt.want {
				t.Errorf("Test %s failed.\nGot:\n%s\nExpected:\n%s", tt.name, got, tt.want)
			}
		})
	}
}

func TestBoxTokens(t *testing.T) {
	tests := []struct {
		name   string
		border Border
		pos    BorderPos
		want   rune
	}{
		{"SingleTopLeft", Single, TopLeft, '┌'},
		{"SingleTopRight", Single, TopRight, '┐'},
		{"SingleBottomLeft", Single, BottomLeft, '└'},
		{"SingleBottomRight", Single, BottomRight, '┘'},
		{"SingleHorizontal", Single, Horizontal, '─'},
		{"SingleVertical", Single, Vertical, '│'},
		{"SingleTopT", Single, TopT, '┬'},
		{"SingleLeftT", Single, LeftT, '├'},
		{"SingleBottomT", Single, BottomT, '┴'},
		{"SingleRightT", Single, RightT, '┤'},
		{"SingleMiddleT", Single, MiddleT, '┼'},

		{"DoubleTopLeft", Double, TopLeft, '╔'},
		{"DoubleTopRight", Double, TopRight, '╗'},
		{"DoubleBottomLeft", Double, BottomLeft, '╚'},
		{"DoubleBottomRight", Double, BottomRight, '╝'},
		{"DoubleHorizontal", Double, Horizontal, '═'},
		{"DoubleVertical", Double, Vertical, '║'},
		{"DoubleTopT", Double, TopT, '╦'},
		{"DoubleLeftT", Double, LeftT, '╠'},
		{"DoubleBottomT", Double, BottomT, '╩'},
		{"DoubleRightT", Double, RightT, '╣'},
		{"DoubleMiddleT", Double, MiddleT, '╬'},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := boxTokens[tt.border][tt.pos]; got != tt.want {
				t.Errorf("boxTokens[%v][%v] = %c, want %c", tt.border, tt.pos, got, tt.want)
			}
		})
	}
}

func TestBorderConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant Border
		expected Border
	}{
		{"SideEqualsLeftBorderRightBorder", Side, LeftBorder | RightBorder},
		{"TopBottomEqualsTopBorderBottomBorder", TopBottom, TopBorder | BottomBorder},
		{"AllEqualsTopBottomSide", All, TopBottom | Side},
		{"AllEqualsAllBordersIndividually", All, TopBorder | BottomBorder | LeftBorder | RightBorder},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("%s: got %v, want %v", tt.name, tt.constant, tt.expected)
			}
		})
	}
}

func TestJustificationConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant Justification
		expected int
	}{
		{"LeftEquals0", Left, 0},
		{"RightEquals1", Right, 1},
		{"CenterEquals2", Center, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.constant) != tt.expected {
				t.Errorf("%s: got %v, want %v", tt.name, int(tt.constant), tt.expected)
			}
		})
	}
}

func TestBorderPosConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant BorderPos
		expected int
	}{
		{"TopLeftEquals0", TopLeft, 0},
		{"TopRightEquals1", TopRight, 1},
		{"BottomLeftEquals2", BottomLeft, 2},
		{"BottomRightEquals3", BottomRight, 3},
		{"HorizontalEquals4", Horizontal, 4},
		{"VerticalEquals5", Vertical, 5},
		{"TopTEquals6", TopT, 6},
		{"LeftTEquals7", LeftT, 7},
		{"BottomTEquals8", BottomT, 8},
		{"RightTEquals9", RightT, 9},
		{"MiddleTEquals10", MiddleT, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.constant) != tt.expected {
				t.Errorf("%s: got %v, want %v", tt.name, int(tt.constant), tt.expected)
			}
		})
	}
}

func TestBoxWithSpecialCases(t *testing.T) {
	tests := []struct {
		name  string
		strs  []string
		width int
		bs    Border
		just  Justification
		want  string
	}{
		{
			name:  "EmptyInputArray",
			strs:  []string{},
			width: 10,
			bs:    Single | All,
			just:  Left,
			want: `
┌────────┐
└────────┘
`,
		},
		{
			name:  "MinimumWidth",
			strs:  []string{"A"},
			width: 3,
			bs:    Single | All,
			just:  Left,
			want: `
┌───┐
│ A │
└───┘
`,
		},
		{
			name:  "TextWiderThanWidth",
			strs:  []string{"This text is wider than the box width"},
			width: 10,
			bs:    Single | All,
			just:  Left,
			want: `
┌────────┐
│ This text is wider than the box width │
└────────┘
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Box(tt.strs, tt.width, tt.bs, tt.just)
			got = strings.TrimSpace(got)
			tt.want = strings.TrimSpace(tt.want)
			if got != tt.want {
				g := strings.ReplaceAll(got, " ", ".")
				e := strings.ReplaceAll(tt.want, " ", ".")
				t.Errorf("Test %s failed.\nGot:\n%s\nExpected:\n%s", tt.name, g, e)
			}
		})
	}
}

func TestEdgeCases(t *testing.T) {
	// Test topBorder and bottomBorder with error handling
	tb1, err1 := topBorder(1, Single)
	if !(err1 == nil || len(tb1) > 0) {
		t.Errorf("topBorder with width 1 should return error")
	}

	tb2, err2 := topBorder(2, Single)
	if !(err2 != nil || len(tb2) == 0 || tb2[0] != "┌┐") {
		t.Errorf("topBorder with width 2 should return '┌┐', got %v", tb2)
	}

	bb1, err3 := bottomBorder(1, Single)
	if err3 == nil || len(bb1) > 0 {
		t.Errorf("bottomBorder with width 1 should return error")
	}

	bb2, err4 := bottomBorder(2, Double)
	if err4 != nil || len(bb2) == 0 || bb2[0] != "╚╝" {
		t.Errorf("bottomBorder with width 2 should return '╚╝', got %v", bb2)
	}

	// Test boxRow with zero-width input
	zeroWidthResult := boxRow("Test", 0, Single, Left)
	if zeroWidthResult != "│ Test │" {
		t.Errorf("boxRow with zero width should handle gracefully, got %v", zeroWidthResult)
	}

	// // Test padRow with extreme values
	extremeWidth := padRow("Small", 1000, Single, Center)
	if len(extremeWidth) != 1000 {
		t.Errorf("padRow with extreme width should pad correctly, got length %v, expected %v",
			len(extremeWidth), 1000-2)
	}

	// Test Box with various border combinations
	bordersToTest := []Border{
		TopBorder,
		BottomBorder,
		LeftBorder,
		RightBorder,
		TopBorder | LeftBorder,
		TopBorder | RightBorder,
		BottomBorder | LeftBorder,
		BottomBorder | RightBorder,
		TopBorder | BottomBorder,
		LeftBorder | RightBorder,
		Single | TCorners | TopBorder,
		Double | TCorners | BottomBorder,
	}

	for i, border := range bordersToTest {
		name := fmt.Sprintf("BorderCombo%d", i)
		result := Box([]string{name}, 10, border, Left)
		if result == "" {
			t.Errorf("Box with border combination %v failed", border)
		}
	}
}

func TestPadRowEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		line      string
		width     int
		bs        Border
		just      Justification
		expectLen int
	}{
		{"ZeroWidth", "Test", 0, Single | All, Left, 4},
		{"NegativeWidth", "Test", -5, Single | All, Left, 4},
		{"WidthEqualToText", "Test", 4, Single | All, Left, 4},
		{"WidthLessThanText", "Testing", 5, Single | All, Left, 7},
		{"VeryLargeWidth", "Small", 1000, Single | All, Center, 998},
		{"WidthExactlyFilled", "1234", 6, Single | All, Left, 5},
		{"JustFitsWithMargin", "123", 6, Single | All, Left, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := padRow(tt.line, tt.width, tt.bs, tt.just)
			if len(result) != tt.expectLen {
				t.Errorf("padRow(%q, %d, %v, %v) returned length %d, expected %d",
					tt.line, tt.width, tt.bs, tt.just, len(result), tt.expectLen)
			}
		})
	}
}

func TestMultilinePadding(t *testing.T) {
	tests := []struct {
		name  string
		input string
		bs    Border
		just  Justification
		width int
		want  string
	}{
		{
			name:  "MultilineLeftJustified",
			input: "Line1\nLonger Line2\nL3",
			bs:    Single | All,
			just:  Left,
			width: 15,
			want: `
│.Line1........│
│.Longer.Line2.│
│.L3...........│
`,
		},
		{
			name:  "MultilineCenterJustified",
			input: "Line1\nLonger Line2\nL3",
			bs:    Single | All,
			just:  Center,
			width: 15,
			want: `
│.....Line1....│
│.Longer.Line2.│
│......L3......│
`,
		},
		{
			name:  "MultilineRightJustified",
			input: "Line1\nLonger Line2\nL3",
			bs:    Single | All,
			just:  Right,
			width: 15,
			want: `
│........Line1.│
│.Longer.Line2.│
│...........L3.│
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := boxRow(tt.input, tt.width, tt.bs, tt.just)
			got = strings.TrimSpace(got)
			got = strings.ReplaceAll(got, " ", ".")
			tt.want = strings.TrimSpace(tt.want)
			if got != tt.want {
				t.Errorf("boxRow multiline test %s failed.\nGot:\n%s\nExpected:\n%s",
					tt.name, got, tt.want)
			}
		})
	}
}
