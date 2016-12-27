package git

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Clone clones the git repository and places it in a temp directory.
func Clone(repo string) (Git, error) {
	target, err := ioutil.TempDir("", "ver_history")

	if err != nil {
		return Git{}, err
	}

	g := Git{
		Location: target,
	}

	out, err := exec.Command("git", "clone", repo, target).CombinedOutput()

	if err != nil {
		return g, fmt.Errorf("failed to clone repo %s to temp directory %s with err '%v' and message %s",
			repo,
			target,
			err,
			string(out))
	}

	return g, nil
}

// Git is a git repository, cloned from the Web.
type Git struct {
	// Location is the location on disk, e.g. /var/tmp/ver_history_12312321/
	Location string
}

// Commit is the data stored within a git log output.
type Commit struct {
	Hash    string `json:"hash"`
	Subject string `json:"subject"`
	// Name is the author name.
	Name      string `json:"name"`
	Email     string `json:"email"`
	Timestamp int64  `json:"timestamp"`
}

// Date Converts the Unix timestamp into a Time.
func (c Commit) Date() time.Time {
	return time.Unix(c.Timestamp, 0)
}

// CleanUp cleans up the temporary directory where the git repo has been stored.
func (g Git) CleanUp() {
	os.RemoveAll(g.Location)
}

// Log gets the git log of the repository.
func (g Git) Log() ([]Commit, error) {
	history := []Commit{}
	separator := ":ec0c7bc17e1ef95b57f47e6ee9f63f54ac187325:"
	logfmt := "--pretty=format:" +
		"%H" + separator + // Hash
		"%f" + separator + // Subject
		"%aN" + separator + // Author Name
		"%aE" + separator + // Author Email
		"%ad" + separator + // Date
		"%at" // Timestamp
	out, err := exec.Command("git", "--no-pager", "log", "--first-parent", "master", "--reverse", logfmt).CombinedOutput()

	if err != nil {
		return history, fmt.Errorf("failed to get the log of %s with err '%v' and message '%s'", g.Location, err, string(out))
	}

	for _, line := range strings.Split(string(out), "\n") {
		lineParts := strings.Split(line, separator)

		if len(lineParts) != 6 {
			return history, fmt.Errorf("failed to parse log line '%s', unexpected number of commit parts found", line)
		}

		ts, err := strconv.ParseInt(lineParts[5], 10, 64)

		if err != nil {
			return history, fmt.Errorf("failed to parse timestamp value of '%s' for line '%s' with err %v", lineParts[5], line, err)
		}

		h := Commit{
			Hash:    lineParts[0],
			Subject: lineParts[1],
			Name:    lineParts[2],
			Email:   lineParts[3],
			// Date: lineParts[4],
			Timestamp: ts,
		}
		history = append(history, h)
	}

	return history, nil
}

// Get extracts all of the files from the given commit into a directory.
func (g Git) Get(hash string) error {
	os.Chdir(g.Location)

	cmd := exec.Command("git", "checkout", hash, "-f")
	cmd.Dir = g.Location
	out, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("failed to checkout hash %s in repo at %s with err '%v' message '%s'", hash, g.Location, err, string(out))
	}

	return nil
}

// Fetch the history from the remote.
func (g Git) Fetch() error {
	os.Chdir(g.Location)

	cmd := exec.Command("git", "fetch", "--all")
	cmd.Dir = g.Location
	out, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("failed to fetch all in repo at %s with err '%v' message '%s'", g.Location, err, string(out))
	}

	return nil
}

// Revert the temporary repository back to HEAD.
func (g Git) Revert() error {
	os.Chdir(g.Location)

	cmd := exec.Command("git", "checkout", "master", "-f")
	cmd.Dir = g.Location
	out, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("failed to revert %s back to head with err '%v' and message '%s'", g.Location, err, string(out))
	}

	return nil
}
