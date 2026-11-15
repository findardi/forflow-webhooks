---
name: Forflow
description: Webhook delivery dashboard drawn as a quiet Braun instrument panel.
colors:
  key: "#237a45"
  key-press: "#1b6237"
  ground: "#f3f3f1"
  panel: "#ffffff"
  panel-sunk: "#f7f7f5"
  rule: "#e2e2df"
  rule-soft: "#eeeeeb"
  field: "#cfcfca"
  ink: "#2b2b2b"
  ink-2: "#5e5e5a"
  ink-3: "#6f6f6b"
  grille: "#c9c9c4"
  error: "#a23b12"
  error-dot: "#e0662a"
typography:
  headline:
    fontFamily: "Hanken Grotesk, system-ui, sans-serif"
    fontSize: "30px"
    fontWeight: 600
    lineHeight: 1.2
  title:
    fontFamily: "Hanken Grotesk, system-ui, sans-serif"
    fontSize: "18px"
    fontWeight: 600
  body:
    fontFamily: "Hanken Grotesk, system-ui, sans-serif"
    fontSize: "15px"
    fontWeight: 400
    lineHeight: 1.5
  label:
    fontFamily: "Hanken Grotesk, system-ui, sans-serif"
    fontSize: "12px"
    fontWeight: 600
    letterSpacing: "0.08em"
  mono:
    fontFamily: "JetBrains Mono, ui-monospace, monospace"
    fontSize: "13px"
    fontWeight: 400
rounded:
  control: "6px"
  panel: "10px"
spacing:
  page: "40px"
  gutter: "24px"
  panel-x: "28px"
  list-x: "22px"
  mobile: "16px"
components:
  button-key:
    backgroundColor: "{colors.key}"
    textColor: "{colors.panel}"
    rounded: "{rounded.control}"
    height: "42px"
    padding: "0 18px"
  button-key-hover:
    backgroundColor: "{colors.key-press}"
  button-ghost:
    textColor: "{colors.ink}"
    rounded: "{rounded.control}"
    height: "42px"
    padding: "0 18px"
  button-ghost-hover:
    backgroundColor: "{colors.panel-sunk}"
  field:
    backgroundColor: "{colors.panel}"
    textColor: "{colors.ink}"
    rounded: "{rounded.control}"
    height: "42px"
    padding: "0 12px"
  panel:
    backgroundColor: "{colors.panel}"
    rounded: "{rounded.panel}"
---

# Design System: Forflow

## Overview

**Creative North Star: "Braun Panel, less colour"**

A quiet instrument face after Dieter Rams' Braun product graphics: matte gray ground, white panels, graphite type, and one green key. Dense enough to see every application at once, calm enough to read as equipment rather than marketing.

**Key Characteristics:**
- One coloured element per screen: the primary action.
- Flat white panels on warm gray, separated by 1px rules.
- Small uppercase labels name panels and table columns.
- A dot grille is the only ornament.

## Colors

Graphite and warm grays, with a single green key and an orange error signal.

### Primary
- **Key Green** (`key`): the one primary action per screen (e.g. "Add endpoint"). Darkens to **Key Pressed** (`key-press`) on hover. Also the native `accent-color`.

### Neutral
- **Matte Ground** (`ground`): page background and the selected list row.
- **Panel White** (`panel`): panels, fields. **Sunk Panel** (`panel-sunk`): row and ghost-button hover.
- **Warm Rule** (`rule`): panel borders and section dividers. **Soft Rule** (`rule-soft`): row dividers inside lists and tables. **Field Stroke** (`field`): input borders.
- **Graphite** (`ink`): text, ghost-button stroke, focus ring. **Graphite 2** (`ink-2`): labels, metadata, secondary nav. **Graphite 3** (`ink-3`): placeholders, field hover stroke.
- **Grille Dot** (`grille`): the dot grille and scrollbar thumb.

### Signal
- **Error Rust** (`error`) text with an **Error Dot** (`error-dot`) 8px circle before it; also the invalid-field border.

### Named Rules
**The One Key Rule.** Green appears on exactly one control per screen, the primary action. List "Create" holds the key on the index screen, where it is the only action, and becomes a graphite ghost button once an application is open and "Add endpoint" takes the key. Nothing else is green.

