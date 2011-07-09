package gobdd

import (
  "testing"
  "bytes"
  "reflect"
  "fmt"
  "strings"
)

func TestDescribe(t *testing.T) {
  i := 0
  
  Describe("foo", func() {
    assertEqualObjects(t, len(testingContexts), 1)
    assertEqualObjects(t, testingContexts[0].Description, "foo")
    i++
    
    Describe("bar", func() {
      assertEqualObjects(t, len(testingContexts), 2)
      assertEqualObjects(t, testingContexts[0].Description, "foo")
      assertEqualObjects(t, testingContexts[1].Description, "bar")
      i++
    })
  })
  
  assertEqualObjects(t, i, 2)
  assertEqualObjects(t, len(testingContexts), 0)
}

func TestIt(t *testing.T) {
  i := 0
  
  Describe("foo", func() {
    It("is very cool", func() {
      assertEqualObjects(t, testingCurrentIt, "is very cool")
      i++
    })
  })
  
  assertEqualObjects(t, i, 1)
}

func TestBeforeEach(t *testing.T) {
  Describe("foo", func() {
    i := 0
    
    BeforeEach(func() {
      i++
    })
    
    It("will run", func() {
      assertEqualObjects(t, i, 1)
    })
    
    It("will run here too", func() {
      assertEqualObjects(t, i, 2)
    })
    
    Describe("bar", func() {
      j := 0
      
      BeforeEach(func() {
        j++
      })
      
      It("will run", func() {
        assertEqualObjects(t, j, 1)
        assertEqualObjects(t, i, 3)
      })
    
      It("will run here too", func() {
        assertEqualObjects(t, j, 2)
        assertEqualObjects(t, i, 4)
      })
    })
  })
}

func TestEqualAssertion(t *testing.T) {
  Describe("foo", func() {
    Describe("bar", func() {
      It("is good", func() {
        Expect(42).toEqual(42)
        Expect(24).toEqual(23)
      })
    })
  })
  
  assertEqualObjects(t, len(testingErrors), 1)
  assertEqualObjects(t, testingErrors[0].String, "expected: 23\n     got: 24\n")
  assertDeepEqualObjects(t, testingErrors[0].Contexts, []string{"foo", "bar"})
  fmt.Println(testingErrors[0].ErrorLine)
  assertEqualObjects(t, strings.HasSuffix(testingErrors[0].ErrorLine, "bdd_test.go:83"), true)
  
  testingErrors = testingErrors[0:0] // cleanup.. bleh
}

func TestPrintSpecReport(t *testing.T) {
  Describe("foo", func() {
    Describe("bar", func() {
      It("is cool", func() {
        Expect(23).toEqual(24)
      })
      It("is lame", func() {
        Expect(23).toEqual(23)
      })
    })
  })
  
  report, ok := BuildSpecReport()
  assertEqualObjects(t, ok, false)
  assertEqualObjects(t,
    strings.Contains(report, testingErrors[0].String) &&
    strings.Contains(report, testingErrors[0].ErrorLine) &&
    strings.Contains(report, testingErrors[0].Contexts[0]) &&
    strings.Contains(report, testingErrors[0].Contexts[1]) &&
    strings.Contains(report, "is cool"),
    true)
  
  stream := bytes.NewBufferString("")
  specReportStream = stream
  debugTesting = true
  PrintSpecReport()
  assertEqualObjects(t, stream.String(), report)
}

func TestPrintGreen(t *testing.T) {
  Describe("foo", func() {
    Describe("bar", func() {
      It("is cool", func() {
        Expect(23).toEqual(23)
      })
      It("is cool", func() {
        Expect(23).toEqual(23)
      })
    })
  })
  
  report, ok := BuildSpecReport()
  assertEqualObjects(t, ok, true)
  assertEqualObjects(t,
    strings.Contains(report, "All tests passed.") &&
    strings.Contains(report, "2 examples") &&
    strings.Contains(report, "0 failures"),
    true)
}

func init() {
  defer PrintSpecReport()
  
  type MyGreatTestType struct {
    Name string
    Age int
  }
  
  MyNil := func() *MyGreatTestType {
    return nil
  }
  
  MyNonNil := func() *MyGreatTestType {
    return &MyGreatTestType{}
  }
  
  Describe("matchers", func() {
    Describe("not equals", func() {
      It("matches on simple objects", func() {
        Expect(&MyGreatTestType{"john", 23}).toNotEqual(&MyGreatTestType{"john", 23})
        Expect("foo").toEqual("foo")
        Expect("foo").toNotEqual("bar")
      })
      It("matches for typed-nil", func() {
        Expect(MyNil()).toBeNil()
        Expect(MyNonNil()).toNotBeNil()
      })
      It("matches for nil", func() {
        Expect(nil).toBeNil()
        Expect(true).toNotBeNil()
      })
    })
    Describe("deep equals matcher", func() {
      It("matches what equals does not", func() {
        Expect(&MyGreatTestType{"john", 23}).toDeepEqual(&MyGreatTestType{"john", 23})
        Expect("foo").toDeepEqual("foo")
      })
    })
    Describe("exception-rescuing matchers", func() {
      It("is super cool", func() {
        Expect(func() { panic("foobar!") }).toPanicWith("foobar!")
        Expect(func() {}).toNotPanic()
      })
    })
  })
}

func assertDeepEqualObjects(t *testing.T, obj interface{}, expected interface{}) {
  if !reflect.DeepEqual(obj, expected) {
    t.Errorf("expected [%v] to equal [%v]", expected, obj)
  }
}

func assertEqualObjects(t *testing.T, obj interface{}, expected interface{}) {
  if obj != expected {
    t.Errorf("expected [%v] to equal [%v]", expected, obj)
  }
}
