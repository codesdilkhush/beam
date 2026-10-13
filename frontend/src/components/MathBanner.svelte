<script lang="ts">
  import type { CalcItem } from '../types/index.js';

  interface Props {
    mathResult: CalcItem | null;
    onCopy: (result: string) => void;
  }

  let { mathResult, onCopy }: Props = $props();
</script>

{#if mathResult && mathResult.success}
  <button
    type="button"
    class="math-banner"
    onclick={() => onCopy(mathResult?.result || '')}
  >
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
      <span class="copy-hint"><kbd>↵</kbd> Copy Result</span>
    </div>
  </button>
{/if}

<style>
  .math-banner {
    width: 100%;
    padding: 9px 18px;
    background: rgba(16, 185, 129, 0.08);
    border: none;
    border-bottom: 1px solid rgba(16, 185, 129, 0.2);
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

  .math-icon {
    color: var(--accent-emerald);
    display: flex;
    align-items: center;
  }

  .math-expr {
    color: var(--text-secondary);
    font-size: 14px;
  }

  .math-res {
    color: var(--accent-emerald);
    font-size: 17px;
    font-weight: 700;
  }

  .math-right {
    font-size: 11px;
    color: var(--text-muted);
  }

  .copy-hint kbd {
    background: rgba(255, 255, 255, 0.08);
    border: 1px solid rgba(255, 255, 255, 0.15);
    border-radius: 4px;
    padding: 2px 6px;
    font-size: 11px;
    color: var(--text-secondary);
  }
</style>
