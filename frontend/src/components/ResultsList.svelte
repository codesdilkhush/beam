<script lang="ts">
  import type { AppItem, ClipboardEntry, TodoEntry } from '../types/index.js';

  interface Props {
    items: (AppItem | ClipboardEntry | TodoEntry)[];
    selectedIndex: number;
    isClipboardMode: boolean;
    isTodoMode: boolean;
    searchQuery: string;
    listElement?: HTMLDivElement | null;
    onSelect: (index: number) => void;
    onHover: (index: number) => void;
  }

  let {
    items,
    selectedIndex,
    isClipboardMode,
    isTodoMode,
    searchQuery,
    listElement = $bindable(null),
    onSelect,
    onHover
  }: Props = $props();
</script>

<div bind:this={listElement} class="results-list {isTodoMode ? 'full-width' : ''}">
  {#if items.length === 0}
    <div class="empty-state">
      {#if isTodoMode}
        <p>No tasks found {searchQuery ? `for "${searchQuery}"` : ''}</p>
        {#if searchQuery.trim()}
          <span class="empty-hint">Press <kbd>↵</kbd> to add "{searchQuery.trim()}" as a task</span>
        {:else}
          <span class="empty-hint">Type a task description and press <kbd>↵</kbd></span>
        {/if}
      {:else}
        <p>No results found for "{searchQuery}"</p>
        <span class="empty-hint">Press <kbd>↵</kbd> to search on Google</span>
      {/if}
    </div>
  {:else}
    {#each items as item, idx}
      {@const isSelected = idx === selectedIndex}
      <button
        type="button"
        class="list-item {isSelected ? 'selected' : ''}"
        onclick={() => onSelect(idx)}
        onmouseenter={() => onHover(idx)}
      >
        <div class="item-left">
          {#if isTodoMode}
            {@const todo = item as TodoEntry}
            <div class="icon-box">
              <div class="todo-checkbox {todo.completed ? 'checked' : ''}">
                {#if todo.completed}
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3" stroke-linecap="round" stroke-linejoin="round">
                    <polyline points="20 6 9 17 4 12"></polyline>
                  </svg>
                {/if}
              </div>
            </div>
            <div class="item-info clip-info">
              <span class="item-title truncate {todo.completed ? 'todo-done' : ''}">{todo.title}</span>
            </div>
          {:else if isClipboardMode}
            {@const clip = item as ClipboardEntry}
            <div class="icon-box">
              <div class="icon-clipboard">
                <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"></path>
                  <rect x="8" y="2" width="8" height="4" rx="1" ry="1"></rect>
                </svg>
              </div>
            </div>
            <div class="item-info clip-info">
              <span class="item-title truncate">{clip.preview}</span>
            </div>
          {:else}
            {@const app = item as AppItem}
            <div class="icon-box">
              {#if app.resolvedIcon}
                <img src={app.resolvedIcon} alt="" class="app-icon" />
              {:else}
                <div class="letter-icon">{app.name.charAt(0).toUpperCase()}</div>
              {/if}
            </div>
            <div class="item-info app-info">
              <span class="item-title truncate">{app.name}</span>
              {#if app.launchCount > 0}
                <span class="launch-badge" title="Launched {app.launchCount} times">
                  {app.launchCount} {app.launchCount === 1 ? 'launch' : 'launches'}
                </span>
              {/if}
            </div>
          {/if}
        </div>

        {#if isSelected}
          <div class="item-right">
            <span class="action-badge">
              ↵ {#if isTodoMode}Toggle{:else if isClipboardMode}Copy{:else}Open{/if}
            </span>
          </div>
        {/if}
      </button>
    {/each}
  {/if}
</div>

<style>
  .results-list {
    flex: 6;
    overflow-y: auto;
    padding: 10px;
    border-right: 1px solid var(--border-subtle);
    position: relative;
    scroll-behavior: auto;
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .results-list.full-width {
    flex: 1;
    border-right: none;
  }

  .list-item {
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 10px 14px;
    min-height: 48px;
    border-radius: 9px;
    border: 1px solid transparent;
    background: transparent;
    color: var(--text-primary);
    cursor: pointer;
    font: inherit;
    text-align: left;
    transition: background 0.1s ease, border-color 0.1s ease;
    outline: none;
  }

  .list-item:hover {
    background: var(--surface-hover);
  }

  .list-item.selected {
    background: var(--surface-selected);
    border-color: var(--border-selected);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
  }

  .item-left {
    display: flex;
    align-items: center;
    gap: 12px;
    min-width: 0;
    flex: 1;
    margin-right: 12px;
  }

  .icon-box {
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    border-radius: 6px;
  }

  .icon-clipboard {
    width: 28px;
    height: 28px;
    border-radius: 6px;
    background: linear-gradient(135deg, rgba(255, 255, 255, 0.12), rgba(255, 255, 255, 0.04));
    border: 1px solid rgba(255, 255, 255, 0.1);
    color: var(--accent-sky);
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: inset 0 1px 2px rgba(255, 255, 255, 0.1);
  }

  .app-icon {
    width: 28px;
    height: 28px;
    object-fit: contain;
  }

  .letter-icon {
    width: 28px;
    height: 28px;
    border-radius: 6px;
    background: linear-gradient(135deg, rgba(255, 255, 255, 0.12), rgba(255, 255, 255, 0.04));
    border: 1px solid rgba(255, 255, 255, 0.1);
    color: var(--accent-sky);
    font-weight: 700;
    font-size: 13px;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: inset 0 1px 2px rgba(255, 255, 255, 0.1);
  }

  .item-info {
    min-width: 0;
    flex: 1;
  }

  .app-info,
  .clip-info {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .todo-checkbox {
    width: 20px;
    height: 20px;
    border-radius: 5px;
    border: 1.5px solid var(--accent-sky);
    background: transparent;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #ffffff;
    transition: all 0.15s ease;
  }

  .todo-checkbox.checked {
    background: var(--accent-blue);
    border-color: var(--accent-blue);
  }

  .todo-done {
    text-decoration: line-through;
    color: var(--text-muted);
  }



  .item-title {
    font-size: 13.5px;
    font-weight: 500;
    color: var(--text-primary);
  }

  .launch-badge {
    font-size: 10px;
    font-weight: 600;
    background: rgba(56, 117, 246, 0.2);
    color: var(--accent-sky);
    padding: 1px 6px;
    border-radius: 10px;
    border: 1px solid rgba(56, 117, 246, 0.3);
    flex-shrink: 0;
  }

  .truncate {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .item-right {
    flex-shrink: 0;
    display: flex;
    align-items: center;
  }

  .action-badge {
    font-size: 11px;
    font-weight: 600;
    color: #ffffff;
    background: var(--accent-blue);
    padding: 3px 9px;
    border-radius: 5px;
    box-shadow: 0 2px 6px var(--accent-blue-glow);
    display: flex;
    align-items: center;
    gap: 4px;
  }

  .empty-state {
    padding: 48px 20px;
    text-align: center;
    color: var(--text-muted);
    font-size: 13.5px;
  }

  .empty-hint {
    display: block;
    margin-top: 8px;
    font-size: 11.5px;
    color: var(--text-muted);
  }
</style>
