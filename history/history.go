package history

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

	out, err := exec.Command("git", "clone", repo, target).Output()

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

// History is the data stored within a git log output.
type History struct {
	Commit  string `json:"commit"`
	Subject string `json:"subject"`
	// Name is the author name.
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Date  time.Time `json:"date"`
}

// CleanUp cleans up the temporary directory where the git repo has been stored.
func (g Git) CleanUp() {
	// os.RemoveAll(g.Location)
}

// Log gets the git log of the repository.
func (g Git) Log() ([]History, error) {
	history := []History{}

	logfmt := `--pretty=format:{ "commit": "%H", "subject": "%f", "name": "%aN", "email": "%aE", "date": "%aI"}`
	out, err := exec.Command("git", "log", "--reverse", logfmt).Output()

	if err != nil {
		return history, fmt.Errorf("failed to get the log of %s with err '%v' and message '%s'", g.Location, err, string(out))
	}

	for _, line := range strings.Split(string(out), "\n") {
		h := &History{}
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

	cmd := exec.Command("git", "reset", "--hard", hash)
	cmd.Dir = g.Location
	out, err := cmd.Output()

	if err != nil {
		return fmt.Errorf("failed to checkout a specific hash of %s with err '%v' and message '%s'", g.Location, err, string(out))
	}

	return nil
}

func (g Git) Revert() error {
	os.Chdir(g.Location)
}
