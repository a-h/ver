Ver
===
Ver is a tool which analyses Go code stored in Git repositories and determines an 
appropriate version number.

Made in response to the request at https://blog.gopheracademy.com/advent-2016/saga-go-dependency-management/
for:

```
A tool that can statically analyze a project and suggest the next SemVer tag to use 
for the next release
```

Version numbers are generated in the form `Major.Minor.Build` based on a simple algorithm:

## Major
Incremented when binary compatibility is broken (e.g. by removing a function, changing a
function signature, removing a package, or changing the structure of a struct).

## Minor
Incremented when new exported interfaces, functions, constants, structs and their fields 
are added to the package.

## Build
Incremented on any commit, regardless of whether the syntax of the Go can be parsed.

# Usage and output

```
./ver -r https://github.com/a-h/terminator
```

Running this will trigger `ver` to:
 * Clone the repository to a temporary directory.
   * See the `git` package for this. It shells out to the command line as per https://golang.org/src/cmd/go/vcs.go
 * Work through the `git log`, creating a signature of exported items in each commit.
   * (See the `signature` package. It uses the `go/loader` package to parse the code.) 
 * Once a set of signatures are calculated, calculate the version delta 
   according to the algorithm above.
   * See `calculateVersionDelta` and `TestThatVersionDeltasCanBeCalculated`

## Example output

