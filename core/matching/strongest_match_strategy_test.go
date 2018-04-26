package matching_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/models"
	. "github.com/SpectoLabs/hoverfly/core/util"
	. "github.com/onsi/gomega"
)

func Test_ClosestRequestMatcherRequestMatcher_EmptyRequestMatchersShouldMatchOnAnyRequest(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{},
		Response:       testResponse,
	})

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "somehost.com",
		Headers: map[string][]string{
			"sdv": {"ascd"},
		},
	}
	result := matching.StrongestMatchStrategy(r, false, simulation, map[string]string{})

	Expect(result.Pair).ToNot(BeNil())
	Expect(result.Pair.Response.Body).To(Equal("request matched"))
}

func Test_ClosestRequestMatcherRequestMatcher_RequestMatchersShouldMatchOnBody(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("body"),
			},
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Body: "body",
	}
	result := matching.StrongestMatchStrategy(r, false, simulation, map[string]string{})
	Expect(result.Error).To(BeNil())

	Expect(result.Pair.Response.Body).To(Equal("request matched"))
}

func Test_ClosestRequestMatcherRequestMatcher_ReturnResponseWhenAllHeadersMatch(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": {"val1"},
		"header2": {"val2"},
	}

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Headers: headers,
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "http://somehost.com",
		Headers: map[string][]string{
			"header1": {"val1"},
			"header2": {"val2"},
		},
	}

	result := matching.StrongestMatchStrategy(r, false, simulation, map[string]string{})

	Expect(result.Pair.Response.Body).To(Equal("request matched"))
}

func Test_ClosestRequestMatcherRequestMatcher_ReturnNilWhenOneHeaderNotPresentInRequest(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": {"val1"},
		"header2": {"val2"},
	}

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Headers: headers,
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "http://somehost.com",
		Headers: map[string][]string{
			"header1": {"val1"},
		},
	}

	result := matching.StrongestMatchStrategy(r, false, simulation, map[string]string{})

	Expect(result.Pair).To(BeNil())
}

func Test_ClosestRequestMatcherRequestMatcher_ReturnNilWhenOneHeaderValueDifferent(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": {"val1"},
		"header2": {"val2"},
	}

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Headers: headers,
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "somehost.com",
		Headers: map[string][]string{
			"header1": {"val1"},
			"header2": {"different"},
		},
	}
	result := matching.StrongestMatchStrategy(r, false, simulation, map[string]string{})

	Expect(result.Pair).To(BeNil())
}

func Test_ClosestRequestMatcherRequestMatcher_ReturnResponseWithMultiValuedHeaderMatch(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": {"val1-a", "val1-b"},
		"header2": {"val2"},
	}

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Headers: headers,
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "http://somehost.com",
		Body:        "test-body",
		Headers: map[string][]string{
			"header1": {"val1-a", "val1-b"},
			"header2": {"val2"},
		},
	}
	result := matching.StrongestMatchStrategy(r, false, simulation, map[string]string{})

	Expect(result.Pair.Response.Body).To(Equal("request matched"))
}

func Test_ClosestRequestMatcherRequestMatcher_ReturnNilWithDifferentMultiValuedHeaders(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": {"val1-a", "val1-b"},
		"header2": {"val2"},
	}

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Headers: headers,
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "http://somehost.com",
		Headers: map[string][]string{
			"header1": {"val1-a", "val1-differnet"},
			"header2": {"val2"},
		},
	}

	result := matching.StrongestMatchStrategy(r, false, simulation, map[string]string{})

	Expect(result.Pair).To(BeNil())
}

func Test_ClosestRequestMatcherRequestMatcher_EndpointMatchWithHeaders(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": {"val1-a", "val1-b"},
		"header2": {"val2"},
	}

	destination := "testhost.com"
	method := "GET"
	path := "/a/1"
	query := "q=test"

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Headers: headers,
			Destination: &models.RequestFieldMatchers{
				ExactMatch: &destination,
			},
			Path: &models.RequestFieldMatchers{
				ExactMatch: &path,
			},
			Method: &models.RequestFieldMatchers{
				ExactMatch: &method,
			},
			Query: &models.RequestFieldMatchers{
				ExactMatch: &query,
			},
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/a/1",
		Query: map[string][]string{
			"q": []string{"test"},
		},
		Headers: map[string][]string{
			"header1": {"val1-a", "val1-b"},
			"header2": {"val2"},
		},
	}
	result := matching.StrongestMatchStrategy(r, false, simulation, map[string]string{})

	Expect(result.Pair.Response.Body).To(Equal("request matched"))
}

