package junitparser

import (
	"encoding/json"
	"encoding/xml"
	"strings"
	"testing"
)

// TestJUnitXMLParser tests the parseJUnitXML function with a simple JUnit XML string.
func TestJUnitXMLParser(t *testing.T) {
	// Define a sample JUnit XML string
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
	<testsuites name="example" tests="2" failures="1" time="1.506">
	  <testsuite name="exampleTestSuite" tests="2" failures="1" time="1.206">
	    <testcase classname="exampleTestClass" name="testExample1" time="0.603"/>
	    <testcase classname="exampleTestClass" name="testExample2" time="0.603">
	      <failure message="AssertionError">Expected value did not match actual value</failure>
	    </testcase>
	  </testsuite>
	</testsuites>`

	// Use an XML reader for the test string
	xmlReader := strings.NewReader(xmlData)
	decoder := xml.NewDecoder(xmlReader)

	// Parse the XML
	var testSuites TestSuites
	err := decoder.Decode(&testSuites)
	if err != nil {
		t.Fatalf("Error parsing XML: %v", err)
	}

	// Convert the result to JSON to easily compare against the expected output
	jsonResult, err := json.MarshalIndent(testSuites, "", "  ")
	if err != nil {
		t.Fatalf("Error converting parsed XML to JSON: %v", err)
	}

	// Define the expected JSON output
	expectedJSON := `{
		"name": "example",
		"tests": 2,
		"failures": 1,
		"time": 1.506,
		"testsuite": [
		  {
		    "name": "exampleTestSuite",
		    "tests": 2,
		    "failures": 1,
		    "time": 1.206,
		    "testcase": [
		      {
		        "classname": "exampleTestClass",
		        "name": "testExample1",
		        "time": 0.603
		      },
		      {
		        "classname": "exampleTestClass",
		        "name": "testExample2",
		        "time": 0.603,
		        "failure": {
		          "message": "AssertionError",
		          "content": "Expected value did not match actual value"
		        }
		      }
		    ]
		  }
		]
	}`

	// Compare the output with the expected JSON
	if string(jsonResult) != expectedJSON {
		t.Errorf("Parsed output does not match expected output.\nGot: %s\nWant: %s", jsonResult, expectedJSON)
	}
}
