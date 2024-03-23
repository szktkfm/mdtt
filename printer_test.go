package mdtt

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/charmbracelet/log"
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
		wnt  string
	}{
		{
			name: "Test Case 1",
			src:  "testdata/replace01.md",
			wnt:  "testdata/replace_want01.md",
		},
		// Add more test cases here
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			want, got := testUtilReplaceTable(tc.src, tc.wnt)

			if diff := cmp.Diff(string(want), string(got)); diff != "" {
				t.Errorf("differs: (-want +got)\n%s", diff)
			}
		})
	}
}

func testUtilReplaceTable(src, wnt string) ([]byte, []byte) {
	tw := TableWriter{}
	fp, _ := os.Open(src)
	defer fp.Close()

	m := NewRoot(
		WithMDFile(wnt),
	)

	tw.render(m.table)
	got := tw.replaceTable(fp)

	log.Debug(string(got))
	fp2, _ := os.Open(wnt)
	defer fp2.Close()
	want, _ := io.ReadAll(fp2)
	return got, want
}