func Test_ClosestRequestMatcherRequestMatcher_EndpointMismatchWithHeadersReturnsNil(t *testing.T) {
	RegisterTestingT(t)

	headers := map[string][]string{
		"header1": {"val1-a", "val1-b"},
		"header2": {"val2"},
	}

	destination := "testhost.com"
	method := "GET"
	path := "/a/1"
	query := "q=test"

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Headers: headers,
			Destination: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "testhost.com",
				ExactMatch: &destination,
			},
			Path: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "/a/1",
				ExactMatch: &path,
			},
			Method: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "GET",
				ExactMatch: &method,
			},
			Query: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "q=test",
				ExactMatch: &query,
			},
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "http://testhost.com",
		Path:        "/a/1",
		Query: map[string][]string{
			"q": []string{"different"},
		},
		Headers: map[string][]string{
			"header1": {"val1-a", "val1-b"},
			"header2": {"val2"},
		},
	}

	result := matching.StrongestMatchStrategy(r, false, simulation, map[string]string{})

	Expect(result.Pair).To(BeNil())
}

func Test_ClosestRequestMatcherRequestMatcher_AbleToMatchAnEmptyPathInAReasonableWay(t *testing.T) {
	RegisterTestingT(t)

	destination := "testhost.com"
	method := "GET"
	path := ""
	query := "q=test"
	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Destination: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "testhost.com",
				ExactMatch: &destination,
			},
			Path: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "",
				ExactMatch: &path,
			},
			Method: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "GET",
				ExactMatch: &method,
			},
			Query: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "q=test",
				ExactMatch: &query,
			},
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Query: map[string][]string{
			"q": []string{"test"},
		},
	}
	result := matching.StrongestMatchStrategy(r, false, simulation, map[string]string{})

	Expect(result.Pair.Response.Body).To(Equal("request matched"))

	r = models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/a/1",
		Query: map[string][]string{
			"q": []string{"test"},
		},
	}

	result = matching.StrongestMatchStrategy(r, false, simulation, map[string]string{})

	Expect(result.Pair).To(BeNil())
}

func Test_ClosestRequestMatcherRequestMatcher_RequestMatchersCanUseGlobsAndBeMatched(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Destination: &models.RequestFieldMatchers{
				GlobMatch: StringToPointer("*.com"),
			},
		},
		Response: testResponse,
	})

	request := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/api/1",
	}

	result := matching.StrongestMatchStrategy(request, false, simulation, map[string]string{})
	Expect(result.Error).To(BeNil())

	Expect(result.Pair.Response.Body).To(Equal("request matched"))
}

func Test_ClosestRequestMatcherRequestMatcher_RequestMatchersCanUseGlobsOnSchemeAndBeMatched(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Scheme: &models.RequestFieldMatchers{
				GlobMatch: StringToPointer("H*"),
			},
		},
		Response: testResponse,
	})

	request := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Scheme:      "http",
		Path:        "/api/1",
	}

	result := matching.StrongestMatchStrategy(request, false, simulation, map[string]string{})
	Expect(result.Error).To(BeNil())

	Expect(result.Pair.Response.Body).To(Equal("request matched"))
}

func Test_ClosestRequestMatcherRequestMatcher_RequestMatchersCanUseGlobsOnHeadersAndBeMatched(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Headers: map[string][]string{
				"unique-header": {"*"},
			},
		},
		Response: testResponse,
	})

	request := models.RequestDetails{
		Method:      "GET",
		Destination: "testhost.com",
		Path:        "/api/1",
		Headers: map[string][]string{
			"unique-header": {"totally-unique"},
		},
	}

	result := matching.StrongestMatchStrategy(request, false, simulation, map[string]string{})
	Expect(result.Error).To(BeNil())

	Expect(result.Pair.Response.Body).To(Equal("request matched"))
}