## Typography

**Body Font:** Hanken Grotesk (with system-ui), weights 400/500/600.
**Mono Font:** JetBrains Mono (13px, ligatures off).

**Character:** A plain grotesk doing all the talking; mono only where the characters themselves are the data.

### Hierarchy
- **Headline** (600, 30px, 1.2; 24px under 800px): the selected record's name.
- **Title** (600, 18px): wordmark.
- **Body** (400, 15px, 1.5): everything else; 600 marks current nav and selected row.
- **Label** (600, 12px, 0.08em, uppercase, Graphite 2): panel headers and table column headers.
- **Mono** (400, 13px): URLs and IDs only.

### Named Rules
**The Mono-Is-Data Rule.** JetBrains Mono is used for URLs and IDs and nothing else.
**The UTC Rule.** Dates render in UTC and say so ("Created (UTC)", "… UTC"). Counts and times use tabular numerals.

## Layout

64px top bar (40px side padding), main padded 28px 40px 40px. Applications use a two-column split: list `minmax(280px, 380px)` and detail `minmax(0, 1fr)`, 24px gap, max 1600px centered. Panels pad 22px (list) and 28px (detail) horizontally. Fields and their button sit in one row with a 10px gap.

Under 800px: single column, 16px padding. When an application is selected the detail replaces the list and shows a back link; field rows stack vertically.

## Elevation & Depth

Flat. No shadows. Depth comes from tone (white panel on matte ground) and 1px rules.

## Shapes

Modest rounding: 10px on panels, 6px on controls. Borders are always 1px. Error dots are full circles.

## Components

### Buttons
- **Shape:** 42px tall, 6px radius, 600 15px text, 18px side padding.
- **Key:** green fill, white text; the one primary action.
- **Ghost:** 1px graphite stroke, transparent; Sunk Panel on hover.
- **Danger:** Error Rust fill, white text, darkens to `#872f0c` on hover. Only at the last step of a delete, never as a resting control. Disabled shows `not-allowed`.
- **States:** press nudges down 1px (80ms); disabled drops to 0.6 opacity with a progress cursor. Transitions are off under reduced motion.

### Inputs / Fields
- **Style:** 42px, 1px Field Stroke, white, 6px radius. Hover strokes Graphite 3; invalid strokes Error Rust.
- **Focus:** 2px graphite outline, 2px offset (global focus-visible).
- **Error:** plain-words message under the field with the orange dot.

### Cards / Containers
- **Panel:** white, 1px Warm Rule, 10px radius, no shadow. Header row carries a Label and a count.

### Lists and Tables
- Rows divided by Soft Rule; hover Sunk Panel. Table headers are Labels; the date column right-aligns.
- Selected list row: Matte Ground fill, name in 600. No side bar.

### Navigation
- Top bar: wordmark left, text links in Graphite 2; current page Graphite 600. 1px rule below.

### Dot Grille
- Radial dots (1.6px, Grille Dot) on an 8px pitch. Used only as the empty-state and error-page mark, decorative (`aria-hidden`).

### Destructive actions
- **Endpoint row:** a quiet "Delete" text button arms the row (Matte Ground fill). Confirm (Danger) appears to its left and Cancel takes the trigger's spot, so a double click cancels. Escape cancels.
- **Application:** a ghost "Delete…" sits beside "Copy ID" in the header. It opens a Sunk Panel strip under the header that names what goes with it (endpoint count) and asks for the exact name before the Danger button enables. The server checks the name too.

## Do's and Don'ts

### Do:
- **Do** give each screen exactly one green key, and make it the primary action.
- **Do** separate surfaces with 1px rules and tone, not shadows.
- **Do** set URLs and IDs in mono, everything else in Hanken Grotesk.
- **Do** show dates in UTC, labelled as such.
- **Do** pair every error with the orange dot and plain wording beside the field.

### Don't:
- **Don't** use green for links, status, selection, or decoration.
- **Don't** add ornament beyond the dot grille, and don't use the grille outside empty or error states.
- **Don't** add drop shadows or radii other than 10px and 6px.
