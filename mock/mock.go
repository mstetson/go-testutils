/*
Package mock contains simple utilities for testing with mocks.

A Mock can be installed and restored.
This package contains types and functions
for creating and using such mocks.

The main idea is that instead of this

	m1 := &mockWhatever{
		Foo: "bar",
	}
	m1.Install()
	defer m1.Restore()

	m2 := &mockAnother{
		Baz: "Quux",
	}
	m2.Install()
	defer m2.Restore()

we can do this

	mocks := mock.Group{
		&mockWhatever{
			Foo: "bar",
		},
		&mockAnother{
			Baz: "Quux",
		},
	}
	mocks.Install()
	defer mocks.Restore()

or

	func TestSomething(t *testing.T) {
		mock.UntilCleanup(t,
			&mockWhatever{
				Foo: "bar",
			},
			&mockAnother{
				Baz: "Quux",
			},
		)

		...
	}
*/
package mock

// A Mock can be installed and restored.
type Mock interface {
	Install() // Install sets up the Mock for use in tests.
	Restore() // Restore undoes changes made by Install.
}

// Install installs each of the given Mocks.
func Install(ms ...Mock) {
	for _, m := range ms {
		m.Install()
	}
}

// Restore restores each of the given Mocks.
// They are restored in reverse order,
// in case they have ordering dependencies.
func Restore(ms ...Mock) {
	for i := len(ms) - 1; i >= 0; i-- {
		ms[i].Restore()
	}
}

// A Group allows a collection of Mocks to be treated as one.
type Group []Mock

// Install installs each mock in g.
func (g Group) Install() {
	Install(g...)
}

// Restore restores each mock in g.
func (g Group) Restore() {
	Restore(g...)
}

// Set returns a Mock that sets *mockable to mock at Install,
// saving *mockable's original value.
// On Restore, *mockable is reset to its original value.
func Set[T any](mockable *T, mock T) Mock {
	return &setMock[T]{mockable: mockable, mock: mock}
}

type setMock[T any] struct {
	mockable *T
	mock     T
	orig     T
}

func (m *setMock[T]) Install() {
	m.orig = *m.mockable
	*m.mockable = m.mock
}

func (m *setMock[T]) Restore() {
	*m.mockable = m.orig
}

// A Cleanupper can be asked to call a cleanup function.
// The most common Cleanupper is *testing.T.
type Cleanupper interface {
	Cleanup(func())
}

// UntilCleanup installs the given mocks and tells t to restore them on Cleanup.
// Usually, t is a *testing.T.
func UntilCleanup(t Cleanupper, ms ...Mock) {
	Group(ms).Install()
	t.Cleanup(Group(ms).Restore)
}