func Test_ShouldReturnClosestMissIfMatchIsNotFound(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "completemiss",
				ExactMatch: StringToPointer("completemiss"),
			},
			Path: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "completemiss",
				ExactMatch: StringToPointer("completemiss"),
			},
		},
		Response: models.ResponseDetails{
			Body: "one",
		},
	})

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "body",
				ExactMatch: StringToPointer("body"),
				GlobMatch:  StringToPointer("bod*"),
			},
			Path: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "path",
				ExactMatch: StringToPointer("path"),
			},
		},
		Response: models.ResponseDetails{
			Body: "two",
		},
	})

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				Matcher:    "glob",
				Value:      "body",
				ExactMatch: StringToPointer("body"),
			},
			Path: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "path",
				ExactMatch: StringToPointer("path"),
			},
		},
		Response: models.ResponseDetails{
			Body: "three",
		},
	})

	r := models.RequestDetails{
		Body: "body",
		Path: "nomatch",
	}

	result := matching.StrongestMatchStrategy(r, false, simulation, map[string]string{})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Pair).To(BeNil())
	Expect(result.Error.ClosestMiss).ToNot(BeNil())
	Expect(result.Error.ClosestMiss.RequestMatcher.Body[0].Matcher).To(Equal(`exact`))
	Expect(result.Error.ClosestMiss.RequestMatcher.Body[0].Value).To(Equal(`body`))
	Expect(result.Error.ClosestMiss.RequestMatcher.Path[0].Matcher).To(Equal(`exact`))
	Expect(result.Error.ClosestMiss.RequestMatcher.Path[0].Value).To(Equal(`path`))
	Expect(result.Error.ClosestMiss.Response.Body).To(Equal(`two`))
	Expect(result.Error.ClosestMiss.RequestDetails.Body).To(Equal(`body`))
}

func Test_ShouldReturnClosestMissIfMatchIsNotFoundAgain(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				Matcher:    "regex",
				Value:      ".*",
				RegexMatch: StringToPointer(".*"),
			},
			Path: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "miss",
				ExactMatch: StringToPointer("miss"),
			},
			Method: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "GET",
				ExactMatch: StringToPointer("GET"),
			},
		},
		Response: models.ResponseDetails{
			Body: "one",
		},
	})

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      ".*",
				ExactMatch: StringToPointer(".*"),
				// GlobMatch:  StringToPointer("miss"),
			},
			Path: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "miss",
				ExactMatch: StringToPointer("miss"),
			},
		},
		Response: models.ResponseDetails{
			Body: "two",
		},
	})

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "miss",
				ExactMatch: StringToPointer("miss"),
			},
			Path: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "miss",
				ExactMatch: StringToPointer("miss"),
			},
		},
		Response: models.ResponseDetails{
			Body: "three",
		},
	})

	r := models.RequestDetails{
		Body:   "foo",
		Method: "GET",
	}

	result := matching.StrongestMatchStrategy(r, false, simulation, map[string]string{})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Pair).To(BeNil())
	Expect(result.Error.ClosestMiss).ToNot(BeNil())
	Expect(result.Error.ClosestMiss.RequestMatcher.Body[0].Matcher).To(Equal(`regex`))
	Expect(result.Error.ClosestMiss.RequestMatcher.Body[0].Value).To(Equal(`.*`))
	Expect(result.Error.ClosestMiss.RequestMatcher.Path[0].Matcher).To(Equal(`exact`))
	Expect(result.Error.ClosestMiss.RequestMatcher.Path[0].Value).To(Equal(`miss`))
	Expect(result.Error.ClosestMiss.RequestMatcher.Method[0].Matcher).To(Equal(`exact`))
	Expect(result.Error.ClosestMiss.RequestMatcher.Method[0].Value).To(Equal(`GET`))
	Expect(result.Error.ClosestMiss.Response.Body).To(Equal(`one`))
}

