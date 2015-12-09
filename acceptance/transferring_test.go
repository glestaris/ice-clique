package acceptance_test

import (
	"fmt"
	"net"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/glestaris/ice-clique/acceptance/runner"
	"github.com/glestaris/ice-clique/config"
	"github.com/glestaris/ice-clique/transfer"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Single transferring", func() {
	var (
		client transfer.Transferer
		proc   *runner.ClqProcess
		port   uint16
	)

	BeforeEach(func() {
		var err error

		client = transfer.NewClient(&logrus.Logger{
			Out:       GinkgoWriter,
			Formatter: new(logrus.TextFormatter),
			Level:     logrus.InfoLevel,
		})

		port = 5000 + uint16(GinkgoParallelNode())

		proc, err = startClique(config.Config{
			TransferPort: port,
			RemoteHosts:  []string{},
		})
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(proc.Stop()).To(Succeed())
	})

	It("should accept transfers", func() {
		_, err := client.Transfer(transfer.TransferSpec{
			IP:   net.ParseIP("127.0.0.1"),
			Port: port,
			Size: 10 * 1024 * 1024,
		})
		Expect(err).NotTo(HaveOccurred())
	})

	Context("when trying to run two transfers to the same server", func() {
		var transferStateCh chan string

		BeforeEach(func() {
			transferStateCh = make(chan string)

			go func(c chan string) {
				defer GinkgoRecover()

				c <- "started"

				_, err := client.Transfer(transfer.TransferSpec{
					IP:   net.ParseIP("127.0.0.1"),
					Port: port,
					Size: 20 * 1024 * 1024,
				})
				Expect(err).NotTo(HaveOccurred())

				close(c)
			}(transferStateCh)
		})

		It("should reject the second one", func() {
			Eventually(transferStateCh).Should(Receive())
			time.Sleep(time.Millisecond * 200) // synchronizing the requests

			_, err := client.Transfer(transfer.TransferSpec{
				IP:   net.ParseIP("127.0.0.1"),
				Port: port,
				Size: 10 * 1024 * 1024,
			})
			Expect(err).To(HaveOccurred())

			Eventually(transferStateCh, "15s").Should(BeClosed())
		})
	})
})

var _ = Describe("Clique", func() {
	var (
		procA, procB *runner.ClqProcess
		hosts        []string
	)

	BeforeEach(func() {
		var err error

		hosts = []string{
			fmt.Sprintf("127.0.0.1:%d", 6001+GinkgoParallelNode()),
			fmt.Sprintf("127.0.0.1:%d", 6002+GinkgoParallelNode()),
		}

		procA, err = startClique(config.Config{
			TransferPort:     6001 + uint16(GinkgoParallelNode()),
			RemoteHosts:      []string{hosts[1]},
			InitTransferSize: 1 * 1024 * 1024,
		}, "-debug")
		Expect(err).NotTo(HaveOccurred())

		procB, err = startClique(config.Config{
			TransferPort:     6002 + uint16(GinkgoParallelNode()),
			RemoteHosts:      []string{hosts[0]},
			InitTransferSize: 1 * 1024 * 1024,
		}, "-debug")
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		Expect(procA.Stop()).To(Succeed())
		Expect(procB.Stop()).To(Succeed())
	})

	It("should print the clique result", func() {
		Eventually(procA.Buffer, 2.0).Should(gbytes.Say("new_state=done"))
		Eventually(procB.Buffer).Should(gbytes.Say("new_state=done"))

		procACont := procA.Buffer.Contents()
		Expect(procACont).To(ContainSubstring("Incoming transfer is completed"))
		Expect(procACont).To(ContainSubstring("Outgoing transfer is completed"))
		Expect(procACont).To(ContainSubstring(
			fmt.Sprintf("127.0.0.1:%d", 6002+GinkgoParallelNode()),
		))

		procBCont := procB.Buffer.Contents()
		Expect(procBCont).To(ContainSubstring("Incoming transfer is completed"))
		Expect(procBCont).To(ContainSubstring("Outgoing transfer is completed"))
		Expect(procBCont).To(ContainSubstring(
			fmt.Sprintf("127.0.0.1:%d", 6001+GinkgoParallelNode()),
		))
	})
})