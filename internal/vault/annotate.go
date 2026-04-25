package vault

import (
	"fmt"
	"time"
)

// Annotation holds a user-defined note attached to a secret path.
type Annotation struct {
	Path      string    `json:"path"`
	Note      string    `json:"note"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
}

// AnnotationIndex maps secret paths to their annotations.
type AnnotationIndex struct {
	Annotations map[string]Annotation `json:"annotations"`
}

// NewAnnotationIndex creates an empty AnnotationIndex.
func NewAnnotationIndex() *AnnotationIndex {
	return &AnnotationIndex{
		Annotations: make(map[string]Annotation),
	}
}

// Add attaches or replaces an annotation for the given path.
func (idx *AnnotationIndex) Add(path, note, author string) {
	idx.Annotations[path] = Annotation{
		Path:      path,
		Note:      note,
		Author:    author,
		CreatedAt: time.Now().UTC(),
	}
}

// Remove deletes the annotation for the given path, if present.
func (idx *AnnotationIndex) Remove(path string) bool {
	if _, ok := idx.Annotations[path]; !ok {
		return false
	}
	delete(idx.Annotations, path)
	return true
}

// Get retrieves the annotation for the given path.
func (idx *AnnotationIndex) Get(path string) (Annotation, bool) {
	a, ok := idx.Annotations[path]
	return a, ok
}

// Summary returns a human-readable overview of all annotations.
func (idx *AnnotationIndex) Summary() string {
	if len(idx.Annotations) == 0 {
		return "no annotations"
	}
	return fmt.Sprintf("%d annotation(s) recorded", len(idx.Annotations))
}