func Test_ShouldNotReturnClosestMissWhenThereIsAMatch(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				RegexMatch: StringToPointer(".*"),
			},
			Method: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("GET"),
			},
		},
		Response: models.ResponseDetails{
			Body: "one",
		},
	})

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("miss"),
			},
			Path: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("GET"),
			},
		},
		Response: models.ResponseDetails{
			Body: "two",
		},
	})

	r := models.RequestDetails{
		Body:   "foo",
		Method: "GET",
	}

	result := matching.StrongestMatchStrategy(r, false, simulation, map[string]string{})

	Expect(result.Error).To(BeNil())
	Expect(result.Pair).ToNot(BeNil())
}

func Test__NotBeCachableIfMatchedOnEverythingApartFromHeadersAtLeastOnce(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "POST",
				ExactMatch: StringToPointer("POST"),
			},
			Body: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "body",
				ExactMatch: StringToPointer("body"),
			},
			Scheme: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "http",
				ExactMatch: StringToPointer("http"),
			},
			Query: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "foo=bar",
				ExactMatch: StringToPointer("foo=bar"),
			},
			Path: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "/foo",
				ExactMatch: StringToPointer("/foo"),
			},
			Destination: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "www.test.com",
				ExactMatch: StringToPointer("www.test.com"),
			},
			Headers: map[string][]string{
				"foo": {"bar"},
			},
		},
		Response: testResponse,
	})

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "GET",
				ExactMatch: StringToPointer("GET"),
			},
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Method:      "POST",
		Destination: "www.test.com",
		Query: map[string][]string{
			"foo": []string{"bar"},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "/foo",
		Headers: map[string][]string{
			"miss": {"me"},
		},
	}

	result := matching.StrongestMatchStrategy(r, false, simulation, map[string]string{})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeFalse())
}

func Test__ShouldBeCachableIfMatchedOnEverythingApartFromHeadersZeroTimes(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "POST",
				ExactMatch: StringToPointer("POST"),
			},
			Body: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "body",
				ExactMatch: StringToPointer("body"),
			},
			Scheme: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "http",
				ExactMatch: StringToPointer("http"),
			},
			Query: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "?foo=bar",
				ExactMatch: StringToPointer("?foo=bar"),
			},
			Path: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "/foo",
				ExactMatch: StringToPointer("/foo"),
			},
			Destination: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "www.test.com",
				ExactMatch: StringToPointer("www.test.com"),
			},
			Headers: map[string][]string{
				"foo": {"bar"},
			},
		},
		Response: testResponse,
	})

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "GET",
				ExactMatch: StringToPointer("GET"),
			},
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Method:      "MISS",
		Destination: "www.test.com",
		Query: map[string][]string{
			"foo": []string{"bar"},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "/foo",
		Headers: map[string][]string{
			"miss": {"me"},
		},
	}

	result := matching.StrongestMatchStrategy(r, false, simulation, map[string]string{})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())

	r = models.RequestDetails{
		Method:      "POST",
		Destination: "miss",
		Query: map[string][]string{
			"foo": []string{"bar"},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "/foo",
		Headers: map[string][]string{
			"miss": {"me"},
		},
	}

	result = matching.StrongestMatchStrategy(r, false, simulation, map[string]string{})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())

	r = models.RequestDetails{
		Method:      "POST",
		Destination: "www.test.com",
		Query: map[string][]string{
			"miss": []string{""},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "/foo",
		Headers: map[string][]string{
			"miss": {"me"},
		},
	}

	result = matching.StrongestMatchStrategy(r, false, simulation, map[string]string{})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())

	r = models.RequestDetails{
		Method:      "POST",
		Destination: "www.test.com",
		Query: map[string][]string{
			"foo": []string{"bar"},
		},
		Scheme: "http",
		Body:   "miss",
		Path:   "/foo",
		Headers: map[string][]string{
			"miss": {"me"},
		},
	}

	result = matching.StrongestMatchStrategy(r, false, simulation, map[string]string{})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())

	r = models.RequestDetails{
		Method:      "POST",
		Destination: "www.test.com",
		Query: map[string][]string{
			"foo": []string{"bar"},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "miss",
		Headers: map[string][]string{
			"miss": {"me"},
		},
	}

	result = matching.StrongestMatchStrategy(r, false, simulation, map[string]string{})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())
}

