# Data Processing

This folder contains everything to do with gathering and importing of previously collected data.

DO NOT UPLOAD TO GIT

## IndexedDB Collection

As all previous data in Antonium was stored per domain inside of indexedDB, a
way oif extracting this information was needed.
[This](https://dfahlander.medium.com/export-indexeddb-from-a-web-app-using-devtools-62c55a8996a1)
information was found to guide me through this process. The data is now located in data.json for later use.

Note: This data is what is stored in Antonium. Any filters that Antonium does to
this data will be inherited, e.g. single words only.

## Processing of Collected Data

This collected data now needs to be processed into a Redis consumable format.
