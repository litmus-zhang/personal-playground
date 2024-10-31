package scan_test

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/litmus-zhang/pScan/scan"
)

func TestAdd(t *testing.T) {
	testCases := []struct {
		name      string
		host      string
		expectLen int
		expectErr error
	}{
		{"AddNew", "host2", 2, nil},
		{"AddExisting", "host1", 2, scan.ErrExists},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hl := &scan.HostsList{}
			if err := hl.Add("host1"); err != nil {
				t.Fatal(err)
			}
			err := hl.Add(tc.host)
			if tc.expectErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tc.expectErr)
				}
				if !errors.Is(err, tc.expectErr) {
					t.Errorf("expected error %v, got %v", tc.expectErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(hl.Hosts) != tc.expectLen {
				t.Errorf("expected %d hosts, got %d", tc.expectLen, len(hl.Hosts))
			}
			if hl.Hosts[1] != tc.host {
				t.Errorf("expected %s hosts, got %s", tc.host, hl.Hosts[1])
			}
		})
	}
}

func TestRemove(t *testing.T) {
	testCases := []struct {
		name      string
		host      string
		expectLen int
		expectErr error
	}{
		{"RemoveExisiting", "host1", 1, nil},
		{"RemoveNotExisting", "host3", 1, scan.ErrNotExists},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hl := &scan.HostsList{}
			for _, h := range []string{"host1", "host2"} {
				if err := hl.Add(h); err != nil {
					t.Fatal(err)
				}
			}
			err := hl.Remove(tc.host)

			if tc.expectErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tc.expectErr)
				}
				if !errors.Is(err, tc.expectErr) {
					t.Errorf("expected error %v, got %v", tc.expectErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(hl.Hosts) != tc.expectLen {
				t.Errorf("expected %d hosts, got %d", tc.expectLen, len(hl.Hosts))
			}
			if hl.Hosts[0] == tc.host {
				t.Errorf("expected %s to be removed, but it's still there", tc.host)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	hl1 := &scan.HostsList{}
	hl2 := &scan.HostsList{}

	hostName := "host1"
	hl1.Add(hostName)
	tf, err := ioutil.TempFile("", "")

	if err != nil {
		t.Fatalf("unexpected error creating temp file: %v", err)
	}
	defer os.Remove(tf.Name())

	if err := hl1.Save(tf.Name()); err != nil {
		t.Fatalf("unexpected error saving hosts to file: %v", err)
	}
	if err := hl2.Load((tf.Name())); err != nil {
		t.Fatalf("unexpected error loading hosts from file: %v", err)
	}
	if hl1.Hosts[0] != hl2.Hosts[0] {
		t.Errorf("expected %s, got %s", hl1.Hosts[0], hl2.Hosts[0])
	}
}

func TestLoadNoFile(t *testing.T) {
	tf, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("Error creating temp file: %v", err)
	}
	if err := os.Remove(tf.Name()); err != nil {
		t.Fatalf("Error removing temp file: %v", err)
	}
	hl := &scan.HostsList{}
	if err := hl.Load(tf.Name()); err != nil {
		t.Fatalf("unexpected error loading hosts from file: %v", err)
	}
}
