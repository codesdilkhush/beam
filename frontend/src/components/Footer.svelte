<script lang="ts">
  import type { AppItem } from '../types/index.js';

  interface Props {
    isClipboardMode: boolean;
    selectedApp: AppItem | null;
  }

  let { isClipboardMode, selectedApp }: Props = $props();
</script>

<footer class="footer">
  <div class="footer-left">
    {#if isClipboardMode}
      <span class="active-action">
        <kbd>↵</kbd> Copy to Clipboard
      </span>
    {:else if selectedApp}
      <span class="active-action">
        <kbd>↵</kbd> Open Application
      </span>
    {:else}
      <span class="active-action">
        <kbd>↵</kbd> Search Online
      </span>
    {/if}
  </div>
  <div class="footer-right">
    <span class="footer-hint"><kbd>Ctrl+K</kbd> Actions</span>
    <span class="footer-hint"><kbd>Tab</kbd> Mode</span>
    <span class="footer-hint"><kbd>↑</kbd><kbd>↓</kbd> Navigate</span>
    <span class="footer-hint"><kbd>Esc</kbd> Close</span>
  </div>
</footer>

<style>
  .footer {
    height: 38px;
    border-top: 1px solid var(--border-subtle);
    background: var(--surface-header);
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0 16px;
    font-size: 11.5px;
    color: var(--text-muted);
  }

  .footer-left {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .active-action {
    display: flex;
    align-items: center;
    gap: 6px;
    color: var(--text-secondary);
    font-weight: 500;
  }

  .footer-right {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .footer-hint {
    display: flex;
    align-items: center;
    gap: 4px;
  }

  kbd {
    background: rgba(255, 255, 255, 0.08);
    border: 1px solid rgba(255, 255, 255, 0.14);
    border-radius: 4px;
    padding: 1px 5px;
    font-size: 10.5px;
    color: var(--text-primary);
  }
</style>
