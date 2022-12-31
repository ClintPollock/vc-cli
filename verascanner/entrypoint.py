#!/usr/bin/python
import sys
import subprocess
import io
import os
import shutil
import readline
import json
import argparse
import logging
import zipfile
import tarfile
import tempfile
import hashlib


knownRunCommands= [
                    "bash", "echo",
                    "trivy", "grype", "syft",
                    "dockle", "docker", "skopeo",
                    "opa", "conftest", "yor","checkov"
                  ]

#
# Just exec out to a shell process for a limited set of commands
#
def handleRun(args):

    command = args[0]
    cmdArgs = args[1:]

    if not command in knownRunCommands:
        raise Exception("unknown run command: \"%s\"" % ( command) )

    shell_results = None
    if command == "bash":
        shell_results = subprocess.run(cmdArgs)
    else:
        shell_results = subprocess.run(args)

    sys.exit(shell_results.returncode)


def setupArgParser():

    parser = argparse.ArgumentParser(description="Container entrypoint driver", exit_on_error=False)
    parser.add_argument("command", choices=['inspect','sbom', 'scan', 'research', 'kitchen-sink'], help="run a specific set of local commands directly")
    parser.add_argument("type", help="run a specific set of local commands directly", choices=['image', 'repo', 'archive', 'directory'], default="image")
    parser.add_argument("target", help="target object to scan")

    parser.add_argument("--prettyprint", action=argparse.BooleanOptionalAction, default=False, required=False)
    parser.set_defaults(prettyprint=False)

    parser.add_argument("--verbose", action=argparse.BooleanOptionalAction, type=bool, help="verbosity of scan process", default=False, required=False)
    parser.set_defaults(verbose=False)

    parser.add_argument("--everything", help="use all the tools available when scanning or analyzing", default=False, required=False)
    parser.add_argument(
       "--format",
       default="json",
       choices=[
                 'json',
                 'spdx', 'spdx-tag-value', 'spdx-json',
                 'cyclonedx', 'cyclonedx-xml', 'cyclonedx-json',
                 'github', 'github-json',
                 'directory', 'table', 'text', 'sarif'
               ],
       required=False,
       help="format to use when printing out SBOM or results",
       )
    parser.set_defaults(format="json")
    parser.add_argument("--out", help="output file", default="", required=False)

    return parser

def printMe(d, pretty):
    if pretty:
        print( json.dumps(d, indent=2) )
    else:
        print( json.dumps(d))

ORIG_WORKING_DIR = os.getcwd()
def setWorkingDir(target):

    if target["type"] == "directory":
        os.chdir("/local-context")

def resetWorkingDir():
    os.chdir(ORIG_WORKING_DIR)

#------------------------
# Target handlers
#------------------------

# {
#     "type" : "type of target (image, archive, repo, directory"
#     "raw_name" : "literal string of target provided to command"
#     "id" : "unique identifier generated from the object in question"
#     "hash" : "sha256 hash of target object, or git hash"
#     "name" : "local name of target object"
# }

def normalizeName(string):
    return string.replace( "://", "_" ).replace( "/", "_").replace(":", "_")

def targetCacheExists(target):
    return os.path.exists(target["id"] + ".target.json")

def saveTargetToCache(target):

    # serialize target as JSON
    fn = target["id"] + ".target.json"
    with open(fn, 'w') as outfile:
        json.dump(target, outfile)

    # Capture ID into text file
    fn = target["name"] + ".txt"
    if fn == "..txt": fn = "dot.txt"
    with open(fn, 'w') as outfile:
        outfile.write(target["id"])

def loadTargetFromCache(target):

    target = None
    with open(target["id"] + ".target.json", 'r') as infile:
        target = json.load(infile)

    return target

def dockerPull(target, args):

    # make sure the image is local
    cmd = ["docker", "pull", target["raw_name"] ]
    proc = subprocess.Popen(cmd,stdout=subprocess.PIPE)
    results = ""
    for line in io.TextIOWrapper(proc.stdout, encoding="utf-8"):
        results += line
    if args.verbose:
        print( results )

def inspectImageWithDocker(target):

    # read in info about the image from docker
    cmd = ["docker", "inspect", target["raw_name"]]
    proc = subprocess.Popen(cmd,stdout=subprocess.PIPE)
    results = ""
    for line in io.TextIOWrapper(proc.stdout, encoding="utf-8"):
        results += line

    data = {}
    try:
        data = json.loads(results)[0]
    except:
        print("ERROR: could not retrieve data on container image %s" % target["raw_name"])
        sys.exit(1)

    return data

