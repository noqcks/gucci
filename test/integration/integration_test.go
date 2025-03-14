package integration_test

import (
	"os"
	"os/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("gucci", func() {

	Describe("template source", func() {

		It("reads stdin", func() {
			gucciCmd := exec.Command(gucciPath)

			tpl, err := os.Open(FixturePath("simple.tpl"))
			defer tpl.Close()
			Expect(err).NotTo(HaveOccurred())
			gucciCmd.Stdin = tpl

			session := RunWithError(gucciCmd, 1)

			Expect(string(session.Err.Contents())).To(Equal("Failed to parse standard input: template: -:1:8: executing \"-\" at <.FOO>: map has no entry for key \"FOO\"\n"))
		})

		It("loads file", func() {
			gucciCmd := exec.Command(gucciPath, FixturePath("simple.tpl"))

			session := RunWithError(gucciCmd, 1)

			Expect(string(session.Err.Contents())).To(Equal("Failed to parse standard input: template: simple.tpl:1:8: executing \"simple.tpl\" at <.FOO>: map has no entry for key \"FOO\"\n"))
		})

	})

	Describe("variable source", func() {

		It("reads env vars", func() {
			gucciCmd := exec.Command(gucciPath, FixturePath("simple.tpl"))
			gucciCmd.Env = []string{
				"FOO=bar",
			}

			session := Run(gucciCmd)

			Expect(string(session.Out.Contents())).To(Equal("text bar text\n"))
		})

		It("loads vars file", func() {
			gucciCmd := exec.Command(gucciPath,
				"-f", FixturePath("simple_vars.yaml"),
				FixturePath("simple.tpl"))

			session := Run(gucciCmd)

			Expect(string(session.Out.Contents())).To(Equal("text bar text\n"))
		})

		It("uses vars options", func() {
			gucciCmd := exec.Command(gucciPath,
				"-s", "FOO=bar",
				FixturePath("simple.tpl"))

			session := Run(gucciCmd)

			Expect(string(session.Out.Contents())).To(Equal("text bar text\n"))
		})
	})

	Describe("variable precedence", func() {

		It("should override variables sources", func() {
			gucciCmd := exec.Command(gucciPath,
				"-s", "C=from_opt",
				"-f", FixturePath("precedence_vars.yaml"),
				FixturePath("precedence.tpl"))
			gucciCmd.Env = []string{
				"B=from_env",
				"C=from_env",
			}

			session := Run(gucciCmd)

			Expect(string(session.Out.Contents())).To(Equal("A=from_file\nB=from_env\nC=from_opt\n"))
		})

	})

	Describe("variable nesting", func() {

		It("should nest file variables", func() {
			gucciCmd := exec.Command(gucciPath,
				"-f", FixturePath("nesting_vars.yaml"),
				FixturePath("nesting.tpl"))

			session := Run(gucciCmd)

			Expect(string(session.Out.Contents())).To(Equal("yep\n"))
		})

		It("should nest option variables", func() {
			gucciCmd := exec.Command(gucciPath,
				"-s", "foo.bar.baz=yep",
				FixturePath("nesting.tpl"))

			session := Run(gucciCmd)

			Expect(string(session.Out.Contents())).To(Equal("yep\n"))
		})

	})

	Describe("toJson and mustToJson functions", func() {
		It("should handle map[interface {}]interface {} in toJson and mustToJson", func() {
			gucciCmd := exec.Command(gucciPath,
				"-f", FixturePath("issue63.yaml"),
				FixturePath("issue63.tpl"))

			session := Run(gucciCmd)

			Expect(session.ExitCode()).To(Equal(0))
		})
	})

})
