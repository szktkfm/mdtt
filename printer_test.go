package mdtt

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFindSegment(t *testing.T) {

	tl := tableLocator{}

	md := `# Title
| foo | bar |
| --- | --- |
| baz | bim |
`
	b := bytes.NewBuffer([]byte(md))
	tl.findLocations(b)
	got := tl.locs[0]
	want := TableLocation{
		Start: 9,
		End:   50,
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("differs: (-want +got)\n%s", diff)
	}
}

func TestReplaceTable(t *testing.T) {
	testCases := []struct {
		name string
		src  string
		idx  int
		wnt  string
	}{
		{
			name: "normal",
			src:  "testdata/replace01.md",
			idx:  0,
			wnt:  "testdata/replace01_want.md",
		},
		{
			name: "thematic break",
			src:  "testdata/replace02.md",
			idx:  0,
			wnt:  "testdata/replace02_want.md",
		},
		{
			name: "fenced code block",
			src:  "testdata/replace03.md",
			idx:  0,
			wnt:  "testdata/replace03_want.md",
		},
		{
			name: "multiple tables",
			src:  "testdata/replace04.md",
			idx:  1,
			wnt:  "testdata/replace04_want.md",
		},
		{
			name: "indented code block",
			src:  "testdata/replace05.md",
			idx:  0,
			wnt:  "testdata/replace05_want.md",
		},
		{
			name: "left/center/right aligned",
			src:  "testdata/replace06.md",
			idx:  0,
			wnt:  "testdata/replace06_want.md",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			want, got := testUtilReplaceTable(tc.src, tc.wnt, tc.idx)

			if diff := cmp.Diff(string(want), string(got)); diff != "" {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}
		})
	}
}

func testUtilReplaceTable(src, wnt string, idx int) ([]byte, []byte) {
	tw := TableWriter{}
	fpSrc, _ := os.Open(src)
	defer fpSrc.Close()

	fp_, _ := os.Open(wnt)
	defer fp_.Close()
	md, _ := io.ReadAll(fp_)
	m := NewRoot(
		WithMarkdown(md),
		WithFilePath(wnt),
	)
	m.table = m.tables[idx]
	m.choose = idx

	tw.render(m.table)
	got := tw.replaceTable(fpSrc, idx)

	fpWnt, _ := os.Open(wnt)
	defer fpWnt.Close()
	want, _ := io.ReadAll(fpWnt)
	return got, want
}