def inspectImageWithSkopeo(target):

    cmd = ["skopeo", "inspect", "docker-daemon:%s" % ( target["raw_name"]) ]
    proc = subprocess.Popen(cmd,stdout=subprocess.PIPE)
    results = ""
    for line in io.TextIOWrapper(proc.stdout, encoding="utf-8"):
        results += line

    return json.loads(results)

def imageHandler(args):
    target = { "type" : "image" }
    target["raw_name"] = args.target
    if not ":" in target["raw_name"]:
        target["raw_name"] += ":latest"

    target["name"] = normalizeName( args.target )

    # make sure the image is local
    dockerPull(target, args)

    # use docker inspect  and use ID (SHA256) as the overall ID

    target["docker"] = inspectImageWithDocker(target)
    target["id"] = target["docker"]["Id"]

    # read in data about it from skopeo too

    #
    # Run other tools as well
    #
    if args.everything:
        try:
            target["skopeo"] = inspectImageWithSkopeo(target)
        except:
            pass

    saveTargetToCache(target)
    return target

def repoHandler(args):
    target = { "type" : "repo" }
    target["raw_name"] = args.target
    target["local_working_copy"] = os.path.normpath( "git+" + normalizeName( args.target ).replace( ".git", "") )
    target["local_git_dir"] = os.path.normpath( "%s.git" % ( target["local_working_copy"] ) )

    # This needs to be pythonified ... hacky shell callouts
    subprocess.run([ "rm", "-rf", target["local_working_copy"], target["local_git_dir"] ] )

    subprocess.run([
                    "git",
                    "clone",
                    "-q",
                    "--separate-git-dir=%s"% target["local_git_dir"],
                    target["raw_name"],
                    target["local_working_copy"]] )

    wd =  os.getcwd()
    os.chdir(target["local_git_dir"])
    cmd = [ "git",  "rev-parse", "HEAD"]
    proc = subprocess.Popen(cmd,stdout=subprocess.PIPE)
    results = ""
    for line in io.TextIOWrapper(proc.stdout, encoding="utf-8"):
        results += line
    os.chdir(wd)

    target["githash"] = results.rstrip()
    target["id"] = target["githash"]

    subprocess.run([ "rm", "-rf", target["local_git_dir"] ] )
    subprocess.run([ "rm", "-f", "./%s/.git" % ( target["local_working_copy"] ) ] )

    target["name"] = target["local_working_copy"]+"_"+target["githash"]
    subprocess.run([ "rm", "-rf", target["name"] ] )
    os.rename(target["local_working_copy"], target["name"])
    target["local_working_copy"] = target["name"]

    saveTargetToCache(target)
    return target

def archiveHandler(args):

    target = { "type" : "archive" }
    target["raw_name"] = os.path.normpath( args.target )
    target["localpath"] = os.path.normpath( "/local-context/" + args.target )

    # calculate SHA256 hash
    sha256_hash = hashlib.sha256()
    with open(target["localpath"],"rb") as f:
        # Read and update hash string value in blocks of 4K
        for byte_block in iter(lambda: f.read(4096),b""):
            sha256_hash.update(byte_block)
    target["hash"] = sha256_hash.hexdigest()
    target["id"] = target["hash"]

    target["name"] = normalizeName( target["raw_name"] )
    target["name"] = target["name"].replace(".zip", "_zip").replace(".tar.gz", "_tar_gz").replace(".tar", "_tar")
    target["name"] = target["name"] + "_" + target["hash"]

    try:
        os.mkdir( target["name"] )
    except:
        pass

    try:
        with zipfile.ZipFile(target["localpath"], 'r') as zip_ref:
            zip_ref.extractall(target["name"])
    except:
        pass

    try:
        tar_file = tarfile.open(target["localpath"])
        tar_file.extractall(target["name"])
        tar_file.close()
    except:
        pass

    target["local_working_copy"] = target["name"]

    saveTargetToCache(target)
    return target

def directoryHandler(args):

    target_path = os.path.normpath( args.target )
    if target_path.startswith("/") or target_path.startswith("../"):
        print("ERROR: when scanning a directory, the path must be relative to your current working directory.")
        sys.exit(1)

    target = { "type" : "directory" }
    target["raw_name"] = target_path
    target["name"]     = normalizeName( target["raw_name"] )

    statString = target["raw_name"] + ":" + str(os.stat( os.path.join("/local-context", target["raw_name"]) ))
    sha256_hash = hashlib.sha256()
    sha256_hash.update(statString.encode())
    target["id"] = target["hash"] = sha256_hash.hexdigest()

    saveTargetToCache(target)
    return target


