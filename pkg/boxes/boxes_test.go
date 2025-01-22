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
			name:    "Single border, valid width",
			width:   5,
			border:  Single,
			want:    "┌───┐",
			wantErr: false,
		},
		{
			name:    "Double border, valid width",
			width:   6,
			border:  Double,
			want:    "╔════╗",
			wantErr: false,
		},
		{
			name:    "Unsupported border style",
			width:   5,
			border:  NoBorder,
			want:    "",
			wantErr: true,
		},
		{
			name:    "Width less than minimum",
			width:   1,
			border:  Single,
			want:    "",
			wantErr: true,
		},
		{
			name:    "Exact minimum width",
			width:   2,
			border:  Single,
			want:    "┌┐",
			wantErr: false,
		},
		{
			name:    "Width less than minimum with Double border",
			width:   1,
			border:  Double,
			want:    "",
			wantErr: true,
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
			name:    "Single border, valid width",
			width:   5,
			border:  Single,
			want:    "└───┘",
			wantErr: false,
		},
		{
			name:    "Double border, valid width",
			width:   6,
			border:  Double,
			want:    "╚════╝",
			wantErr: false,
		},
		{
			name:    "Unsupported border style",
			width:   5,
			border:  NoBorder,
			want:    "",
			wantErr: true,
		},
		{
			name:    "Width less than minimum",
			width:   1,
			border:  Single,
			want:    "",
			wantErr: true,
		},
		{
			name:    "Exact minimum width",
			width:   2,
			border:  Single,
			want:    "└┘",
			wantErr: false,
		},
		{
			name:    "Width less than minimum with Double border",
			width:   1,
			border:  Double,
			want:    "",
			wantErr: true,
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

func TestInnerBorder(t *testing.T) {
	tests := []struct {
		name    string
		width   int
		border  Border
		want    string
		wantErr bool
	}{
		{
			name:    "Single border, valid width",
			width:   5,
			border:  Single | All | TCorners,
			want:    "├───┤",
			wantErr: false,
		},
		{
			name:    "Double border, valid width",
			width:   6,
			border:  Double | All | TCorners,
			want:    "╠════╣",
			wantErr: false,
		},
		{
			name:    "Unsupported border style",
			width:   5,
			border:  NoBorder | All | TCorners,
			want:    "",
			wantErr: true,
		},
		{
			name:    "Width less than minimum",
			width:   1,
			border:  Single | All | TCorners,
			want:    "",
			wantErr: true,
		},
		{
			name:    "Exact minimum width",
			width:   2,
			border:  Double | All | TCorners,
			want:    "╠╣",
			wantErr: false,
		},
		{
			name:    "Width less than minimum with Double border",
			width:   1,
			border:  Double | All | TCorners,
			want:    "",
			wantErr: true,
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
		{"Single_A_Left", "A", Single | All, Left, 8, ".A...."},
		{"Single_A_Center", "A", Single | All, Center, 8, "...A.."},
		{"Single_A_Right", "A", Single | All, Right, 8, "....A."},

		{"Single_AB_Left", "AB", Single | All, Left, 8, ".AB..."},
		{"Single_AB_Center", "AB", Single | All, Center, 8, "..AB.."},
		{"Single_AB_Right", "AB", Single | All, Right, 8, "...AB."},

		{"Single_ABCDE_Left", "ABCDE", Single | All, Left, 8, ".ABCDE"},
		{"Single_ABCDE_Center", "ABCDE", Single | All, Center, 8, ".ABCDE"},
		{"Single_ABCDE_Right", "ABCDE", Single | All, Right, 8, "ABCDE."},

		{"Single_ABCDEFG_Left", "ABCDEFG", Single | All, Left, 8, ".ABCDEFG"},
		{"Single_ABCDEFG_Center", "ABCDEFG", Single | All, Center, 8, "ABCDEFG"},
		{"Single_ABCDEFG_Right", "ABCDEFG", Single | All, Right, 8, "ABCDEFG."},

		{"Single_Long_Left", "ABCDEFGABCDEFG", Single | All, Left, 8, "ABCDEFGABCDEFG"},
		{"Single_Long_Center", "ABCDEFGABCDEFG", Single | All, Center, 8, "ABCDEFGABCDEFG"},
		{"Single_Long_Right", "ABCDEFGABCDEFG", Single | All, Right, 8, "ABCDEFGABCDEFG"},

		{"Single_Empty_Left", "", Single | All, Left, 8, "......"},
		{"Single_Empty_Center", "", Single | All, Center, 8, "......"},
		{"Single_Empty_Right", "", Single | All, Right, 8, "......"},

		// Double Border
		{"Double_A_Left", "A", Double | All, Left, 8, ".A...."},
		{"Double_A_Center", "A", Double | All, Center, 8, "...A.."},
		{"Double_A_Right", "A", Double | All, Right, 8, "....A."},

		{"Double_AB_Left", "AB", Double | All, Left, 8, ".AB..."},
		{"Double_AB_Center", "AB", Double | All, Center, 8, "..AB.."},
		{"Double_AB_Right", "AB", Double | All, Right, 8, "...AB."},

		{"Double_ABCDE_Left", "ABCDE", Double | All, Left, 8, ".ABCDE"},
		{"Double_ABCDE_Center", "ABCDE", Double | All, Center, 8, ".ABCDE"},
		{"Double_ABCDE_Right", "ABCDE", Double | All, Right, 8, "ABCDE."},

		{"Double_ABCDEFG_Left", "ABCDEFG", Double | All, Left, 8, ".ABCDEFG"},
		{"Double_ABCDEFG_Center", "ABCDEFG", Double | All, Center, 8, "ABCDEFG"},
		{"Double_ABCDEFG_Right", "ABCDEFG", Double | All, Right, 8, "ABCDEFG."},

		{"Double_Long_Left", "ABCDEFGABCDEFG", Double | All, Left, 8, "ABCDEFGABCDEFG"},
		{"Double_Long_Center", "ABCDEFGABCDEFG", Double | All, Center, 8, "ABCDEFGABCDEFG"},
		{"Double_Long_Right", "ABCDEFGABCDEFG", Double | All, Right, 8, "ABCDEFGABCDEFG"},

		{"Double_Empty_Left", "", Double | All, Left, 8, "......"},
		{"Double_Empty_Center", "", Double | All, Center, 8, "......"},
		{"Double_Empty_Right", "", Double | All, Right, 8, "......"},

		// NoBorder Border
		{"NoBorder_A_Left", "A", NoBorder, Left, 8, ".A......"},
		{"NoBorder_A_Center", "A", NoBorder, Center, 8, "....A..."},
		{"NoBorder_A_Right", "A", NoBorder, Right, 8, "......A."},

		{"NoBorder_AB_Left", "AB", NoBorder, Left, 8, ".AB....."},
		{"NoBorder_AB_Center", "AB", NoBorder, Center, 8, "...AB..."},
		{"NoBorder_AB_Right", "AB", NoBorder, Right, 8, ".....AB."},

		{"NoBorder_ABCDE_Left", "ABCDE", NoBorder, Left, 8, ".ABCDE.."},
		{"NoBorder_ABCDE_Center", "ABCDE", NoBorder, Center, 8, "..ABCDE."},
		{"NoBorder_ABCDE_Right", "ABCDE", NoBorder, Right, 8, "..ABCDE."},

		{"NoBorder_ABCDEFG_Left", "ABCDEFG", NoBorder, Left, 8, ".ABCDEFG"},
		{"NoBorder_ABCDEFG_Center", "ABCDEFG", NoBorder, Center, 8, ".ABCDEFG"},
		{"NoBorder_ABCDEFG_Right", "ABCDEFG", NoBorder, Right, 8, "ABCDEFG."},

		{"NoBorder_Long_Left", "ABCDEFGABCDEFG", NoBorder, Left, 8, "ABCDEFGABCDEFG"},
		{"NoBorder_Long_Center", "ABCDEFGABCDEFG", NoBorder, Center, 8, "ABCDEFGABCDEFG"},
		{"NoBorder_Long_Right", "ABCDEFGABCDEFG", NoBorder, Right, 8, "ABCDEFGABCDEFG"},

		{"NoBorder_Empty_Left", "", NoBorder, Left, 8, "........"},
		{"NoBorder_Empty_Center", "", NoBorder, Center, 8, "........"},
		{"NoBorder_Empty_Right", "", NoBorder, Right, 8, "........"},
	}

	for _, tt := range tests {
		t.Run(tt.str+"_"+tt.want, func(t *testing.T) {
			got := padRow(tt.str, tt.width, tt.bs, tt.just)
			got = strings.ReplaceAll(got, " ", ".")
			// fmt.Println(tt)
			// fmt.Println(tt.want, len(tt.want))
			// fmt.Println(got, len(got))
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
			name:  "Single line, Left Justified, Single Border",
			str:   "Hello",
			width: 10,
			bs:    Single | All,
			just:  Left,
			want: `
│ Hello  │
`,
		},
		{
			name:  "Single line, Center Justified, Single Border",
			str:   "Hello",
			width: 10,
			bs:    Single | All,
			just:  Center,
			want: `
│  Hello │
`,
		},
		{
			name:  "Single line, Right Justified, Single Border",
			str:   "Hello",
			width: 10,
			bs:    Single | All,
			just:  Right,
			want: `
│  Hello │
`,
		},
		{
			name:  "Multi-line, Center Justified, Single Border",
			str:   "Hello\nWorld",
			width: 12,
			bs:    Single | All,
			just:  Center,
			want: `
│   Hello  │
│   World  │
		`,
		},
		{
			name:  "Empty string, Single Border",
			str:   "",
			width: 8,
			bs:    Single | All,
			just:  Center,
			want: `
│      │
`,
		},
		{
			name:  "Single line, Left Justified, Double Border",
			str:   "Hello",
			width: 10,
			bs:    Double | All,
			just:  Left,
			want: `
║ Hello  ║
`,
		},
		{
			name:  "Single line, Right Justified, Double Border",
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
			tt.want = strings.TrimSpace(tt.want)
			tt.want = strings.ReplaceAll(tt.want, " ", ".")
			got = strings.ReplaceAll(got, " ", ".")
			fmt.Println("test:", tt)
			fmt.Println("want:\n", tt.want, len(tt.want))
			fmt.Println("got :\n", got, len(got))
			if got != tt.want {
				t.Errorf("boxRow(%s, %q, %d, %v, %v) = %q, want %q", tt.name, tt.str, tt.width, tt.bs, tt.just, got, tt.want)
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
			name:  "Single All Borders",
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
			name:  "Double All Borders",
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
			name:  "Double Side Only",
			strs:  []string{"Side", "Only"},
			width: 10,
			bs:    Double | LeftBorder | RightBorder,
			just:  Left,
			want: `
║ Side   ║
║ Only   ║
`,
		},
		// Test cases for justification
		{
			name:  "Single Justify Left",
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
			name:  "Single Justify Right",
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
			name:  "Single Justify Center",
			strs:  []string{"Centered"},
			width: 20,
			bs:    Single | All,
			just:  Center,
			want: `
┌──────────────────┐
│     Centered     │
└──────────────────┘
`,
		},
		// Mixed scenarios
		{
			name:  "Single with TCorners",
			strs:  []string{"Underline"},
			width: 15,
			bs:    Single | BottomBorder | Side | TCorners,
			just:  Left,
			want: `
│ Underline   │
├─────────────┤
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Box(tt.strs, tt.width, tt.bs, tt.just)
			got = strings.TrimSpace(got)
			tt.want = strings.TrimSpace(tt.want)
			// fmt.Println(tt)
			// fmt.Println(tt.want, len(tt.want))
			// fmt.Println(got, len(got))
			if got != tt.want {
				t.Errorf("Test %s failed.\nGot:\n%s\nExpected:\n%s", tt.name, got, tt.want)
			}
		})
	}
}
