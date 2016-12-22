package git

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
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
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Date  time.Time `json:"date"`
}

// CleanUp cleans up the temporary directory where the git repo has been stored.
func (g Git) CleanUp() {
	os.RemoveAll(g.Location)
}

// Log gets the git log of the repository.
func (g Git) Log() ([]Commit, error) {
	history := []Commit{}

	logfmt := `--pretty=format:{ "hash": "%H", "subject": "%f", "name": "%aN", "email": "%aE", "date": "%aI"}`
	out, err := exec.Command("git", "log", "--first-parent", "master", "--reverse", logfmt).CombinedOutput()

	if err != nil {
		return history, fmt.Errorf("failed to get the log of %s with err '%v' and message '%s'", g.Location, err, string(out))
	}

	for _, line := range strings.Split(string(out), "\n") {
		h := &Commit{}
		if err := json.Unmarshal([]byte(line), &h); err != nil {
			return history, fmt.Errorf("failed to parse log line '%s' with err: %v", line, err)
		}
		history = append(history, *h)
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
