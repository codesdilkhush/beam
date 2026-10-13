<script lang="ts">
  import type { EmojiEntry } from '../types/index.js';

  interface Props {
    items: EmojiEntry[];
    selectedIndex: number;
    searchQuery: string;
    gridElement?: HTMLDivElement | null;
    onSelect: (index: number) => void;
    onHover: (index: number) => void;
  }

  let {
    items,
    selectedIndex,
    searchQuery,
    gridElement = $bindable(null),
    onSelect,
    onHover
  }: Props = $props();
</script>

<div bind:this={gridElement} class="emoji-container">
  {#if items.length === 0}
    <div class="empty-state">
      <p>No emojis found for "{searchQuery}"</p>
      <span class="empty-hint">Try searching by keyword like "smile", "fire", "cat", "heart"</span>
    </div>
  {:else}
    <div class="section-header">
      <span class="section-title">Results</span>
      <span class="section-count">{items.length} emojis</span>
    </div>

    <div class="emoji-grid">
      {#each items as item, idx}
        {@const isSelected = idx === selectedIndex}
        <button
          type="button"
          class="emoji-tile {isSelected ? 'selected' : ''}"
          title="{item.name} ({item.shortcode})"
          onclick={() => onSelect(idx)}
          onmouseenter={() => onHover(idx)}
        >
          <span class="emoji-glyph">{item.emoji}</span>
        </button>
      {/each}
    </div>
  {/if}
</div>

<style>
  .emoji-container {
    flex: 1;
    overflow-y: auto;
    padding: 14px 16px;
    position: relative;
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0 4px 4px 4px;
  }

  .section-title {
    font-size: 11.5px;
    font-weight: 600;
    color: var(--text-muted);
    text-transform: capitalize;
    letter-spacing: 0.3px;
  }

  .section-count {
    font-size: 11px;
    color: var(--text-muted);
    opacity: 0.7;
  }

  .emoji-grid {
    display: grid;
    grid-template-columns: repeat(8, 1fr);
    gap: 10px;
    width: 100%;
  }

  .emoji-tile {
    aspect-ratio: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(255, 255, 255, 0.04);
    border: 1.5px solid transparent;
    border-radius: 10px;
    cursor: pointer;
    transition: background 0.12s ease, border-color 0.12s ease;
    user-select: none;
    padding: 0;
  }

  .emoji-tile:hover {
    background: rgba(255, 255, 255, 0.08);
  }

  .emoji-tile.selected {
    background: rgba(255, 255, 255, 0.12);
    border-color: rgba(255, 255, 255, 0.85);
    box-shadow: 0 0 10px rgba(255, 255, 255, 0.12);
  }

  .emoji-glyph {
    font-size: 30px;
    line-height: 1;
    pointer-events: none;
    font-family: "Apple Color Emoji", "Segoe UI Emoji", "Noto Color Emoji", sans-serif;
  }

  .empty-state {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    text-align: center;
    padding: 40px 20px;
    color: var(--text-muted);
    gap: 6px;
  }

  .empty-state p {
    font-size: 14px;
    color: var(--text-secondary);
  }

  .empty-hint {
    font-size: 12px;
    opacity: 0.8;
  }
</style>
