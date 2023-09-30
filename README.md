# gitglean action

Automatically generate a textual profile README.md with your GitHub info.

## Inputs

### `username`

**Required** Username to use in GitHub API requests.

### `template-path`

**Required** Path to template file. It must follow Go's [text/template](https://pkg.go.dev/text/template) rules.

### `output-path`

**Required** Path to output file.

## Example usage

``` yaml
uses: actions/checkout@v3
uses: IsaacDennis/gitglean@<version>
with:
  username: IsaacDennis
  template-path: ${{ github.workspace }}/README.template
  output-path: ${{ github.workspace }}/README.md
```
