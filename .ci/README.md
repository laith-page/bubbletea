# Own runner validation

Run `mise run check` before pushing; the installed pre-push hook validates workflow syntax and static L3 routing. The build workflow then tests the root, examples, and tutorials modules and requires every build, lint, and coverage job through `Check`.

Repository Actions must be enabled before pushing the branch. Enabling Actions does not replay earlier pull request events: push a reviewed change and verify the run uses the exact new commit. Native Linux ARM64 Go 1.27.3 build, tests, and race tests passed for all three modules on 2026-11-11; the workflow's other versions and checks still need their own results.