targetHandlers = { "image":imageHandler, "repo":repoHandler, "archive":archiveHandler, "directory":directoryHandler }

# {
#     target_id : "sha26:...",
#     target_type : .. "image" or "repo" or "archive" or "directory"
#     files: [],
#     sbom: { "syft schema SBOM ..." }
# }

def generateInventory(target):

    inventory = {}
    inventoryCacheFile = target["id"]+".inventory.json"
    sbomCacheFile = target["id"]+".sbom.json"

    if os.path.exists(inventoryCacheFile):
        with open(inventoryCacheFile, 'r') as infile:
            inventory = json.load(infile)

    else:

        setWorkingDir(target)

        exclude_options = ""
        tgt = "dir:" + target["name"]
        if target["type"] == "image":
            tgt = "docker:" + target["raw_name"]
        if target["type"] == "directory":
            tgt = "dir:" + target["raw_name"]

        cmd = ["syft", "packages", "-o", "json", tgt]

        if target["type"] == "directory":
            cmd += [
                     "--exclude", "./.DS_Store",
                     "--exclude", "*/.DS_Store",
                     "--exclude", "./.git/**",
                     "--exclude", "./.git/*",
                     "--exclude", "**/.git/**",
                     "--exclude", "**/.git/*"
                   ]

        proc = subprocess.Popen(cmd,stdout=subprocess.PIPE)
        results = ""
        for line in io.TextIOWrapper(proc.stdout, encoding="utf-8"):
            results += line

        resetWorkingDir()

        inventory = { "target_id" : target["id"], "target_type" : target["type"] }
        inventory [ "cache_file" ] = inventoryCacheFile
        inventory [ "sbom_file" ] = sbomCacheFile
        try:
            inventory ["sbom"] = json.loads(results)
        except:
            print("ERROR: could load generate SBOM against " + tgt )
            print(results)

        inventory ["files"] = inventory ["sbom"]["files"]
        #inventory ["sbom"]["files"] = []

        # serialize inventory as JSON
        with open(inventoryCacheFile, 'w') as outfile:
            json.dump(inventory, outfile)

        with open(sbomCacheFile, 'w') as outfile:
            json.dump(inventory["sbom"], outfile)

    return inventory

# Maps our format names to that in syft
formatMap = {
   'json':'syft-json',
   'spdx':'spdx-tag-value',
   'spdx-json':'spdx-json',
   'spdx-tag-value':'spdx-tag-value',
   'cyclonedx-xml':'cyclonedx-xml',
   'cyclonedx':'cyclonedx-xml',
   'cyclonedx-json':'cyclonedx-json',
   'github':'github',
   'github-json':'github-json',
   'table':'table',
   'text':'text',
   'sarif':'syft-json'
   }

def printToSBOMFormat(inventory, format, pp=None):

    cmd = ["syft", "convert", inventory[ "sbom_file" ], "-o", formatMap[format] ]
    proc = subprocess.Popen(cmd,stdout=subprocess.PIPE)
    results = ""
    for line in io.TextIOWrapper(proc.stdout, encoding="utf-8"):
        results += line

    if "json" in format:
        printMe( json.loads(results), pp )
    else:
        print(results)


def doGrypeVulnerabilityScan(inventory):

    # read in sbom and return vuln report
    cmd = ["grype",  inventory [ "sbom_file" ], "-o", "json"  ]
    proc = subprocess.Popen(cmd,stdout=subprocess.PIPE)
    results = ""
    for line in io.TextIOWrapper(proc.stdout, encoding="utf-8"):
        results += line

    scanResult = json.loads(results)

    return scanResult


target2trivyTypeMap = { "image" : "image", "directory" : "fs", "repo" : "fs", "archive" : "fs" }

def doTrivyScan( target, mode ):

    setWorkingDir(target)

    type = target2trivyTypeMap[target["type"]]
    tgt = target["name"]
    if target["type"] in ["image", "directory"]: tgt = target["raw_name"]

    cmd = ["trivy", "-q", type, "--security-checks", mode, "-f", "json", tgt ]
    proc = subprocess.Popen(cmd,stdout=subprocess.PIPE)
    results = ""
    for line in io.TextIOWrapper(proc.stdout, encoding="utf-8"):
        results += line

    resetWorkingDir()

    return json.loads(results)

