# Bench test for remote dialer

Test case is to compare the preformance between direct access and tunneling.

```golang
//======================
//
//                      (direct)
//         .--------------->----------------.
//        /              tunnel              \
// request--->server:8125------>client---->fileserver:8124
//
//
//======================
```

The test client will request following sizes file from file server in each request.

- 4kb
- 256kb
- 1Mb
- 4Mb
- 16Mb
- 256Mb

My test result is following:

Here is the result of buffer size 4k

| mod    | size  | content transfer | total        |
| ------ | ----- | ---------------- | ------------ |
| direct | 4kb   | 60.58µs          | 668.504µs    |
| direct | 256kb | 550.815µs        | 914.64µs     |
| direct | 1Mb   | 3.224235ms       | 3.616695ms   |
| direct | 4Mb   | 8.22066ms        | 9.15677ms    |
| direct | 16Mb  | 30.850098ms      | 32.328487ms  |
| direct | 256Mb | 477.86015ms      | 482.167926ms |
| proxy  | 4kb   | 40.023µs         | 1.753296ms   |
| proxy  | 256kb | 329.535µs        | 2.156589ms   |
| proxy  | 1Mb   | 2.637765ms       | 3.876888ms   |
| proxy  | 4Mb   | 9.944889ms       | 11.684215ms  |
| proxy  | 16Mb  | 41.398212ms      | 43.768767ms  |
| proxy  | 256Mb | 925.522907ms     | 946.434531ms |

Here is the result of buffer size 64k

| mod    | size  | content transfer | total        |
| ------ | ----- | ---------------- | ------------ |
| direct | 4kb   | 97.498µs         | 866.299µs    |
| direct | 256kb | 624.417µs        | 1.02647ms    |
| direct | 1Mb   | 2.05906ms        | 2.990424ms   |
| direct | 4Mb   | 6.837846ms       | 7.760798ms   |
| direct | 16Mb  | 30.671629ms      | 31.881404ms  |
| direct | 256Mb | 488.276015ms     | 492.736225ms |
| proxy  | 4kb   | 35.241µs         | 1.137234ms   |
| proxy  | 256kb | 256.289µs        | 1.090882ms   |
| proxy  | 1Mb   | 1.331909ms       | 2.269982ms   |
| proxy  | 4Mb   | 5.966625ms       | 7.323062ms   |
| proxy  | 16Mb  | 31.34299ms       | 33.883659ms  |
| proxy  | 256Mb | 604.659363ms     | 623.715397ms |

> PS:
> Here is the log to table regexp `INFO\[[0-9]+\] mod: ([a-z]+), +size: ([0-9]+) +content transfer: ([0-9]+.[0-9]+[µ|m]s) +total: ([0-9]+.[0-9]+[µ|m]s)`
