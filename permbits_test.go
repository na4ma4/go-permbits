package permbits_test

import (
	"os"

	"github.com/na4ma4/permbits"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Permbits Test", func() {
	modeAll := os.FileMode(0777)
	modeNone := os.FileMode(0000)

	Describe("Is()", func() {
		DescribeTable(
			"Compares Single Modes True",
			func(mode, is os.FileMode) {
				Expect(permbits.Is(mode, permbits.UserRead)).To(BeTrue())
			},
			Entry("will compare UserRead", modeAll, permbits.UserRead),
			Entry("will compare UserWrite", modeAll, permbits.UserWrite),
			Entry("will compare UserExecute", modeAll, permbits.UserExecute),
			Entry("will compare GroupRead", modeAll, permbits.GroupRead),
			Entry("will compare GroupWrite", modeAll, permbits.GroupWrite),
			Entry("will compare GroupExecute", modeAll, permbits.GroupExecute),
			Entry("will compare OtherRead", modeAll, permbits.OtherRead),
			Entry("will compare OtherWrite", modeAll, permbits.OtherWrite),
			Entry("will compare OtherExecute", modeAll, permbits.OtherExecute),
		)

		DescribeTable(
			"Compares Single Modes False",
			func(mode, is os.FileMode) {
				Expect(permbits.Is(mode, permbits.UserRead)).To(BeFalse())
			},
			Entry("will compare UserRead", modeNone, permbits.UserRead),
			Entry("will compare UserWrite", modeNone, permbits.UserWrite),
			Entry("will compare UserExecute", modeNone, permbits.UserExecute),
			Entry("will compare GroupRead", modeNone, permbits.GroupRead),
			Entry("will compare GroupWrite", modeNone, permbits.GroupWrite),
			Entry("will compare GroupExecute", modeNone, permbits.GroupExecute),
			Entry("will compare OtherRead", modeNone, permbits.OtherRead),
			Entry("will compare OtherWrite", modeNone, permbits.OtherWrite),
			Entry("will compare OtherExecute", modeNone, permbits.OtherExecute),
		)

		It("Compare Multiple Modes", func() {
			mode := os.FileMode(0777)
			Expect(permbits.Is(mode, permbits.UserAll+permbits.GroupAll+permbits.OtherAll)).To(BeTrue())

			mode = os.FileMode(0775)
			Expect(permbits.Is(mode, permbits.UserAll+permbits.GroupAll+permbits.OtherAll)).To(BeFalse())
		})
	})
})
