package lmsp

import (
	"math"
	"os"
	"testing"
	"time"
)

func TestReadingFile(t *testing.T) {
	run := func(filename, desc string, tf func(*testing.T, *Reader)) {
		t.Run(filename+": "+desc, func(t *testing.T) {
			t.Parallel()

			f, err := os.Open(filename)
			if err != nil {
				t.Errorf("error opening file: %v", err)
				return
			}
			defer f.Close()

			r, err := ReadFile(f)
			if err != nil {
				t.Errorf("error reading file: %v", err)
			}

			tf(t, r)
		})
	}

	//t.Run("no file", func(t *testing.T) {
	//	_, err := Read("does-not-exist.lmsp")
	//	if !os.IsNotExist(err) {
	//		t.Errorf("expected IsNotExist error but got %v", err)
	//	}
	//})

	run("testdata/my-block-with-no-refs.lmsp", "manifest", func(t *testing.T, r *Reader) {
		compareManifest(t, r, Manifest{
			Type:          "word-blocks",
			AutoDelete:    false,
			Created:       time.Date(2021, time.January, 1, 21, 36, 31, 629000000, time.UTC),
			ID:            "UIcJW0lhqJhf",
			LastSaved:     time.Date(2021, time.May, 28, 12, 48, 9, 124000000, time.UTC),
			Size:          0,
			Name:          "oh yessssss",
			SlotIndex:     0,
			WorkspaceX:    346,
			WorkspaceY:    329,
			ZoomLevel:     0.675,
			ShowAllBlocks: false,
			Version:       5,
			//Hardware:      {},
		})
	})

	run("testdata/Gyro drive.lmsp", "manifest", func(t *testing.T, r *Reader) {
		compareManifest(t, r, Manifest{
			Type:          "word-blocks",
			AutoDelete:    false,
			Created:       time.Date(2020, time.September, 5, 16, 21, 51, 720000000, time.UTC),
			ID:            "nsnxhdh9RxVt",
			LastSaved:     time.Date(2021, time.May, 28, 12, 49, 7, 638000000, time.UTC),
			Size:          0,
			Name:          "Gyro drive",
			SlotIndex:     0,
			WorkspaceX:    120,
			WorkspaceY:    212.425,
			ZoomLevel:     0.675,
			ShowAllBlocks: false,
			Version:       5,
			//Hardware:      {},
		})
	})
}

func compareManifest(t *testing.T, r *Reader, expected Manifest) {
	actual, err := r.Manifest()
	if err != nil {
		t.Error(err)
		return
	}
	if expected.Type != actual.Type {
		t.Errorf("manifest.Type: expected %v but got %v", expected.Type, actual.Type)
	}
	if expected.AutoDelete != actual.AutoDelete {
		t.Errorf("manifest.AutoDelete: expected %v but got %v", expected.AutoDelete, actual.AutoDelete)
	}
	if expected.Created != actual.Created {
		t.Errorf("manifest.Created: expected %v but got %v", expected.Created, actual.Created)
	}
	if expected.ID != actual.ID {
		t.Errorf("manifest.ID: expected %v but got %v", expected.ID, actual.ID)
	}
	if expected.LastSaved != actual.LastSaved {
		t.Errorf("manifest.LastSaved: expected %v but got %v", expected.LastSaved, actual.LastSaved)
	}
	if expected.Size != actual.Size {
		t.Errorf("manifest.Size: expected %v but got %v", expected.Size, actual.Size)
	}
	if expected.Name != actual.Name {
		t.Errorf("manifest.Name: expected %q but got %q", expected.Name, actual.Name)
	}
	if expected.SlotIndex != actual.SlotIndex {
		t.Errorf("manifest.SlotIndex: expected %v but got %v", expected.SlotIndex, actual.SlotIndex)
	}
	if !closeEnoughFloats(expected.WorkspaceX, actual.WorkspaceX, 0.0001) {
		t.Errorf("manifest.WorkspaceX: expected %v but got %v", expected.WorkspaceX, actual.WorkspaceX)
	}
	if !closeEnoughFloats(expected.WorkspaceY, actual.WorkspaceY, 0.0001) {
		t.Errorf("manifest.WorkspaceY: expected %v but got %v", expected.WorkspaceY, actual.WorkspaceY)
	}
	if expected.ZoomLevel != actual.ZoomLevel {
		t.Errorf("manifest.ZoomLevel: expected %v but got %v", expected.ZoomLevel, actual.ZoomLevel)
	}
	if expected.ShowAllBlocks != actual.ShowAllBlocks {
		t.Errorf("manifest.ShowAllBlocks: expected %v but got %v", expected.ShowAllBlocks, actual.ShowAllBlocks)
	}
	if expected.Version != actual.Version {
		t.Errorf("manifest.Version: expected %v but got %v", expected.Version, actual.Version)
	}
}

func closeEnoughFloats(a, b, epsilon float64) bool {
	return math.Abs(a-b) < epsilon
}
