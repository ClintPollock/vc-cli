package intrem

// This type represents the information needed to request an auto-generated
// fix for a flaw. Here we refer to "flaws" explicitly, since we are
// seeking a fix.

// Based on OpenAPI spec:
/*
Flaw:
  type: object
  properties:
    sourceFile:
      type: string
    function:
      type: string
    line:
      type: integer
    CWEId:
      type: string
    flow:
      type: array
      items:
        type: object
        properties:
          line:
            type: integer
          region:            # Expression region, if known
            type: object
            properties:
              startLine:
                type: integer
              endLine:
                type: integer
              startColumn:
                type: integer
              endColumn:
                type: integer
          expression:
            type: string
          expressionType:    # Compile-time class, if known
            type: string
*/

type FlawToFix struct {
  SourceFile string               `json:"sourceFile"`
  Function string                 `json:"function"`
  Line int                        `json:"line"`
  CWEId string                    `json:"CWEId"`
  Flow []Step                     `json:"flow"`
}

type Step struct {
  Line int                        `json:"line"`
  Region Region                   `json:"region"`
  Expression string               `json:"expression"`
  ExpressionType string           `json:"expressionType"`
}

type Region struct {
  StartLine int                 `json:"startLine"`
  EndLine int                   `json:"endLine"`
  StartColumn int               `json:"startColumn"`
  EndColumn int                 `json:"endColumn"`
}
