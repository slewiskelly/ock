---
title: "Schema Reference"
type: "reference"
lastReviewed: "2025-02-06"
---
# Schema Reference

## Errors

Custom error messages can be added to fields of a schema by using the `@error`
[attribute](https://cuelang.org/docs/reference/spec/#attributes). For example:

```cue
#Metadata: {
	title:  string & strings.MinRunes(1)   @error(a title is required)
	status: #Status                        @error(invalid status; must be either "archived", "draft", or "published")
	type:   #Type                          @error(invalid type; must be either "concept", "guide", or "reference")
	owner:  #Owner                         @error(invalid owner; must be either "alice", "bob", or "clyde")
	tags:   [...string] & list.MinItems(1) @error(must contain at least one tag)
}
```

The rationale behind this is that the errors reported by CUE can sometimes be
verbose, leading to some confusion as to how to resolve. A succint error message
may prove useful in helping an author update a page's metadata more easily than
first navigating an error provided by CUE directly.

The use of custom error messages is optional, if a field has no `@error`
attribute, or is unpopulated, then the native error reported via CUE's
validation will be presented.

## Warnings

Expressing warnings is identical to that of [errors](#errors), with the
exception that the `@warning` attribute is used in place of `@error`.

Note that if a field has neither an `@error` or `@warning` attribute, any
validation error will be considered an error.
