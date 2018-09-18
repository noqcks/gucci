package integration_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

var gucciPath string

// For determining package name
type Noop struct{}

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

func getGucciPackage() string {
	thisPkg := reflect.TypeOf(Noop{}).PkgPath()
	parts := strings.Split(thisPkg, "/")
	return strings.Join(parts[0:len(parts)-2], "/")
}

var _ = SynchronizedBeforeSuite(func() []byte {
	binPath, err := gexec.Build(getGucciPackage())
	Expect(err).NotTo(HaveOccurred())
	return []byte(binPath)
}, func(data []byte) {
	gucciPath = string(data)
})

func Run(gucciCmd *exec.Cmd) *gexec.Session {
	session, err := gexec.Start(gucciCmd, GinkgoWriter, GinkgoWriter)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session).Should(gexec.Exit(0))
	return session
}

func FixturePath(fixture string) string {
	_, basedir, _, ok := runtime.Caller(0)
	if !ok {
		// Don't assert here because it can be called outside of an It()
		panic(fmt.Errorf("Fixture not found: %s", fixture))
	}

	f := filepath.Join(basedir, "../fixtures/", fixture)
	return f
}
