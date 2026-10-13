<script lang="ts">
  import { onMount } from 'svelte';
  import {
    GetApplications,
    LaunchApp,
    GetClipboard,
    CopyToClipboard,
    EvaluateMath,
    OpenURL,
    HideWindow,
    QuitApp,
    LoadIconAsDataURL
  } from '../wailsjs/go/main/App.js';

  interface AppItem {
    id: string;
    name: string;
    comment: string;
    exec: string;
    icon: string;
    iconPath: string;
    terminal: boolean;
    desktopPath: string;
    resolvedIcon?: string;
  }

  interface ClipboardEntry {
    id: string;
    content: string;
    preview: string;
    timestamp: string;
    charCount: number;
  }

  interface CalcItem {
    expression: string;
    result: string;
    success: boolean;
  }

  let searchQuery = $state('');
  let activeTab = $state<'all' | 'apps' | 'clipboard'>('all');
  let selectedIndex = $state(0);

  let allApps = $state<AppItem[]>([]);
  let clipboardItems = $state<ClipboardEntry[]>([]);
  let mathResult = $state<CalcItem | null>(null);

  let searchInputRef = $state<HTMLInputElement | null>(null);

  onMount(async () => {
    await refreshApps();
    await refreshClipboard();
    if (searchInputRef) {
      searchInputRef.focus();
    }
  });

  async function refreshApps() {
    try {
      const apps = await GetApplications();
      allApps = apps || [];
      // Lazy load some icons
      for (const app of allApps.slice(0, 30)) {
        if (app.iconPath) {
          LoadIconAsDataURL(app.iconPath).then((dataUrl: string) => {
            if (dataUrl) app.resolvedIcon = dataUrl;
          });
        }
      }
    } catch (e) {
      console.error('Failed to load apps:', e);
    }
  }

  async function refreshClipboard() {
    try {
      const clips = await GetClipboard();
      clipboardItems = clips || [];
    } catch (e) {
      console.error('Failed to load clipboard:', e);
    }
  }

  // Handle live calculation
  $effect(() => {
    const q = searchQuery.trim();
    if (q.length > 1 && /[0-9\+\-\*\/\%\^\(\)]/.test(q)) {
      EvaluateMath(q).then((res) => {
        if (res && res.success) {
          mathResult = res;
        } else {
          mathResult = null;
        }
      }).catch(() => {
        mathResult = null;
      });
    } else {
      mathResult = null;
    }
  });

  // Filter items
  let filteredApps = $derived.by(() => {
    const q = searchQuery.toLowerCase().trim();
    if (!q || q === '/clip' || q.startsWith('/clip ')) return allApps;
    return allApps.filter(
      (a) =>
        a.name.toLowerCase().includes(q) ||
        (a.comment && a.comment.toLowerCase().includes(q)) ||
        (a.exec && a.exec.toLowerCase().includes(q))
    );
  });

  let isClipboardMode = $derived.by(() => {
    return activeTab === 'clipboard' || searchQuery.startsWith('/clip');
  });

  let displayItems = $derived.by(() => {
    if (isClipboardMode) {
      const clipQuery = searchQuery.replace(/^\/clip\s*/, '').toLowerCase().trim();
      if (!clipQuery) return clipboardItems;
      return clipboardItems.filter((c) => c.content.toLowerCase().includes(clipQuery));
    }
    return filteredApps;
  });

  let selectedApp = $derived(
    !isClipboardMode && displayItems.length > 0
      ? (displayItems[selectedIndex] as AppItem)
      : null
  );

  let selectedClip = $derived(
    isClipboardMode && displayItems.length > 0
      ? (displayItems[selectedIndex] as ClipboardEntry)
      : null
  );

  function handleKeyDown(e: KeyboardEvent) {
    if (e.key === 'ArrowDown') {
      e.preventDefault();
      if (selectedIndex < displayItems.length - 1) {
        selectedIndex++;
      }
    } else if (e.key === 'ArrowUp') {
      e.preventDefault();
      if (selectedIndex > 0) {
        selectedIndex--;
      }
    } else if (e.key === 'Enter') {
      e.preventDefault();
      executeCurrentSelection();
    } else if (e.key === 'Escape') {
      e.preventDefault();
      if (searchQuery) {
        searchQuery = '';
        selectedIndex = 0;
      } else {
        HideWindow();
      }
    } else if (e.key === 'Tab') {
      e.preventDefault();
      if (activeTab === 'all') activeTab = 'clipboard';
      else if (activeTab === 'clipboard') activeTab = 'apps';
      else activeTab = 'all';
      selectedIndex = 0;
      if (activeTab === 'clipboard') refreshClipboard();
    }
  }

  async function executeCurrentSelection() {
    if (mathResult && mathResult.success && selectedIndex === 0 && !isClipboardMode) {
      await CopyToClipboard(mathResult.result);
      return;
    }

    if (isClipboardMode) {
      if (selectedClip) {
        await CopyToClipboard(selectedClip.content);
      }
    } else {
      if (selectedApp) {
        await LaunchApp(selectedApp.exec, selectedApp.terminal);
      } else if (searchQuery.trim()) {
        // Fallback to Google search
        await OpenURL(`https://www.google.com/search?q=${encodeURI(searchQuery.trim())}`);
      }
    }
  }

  async function selectItem(idx: number) {
    selectedIndex = idx;
    executeCurrentSelection();
  }
