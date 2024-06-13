package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSapConvergedCloudRegionMappings(t *testing.T) {
	// this test is used for human-testing the catalog response
	prettify := func(content []byte) *bytes.Buffer {
		var prettyJSON bytes.Buffer
		err := json.Indent(&prettyJSON, content, "", "    ")
		assert.NoError(t, err)
		return &prettyJSON
	}

	// given
	suite := NewBrokerSuiteTest(t)
	defer suite.TearDown()

	t.Run("Create catalog - test converged cloud plan not render if not region in path", func(t *testing.T) {
		// when
		resp := suite.CallAPI("GET", fmt.Sprintf("oauth/v2/catalog"), ``)

		content, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		defer resp.Body.Close()

		// then
		assert.NotContains(t, prettify(content).String(), "sap-converged-cloud")
		assert.NotContains(t, prettify(content).String(), "non-existing-1")
		assert.NotContains(t, prettify(content).String(), "non-existing-2")
		assert.NotContains(t, prettify(content).String(), "non-existing-3")
		assert.NotContains(t, prettify(content).String(), "non-existing-4")
	})

	t.Run("Create catalog - test converged cloud plan not render if invalid region in path", func(t *testing.T) {
		// when
		resp := suite.CallAPI("GET", fmt.Sprintf("oauth/non-existing/v2/catalog"), ``)

		content, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		defer resp.Body.Close()

		// then
		assert.NotContains(t, prettify(content).String(), "sap-converged-cloud")
		assert.NotContains(t, prettify(content).String(), "non-existing-1")
		assert.NotContains(t, prettify(content).String(), "non-existing-2")
		assert.NotContains(t, prettify(content).String(), "non-existing-3")
		assert.NotContains(t, prettify(content).String(), "non-existing-4")
	})

	t.Run("Create catalog - test converged cloud plan render if correct region in path", func(t *testing.T) {
		// when
		resp := suite.CallAPI("GET", fmt.Sprintf("oauth/cf-eu20-staging/v2/catalog"), ``)

		content, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		defer resp.Body.Close()

		// then
		assert.Contains(t, prettify(content).String(), "sap-converged-cloud")
		assert.Contains(t, prettify(content).String(), "non-existing-1")
		assert.Contains(t, prettify(content).String(), "non-existing-2")
		assert.NotContains(t, prettify(content).String(), "non-existing-3")
		assert.NotContains(t, prettify(content).String(), "non-existing-4")
	})

}
