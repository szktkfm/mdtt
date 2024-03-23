package mdtt

import (
	"bytes"
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
