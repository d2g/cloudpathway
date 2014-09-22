package datastore

import (
	"bytes"
	"encoding/json"
	"net"
	"testing"
)

func TestConfigurationJSONMarshalling(test *testing.T) {
	var err error

	exampleDevice := Device{}
	exampleDevice.MACAddress, err = net.ParseMAC("00:00:00:00:00:00")
	exampleDevice.IPAddress = net.IPv4(192, 168, 0, 1)
	exampleDevice.Hostname = "ExampleHostname"
	exampleDevice.Nickname = "Example Nickname"
	exampleDevice.CurrentUser = nil
	exampleDevice.DefaultUser = nil

	test.Logf("Configuration Object:%v\n", exampleDevice)

	byteExampleDevice, err := json.Marshal(exampleDevice)
	if err != nil {
		test.Error("Error Marshaling to JSON:" + err.Error())
	}

	test.Log("As JSON:" + string(byteExampleDevice))

	endExampleDevice := Device{}
	err = json.Unmarshal(byteExampleDevice, &endExampleDevice)
	if err != nil {
		test.Error("Error Unmarshaling to JSON:" + err.Error())
	}

	test.Logf("Configuration Object:%v\n", endExampleDevice)
}

func TestCreateAndRecall(test *testing.T) {
	helper, err := GetDeviceHelper()
	if err != nil {
		test.Error("Error Getting Helper:" + err.Error())
	}

	exampleDevice := Device{}
	exampleDevice.MACAddress, err = net.ParseMAC("00:00:00:00:00:00")
	if err != nil {
		test.Error("Error Parsing Mac:" + err.Error())
	}

	exampleDevice.IPAddress = net.IPv4(192, 168, 0, 1)
	exampleDevice.Hostname = "ExampleHostname"
	exampleDevice.Nickname = "Example Nickname"
	exampleDevice.CurrentUser = nil
	exampleDevice.DefaultUser = nil

	err = helper.SetDevice(&exampleDevice)
	if err != nil {
		test.Error("Error Saving Device:" + err.Error())
	}

	recalledDevice, err := helper.GetDevice(exampleDevice.MACAddress)
	if err != nil {
		test.Error("Error Recalling Device:" + err.Error())
	}

	if !bytes.Equal(exampleDevice.MACAddress, recalledDevice.MACAddress) ||
		!bytes.Equal(exampleDevice.IPAddress, recalledDevice.IPAddress) ||
		exampleDevice.Hostname != recalledDevice.Hostname ||
		exampleDevice.Nickname != recalledDevice.Nickname {
		test.Logf("Saved:%v\n", exampleDevice)
		test.Logf("Recalled:%v\n", recalledDevice)
		test.Error("Recalled Device Doesn't Match Saved:" + err.Error())
	}

	recalledbyip, err := helper.GetDeviceByIP(net.IPv4(192, 168, 0, 1))
	if err != nil {
		test.Error("Error Recalling Device By IP:" + err.Error())
	}

	if !bytes.Equal(exampleDevice.MACAddress, recalledbyip.MACAddress) ||
		!bytes.Equal(exampleDevice.IPAddress, recalledbyip.IPAddress) ||
		exampleDevice.Hostname != recalledbyip.Hostname ||
		exampleDevice.Nickname != recalledbyip.Nickname {
		test.Logf("Saved:%v\n", exampleDevice)
		test.Logf("Recalled By IP:%v\n", recalledbyip)
		test.Error("Recalled By IP Device Doesn't Match Saved:" + err.Error())
	}

}
