Create directory:
https://www.jfrog.com/confluence/display/JFROG/Artifactory+REST+API#ArtifactoryRESTAPI-CreateDirectory

Upload? Usually the maven-deploy-plugin is used as a client, so let's have a look at how this client works:
https://maven.apache.org/plugins/maven-deploy-plugin/
https://github.com/apache/maven-deploy-plugin/

Class ArtifactRepository is used for configuration. See class MavenProject (org.apache.maven:maven-core:3.0) Line 1845ff for authentication configuration

org.eclipse.aether.repository.RemoteRepository:
- String url
- Authentication authentication

path = "${groupId.replace('.', '/')}/$artifactId/$version/$artifactId-$version"

Looks like an upload is quite simple:
PUT $uri

See https://www.jfrog.com/confluence/display/JFROG/Artifactory+REST+API#ArtifactoryRESTAPI-DeployArtifact

```bash
curl -X PUT "http://localhost:8081/artifactory/my-repository/my/new/artifact/directory/file.txt" \
     -T Desktop/myNewFile.txt \
     -u myUser:myP455w0rd!
```
HINT: The same API applies to Nexus 3 according to a [forum post](https://support.sonatype.com/hc/en-us/articles/115006744008-How-can-I-programmatically-upload-files-into-Nexus-3-)

CONFIRMED BY TESTING! The API is actually a simple PUT on a context path beginning with the repository name.
Prefix is by-default "/artifactory/", schema, host and port for a local installation are http://localhost:8081 by default.

Example of an actual event message:

```http
POST /webhook HTTP/1.1
User-Agent: JFrog Event/7.15.3 (revision: 3ac54edd87, build date: 2021-03-02T13:15:20Z)
Accept-Encoding: gzip
Content-Type: application/json
Content-Length: 237
Uber-Trace-Id: 034adca87246b45d:0ef6757aefed8e21:77537024be3d4d40:0
X-Jfrog-Node-Id: Wicked
X-Jfrog-Service-Id: jfevt@01f0662jmb5wf21q9cmnzw19bt
X-Jfrog-Event-Auth: secret

{
  "domain": "artifact",
  "event_type": "copied",
  "data": {
    "name": "sample.txt",
    "path": "sample_dir/sample.txt",
    "repo_key": "sample_repo",
    "sha256": "sample_checksum",
    "size": 0,
    "source_repo_path": "sample_repo",
    "target_repo_path": "sample_target_repo"
  }
}
```

The deployment returns a json response like this one:

```json
{
  "repo" : "libs-release-local",
  "path" : "/com/example/demo/0.0.1/README-0.0.1.md",
  "created" : "2021-03-09T17:38:20.434+01:00",
  "createdBy" : "andi",
  "downloadUri" : "http://localhost:8081/artifactory/libs-release-local/com/example/demo/0.0.1/README-0.0.1.md",
  "mimeType" : "application/octet-stream",
  "size" : "1923",
  "checksums" : {
    "sha1" : "5c9c3a38bc28ceebd784abc7af9ee7f3ec3838e4",
    "md5" : "9fad8c7c4571238e88cf0868fd0bfe2a",
    "sha256" : "ce31c10b09c556c5dd895e29314de204038c507d116e866d46c91f47be290f29"
  },
  "originalChecksums" : {
    "sha256" : "ce31c10b09c556c5dd895e29314de204038c507d116e866d46c91f47be290f29"
  },
  "uri" : "http://localhost:8081/artifactory/libs-release-local/com/example/demo/0.0.1/README-0.0.1.md"
}
```