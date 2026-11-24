#!/usr/bin/env python3
"""Mechanical compile-risk gate for generated Go files.

Two classes of defect that the Go compiler / `go vet` either do not report
(unused unexported symbols) or only report after the whole package is
assembled (a name declared twice by two generated files of one product):

  1. UNUSED  - an unexported top-level func / const / var / type that nothing
               else in the package references. The compiler ignores it, only
               staticcheck U1000 reports it, and staticcheck is not guaranteed
               to be installed. Generated code must never contain one.
  2. COLLIDE - a top-level name declared in a checked file that is ALSO
               declared by another file of the same package (the product
               directory is shared by every generated resource and by every
               external `<product>_test` test file, so a bare helper name such
               as `ptrStr` collides as soon as a second resource is generated).

Usage (repository root):

    python3 tencentcloud/framework/templates/check_go_symbols.py FILE.go [FILE.go ...]

Every FILE is checked against the other *.go files of its own directory,
grouped by their `package` clause (so `demo` and `demo_test` stay separate).
Exit code 0 = clean, 1 = at least one problem (all are printed).
"""
import pathlib
import re
import sys

_FUNC = re.compile(r"^func ([A-Za-z_]\w*)\s*[\[(]", re.M)
_SINGLE = re.compile(r"^(?:const|var|type) ([A-Za-z_]\w*)\b", re.M)
_BLOCK = re.compile(r"^(?:const|var|type) \(\n(.*?)^\)", re.M | re.S)
_ITEM = re.compile(r"^\t([A-Za-z_]\w*)\b", re.M)
_PACKAGE = re.compile(r"^package (\w+)", re.M)
_COMMENT = re.compile(r"//[^\n]*")


def strip_comments(src: str) -> str:
    return _COMMENT.sub("", src)


def top_level_names(src: str) -> set:
    names = set(_FUNC.findall(src)) | set(_SINGLE.findall(src))
    for block in _BLOCK.findall(src):
        names.update(_ITEM.findall(block))
    names.discard("_")
    return names


def package_of(src: str) -> str:
    m = _PACKAGE.search(src)
    return m.group(1) if m else ""


def check(files):
    """Return a list of human-readable problems for the given files."""
    problems = []
    for f in files:
        f = pathlib.Path(f).resolve()
        src = strip_comments(f.read_text())
        pkg = package_of(src)
        siblings = [
            p for p in sorted(f.parent.glob("*.go"))
            if p.resolve() != f and package_of(strip_comments(p.read_text())) == pkg
        ]
        sibling_src = {p: strip_comments(p.read_text()) for p in siblings}
        corpus = src + "\n" + "\n".join(sibling_src.values())
        mine = top_level_names(src)

        for name in sorted(mine):
            # UNUSED: unexported symbol referenced nowhere but its declaration.
            if name[0].islower() and name not in ("init", "main"):
                if len(re.findall(r"\b%s\b" % re.escape(name), corpus)) <= 1:
                    problems.append(f"{f.name}: UNUSED private symbol `{name}`")
            # COLLIDE: the same top-level name declared by a sibling file.
            for p, s in sibling_src.items():
                if name in top_level_names(s):
                    problems.append(f"{f.name}: COLLIDES with `{name}` already declared in {p.name}")
    return problems


def main(argv):
    if not argv:
        print(__doc__)
        return 2
    problems = check(argv)
    for p in problems:
        print(p)
    print(f"check_go_symbols: {len(problems)} problem(s) in {len(argv)} file(s)")
    return 1 if problems else 0


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