func Test_ShouldSetClosestMissBackToNilIfThereIsAMatchLaterOn(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer(`body`),
			},
			Method: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("GET"),
			},
		},
		Response: models.ResponseDetails{
			Body: "one",
		},
	})

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer(`body`),
			},
			Method: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("POST"),
			},
		},
		Response: models.ResponseDetails{
			Body: "two",
		},
	})

	r := models.RequestDetails{
		Body:   `body`,
		Method: "POST",
	}

	result := matching.StrongestMatchStrategy(r, false, simulation, map[string]string{})

	Expect(result.Error).To(BeNil())
}

func Test_ShouldIncludeHeadersInCalculationForStrongestMatch(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				RegexMatch: StringToPointer(".*"),
			},
			Method: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("GET"),
			},
			Headers: map[string][]string{
				"one":   {"one"},
				"two":   {"one"},
				"three": {"one"},
			},
		},
		Response: models.ResponseDetails{
			Body: "one",
		},
	})

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				RegexMatch: StringToPointer(".*"),
			},
			Method: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("GET"),
				RegexMatch: StringToPointer(".*"),
			},
		},
		Response: models.ResponseDetails{
			Body: "two",
		},
	})

	r := models.RequestDetails{
		Body:   "foo",
		Method: "GET",
		Headers: map[string][]string{
			"one":   {"one"},
			"two":   {"one"},
			"three": {"one"},
		},
	}

	result := matching.StrongestMatchStrategy(r, false, simulation, map[string]string{})

	Expect(result.Error).To(BeNil())
	Expect(result.Pair).ToNot(BeNil())
	Expect(result.Pair.Response.Body).To(Equal("one"))
}

func Test_ShouldIncludeHeadersInCalculationForClosestMiss(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "GET",
				ExactMatch: StringToPointer("GET"),
			},
			Headers: map[string][]string{
				"one":   {"one"},
				"two":   {"one"},
				"three": {"one"},
			},
		},
		Response: models.ResponseDetails{
			Body: "one",
		},
	})

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: &models.RequestFieldMatchers{
				Matcher: "regex",
				Value:   "GET",
				// ExactMatch: StringToPointer("GET"),
				ExactMatch: StringToPointer("GET"),
			},
		},
		Response: models.ResponseDetails{
			Body: "two",
		},
	})

	r := models.RequestDetails{
		Body:   "foo",
		Method: "MISS",
		Headers: map[string][]string{
			"one":   {"one"},
			"two":   {"one"},
			"three": {"one"},
		},
	}

	result := matching.StrongestMatchStrategy(r, false, simulation, map[string]string{})

	// Expect(result.Error).ToNot(BeNil())
	Expect(result.Pair).To(BeNil())
	// Expect(result.Error.ClosestMiss).ToNot(BeNil())
	Expect(result.Error.ClosestMiss.Response.Body).To(Equal("one"))
}

func Test_ShouldReturnFieldsMissedInClosestMiss(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				Matcher:   "glob",
				Value:     "miss",
				GlobMatch: StringToPointer("miss"),
			},
			Path: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "miss",
				ExactMatch: StringToPointer("miss"),
			},
			Method: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "hit",
				ExactMatch: StringToPointer("hit"),
			},
			Destination: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "hit",
				ExactMatch: StringToPointer("hit"),
			},
			Query: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "miss",
				ExactMatch: StringToPointer("miss"),
			},

			Scheme: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "miss",
				ExactMatch: StringToPointer("miss"),
			},
			Headers: map[string][]string{
				"hitKey": {"hitValue"},
			},
		},
		Response: models.ResponseDetails{
			Body: "two",
		},
	})

	r := models.RequestDetails{
		Method:      "hit",
		Destination: "hit",
		Headers: map[string][]string{
			"hitKey": {"hitValue"},
		},
	}

	result := matching.StrongestMatchStrategy(r, false, simulation, map[string]string{})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Pair).To(BeNil())
	Expect(result.Error.ClosestMiss).ToNot(BeNil())
	//TODO: Scheme matching?
	Expect(result.Error.ClosestMiss.MissedFields).To(ConsistOf(`body`, `path`, `query`))
}

