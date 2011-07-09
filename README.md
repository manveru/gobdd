# BDD testing for Go

It's pretty cool.

## Installation

	make test # assurance that it works
	make install # installs to your $GOROOT

## Usage

Inside a normal gotest-style file (ie, *_test.go) add the following:

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
  
	  var anObject *MyGreatTestType
  
	  Describe("matchers", func() {
    
	    BeforeEach(func() {
	      // this is run at this level and every describe within it, however deeply nested
	      anObject = new(MyGreatTestType)
	      anObject.Name = "john"
	      anObject.Age = 23
	    })
    
	    Describe("not equals", func() {
      
	      It("matches on simple objects", func() {
	        Expect(&MyGreatTestType{"john", 23}, ToNotEqual, anObject)
	        Expect("foo", ToEqual, "foo")
	        Expect("foo", ToNotEqual, "bar")
	      })
      
	      It("matches for typed-nil", func() {
	        Expect(MyNil(), ToBeNil)
	        Expect(MyNonNil(), ToNotBeNil)
	      })
      
	      It("matches for nil", func() {
	        // Expect(nil, ToBeNil)
	        Expect(true, ToNotBeNil)
	      })
      
	    })
    
	    Describe("deep equals matcher", func() {
      
	      It("matches what equals does not", func() {
	        Expect(&MyGreatTestType{"john", 23}, ToDeepEqual, anObject)
	        Expect("foo", ToDeepEqual, "foo")
	      })
      
	    })
    
	    Describe("exception-rescuing matchers", func() {
      
	      It("is super cool", func() {
	        Expect(func() { panic("foobar!") }, ToPanicWith, "foobar!")
	        Expect(func() {}, ToNotPanic)
	      })
      
	    })
    
	  })
	}
