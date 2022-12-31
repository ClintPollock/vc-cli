# conftest policy to evaluate results for issues
package main


deny[msg] {
  input.findings.docker.details[i].level == "FATAL"

  msg = sprintf(
      "Found FATAL docker image issue: %s: %s",
      [ input.findings.docker.details[i].code,
        input.findings.docker.details[i].alerts[0] ] )
}
