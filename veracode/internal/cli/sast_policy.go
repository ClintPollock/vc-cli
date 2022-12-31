package cli

import (
	"context"

	"github.com/open-policy-agent/opa/rego"
)

const BASIC_SAST_POLICY = `
package sast.basic

import future.keywords

default allow := false

allow {
    print("Entering allow")
    print("input: ", input)
    #finding := input
    #print("Found finding with Id", finding)
    #input.Finding.Issue < 9

}

`

func newPolicyQuery(ctx context.Context, policy string) (rego.PreparedEvalQuery, error) {

	query, err := rego.New(
		rego.Query("x = data.sast.basic.allow"),
		rego.Module("policy.rego", policy),
	).PrepareForEval(ctx)

	if err != nil {
		// Handle error.
	}

	return query, err
}

//
// input := map[string]interface{}{
//     "method": "GET",
//     "path": []interface{}{"salary", "bob"},
//     "subject": map[string]interface{}{
//         "user": "bob",
//         "groups": []interface{}{"sales", "marketing"},
//     },
// }
//
// ctx := context.TODO()
// results, err := query.Eval(ctx, rego.EvalInput(input))
