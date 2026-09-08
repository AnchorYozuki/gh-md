package writer

import (
	"strings"
	"testing"
	"time"

	"github.com/jackchuka/gh-md/internal/github"
)

func TestIssueCommentMetadataIncludesUpdatedAndURL(t *testing.T) {
	created := time.Date(2026, 9, 8, 2, 10, 16, 0, time.UTC)
	updated := time.Date(2026, 9, 8, 2, 10, 59, 0, time.UTC)
	commentURL := "https://github.com/owner/repo/issues/21#issuecomment-123"

	issue := &github.Issue{
		ID:        "I_test",
		URL:       "https://github.com/owner/repo/issues/21",
		Number:    21,
		Owner:     "owner",
		Repo:      "repo",
		Title:     "Comment metadata",
		Body:      "body",
		State:     "open",
		Author:    "author",
		CreatedAt: created,
		UpdatedAt: updated,
		Comments: []github.Comment{
			{
				ID:        "IC_test",
				URL:       commentURL,
				Author:    "commenter",
				Body:      "edited comment",
				CreatedAt: created,
				UpdatedAt: updated,
			},
		},
	}

	got, err := IssueToMarkdown(issue)
	if err != nil {
		t.Fatalf("IssueToMarkdown() error = %v", err)
	}

	for _, want := range []string{
		"id: IC_test",
		"created: 2026-09-08T02:10:16Z",
		"updated: 2026-09-08T02:10:59Z",
		"url: " + commentURL,
		"edited comment",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("IssueToMarkdown() missing %q\nGot:\n%s", want, got)
		}
	}
}

func TestCommentMetadataOmitsUnavailableOptionalFields(t *testing.T) {
	created := time.Date(2026, 9, 8, 2, 10, 16, 0, time.UTC)

	issue := &github.Issue{
		ID:        "I_test",
		URL:       "https://github.com/owner/repo/issues/21",
		Number:    21,
		Owner:     "owner",
		Repo:      "repo",
		Title:     "Comment metadata",
		Body:      "body",
		State:     "open",
		CreatedAt: created,
		UpdatedAt: created,
		Comments: []github.Comment{
			{
				ID:        "IC_test",
				Author:    "commenter",
				Body:      "comment",
				CreatedAt: created,
			},
		},
	}

	got, err := IssueToMarkdown(issue)
	if err != nil {
		t.Fatalf("IssueToMarkdown() error = %v", err)
	}

	commentStart := strings.Index(got, "<!-- gh-md:comment")
	commentEnd := strings.Index(got[commentStart:], "-->")
	if commentStart == -1 || commentEnd == -1 {
		t.Fatalf("comment metadata block not found\nGot:\n%s", got)
	}
	metadata := got[commentStart : commentStart+commentEnd]

	if strings.Contains(metadata, "updated:") {
		t.Errorf("comment metadata unexpectedly contains zero updated time: %s", metadata)
	}
	if strings.Contains(metadata, "url:") {
		t.Errorf("comment metadata unexpectedly contains empty URL: %s", metadata)
	}
}
