package mcp

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"
)

func GetPrompt() mcp.Prompt {
	return mcp.NewPrompt("extract-relations-from-text",
		mcp.WithPromptDescription("Analyzes text to extract Person-Manager-Team relationships according to a specific ontology and inserts them as triples using InsertTriple."),
		mcp.WithArgument("input",
			mcp.ArgumentDescription("The text containing information about people, managers, teams, and their relationships."),
			mcp.RequiredArgument(),
		),
	)
}

func GetPromptHandler(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	arguments := request.Params.Arguments
	return &mcp.GetPromptResult{
		Description: "A complex prompt with arguments",
		Messages: []mcp.PromptMessage{
			{
				Role: mcp.RoleUser,
				Content: mcp.TextContent{
					Type: "text",
					Text: `Your task is to analyze the following text and extract relationships according to a specific ontology. 

**Ontology Rules:**
1.  Identify entities of type 'Person', 'Manager', and 'Team'.
2.  Extract relationships *only* if they match these patterns:
    *   Subject Type: Person, Predicate: 'worksFor', Object Type: Manager
    *   Subject Type: Person, Predicate: 'isMemberOf', Object Type: Team
    *   Subject Type: Team, Predicate: 'hasLeader', Object Type: Manager
3.  Represent entity names as accurately as possible from the text.

**Action Required:**
For *each* valid relationship you identify that strictly conforms to the ontology rules above, you **must** call the ` + "`InsertTriple`" + ` action. 
*   Use the identified entity names for ` + "`subject`" + ` and ` + "`object`" + `.
*   Use the corresponding predicate ('worksFor', 'isMemberOf', 'hasLeader').

Do not insert triples for relationships or entity types not explicitly mentioned in the ontology rules. Ensure the types of the subject and object match the rule for the predicate used.

Here is the text to analyze:`,
				},
			},
			{
				Role: mcp.RoleUser,
				Content: mcp.TextContent{
					Type: "text",
					Text: arguments["input"],
				},
			},
		},
	}, nil
}
