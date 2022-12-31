# conftest policy to test docker image for CIS Docker Benchmark compliance
package main


deny[msg] {
  setuid_permissions_used
  msg = sprintf(
      "4.8 Level 1 benchmark - setuid file permission validation: (%s, %d)",
      [ input.location.path, input.metadata.mode ] )
}

setuid_permissions_used {

  input.metadata.type == "RegularFile"
  permission := format_int(input.metadata.mode, 10)
  regex.match("20......", permission)

}

deny[msg] {

  setgid_permissions_used
  msg = sprintf(
      "4.8 Level 1 benchmark - setgid file permission validation: (%s, %d)",
      [ input.location.path, input.metadata.mode ] )

}

setgid_permissions_used {

  input.metadata.type == "RegularFile"
  permission := format_int(input.metadata.mode, 10)
  regex.match("40......", permission)

}
