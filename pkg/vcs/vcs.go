// Copyright 2018 syzkaller project authors. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

// Package vcs provides helper functions for working with various repositories (e.g. git).
package vcs

import (
	"bytes"
	"fmt"
	"net/mail"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/syzkaller/dashboard/dashapi"
	"github.com/google/syzkaller/pkg/debugtracer"
	"github.com/google/syzkaller/pkg/osutil"
	"github.com/google/syzkaller/pkg/report/crash"
	"github.com/google/syzkaller/sys/targets"
)

type Repo interface {
	// Poll checkouts the specified repository/branch.
	// This involves fetching/resetting/cloning as necessary to recover from all possible problems.
	// Returns hash of the HEAD commit in the specified branch.
	Poll(repo, branch string) (*Commit, error)

	// CheckoutBranch checkouts the specified repository/branch.
	CheckoutBranch(repo, branch string) (*Commit, error)

	// CheckoutCommit checkouts the specified repository on the specified commit.
	CheckoutCommit(repo, commit string) (*Commit, error)

	// SwitchCommit checkouts the specified commit without fetching.
	SwitchCommit(commit string) (*Commit, error)

	// Commit returns info about the specified commit hash.
	// The commit may be the special value HEAD for the current commit.
	Commit(com string) (*Commit, error)

	// GetCommitByTitle finds commit info by the title. If the commit is not found, nil is returned.
	// Remote is not fetched and only commits reachable from the checked out HEAD are searched
	// (e.g. do CheckoutBranch before).
	GetCommitByTitle(title string) (*Commit, error)

	// GetCommitsByTitles is a batch version of GetCommitByTitle.
	// Returns list of commits and titles of commits that are not found.
	GetCommitsByTitles(titles []string) ([]*Commit, []string, error)

	// ExtractFixTagsFromCommits extracts fixing tags for bugs from git log.
	// Given email = "user@domain.com", it searches for tags of the form "user+tag@domain.com"
	// and returns commits with these tags.
	ExtractFixTagsFromCommits(baseCommit, email string) ([]*Commit, error)

	// ReleaseTag returns the latest release tag that is reachable from the given commit.
	ReleaseTag(commit string) (string, error)

	// Returns true if the current tree contains the specified commit.
	// Remote is not fetched and only commits reachable from the checked out HEAD are searched
	// (e.g. do CheckoutBranch before).
	Contains(commit string) (bool, error)

	// LatestCommits lists all latest commit hashes well as their commit dates.
	// If afterCommit is specified, the output only includes the commits from which afterCommit is reachable.
	// If afterDate is specified, the output only includes the newe commits.
	LatestCommits(afterCommit string, afterDate time.Time) ([]CommitShort, error)

	// Object returns the contents of a git repository object at the particular moment in history.
	Object(name, commit string) ([]byte, error)

	// MergeBases returns good common ancestors of the two commits.
	MergeBases(firstCommit, secondCommit string) ([]*Commit, error)

	// CommitExists check for the commit presence in local checkout.
	CommitExists(commit string) (bool, error)

	// PushCommit is used to store commit in remote repo.
	PushCommit(repo, commit string) error
}

// Bisecter may be optionally implemented by Repo.
type Bisecter interface {
	// Can be used for last minute preparations like pulling release tags into the bisected repo, which
	// is required to determin the compiler version to use on linux. Can be an empty function.
	PrepareBisect() error

	// Bisect bisects good..bad commit range against the provided predicate (wrapper around git bisect).
	// The predicate should return an error only if there is no way to proceed
	// (it will abort the process), if possible it should prefer to return BisectSkip.
	// Progress of the process is streamed to the provided trace.
	// Returns the first commit on which the predicate returns BisectBad,
	// or multiple commits if bisection is inconclusive due to BisectSkip.
	Bisect(bad, good string, dt debugtracer.DebugTracer, pred func() (BisectResult, error)) ([]*Commit, error)

	// PreviousReleaseTags returns list of preceding release tags that are reachable from the given commit.
	// If the commit itself has a release tag, this tag is not included.
	PreviousReleaseTags(commit, compilerType string) ([]string, error)

	IsRelease(commit string) (bool, error)

	EnvForCommit(defaultCompiler, compilerType, binDir, commit string,
		kernelConfig []byte, backports []BackportCommit) (*BisectEnv, error)
}

