package history

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
		repo: target,
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
	repo string
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
	os.RemoveAll(g.repo)
}

// Log gets the git log of the repository.
func (g Git) Log() ([]History, error) {
	history := []History{}

	logfmt := `--pretty=format:{ "commit": "%H", "subject": "%f", "name": "%aN", "email": "%aE", "date": "%aI"}`
	out, err := exec.Command("git", "log", "--reverse", logfmt).Output()

	if err != nil {
		return history, fmt.Errorf("failed to get the log of %s with err '%v' and message '%s'", g.repo, err, string(out))
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
func Get(id string) (directory string, clearup func()) {
	//TODO: Implement.
	return "", func() { log.Print("hello") }
}
