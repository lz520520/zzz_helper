package file2

import (
	"testing"
)

func TestPath(t *testing.T) {
	err := CopyDir("E:\\code\\go\\Level6WorkSpace\\Level6\\tmpout\\compile_server\\new\\980269c2-1e3f-437a-9fcc-48e4763ef14b", "E:\\code\\go\\Level6WorkSpace\\Level6\\tmpout\\compile_server\\code\\980269c2-1e3f-437a-9fcc-48e4763ef14b")
	if err != nil {
		t.Log(err)
	}
	//t.Log(filepath.Join("c:\\users", "d:\\users"))
}