def doTrivyVulnScan( target ): return doTrivyScan( target, "vuln" )
def doTrivySecretScan( target ): return doTrivyScan( target, "secret" )
def doTrivyIaCScan( target ): return doTrivyScan( target, "config" )

def doDockerLint( target ):

    # dockle looks like intends to support a convert function but it doesnt work
    # so rerun for knowledge
    cmd = ["dockle", "-q", "-f", "json", target["raw_name"]  ]

    proc = subprocess.Popen(cmd,stdout=subprocess.PIPE)
    results = ""
    for line in io.TextIOWrapper(proc.stdout, encoding="utf-8"):
        results += line

    return json.loads(results)


def doSecretScan( target ):
    return doTrivySecretScan( target )

def doCheckovSecretScan( target ):

    setWorkingDir(target)

    # dockle looks like intends to support a convert function but it doesnt work
    # so rerun for knowledge
    #checkov --directory /data -o json
    if target["type"] == "image": return None

    tgt = target["name"]

    cmd = ["checkov", "-o", "json", "-d",  tgt ]
    proc = subprocess.Popen(cmd,stdout=subprocess.PIPE)
    results = ""
    for line in io.TextIOWrapper(proc.stdout, encoding="utf-8"):
        results += line

    resetWorkingDir()

    return json.loads(results)

def doIaCScan(target):
    return doTrivyIaCScan(target)

opaPermissionsPolicyPath = "/opa/policy/permissions"

def doFilePermissionsScan(inventory):

    # conftest test --policy /opa/policy/permissions -o json ${source_id}.manifest.json
    # serialize target as JSON
    fn = inventory["target_id"] + ".files.json"
    with open(fn, 'w') as outfile:
        json.dump(inventory["files"] , outfile)

    cmd = ["conftest", "test",  fn , "--policy", opaPermissionsPolicyPath, "-o", "json" ]
    proc = subprocess.Popen(cmd,stdout=subprocess.PIPE)
    results = ""
    for line in io.TextIOWrapper(proc.stdout, encoding="utf-8"):
        results += line

    return json.loads(results)

def doScan(target, inventory):

    data = { "target_id" : inventory["target_id"], "findings" : {} }

    data["findings"]["vulnerabilities"] = doGrypeVulnerabilityScan(inventory)

    if inventory["target_type"] == "image":
        data["findings"]["docker"] = doDockerLint(target)

    #scanResult["trivyVulnerabilities"] = doTrivyVulnScan(target)
    data["findings"]["secrets"] = doSecretScan(target)
    data["findings"]["iac"] = doIaCScan(target)

    #scanResult["checkov"] = doCheckovSecretScan(target)

    data["findings"]["permissions"] = doFilePermissionsScan(inventory)

    return data

opaPolicyPath = "/opa/policy/results"

def evaluatePolicy(inventory, results):

    fn = inventory["target_id"] + ".results.json"
    with open(fn, 'w') as outfile:
        json.dump(results , outfile)

    cmd = ["conftest", "test",  fn , "--policy", opaPolicyPath, "-o", "json" ]
    proc = subprocess.Popen(cmd,stdout=subprocess.PIPE)
    results = ""
    for line in io.TextIOWrapper(proc.stdout, encoding="utf-8"):
        results += line

    retval = proc.wait()
    return ( json.loads(results), retval ) # json.loads(results)


def printScanToFormat( inventory, format, pp=None ):

    # read in sbom and return vuln report
    cmd = ["grype",  inventory [ "sbom_file" ], "-o", format ]
    proc = subprocess.Popen(cmd,stdout=subprocess.PIPE)
    results = ""
    for line in io.TextIOWrapper(proc.stdout, encoding="utf-8"):
        results += line

    print ( results )

def cleanup( target ):

    if "local_working_copy" in target.keys() :
        shutil.rmtree(target["local_working_copy"])


#------------------------
# Command Handlers
#------------------------

def inspectHandler(args):

    target = targetHandlers[args.type](args)

    printMe(target, args.prettyprint)

    cleanup(target)


def sbomHandler(args):

    target = targetHandlers[args.type](args)
    inventory = generateInventory( target )

    if args.format == "json" or args.format == "json":
        printMe(inventory["sbom"], args.prettyprint)
    else:
        printToSBOMFormat(inventory, args.format, args.prettyprint )

    cleanup(target)

