package permbits_test

import (
	"os"
	"os/exec"

	"github.com/na4ma4/permbits"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

//nolint:gosec,wrapcheck // test function
func shellTest(mode string) (os.FileMode, error) {
	f, err := os.CreateTemp("", "permbits-test-*")
	if err != nil {
		return 0, err
	}

	name := f.Name()

	f.Close()

	cmd := exec.Command(`chmod`, "ugo=", name)
	if err = cmd.Run(); err != nil {
		return 0, err
	}

	cmd = exec.Command(`chmod`, mode, name)
	if err = cmd.Run(); err != nil {
		return 0, err
	}

	i, err := os.Stat(name)
	if err != nil {
		return 0, err
	}

	if err = os.Remove(name); err != nil {
		return 0, err
	}

	return i.Mode(), nil
}

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
