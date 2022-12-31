package intrem

type Rule struct {
  ID string                 `json:"id"`
  Name string               `json:"name"`
  ShortDescription struct {
    Text string             `json:"text"`
  }                         `json:"shortDescription"`
  HelpURI string            `json:"helpUri"`
  Properties struct {
    Category string         `json:"category"`
    Tags []string           `json:"tags"`
  }                         `json:"properties"`
  DefaultConfiguration struct {
    Level string            `json:"level"`
  }                         `json:"defaultConfiguration"`
}

type LogicalLocation struct {
  Name               string`json:"name"`
  FullyQualifiedName string`json:"fullyQualifiedName"`
  Kind               string`json:"function"`
}

type Location struct {
  PhysicalLocation struct {
    ArtifactLocation struct {
      URI string            `json:"uri"`
    }                       `json:"artifactLocation"`
    Region struct {
      StartLine int         `json:"startLine"`
      EndLine int           `json:"endLine"`
    }                       `json:"region"`
  }                         `json:"physicalLocation"`
  LogicalLocations []LogicalLocation  `json:"logicalLocations"`
}

type Result struct {
  Level   string              `json:"level"`
  Rank    int                 `json:"rank"`
  Message struct {
    Text string               `json:"text"`
  }                           `json:"message"`
  Locations []Location         `json:"locations"`
  RuleID int                   `json:"ruleId"`
  Fingerprints struct {
    FlawHash string            `json:"flawHash"`
    CauseHash string           `json:"causeHash"`
    ProcedureHash string       `json:"procedureHash"`
    PrototypeHash string       `json:"prototypeHash"`
  }                            `json:"fingerprints"`
}

type Run struct {
  Tool struct {
    Driver struct {
      Name   string             `json:"name"`
      Rules  []*Rule            `json:"rules"`
    }                           `json:"driver"`
  }                             `json:"tool"`
  Results []*Result             `json:"results"`
}

type SARIFSimple struct {
  Schema       string             `json:"$schema"`
  Version      string             `json:"version"`
  Runs []Run                      `json:"runs"`
}