```
Cloned repo https://github.com/a-h/terminator into /var/folders/v0/gv8rbbt9157g5599sh8qljpr0000gn/T/ver_history534925596
Processing git log entry: {7af0cf7ba209d7fd5e5a9ae41122ac746c8d30fc First-commit Adrian Hesketh adrianhesketh@hushmail.com 2016-08-05 17:34:04 +0100 BST}
Processing git log entry: {8e162fcb8302f4c5d216552494b02cb898f1841b Update-README.md Adrian Hesketh a-h@users.noreply.github.com 2016-08-05 17:36:03 +0100 BST}
Processing git log entry: {3158796d558358c20193adae633de31117166d78 Update-launch.sh Adrian Hesketh a-h@users.noreply.github.com 2016-08-05 17:37:00 +0100 BST}
Processing git log entry: {eae6f41ab74f40b34d7a0e095d327c27a90e42b7 Updated-onscreen-messages Adrian Hesketh adrianhesketh@hushmail.com 2016-08-05 18:28:28 +0100 BST}
Processing git log entry: {62b3596f0283bf133e97f7cdf727d460c92c7b32 Ordered-by-version-and-age-for-termination Adrian Hesketh adrianhesketh@hushmail.com 2016-08-06 22:22:30 +0100 BST}
Processing git log entry: {14e8df7c71c08a5d3c4cb15fb5e320d757af4827 Removing-case-sensitivity-on-health-status Adrian Hesketh adrianhesketh@hushmail.com 2016-08-06 22:41:15 +0100 BST}
Processing git log entry: {d8b09b2b12d49d297ebdfa66bb7e6ae531014d2b Added-support-for-help-flags Adrian Hesketh adrianhesketh@hushmail.com 2016-09-01 09:06:39 +0100 BST}
Processing git log entry: {29706f9ed72aa4c3ce51907c5615c8c8e4b00b95 Also-accept-the-v1.0.0-format Adrian Hesketh adrianhesketh@hushmail.com 2016-09-01 14:32:28 +0100 BST}
Processing git log entry: {86fc3f14abb77329a1bfcd8a04e6ed0e8581024b Merge-branch-master-of-https-github.com-a-h-terminator Adrian Hesketh adrianhesketh@hushmail.com 2016-09-01 14:32:35 +0100 BST}
Processing git log entry: {388290610653a31bb0ea186be947870655366f36 Improved-logging Adrian Hesketh adrianhesketh@hushmail.com 2016-09-01 16:43:26 +0100 BST}
Processing git log entry: {e7b613cd0f73e7c6369feee4d26c6dff15257c58 Added-ability-to-filter-by-group-name Adrian Hesketh adrianhesketh@hushmail.com 2016-09-02 16:20:03 +0100 BST}
Processing git log entry: {de81f907e074e8234a6a8530bbdbf3123f48e7e4 Added-name-to-autoscaling-groups-flag Adrian Hesketh adrianhesketh@hushmail.com 2016-09-07 11:24:11 +0100 BST}
Processing git log entry: {9e8b548513f59c0cd10e6472c3f3b66fbe639f01 Added-travis-release Adrian Hesketh adrianhesketh@hushmail.com 2016-09-08 14:07:29 +0100 BST}
Processing git log entry: {5bf462409d5dbac782f2d3cd25f3b0d3bc92fc01 Added-github-deployment Adrian Hesketh adrianhesketh@hushmail.com 2016-09-08 14:18:05 +0100 BST}
Processing git log entry: {50497799c0df6310218dd85e8674256f0264c3a7 Removed-legacy-Go-build Adrian Hesketh adrianhesketh@hushmail.com 2016-09-08 14:24:11 +0100 BST}
Processing git log entry: {851216b760e82fccf172550ffd5dda9764a64418 Merge-branch-master-of-https-github.com-a-h-terminator Adrian Hesketh adrianhesketh@hushmail.com 2016-09-08 14:24:15 +0100 BST}
Processing git log entry: {52a789e656d79087d72b4aa5dcfc5204af8dd631 Skip-cleanup-on-deploy Adrian Hesketh adrianhesketh@hushmail.com 2016-09-08 14:26:23 +0100 BST}
Processing git log entry: {d430e9815b484d052d9c7c9b720bb490eadcdf09 Only-release-on-updated-tags Adrian Hesketh adrianhesketh@hushmail.com 2016-09-08 14:30:09 +0100 BST}
Processing git log entry: {d7eede9a865556de792c0c1cb4adddf1c9083c7d Show-the-result-of-filtration Adrian Hesketh adrianhesketh@hushmail.com 2016-09-14 15:01:32 +0100 BST}
Processing git log entry: {e6d0e02001850e7cf4cddf83b0205331ea85d0f5 Trimmed-quotes-from-version-endpoints-that-produce-JSON Adrian Hesketh adrianhesketh@hushmail.com 2016-09-14 16:30:20 +0100 BST}

Commit 
Version: 0.0.1

Commit 
Version: 0.0.2

Commit 
Version: 0.0.3

Commit 62b3596f0283bf133e97f7cdf727d460c92c7b32
Version: 0.1.4

Commit 14e8df7c71c08a5d3c4cb15fb5e320d757af4827
Version: 0.2.5

Commit d8b09b2b12d49d297ebdfa66bb7e6ae531014d2b
Version: 0.3.6

Commit 29706f9ed72aa4c3ce51907c5615c8c8e4b00b95
Version: 0.4.7

Commit 86fc3f14abb77329a1bfcd8a04e6ed0e8581024b
Version: 0.5.8

Commit 388290610653a31bb0ea186be947870655366f36
Version: 0.6.9

Commit e7b613cd0f73e7c6369feee4d26c6dff15257c58
Version: 0.7.10

Commit de81f907e074e8234a6a8530bbdbf3123f48e7e4
Version: 0.8.11

Commit 9e8b548513f59c0cd10e6472c3f3b66fbe639f01
Version: 0.9.12

Commit 5bf462409d5dbac782f2d3cd25f3b0d3bc92fc01
Version: 0.10.13

Commit 50497799c0df6310218dd85e8674256f0264c3a7
Version: 0.11.14

Commit 851216b760e82fccf172550ffd5dda9764a64418
Version: 0.12.15

Commit 52a789e656d79087d72b4aa5dcfc5204af8dd631
Version: 0.13.16

Commit d430e9815b484d052d9c7c9b720bb490eadcdf09
Version: 0.14.17

Commit d7eede9a865556de792c0c1cb4adddf1c9083c7d
Version: 0.15.18

Commit e6d0e02001850e7cf4cddf83b0205331ea85d0f5
Version: 0.16.19
```

# Possible improvements

 * Avoid having to clone the repo, have the option to take in an existing repo.
 * Reduce the version delta severity of adding new fields to structs.
 * Simplify the output of the tool, especially when a commit can't be parsed 
   due to errors in the source code (e.g. by having verbose / non-verbose mode).
 * Increase processing speed, e.g. by storing the signature of past commits to disk to avoid recalculation.