type ConfigMinimizer interface {
	Minimize(target *targets.Target, original, baseline []byte, types []crash.Type,
		dt debugtracer.DebugTracer, pred func(test []byte) (BisectResult, error)) ([]byte, error)
}

type Commit struct {
	Hash       string
	Title      string
	Author     string
	AuthorName string
	Recipients Recipients
	Tags       []string
	Parents    []string
	Date       time.Time
	CommitDate time.Time
	Patch      []byte
}

type CommitShort struct {
	Hash       string
	CommitDate time.Time
}

type RecipientType int

const (
	To RecipientType = iota
	Cc
)

func (t RecipientType) String() string {
	return [...]string{"To", "Cc"}[t]
}

type RecipientInfo struct {
	Address mail.Address
	Type    RecipientType
}

type Recipients []RecipientInfo

func (r Recipients) GetEmails(filter RecipientType) []string {
	emails := []string{}
	for _, user := range r {
		if user.Type == filter {
			emails = append(emails, user.Address.Address)
		}
	}
	sort.Strings(emails)
	return emails
}

func NewRecipients(emails []string, t RecipientType) Recipients {
	r := Recipients{}
	for _, e := range emails {
		r = append(r, RecipientInfo{mail.Address{Address: e}, t})
	}
	sort.Sort(r)
	return r
}

func (r Recipients) Len() int           { return len(r) }
func (r Recipients) Less(i, j int) bool { return r[i].Address.Address < r[j].Address.Address }
func (r Recipients) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

func (r Recipients) ToDash() dashapi.Recipients {
	d := dashapi.Recipients{}
	for _, user := range r {
		d = append(d, dashapi.RecipientInfo{Address: user.Address, Type: dashapi.RecipientType(user.Type)})
	}
	return d
}

type BisectResult int

const (
	BisectBad BisectResult = iota
	BisectGood
	BisectSkip
)

type BisectEnv struct {
	Compiler     string
	KernelConfig []byte
}

type RepoOpt int

const (
	// RepoPrecious is intended for command-line tools that work with a user-provided repo.
	// Such repo won't be re-created to recover from errors, but rather return errors.
	// If this option is not specified, the repo can be re-created from scratch to recover from any errors.
	OptPrecious RepoOpt = iota
	// Don't use sandboxing suitable for pkg/build.
	OptDontSandbox
)

func NewRepo(os, vmType, dir string, opts ...RepoOpt) (Repo, error) {
	switch os {
	case targets.Linux:
		if vmType == targets.Starnix {
			return newFuchsia(dir, opts), nil
		}
		return newLinux(dir, opts, vmType), nil
	case targets.Fuchsia:
		return newFuchsia(dir, opts), nil
	case targets.OpenBSD:
		return newGitRepo(dir, nil, opts), nil
	case targets.NetBSD:
		return newGitRepo(dir, nil, opts), nil
	case targets.FreeBSD:
		return newGitRepo(dir, nil, opts), nil
	case targets.TestOS:
		return newTestos(dir, opts), nil
	}
	return nil, fmt.Errorf("vcs is unsupported for %v", os)
}

func NewSyzkallerRepo(dir string, opts ...RepoOpt) Repo {
	git := newGitRepo(dir, nil, append(opts, OptDontSandbox))
	return git
}

func NewLKMLRepo(dir string) Repo {
	return newGitRepo(dir, nil, []RepoOpt{OptDontSandbox})
}

