# Data Processing

This folder contains everything to do with gathering and importing of previously collected data.

> [!CAUTION]
> DATA SHOULD NOT BE UPLOADED TO GIT

## IndexedDB Collection

As all previous data in Antonium was stored per web domain inside of indexedDB, a
way of extracting this information was needed.
[This](https://dfahlander.medium.com/export-indexeddb-from-a-web-app-using-devtools-62c55a8996a1)
information was found to guide me through this process. The data is now located in data.json for later use.

> [!NOTE]
> This data is what is stored in Antonium. Any filters that Antonium does to
> this data will be inherited, e.g. single words only.

## Processing of Collected Data

This collected data now needs to be processed into a SQLite consumable format.
This is handled by the bulk import endpoint. This takes in the raw json that
process.json produces.
