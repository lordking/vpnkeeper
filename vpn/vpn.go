package vpn

import (
	"errors"
	"log"
	"regexp"
	"time"

	"github.com/lordking/toolbox/common"
)

const (
	StatusNone int = iota
	StatusConnecting
	StatusConnected
	StatusDisconnected
)

type (
	Service struct {
		ID     string
		Status int
		Name   string
	}

	VPN struct {
		Services  []Service
		Selected  Service
		interrupt chan bool
	}
)

func (vpn *VPN) Fetch() ([]Service, error) {

	var lines []string
	var err error
	if lines, err = common.ExecCommand("scutil", "--nc", "list"); err != nil {
		return nil, err
	}

	//* (Disconnected)   9901A465-C7EA-4592-B97D-29DAA672C158 IPSec              "[VIP] JP - Tokyo 194-Nydus-IPSe [IPSec]
	re1 := regexp.MustCompile("(Disconnected|Connected|Connecting)")
	re2 := regexp.MustCompile("([A-F0-9]{8}\\-[A-F0-9]{4}\\-[A-F0-9]{4}\\-[A-F0-9]{4}\\-[A-F0-9]{12})")
	re3 := regexp.MustCompile("\\p{Z}\"(.*)")
	vpn.Services = make([]Service, 0)
	for _, line := range lines[0:] {

		service := Service{}

		ss1 := re1.FindStringSubmatch(line)
		if len(ss1) > 0 {
			status := ss1[1]
			if status == "Connecting" {
				service.Status = StatusConnecting
			} else if status == "Connected" {
				service.Status = StatusConnected
			} else if status == "Disconnected" {
				service.Status = StatusDisconnected
			} else {
				service.Status = StatusNone
			}
		}

		ss2 := re2.FindStringSubmatch(line)
		if len(ss2) > 0 {
			service.ID = ss2[1]
		}

		ss3 := re3.FindStringSubmatch(line)
		if len(ss3) > 0 {
			service.Name = ss3[1]
			vpn.Services = append(vpn.Services, service)
		}
	}

	return vpn.Services, err
}

func (vpn *VPN) Select(seq int) error {
	listLength := len(vpn.Services)
	if listLength < 0 {
		return errors.New("Not found vpn.")
	}

	if seq >= listLength || seq < 0 {
		return errors.New("Not found this vpn sequence number in current vpn.")
	}

	vpn.Selected = vpn.Services[seq]

	return nil
}

func (vpn *VPN) Stop(service Service) error {
	log.Printf("stop %s", service.Name)
	_, err := common.ExecCommand("scutil", "--nc", "stop", service.ID)
	return err
}

func (vpn *VPN) Start(service Service) error {
	log.Printf("start %s", service.Name)
	_, err := common.ExecCommand("scutil", "--nc", "start", service.ID)
	return err
}

func (vpn *VPN) Status(service Service) (int, error) {

	lines, err := common.ExecCommand("scutil", "--nc", "status", service.ID)
	if err != nil {
		return -1, err
	}

	line := lines[0]
	re := regexp.MustCompile("(.*)\\n")
	ss := re.FindStringSubmatch(line)
	if len(ss) == 0 {
		return -1, errors.New("Not found status")
	}

	var status int
	if len(ss) > 0 {
		statusStr := ss[1]
		if statusStr == "Connecting" {
			status = StatusConnecting
		} else if statusStr == "Connected" {
			status = StatusConnected
		} else if statusStr == "Disconnected" {
			status = StatusDisconnected
		} else {
			status = StatusNone
		}
	}

	return status, nil
}

func (vpn *VPN) RunServ() {

	vpn.interrupt <- false

	for {

		if <-vpn.interrupt {
			vpn.Stop(vpn.Selected)
			return
		}

		status, err := vpn.Status(vpn.Selected)
		if err != nil {
			log.Fatalf("query vpn error: %s", err.Error())
			vpn.interrupt <- true
			continue
		}

		if status == StatusConnecting {
			vpn.interrupt <- false
			time.Sleep(2 * time.Second)

		} else if status == StatusConnected {
			vpn.interrupt <- false
			time.Sleep(2 * time.Second)

		} else {

			err := vpn.Start(vpn.Selected)

			if err == nil {
				vpn.interrupt <- false
				time.Sleep(30 * time.Second)

			} else {
				log.Fatalf("run vpn error: %s", err.Error())
				vpn.interrupt <- true
			}

		}

	}

}

func New() *VPN {
	return &VPN{
		interrupt: make(chan bool, 1),
	}
}
