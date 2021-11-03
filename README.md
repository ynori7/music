# Music
This repository scrapes data from sites like Allmusic to give tailored information 
such as notifactions about interesting new releases.

## Features

### New Releases
The new releases command gathers information about new releases from Allmusic based on
the configured genres, and for each one checks if it's interesting based on sub-genres
and ratings. It can then generate an HTML report which is sent by email.

**Additional Details:**
- Artists whose best rating isn't at least 4 stars will be filtered out
- Emails are sent using Mailjet

**Usage:**

```
go run cmd/newreleases/main.go --config config.yaml \
    --new-release-week 20200327 --output out
```

- `--config` This flag is required and is the path to the configuration YAML.
- `--new-release-week` This flag is optional and indicates which specific week should be fetched 
(it's always a Thursday in the format yyyyMMdd). By default it uses the current week.
- `--output` This is an optional flag to indicate where html files should be saved. By default it's `./out`

Be sure to first copy `config.yaml.dist` to `config.yaml` and fill in the missing blanks

**Set up cronjob:**

First, build the binary:
```
go build -o newreleasesmailer cmd/newreleases/main.go
```

Then set up the cronjob:
```
0 14 * * 5 /path/to/goprojects/src/github.com/ynori7/music/newreleasesmailer --config /path/to/goprojects/src/github.com/ynori7/music/config.yaml --output /path/to/goprojects/src/github.com/ynori7/music/out
```

Note that the new releases page is available on Allmusic on Friday afternoons, so that's when we schedule the job to run.

## Project Structure

Commands are located in `cmd` and are the main entry points.

The `newreleases` command gathers configuration from the `config` package, then sets up
a `newreleases/handler` which orchestrates fetching data from `allmusic` and then filtering 
using the `filter` worker pool.