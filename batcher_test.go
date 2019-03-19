package batcher

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGet_Basic(t *testing.T) {
	const (
		expected1 = "foo"
		expected2 = "bar"
		expected3 = "baz"
	)
	ctx := context.Background()
	called := false

	b := New(3*time.Second, func(params []interface{}) ([]interface{}, error) {
		if called {
			// this should only be called once
			t.Fatal("already called bulk fetch")
		}
		called = true

		results := make([]interface{}, len(params))
		for i, p := range params {
			if p != "notfound" {
				results[i] = p
			}
		}

		return results, nil
	})

	wg := sync.WaitGroup{}
	wg.Add(4)

	go func() {
		defer wg.Done()

		time.Sleep(2 * time.Second)
		v, err := b.Get(ctx, expected1)
		assert.NoError(t, err)
		assert.Equal(t, expected1, v)
	}()
	go func() {
		defer wg.Done()

		time.Sleep(1 * time.Second)
		v, err := b.Get(ctx, expected2)
		assert.NoError(t, err)
		assert.Equal(t, expected2, v)
	}()
	go func() {
		defer wg.Done()

		v, err := b.Get(ctx, expected3)
		assert.NoError(t, err)
		assert.Equal(t, expected3, v)
	}()
	go func() {
		defer wg.Done()

		v, err := b.Get(ctx, "notafound")
		assert.NoError(t, err)
		assert.Nil(t, v)
	}()

	wg.Wait()
}

func TestGet_Error(t *testing.T) {
	ctx := context.Background()
	called := false
	expectedErr := errors.New("expected error")

	b := New(1*time.Second, func([]interface{}) ([]interface{}, error) {
		if called {
			// this should only be called once
			t.Fatal("already called bulk fetch")
		}
		called = true
		return nil, expectedErr
	})

	_, err := b.Get(ctx, "doesn't matter")
	assert.EqualError(t, err, expectedErr.Error())
}

func TestGet_Context(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	expectedErr := "context canceled"

	b := New(1*time.Second, func([]interface{}) ([]interface{}, error) {
		t.Fatal("this should not be called")
		panic("shouldn't get here")
	})

	cancel()
	_, err := b.Get(ctx, "doesn't matter")
	assert.EqualError(t, err, expectedErr)
}