func Test_ShouldReturnFieldsMissedInClosestMissAgain(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Body: &models.RequestFieldMatchers{
				Matcher:   "glob",
				Value:     "hit",
				GlobMatch: StringToPointer("hit"),
			},
			Path: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "hit",
				ExactMatch: StringToPointer("hit"),
			},
			Method: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "miss",
				ExactMatch: StringToPointer("miss"),
			},
			Destination: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "miss",
				ExactMatch: StringToPointer("miss"),
			},
			Query: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "hit=",
				ExactMatch: StringToPointer("hit="),
			},
			Scheme: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "hit",
				ExactMatch: StringToPointer("hit"),
			},
			Headers: map[string][]string{
				"miss": {"miss"},
			},
		},
		Response: models.ResponseDetails{
			Body: "two",
		},
	})

	r := models.RequestDetails{
		Body: "hit",
		Path: "hit",
		Query: map[string][]string{
			"hit": []string{""},
		},
	}

	result := matching.StrongestMatchStrategy(r, false, simulation, map[string]string{})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Pair).To(BeNil())
	Expect(result.Error.ClosestMiss).ToNot(BeNil())
	//TODO: Scheme matching?
	Expect(result.Error.ClosestMiss.MissedFields).To(ConsistOf(`method`, `destination`, `headers`))
}

func Test_ShouldReturnMessageForClosestMiss(t *testing.T) {
	RegisterTestingT(t)

	miss := &models.ClosestMiss{
		RequestDetails: models.RequestDetails{
			Path:        "path",
			Method:      "method",
			Destination: "destination",
			Scheme:      "scheme",
			Query: map[string][]string{
				"query": {""},
			},
			Body: "body",
			Headers: map[string][]string{
				"miss": {"miss"},
			},
		},
		State: map[string]string{
			"key1": "value2",
			"key3": "value4",
		},
		Response: v2.ResponseDetailsViewV5{
			Body: "hello world",
			Headers: map[string][]string{
				"hello": {"world"},
			},
			Status: 200,
		},
		RequestMatcher: v2.RequestMatcherViewV5{
			Body: []v2.MatcherViewV5{
				{
					Matcher: "glob",
					Value:   "hit",
				},
			},
			Path: []v2.MatcherViewV5{
				{
					Matcher: "exact",
					Value:   "hit",
				},
			},
			Method: []v2.MatcherViewV5{
				{
					Matcher: "exact",
					Value:   "miss",
				},
			},
			Destination: []v2.MatcherViewV5{
				{
					Matcher: "exact",
					Value:   "miss",
				},
			},
			Query: []v2.MatcherViewV5{
				{
					Matcher: "exact",
					Value:   "hit",
				},
			},
			Scheme: []v2.MatcherViewV5{
				{
					Matcher: "exact",
					Value:   "hit",
				},
			},
			Headers: map[string][]string{
				"miss": {"miss"},
			},
		},
		MissedFields: []string{"body", "path", "method"},
	}

	message := miss.GetMessage()
	Expect(message).To(Equal(
		`

The following request was made, but was not matched by Hoverfly:

{
    "Path": "path",
    "Method": "method",
    "Destination": "destination",
    "Scheme": "scheme",
    "Query": {
        "query": [
            ""
        ]
    },
    "Body": "body",
    "Headers": {
        "miss": [
            "miss"
        ]
    }
}

Whilst Hoverfly has the following state:

{
    "key1": "value2",
    "key3": "value4"
}

The matcher which came closest was:

{
    "path": [
        {
            "matcher": "exact",
            "value": "hit",
            "config": null
        }
    ],
    "method": [
        {
            "matcher": "exact",
            "value": "miss",
            "config": null
        }
    ],
    "destination": [
        {
            "matcher": "exact",
            "value": "miss",
            "config": null
        }
    ],
    "scheme": [
        {
            "matcher": "exact",
            "value": "hit",
            "config": null
        }
    ],
    "query": [
        {
            "matcher": "exact",
            "value": "hit",
            "config": null
        }
    ],
    "body": [
        {
            "matcher": "glob",
            "value": "hit",
            "config": null
        }
    ],
    "headers": {
        "miss": [
            "miss"
        ]
    }
}

But it did not match on the following fields:

[body, path, method]

Which if hit would have given the following response:

{
    "status": 200,
    "body": "hello world",
    "encodedBody": false,
    "headers": {
        "hello": [
            "world"
        ]
    },
    "templated": false
}`))
}

