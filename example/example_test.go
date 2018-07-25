package example

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/openebs/CITF"
	citfoptions "github.com/openebs/CITF/citf_options"
)

var CitfInstance citf.CITF

func TestIntegrationExample(t *testing.T) {
	RegisterFailHandler(Fail)

	var err error
	// Initializing CITF without config file.
	// Also We should not include K8S as currently we don't have kubernetes environment setup
	CitfInstance, err = citf.NewCITF(citfoptions.CreateOptionsIncludeAllButK8s(""))
	Expect(err).NotTo(HaveOccurred())

	RunSpecs(t, "Integration Test Suite")
}

var _ = BeforeSuite(func() {

	// Setting up the default Platform i.e minikube
	err := CitfInstance.Environment.Setup()
	Expect(err).NotTo(HaveOccurred())

	// You have to update the K8s config when environment has been set up
	// this extra step will be unsolicited in upcoming changes.
	err = CitfInstance.Reload(citfoptions.CreateOptionsIncludeAll(""))
	Expect(err).NotTo(HaveOccurred())

	// Wait until platform is up
	time.Sleep(30 * time.Second)

	err = CitfInstance.K8S.YAMLApply("./nginx-rc.yaml")
	Expect(err).NotTo(HaveOccurred())

	// Wait until the pod is up and running
	time.Sleep(30 * time.Second)
})

var _ = AfterSuite(func() {

	// Tear Down the Platform
	err := CitfInstance.Environment.Teardown()
	Expect(err).NotTo(HaveOccurred())
})

var _ = Describe("Integration Test", func() {
	When("We check the log", func() {
		It("has `started the controller` in the log", func() {
			pods, err := CitfInstance.K8S.GetPods("default", "nginx")
			Expect(err).NotTo(HaveOccurred())

			// Give pods some time to generate logs
			time.Sleep(2 * time.Second)

			// Assuming that only 1 nginx pod is running
			for _, v := range pods {
				log, err := CitfInstance.K8S.GetLog(v.GetName(), "default")
				Expect(err).NotTo(HaveOccurred())

				Expect(log).Should(ContainSubstring("started the controller"))
			}
		})
	})
})