def scanHandler(args):

    target = targetHandlers[args.type](args)
    inventory = generateInventory( target )
    data = doScan(target, inventory)

    (policyResults, retval) = evaluatePolicy(inventory, data["findings"])
    data["policy-results"] = policyResults
    data["policy-passed"]  = True if retval == 0 else False

    if args.format in ["json", "spdx" ]:
        printMe(data, args.prettyprint)
    else:
        printScanToFormat(inventory, args.format )

    cleanup(target)
    sys.exit( retval )


#
# Assemble all the data we can find into one uber report
#
def kitchenSinkDataGenerator(args):

    data = {}
    target = targetHandlers[args.type](args)
    inventory = generateInventory( target )
    scan = doScan(target, inventory)

    data["target"] = target
    data["inventory"] = inventory
    data["findings"] = scan["findings"]
    data["findings-by-tool"] = {}

    if data["inventory"]["target_type"] == "image":
        data["findings-by-tool"]["dockle"] = {}
        data["findings-by-tool"]["dockle"]["docker"] = doDockerLint(target)
        data["findings-by-tool"]["skopeo"] = {}
        data["findings-by-tool"]["skopeo"]["inventory"]        = inspectImageWithSkopeo(target)

    # Grype
    data["findings-by-tool"]["grype"] = {}
    data["findings-by-tool"]["grype"]["vulnerabilities"] = data["findings"]["vulnerabilities"]

    # Trivy
    data["findings-by-tool"]["trivy"] = {}
    data["findings-by-tool"]["trivy"]["vulnerabilities"] = doTrivyVulnScan(target)
    data["findings-by-tool"]["trivy"]["secrets"]   = data["findings"]["secrets"]
    data["findings-by-tool"]["trivy"]["iac"] = data["findings"]["iac"]

    data["findings-by-tool"]["checkov"] = {}
    data["findings-by-tool"]["checkov"]["secrets"] = doCheckovSecretScan(target)


    (policyResults, retval) = evaluatePolicy(inventory, scan)
    data["policy-results"] = policyResults
    data["policy-passed"] = True if retval == 0 else False

    return data

def kitchenSinkHandler(args):

    printMe(kitchenSinkDataGenerator(args), args.prettyprint)


