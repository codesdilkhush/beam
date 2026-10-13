<script lang="ts">
  import { onMount } from 'svelte';
  import {
    SearchApps,
    LaunchApp,
    GetClipboard,
    CopyToClipboard,
    EvaluateMath,
    OpenURL,
    HideWindow,
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
    score: number;
    launchCount: number;
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

  interface ActionOption {
    id: string;
    label: string;
    shortcut: string;
    run: () => Promise<void> | void;
  }

  let searchQuery = $state('');
  let activeTab = $state<'all' | 'apps' | 'clipboard'>('all');
  let selectedIndex = $state(0);
  let showActionPalette = $state(false);
  let actionPaletteIndex = $state(0);
  let currentTheme = $state<'default' | 'theme-catppuccin' | 'theme-tokyo' | 'theme-nord'>('default');

  let searchResults = $state<AppItem[]>([]);
  let clipboardItems = $state<ClipboardEntry[]>([]);
  let mathResult = $state<CalcItem | null>(null);

  let searchInputRef = $state<HTMLInputElement | null>(null);

  onMount(async () => {
    const savedTheme = localStorage.getItem('beam_theme');
    if (savedTheme && ['default', 'theme-catppuccin', 'theme-tokyo', 'theme-nord'].includes(savedTheme)) {
      currentTheme = savedTheme as typeof currentTheme;
    }

    await performAppSearch('');
    await refreshClipboard();
    if (searchInputRef) {
      searchInputRef.focus();
    }
  });

  function setTheme(t: typeof currentTheme) {
    currentTheme = t;
    localStorage.setItem('beam_theme', t);
    showActionPalette = false;
  }

  async function performAppSearch(query: string) {
    try {
      const q = query.replace(/^\/clip\s*/, '').trim();
      const apps = await SearchApps(q);
      searchResults = apps || [];

      // Lazy load icons for remaining
      for (const app of searchResults.slice(0, 25)) {
        if (app.iconPath && !app.resolvedIcon) {
          LoadIconAsDataURL(app.iconPath).then((dataUrl: string) => {
            if (dataUrl) app.resolvedIcon = dataUrl;
          });
        }
      }
    } catch (e) {
      console.error('Failed to search apps:', e);
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

  $effect(() => {
    const q = searchQuery;
    if (!isClipboardMode) {
      performAppSearch(q);
    }
    selectedIndex = 0;
  });

  $effect(() => {
    const q = searchQuery.trim();
    if (q.length > 1 && /[0-9\+\-\*\/\%\^\(\)]/.test(q)) {
      EvaluateMath(q)
        .then((res) => {
          if (res && res.success) {
            mathResult = res;
          } else {
            mathResult = null;
          }
        })
        .catch(() => {
          mathResult = null;
        });
    } else {
      mathResult = null;
    }
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
    return searchResults;
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

  let currentActions = $derived.by((): ActionOption[] => {
    const list: ActionOption[] = [];
    if (!isClipboardMode && selectedApp) {
      list.push({
        id: 'launch',
        label: `Open ${selectedApp.name}`,
        shortcut: '↵',
        run: () => LaunchApp(selectedApp!.exec, selectedApp!.terminal, selectedApp!.id)
      });
      list.push({
        id: 'term',
        label: `Run in Terminal`,
        shortcut: 'Ctrl+T',
        run: () => LaunchApp(selectedApp!.exec, true, selectedApp!.id)
      });
      list.push({
        id: 'copy-exec',
        label: 'Copy Command Line',
        shortcut: 'Ctrl+C',
        run: () => CopyToClipboard(selectedApp!.exec)
      });
      list.push({
        id: 'copy-path',
        label: 'Copy Desktop File Path',
        shortcut: 'Ctrl+P',
        run: () => CopyToClipboard(selectedApp!.desktopPath)
      });
    } else if (isClipboardMode && selectedClip) {
      list.push({
        id: 'copy',
        label: 'Copy to Clipboard',
        shortcut: '↵',
        run: () => CopyToClipboard(selectedClip!.content)
      });
    }

    // Universal theme actions
    list.push({
      id: 'theme-default',
      label: 'Switch to Default Dark Theme',
      shortcut: 'T1',
      run: () => setTheme('default')
    });
    list.push({
      id: 'theme-catppuccin',
      label: 'Switch to Catppuccin Mocha Theme',
      shortcut: 'T2',
      run: () => setTheme('theme-catppuccin')
    });
    list.push({
      id: 'theme-tokyo',
      label: 'Switch to Tokyo Night Theme',
      shortcut: 'T3',
      run: () => setTheme('theme-tokyo')
    });
    list.push({
      id: 'theme-nord',
      label: 'Switch to Nord Theme',
      shortcut: 'T4',
      run: () => setTheme('theme-nord')
    });

    return list;
  });

  function handleKeyDown(e: KeyboardEvent) {
    if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 'k') {
      e.preventDefault();
      showActionPalette = !showActionPalette;
      actionPaletteIndex = 0;
      return;
    }

    if (showActionPalette) {
      if (e.key === 'ArrowDown') {
        e.preventDefault();
        if (actionPaletteIndex < currentActions.length - 1) actionPaletteIndex++;
      } else if (e.key === 'ArrowUp') {
        e.preventDefault();
        if (actionPaletteIndex > 0) actionPaletteIndex--;
      } else if (e.key === 'Enter') {
        e.preventDefault();
        if (currentActions[actionPaletteIndex]) {
          currentActions[actionPaletteIndex].run();
          showActionPalette = false;
        }
      } else if (e.key === 'Escape') {
        e.preventDefault();
        showActionPalette = false;
      }
      return;
    }

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
        await LaunchApp(selectedApp.exec, selectedApp.terminal, selectedApp.id);
      } else if (searchQuery.trim()) {
        await OpenURL(`https://www.google.com/search?q=${encodeURI(searchQuery.trim())}`);
      }
    }
  }

  async function selectItem(idx: number) {
    selectedIndex = idx;
    executeCurrentSelection();
  }

  function formatTime(iso: string) {
    if (!iso) return '';
    try {
      const d = new Date(iso);
      const diffSec = Math.floor((Date.now() - d.getTime()) / 1000);
      if (diffSec < 60) return 'Just now';
      if (diffSec < 3600) return `${Math.floor(diffSec / 60)}m ago`;
      if (diffSec < 86400) return `${Math.floor(diffSec / 3600)}h ago`;
      return `${Math.floor(diffSec / 86400)}d ago`;
    } catch {
      return '';
    }
  }
</script>

<svelte:window onkeydown={handleKeyDown} />

<div class="beam-window {currentTheme}">
  <!-- Top Search Bar -->
  <header class="beam-header">
    <div class="search-icon">
      <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2">
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
      <span class="mode-badge clip-badge">
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2">
          <path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"></path>
          <rect x="8" y="2" width="8" height="4" rx="1" ry="1"></rect>
        </svg>
        Clipboard
      </span>
    {:else if activeTab === 'apps'}
      <span class="mode-badge app-badge">
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2">
          <circle cx="12" cy="12" r="10"></circle>
          <polygon points="10 8 16 12 10 16 10 8"></polygon>
        </svg>
        Applications
      </span>
    {/if}
  </header>

  <!-- Math Inline Banner (Emerald Glow) -->
  {#if mathResult && mathResult.success && !isClipboardMode}
    <button type="button" class="math-banner" onclick={() => CopyToClipboard(mathResult?.result || '')}>
      <div class="math-left">
        <span class="math-icon">
          <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.2">
            <rect x="4" y="2" width="16" height="20" rx="2"></rect>
            <line x1="8" y1="6" x2="16" y2="6"></line>
            <line x1="16" y1="14" x2="16" y2="18"></line>
            <path d="M16 10h.01"></path>
            <path d="M12 10h.01"></path>
            <path d="M8 10h.01"></path>
          </svg>
        </span>
        <span class="math-expr">{mathResult.expression} =</span>
        <span class="math-res">{mathResult.result}</span>
      </div>
      <div class="math-right">
        <span class="key-hint"><kbd>↵</kbd> Copy Result</span>
      </div>
    </button>
  {/if}

  <!-- Split View: Results (60%) + Details (40%) -->
  <div class="beam-body">
    <!-- Results List -->
    <div class="results-list">
      {#if displayItems.length === 0}
        <div class="empty-state">
          <p>No results found for "{searchQuery}"</p>
          <span class="empty-hint">Press <kbd>↵</kbd> to search on Google</span>
        </div>
      {:else}
        {#each displayItems as item, idx}
          {@const isSelected = idx === selectedIndex}
          <button
            type="button"
            class="list-item {isSelected ? 'selected' : ''}"
            onclick={() => selectItem(idx)}
            onmouseenter={() => (selectedIndex = idx)}
          >
            {#if isClipboardMode}
              {@const clip = item as ClipboardEntry}
              <div class="item-icon-box clip-icon">
                <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"></path>
                  <rect x="8" y="2" width="8" height="4" rx="1" ry="1"></rect>
                </svg>
              </div>
              <div class="item-info">
                <span class="item-title truncate">{clip.preview}</span>
                <span class="item-subtitle">
                  {clip.charCount} chars {clip.timestamp ? `• ${formatTime(clip.timestamp)}` : ''}
                </span>
              </div>
              {#if isSelected}
                <span class="item-action-pill">↵ Copy</span>
              {/if}
            {:else}
              {@const app = item as AppItem}
              <div class="item-icon-box">
                {#if app.resolvedIcon}
                  <img src={app.resolvedIcon} alt="" class="app-icon" />
                {:else}
                  <div class="generic-icon">{app.name.charAt(0)}</div>
                {/if}
              </div>
              <div class="item-info">
                <div class="title-row">
                  <span class="item-title">{app.name}</span>
                  {#if app.launchCount > 0}
                    <span class="usage-badge" title="Launched {app.launchCount} times">{app.launchCount} runs</span>
                  {/if}
                </div>
                <span class="item-subtitle truncate">{app.comment || app.exec}</span>
              </div>
              {#if isSelected}
                <span class="item-action-pill">↵ Open</span>
              {/if}
            {/if}
          </button>
        {/each}
      {/if}
    </div>

    <!-- Detail Panel (40%) -->
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
            <span class="spec-val truncate code-font">{selectedApp.exec}</span>
          </div>
          {#if selectedApp.comment}
            <div class="spec-row">
              <span class="spec-label">Description</span>
              <span class="spec-val">{selectedApp.comment}</span>
            </div>
          {/if}
          <div class="spec-row">
            <span class="spec-label">Desktop File</span>
            <span class="spec-val truncate code-font">{selectedApp.id}</span>
          </div>
          {#if selectedApp.launchCount > 0}
            <div class="spec-row">
              <span class="spec-label">Frequency</span>
              <span class="spec-val">Launched {selectedApp.launchCount} times</span>
            </div>
          {/if}
        </div>
      {:else}
        <div class="detail-placeholder">
          <div class="placeholder-icon">
            <svg width="30" height="30" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2"></polygon>
            </svg>
          </div>
          <p>Beam Launcher</p>
          <span>Raycast-grade speed on Linux</span>
        </div>
      {/if}
    </div>
  </div>

  <!-- Action Palette Modal (Ctrl + K) -->
  {#if showActionPalette}
    <div
      class="action-palette-backdrop"
      role="button"
      tabindex="0"
      aria-label="Close actions"
      onclick={() => (showActionPalette = false)}
      onkeydown={(e) => { if (e.key === 'Escape') showActionPalette = false; }}
    >
      <div
        class="action-palette-modal"
        role="dialog"
        aria-label="Actions Menu"
        tabindex="-1"
        onclick={(e) => e.stopPropagation()}
        onkeydown={(e) => e.stopPropagation()}
      >
        <div class="action-palette-header">
          <span>Actions</span>
          <span class="palette-sub">Press Esc to dismiss</span>
        </div>
        <div class="action-palette-list">
          {#each currentActions as action, aIdx}
            {@const isActSelected = aIdx === actionPaletteIndex}
            <button
              type="button"
              class="action-palette-item {isActSelected ? 'selected' : ''}"
              onclick={() => { action.run(); showActionPalette = false; }}
              onmouseenter={() => (actionPaletteIndex = aIdx)}
            >
              <span class="action-label">{action.label}</span>
              <span class="action-shortcut"><kbd>{action.shortcut}</kbd></span>
            </button>
          {/each}
        </div>
      </div>
    </div>
  {/if}

  <!-- Footer Action Bar -->
  <footer class="beam-footer">
    <div class="footer-left">
      {#if isClipboardMode}
        <span class="key-pill"><kbd>↵</kbd> Copy to Clipboard</span>
      {:else if selectedApp}
        <span class="key-pill"><kbd>↵</kbd> Open Application</span>
      {:else}
        <span class="key-pill"><kbd>↵</kbd> Search Online</span>
      {/if}
    </div>
    <div class="footer-right">
      <span class="footer-hint"><kbd>Ctrl+K</kbd> Actions</span>
      <span class="footer-hint"><kbd>Tab</kbd> Mode</span>
      <span class="footer-hint"><kbd>↑</kbd><kbd>↓</kbd> Navigate</span>
      <span class="footer-hint"><kbd>Esc</kbd> Close</span>
    </div>
  </footer>
</div>

<style>
  .beam-window {
    width: 760px;
    height: 500px;
    background: var(--surface-canvas);
    backdrop-filter: blur(24px) saturate(180%);
    -webkit-backdrop-filter: blur(24px) saturate(180%);
    border: 1px solid var(--border-glass);
    border-radius: 14px;
    box-shadow: 0 24px 60px -12px rgba(0, 0, 0, 0.75), 0 0 0 1px rgba(255, 255, 255, 0.08);
    display: flex;
    flex-direction: column;
    overflow: hidden;
    position: relative;
  }

  /* Header */
  .beam-header {
    height: 60px;
    border-bottom: 1px solid var(--border-subtle);
    display: flex;
    align-items: center;
    padding: 0 18px;
    gap: 12px;
    background: var(--surface-header);
    backdrop-filter: blur(16px);
  }

  .search-icon {
    color: var(--accent-sky);
    display: flex;
    align-items: center;
    opacity: 0.85;
  }

  .search-input {
    flex: 1;
    background: transparent;
    border: none;
    outline: none;
    color: var(--text-primary);
    font-size: 19px;
    font-weight: 500;
  }

  .search-input::placeholder {
    color: var(--text-muted);
  }

  .mode-badge {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    font-size: 12px;
    padding: 4px 10px;
    border-radius: 20px;
    font-weight: 600;
    backdrop-filter: blur(8px);
  }

  .clip-badge {
    background: rgba(56, 117, 246, 0.16);
    color: var(--accent-sky);
    border: 1px solid rgba(56, 117, 246, 0.35);
  }

  .app-badge {
    background: rgba(168, 85, 247, 0.16);
    color: #c084fc;
    border: 1px solid rgba(168, 85, 247, 0.35);
  }

  /* Math Banner (Emerald Glow) */
  .math-banner {
    width: 100%;
    background: rgba(16, 185, 129, 0.08);
    border: none;
    border-bottom: 1px solid rgba(16, 185, 129, 0.22);
    padding: 10px 20px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    cursor: pointer;
    font: inherit;
    text-align: left;
    transition: background 0.15s ease;
  }

  .math-banner:hover {
    background: rgba(16, 185, 129, 0.14);
  }

  .math-left {
    display: flex;
    align-items: center;
    gap: 10px;
  }

  .math-expr {
    color: var(--text-secondary);
    font-size: 15px;
  }

  .math-res {
    color: var(--accent-emerald);
    font-size: 18px;
    font-weight: 700;
  }

  /* Main Body */
  .beam-body {
    flex: 1;
    display: flex;
    overflow: hidden;
  }

  .results-list {
    flex: 6;
    overflow-y: auto;
    padding: 8px;
    border-right: 1px solid var(--border-subtle);
  }

  .list-item {
    width: 100%;
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 9px 12px;
    border-radius: 8px;
    cursor: pointer;
    transition: background 0.1s ease, border-color 0.1s ease;
    border: 1px solid transparent;
    background: transparent;
    font: inherit;
    text-align: left;
  }

  .list-item:hover {
    background: var(--surface-hover);
  }

  .list-item.selected {
    background: var(--surface-selected);
    border-color: var(--border-selected);
    border-left: 3px solid var(--accent-blue);
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
    background: rgba(255, 255, 255, 0.08);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 6px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 700;
    color: var(--text-primary);
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

  .title-row {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .item-title {
    font-size: 14.5px;
    font-weight: 500;
    color: var(--text-primary);
  }

  .usage-badge {
    font-size: 10px;
    background: rgba(255, 255, 255, 0.08);
    color: var(--accent-sky);
    padding: 1px 5px;
    border-radius: 10px;
    font-weight: 600;
  }

  .item-subtitle {
    font-size: 12px;
    color: var(--text-muted);
  }

  .truncate {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .item-action-pill {
    font-size: 11px;
    background: var(--accent-blue);
    color: #ffffff;
    padding: 3px 8px;
    border-radius: 5px;
    font-weight: 600;
    box-shadow: 0 2px 8px var(--accent-blue-glow);
  }

  .empty-state {
    padding: 40px 20px;
    text-align: center;
    color: var(--text-muted);
  }

  .empty-hint {
    display: block;
    margin-top: 8px;
    font-size: 12px;
    color: var(--text-muted);
  }

  /* Detail Panel */
  .detail-panel {
    flex: 4;
    padding: 16px;
    overflow-y: auto;
    background: var(--surface-card);
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
    color: var(--text-muted);
  }

  .detail-badge {
    font-size: 11px;
    background: rgba(255, 255, 255, 0.08);
    color: var(--text-secondary);
    padding: 3px 8px;
    border-radius: 12px;
  }

  .clip-preview-box {
    background: rgba(0, 0, 0, 0.35);
    border: 1px solid var(--surface-card-border);
    padding: 12px;
    border-radius: 8px;
    max-height: 280px;
    overflow-y: auto;
  }

  .clip-preview-box pre {
    white-space: pre-wrap;
    word-break: break-all;
    font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
    font-size: 12.5px;
    color: #e5e7eb;
  }

  .detail-specs {
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .spec-row {
    display: flex;
    flex-direction: column;
    gap: 3px;
  }

  .spec-label {
    font-size: 11px;
    text-transform: uppercase;
    letter-spacing: 0.6px;
    color: var(--text-muted);
    font-weight: 600;
  }

  .spec-val {
    font-size: 13px;
    color: #d1d5db;
  }

  .code-font {
    font-family: ui-monospace, SFMono-Regular, Menlo, monospace;
    font-size: 12px;
  }

  .detail-placeholder {
    height: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    color: var(--text-muted);
    text-align: center;
    gap: 8px;
  }

  .placeholder-icon {
    font-size: 30px;
    color: var(--accent-blue);
    opacity: 0.8;
  }

  .detail-placeholder p {
    font-size: 15px;
    color: var(--text-primary);
    font-weight: 600;
  }

  .detail-placeholder span {
    font-size: 12px;
  }

  /* Action Palette Modal */
  .action-palette-backdrop {
    position: absolute;
    inset: 0;
    background: rgba(0, 0, 0, 0.65);
    backdrop-filter: blur(8px);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 50;
  }

  .action-palette-modal {
    width: 480px;
    background: var(--surface-header);
    border: 1px solid var(--border-glass);
    border-radius: 12px;
    box-shadow: 0 20px 40px rgba(0, 0, 0, 0.8);
    overflow: hidden;
  }

  .action-palette-header {
    padding: 12px 16px;
    border-bottom: 1px solid var(--border-subtle);
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-weight: 600;
    font-size: 14px;
    color: var(--text-primary);
  }

  .palette-sub {
    font-size: 11px;
    color: var(--text-muted);
    font-weight: 400;
  }

  .action-palette-list {
    max-height: 260px;
    overflow-y: auto;
    padding: 6px;
  }

  .action-palette-item {
    width: 100%;
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 9px 12px;
    border-radius: 6px;
    background: transparent;
    border: none;
    color: var(--text-primary);
    font-size: 13.5px;
    cursor: pointer;
    text-align: left;
    transition: background 0.1s ease;
  }

  .action-palette-item:hover,
  .action-palette-item.selected {
    background: var(--surface-selected);
  }

  /* Footer */
  .beam-footer {
    height: 38px;
    border-top: 1px solid var(--border-subtle);
    background: rgba(12, 14, 18, 0.95);
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 16px;
  }

  .key-pill {
    font-size: 12px;
    font-weight: 500;
    color: var(--text-secondary);
    display: flex;
    align-items: center;
    gap: 6px;
  }

  .footer-right {
    display: flex;
    gap: 14px;
    align-items: center;
  }

  .footer-hint {
    font-size: 11px;
    color: var(--text-muted);
    display: flex;
    align-items: center;
    gap: 4px;
  }
</style>
