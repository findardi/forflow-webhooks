---
version: 1
slug: "web-src-routes"
primary_target: "web/src/routes"
related_targets: []
---

# Dashboard (web/src/routes)

Mode: Operate. Audience: developer integrators and the project owner. Task: create and list applications; create and list endpoints per application. Only the four real APIs; no auth yet.

## Direction contract

THESIS: Braun Panel, less colour. A quiet instrument face where the only colour on screen is the one key you can press. Refuses the gray-sidebar-plus-blue-button category default and its dark neon opposite.

OWN-WORLD: Light matte gray ground (#F3F3F1), white panels with 1px warm-gray rules and modest radius, graphite text. Hanken Grotesk for UI, JetBrains Mono only for URLs and IDs. Small uppercase panel labels. One green key (#237A45) per screen for the primary action; error is a small orange dot plus dark orange text. A fine dot grille is the only ornament and doubles as the empty state.

STORY: The developer sees every application at once, opens one, sees its endpoints, adds one in a single field. Errors from the API read in plain words beside the field.

FIRST VIEWPORT: 64px top bar (wordmark, Applications, org name). Below, two panels: left 380px application list with counts and an inline create field; right the selected application: name, full ID and created date, endpoint field with the green key, endpoint table (URL, ID, created).

FORM: Braun / Rams product graphics, my top-ranked grounded candidate, chosen by the user over the roll. seed:34eb9f5a18ad6cf5

FINISH: unreviewed and undocumented is unfinished; this build ends with the finish review, the verdict, DESIGN.md, and every shipping raster carrying its provenance