def researchHandler(args):

    data = kitchenSinkDataGenerator(args)

    vulnerabilities = { "grype": {}, "trivy" : {} }

    grype_vulns = data["findings-by-tool"]["grype"]["vulnerabilities"]
    vulnerabilities["grype"]["distro"] = grype_vulns["distro"]["name"] + " " + grype_vulns["distro"]["version"]
    vulns = []
    sev_vulns = []
    cvss_vulns = []
    for m in grype_vulns["matches"]:
        id = m["vulnerability"]["id"]
        vulns.append( id )
        sev_vulns.append( m["vulnerability"]["severity"].upper() + " - " +  id )
        cvss = "unknown"
        if "cvss" in m.keys():
            for cvss_instance in m["cvss"]:
                if cvss_instance["version"].startswith("3"):
                    cvss = cvss_instance["metrics"]["baseScore"]

        else:
            for rv in m["relatedVulnerabilities"]:
                if "cvss" in rv.keys():
                    for cvss_j in rv["cvss"]:
                        if cvss_j["version"].startswith("3"):
                            cvss = cvss_j["metrics"]["baseScore"]

        cvss_vulns.append( m["vulnerability"]["id"] + " - " + cvss )

    # convert to set to get a unique list, then sort
    vulns = list(set(vulns))
    vulns.sort()
    sev_vulns = list(set(sev_vulns))
    sev_vulns.sort()
    cvss_vulns = list(set(cvss_vulns))
    cvss_vulns.sort()

    vulnerabilities["grype"]["vulnerabilities"] = vulns
    vulnerabilities["grype"]["severity_vulnerabilities"] = sev_vulns
    vulnerabilities["grype"]["cvss_vulnerabilities"] = cvss_vulns
    vulnerabilities["grype"]["N"] = len(vulns)

    trivy_vulns = data["findings-by-tool"]["trivy"]["vulnerabilities"]
    #CVEs["trivy"]["distro"] = trivy_vulns["distro"]["name"] + " " + grype_vulns["distro"]["version"]
    vulns = []
    sev_vulns = []
    cvss_vulns = []
    for r in trivy_vulns["Results"]:
        if r != None and "Vulnerabilities" in r.keys():
            for v in r["Vulnerabilities"]:
                vulns.append( v["VulnerabilityID"] )
                sev_vulns.append( v["Severity"].upper() + " - " + v["VulnerabilityID"]  )
                for vendor in v["CVSS"].keys():
                    cvss_vulns.append( v["VulnerabilityID"] + " - " + v["CVSS"][vendor]["V3Score"] )


    # convert to set to get a unique list, then sort
    vulns = list(set(vulns))
    vulns.sort()
    sev_vulns = list(set(sev_vulns))
    sev_vulns.sort()
    cvss_vulns = list(set(cvss_vulns))
    cvss_vulns.sort()

    vulnerabilities["trivy"]["vulnerabilities"] = vulns
    vulnerabilities["trivy"]["severity_vulnerabilities"] = sev_vulns
    vulnerabilities["trivy"]["cvss_vulnerabilities"] = cvss_vulns
    vulnerabilities["trivy"]["N"] = len(vulns)

    total = set(
                 vulnerabilities["grype"]["vulnerabilities"]
                ).union (
                  set(vulnerabilities["trivy"]["vulnerabilities"])
                )

    agreed = set(
                  vulnerabilities["grype"]["vulnerabilities"]
                ).intersection(
                  set(vulnerabilities["trivy"]["vulnerabilities"]))

    diffs = total - agreed

    total = list(total)
    total.sort()

    agreed = list(agreed)
    agreed.sort()

    diffs = list(diffs)
    diffs.sort()

    sv_total = set(
                 vulnerabilities["grype"]["severity_vulnerabilities"]
                ).union (
                  set(vulnerabilities["trivy"]["severity_vulnerabilities"])
                )

    sv_agreed = set(
                  vulnerabilities["grype"]["severity_vulnerabilities"]
                ).intersection(
                  set(vulnerabilities["trivy"]["severity_vulnerabilities"]))

    sv_diffs = sv_total - sv_agreed

    sv_total = list(sv_total)
    sv_total.sort()

    sv_agreed = list(sv_agreed)
    sv_agreed.sort()

    sv_diffs = list(sv_diffs)
    sv_diffs.sort()

    cvss_total = set(
                 vulnerabilities["grype"]["cvss_vulnerabilities"]
                ).union (
                  set(vulnerabilities["trivy"]["cvss_vulnerabilities"])
                )

    cvss_agreed = set(
                  vulnerabilities["grype"]["cvss_vulnerabilities"]
                ).intersection(
                  set(vulnerabilities["trivy"]["cvss_vulnerabilities"]))

    cvss_diffs = cvss_total - cvss_agreed

    cvss_total = list(cvss_total)
    cvss_total.sort()

    cvss_agreed = list(cvss_agreed)
    cvss_agreed.sort()

    cvss_diffs = list(cvss_diffs)
    cvss_diffs.sort()

    vulnerabilities["all"] = total
    vulnerabilities["N_all"] = len(total)
    vulnerabilities["shared"] = agreed
    vulnerabilities["N_shared"] = len(agreed)
    vulnerabilities["disputed"] = diffs
    vulnerabilities["N_disputed"] = len(diffs)

    vulnerabilities["sv_all"] = sv_total
    vulnerabilities["sv_N_all"] = len(sv_total)
    vulnerabilities["sv_shared"] = sv_agreed
    vulnerabilities["sv_N_shared"] = len(sv_agreed)
    vulnerabilities["sv_disputed"] = sv_diffs
    vulnerabilities["sv_N_disputed"] = len(sv_diffs)

    vulnerabilities["cvss_all"] = cvss_total
    vulnerabilities["cvss_N_all"] = len(cvss_total)
    vulnerabilities["cvss_shared"] = cvss_agreed
    vulnerabilities["cvss_N_shared"] = len(cvss_agreed)
    vulnerabilities["cvss_disputed"] = cvss_diffs
    vulnerabilities["cvss_N_disputed"] = len(cvss_diffs)

    data["vulnerability_ids"] = vulnerabilities
    printMe(data, args.prettyprint)


commandHandlers = { "inspect":inspectHandler, "scan":scanHandler, "sbom":sbomHandler, "research":researchHandler, "kitchen-sink":kitchenSinkHandler }



def main():

    #
    # Bypass argparse if the command is "run"
    #
    if "run" == sys.argv[1]:
        handleRun( sys.argv[2:] )

    else:
        parser = setupArgParser()
        args = parser.parse_args()

        commandHandlers[args.command](args)


if __name__ == "__main__":
    try:
        main()
    # except ValidationException as e:
    #     logging.error(e)
    except Exception as e:
        logging.error(e)
        logging.exception(e)
