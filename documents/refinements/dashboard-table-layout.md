# Dashboard Table Layout Refinement

## Changes

### Namespace column removed
The `Namespace` table column was removed since all data within a single query belongs to the same namespace. Namespace info is now displayed as a subtitle above the table (`table-card-header`).

### Column width adjustments
- **Image tag**: increased from `max-width: 250px` to `min-width: 350px` for better visibility of long image references.
- **Container name**: removed fixed width, naturally shrinks with monospace font.

### Action buttons inline
Env and Logs buttons now sit side by side instead of stacking vertically. Removed fixed `width: 90px` from `.cell-actions`, added `white-space: nowrap` and horizontal spacing.

### CSS cleanup
- Removed unused `.ns-badge` styles (namespace badge no longer rendered in table).
- Added `.table-card-header` and `.table-card-subtitle` for the new namespace subtitle.

### Template changes
- All `colspan="7"` updated to `colspan="6"` to reflect the removed column.
- JS `filterTableRows` no longer searches against namespace (removed `.cell-ns` reference).

## Files modified
- `front/table.html` — template and JS
- `front/static/page.css` — styles
