/**
 * Synchronizes the scroll position of a scrollable container
 * so that the selected child item remains inside the visible viewport.
 */
export function scrollSelectedIntoView(container: HTMLElement | null, selector: string, index: number): void {
  if (!container) return;
  if (index <= 0) {
    container.scrollTop = 0;
    return;
  }
  const items = container.querySelectorAll(selector);
  const item = items[index] as HTMLElement | undefined;
  if (!item) return;

  const containerTop = container.scrollTop;
  const containerBottom = containerTop + container.clientHeight;
  const itemTop = item.offsetTop;
  const itemBottom = itemTop + item.offsetHeight;

  if (itemTop < containerTop) {
    container.scrollTop = Math.max(0, itemTop - 10);
  } else if (itemBottom > containerBottom) {
    container.scrollTop = itemBottom - container.clientHeight + 10;
  }
}
