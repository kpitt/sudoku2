package board

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoard_IsValid(t *testing.T) {
	t.Run("valid board", func(t *testing.T) {
		board := Board{
			{5, 3, 0, 0, 7, 0, 0, 0, 0},
			{6, 0, 0, 1, 9, 5, 0, 0, 0},
			{0, 9, 8, 0, 0, 0, 0, 6, 0},
			{8, 0, 0, 0, 6, 0, 0, 0, 3},
			{4, 0, 0, 8, 0, 3, 0, 0, 1},
			{7, 0, 0, 0, 2, 0, 0, 0, 6},
			{0, 6, 0, 0, 0, 0, 2, 8, 0},
			{0, 0, 0, 4, 1, 9, 0, 0, 5},
			{0, 0, 0, 0, 8, 0, 0, 7, 9},
		}
		assert.True(t, board.IsValid())
	})

	t.Run("duplicate in row", func(t *testing.T) {
		board := Board{
			{5, 5, 0, 0, 7, 0, 0, 0, 0},
			{0, 0, 0, 0, 0, 0, 0, 0, 0},
		}
		assert.False(t, board.IsValid())
	})

	t.Run("duplicate in column", func(t *testing.T) {
		board := Board{
			{5, 0, 0, 0, 0, 0, 0, 0, 0},
			{5, 0, 0, 0, 0, 0, 0, 0, 0},
		}
		assert.False(t, board.IsValid())
	})

	t.Run("duplicate in box", func(t *testing.T) {
		board := Board{
			{5, 0, 0, 0, 0, 0, 0, 0, 0},
			{0, 5, 0, 0, 0, 0, 0, 0, 0},
		}
		assert.False(t, board.IsValid())
	})
}

func TestBoard_Getters(t *testing.T) {
	board := Board{
		{1, 2, 3, 4, 5, 6, 7, 8, 9},
		{4, 5, 6, 7, 8, 9, 1, 2, 3},
		{7, 8, 9, 1, 2, 3, 4, 5, 6},
		{2, 3, 1, 5, 6, 4, 8, 9, 7},
		{5, 6, 4, 8, 9, 7, 2, 3, 1},
		{8, 9, 7, 2, 3, 1, 5, 6, 4},
		{3, 1, 2, 6, 4, 5, 9, 7, 8},
		{6, 4, 5, 9, 7, 8, 3, 1, 2},
		{9, 7, 8, 3, 1, 2, 6, 4, 5},
	}

	t.Run("GetRow", func(t *testing.T) {
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, board.GetRow(0))
	})

	t.Run("GetCol", func(t *testing.T) {
		assert.Equal(t, []int{1, 4, 7, 2, 5, 8, 3, 6, 9}, board.GetCol(0))
	})

	t.Run("GetBox", func(t *testing.T) {
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, board.GetBox(0, 0))
	})
}

func TestBoard_IsSolved(t *testing.T) {
	t.Run("solved board", func(t *testing.T) {
		board := Board{
			{1, 2, 3, 4, 5, 6, 7, 8, 9},
			{4, 5, 6, 7, 8, 9, 1, 2, 3},
			{7, 8, 9, 1, 2, 3, 4, 5, 6},
			{2, 3, 1, 5, 6, 4, 8, 9, 7},
			{5, 6, 4, 8, 9, 7, 2, 3, 1},
			{8, 9, 7, 2, 3, 1, 5, 6, 4},
			{3, 1, 2, 6, 4, 5, 9, 7, 8},
			{6, 4, 5, 9, 7, 8, 3, 1, 2},
			{9, 7, 8, 3, 1, 2, 6, 4, 5},
		}
		assert.True(t, board.IsSolved())
	})

	t.Run("unsolved valid board", func(t *testing.T) {
		board := Board{
			{1, 2, 3, 4, 5, 6, 7, 8, 0},
			{4, 5, 6, 7, 8, 9, 1, 2, 3},
			{7, 8, 9, 1, 2, 3, 4, 5, 6},
			{2, 3, 1, 5, 6, 4, 8, 9, 7},
			{5, 6, 4, 8, 9, 7, 2, 3, 1},
			{8, 9, 7, 2, 3, 1, 5, 6, 4},
			{3, 1, 2, 6, 4, 5, 9, 7, 8},
			{6, 4, 5, 9, 7, 8, 3, 1, 2},
			{9, 7, 8, 3, 1, 2, 6, 4, 5},
		}
		assert.False(t, board.IsSolved())
	})

	t.Run("invalid board", func(t *testing.T) {
		board := Board{
			{1, 1, 3, 4, 5, 6, 7, 8, 9},
			{4, 5, 6, 7, 8, 9, 1, 2, 3},
			{7, 8, 9, 1, 2, 3, 4, 5, 6},
			{2, 3, 1, 5, 6, 4, 8, 9, 7},
			{5, 6, 4, 8, 9, 7, 2, 3, 1},
			{8, 9, 7, 2, 3, 1, 5, 6, 4},
			{3, 1, 2, 6, 4, 5, 9, 7, 8},
			{6, 4, 5, 9, 7, 8, 3, 1, 2},
			{9, 7, 8, 3, 1, 2, 6, 4, 5},
		}
		assert.False(t, board.IsSolved())
	})
}

func TestBoard_IsSafe(t *testing.T) {
	board := Board{
		{5, 3, 0, 0, 7, 0, 0, 0, 0},
		{6, 0, 0, 1, 9, 5, 0, 0, 0},
		{0, 9, 8, 0, 0, 0, 0, 6, 0},
		{8, 0, 0, 0, 6, 0, 0, 0, 3},
		{4, 0, 0, 8, 0, 3, 0, 0, 1},
		{7, 0, 0, 0, 2, 0, 0, 0, 6},
		{0, 6, 0, 0, 0, 0, 2, 8, 0},
		{0, 0, 0, 4, 1, 9, 0, 0, 5},
		{0, 0, 0, 0, 8, 0, 0, 7, 9},
	}

	t.Run("safe placement", func(t *testing.T) {
		assert.True(t, board.IsSafe(0, 2, 1))
	})

	t.Run("unsafe in row", func(t *testing.T) {
		assert.False(t, board.IsSafe(0, 2, 3))
	})

	t.Run("unsafe in col", func(t *testing.T) {
		assert.False(t, board.IsSafe(0, 2, 8))
	})

	t.Run("unsafe in box", func(t *testing.T) {
		assert.False(t, board.IsSafe(0, 2, 9))
	})
}
