"""Pilot boundary contracts; --real-go evaluates native bootstrap without GitHub writes."""

import json
from pathlib import Path
import shutil
import subprocess
import sys
import tempfile
import unittest

import yaml

ROOT = Path(__file__).resolve().parents[2]
REVISION = "f3e7a0168da9e3882840302add3ac8d162bac81e"
REAL_GO = "--real-go" in sys.argv
if REAL_GO:
    sys.argv.remove("--real-go")

PR_GUARD = "${{ (github.event_name == 'pull_request' || github.event_name == 'pull_request_review') && (github.event.pull_request.base.ref == 'develop' || github.event.pull_request.base.ref == 'master') }}"
PUBLISH_GUARD = "${{ github.event_name == 'push' && github.ref_name == 'master' }}"


class CallerContract(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.workflow = yaml.load(
            (ROOT / ".github/workflows/ci.yml").read_text(encoding="utf-8-sig"),
            Loader=yaml.BaseLoader,
        )

    def test_published_workflows_selected_at_same_reviewed_revision(self):
        for job, workflow in [("validate-pr", "version-preview.yml"), ("publish", "version-publish.yml")]:
            self.assertEqual(
                self.workflow["jobs"][job]["uses"],
                f"aplicacoesBoilerplate/.github/.github/workflows/{workflow}@{REVISION}",
            )

    def test_versioning_inputs_only_describe_adapter_project_and_target(self):
        for job in ["validate-pr", "publish"]:
            self.assertEqual(self.workflow["jobs"][job]["with"], {
                "adapter": "go-gitsemver", "project_path": ".", "target_branch": "master",
            })

    def test_preview_does_not_fail_feature_to_release_ci(self):
        self.assertEqual(self.workflow["jobs"]["validate-pr"]["if"], PR_GUARD)

    def test_preview_permissions_are_read_only_and_include_checks(self):
        self.assertEqual(self.workflow["jobs"]["validate-pr"]["permissions"], {
            "contents": "read", "pull-requests": "read", "issues": "read", "checks": "read",
        })

    def test_publication_only_after_master_push_and_successful_application_ci(self):
        job = self.workflow["jobs"]["publish"]
        self.assertEqual(job["if"], PUBLISH_GUARD)
        self.assertEqual(job["needs"], "go-ci")
        self.assertEqual(job["permissions"], {
            "contents": "write", "pull-requests": "read", "issues": "read",
        })

    def test_all_policy_invalidation_events_subscribed(self):
        self.assertEqual(self.workflow["on"]["pull_request"]["types"], [
            "opened", "reopened", "synchronize", "edited", "labeled", "unlabeled", "milestoned", "demilestoned",
        ])
        self.assertEqual(self.workflow["on"]["pull_request_review"]["types"], ["submitted", "dismissed"])

    def test_ci_keeps_release_develop_master_and_review_events(self):
        self.assertEqual(self.workflow["on"]["push"]["branches"], ["master", "develop", "release/**"])
        self.assertEqual(self.workflow["on"]["pull_request"]["branches"], ["master", "develop", "release/**", "feature/issue-8"])

    def test_ci_inline_and_gated_on_go_mod_not_hardcoded_go_version(self):
        job = self.workflow["jobs"]["go-ci"]
        self.assertNotIn("uses", job)
        self.assertEqual(job["if"], "${{ github.event_name != 'pull_request_review' }}")
        steps = job["steps"]
        setup = next(step for step in steps if step.get("uses", "").startswith("actions/setup-go@"))
        self.assertEqual(setup["with"]["go-version-file"], "go.mod")
        self.assertNotIn("go-version", setup["with"])
        runs = [step["run"] for step in steps if "run" in step]
        for command in ["go test ./...", "go vet ./...", "go build -o boilerplate .", "./boilerplate --help", "python tests/versioning/test_caller.py"]:
            self.assertIn(command, runs)
        self.assertTrue(any(step.get("uses", "").startswith("golangci/golangci-lint-action@") for step in steps))
        self.assertTrue(any("scripts/validate-commit-msg.sh" in run for run in runs))
        self.assertEqual(self.workflow["permissions"], {"contents": "read"})

    def test_versioning_has_no_application_steps_local_secret_or_distribution_dependency(self):
        for job_name in ["validate-pr", "publish"]:
            job = self.workflow["jobs"][job_name]
            self.assertEqual(set(job), {"if", "permissions", "uses", "with"} | ({"needs"} if job_name == "publish" else set()))
        self.assertEqual(set(self.workflow["jobs"]), {"go-ci", "validate-pr", "publish"})
        self.assertFalse((ROOT / ".github/workflows/go-ci.yml").exists())


def run_native_bootstrap():
    """Real pilot history in a clone avoids native-library linked-worktree quirks."""
    binary = shutil.which("go-gitsemver")
    if not binary:
        raise RuntimeError("Install the pinned go-gitsemver before --real-go")
    output = subprocess.check_output(["go", "version", "-m", binary], text=True)
    if "680c1c12d9a4" not in output:
        raise AssertionError("Native binary must match pinned revision 680c1c12d9a4")
    sha = subprocess.check_output(["git", "rev-parse", "HEAD"], cwd=ROOT, text=True).strip()
    with tempfile.TemporaryDirectory(prefix="boilerplate-versioning-") as directory:
        repo = Path(directory) / "pilot"
        subprocess.run(["git", "clone", "--quiet", "--no-hardlinks", str(ROOT), str(repo)], check=True)
        subprocess.run(["git", "checkout", "--quiet", "--detach", sha], cwd=repo, check=True)
        subprocess.run(["git", "branch", "-f", "master", sha], cwd=repo, check=True)
        shutil.copyfile(ROOT / "go-gitsemver.yml", repo / "go-gitsemver.yml")
        before = subprocess.check_output(["git", "tag"], cwd=repo, text=True)
        result = subprocess.run([binary, "--branch", "master", "--commit", sha, "-o", "json", "--explain"], cwd=repo, text=True, capture_output=True, check=True)
        native = json.loads(result.stdout)
        if native["SemVer"] != "0.0.1" or native["Sha"] != sha:
            raise AssertionError(f"Expected stable 0.0.1 at {sha}, received {native}")
        if not result.stderr.strip():
            raise AssertionError("Native explanation must be available")
        if subprocess.check_output(["git", "tag"], cwd=repo, text=True) != before:
            raise AssertionError("Native calculation must not create tags")
        print(f"Real pinned adapter: 0.0.1 at {sha}; explanation present, tags unchanged")


if __name__ == "__main__":
    if REAL_GO:
        run_native_bootstrap()
    else:
        unittest.main()
