package permbits_test

import (
	"os"
	"os/exec"

	"github.com/na4ma4/go-permbits"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

//nolint:gosec,wrapcheck // test function
func shellTest(mode string) (os.FileMode, error) {
	tmpFile, err := os.CreateTemp("", "permbits-test-*")
	if err != nil {
		return 0, err
	}

	name := tmpFile.Name()

	tmpFile.Close()

	cmd := exec.Command(`chmod`, "ugo=", name)
	if err = cmd.Run(); err != nil {
		return 0, err
	}

	cmd = exec.Command(`chmod`, mode, name)
	if err = cmd.Run(); err != nil {
		return 0, err
	}

	stv, err := os.Stat(name)
	if err != nil {
		return 0, err
	}

	if err = os.Remove(name); err != nil {
		return 0, err
	}

	return stv.Mode(), nil
}

// func testSymbolicModes(t *testing.T, modeString string, is uint32) {
// 	mode, err := permbits.FromString(modeString)
// 	if err != nil {
// 		t.Errorf("[%s] error occurred: %s", modeString, err)
// 		return
// 	}

// 	if uint32(mode) != is {
// 		t.Errorf("[%s] expected %o, go %o", modeString, mode, is)
// 	}
// }

// func TestParseSymbolicModes(t *testing.T) {
// 	list := map[string]uint32{
// 		"a+r":         0o444,
// 		"a-x":         0o000,
// 		"a+rx":        0o555,
// 		"u=rw,g=r,o=": 0o640,
// 		"u+w,go-w":    0o200,
// 		"a+w,go-w":    0o200,
// 		"ug=rw":       0o660,
// 	}
// 	for k, v := range list {
// 		testSymbolicModes(t, k, v)
// 	}
// }

// func TestParseSymbolicMode1(t *testing.T) {
// 	modeString := "u+r,g+r,o+r"
// 	mode, err := permbits.FromString(modeString)
// 	if err != nil {
// 		t.Errorf("[%s] error occurred: %s", modeString, err)
// 		return
// 	}

// 	if uint32(mode) != 0o444 {
// 		t.Errorf("[%s] expected %o, go 0o444", modeString, mode)
// 	}

// 	if !permbits.Is(mode, permbits.UserRead) {
// 		t.Errorf("[%s] expected u+r to be true", modeString)
// 	}
// 	if permbits.Is(mode, permbits.UserWrite) {
// 		t.Errorf("[%s] expected u+w to be false", modeString)
// 	}
// 	if permbits.Is(mode, permbits.UserExecute) {
// 		t.Errorf("[%s] expected u+x to be false", modeString)
// 	}

// 	if !permbits.Is(mode, permbits.GroupRead) {
// 		t.Errorf("[%s] expected g+r to be true", modeString)
// 	}
// 	if permbits.Is(mode, permbits.GroupWrite) {
// 		t.Errorf("[%s] expected g+w to be false", modeString)
// 	}
// 	if permbits.Is(mode, permbits.GroupExecute) {
// 		t.Errorf("[%s] expected g+x to be false", modeString)
// 	}

// 	if !permbits.Is(mode, permbits.OtherRead) {
// 		t.Errorf("[%s] expected o+r to be true", modeString)
// 	}
// 	if permbits.Is(mode, permbits.OtherWrite) {
// 		t.Errorf("[%s] expected o+w to be false", modeString)
// 	}
// 	if permbits.Is(mode, permbits.OtherExecute) {
// 		t.Errorf("[%s] expected o+x to be false", modeString)
// 	}
// }

// func TestParseWeirdValue(t *testing.T) {
// 	modeString := "u+rg+ro+r"
// 	_, err := permbits.FromString(modeString)
// 	if err != nil {
// 		t.Errorf("[%s] error occurred: %s", modeString, err)
// 	}
// }

// func TestParseReturnError(t *testing.T) {
// 	modeString := "a:r"
// 	_, err := permbits.FromString(modeString)
// 	if err == nil {
// 		t.Errorf("[%s] no error occurred", modeString)
// 	}
// }

// func testShellChmod(t *testing.T, modeString string) {
// 	pbMode, err := permbits.FromString(modeString)
// 	if err != nil {
// 		t.Errorf("[%s] error occurred: %s", modeString, err)
// 		return
// 	}

// 	chMode, err := shellTest(modeString)
// 	if err != nil {
// 		t.Errorf("[%s] error occurred: %s", modeString, err)
// 		return
// 	}

// 	if chMode != pbMode {
// 		t.Errorf("[%s] shell mode (%o) and permbits mode (%o) are not identical", modeString, chMode, pbMode)
// 	}
// }

// func TestShellChangeMode(t *testing.T) {
// 	list := []string{
// 		"a+r",
// 		"a-x",
// 		"a+rx",
// 		"u=rw,g=r,o=",
// 		"u+w,go-w",
// 		"ug=rw",
// 		"a=r,g-w",
// 	}
// 	for _, v := range list {
// 		testShellChmod(t, v)
// 	}
// }

var _ = Describe("Parse Test", func() {
	Describe("FromString()", func() {
		DescribeTable(
			"Resolves Symbolic Modes",
			func(modeString string, is int) {
				mode, err := permbits.FromString(modeString)
				Expect(err).NotTo(HaveOccurred())
				Expect(mode).To(BeEquivalentTo(is))
			},
			Entry("a+r", "a+r", 0o444),
			Entry("a-x", "a-x", 0o000),
			Entry("a+rx", "a+rx", 0o555),
			Entry("u=rw,g=r,o=", "u=rw,g=r,o=", 0o640),
			Entry("u+w,go-w", "u+w,go-w", 0o200),
			Entry("a+w,go-w", "a+w,go-w", 0o200),
			Entry("ug=rw", "ug=rw", 0o660),
		)

		It("Resolves Symbolic Modes", func() {
			mode, err := permbits.FromString("u+r,g+r,o+r")
			Expect(err).NotTo(HaveOccurred())
			Expect(mode).To(BeEquivalentTo(0o444))

			Expect(permbits.Is(mode, permbits.UserRead)).To(BeTrue())
			Expect(permbits.Is(mode, permbits.UserWrite)).To(BeFalse())
			Expect(permbits.Is(mode, permbits.UserExecute)).To(BeFalse())

			Expect(permbits.Is(mode, permbits.GroupRead)).To(BeTrue())
			Expect(permbits.Is(mode, permbits.GroupWrite)).To(BeFalse())
			Expect(permbits.Is(mode, permbits.GroupExecute)).To(BeFalse())

			Expect(permbits.Is(mode, permbits.OtherRead)).To(BeTrue())
			Expect(permbits.Is(mode, permbits.OtherWrite)).To(BeFalse())
			Expect(permbits.Is(mode, permbits.OtherExecute)).To(BeFalse())
		})

		DescribeTable(
			"Weird valid parse without error",
			func(mode string) {
				_, err := permbits.FromString(mode)
				Expect(err).NotTo(HaveOccurred())
			},
			Entry("u+rg+ro+r", "u+rg+ro+r"),
		)

		DescribeTable(
			"Returns Errors on invalid mode strings",
			func(mode string) {
				_, err := permbits.FromString(mode)
				Expect(err).To(HaveOccurred())
			},
			Entry("a:r", "a:r"),
		)

		DescribeTable(
			"Compare with shell chmod",
			func(modeString string) {
				pbMode, err := permbits.FromString(modeString)
				Expect(err).NotTo(HaveOccurred())
				chMode, err := shellTest(modeString)
				Expect(err).NotTo(HaveOccurred())
				Expect(chMode).To(BeEquivalentTo(pbMode))
			},
			Entry("a+r", "a+r"),
			Entry("a-x", "a-x"),
			Entry("a+rx", "a+rx"),
			Entry("u=rw,g=r,o=", "u=rw,g=r,o="),
			Entry("u+w,go-w", "u+w,go-w"),
			Entry("ug=rw", "ug=rw"),
			Entry("a=r,g-w", "a=r,g-w"),
		)
	})
})
