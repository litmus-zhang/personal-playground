package calculator_test

import (
	"log"
	"os"
	"testing"

	"github.com/litmus-zhang/go-tdd/ch-2/calculator"
)

func TestMain(m *testing.M) {
	// setup
	setup()
	e := m.Run()
	// teardown
	teardown()
	os.Exit(e)
}
func setup() {
	log.Println("Before all tests")
}
func teardown() {
	log.Println("Tearing down after all tests")
}
func init() {
	log.Println("Initial setup before all tests")
}
func TestAdd(t *testing.T) {
	e := calculator.Engine{}

	actAssert := func(x, y, want float64) {
		got := e.Add(x, y)
		if got != want {
			t.Errorf("got %.2f want %.2f", got, want)
		}

	}
	t.Run("Positive input", func(t *testing.T) {
		x, y := 2.5, 3.5
		actAssert(x, y, 6.0)
	})
	t.Run("Negative input", func(t *testing.T) {
		x, y := -2.5, -3.5
		actAssert(x, y, -6.0)
	})
}

func BenchmarkAdd(b *testing.B) {
	e := calculator.Engine{}
	for i := 0; i < b.N; i++ {
		e.Add(2.5, 3.5)
	}
}

func TestSubtract(t *testing.T) {
	e := calculator.Engine{}

	actAssert := func(x, y, want float64) {
		got := e.Subtract(x, y)
		if got != want {
			t.Errorf("got %.2f want %.2f", got, want)
		}

	}
	t.Run("Positive input", func(t *testing.T) {
		x, y := 3.5, 2.5
		actAssert(x, y, 1.0)
	})
	t.Run("Negative input", func(t *testing.T) {
		x, y := -3.5, -2.5
		actAssert(x, y, -1.0)
	})
}
func TestMultiply(t *testing.T) {
	e := calculator.Engine{}

	actAssert := func(x, y, want float64) {
		got := e.Multiply(x, y)
		if got != want {
			t.Errorf("got %.2f want %.2f", got, want)
		}

	}
	t.Run("Positive input", func(t *testing.T) {
		x, y := 3.0, 2.5
		actAssert(x, y, 7.5)
	})
	t.Run("Negative input", func(t *testing.T) {
		x, y := -3.0, -2.0
		actAssert(x, y, 6.0)
	})
}

func TestDivide(t *testing.T) {
	e := calculator.Engine{}
	actAssert := func(x, y, want float64) {
		got, err := e.Divide(x, y)
		if err != nil {
			t.Error(err)
		}
		if got != want {
			t.Errorf("got %.2f want %.2f", got, want)
		}
	}
	t.Run("Positive input", func(t *testing.T) {
		x, y := 3.0, 2.0
		actAssert(x, y, 1.5)
	})
	t.Run("Negative input", func(t *testing.T) {
		x, y := -3.0, -2.0
		actAssert(x, y, 1.5)
	})
	t.Run("Divide by zero", func(t *testing.T) {
		x, y := 3.0, 0.0
		_, err := e.Divide(x, y)
		if err == nil {
			t.Error("Expected error")
		}
	})

}
