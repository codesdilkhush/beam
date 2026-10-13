<script lang="ts">
  import type { AppItem, ClipboardEntry } from '../types/index.js';
  import { formatTime } from '../utils/format.js';

  interface Props {
    isClipboardMode: boolean;
    selectedClip: ClipboardEntry | null;
    selectedApp: AppItem | null;
  }

  let { isClipboardMode, selectedClip, selectedApp }: Props = $props();

  function countWords(str: string): number {
    if (!str) return 0;
    const matches = str.trim().match(/\S+/g);
    return matches ? matches.length : 0;
  }
</script>

<div class="detail-panel">
  {#if isClipboardMode && selectedClip}
    <!-- Clipboard Detail View -->
    <div class="panel-section">
      <div class="clip-header">
        <div class="clip-title-row">
          <span class="clip-badge">Clipboard Item</span>
          {#if selectedClip.timestamp}
            <span class="time-tag">{formatTime(selectedClip.timestamp)}</span>
          {/if}
        </div>
        <div class="clip-stats">
          <span>{selectedClip.charCount} characters</span>
          <span>•</span>
          <span>{countWords(selectedClip.content)} words</span>
        </div>
      </div>

      <div class="clip-preview-container">
        <pre class="clip-text">{selectedClip.content}</pre>
      </div>
    </div>
  {:else if selectedApp}
    <!-- End-User Focused App View -->
    <div class="app-hero">
      <div class="icon-wrapper">
        {#if selectedApp.resolvedIcon}
          <img src={selectedApp.resolvedIcon} alt="" class="hero-icon" />
        {:else}
          <div class="hero-letter-icon">{selectedApp.name.charAt(0).toUpperCase()}</div>
        {/if}
      </div>

      <div class="hero-info">
        <h2 class="app-heading truncate">{selectedApp.name}</h2>
        {#if selectedApp.launchCount > 0}
          <span class="activity-badge">
            Frequently Used • {selectedApp.launchCount} {selectedApp.launchCount === 1 ? 'launch' : 'launches'}
          </span>
        {/if}
      </div>
    </div>

    <!-- What it does -->
    <div class="info-block">
      <span class="block-label">About</span>
      <p class="app-about">
        {selectedApp.comment || 'Desktop application for Linux.'}
      </p>
    </div>
  {:else}
    <!-- Empty / Idle State -->
    <div class="placeholder">
      <div class="placeholder-icon">
        <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.8">
          <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"></polygon>
        </svg>
      </div>
      <p class="placeholder-title">Beam Spotlight</p>
      <span class="placeholder-subtitle">Type to search apps, calculate math, or manage clipboard</span>
    </div>
  {/if}
</div>

<style>
  .detail-panel {
    flex: 4;
    padding: 22px 20px;
    overflow-y: auto;
    background: var(--surface-card);
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  /* App Hero */
  .app-hero {
    display: flex;
    align-items: center;
    gap: 14px;
    padding-bottom: 16px;
    border-bottom: 1px solid var(--border-subtle);
  }

  .icon-wrapper {
    flex-shrink: 0;
  }

  .hero-icon {
    width: 46px;
    height: 46px;
    object-fit: contain;
    border-radius: 10px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  }

  .hero-letter-icon {
    width: 46px;
    height: 46px;
    border-radius: 10px;
    background: linear-gradient(135deg, rgba(255, 255, 255, 0.14), rgba(255, 255, 255, 0.04));
    border: 1px solid rgba(255, 255, 255, 0.12);
    color: var(--accent-sky);
    font-weight: 700;
    font-size: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.25);
  }

  .hero-info {
    min-width: 0;
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .app-heading {
    font-size: 16px;
    font-weight: 600;
    color: var(--text-primary);
    letter-spacing: -0.2px;
  }

  .activity-badge {
    font-size: 11px;
    color: var(--accent-sky);
    font-weight: 500;
  }

  /* Info Blocks */
  .info-block {
    display: flex;
    flex-direction: column;
    gap: 6px;
  }

  .block-label {
    font-size: 10.5px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.6px;
    color: var(--text-muted);
  }

  .app-about {
    font-size: 13px;
    line-height: 1.55;
    color: var(--text-secondary);
  }

  /* Clipboard Detail */
  .panel-section {
    height: 100%;
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .clip-header {
    display: flex;
    flex-direction: column;
    gap: 4px;
    padding-bottom: 10px;
    border-bottom: 1px solid var(--border-subtle);
  }

  .clip-title-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .clip-badge {
    font-size: 11px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.6px;
    color: var(--accent-amber);
  }

  .time-tag {
    font-size: 11px;
    color: var(--text-muted);
  }

  .clip-stats {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 11.5px;
    color: var(--text-muted);
  }

  .clip-preview-container {
    flex: 1;
    overflow-y: auto;
    background: rgba(0, 0, 0, 0.35);
    border: 1px solid var(--border-subtle);
    border-radius: 8px;
    padding: 12px;
  }

  .clip-text {
    font-family: ui-monospace, SFMono-Regular, "JetBrains Mono", Menlo, monospace;
    font-size: 12.5px;
    line-height: 1.55;
    color: var(--text-primary);
    white-space: pre-wrap;
    word-break: break-all;
    user-select: text;
  }

  /* Placeholder */
  .placeholder {
    height: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    text-align: center;
    gap: 10px;
    color: var(--text-muted);
  }

  .placeholder-icon {
    color: var(--accent-sky);
    opacity: 0.7;
  }

  .placeholder-title {
    font-size: 15px;
    font-weight: 600;
    color: var(--text-primary);
  }

  .placeholder-subtitle {
    font-size: 12px;
    max-width: 200px;
    line-height: 1.4;
  }

  .truncate {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
</style>
