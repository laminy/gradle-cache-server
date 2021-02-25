# gradle-cache-server
Golang simple gradle cache server. Used to increase gradle build speed
## Configure server
Create sample directory structure
```text
. /opt
├── data
│   ├── cache
│   ├── gradle-cache-server
│   └── config.json
```
Where

| object | type | comment |
|---|---|---|
| cache | directory | directory for cache |
| gradle-cache-server | file | server executable |
| config.json | file | configuration file |

### Server settings
Server settings stored in file `config.json`
```json
{
  "path": "/opt/data/cache",
  "port": 32555,
  "scan": "1h",
  "alive": "7d"
}
```
Where

| property | type | comment |
|---|---|---|
| path | string | cache directory |
| port | int | server port |
| scan | string | cache clear: scan interval |
| alive | string | cache clear: max time old file alive before delete |

#### Scan and alive formats
Format is `{number}{literal}`, where

| literal | comment |
|---|---|
| s | second |
| m | minute |
| h | hour |
| d | day |
| w | week |

## Start server
```shell script
cd /opt/data
./gradle-cache-server
```
Note: server reads file `config.json` from current workdir

## Configure gradle project
### Configure gradle settings file
In this sample gradle cache server started at the same computer, where gradle builds. But it can be started on a remote server.
Add next lines to `settings.gradle`
```groovy
buildCache {
    remote(HttpBuildCache) {
        url = 'http://localhost:32555/'
        push = true
    }
}
```
For full list of configuration properties see [documentation](https://docs.gradle.org/current/userguide/build_cache.html#sec:build_cache_configure).

Note: server supports a cache subpath, like http://localhost:32555/team1/
### Start build with cache
From command line
```shell script
gradle build --build-cache
```
Or from `gradle.properties` file
```properties
org.gradle.caching=true
```

## Monitoring
Prometheus metrics available at `/metrics` endpoint on the same port.
Metrics

| name | type | comment |
|---|---|---|
| cache_get | counter | Total number of requested items |
| cache_put | counter | Total number of saved items |
| cache_del | counter | Total number of deleted items |

#### Metrics attributes (labels)

| name | comment |
|---|---|
| code | http status code (200 or 404) |
| path0 | part#0 of url  |
| path1 | part#1 of url |

#### Part examples

| Request URL | path0 | path1 |
|---|---|---|
| /9af937214740e6369e0d9c47dbea521d| empty string | empty string |
| /service/9af937214740e6369e0d9c47dbea521d| service | empty string |
| /team/service/9af937214740e6369e0d9c47dbea521d| team | service |
| /space/team/service/9af937214740e6369e0d9c47dbea521d| space | team |