func Patch(dir string, patch []byte) error {
	// Do --dry-run first to not mess with partially consistent state.
	cmd := osutil.Command("patch", "-p1", "--force", "--ignore-whitespace", "--dry-run")
	if err := osutil.Sandbox(cmd, true, true); err != nil {
		return err
	}
	cmd.Stdin = bytes.NewReader(patch)
	cmd.Dir = dir
	if output, err := cmd.CombinedOutput(); err != nil {
		// If it reverses clean, then it's already applied
		// (seems to be the easiest way to detect it).
		cmd = osutil.Command("patch", "-p1", "--force", "--ignore-whitespace", "--reverse", "--dry-run")
		if err := osutil.Sandbox(cmd, true, true); err != nil {
			return err
		}
		cmd.Stdin = bytes.NewReader(patch)
		cmd.Dir = dir
		if _, err := cmd.CombinedOutput(); err == nil {
			return fmt.Errorf("patch is already applied")
		}
		return fmt.Errorf("failed to apply patch:\n%s", output)
	}
	// Now apply for real.
	cmd = osutil.Command("patch", "-p1", "--force", "--ignore-whitespace")
	if err := osutil.Sandbox(cmd, true, true); err != nil {
		return err
	}
	cmd.Stdin = bytes.NewReader(patch)
	cmd.Dir = dir
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to apply patch after dry run:\n%s", output)
	}
	return nil
}

// CheckRepoAddress does a best-effort approximate check of a git repo address.
func CheckRepoAddress(repo string) bool {
	return gitLocalRepoRe.MatchString(repo) ||
		gitRemoteRepoRe.MatchString(repo) ||
		gitSSHRepoRe.MatchString(repo)
}

// CheckBranch does a best-effort approximate check of a git branch name.
func CheckBranch(branch string) bool {
	return gitBranchRe.MatchString(branch)
}

func CheckCommitHash(hash string) bool {
	return gitHashRe.MatchString(hash)
}

func ParseReleaseTag(tag string) (v1, v2, rc, v3 int) {
	invalid := func() {
		v1, v2, rc, v3 = -1, -1, -1, -1
	}
	invalid()
	matches := releaseTagRe.FindStringSubmatch(tag)
	if matches == nil {
		return
	}
	for ptr, idx := range map[*int]int{
		&v1: 1, &v2: 2, &rc: 3, &v3: 4,
	} {
		if matches[idx] == "" {
			continue
		}
		var err error
		*ptr, err = strconv.Atoi(matches[idx])
		if err != nil {
			invalid()
			return
		}
	}
	return
}

func runSandboxed(dir, command string, args ...string) ([]byte, error) {
	cmd := osutil.Command(command, args...)
	cmd.Dir = dir
	if err := osutil.Sandbox(cmd, true, false); err != nil {
		return nil, err
	}
	return osutil.Run(time.Hour, cmd)
}

var (
	// nolint: lll
	gitLocalRepoRe = regexp.MustCompile(`^file:///[a-zA-Z0-9-_./~]+(/)?$`)
	// nolint: lll
	gitRemoteRepoRe = regexp.MustCompile(`^(git|ssh|http|https|ftp|ftps|sso)://[a-zA-Z0-9-_.]+(:[0-9]+)?(/[a-zA-Z0-9-_./~]+)?(/)?$`)
	// nolint: lll
	gitSSHRepoRe = regexp.MustCompile(`^(git|ssh|http|https|ftp|ftps|sso)@[a-zA-Z0-9-_.]+(:[a-zA-Z0-9-_]+)?(/[a-zA-Z0-9-_./~]+)?(/)?$`)
	gitBranchRe  = regexp.MustCompile("^[a-zA-Z0-9-_/.]{2,200}$")
	gitHashRe    = regexp.MustCompile("^[a-f0-9]{8,40}$")
	releaseTagRe = regexp.MustCompile(`^v([0-9]+).([0-9]+)(?:-rc([0-9]+))?(?:\.([0-9]+))?$`)
	// CC: is intentionally not on this list, see #1441.
	ccRes = []*regexp.Regexp{
		regexp.MustCompile(`^Reviewed\-.*: (.*)$`),
		regexp.MustCompile(`^[A-Za-z-]+\-and\-[Rr]eviewed\-.*: (.*)$`),
		regexp.MustCompile(`^Acked\-.*: (.*)$`),
		regexp.MustCompile(`^[A-Za-z-]+\-and\-[Aa]cked\-.*: (.*)$`),
		regexp.MustCompile(`^Tested\-.*: (.*)$`),
		regexp.MustCompile(`^[A-Za-z-]+\-and\-[Tt]ested\-.*: (.*)$`),
		regexp.MustCompile(`^Signed-off-by: (.*)$`),
	}
)

