package junitparser

import (
	"bytes"
	"encoding/xml"
)

type TestSuites struct {
	XMLName   xml.Name    `xml:"testsuites"`
	Name      string      `xml:"name,attr"`
	Tests     int         `xml:"tests,attr"`
	Failures  int         `xml:"failures,attr"`
	Time      float64     `xml:"time,attr"`
	TestSuite []TestSuite `xml:"testsuite"`
}

type TestSuite struct {
	XMLName   xml.Name   `xml:"testsuite"`
	Name      string     `xml:"name,attr"`
	Tests     int        `xml:"tests,attr"`
	Failures  int        `xml:"failures,attr"`
	Time      float64    `xml:"time,attr"`
	TestCases []TestCase `xml:"testcase"`
}

type TestCase struct {
	XMLName   xml.Name `xml:"testcase"`
	ClassName string   `xml:"classname,attr"`
	Name      string   `xml:"name,attr"`
	Time      float64  `xml:"time,attr"`
	Failure   *Failure `xml:"failure,omitempty"`
}

type Failure struct {
	Message string `xml:"message,attr"`
	Content string `xml:",chardata"`
}

func ParseJUnitXML(data []byte) (*TestSuites, error) {
	var testSuites TestSuites
	decoder := xml.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&testSuites)
	if err != nil {
		return nil, err
	}

	return &testSuites, nil
}
