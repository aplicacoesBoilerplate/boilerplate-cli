"""Boundary contracts for centralized Release Please; no GitHub writes."""

from pathlib import Path
import unittest

import yaml

ROOT = Path(__file__).resolve().parents[2]
REVISION = "1d874d946e5f5ad233a6b1f6bf278fef78f3b46d"
MERGED_GUARD = "github.event.pull_request.merged == true"


class CallerContract(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.workflow = yaml.load(
            (ROOT / ".github/workflows/ci.yml").read_text(encoding="utf-8-sig"),
            Loader=yaml.BaseLoader,
        )
        cls.caller = yaml.load(
            (ROOT / ".github/workflows/release.yml").read_text(encoding="utf-8-sig"),
            Loader=yaml.BaseLoader,
        )

    def test_published_workflow_and_scripts_use_same_reviewed_revision(self):
        job = self.caller["jobs"]["release"]
        self.assertEqual(job["uses"], f"aplicacoesBoilerplate/.github/.github/workflows/version-release-please-called.yml@{REVISION}")
        self.assertEqual(job["with"]["central-ref"], REVISION)

    def test_versioning_selects_go_without_binary_distribution(self):
        self.assertEqual(self.caller["jobs"]["release"]["with"], {
            "project-type": "go", "central-ref": REVISION, "publish-artifacts": "false",
        })

    def test_only_closed_master_prs_can_start_release(self):
        self.assertEqual(self.caller["on"], {
            "pull_request": {"types": ["closed"], "branches": ["master"]},
        })
        for job in self.caller["jobs"].values():
            self.assertEqual(job["if"], MERGED_GUARD)

    def test_validation_is_read_only_and_uses_application_ci(self):
        self.assertEqual(self.caller["permissions"], {"contents": "read"})
        self.assertEqual(self.caller["jobs"]["validate"], {
            "if": MERGED_GUARD, "uses": "./.github/workflows/ci.yml",
        })
        self.assertIn("workflow_call", self.workflow["on"])

    def test_release_requires_ci_and_minimum_writer_permissions(self):
        job = self.caller["jobs"]["release"]
        self.assertEqual(job["needs"], "validate")
        self.assertEqual(job["permissions"], {
            "contents": "write", "pull-requests": "write", "issues": "write",
        })
        self.assertEqual(self.workflow["permissions"], {"contents": "read"})

    def test_legacy_preview_and_publisher_are_retired(self):
        self.assertEqual(set(self.workflow["jobs"]), {"go-ci"})
        workflows = "\n".join(p.read_text(encoding="utf-8-sig") for p in (ROOT / ".github/workflows").glob("*.yml"))
        for legacy in ["version-preview.yml@", "version-publish.yml@", "adapter: go-gitsemver", "goreleaser/goreleaser-action@", "packages: write"]:
            self.assertNotIn(legacy, workflows)

    def test_ci_keeps_release_develop_master_and_pilot_branches(self):
        self.assertEqual(self.workflow["on"]["push"]["branches"], ["master", "develop", "release/**"])
        self.assertEqual(self.workflow["on"]["pull_request"]["branches"], ["master", "develop", "release/**", "feature/issue-8"])
        self.assertEqual(self.workflow["on"]["pull_request"]["types"], [
            "opened", "reopened", "synchronize", "edited", "labeled", "unlabeled", "milestoned", "demilestoned",
        ])

    def test_ci_inline_and_gated_on_go_mod_not_hardcoded_go_version(self):
        job = self.workflow["jobs"]["go-ci"]
        self.assertNotIn("uses", job)
        self.assertEqual(job["if"], "${{ github.event_name != 'pull_request_review' }}")
        steps = job["steps"]
        checkout = next(step for step in steps if step.get("uses", "").startswith("actions/checkout@"))
        self.assertEqual(checkout["with"]["ref"], "${{ github.event.pull_request.merged == true && github.event.pull_request.merge_commit_sha || github.event.pull_request.head.sha || github.sha }}")
        setup = next(step for step in steps if step.get("uses", "").startswith("actions/setup-go@"))
        self.assertEqual(setup["with"]["go-version-file"], "go.mod")
        self.assertNotIn("go-version", setup["with"])
        runs = [step["run"] for step in steps if "run" in step]
        for command in ["go test ./...", "go vet ./...", "go build -o boilerplate .", "./boilerplate --help", "python tests/versioning/test_caller.py"]:
            self.assertIn(command, runs)
        self.assertTrue(any(step.get("uses", "").startswith("golangci/golangci-lint-action@") for step in steps))
        self.assertTrue(any("scripts/validate-commit-msg.sh" in run for run in runs))

    def test_caller_has_no_application_steps_secrets_or_distribution_dependency(self):
        self.assertEqual(set(self.caller["jobs"]), {"validate", "release"})
        self.assertEqual(set(self.caller["jobs"]["release"]), {"if", "needs", "permissions", "uses", "with"})
        self.assertNotIn("secrets", self.caller["jobs"]["validate"])
        self.assertFalse((ROOT / ".github/workflows/go-ci.yml").exists())
        self.assertTrue((ROOT / ".goreleaser.yml").is_file())


if __name__ == "__main__":
    unittest.main()