</script>

<svelte:window onkeydown={handleKeyDown} />

<div class="beam-window">
  <!-- Top Bar: Search Input & Mode Chips -->
  <header class="beam-header">
    <div class="search-icon">
      <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <circle cx="11" cy="11" r="8"></circle>
        <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
      </svg>
    </div>

    <input
      bind:this={searchInputRef}
      bind:value={searchQuery}
      type="text"
      class="search-input"
      placeholder={isClipboardMode ? "Search clipboard history..." : "Search apps, calculate, or type /clip..."}
      spellcheck="false"
      autocomplete="off"
    />

    {#if isClipboardMode}
      <span class="mode-badge clip-badge">📋 Clipboard</span>
    {:else if activeTab === 'apps'}
      <span class="mode-badge app-badge">🚀 Apps</span>
    {/if}
  </header>

  <!-- Math Inline Banner if active -->
  {#if mathResult && mathResult.success && !isClipboardMode}
    <div class="math-banner" onclick={() => CopyToClipboard(mathResult?.result || '')}>
      <div class="math-left">
        <span class="math-icon">🧮</span>
        <span class="math-expr">{mathResult.expression} =</span>
        <span class="math-res">{mathResult.result}</span>
      </div>
      <div class="math-right">
        <span class="key-hint">↵ Copy Result</span>
      </div>
    </div>
  {/if}

  <!-- Split Main Body: Results List + Detail Pane -->
  <div class="beam-body">
    <!-- Results List (Left) -->
    <div class="results-list">
      {#if displayItems.length === 0}
        <div class="empty-state">
          <p>No results found for "{searchQuery}"</p>
          <span class="empty-hint">Press ↵ to search on Google</span>
        </div>
      {:else}
        {#each displayItems as item, idx}
          {@const isSelected = idx === selectedIndex}
          <div
            class="list-item {isSelected ? 'selected' : ''}"
            onclick={() => selectItem(idx)}
            onmouseenter={() => (selectedIndex = idx)}
          >
            {#if isClipboardMode}
              <!-- Clipboard Row -->
              {@const clip = item as ClipboardEntry}
              <div class="item-icon-box clip-icon">
                <span>📋</span>
              </div>
              <div class="item-info">
                <span class="item-title truncate">{clip.preview}</span>
                <span class="item-subtitle">{clip.charCount} characters</span>
              </div>
              {#if isSelected}
                <span class="item-action-pill">↵ Copy</span>
              {/if}
            {:else}
              <!-- App Row -->
              {@const app = item as AppItem}
              <div class="item-icon-box">
                {#if app.resolvedIcon}
                  <img src={app.resolvedIcon} alt="" class="app-icon" />
                {:else}
                  <div class="generic-icon">{app.name.charAt(0)}</div>
                {/if}
              </div>
              <div class="item-info">
                <span class="item-title">{app.name}</span>
                <span class="item-subtitle truncate">{app.comment || app.exec}</span>
              </div>
              {#if isSelected}
                <span class="item-action-pill">↵ Open</span>
              {/if}
            {/if}
          </div>
        {/each}
      {/if}
    </div>

    <!-- Right Detail Panel -->
    <div class="detail-panel">
      {#if isClipboardMode && selectedClip}
        <div class="detail-header">
          <h3>Clipboard Content</h3>
          <span class="detail-badge">{selectedClip.charCount} chars</span>
        </div>
        <div class="clip-preview-box">
          <pre>{selectedClip.content}</pre>
        </div>
      {:else if selectedApp}
        <div class="detail-header">
          <div class="detail-app-title">
            {#if selectedApp.resolvedIcon}
              <img src={selectedApp.resolvedIcon} alt="" class="detail-app-icon" />
            {:else}
              <div class="generic-icon large">{selectedApp.name.charAt(0)}</div>
            {/if}
            <div>
              <h3>{selectedApp.name}</h3>
              <p class="detail-app-subtitle">Application</p>
            </div>
          </div>
        </div>

        <div class="detail-specs">
          <div class="spec-row">
            <span class="spec-label">Command</span>
            <span class="spec-val truncate">{selectedApp.exec}</span>
          </div>
          {#if selectedApp.comment}
            <div class="spec-row">
              <span class="spec-label">About</span>
              <span class="spec-val">{selectedApp.comment}</span>
            </div>
          {/if}
          <div class="spec-row">
            <span class="spec-label">Desktop File</span>
            <span class="spec-val truncate">{selectedApp.id}</span>
          </div>
        </div>
      {:else}
        <div class="detail-placeholder">
          <p>Beam Launcher</p>
          <span>Fast, keyboard-driven workflow for Linux</span>
        </div>
      {/if}
    </div>
  </div>

  <!-- Bottom Footer -->
  <footer class="beam-footer">
    <div class="footer-left">
      {#if isClipboardMode}
        <span class="key-pill">↵ Copy to Clipboard</span>
      {:else if selectedApp}
        <span class="key-pill">↵ Open Application</span>
      {:else}
        <span class="key-pill">↵ Search Online</span>
      {/if}
    </div>
    <div class="footer-right">
      <span class="footer-hint">Tab: Mode</span>
      <span class="footer-hint">↑↓: Navigate</span>
      <span class="footer-hint">Esc: Close</span>
    </div>
  </footer>
</div>

<style>
  .beam-window {
    width: 760px;
    height: 500px;
    background: #14151b;
    border: 1px solid #282a36;
    border-radius: 14px;
    box-shadow: 0 20px 45px rgba(0, 0, 0, 0.7);
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  /* Header */
  .beam-header {
    height: 60px;
    border-bottom: 1px solid #232530;
    display: flex;
    align-items: center;
    padding: 0 18px;
    gap: 12px;
    background: #171821;
  }

  .search-icon {
    color: #6b7280;
    display: flex;
    align-items: center;
  }

  .search-input {
    flex: 1;
    background: transparent;
    border: none;
    outline: none;
    color: #f3f4f6;
    font-size: 19px;
    font-weight: 500;
  }

  .search-input::placeholder {
    color: #4b5563;
  }

  .mode-badge {
    font-size: 12px;
    padding: 4px 10px;
    border-radius: 20px;
    font-weight: 600;
  }

  .clip-badge {
    background: rgba(59, 130, 246, 0.15);
    color: #60a5fa;
    border: 1px solid rgba(59, 130, 246, 0.3);
  }

  .app-badge {
    background: rgba(168, 85, 247, 0.15);
    color: #c084fc;
    border: 1px solid rgba(168, 85, 247, 0.3);
  }

  /* Math Banner */
  .math-banner {
    background: #1e202c;
    border-bottom: 1px solid #2d3042;
    padding: 10px 20px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    cursor: pointer;
  }

  .math-left {
    display: flex;
    align-items: center;
    gap: 10px;
  }

  .math-expr {
    color: #9ca3af;
    font-size: 15px;
  }

  .math-res {
    color: #10b981;
    font-size: 18px;
    font-weight: 700;
  }

  .key-hint {
    background: #2a2d3e;
    color: #9ca3af;
    font-size: 12px;
    padding: 3px 8px;
    border-radius: 4px;
  }

  /* Body */
  .beam-body {
    flex: 1;
    display: flex;
    overflow: hidden;
  }

  .results-list {
    flex: 6;
    overflow-y: auto;
    padding: 8px;
    border-right: 1px solid #232530;
  }

  .list-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 9px 12px;
    border-radius: 8px;
    cursor: pointer;
    transition: background 0.12s ease;
  }

  .list-item:hover {
    background: #1e202a;
  }

  .list-item.selected {
    background: #292b3a;
  }

  .item-icon-box {
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 6px;
    flex-shrink: 0;
  }

  .app-icon {
    width: 30px;
    height: 30px;
    object-fit: contain;
  }

  .generic-icon {
    width: 30px;
    height: 30px;
    background: #374151;
    border-radius: 6px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 700;
    color: #d1d5db;
    font-size: 14px;
  }

  .generic-icon.large {
    width: 44px;
    height: 44px;
    font-size: 20px;
  }

  .item-info {
    flex: 1;
    display: flex;
    flex-direction: column;
    min-width: 0;
  }

  .item-title {
    font-size: 14.5px;
    font-weight: 500;
    color: #f3f4f6;
  }

  .item-subtitle {
    font-size: 12px;
    color: #6b7280;
  }

  .truncate {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .item-action-pill {
    font-size: 11px;
    background: #3b82f6;
    color: #fff;
    padding: 3px 8px;
    border-radius: 4px;
    font-weight: 600;
  }

  .empty-state {
    padding: 40px 20px;
    text-align: center;
    color: #6b7280;
  }

  .empty-hint {
    display: block;
    margin-top: 8px;
    font-size: 12px;
    color: #4b5563;
  }

  /* Detail Panel */
  .detail-panel {
    flex: 4;
    padding: 16px;
    overflow-y: auto;
    background: #16171f;
    display: flex;
    flex-direction: column;
    gap: 16px;
  }

  .detail-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .detail-app-title {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .detail-app-icon {
    width: 44px;
    height: 44px;
    object-fit: contain;
  }

  .detail-app-subtitle {
    font-size: 12px;
    color: #6b7280;
  }

  .detail-badge {
    font-size: 11px;
    background: #2a2c3a;
    color: #9ca3af;
    padding: 3px 8px;
    border-radius: 12px;
  }

  .clip-preview-box {
    background: #1c1d27;
    border: 1px solid #282a39;
    padding: 12px;
    border-radius: 8px;
    max-height: 280px;
    overflow-y: auto;
  }

  .clip-preview-box pre {
    white-space: pre-wrap;
    word-break: break-all;
    font-family: monospace;
    font-size: 12.5px;
    color: #d1d5db;
  }

  .detail-specs {
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  .spec-row {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  .spec-label {
    font-size: 11px;
    text-transform: uppercase;
    letter-spacing: 0.5px;
    color: #4b5563;
    font-weight: 600;
  }

  .spec-val {
    font-size: 13px;
    color: #d1d5db;
  }

  .detail-placeholder {
    height: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    color: #4b5563;
    text-align: center;
    gap: 6px;
  }

  .detail-placeholder p {
    font-size: 15px;
    color: #6b7280;
    font-weight: 600;
  }

  .detail-placeholder span {
    font-size: 12px;
  }

  /* Footer */
  .beam-footer {
    height: 38px;
    border-top: 1px solid #232530;
    background: #111218;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 16px;
  }

  .key-pill {
    font-size: 12px;
    font-weight: 600;
    color: #d1d5db;
  }

  .footer-right {
    display: flex;
    gap: 12px;
  }

  .footer-hint {
    font-size: 11px;
    color: #6b7280;
  }
</style>
