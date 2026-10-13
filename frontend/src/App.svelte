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
    LoadIconAsDataURL,
    GetTodos,
    AddTodo,
    ToggleTodo,
    DeleteTodo,
    ClearCompletedTodos,
    GetEmojis,
    SearchEmojis,
    ConvertQuery
  } from '../wailsjs/go/main/App.js';

  import type {
    AppItem,
    ClipboardEntry,
    CalcItem,
    TodoEntry,
    EmojiEntry,
    ConvertItem,
    ActionOption,
    ThemeType,
    TabType
  } from './types/index.js';

  import { scrollSelectedIntoView } from './utils/scroll.js';

  import Header from './components/Header.svelte';
  import MathBanner from './components/MathBanner.svelte';
  import ResultsList from './components/ResultsList.svelte';
  import EmojiGrid from './components/EmojiGrid.svelte';
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
  let todoItems = $state<TodoEntry[]>([]);
  let emojiItems = $state<EmojiEntry[]>([]);
  let mathResult = $state<CalcItem | null>(null);
  let convertResult = $state<ConvertItem | null>(null);

  let searchInputRef = $state<HTMLInputElement | null>(null);
  let resultsListRef = $state<HTMLDivElement | null>(null);
  let actionListRef = $state<HTMLDivElement | null>(null);

  $effect(() => {
    const idx = selectedIndex;
    if (resultsListRef) {
      const selector = isEmojiMode ? '.emoji-tile' : '.list-item';
      requestAnimationFrame(() => {
        scrollSelectedIntoView(resultsListRef, selector, idx);
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
    await refreshTodos();
    await performEmojiSearch('');
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

  async function refreshTodos() {
    try {
      const todos = await GetTodos();
      todoItems = todos || [];
    } catch {
      todoItems = [];
    }
  }

  async function performEmojiSearch(q: string) {
    try {
      const emojis = q ? await SearchEmojis(q) : await GetEmojis();
      emojiItems = emojis || [];
    } catch {
      emojiItems = [];
    }
  }

  const isClipboardMode = $derived(
    activeTab === 'clipboard' || searchQuery.startsWith('/clip')
  );

  const isTodoMode = $derived(
    activeTab === 'todo' || searchQuery.startsWith('/todo') || searchQuery.startsWith('/t ')
  );

  const isEmojiMode = $derived(
    activeTab === 'emoji' || searchQuery.startsWith('/emoji') || searchQuery.startsWith('/e ') || searchQuery.startsWith(':')
  );

  function getCleanTodoQuery(query: string): string {
    if (query.startsWith('/todo ')) return query.slice(6).trim();
    if (query.startsWith('/todo')) return query.slice(5).trim();
    if (query.startsWith('/t ')) return query.slice(3).trim();
    if (activeTab === 'todo') return query.trim();
    return '';
  }

  function getCleanEmojiQuery(query: string): string {
    if (query.startsWith('/emoji ')) return query.slice(7).trim();
    if (query.startsWith('/emoji')) return query.slice(6).trim();
    if (query.startsWith('/e ')) return query.slice(3).trim();
    if (query.startsWith(':')) return query.slice(1).trim();
    if (activeTab === 'emoji') return query.trim();
    return '';
  }

  $effect(() => {
    const q = searchQuery;
    if (!isClipboardMode && !isTodoMode && !isEmojiMode) {
      performAppSearch(q);
    }
    selectedIndex = 0;
  });

  $effect(() => {
    if (isEmojiMode) {
      const q = getCleanEmojiQuery(searchQuery);
      performEmojiSearch(q);
    }
  });

  $effect(() => {
    if (isTodoMode) {
      refreshTodos();
    }
  });

  $effect(() => {
    if (isClipboardMode) {
      refreshClipboard();
    }
  });

  $effect(() => {
    const q = searchQuery.trim();
    if (isClipboardMode || isTodoMode || isEmojiMode || q.length < 2) {
      convertResult = null;
      mathResult = null;
      return;
    }

    let active = true;

    ConvertQuery(q)
      .then((res: ConvertItem) => {
        if (!active) return;
        if (res && res.success) {
          convertResult = res;
          mathResult = null;
          return;
        }
        convertResult = null;

        if (/[0-9\+\-\*\/\%\^\(\)]/.test(q)) {
          EvaluateMath(q)
            .then((mRes: CalcItem) => {
              if (!active) return;
              if (mRes && mRes.success) {
                mathResult = mRes;
              } else {
                mathResult = null;
              }
            })
            .catch(() => {
              if (active) mathResult = null;
            });
        } else {
          mathResult = null;
        }
      })
      .catch(() => {
        if (!active) return;
        convertResult = null;
        if (/[0-9\+\-\*\/\%\^\(\)]/.test(q)) {
          EvaluateMath(q)
            .then((mRes: CalcItem) => {
              if (!active) return;
              if (mRes && mRes.success) {
                mathResult = mRes;
              } else {
                mathResult = null;
              }
            })
            .catch(() => {
              if (active) mathResult = null;
            });
        } else {
          mathResult = null;
        }
      });

    return () => {
      active = false;
    };
  });

  const displayItems = $derived.by((): (AppItem | ClipboardEntry | TodoEntry)[] => {
    if (isTodoMode) {
      const q = getCleanTodoQuery(searchQuery).toLowerCase();
      if (!q) return todoItems;
      return todoItems.filter((t) => t.title.toLowerCase().includes(q));
    }

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
    if (!isClipboardMode && !isTodoMode && !isEmojiMode && displayItems.length > 0 && selectedIndex < displayItems.length) {
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

  const selectedTodo = $derived.by(() => {
    if (isTodoMode && displayItems.length > 0 && selectedIndex < displayItems.length) {
      return displayItems[selectedIndex] as TodoEntry;
    }
    return null;
  });

  const selectedEmoji = $derived.by(() => {
    if (isEmojiMode && emojiItems.length > 0 && selectedIndex < emojiItems.length) {
      return emojiItems[selectedIndex] as EmojiEntry;
    }
    return null;
  });

  const currentActions = $derived.by((): ActionOption[] => {
    const list: ActionOption[] = [];

    if (isTodoMode) {
      const cleanTitle = getCleanTodoQuery(searchQuery);
      if (cleanTitle) {
        list.push({
          id: 'add-todo',
          label: `Add Task "${cleanTitle}"`,
          shortcut: '↵',
          run: async () => {
            await AddTodo(cleanTitle);
            if (activeTab === 'todo') searchQuery = '';
            await refreshTodos();
          }
        });
      }
      if (selectedTodo) {
        list.push({
          id: 'toggle-todo',
          label: selectedTodo.completed ? 'Mark Task Incomplete' : 'Mark Task Completed',
          shortcut: 'Space',
          run: async () => {
            await ToggleTodo(selectedTodo!.id);
            await refreshTodos();
          }
        });
        list.push({
          id: 'delete-todo',
          label: 'Delete Task',
          shortcut: 'Ctrl+D',
          run: async () => {
            await DeleteTodo(selectedTodo!.id);
            await refreshTodos();
          }
        });
      }
      const hasCompleted = todoItems.some((t) => t.completed);
      if (hasCompleted) {
        list.push({
          id: 'clear-completed',
          label: 'Clear Completed Tasks',
          shortcut: 'Shift+Del',
          run: async () => {
            await ClearCompletedTodos();
            await refreshTodos();
          }
        });
      }
    } else if (isEmojiMode && selectedEmoji) {
      list.push({
        id: 'copy-emoji',
        label: `Copy Emoji "${selectedEmoji.emoji}"`,
        shortcut: '↵',
        run: () => CopyToClipboard(selectedEmoji!.emoji)
      });
      list.push({
        id: 'copy-shortcode',
        label: `Copy Shortcode "${selectedEmoji.shortcode}"`,
        shortcut: 'C',
        run: () => CopyToClipboard(selectedEmoji!.shortcode)
      });
      list.push({
        id: 'copy-codepoint',
        label: 'Copy Codepoint',
        shortcut: 'P',
        run: () => CopyToClipboard(selectedEmoji!.emoji.codePointAt(0)?.toString(16).toUpperCase() || '')
      });
    } else if (isClipboardMode && selectedClip) {
      list.push({
        id: 'copy',
        label: 'Copy to Clipboard',
        shortcut: '↵',
        run: () => CopyToClipboard(selectedClip!.content)
      });
    } else if (!isClipboardMode && !isTodoMode && !isEmojiMode && selectedApp) {
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

    if (isTodoMode && selectedTodo && (e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 'd') {
      e.preventDefault();
      DeleteTodo(selectedTodo.id).then(() => {
        refreshTodos();
        if (selectedIndex >= displayItems.length - 1 && selectedIndex > 0) {
          selectedIndex--;
        }
      });
      return;
    }

    if (
      isTodoMode &&
      selectedTodo &&
      ((e.ctrlKey && e.code === 'Space') ||
        (e.key === ' ' && (searchQuery === '' || searchQuery === '/todo' || searchQuery === '/t')))
    ) {
      e.preventDefault();
      ToggleTodo(selectedTodo.id).then(() => refreshTodos());
      return;
    }

    if (
      isTodoMode &&
      selectedTodo &&
      e.key === 'Delete' &&
      (searchQuery === '' || searchQuery === '/todo' || searchQuery === '/t')
    ) {
      e.preventDefault();
      DeleteTodo(selectedTodo.id).then(() => {
        refreshTodos();
        if (selectedIndex >= displayItems.length - 1 && selectedIndex > 0) {
          selectedIndex--;
        }
      });
      return;
    }

    if (isEmojiMode) {
      const COLS = 8;
      if (e.key === 'ArrowRight') {
        e.preventDefault();
        if (selectedIndex < displayItems.length - 1) selectedIndex++;
        return;
      } else if (e.key === 'ArrowLeft') {
        e.preventDefault();
        if (selectedIndex > 0) selectedIndex--;
        return;
      } else if (e.key === 'ArrowDown') {
        e.preventDefault();
        if (selectedIndex + COLS < displayItems.length) {
          selectedIndex += COLS;
        } else {
          selectedIndex = displayItems.length - 1;
        }
        return;
      } else if (e.key === 'ArrowUp') {
        e.preventDefault();
        if (selectedIndex - COLS >= 0) {
          selectedIndex -= COLS;
        } else {
          selectedIndex = 0;
        }
        return;
      } else if (e.key === 'PageDown') {
        e.preventDefault();
        selectedIndex = Math.min(displayItems.length - 1, selectedIndex + COLS * 3);
        return;
      } else if (e.key === 'PageUp') {
        e.preventDefault();
        selectedIndex = Math.max(0, selectedIndex - COLS * 3);
        return;
      } else if (e.key === 'Home') {
        e.preventDefault();
        selectedIndex = 0;
        return;
      } else if (e.key === 'End') {
        e.preventDefault();
        if (displayItems.length > 0) selectedIndex = displayItems.length - 1;
        return;
      }
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
      else if (activeTab === 'clipboard') activeTab = 'todo';
      else if (activeTab === 'todo') activeTab = 'emoji';
      else if (activeTab === 'emoji') activeTab = 'apps';
      else activeTab = 'all';

      selectedIndex = 0;
      if (activeTab === 'clipboard') refreshClipboard();
      else if (activeTab === 'todo') refreshTodos();
      else if (activeTab === 'emoji') performEmojiSearch('');
    }
  }

  async function executeCurrentSelection() {
    if (convertResult && convertResult.success && !isClipboardMode && !isTodoMode && !isEmojiMode) {
      await CopyToClipboard(convertResult.result);
      return;
    }

    if (mathResult && mathResult.success && selectedIndex === 0 && !isClipboardMode && !isTodoMode && !isEmojiMode) {
      await CopyToClipboard(mathResult.result);
      return;
    }

    if (isTodoMode) {
      const cleanTitle = getCleanTodoQuery(searchQuery);
      if (displayItems.length === 0 && cleanTitle) {
        await AddTodo(cleanTitle);
        if (activeTab === 'todo') {
          searchQuery = '';
        }
        await refreshTodos();
      } else if (selectedTodo) {
        await ToggleTodo(selectedTodo.id);
        await refreshTodos();
      }
      return;
    }

    if (isEmojiMode) {
      if (selectedEmoji) {
        await CopyToClipboard(selectedEmoji.emoji);
      }
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
    {isTodoMode}
    {isEmojiMode}
    {activeTab}
    bind:inputElement={searchInputRef}
    onQueryChange={(q) => (searchQuery = q)}
  />

  {#if !isClipboardMode && !isTodoMode && !isEmojiMode}
    <MathBanner
      {mathResult}
      {convertResult}
      onCopy={(val) => CopyToClipboard(val)}
    />
  {/if}

  <div class="beam-body">
    {#if isEmojiMode}
      <EmojiGrid
        items={emojiItems}
        {selectedIndex}
        {searchQuery}
        bind:gridElement={resultsListRef}
        onSelect={(idx) => selectItem(idx)}
        onHover={(idx) => (selectedIndex = idx)}
      />
    {:else}
      <ResultsList
        items={displayItems}
        {selectedIndex}
        {isClipboardMode}
        {isTodoMode}
        {searchQuery}
        bind:listElement={resultsListRef}
        onSelect={(idx) => selectItem(idx)}
        onHover={(idx) => (selectedIndex = idx)}
      />
    {/if}

    {#if !isTodoMode && !isEmojiMode}
      <DetailPanel
        {isClipboardMode}
        {isEmojiMode}
        {selectedClip}
        {selectedEmoji}
        {selectedApp}
      />
    {/if}
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
    {isTodoMode}
    {isEmojiMode}
    {selectedApp}
    {selectedTodo}
    {selectedEmoji}
    {convertResult}
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
