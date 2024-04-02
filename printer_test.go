package mdtt

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFindSegment(t *testing.T) {
	tw := TableWriter{}

	md := `# Title
| foo | bar |
| --- | --- |
| baz | bim |
`
	b := bytes.NewBuffer([]byte(md))
	tw.findSegment(b)
	got := tw.seg
	want := TableSegment{
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
			name: "Test Case 1",
			src:  "testdata/replace01.md",
			idx:  0,
			wnt:  "testdata/replace01_want.md",
		},
		{
			name: "Test Case 2",
			src:  "testdata/replace02.md",
			idx:  0,
			wnt:  "testdata/replace02_want.md",
		},
		{
			name: "Test Case 3",
			src:  "testdata/replace03.md",
			idx:  0,
			wnt:  "testdata/replace03_want.md",
		},
		{
			name: "Test Case 4",
			src:  "testdata/replace04.md",
			idx:  1,
			wnt:  "testdata/replace04_want.md",
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
	fp, _ := os.Open(src)
	defer fp.Close()

	m := NewRoot(
		WithMDFile(wnt),
	)
	m.table = m.tables[idx]
	m.choose = idx

	tw.render(m.table)
	got := tw.replaceTable(fp, idx)

	fp2, _ := os.Open(wnt)
	defer fp2.Close()
	want, _ := io.ReadAll(fp2)
	return got, want
}
