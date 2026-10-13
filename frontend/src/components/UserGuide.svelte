<script lang="ts">
  import type { GuideTopic } from '../types/index.js';
  import { CopyToClipboard } from '../../wailsjs/go/main/App.js';

  interface Props {
    topics: GuideTopic[];
    selectedIndex: number;
    searchQuery: string;
    listElement?: HTMLDivElement | null;
    onSelect: (index: number) => void;
    onHover: (index: number) => void;
  }

  let {
    topics,
    selectedIndex,
    searchQuery,
    listElement = $bindable(null),
    onSelect,
    onHover
  }: Props = $props();

  let copiedExample = $state<string | null>(null);

  const selectedTopic = $derived.by(() => {
    if (topics.length > 0 && selectedIndex >= 0 && selectedIndex < topics.length) {
      return topics[selectedIndex];
    }
    return topics[0] || null;
  });

  async function copyText(text: string) {
    await CopyToClipboard(text);
    copiedExample = text;
    setTimeout(() => {
      if (copiedExample === text) {
        copiedExample = null;
      }
    }, 1800);
  }
</script>

<div class="guide-container">
  <!-- Left: Topic list -->
  <div bind:this={listElement} class="guide-sidebar">
    {#if topics.length === 0}
      <div class="empty-state">
        <p>No guide topics found for "{searchQuery}"</p>
        <span class="empty-hint">Try searching 'apps', 'todo', 'emoji', or 'convert'</span>
      </div>
    {:else}
      <div class="topics-list">
        {#each topics as topic, idx}
          {@const isSelected = idx === selectedIndex}
          <button
            type="button"
            class="topic-item {isSelected ? 'selected' : ''}"
            onclick={() => onSelect(idx)}
            onmouseenter={() => onHover(idx)}
          >
            <div class="topic-item-header">
              <span class="topic-category">{topic.category}</span>
            </div>
            <span class="topic-title truncate">{topic.title}</span>
            <span class="topic-summary truncate">{topic.summary}</span>
          </button>
        {/each}
      </div>
    {/if}
  </div>

  <!-- Right: Selected topic content -->
  <div class="guide-content">
    {#if selectedTopic}
      <div class="content-header">
        <div class="header-badge-row">
          <span class="category-badge">{selectedTopic.category}</span>
          <span class="topic-id-tag">topic: {selectedTopic.id}</span>
        </div>
        <h2 class="guide-heading">{selectedTopic.title}</h2>
        <p class="guide-lead">{selectedTopic.summary}</p>
      </div>

      <div class="content-body">
        <div class="body-paragraphs">
          {#each selectedTopic.content as paragraph}
            <p class="guide-p">{paragraph}</p>
          {/each}
        </div>

        {#if selectedTopic.shortcuts && selectedTopic.shortcuts.length > 0}
          <div class="guide-section">
            <h3 class="section-title">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
                <rect x="2" y="4" width="20" height="16" rx="2"></rect>
                <line x1="6" y1="8" x2="6" y2="8"></line>
                <line x1="10" y1="8" x2="10" y2="8"></line>
                <line x1="14" y1="8" x2="14" y2="8"></line>
                <line x1="18" y1="8" x2="18" y2="8"></line>
                <line x1="6" y1="12" x2="6" y2="12"></line>
                <line x1="18" y1="12" x2="18" y2="12"></line>
                <line x1="10" y1="16" x2="14" y2="16"></line>
              </svg>
              Keyboard Shortcuts
            </h3>
            <div class="shortcuts-table">
              {#each selectedTopic.shortcuts as sc}
                <div class="shortcut-row">
                  <span class="shortcut-key"><kbd>{sc.key}</kbd></span>
                  <span class="shortcut-desc">{sc.desc}</span>
                </div>
              {/each}
            </div>
          </div>
        {/if}

        {#if selectedTopic.examples && selectedTopic.examples.length > 0}
          <div class="guide-section">
            <h3 class="section-title">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="9 11 12 14 22 4"></polyline>
                <path d="M21 12v7a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11"></path>
              </svg>
              Try Examples (Click to copy)
            </h3>
            <div class="examples-list">
              {#each selectedTopic.examples as example}
                <button
                  type="button"
                  class="example-pill {copiedExample === example ? 'copied' : ''}"
                  onclick={() => copyText(example)}
                >
                  <code class="example-code">{example}</code>
                  <span class="copy-status">
                    {#if copiedExample === example}
                      Copied!
                    {:else}
                      Copy
                    {/if}
                  </span>
                </button>
              {/each}
            </div>
          </div>
        {/if}
      </div>
    {:else}
      <div class="empty-state">
        <p>Select a topic to read</p>
      </div>
    {/if}
  </div>
</div>

<style>
  .guide-container {
    flex: 1;
    display: flex;
    overflow: hidden;
    height: 100%;
  }

  .guide-sidebar {
    flex: 5;
    overflow-y: auto;
    padding: 10px;
    border-right: 1px solid var(--border-subtle);
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .topics-list {
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .topic-item {
    width: 100%;
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    padding: 9px 12px;
    border-radius: 9px;
    background: transparent;
    border: 1px solid transparent;
    cursor: pointer;
    text-align: left;
    transition: background 0.1s ease, border-color 0.1s ease;
  }

  .topic-item:hover {
    background: rgba(255, 255, 255, 0.04);
  }

  .topic-item.selected {
    background: var(--surface-card);
    border-color: var(--border-glass);
    box-shadow: 0 4px 14px rgba(0, 0, 0, 0.2);
  }

  .topic-item-header {
    margin-bottom: 2px;
  }

  .topic-category {
    font-size: 10px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    color: var(--accent-sky);
  }

  .topic-title {
    font-size: 13px;
    font-weight: 600;
    color: var(--text-primary);
    width: 100%;
    margin-bottom: 2px;
  }

  .topic-summary {
    font-size: 11px;
    color: var(--text-muted);
    width: 100%;
    line-height: 1.35;
  }

  .guide-content {
    flex: 7;
    overflow-y: auto;
    padding: 20px 22px;
    background: var(--surface-card);
    display: flex;
    flex-direction: column;
    gap: 18px;
  }

  .content-header {
    display: flex;
    flex-direction: column;
    gap: 6px;
    padding-bottom: 14px;
    border-bottom: 1px solid var(--border-subtle);
  }

  .header-badge-row {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .category-badge {
    font-size: 10.5px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    color: var(--accent-sky);
    background: rgba(56, 117, 246, 0.14);
    padding: 2px 7px;
    border-radius: 4px;
    border: 1px solid rgba(56, 117, 246, 0.25);
  }

  .topic-id-tag {
    font-family: ui-monospace, SFMono-Regular, monospace;
    font-size: 10.5px;
    color: var(--text-muted);
  }

  .guide-heading {
    font-size: 18px;
    font-weight: 700;
    color: var(--text-primary);
    margin: 2px 0 0 0;
    letter-spacing: -0.2px;
  }

  .guide-lead {
    font-size: 12.5px;
    color: var(--text-secondary);
    line-height: 1.45;
    margin: 0;
  }

  .content-body {
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .body-paragraphs {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .guide-p {
    font-size: 12.5px;
    color: var(--text-secondary);
    line-height: 1.55;
    margin: 0;
  }

  .guide-section {
    display: flex;
    flex-direction: column;
    gap: 10px;
    padding-top: 10px;
    border-top: 1px solid var(--border-subtle);
  }

  .section-title {
    font-size: 12px;
    font-weight: 600;
    color: var(--text-primary);
    display: flex;
    align-items: center;
    gap: 6px;
    margin: 0;
    letter-spacing: 0.2px;
  }

  .section-title svg {
    color: var(--accent-sky);
  }

  .shortcuts-table {
    display: flex;
    flex-direction: column;
    gap: 6px;
    background: rgba(0, 0, 0, 0.2);
    padding: 10px 12px;
    border-radius: 8px;
    border: 1px solid var(--border-subtle);
  }

  .shortcut-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    font-size: 12px;
    gap: 12px;
  }

  .shortcut-key {
    flex-shrink: 0;
  }

  .shortcut-desc {
    color: var(--text-secondary);
    text-align: right;
  }

  .examples-list {
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .example-pill {
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 6px 10px;
    background: rgba(255, 255, 255, 0.04);
    border: 1px solid var(--border-subtle);
    border-radius: 6px;
    cursor: pointer;
    transition: background 0.12s ease, border-color 0.12s ease;
    text-align: left;
  }

  .example-pill:hover {
    background: rgba(255, 255, 255, 0.08);
    border-color: rgba(255, 255, 255, 0.16);
  }

  .example-pill.copied {
    border-color: #10b981;
    background: rgba(16, 185, 129, 0.1);
  }

  .example-code {
    font-family: ui-monospace, SFMono-Regular, monospace;
    font-size: 11.5px;
    color: var(--accent-sky);
  }

  .copy-status {
    font-size: 11px;
    color: var(--text-muted);
    font-weight: 500;
  }

  .example-pill.copied .copy-status {
    color: #10b981;
    font-weight: 600;
  }

  kbd {
    background: rgba(255, 255, 255, 0.08);
    border: 1px solid rgba(255, 255, 255, 0.14);
    border-radius: 4px;
    padding: 2px 6px;
    font-size: 11px;
    color: var(--text-primary);
  }

  .empty-state {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    text-align: center;
    padding: 30px;
    color: var(--text-muted);
    gap: 6px;
  }

  .empty-state p {
    font-size: 13.5px;
    color: var(--text-secondary);
    margin: 0;
  }

  .empty-hint {
    font-size: 11.5px;
    opacity: 0.8;
  }

  .truncate {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
</style>
