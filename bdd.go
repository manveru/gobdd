package gobdd

import (
  "fmt"
  "runtime"
  "reflect"
  "strings"
  "io"
  "os"
)


type context struct {
  Description string
  BeforeEach func()
}

type testingError struct {
  String string
  ItString string
  Contexts []string
  ErrorLine string
}

type Expectation struct {
  value interface{}
}

var debugTesting = false
var testingContexts []*context
var testingCurrentIt string
var testingExamples int
var testingErrors []*testingError
var specReportStream io.ReadWriter

var redColor = fmt.Sprintf("%c[31m", 27)
var greenColor = fmt.Sprintf("%c[32m", 27)
var resetColors = fmt.Sprintf("%c[0m", 27)

func Describe(s string, f func()) {
  testingContexts = append(testingContexts, &context{Description: s})
  f()
  testingContexts = testingContexts[0:len(testingContexts)-1]
}

func It(s string, f func()) {
  for _, testingContext := range testingContexts {
    if beforeFunc := testingContext.BeforeEach; beforeFunc != nil {
      beforeFunc()
    }
  }
  
  testingCurrentIt = s
  testingExamples++
  
  f()
}

func BeforeEach(f func()) {
  testingContexts[len(testingContexts)-1].BeforeEach = f
}

func getErrorLine() string {
  pc, _, _, _ := runtime.Caller(3)
  file, line := runtime.FuncForPC(pc).FileLine(pc)
  return fmt.Sprintf("%s:%d", file, line)
}

func Expect(obj interface{}) *Expectation {
  return &Expectation{obj}
}

func (e *Expectation) toEqual(obj interface{}) {
  if e.value != obj {
    addErrorObject("expected: %v\n     got: %v\n", obj, e.value)
  }
}

func (e *Expectation) toNotEqual(obj interface{}) {
  if e.value == obj {
    addErrorObject(" expected: %v\nto not be: %v\n", obj, e.value)
  }
}

func (e *Expectation) toDeepEqual(obj interface{}) {
  if !reflect.DeepEqual(e.value, obj) {
    addErrorObject("    expected: %v\nto deeply be: %v\n", obj, e.value)
  }
}

func (e *Expectation) toBeNil() {
  if !reflect.DeepEqual(reflect.ValueOf(nil), reflect.Indirect(reflect.ValueOf(e.value))) {
    addErrorObject("expected to be nil,\n           but got: %v\n", e.value)
  }
}

func (e *Expectation) toNotBeNil() {
  if reflect.DeepEqual(reflect.ValueOf(nil), reflect.Indirect(reflect.ValueOf(e.value))) {
    addErrorObject("expected to not be nil,\n               but got: nil\n")
  }
}

func (e *Expectation) toPanicWith(obj interface{}) {
  fn := e.value.(func())
  actual := rescueException(fn)
  
  if actual != obj {
    addErrorObject("expected panic: %v\n           got: %v\n", obj, actual)
  }
}

func (e *Expectation) toNotPanic() {
  fn := e.value.(func())
  actual := rescueException(fn)
  
  if actual != nil {
    addErrorObject("expected no panic,\n          but got: %v\n", actual)
  }
}

func rescueException(try func()) (out interface{}) {
  defer func() {
    out = recover()
  }()
  out = recover()
  try()
  return nil
}

func addErrorObject(s string, args ...interface{}) {
  s = fmt.Sprintf(s, args...)
  
  var contexts []string
  for _, testingContext := range testingContexts {
    contexts = append(contexts, testingContext.Description)
  }
  
  testingErrors = append(testingErrors, &testingError{
    String: s,
    ItString: testingCurrentIt,
    Contexts: contexts,
    ErrorLine: getErrorLine(),
  })
}



func BuildSpecReport() (string, bool) {
  var s string
  
  ok := len(testingErrors) == 0
  
  if !ok {
    s += redColor

    for _, error := range testingErrors {
      indents := 0

      for _, contextStr := range error.Contexts {
        s += fmt.Sprintf("%s- %s\n", strings.Repeat("  ", indents), contextStr)
        indents++
      }

      s += fmt.Sprintf("%s  %s\n\n", strings.Repeat("  ", indents), error.ItString)
      s += fmt.Sprintf("%s\n\t%s\n", error.String, error.ErrorLine)
    }

    s += resetColors
  } else {
    s += greenColor
    s += fmt.Sprintf("All tests passed. %d examples. 0 failures.\n", testingExamples)
    s += resetColors
  }
  
  return s, ok
}

func PrintSpecReport() {
  report, ok := BuildSpecReport()
  
  stream := specReportStream
  if stream == nil {
    stream = os.Stdout
  }
  fmt.Fprintf(stream, report)
  
  if !ok && !debugTesting {
    os.Exit(1)
  }
  
  testingContexts = testingContexts[0:0]
  testingErrors = testingErrors[0:0]
  testingExamples = 0
}
