package brats_test

import (
	"github.com/cloudfoundry/libbuildpack/bratshelper"
	"github.com/cloudfoundry/libbuildpack/cutlass"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// TODO The following test is pending because they currently fail.
// You need to make them pass by implementing the CopyBrats function in brats/brats_suite_test.go
var _ = PDescribe("crystal buildpack", func() {
	bratshelper.UnbuiltBuildpack("crystal", CopyBrats)
	bratshelper.DeployingAnAppWithAnUpdatedVersionOfTheSameBuildpack(CopyBrats)
	bratshelper.StagingWithBuildpackThatSetsEOL("crystal", CopyBrats)
	bratshelper.StagingWithADepThatIsNotTheLatest("crystal", CopyBrats)
	bratshelper.StagingWithCustomBuildpackWithCredentialsInDependencies(`crystal\-[\d\.]+\-linux\-x64\-[\da-f]+\.tgz`, CopyBrats)
	bratshelper.DeployAppWithExecutableProfileScript("crystal", CopyBrats)
	bratshelper.DeployAnAppWithSensitiveEnvironmentVariables(CopyBrats)
	bratshelper.ForAllSupportedVersions("crystal", CopyBrats, func(version string, app *cutlass.App) {
		PushApp(app)

		By("does a thing", func() {
			Expect(app).ToNot(BeNil())
		})
	})
})
