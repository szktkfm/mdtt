package mdtt

import (
	"bytes"
	"fmt"
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
	tw := TableWriter{}
	fp, _ := os.Open("testdata/replace01.md")
	defer fp.Close()

	m := NewRoot(
		WithMDFile("testdata/replace_want01.md"),
	)

	tw.render(m.table)
	b := tw.replaceTable(fp)

	fmt.Println(string(b))
	// fp2, _ := os.Open("testdata/replace_want01.md")
	// defer fp2.Close()
	// want, _ := io.ReadAll(fp2)
	want := `# Title
| foo2                               | bar                                 |
| -----------------------------------| ----------------------------------- |
| bazaaaa                            | bim                                 |
`

	if diff := cmp.Diff(string(want), string(b)); diff != "" {
		t.Errorf("differs: (-want +got)\n%s", diff)
	}

}
