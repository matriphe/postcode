# Postcode rule data

`rules.jsonc` contains country-level postcode format patterns, examples, and normalization/display replacements. Its Google-derived patterns and examples come from the [Google libaddressinput metadata service](https://chromium-i18n.appspot.com/ssl-address/data/) and are licensed under [CC BY 4.0](https://creativecommons.org/licenses/by/4.0/). The Google project identifies its data license in its [repository README](https://github.com/google/libaddressinput).

Each country entry includes the metadata URL and verification date. Luxembourg's prefix forms and the United Kingdom's postcode reference also cite official postal guidance. Norway's three-digit shorthand is a package-specific normalization rule; Google's metadata describes the four-digit postcode format.

The snapshot was collected on 2026-10-21. It includes 179 ISO alpha-2 country and territory entries for which Google metadata supplied a postcode pattern. Countries without a usable pattern are omitted rather than assigned a guessed rule. Matching validates format only, not whether a code is assigned to an address.
