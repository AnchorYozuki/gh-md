package parser

import "testing"

func TestParseCommentIgnoresExtendedMetadata(t *testing.T) {
	content := `---
id: I_test
url: https://github.com/owner/repo/issues/21
number: 21
owner: owner
repo: repo
title: Test
state: open
created: 2026-09-08T02:00:00Z
updated: 2026-09-08T02:20:00Z
last_pulled: 2026-09-08T02:30:00Z
---

<!-- gh-md:content -->
# Test

body
<!-- /gh-md:content -->

---

<!-- gh-md:comment
id: IC_test
author: commenter
created: 2026-09-08T02:10:16Z
updated: 2026-09-08T02:10:59Z
url: https://github.com/owner/repo/issues/21#issuecomment-123
-->
### @commenter (2026-09-08)

edited comment
<!-- /gh-md:comment -->

<!-- gh-md:new-comment -->

<!-- /gh-md:new-comment -->
`

	parsed, err := parseContent(content, "/tmp/issues/21.md")
	if err != nil {
		t.Fatalf("parseContent() error = %v", err)
	}
	if len(parsed.Comments) != 1 {
		t.Fatalf("got %d comments, want 1", len(parsed.Comments))
	}

	comment := parsed.Comments[0]
	if comment.ID != "IC_test" {
		t.Errorf("comment ID = %q, want IC_test", comment.ID)
	}
	if comment.Author != "commenter" {
		t.Errorf("comment Author = %q, want commenter", comment.Author)
	}
	if comment.Body != "edited comment" {
		t.Errorf("comment Body = %q, want edited comment", comment.Body)
	}
}
