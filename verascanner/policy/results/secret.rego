# conftest policy to evaluate results for issues
package main


deny[msg] {

  severities := ["HIGH", "CRITICAL"]
  input.findings.secrets.Results[i].Secrets[j].Severity == severities[_]


  msg = sprintf(
      "Found %s secret: %s: %s",
      [ input.findings.secrets.Results[i].Secrets[j].Severity,
        input.findings.secrets.Results[i].Target,
        input.findings.secrets.Results[i].Secrets[j].Title
      ] )
}