func Test_StrongestMatch_ShouldNotBeCachableIfMatchedOnEverythingApartFromHeadersAtLeastOnce(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("POST"),
			},
			Body: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("body"),
			},
			Scheme: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("http"),
			},
			Query: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("foo=bar"),
			},
			Path: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("/foo"),
			},
			Destination: &models.RequestFieldMatchers{
				ExactMatch: StringToPointer("www.test.com"),
			},
			Headers: map[string][]string{
				"foo": {"bar"},
			},
		},
		Response: testResponse,
	})

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "GET",
				ExactMatch: StringToPointer("GET"),
			},
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Method:      "POST",
		Destination: "www.test.com",
		Query: map[string][]string{
			"foo": []string{"bar"},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "/foo",
		Headers: map[string][]string{
			"miss": {"me"},
		},
	}

	result := matching.StrongestMatchStrategy(r, false, simulation, make(map[string]string))

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeFalse())
}

func Test_StrongestMatch__ShouldBeCachableIfMatchedOnEverythingApartFromHeadersZeroTimes(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "POST",
				ExactMatch: StringToPointer("POST"),
			},
			Body: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "body",
				ExactMatch: StringToPointer("body"),
			},
			Scheme: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "http",
				ExactMatch: StringToPointer("http"),
			},
			Query: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "?foo=bar",
				ExactMatch: StringToPointer("?foo=bar"),
			},
			Path: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "/foo",
				ExactMatch: StringToPointer("/foo"),
			},
			Destination: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "www.test.com",
				ExactMatch: StringToPointer("www.test.com"),
			},
			Headers: map[string][]string{
				"foo": {"bar"},
			},
		},
		Response: testResponse,
	})

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "GET",
				ExactMatch: StringToPointer("GET"),
			},
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Method:      "MISS",
		Destination: "www.test.com",
		Query: map[string][]string{
			"foo": []string{"bar"},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "/foo",
		Headers: map[string][]string{
			"miss": {"me"},
		},
	}

	result := matching.StrongestMatchStrategy(r, false, simulation, make(map[string]string))

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())

	r = models.RequestDetails{
		Method:      "POST",
		Destination: "miss",
		Query: map[string][]string{
			"foo": {"bar"},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "/foo",
		Headers: map[string][]string{
			"miss": {"me"},
		},
	}

	result = matching.StrongestMatchStrategy(r, false, simulation, make(map[string]string))

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())

	r = models.RequestDetails{
		Method:      "POST",
		Destination: "www.test.com",
		Query: map[string][]string{
			"miss": []string{""},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "/foo",
		Headers: map[string][]string{
			"miss": {"me"},
		},
	}

	result = matching.StrongestMatchStrategy(r, false, simulation, make(map[string]string))

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())

	r = models.RequestDetails{
		Method:      "POST",
		Destination: "www.test.com",
		Query: map[string][]string{
			"foo": {"bar"},
		},
		Scheme: "http",
		Body:   "miss",
		Path:   "/foo",
		Headers: map[string][]string{
			"miss": {"me"},
		},
	}

	result = matching.StrongestMatchStrategy(r, false, simulation, make(map[string]string))

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())

	r = models.RequestDetails{
		Method:      "POST",
		Destination: "www.test.com",
		Query: map[string][]string{
			"foo": {"bar"},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "miss",
		Headers: map[string][]string{
			"miss": {"me"},
		},
	}

	result = matching.StrongestMatchStrategy(r, false, simulation, make(map[string]string))

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())
}

func Test_StrongestMatchStrategy_RequestMatchersShouldMatchOnStateAndNotBeCachable(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			RequiresState: map[string]string{"key1": "value1", "key2": "value2"},
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Body: "body",
	}

	result := matching.StrongestMatchStrategy(
		r,
		false,
		simulation,
		map[string]string{"key1": "value1", "key2": "value2"})

	Expect(result.Error).To(BeNil())
	Expect(result.Cachable).To(BeFalse())
	Expect(result.Pair.Response.Body).To(Equal("request matched"))
}

