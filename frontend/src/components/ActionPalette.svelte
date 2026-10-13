<script lang="ts">
  import type { ActionOption } from '../types/index.js';

  interface Props {
    show: boolean;
    actions: ActionOption[];
    selectedIndex: number;
    listElement?: HTMLDivElement | null;
    onSelect: (action: ActionOption) => void;
    onHover: (index: number) => void;
    onClose: () => void;
  }

  let {
    show,
    actions,
    selectedIndex,
    listElement = $bindable(null),
    onSelect,
    onHover,
    onClose
  }: Props = $props();
</script>

{#if show}
  <div
    class="palette-backdrop"
    role="button"
    tabindex="0"
    aria-label="Close actions"
    onclick={onClose}
    onkeydown={(e) => { if (e.key === 'Escape') onClose(); }}
  >
    <div
      class="palette-modal"
      role="dialog"
      aria-label="Actions Menu"
      tabindex="-1"
      onclick={(e) => e.stopPropagation()}
      onkeydown={(e) => e.stopPropagation()}
    >
      <div class="palette-header">
        <span class="palette-title">Actions</span>
        <span class="palette-sub">Press Esc to dismiss</span>
      </div>
      <div bind:this={listElement} class="palette-list">
        {#each actions as action, aIdx}
          {@const isActSelected = aIdx === selectedIndex}
          <button
            type="button"
            class="palette-item {isActSelected ? 'selected' : ''}"
            onclick={() => onSelect(action)}
            onmouseenter={() => onHover(aIdx)}
          >
            <span class="action-label">{action.label}</span>
            <span class="action-key"><kbd>{action.shortcut}</kbd></span>
          </button>
        {/each}
      </div>
    </div>
  </div>
{/if}

<style>
  .palette-backdrop {
    position: absolute;
    inset: 0;
    background: rgba(0, 0, 0, 0.65);
    backdrop-filter: blur(8px);
    -webkit-backdrop-filter: blur(8px);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 50;
  }

  .palette-modal {
    width: 480px;
    background: var(--surface-header);
    border: 1px solid var(--border-glass);
    border-radius: 12px;
    box-shadow: 0 20px 40px rgba(0, 0, 0, 0.7);
    overflow: hidden;
  }

  .palette-header {
    padding: 12px 16px;
    border-bottom: 1px solid var(--border-subtle);
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .palette-title {
    font-size: 13.5px;
    font-weight: 600;
    color: var(--text-primary);
  }

  .palette-sub {
    font-size: 11px;
    color: var(--text-muted);
  }

  .palette-list {
    max-height: 260px;
    overflow-y: auto;
    padding: 6px;
    display: flex;
    flex-direction: column;
    gap: 2px;
    position: relative;
    scroll-behavior: auto;
  }

  .palette-item {
    width: 100%;
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px 12px;
    border-radius: 6px;
    border: none;
    background: transparent;
    color: var(--text-primary);
    font-size: 13px;
    cursor: pointer;
    text-align: left;
    transition: background 0.1s ease;
    outline: none;
  }

  .palette-item:hover {
    background: var(--surface-hover);
  }

  .palette-item.selected {
    background: var(--surface-selected);
    color: var(--accent-sky);
    font-weight: 500;
  }

  .action-key {
    flex-shrink: 0;
  }

  .action-key kbd {
    background: rgba(255, 255, 255, 0.08);
    border: 1px solid rgba(255, 255, 255, 0.15);
    border-radius: 4px;
    padding: 2px 6px;
    font-size: 11px;
  }
</style>
