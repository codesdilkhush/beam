<script lang="ts">
  import type { TabType } from '../types/index.js';

  interface Props {
    searchQuery: string;
    isClipboardMode: boolean;
    activeTab: TabType;
    onQueryChange: (q: string) => void;
    inputElement?: HTMLInputElement | null;
  }

  let {
    searchQuery,
    isClipboardMode,
    activeTab,
    onQueryChange,
    inputElement = $bindable(null)
  }: Props = $props();
</script>

<header class="header">
  <div class="search-icon">
    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2">
      <circle cx="11" cy="11" r="8"></circle>
      <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
    </svg>
  </div>

  <input
    bind:this={inputElement}
    value={searchQuery}
    oninput={(e) => onQueryChange((e.target as HTMLInputElement).value)}
    type="text"
    class="search-input"
    placeholder={isClipboardMode ? "Search clipboard history..." : "Search apps, calculate, or type /clip..."}
    spellcheck="false"
    autocomplete="off"
  />

  {#if isClipboardMode}
    <span class="mode-badge badge-clipboard">
      <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2">
        <path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1 2-2h2"></path>
        <rect x="8" y="2" width="8" height="4" rx="1" ry="1"></rect>
      </svg>
      Clipboard
    </span>
  {:else if activeTab === 'apps'}
    <span class="mode-badge badge-apps">
      <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2">
        <circle cx="12" cy="12" r="10"></circle>
        <polygon points="10 8 16 12 10 16 10 8"></polygon>
      </svg>
      Applications
    </span>
  {/if}
</header>

<style>
  .header {
    height: 58px;
    border-bottom: 1px solid var(--border-subtle);
    display: flex;
    align-items: center;
    padding: 0 18px;
    gap: 12px;
    background: var(--surface-header);
    backdrop-filter: blur(16px);
    -webkit-backdrop-filter: blur(16px);
  }

  .search-icon {
    color: var(--accent-sky);
    display: flex;
    align-items: center;
    opacity: 0.9;
  }

  .search-input {
    flex: 1;
    background: transparent;
    border: none;
    outline: none;
    font-size: 18px;
    font-weight: 450;
    color: var(--text-primary);
    letter-spacing: -0.2px;
  }

  .search-input::placeholder {
    color: var(--text-muted);
    font-weight: 400;
  }

  .mode-badge {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 11.5px;
    font-weight: 500;
    padding: 3px 10px;
    border-radius: 20px;
    border: 1px solid transparent;
  }

  .badge-clipboard {
    color: var(--accent-amber);
    background: rgba(245, 158, 11, 0.12);
    border-color: rgba(245, 158, 11, 0.25);
  }

  .badge-apps {
    color: var(--accent-sky);
    background: rgba(96, 165, 250, 0.12);
    border-color: rgba(96, 165, 250, 0.25);
  }
</style>