func Test_StrongestMatch_ShouldNotBeCachableIfMatchedOnEverythingApartFromStateAtLeastOnce(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "POST",
				ExactMatch: StringToPointer("POST"),
			},
			Body: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "body",
				ExactMatch: StringToPointer("body"),
			},
			Scheme: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "http",
				ExactMatch: StringToPointer("http"),
			},
			Query: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "foo=bar",
				ExactMatch: StringToPointer("foo=bar"),
			},
			Path: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "/foo",
				ExactMatch: StringToPointer("/foo"),
			},
			Destination: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "www.test.com",
				ExactMatch: StringToPointer("www.test.com"),
			},
			RequiresState: map[string]string{
				"foo": "bar",
			},
		},
		Response: testResponse,
	})

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "GET",
				ExactMatch: StringToPointer("GET"),
			},
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Method:      "POST",
		Destination: "www.test.com",
		Query: map[string][]string{
			"foo": {"bar"},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "/foo",
	}

	result := matching.StrongestMatchStrategy(r, false, simulation, map[string]string{"miss": "me"})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeFalse())
}

func Test_StrongestMatch__ShouldBeCachableIfMatchedOnEverythingApartFromStateZeroTimes(t *testing.T) {
	RegisterTestingT(t)

	simulation := models.NewSimulation()

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "POST",
				ExactMatch: StringToPointer("POST"),
			},
			Body: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "body",
				ExactMatch: StringToPointer("body"),
			},
			Scheme: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "http",
				ExactMatch: StringToPointer("http"),
			},
			Query: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "?foo=bar",
				ExactMatch: StringToPointer("?foo=bar"),
			},
			Path: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "/foo",
				ExactMatch: StringToPointer("/foo"),
			},
			Destination: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "www.test.com",
				ExactMatch: StringToPointer("www.test.com"),
			},
			RequiresState: map[string]string{
				"foo": "bar",
			},
		},
		Response: testResponse,
	})

	simulation.AddRequestMatcherResponsePair(&models.RequestMatcherResponsePair{
		RequestMatcher: models.RequestMatcher{
			Method: &models.RequestFieldMatchers{
				Matcher:    "exact",
				Value:      "GET",
				ExactMatch: StringToPointer("GET"),
			},
		},
		Response: testResponse,
	})

	r := models.RequestDetails{
		Method:      "MISS",
		Destination: "www.test.com",
		Query: map[string][]string{
			"foo": {"bar"},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "/foo",
	}

	result := matching.StrongestMatchStrategy(r, false, simulation, map[string]string{"miss": "me"})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())

	r = models.RequestDetails{
		Method:      "POST",
		Destination: "miss",
		Query: map[string][]string{
			"foo": {"bar"},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "/foo",
	}

	result = matching.StrongestMatchStrategy(r, false, simulation, map[string]string{"miss": "me"})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())

	r = models.RequestDetails{
		Method:      "POST",
		Destination: "www.test.com",
		Query: map[string][]string{
			"miss": []string{""},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "/foo",
	}

	result = matching.StrongestMatchStrategy(r, false, simulation, map[string]string{"miss": "me"})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())

	r = models.RequestDetails{
		Method:      "POST",
		Destination: "www.test.com",
		Query: map[string][]string{
			"foo": {"bar"},
		},
		Scheme: "http",
		Body:   "miss",
		Path:   "/foo",
	}

	result = matching.StrongestMatchStrategy(r, false, simulation, map[string]string{"miss": "me"})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())

	r = models.RequestDetails{
		Method:      "POST",
		Destination: "www.test.com",
		Query: map[string][]string{
			"foo": {"bar"},
		},
		Scheme: "http",
		Body:   "body",
		Path:   "miss",
	}

	result = matching.StrongestMatchStrategy(r, false, simulation, map[string]string{"miss": "me"})

	Expect(result.Error).ToNot(BeNil())
	Expect(result.Cachable).To(BeTrue())
}