// CanonicalizeCommit returns commit title that can be used when checking
// if a particular commit is present in a git tree.
// Some trees add prefixes to commit titles during backporting,
// so we want e.g. commit "foo bar" match "BACKPORT: foo bar".
func CanonicalizeCommit(title string) string {
	for _, prefix := range commitPrefixes {
		if strings.HasPrefix(title, prefix) {
			title = title[len(prefix):]
			break
		}
	}
	return strings.TrimSpace(title)
}

var commitPrefixes = []string{
	"UPSTREAM:",
	"CHROMIUM:",
	"FROMLIST:",
	"BACKPORT:",
	"FROMGIT:",
	"net-backports:",
}

const SyzkallerRepo = "https://github.com/google/syzkaller"

const HEAD = "HEAD"

func CommitLink(url, hash string) string {
	return link(url, hash, "", 0, 0)
}

// Used externally - do not remove.
func TreeLink(url, hash string) string {
	return link(url, hash, "", 0, 1)
}

func LogLink(url, hash string) string {
	return link(url, hash, "", 0, 2)
}

func FileLink(url, hash, file string, line int) string {
	return link(url, hash, file, line, 3)
}

// nolint: goconst
func link(url, hash, file string, line, typ int) string {
	if url == "" || hash == "" {
		return ""
	}
	switch url {
	case "https://fuchsia.googlesource.com":
		// We collect hashes from the fuchsia repo.
		return link(url+"/fuchsia", hash, file, line, typ)
	}
	if strings.HasPrefix(url, "https://github.com/") {
		url = strings.TrimSuffix(url, ".git")
		switch typ {
		case 1:
			return url + "/tree/" + hash
		case 2:
			return url + "/commits/" + hash
		case 3:
			return url + "/blob/" + hash + "/" + file + "#L" + fmt.Sprint(line)
		default:
			return url + "/commit/" + hash
		}
	}
	if strings.HasPrefix(url, "https://git.kernel.org/pub/scm/") ||
		strings.HasPrefix(url, "git://git.kernel.org/pub/scm/") {
		url = strings.TrimPrefix(url, "git")
		url = strings.TrimPrefix(url, "https")
		url = "https" + url
		switch typ {
		case 1:
			return url + "/tree/?id=" + hash
		case 2:
			return url + "/log/?id=" + hash
		case 3:
			return url + "/tree/" + file + "?id=" + hash + "#n" + fmt.Sprint(line)
		default:
			return url + "/commit/?id=" + hash
		}
	}
	for _, cgitHost := range []string{"git.kernel.dk", "git.breakpoint.cc"} {
		if strings.HasPrefix(url, "https://"+cgitHost) ||
			strings.HasPrefix(url, "git://"+cgitHost) {
			url = strings.TrimPrefix(strings.TrimPrefix(url, "git://"), "https://")
			url = strings.TrimPrefix(url, cgitHost)
			url = "https://" + cgitHost + "/cgit" + url
			switch typ {
			case 1:
				return url + "/tree/?id=" + hash
			case 2:
				return url + "/log/?id=" + hash
			case 3:
				return url + "/tree/" + file + "?id=" + hash + "#n" + fmt.Sprint(line)
			default:
				return url + "/commit/?id=" + hash
			}
		}
	}
	if strings.HasPrefix(url, "https://") && strings.Contains(url, ".googlesource.com") {
		switch typ {
		case 1:
			return url + "/+/" + hash + "/"
		case 2:
			return url + "/+log/" + hash
		case 3:
			return url + "/+/" + hash + "/" + file + "#" + fmt.Sprint(line)
		default:
			return url + "/+/" + hash + "^!"
		}
	}
	return ""
}
