package partyservice

import (
	"strings"
	"testing"
	"time"
)

var (
	begin = time.Date(2020, 7, 31, 2, 30, 0, 0, time.UTC)
	end   = time.Date(2020, 7, 31, 4, 30, 0, 0, time.UTC)
)

func TestBuildSQLForQueryPartiesWhereTimeRange(t *testing.T) {
	assert := func(t *testing.T, input *PartyQueryDto, expSQL string, expArgLen int) {
		// Act
		actual, args := buildSQLForQueryPartiesWhereTimeRange(input)
		// Assert
		expSQL = strings.TrimSpace(expSQL)
		expSQL = strings.ReplaceAll(expSQL, "\n", " ")
		if actual != expSQL {
			t.Errorf("\nexpected:\n %s, \ngot:\n %s", expSQL, actual)
		}
		if len(args) != expArgLen {
			t.Errorf("\nexpected:\n %d, \ngot:\n %d", expArgLen, len(args))
		}
	}

	t.Run("begin and end query dto", func(t *testing.T) {
		var (
			expSQL = `
SELECT id, startFrom, endTo, chatRoomId
FROM parties
WHERE startFrom >= ? AND endTo <= ?
`
			expArgLen = 2
		)
		input := &PartyQueryDto{
			beginDateTime: &begin,
			endDateTime:   &end,
		}
		assert(t, input, expSQL, expArgLen)
	})

	t.Run("begin query dto", func(t *testing.T) {
		var (
			expSQL = `
SELECT id, startFrom, endTo, chatRoomId
FROM parties
WHERE startFrom >= ?
`
			expArgLen = 1
		)
		input := &PartyQueryDto{
			beginDateTime: &begin,
		}
		assert(t, input, expSQL, expArgLen)
	})

	t.Run("empty query dto", func(t *testing.T) {
		var (
			expSQL = `
SELECT id, startFrom, endTo, chatRoomId
FROM parties
`
			expArgLen = 0
		)
		input := &PartyQueryDto{}
		assert(t, input, expSQL, expArgLen)
	})
}

func TestBuildSQLForQueryPartiesWhereUserIdAndTimeRange(t *testing.T) {
	assert := func(t *testing.T, input *PartyQueryDto, expSQL string, expArgLen int) {
		userId := "user-id"
		// Act
		actual, args := buildSQLForQueryPartiesWhereUserIdAndTimeRange(userId, input)
		// Assert
		expSQL = strings.TrimSpace(expSQL)
		expSQL = strings.ReplaceAll(expSQL, "\n", " ")
		if actual != expSQL {
			t.Errorf("\nexpected:\n %s, \ngot:\n %s", expSQL, actual)
		}
		if len(args) != expArgLen {
			t.Errorf("\nexpected:\n %d, \ngot:\n %d", expArgLen, len(args))
		}
	}

	t.Run("begin and end query dto", func(t *testing.T) {
		var (
			expSQL = `
SELECT p.id AS id, p.startFrom AS startFrom, p.endTo AS endTo, p.chatRoomId AS chatRoomId
FROM partymembers AS pm
JOIN parties AS p ON pm.partyId = p.id
WHERE pm.userId = ? AND p.startFrom >= ? AND p.endTo <= ?
`
			expArgLen = 3
		)
		input := &PartyQueryDto{
			beginDateTime: &begin,
			endDateTime:   &end,
		}
		assert(t, input, expSQL, expArgLen)
	})

	t.Run("begin query dto", func(t *testing.T) {
		var (
			expSQL = `
SELECT p.id AS id, p.startFrom AS startFrom, p.endTo AS endTo, p.chatRoomId AS chatRoomId
FROM partymembers AS pm
JOIN parties AS p ON pm.partyId = p.id
WHERE pm.userId = ? AND p.startFrom >= ?
`
			expArgLen = 2
		)
		input := &PartyQueryDto{
			beginDateTime: &begin,
		}
		assert(t, input, expSQL, expArgLen)
	})

	t.Run("empty query dto", func(t *testing.T) {
		var (
			expSQL = `
SELECT p.id AS id, p.startFrom AS startFrom, p.endTo AS endTo, p.chatRoomId AS chatRoomId
FROM partymembers AS pm
JOIN parties AS p ON pm.partyId = p.id
WHERE pm.userId = ?
`
			expArgLen = 1
		)
		input := &PartyQueryDto{}
		assert(t, input, expSQL, expArgLen)
	})
}

func TestBuildSQLForQueryPartyReviewMembers(t *testing.T) {
	var (
		partyId  int64 = 1
		reviewer       = "user-id-1"
		reviewee       = "user-id-2"
	)

	assert := func(t *testing.T, input *ReviewMemberQueryDto, expSQL string, expArgLen int) {
		// Act
		actual, args := buildSQLForQueryPartyReviewMembers(input)
		// Assert
		expSQL = strings.TrimSpace(expSQL)
		expSQL = strings.ReplaceAll(expSQL, "\n", " ")
		if actual != expSQL {
			t.Errorf("\nexpected:\n %s, \ngot:\n %s", expSQL, actual)
		}
		if len(args) != expArgLen {
			t.Errorf("\nexpected:\n %d, \ngot:\n %d", expArgLen, len(args))
		}
	}

	t.Run("full where", func(t *testing.T) {
		var (
			expSQL = `
SELECT partyId, reviewer, reviewee, score, comments
FROM partymemberreviews
WHERE partyId = ? AND reviewer = ? AND reviewee = ?
`
			expArgLen = 3
		)
		input := &ReviewMemberQueryDto{
			partyID:  partyId,
			reviewer: reviewer,
			reviewee: reviewee,
		}
		assert(t, input, expSQL, expArgLen)
	})

	t.Run("one where", func(t *testing.T) {
		var (
			expSQL = `
SELECT partyId, reviewer, reviewee, score, comments
FROM partymemberreviews
WHERE reviewer = ?
`
			expArgLen = 1
		)
		input := &ReviewMemberQueryDto{
			reviewer: reviewer,
		}
		assert(t, input, expSQL, expArgLen)
	})

	t.Run("no where", func(t *testing.T) {
		var (
			expSQL = `
SELECT partyId, reviewer, reviewee, score, comments
FROM partymemberreviews
`
			expArgLen = 0
		)
		input := &ReviewMemberQueryDto{}
		assert(t, input, expSQL, expArgLen)
	})
}
