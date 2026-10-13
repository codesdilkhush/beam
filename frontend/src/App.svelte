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

  import type {
    AppItem,
    ClipboardEntry,
    CalcItem,
    ActionOption,
    ThemeType,
    TabType
  } from './types/index.js';

  import { scrollSelectedIntoView } from './utils/scroll.js';

  import Header from './components/Header.svelte';
  import MathBanner from './components/MathBanner.svelte';
  import ResultsList from './components/ResultsList.svelte';
  import DetailPanel from './components/DetailPanel.svelte';
  import ActionPalette from './components/ActionPalette.svelte';
  import Footer from './components/Footer.svelte';

  let searchQuery = $state('');
  let activeTab = $state<TabType>('all');
  let selectedIndex = $state(0);
  let showActionPalette = $state(false);
  let actionPaletteIndex = $state(0);
  let currentTheme = $state<ThemeType>('default');

  let searchResults = $state<AppItem[]>([]);
  let clipboardItems = $state<ClipboardEntry[]>([]);
  let mathResult = $state<CalcItem | null>(null);

  let searchInputRef = $state<HTMLInputElement | null>(null);
  let resultsListRef = $state<HTMLDivElement | null>(null);
  let actionListRef = $state<HTMLDivElement | null>(null);

  $effect(() => {
    const idx = selectedIndex;
    if (resultsListRef) {
      requestAnimationFrame(() => {
        scrollSelectedIntoView(resultsListRef, '.list-item', idx);
      });
    }
  });

  $effect(() => {
    const idx = actionPaletteIndex;
    if (showActionPalette && actionListRef) {
      requestAnimationFrame(() => {
        scrollSelectedIntoView(actionListRef, '.palette-item', idx);
      });
    }
  });

  onMount(async () => {
    const savedTheme = localStorage.getItem('beam_theme');
    if (savedTheme && ['default', 'theme-catppuccin', 'theme-tokyo', 'theme-nord'].includes(savedTheme)) {
      currentTheme = savedTheme as ThemeType;
    }

    await performAppSearch('');
    await refreshClipboard();
    if (searchInputRef) {
      searchInputRef.focus();
    }
  });

  function setTheme(t: ThemeType) {
    currentTheme = t;
    localStorage.setItem('beam_theme', t);
  }

  async function performAppSearch(q: string) {
    try {
      const apps = await SearchApps(q);
      searchResults = apps || [];

      // Lazy load icons for first batch of results
      for (const app of searchResults.slice(0, 25)) {
        if (app.iconPath && !app.resolvedIcon) {
          LoadIconAsDataURL(app.iconPath).then((dataUrl: string) => {
            if (dataUrl) app.resolvedIcon = dataUrl;
          });
        }
      }
    } catch {
      searchResults = [];
    }
  }

  async function refreshClipboard() {
    try {
      const clips = await GetClipboard();
      clipboardItems = clips || [];
    } catch {
      clipboardItems = [];
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
        .then((res: CalcItem) => {
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

  const isClipboardMode = $derived(
    activeTab === 'clipboard' || searchQuery.startsWith('/clip')
  );

  const displayItems = $derived.by(() => {
    if (isClipboardMode) {
      const query = searchQuery.startsWith('/clip')
        ? searchQuery.slice(5).trim().toLowerCase()
        : searchQuery.trim().toLowerCase();

      if (!query) return clipboardItems;
      return clipboardItems.filter((item) =>
        item.content.toLowerCase().includes(query) || item.preview.toLowerCase().includes(query)
      );
    }

    return searchResults;
  });

  const selectedApp = $derived.by(() => {
    if (!isClipboardMode && displayItems.length > 0 && selectedIndex < displayItems.length) {
      return displayItems[selectedIndex] as AppItem;
    }
    return null;
  });

  const selectedClip = $derived.by(() => {
    if (isClipboardMode && displayItems.length > 0 && selectedIndex < displayItems.length) {
      return displayItems[selectedIndex] as ClipboardEntry;
    }
    return null;
  });

  const currentActions = $derived.by((): ActionOption[] => {
    const list: ActionOption[] = [];

    if (!isClipboardMode && selectedApp) {
      list.push({
        id: 'open',
        label: `Open ${selectedApp.name}`,
        shortcut: '↵',
        run: () => LaunchApp(selectedApp!.exec, selectedApp!.terminal, selectedApp!.id)
      });
      list.push({
        id: 'copy-exec',
        label: 'Copy Command Line',
        shortcut: 'C',
        run: () => CopyToClipboard(selectedApp!.exec)
      });
      list.push({
        id: 'copy-path',
        label: 'Copy Desktop File Path',
        shortcut: 'P',
        run: () => CopyToClipboard(selectedApp!.desktopPath)
      });
    }

    if (isClipboardMode && selectedClip) {
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
      } else if (e.key === 'PageDown') {
        e.preventDefault();
        actionPaletteIndex = Math.min(currentActions.length - 1, actionPaletteIndex + 5);
      } else if (e.key === 'PageUp') {
        e.preventDefault();
        actionPaletteIndex = Math.max(0, actionPaletteIndex - 5);
      } else if (e.key === 'Home') {
        e.preventDefault();
        actionPaletteIndex = 0;
      } else if (e.key === 'End') {
        e.preventDefault();
        if (currentActions.length > 0) actionPaletteIndex = currentActions.length - 1;
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
    } else if (e.key === 'PageDown') {
      e.preventDefault();
      selectedIndex = Math.min(displayItems.length - 1, selectedIndex + 5);
    } else if (e.key === 'PageUp') {
      e.preventDefault();
      selectedIndex = Math.max(0, selectedIndex - 5);
    } else if (e.key === 'Home') {
      e.preventDefault();
      selectedIndex = 0;
    } else if (e.key === 'End') {
      e.preventDefault();
      if (displayItems.length > 0) selectedIndex = displayItems.length - 1;
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
</script>

<svelte:window onkeydown={handleKeyDown} />

<div class="beam-window {currentTheme}">
  <Header
    {searchQuery}
    {isClipboardMode}
    {activeTab}
    bind:inputElement={searchInputRef}
    onQueryChange={(q) => (searchQuery = q)}
  />

  {#if !isClipboardMode}
    <MathBanner
      {mathResult}
      onCopy={(val) => CopyToClipboard(val)}
    />
  {/if}

  <div class="beam-body">
    <ResultsList
      items={displayItems}
      {selectedIndex}
      {isClipboardMode}
      {searchQuery}
      bind:listElement={resultsListRef}
      onSelect={(idx) => selectItem(idx)}
      onHover={(idx) => (selectedIndex = idx)}
    />

    <DetailPanel
      {isClipboardMode}
      {selectedClip}
      {selectedApp}
    />
  </div>

  <ActionPalette
    show={showActionPalette}
    actions={currentActions}
    selectedIndex={actionPaletteIndex}
    bind:listElement={actionListRef}
    onSelect={(action) => {
      action.run();
      showActionPalette = false;
    }}
    onHover={(idx) => (actionPaletteIndex = idx)}
    onClose={() => (showActionPalette = false)}
  />

  <Footer
    {isClipboardMode}
    {selectedApp}
  />
</div>

<style>
  .beam-window {
    width: 760px;
    height: 500px;
    background: var(--surface-canvas);
    backdrop-filter: blur(28px) saturate(180%);
    -webkit-backdrop-filter: blur(28px) saturate(180%);
    border: 1px solid var(--border-glass);
    border-radius: 14px;
    box-shadow: 0 24px 60px -12px rgba(0, 0, 0, 0.75), 0 0 0 1px rgba(255, 255, 255, 0.08);
    display: flex;
    flex-direction: column;
    overflow: hidden;
    position: relative;
  }

  .beam-body {
    flex: 1;
    display: flex;
    overflow: hidden;
  }
</style>
