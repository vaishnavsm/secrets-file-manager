# Secrets File Manager

A tool to manage secrets files easily in a git repo.

The intended workflow is:

1. Add config files containing secrets to gitignore (or use the `gen-gitignore` command)
2. Run `secrets-file-manager gen-gitignore >> .gitignore && secrets-file-manager sync` on git commit / post-merge hooks
3. Commit the crypt files to git

There are lots of ways to do this already.
This is just an extremely simple implementation.

## Usage

```bash
secrets-file-manager init > config.secrets-file-manager.yaml
# edit paths
secrets-file-manager gen-gitignore >> .gitignore
secrets-file-manager sync
```
